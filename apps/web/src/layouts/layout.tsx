import { tsr } from "../api";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactNode } from "react";

const queryClient = new QueryClient();

export function Layout(props: { children: ReactNode }) {
  // add other layout components here with its conditions if needed. Example:
  // if (isAdmin) reeturn <RootLayout><AdminLayout>{props.children}</AdminLayout></RootLayout>

  return <RootLayout>{props.children}</RootLayout>;
}

function RootLayout(props: { children: ReactNode }) {
  return (
    <QueryClientProvider client={queryClient}>
      <tsr.ReactQueryProvider>{props.children}</tsr.ReactQueryProvider>
    </QueryClientProvider>
  );
}
