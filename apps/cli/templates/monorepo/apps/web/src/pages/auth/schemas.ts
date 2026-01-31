import { z } from "zod";

export const loginSchema = z.object({
  identifier: z.string().min(1, "Identifier is required"),
  password: z.string().min(1, "Password is required"),
});

export const registerSchema = z.object({
  email: z.email("Enter a valid email"),
  username: z.string().min(3, "Username is too short").max(50),
  password: z
    .string()
    .min(8, "Password must be at least 8 characters")
    .max(128),
});

export const verifyEmailSchema = z.object({
  email: z.string().email("Enter a valid email"),
  code: z.string().min(4, "Code is too short").max(10),
});
