import { z } from "zod";

const envVarsSchema = z.object({
  VITE_API_URL: z.url().default("http://localhost:8080"),
  VITE_GOOGLE_CLIENT_ID: z.string().default(""),
  VITE_ENV: z
    .enum(["production", "development", "staging"])
    .default("development"),
});

const parseResult = envVarsSchema.safeParse(process.env);

if (!parseResult.success) {
  console.error(
    "❌ Invalid environment variables:",
    z.treeifyError(parseResult.error),
  );
  throw new Error("Invalid environment variables");
}

const envVars = parseResult.data;

// export individual variables
export const ENV = envVars.VITE_ENV;
export const API_URL = envVars.VITE_API_URL;
export const GOOGLE_CLIENT_ID = envVars.VITE_GOOGLE_CLIENT_ID;
