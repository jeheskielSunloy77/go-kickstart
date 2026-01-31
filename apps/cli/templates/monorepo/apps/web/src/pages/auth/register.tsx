import { tsr } from "@/api";
import { ThemeDropdown } from "@/components/theme-dropdown";
import { applyFieldErrors, getApiErrorMessage } from "@/lib/api-errors";
import { registerSchema } from "@/pages/auth/schemas";
import {
  Button,
  CardAction,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
  Input,
  Label,
  Spinner,
} from "@go-kickstart/ui";
import { zodResolver } from "@hookform/resolvers/zod";
import { Plus } from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { Link, useNavigate } from "react-router-dom";
import { z } from "zod";

export function RegisterPage() {
  const navigate = useNavigate();
  const [formError, setFormError] = useState<string | null>(null);

  const form = useForm<z.infer<typeof registerSchema>>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      email: "",
      username: "",
      password: "",
    },
  });

  const registerMutation = tsr.auth.register.useMutation({
    onSuccess: (response) => {
      const user = response.body;
      navigate(`/auth/verify-email?email=${encodeURIComponent(user.email)}`);
    },
    onError: (error) => {
      setFormError(getApiErrorMessage(error));
      applyFieldErrors(error, form.setError);
    },
  });

  const onSubmit = form.handleSubmit((values) => {
    setFormError(null);
    registerMutation.mutate({ body: values });
  });

  return (
    <>
      <CardHeader className="space-y-2">
        <div>
          <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
            Join the workspace
          </p>
          <CardTitle className="font-serif text-3xl">
            Create your account
          </CardTitle>
        </div>
        <CardDescription>
          We will send a verification code to your email.
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
            <Label htmlFor="username">Username</Label>
            <Input
              id="username"
              placeholder="yourname"
              {...form.register("username")}
            />
            {form.formState.errors.username?.message && (
              <p className="text-xs text-destructive">
                {form.formState.errors.username.message}
              </p>
            )}
          </div>
          <div className="space-y-2">
            <Label htmlFor="password">Password</Label>
            <Input
              id="password"
              type="password"
              placeholder="Create a strong password"
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
            disabled={registerMutation.isPending}
          >
            {registerMutation.isPending ? (
              <Spinner className="size-4" />
            ) : (
              <>
                <Plus />
                Create account
              </>
            )}
          </Button>
        </form>
      </CardContent>

      <CardFooter className="border-t">
        <div className="text-xs text-muted-foreground">
          Already have an account?{" "}
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
