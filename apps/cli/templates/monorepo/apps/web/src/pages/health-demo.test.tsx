import { HealthDemoPage } from "@/pages/health-demo";
import { renderWithRouter } from "@/testing";
import { screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

const routes = [{ path: "/health-demo", element: <HealthDemoPage /> }];

describe("HealthDemoPage", () => {
  it("renders the health response", async () => {
    renderWithRouter(routes, { initialEntries: ["/health-demo"] });

    expect(await screen.findByText(/"status": "healthy"/)).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Refresh" }))
      .toBeInTheDocument();
  });
});
