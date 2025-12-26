import { VerifyEmailPage } from "@/pages/auth/verify-email";
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
  { path: "/auth/verify-email", element: <VerifyEmailPage /> },
  { path: "/auth/me", element: <div>Me page</div> },
];

describe("VerifyEmailPage", () => {
  it("prefills the email from the query string", async () => {
    renderWithRouter(routes, {
      initialEntries: ["/auth/verify-email?email=test%40example.com"],
    });

    const emailInput = await screen.findByLabelText("Email");
    expect(emailInput).toHaveValue("test@example.com");
  });

  it("navigates to /auth/me after successful verification", async () => {
    server.use(
      http.post(`${apiBase}/api/v1/auth/verify-email`, () => {
        return HttpResponse.json(makeUser());
      }),
    );

    const { router } = renderWithRouter(routes, {
      initialEntries: ["/auth/verify-email?email=verifier%40example.com"],
    });

    const user = userEvent.setup();
    await user.type(screen.getByLabelText("Verification code"), "123456");
    await user.click(screen.getByRole("button", { name: "Verify email" }));

    await waitFor(() => {
      expect(router.state.location.pathname).toBe("/auth/me");
    });
  });
});
