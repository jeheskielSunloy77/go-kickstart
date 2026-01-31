import { tsr } from "@/api";
import {
  Button,
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
  Spinner,
} from "@go-kickstart/ui";
import { isUnauthorizedError } from "@/lib/api-errors";
import { Link } from "react-router-dom";

export function HomePage() {
  const sessionQuery = tsr.auth.me.useQuery({
    queryKey: ["auth", "me"],
    retry: false,
    select: (response) => response.body,
  });

  const isAuthed = sessionQuery.isSuccess;
  const isChecking = sessionQuery.isLoading;
  const showLogin =
    sessionQuery.isError && isUnauthorizedError(sessionQuery.error);

  return (
    <main className="min-h-screen bg-gradient-to-b from-background via-muted/40 to-background">
      <div className="mx-auto flex min-h-screen w-full max-w-4xl items-center px-6 py-12">
        <Card className="w-full">
          <CardHeader className="space-y-3">
            <div>
              <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
                Go Kickstart
              </p>
              <CardTitle className="font-serif text-3xl">
                A Go + React starter template
              </CardTitle>
            </div>
            <CardDescription>
              A monorepo for a Go API and a Vite + React web app with shared
              TypeScript packages, managed with Turborepo and Bun workspaces.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4 text-sm text-muted-foreground">
            <ul className="list-disc space-y-2 pl-5">
              <li>Go API (Fiber) with clean architecture layers.</li>
              <li>Vite + React frontend with shadcn/ui and ts-rest.</li>
              <li>Shared Zod schemas and OpenAPI generation.</li>
            </ul>
            <p>
              For more details, visit the public repository at{" "}
              <a
                href="https://github.com/jeheskielSunloy77/go-kickstart"
                className="font-semibold text-foreground hover:text-muted-foreground"
                target="_blank"
                rel="noreferrer"
              >
                github.com/jeheskielSunloy77/go-kickstart
              </a>
              .
            </p>
          </CardContent>
          <CardFooter className="flex flex-col gap-3 border-t sm:flex-row sm:justify-between">
            <div className="flex w-full flex-col gap-2 sm:w-auto sm:flex-row">
              {isChecking && (
                <Button
                  variant="secondary"
                  disabled
                  className="w-full sm:w-auto"
                >
                  <Spinner className="size-4" />
                  Checking session
                </Button>
              )}
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
              {!isChecking && !isAuthed && !showLogin && (
                <Button asChild className="w-full sm:w-auto">
                  <Link to="/auth/login">Get started</Link>
                </Button>
              )}
            </div>
            <Button asChild variant="outline" className="w-full sm:w-auto">
              <Link to="/health-demo">View health check demo</Link>
            </Button>
          </CardFooter>
        </Card>
      </div>
    </main>
  );
}
