import { LoginPage } from "@/pages/auth/login";
import { makeUser, renderWithRouter } from "@/testing";
import { http, HttpResponse } from "msw";
import { server } from "@/testing/server";
import { screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it } from "vitest";

const apiBase =
  (import.meta.env.VITE_API_URL as string | undefined) ||
  "http://localhost:3001";

const routes = [
  { path: "/auth/login", element: <LoginPage /> },
  { path: "/auth/me", element: <div>Me page</div> },
  { path: "/auth/verify-email", element: <div>Verify page</div> },
];

describe("LoginPage", () => {
  it("shows validation errors on empty submit", async () => {
    renderWithRouter(routes, { initialEntries: ["/auth/login"] });

    const user = userEvent.setup();
    await user.click(screen.getByRole("button", { name: "Sign in" }));

    expect(
      await screen.findByText("Identifier is required"),
    ).toBeInTheDocument();
    expect(
      await screen.findByText("Password is required"),
    ).toBeInTheDocument();
  });

  it("navigates to the profile when email is verified", async () => {
    server.use(
      http.post(`${apiBase}/api/v1/auth/login`, () => {
        return HttpResponse.json(makeUser({ emailVerifiedAt: "2024-01-01T00:00:00.000Z" }));
      }),
    );

    const { router } = renderWithRouter(routes, {
      initialEntries: ["/auth/login"],
    });

    const user = userEvent.setup();
    await user.type(
      screen.getByLabelText("Email or username"),
      "jane",
    );
    await user.type(screen.getByLabelText("Password"), "password123");
    await user.click(screen.getByRole("button", { name: "Sign in" }));

    await waitFor(() => {
      expect(router.state.location.pathname).toBe("/auth/me");
    });
  });

  it("navigates to verify-email when the account is unverified", async () => {
    server.use(
      http.post(`${apiBase}/api/v1/auth/login`, () => {
        return HttpResponse.json(
          makeUser({ emailVerifiedAt: undefined, email: "test@example.com" }),
        );
      }),
    );

    const { router } = renderWithRouter(routes, {
      initialEntries: ["/auth/login"],
    });

    const user = userEvent.setup();
    await user.type(
      screen.getByLabelText("Email or username"),
      "test@example.com",
    );
    await user.type(screen.getByLabelText("Password"), "password123");
    await user.click(screen.getByRole("button", { name: "Sign in" }));

    await waitFor(() => {
      expect(router.state.location.pathname).toBe("/auth/verify-email");
    });
    expect(router.state.location.search).toContain("email=test%40example.com");
  });
});
