import App from "./App.tsx";
import { tsr } from "./api";
import "./index.css";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactNode } from "react";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { Toaster } from "sonner";

const queryClient = new QueryClient();

function RootLayout(props: { children: ReactNode }) {
  return (
    <QueryClientProvider client={queryClient}>
      <tsr.ReactQueryProvider>{props.children}</tsr.ReactQueryProvider>
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
