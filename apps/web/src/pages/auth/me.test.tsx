import { AuthProvider } from "@/auth/auth-context";
import { MePage } from "@/pages/auth/me";
import { makeUser, renderWithRouter } from "@/testing";
import { screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

const routes = [
  {
    path: "/auth/me",
    element: (
      <AuthProvider user={makeUser({ emailVerifiedAt: undefined })}>
        <MePage />
      </AuthProvider>
    ),
  },
];

describe("MePage", () => {
  it("renders user profile details", async () => {
    renderWithRouter(routes, { initialEntries: ["/auth/me"] });

    expect(await screen.findByText("Welcome, janedoe")).toBeInTheDocument();
    expect(screen.getByText("jane.doe@example.com")).toBeInTheDocument();
    expect(screen.getByText("Not verified")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Resend code" }))
      .toBeInTheDocument();
  });

  it("hides the resend card when email is verified", () => {
    const verifiedRoutes = [
      {
        path: "/auth/me",
        element: (
          <AuthProvider user={makeUser({ emailVerifiedAt: "2024-01-01T00:00:00.000Z" })}>
            <MePage />
          </AuthProvider>
        ),
      },
    ];

    renderWithRouter(verifiedRoutes, { initialEntries: ["/auth/me"] });

    expect(screen.queryByText("Resend code")).not.toBeInTheDocument();
  });
});
