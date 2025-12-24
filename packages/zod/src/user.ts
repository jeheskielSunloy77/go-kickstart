import { z } from "zod";

export const ZDeletedAt = z.object({
  Time: z.string().datetime(),
  Valid: z.boolean(),
});

export const ZUser = z.object({
  id: z.string().uuid(),
  email: z.string().email(),
  username: z.string().min(3).max(50),
  googleId: z.string().optional(),
  lastLoginAt: z.string().datetime().optional(),
  createdAt: z.string().datetime(),
  updatedAt: z.string().datetime(),
  deletedAt: ZDeletedAt,
});

export const ZStoreUserRequest = z.object({
  email: z.string().email(),
  username: z.string().min(3).max(50),
  password: z.string().min(8).max(128),
  googleId: z.string().optional(),
});

export const ZUpdateUserRequest = z.object({
  email: z.string().email().optional(),
  username: z.string().min(3).max(50).optional(),
  password: z.string().min(8).max(128).optional(),
});
