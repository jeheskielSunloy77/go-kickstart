import { tsr } from "./api";
import "./index.css";
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
import { isFetchError } from "@ts-rest/react-query/v5";
import { RotateCcw } from "lucide-react";

export default function App() {
  const q = tsr.health.getHealth.useQuery({
    queryKey: ["health"],
    select: (response) => response.body,
  });

  return (
    <main className="max-w-4xl mx-auto h-screen flex items-center">
      <Card className="max-w-xl mx-auto w-full">
        <CardHeader>
          <CardTitle>Example Page</CardTitle>
          <CardDescription>
            This page demonstrates how to use the api client with React Query.
            it will fetch server health status.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="h-96 flex items-center justify-center bg-muted/50 rounded-md p-4">
            {q.isLoading && <Spinner className="size-10" />}
            {q.isError && (
              <div className="text-destructive text-center">
                {getErrorMessage(q.error)}
              </div>
            )}
            {q.isSuccess && (
              <pre className="text-sm whitespace-pre-wrap">
                {JSON.stringify(q.data.data, null, 2)}
              </pre>
            )}
          </div>
        </CardContent>
        {!q.isLoading && (
          <CardFooter>
            <Button onClick={() => q.refetch()} className="w-full">
              Refetch <RotateCcw />
            </Button>
          </CardFooter>
        )}
      </Card>
    </main>
  );
}

function getErrorMessage(error: unknown) {
  if (!error) return null;

  if (isFetchError(error))
    return "We could not retrieve this data. Please check your internet connection.";

  const err = error as {
    status: number;
    body: {
      message: string;
    };
  };

  return err.body.message;
}
