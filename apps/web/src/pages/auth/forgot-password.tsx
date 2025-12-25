import { Button } from "@/components/ui/button";
import {
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { Link } from "react-router-dom";
import { toast } from "sonner";
import { z } from "zod";

const ForgotPasswordSchema = z.object({
  email: z.email(),
});

export function ForgotPasswordPage() {
  const form = useForm<z.infer<typeof ForgotPasswordSchema>>({
    resolver: zodResolver(ForgotPasswordSchema),
    defaultValues: {
      email: "",
    },
  });

  const onSubmit = form.handleSubmit(() => {
    toast.message("Password reset isn't wired yet.", {
      description: "Add backend support to send reset links.",
    });
  });

  return (
    <>
      <CardHeader className="space-y-2">
        <div>
          <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
            Reset password
          </p>
          <CardTitle className="font-display text-3xl">
            Recover your access
          </CardTitle>
        </div>
        <CardDescription>
          Enter your email and we will send a reset link.
        </CardDescription>
      </CardHeader>

      <CardContent className="space-y-4">
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
          <Button className="w-full" type="submit">
            Send reset link
          </Button>
        </form>
      </CardContent>

      <CardFooter className="border-t">
        <div className="text-xs text-muted-foreground">
          Remembered your password?{" "}
          <Link
            to="/auth/login"
            className="font-semibold text-foreground hover:text-muted-foreground"
          >
            Back to sign in
          </Link>
        </div>
      </CardFooter>
    </>
  );
}
