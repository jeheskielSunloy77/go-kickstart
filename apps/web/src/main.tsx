import App from "./App.tsx";
import { tsr } from "./api";
import "./index.css";
import { GOOGLE_CLIENT_ID } from "@/config/env";
import { GoogleOAuthProvider } from "@react-oauth/google";
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
      <Toaster position="top-right" richColors />
    </QueryClientProvider>
  );
}

const app = (
  <RootLayout>
    <App />
  </RootLayout>
);

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    {GOOGLE_CLIENT_ID ? (
      <GoogleOAuthProvider clientId={GOOGLE_CLIENT_ID}>
        {app}
      </GoogleOAuthProvider>
    ) : (
      app
    )}
  </StrictMode>,
);
