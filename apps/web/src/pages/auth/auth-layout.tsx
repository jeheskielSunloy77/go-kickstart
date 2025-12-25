import { Card } from "@/components/ui/card";
import { CheckCircle2, Lock } from "lucide-react";
import { Outlet } from "react-router-dom";

const features = [
  "Short-lived access tokens with silent refresh.",
  "HTTP-only cookies with zero token handling in the browser.",
  "Database-backed sessions with rotation on refresh.",
];

export function AuthLayout() {
  return (
    <main className="min-h-screen bg-[radial-gradient(1200px_circle_at_15%_0%,rgba(253,230,138,0.45),transparent_60%),radial-gradient(900px_circle_at_100%_80%,rgba(125,211,252,0.35),transparent_55%)]">
      <div className="mx-auto flex min-h-screen w-full items-center px-6 py-12">
        <div className="grid w-full gap-10 lg:grid-cols-2 divide-x">
          <section className="hidden lg:block">
            <div className="p-10 max-w-xl mx-auto space-y-10">
              <div className="space-y-4">
                <div>
                  <span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
                    Go Kickstart
                  </span>
                  <h1 className="font-display text-4xl font-semibold text-foreground">
                    Security-first auth.
                    <span className="block text-muted-foreground">
                      Quietly powerful UX.
                    </span>
                  </h1>
                </div>
                <p className="text-sm leading-6 text-muted-foreground">
                  Built for production: rotation-ready refresh tokens,
                  profile-backed sessions, and route protection that stays
                  invisible to users.
                </p>
              </div>
              <div className="space-y-3">
                {features.map((feature) => (
                  <div
                    key={feature}
                    className="flex items-center gap-3 rounded-lg border border-border bg-background/70 px-4 py-3 text-sm text-foreground/90"
                  >
                    <CheckCircle2 className="size-4 text-primary" />
                    {feature}
                  </div>
                ))}
              </div>
              <div className="flex gap-3 rounded-lg border px-6 py-4">
                <Lock className="size-4 mt-0.5" />
                <div>
                  <p className="font-semibold text-sm">
                    No JWTs in storage. Ever.
                  </p>
                  <p className="text-sm text-muted-foreground">
                    Cookies only, protected by the server.
                  </p>
                </div>
              </div>
            </div>
          </section>
          <section className="flex items-center justify-center">
            <Card className="duration-700 fade-in animate-in slide-in-from-bottom-6 min-w-md">
              <Outlet />
            </Card>
          </section>
        </div>
      </div>
    </main>
  );
}
