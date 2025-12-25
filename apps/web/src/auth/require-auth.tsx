import { tsr } from "@/api";
import { Spinner } from "@/components/ui/spinner";
import { getApiErrorMessage, isUnauthorizedError } from "@/lib/api-errors";
import type { ReactNode } from "react";
import { Navigate } from "react-router-dom";

export function RequireAuth(props: { children: ReactNode }) {
  const meQuery = tsr.auth.me.useQuery({
    queryKey: ["auth", "me"],
    retry: false,
    select: (response) => response.body,
  });

  if (meQuery.isLoading) {
    return (
      <div className="flex min-h-60 items-center justify-center">
        <Spinner className="size-8" />
      </div>
    );
  }

  if (meQuery.isError) {
    if (isUnauthorizedError(meQuery.error)) {
      return <Navigate to="/auth/login" replace />;
    }

    return (
      <div className="rounded-2xl border border-destructive/20 bg-destructive/10 p-4 text-sm text-destructive">
        {getApiErrorMessage(meQuery.error)}
      </div>
    );
  }

  return <>{props.children}</>;
}
