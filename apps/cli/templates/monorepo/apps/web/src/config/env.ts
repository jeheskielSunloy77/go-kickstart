import { z } from "zod";

const envVarsSchema = z.object({
  VITE_API_URL: z.url(),
  VITE_ENV: z
    .enum(["production", "development", "staging"])
    .default("development"),
  VITE_GOOGLE_AUTH_ENABLED: z
    .enum(["true", "false"])
    .default("false")
    .transform((val) => val === "true"),
});

const parseResult = envVarsSchema.safeParse(import.meta.env);

if (!parseResult.success) {
  console.error(
    "❌ Invalid environment variables:",
    z.treeifyError(parseResult.error),
  );
  throw new Error("Invalid environment variables");
}

// export individual variables
export const ENV = parseResult.data.VITE_ENV;
export const API_URL = parseResult.data.VITE_API_URL;
export const GOOGLE_AUTH_ENABLED = parseResult.data.VITE_GOOGLE_AUTH_ENABLED;
