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
import {
  applyFieldErrors,
  getApiErrorMessage,
  isUnauthorizedError,
} from "@/lib/api-errors";
import { verifyEmailSchema } from "@/pages/auth/schemas";
import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Link, useNavigate, useSearchParams } from "react-router-dom";
import { toast } from "sonner";
import { z } from "zod";

export function VerifyEmailPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const [formError, setFormError] = useState<string | null>(null);

  const form = useForm<z.infer<typeof verifyEmailSchema>>({
    resolver: zodResolver(verifyEmailSchema),
    defaultValues: {
      email: searchParams.get("email") ?? "",
      code: "",
    },
  });

  useEffect(() => {
    const email = searchParams.get("email");
    if (email) {
      form.setValue("email", email);
    }
  }, [searchParams, form]);

  const verifyMutation = tsr.auth.verifyEmail.useMutation({
    onSuccess: () => {
      toast.success("Email verified. You're all set.");
      navigate("/auth/me");
    },
    onError: (error) => {
      setFormError(getApiErrorMessage(error));
      applyFieldErrors(error, form.setError);
    },
  });

  const resendMutation = tsr.auth.resendVerification.useMutation({
    onSuccess: () => {
      toast.success("Verification email sent.");
    },
    onError: (error) => {
      if (isUnauthorizedError(error)) {
        setFormError("Sign in to resend a verification email.");
      } else {
        setFormError(getApiErrorMessage(error));
      }
    },
  });

  const onSubmit = form.handleSubmit((values) => {
    setFormError(null);
    verifyMutation.mutate({ body: values });
  });

  return (
    <>
      <CardHeader className="space-y-2">
        <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
          Verify your email
        </p>
        <CardTitle className="font-serif text-3xl">Enter your code</CardTitle>
        <CardDescription>
          We sent a code to your inbox. It expires shortly.
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
            <Label htmlFor="email">Email</Label>
            <Input
              id="email"
              placeholder="you@example.com"
              {...form.register("email")}
            />
            {form.formState.errors.email?.message && (
              <p className="text-xs text-destructive">
                {form.formState.errors.email.message}
              </p>
            )}
          </div>
          <div className="space-y-2">
            <Label htmlFor="code">Verification code</Label>
            <Input id="code" placeholder="123456" {...form.register("code")} />
            {form.formState.errors.code?.message && (
              <p className="text-xs text-destructive">
                {form.formState.errors.code.message}
              </p>
            )}
          </div>
          <Button
            className="w-full"
            type="submit"
            disabled={verifyMutation.isPending}
          >
            {verifyMutation.isPending ? (
              <Spinner className="size-4" />
            ) : (
              "Verify email"
            )}
          </Button>
        </form>
      </CardContent>

      <CardFooter className="flex-col items-start gap-3 border-t text-xs text-muted-foreground">
        <button
          type="button"
          className="font-semibold text-foreground hover:text-muted-foreground"
          onClick={() => resendMutation.mutate({ body: {} })}
          disabled={resendMutation.isPending}
        >
          {resendMutation.isPending
            ? "Sending..."
            : "Resend verification email"}
        </button>
        <div>
          Prefer to log in again?{" "}
          <Link
            to="/auth/login"
            className="font-semibold text-foreground hover:text-muted-foreground"
          >
            Sign in
          </Link>
        </div>
      </CardFooter>
    </>
  );
}
