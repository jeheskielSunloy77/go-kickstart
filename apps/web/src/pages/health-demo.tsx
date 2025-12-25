import { tsr } from "@/api";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Spinner } from "@/components/ui/spinner";
import { getApiErrorMessage, isUnauthorizedError } from "@/lib/api-errors";
import { Link } from "react-router-dom";

export function HealthDemoPage() {
  const healthQuery = tsr.health.getHealth.useQuery({
    queryKey: ["health"],
    select: (response) => response.body,
  });

  const sessionQuery = tsr.auth.me.useQuery({
    queryKey: ["auth", "me"],
    retry: false,
    select: (response) => response.body,
  });

  const isAuthed = sessionQuery.isSuccess;
  const showLogin =
    sessionQuery.isError && isUnauthorizedError(sessionQuery.error);

  return (
    <main className="min-h-screen bg-gradient-to-b from-background via-muted/40 to-background">
      <div className="mx-auto flex min-h-screen w-full max-w-4xl items-center px-6 py-12">
        <Card className="w-full">
          <CardHeader className="space-y-3">
            <div>
              <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
                Example fetch
              </p>
              <CardTitle className="font-serif text-3xl">
                Health check demo
              </CardTitle>
            </div>
            <CardDescription>
              This page calls the backend health endpoint and shows the
              response.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="rounded-xl border border-border bg-background p-4">
              {healthQuery.isLoading && (
                <div className="flex min-h-40 items-center justify-center">
                  <Spinner className="size-8" />
                </div>
              )}
              {healthQuery.isError && (
                <div className="text-sm text-destructive">
                  {getApiErrorMessage(healthQuery.error)}
                </div>
              )}
              {healthQuery.isSuccess && (
                <pre className="text-xs text-foreground whitespace-pre-wrap">
                  {JSON.stringify(healthQuery.data.data, null, 2)}
                </pre>
              )}
            </div>
          </CardContent>
          <CardFooter className="flex flex-col gap-3 border-t sm:flex-row sm:justify-between">
            <div className="flex w-full flex-col gap-2 sm:w-auto sm:flex-row">
              <Button asChild variant="outline" className="w-full sm:w-auto">
                <Link to="/">Back to home</Link>
              </Button>
              {isAuthed && (
                <Button asChild className="w-full sm:w-auto">
                  <Link to="/auth/me">Go to profile</Link>
                </Button>
              )}
              {showLogin && (
                <Button asChild className="w-full sm:w-auto">
                  <Link to="/auth/login">Sign in</Link>
                </Button>
              )}
            </div>
            <Button
              onClick={() => healthQuery.refetch()}
              disabled={healthQuery.isLoading}
              className="w-full sm:w-auto"
            >
              Refresh
            </Button>
          </CardFooter>
        </Card>
      </div>
    </main>
  );
}
