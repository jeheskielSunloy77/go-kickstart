import { tsr } from "@/api";
import { useAuth } from "@/auth/auth-context";
import { ThemeDropdown } from "@/components/theme-dropdown";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { getApiErrorMessage } from "@/lib/api-errors";
import { useQueryClient } from "@tanstack/react-query";
import { Link, useNavigate } from "react-router-dom";
import { toast } from "sonner";

export function MePage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const { user } = useAuth();

  const logoutMutation = tsr.auth.logout.useMutation({
    onSuccess: () => {
      queryClient.removeQueries({ queryKey: ["auth", "me"] });
      navigate("/auth/login");
    },
    onError: (error) => {
      toast.error(getApiErrorMessage(error));
    },
  });

  const logoutAllMutation = tsr.auth.logoutAll.useMutation({
    onSuccess: () => {
      queryClient.removeQueries({ queryKey: ["auth", "me"] });
      navigate("/auth/login");
    },
    onError: (error) => {
      toast.error(getApiErrorMessage(error));
    },
  });

  const resendMutation = tsr.auth.resendVerification.useMutation({
    onSuccess: () => {
      toast.success("Verification email sent.");
    },
    onError: (error) => {
      toast.error(getApiErrorMessage(error));
    },
  });

  return (
    <>
      <CardHeader className="space-y-2">
        <div>
          <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
            Your session
          </p>
          <CardTitle className="font-serif text-3xl">
            Welcome, {user.username}
          </CardTitle>
        </div>
        <CardDescription>
          This profile is backed by your database session and HTTP-only cookies.
        </CardDescription>
        <CardAction>
          <ThemeDropdown />
        </CardAction>
      </CardHeader>

      <CardContent className="space-y-4">
        <Card className="shadow-sm">
          <CardContent className="space-y-3 text-sm">
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">Email</span>
              <span className="font-medium text-foreground">{user.email}</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">Username</span>
              <span className="font-medium text-foreground">
                {user.username}
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">User ID</span>
              <span className="font-mono text-xs text-muted-foreground">
                {user.id}
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">Email status</span>
              <span className="font-medium text-foreground">
                {user.emailVerifiedAt ? "Verified" : "Not verified"}
              </span>
            </div>
          </CardContent>
        </Card>

        {!user.emailVerifiedAt && (
          <Card className="gap-0 border-primary/20 bg-primary/10 py-0 shadow-none">
            <CardContent className="space-y-3 px-4 py-3 text-sm text-foreground">
              <p>
                Your email is not verified yet. You can still verify it or
                resend the code.
              </p>
              <div className="flex flex-wrap gap-2">
                <Button
                  type="button"
                  variant="secondary"
                  onClick={() => resendMutation.mutate({ body: {} })}
                  disabled={resendMutation.isPending}
                >
                  {resendMutation.isPending ? "Sending..." : "Resend code"}
                </Button>
                <Button type="button" asChild>
                  <Link
                    to={`/auth/verify-email?email=${encodeURIComponent(user.email)}`}
                  >
                    Verify email
                  </Link>
                </Button>
              </div>
            </CardContent>
          </Card>
        )}
      </CardContent>

      <CardFooter className="border-t">
        <div className="grid w-full gap-3 sm:grid-cols-2">
          <Button
            variant="secondary"
            onClick={() => logoutMutation.mutate({ body: {} })}
            disabled={logoutMutation.isPending}
          >
            {logoutMutation.isPending ? "Signing out..." : "Sign out"}
          </Button>
          <Button
            variant="destructive"
            onClick={() => logoutAllMutation.mutate({ body: {} })}
            disabled={logoutAllMutation.isPending}
          >
            {logoutAllMutation.isPending
              ? "Signing out..."
              : "Sign out everywhere"}
          </Button>
        </div>
      </CardFooter>
    </>
  );
}
