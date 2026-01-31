import { RequireAuth } from "@/auth/require-auth";
import { makeUser, renderWithRouter } from "@/testing";
import { http, HttpResponse } from "msw";
import { server } from "@/testing/server";
import { screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

const apiBase =
  (import.meta.env.VITE_API_URL as string | undefined) ||
  "http://localhost:3001";

const routes = [
  {
    path: "/",
    element: (
      <RequireAuth>
        <div>Protected content</div>
      </RequireAuth>
    ),
  },
  { path: "/auth/login", element: <div>Login page</div> },
];

describe("RequireAuth", () => {
  it("renders protected content for authenticated users", async () => {
    server.use(
      http.get(`${apiBase}/api/v1/auth/me`, () => {
        return HttpResponse.json(makeUser());
      }),
    );

    renderWithRouter(routes, { initialEntries: ["/"] });

    expect(await screen.findByText("Protected content")).toBeInTheDocument();
  });

  it("redirects to login on unauthorized responses", async () => {
    server.use(
      http.get(`${apiBase}/api/v1/auth/me`, () => {
        return HttpResponse.json(
          { message: "Unauthorized" },
          { status: 401 },
        );
      }),
    );

    renderWithRouter(routes, { initialEntries: ["/"] });

    expect(await screen.findByText("Login page")).toBeInTheDocument();
  });

  it("shows an error message on non-auth errors", async () => {
    server.use(
      http.get(`${apiBase}/api/v1/auth/me`, () => {
        return HttpResponse.json(
          { message: "Something broke" },
          { status: 500 },
        );
      }),
    );

    renderWithRouter(routes, { initialEntries: ["/"] });

    expect(await screen.findByText("Something broke")).toBeInTheDocument();
  });
});
