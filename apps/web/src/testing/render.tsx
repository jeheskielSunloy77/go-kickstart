import { tsr } from "@/api";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, type RenderOptions } from "@testing-library/react";
import type { ReactElement } from "react";
import { Toaster } from "sonner";
import {
  RouterProvider,
  createMemoryRouter,
  type RouteObject,
} from "react-router-dom";

type RenderWithRouterOptions = {
  initialEntries?: string[];
};

function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
      mutations: {
        retry: false,
      },
    },
  });
}

export function renderWithProviders(
  ui: ReactElement,
  options?: Omit<RenderOptions, "wrapper">,
) {
  const queryClient = createTestQueryClient();

  const result = render(
    <QueryClientProvider client={queryClient}>
      <tsr.ReactQueryProvider>{ui}</tsr.ReactQueryProvider>
      <Toaster position="top-right" />
    </QueryClientProvider>,
    options,
  );

  return { ...result, queryClient };
}

export function renderWithRouter(
  routes: RouteObject[],
  options: RenderWithRouterOptions = {},
) {
  const queryClient = createTestQueryClient();
  const router = createMemoryRouter(routes, {
    initialEntries: options.initialEntries ?? ["/"],
  });

  const result = render(
    <QueryClientProvider client={queryClient}>
      <tsr.ReactQueryProvider>
        <RouterProvider router={router} />
      </tsr.ReactQueryProvider>
      <Toaster position="top-right" />
    </QueryClientProvider>,
  );

  return { ...result, router, queryClient };
}
