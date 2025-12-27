import App from "./App.tsx";
import { tsr } from "./api";
import { ThemeProvider } from "./hooks/use-theme.tsx";
import "./index.css";
import "@go-kickstart/ui/styles.css";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactNode } from "react";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { Toaster } from "sonner";

const queryClient = new QueryClient();

function RootLayout(props: { children: ReactNode }) {
  return (
    <QueryClientProvider client={queryClient}>
      <tsr.ReactQueryProvider>
        <ThemeProvider>{props.children}</ThemeProvider>
      </tsr.ReactQueryProvider>
      <Toaster position="top-right" />
    </QueryClientProvider>
  );
}

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <RootLayout>
      <App />
    </RootLayout>
  </StrictMode>,
);
