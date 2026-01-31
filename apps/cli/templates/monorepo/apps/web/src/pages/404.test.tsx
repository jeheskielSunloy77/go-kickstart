import { NotFound } from "@/pages/404";
import { renderWithRouter } from "@/testing";
import { screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

const routes = [{ path: "*", element: <NotFound /> }];

describe("NotFound", () => {
  it("renders a 404 message", async () => {
    renderWithRouter(routes, { initialEntries: ["/missing"] });

    expect(await screen.findByText("404")).toBeInTheDocument();
    expect(screen.getByText("Page not found")).toBeInTheDocument();
  });
});
