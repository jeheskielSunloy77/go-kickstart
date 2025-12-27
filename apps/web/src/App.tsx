import "./index.css";
import { router } from "./router";
import { tsr } from "@/api";
import { ThemeProvider } from "@/hooks/use-theme.tsx";
import "@go-kickstart/ui/styles.css";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactNode } from "react";
import { RouterProvider } from "react-router-dom";
import { Toaster } from "sonner";

const queryClient = new QueryClient();

function Providers(props: { children: ReactNode }) {
  return (
    <QueryClientProvider client={queryClient}>
      <tsr.ReactQueryProvider>
        <ThemeProvider>{props.children}</ThemeProvider>
      </tsr.ReactQueryProvider>
      <Toaster position="top-right" />
    </QueryClientProvider>
  );
}

export function App() {
  return (
    <Providers>
      <RouterProvider router={router} />
    </Providers>
  );
}
