import { RegisterPage } from "@/pages/auth/register";
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
  { path: "/auth/register", element: <RegisterPage /> },
  { path: "/auth/verify-email", element: <div>Verify page</div> },
];

describe("RegisterPage", () => {
  it("validates inputs before submitting", async () => {
    renderWithRouter(routes, { initialEntries: ["/auth/register"] });

    const user = userEvent.setup();
    await user.type(screen.getByLabelText("Email"), "bad");
    await user.type(screen.getByLabelText("Username"), "ab");
    await user.type(screen.getByLabelText("Password"), "short");
    await user.click(
      screen.getByRole("button", { name: "Create account" }),
    );

    expect(await screen.findByText("Enter a valid email")).toBeInTheDocument();
    expect(await screen.findByText("Username is too short")).toBeInTheDocument();
    expect(
      await screen.findByText("Password must be at least 8 characters"),
    ).toBeInTheDocument();
  });

  it("navigates to verify-email on success", async () => {
    server.use(
      http.post(`${apiBase}/api/v1/auth/register`, () => {
        return HttpResponse.json(makeUser({ email: "new@example.com" }));
      }),
    );

    const { router } = renderWithRouter(routes, {
      initialEntries: ["/auth/register"],
    });

    const user = userEvent.setup();
    await user.type(
      screen.getByLabelText("Email"),
      "new@example.com",
    );
    await user.type(screen.getByLabelText("Username"), "newuser");
    await user.type(
      screen.getByLabelText("Password"),
      "password123",
    );
    await user.click(
      screen.getByRole("button", { name: "Create account" }),
    );

    await waitFor(() => {
      expect(router.state.location.pathname).toBe("/auth/verify-email");
    });
    expect(router.state.location.search).toContain("email=new%40example.com");
  });
});
