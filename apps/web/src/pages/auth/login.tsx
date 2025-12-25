import { tsr } from "@/api";
import { ThemeDropdown } from "@/components/theme-dropdown";
import { Button } from "@/components/ui/button";
import {
  CardAction,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Spinner } from "@/components/ui/spinner";
import { GOOGLE_CLIENT_ID } from "@/config/env";
import { applyFieldErrors, getApiErrorMessage } from "@/lib/api-errors";
import { loginSchema } from "@/pages/auth/schemas";
import { zodResolver } from "@hookform/resolvers/zod";
import { GoogleLogin } from "@react-oauth/google";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { Link, useNavigate } from "react-router-dom";
import { z } from "zod";

export function LoginPage() {
  const navigate = useNavigate();
  const [formError, setFormError] = useState<string | null>(null);

  const form = useForm<z.infer<typeof loginSchema>>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      identifier: "",
      password: "",
    },
  });

  const loginMutation = tsr.auth.login.useMutation({
    onSuccess: (response) => {
      const user = response.body;
      if (!user.emailVerifiedAt) {
        navigate(`/auth/verify-email?email=${encodeURIComponent(user.email)}`);
        return;
      }
      navigate("/auth/me");
    },
    onError: (error) => {
      setFormError(getApiErrorMessage(error));
      applyFieldErrors(error, form.setError);
    },
  });

  const googleMutation = tsr.auth.googleLogin.useMutation({
    onSuccess: () => {
      navigate("/auth/me");
    },
    onError: (error) => {
      setFormError(getApiErrorMessage(error));
    },
  });

  const onSubmit = form.handleSubmit((values) => {
    setFormError(null);
    loginMutation.mutate({ body: values });
  });

  return (
    <>
      <CardHeader className="space-y-2">
        <div>
          <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
            Welcome back
          </p>
          <CardTitle className="font-display text-3xl">
            Sign in to continue
          </CardTitle>
        </div>
        <CardDescription>
          Access your session without ever touching a token.
        </CardDescription>
        <CardAction>
          <ThemeDropdown />
        </CardAction>
      </CardHeader>

      <CardContent className="space-y-4">
        {formError && (
          <div className="rounded-2xl border border-destructive/20 bg-destructive/10 px-4 py-3 text-sm text-destructive">
            {formError}
          </div>
        )}

        <form className="space-y-4" onSubmit={onSubmit}>
          <div className="space-y-2">
            <Label htmlFor="identifier">Email or username</Label>
            <Input
              id="identifier"
              placeholder="you@example.com"
              {...form.register("identifier")}
            />
            {form.formState.errors.identifier?.message && (
              <p className="text-xs text-destructive">
                {form.formState.errors.identifier.message}
              </p>
            )}
          </div>
          <div className="space-y-2">
            <Label htmlFor="password">Password</Label>
            <Input
              id="password"
              type="password"
              placeholder="••••••••"
              {...form.register("password")}
            />
            {form.formState.errors.password?.message && (
              <p className="text-xs text-destructive">
                {form.formState.errors.password.message}
              </p>
            )}
          </div>
          <Button
            className="w-full"
            type="submit"
            disabled={loginMutation.isPending}
          >
            {loginMutation.isPending ? (
              <Spinner className="size-4" />
            ) : (
              "Sign in"
            )}
          </Button>
        </form>
      </CardContent>

      <CardFooter className="flex-col items-stretch gap-4 border-t">
        <div className="flex items-center justify-between text-xs text-muted-foreground">
          <Link to="/auth/forgot-password" className="hover:text-foreground">
            Forgot password?
          </Link>
          <Link
            to="/auth/register"
            className="font-semibold text-foreground hover:text-muted-foreground"
          >
            Create account
          </Link>
        </div>

        {GOOGLE_CLIENT_ID && (
          <div className="space-y-3">
            <div className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
              Or
            </div>
            <div className="rounded-2xl border border-border bg-background px-3 py-2 shadow-sm">
              <GoogleLogin
                onSuccess={(credentialResponse) => {
                  if (!credentialResponse.credential) return;
                  setFormError(null);
                  googleMutation.mutate({
                    body: { idToken: credentialResponse.credential },
                  });
                }}
                onError={() =>
                  setFormError("Google sign-in failed. Please try again.")
                }
                width="350"
              />
            </div>
          </div>
        )}
      </CardFooter>
    </>
  );
}
