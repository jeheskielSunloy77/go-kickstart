import { ForgotPasswordPage } from "@/pages/auth/forgot-password";
import { renderWithRouter } from "@/testing";
import { screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it } from "vitest";

const routes = [{ path: "/auth/forgot-password", element: <ForgotPasswordPage /> }];

describe("ForgotPasswordPage", () => {
  it("validates the email input", async () => {
    renderWithRouter(routes, { initialEntries: ["/auth/forgot-password"] });

    const user = userEvent.setup();
    await user.type(screen.getByLabelText("Email"), "bad");
    await user.click(
      screen.getByRole("button", { name: "Send reset link" }),
    );

    expect(await screen.findByText(/invalid input/i)).toBeInTheDocument();
  });
});
