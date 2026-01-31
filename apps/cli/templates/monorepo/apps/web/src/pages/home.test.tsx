import { HomePage } from "@/pages/home";
import { makeUser, renderWithRouter } from "@/testing";
import { http, HttpResponse } from "msw";
import { server } from "@/testing/server";
import { screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

const apiBase =
  (import.meta.env.VITE_API_URL as string | undefined) ||
  "http://localhost:3001";

const routes = [{ path: "/", element: <HomePage /> }];

describe("HomePage", () => {
  it("shows the profile link when authenticated", async () => {
    server.use(
      http.get(`${apiBase}/api/v1/auth/me`, () => {
        return HttpResponse.json(makeUser());
      }),
    );

    renderWithRouter(routes, { initialEntries: ["/"] });

    expect(await screen.findByText("Go to profile")).toBeInTheDocument();
  });

  it("shows the sign in CTA when unauthorized", async () => {
    server.use(
      http.get(`${apiBase}/api/v1/auth/me`, () => {
        return HttpResponse.json({ message: "Unauthorized" }, { status: 401 });
      }),
    );

    renderWithRouter(routes, { initialEntries: ["/"] });

    expect(await screen.findByText("Sign in")).toBeInTheDocument();
  });
});
