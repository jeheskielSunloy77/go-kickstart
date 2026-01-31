import { tsr } from "@/api";
import { ThemeProvider } from "@/hooks/use-theme";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, type RenderOptions } from "@testing-library/react";
import type { ReactElement } from "react";
import {
  RouterProvider,
  createMemoryRouter,
  type RouteObject,
} from "react-router-dom";
import { Toaster } from "sonner";

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

  const result = render(<Providers>{ui}</Providers>, options);

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
    <Providers>
      <RouterProvider router={router} />
    </Providers>,
  );

  return { ...result, router, queryClient };
}

export function renderWithCriticalProviders(
  ui: ReactElement,
  options?: Omit<RenderOptions, "wrapper">,
) {
  const result = render(<CriticalProviders>{ui}</CriticalProviders>, options);

  return { ...result };
}

function Providers(props: { children: React.ReactNode }) {
  const queryClient = createTestQueryClient();
  return (
    <QueryClientProvider client={queryClient}>
      <tsr.ReactQueryProvider>
        <CriticalProviders>{props.children}</CriticalProviders>
      </tsr.ReactQueryProvider>
      <Toaster position="top-right" />
    </QueryClientProvider>
  );
}

function CriticalProviders(props: { children: React.ReactNode }) {
  return <ThemeProvider>{props.children}</ThemeProvider>;
}
