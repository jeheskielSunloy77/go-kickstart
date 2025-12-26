import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ArrowLeft, Home, SearchX } from "lucide-react";
import { useNavigate } from "react-router-dom";

export function NotFound() {
  const navigate = useNavigate();

  return (
    <div className="flex min-h-screen items-center justify-center bg-muted/40 px-4">
      <Card className="w-full max-w-md text-center">
        <CardHeader className="space-y-4">
          <SearchX className="mx-auto h-12 w-12 text-muted-foreground" />

          <div>
            <CardTitle className="text-6xl font-bold">404</CardTitle>
            <CardDescription className="text-lg">
              Page not found
            </CardDescription>
          </div>
        </CardHeader>

        <CardContent className="space-y-6">
          <p className="text-muted-foreground">
            Sorry, the page you’re looking for doesn’t exist or has been moved.
          </p>

          <div className="flex justify-center gap-3">
            <Button variant="outline" onClick={() => navigate(-1)}>
              <ArrowLeft className="mr-2 h-4 w-4" />
              Go Back
            </Button>

            <Button onClick={() => navigate("/", { replace: true })}>
              <Home className="mr-2 h-4 w-4" />
              Go Home
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
