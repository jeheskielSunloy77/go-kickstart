import { RootErrorBoundary } from "@/pages/error";
import { render, screen } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

const routerMocks = vi.hoisted(() => ({
  useRouteError: vi.fn(),
  isRouteErrorResponse: vi.fn(),
  navigate: vi.fn(),
}));

vi.mock("react-router-dom", async () => {
  const actual = await vi.importActual<typeof import("react-router-dom")>(
    "react-router-dom",
  );
  return {
    ...actual,
    useRouteError: () => routerMocks.useRouteError(),
    isRouteErrorResponse: (error: unknown) =>
      routerMocks.isRouteErrorResponse(error),
    useNavigate: () => routerMocks.navigate,
  };
});

describe("RootErrorBoundary", () => {
  beforeEach(() => {
    routerMocks.useRouteError.mockReset();
    routerMocks.isRouteErrorResponse.mockReset();
    routerMocks.navigate.mockReset();
  });

  it("renders route errors with status details", () => {
    routerMocks.useRouteError.mockReturnValue({
      status: 404,
      statusText: "Not Found",
      data: "Missing page",
    });
    routerMocks.isRouteErrorResponse.mockReturnValue(true);

    render(<RootErrorBoundary />);

    expect(screen.getByText("404")).toBeInTheDocument();
    expect(screen.getByText("Not Found")).toBeInTheDocument();
    expect(screen.getByText("Missing page")).toBeInTheDocument();
  });

  it("renders unexpected JS errors", () => {
    routerMocks.useRouteError.mockReturnValue(new Error("Boom"));
    routerMocks.isRouteErrorResponse.mockReturnValue(false);

    render(<RootErrorBoundary />);

    expect(screen.getByText("Something went wrong")).toBeInTheDocument();
    expect(screen.getByText("Boom")).toBeInTheDocument();
  });

  it("renders a fallback for unknown errors", () => {
    routerMocks.useRouteError.mockReturnValue("unknown");
    routerMocks.isRouteErrorResponse.mockReturnValue(false);

    render(<RootErrorBoundary />);

    expect(screen.getByText("Unknown error")).toBeInTheDocument();
  });
});
