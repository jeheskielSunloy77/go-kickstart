import { z } from "zod";

import { ZUser } from "./user.js";

export const ZAuthToken = z.object({
  token: z.string(),
  expiresAt: z.string().datetime(),
});

export const ZAuthResult = z.object({
  user: ZUser,
  token: ZAuthToken,
});

export const ZAuthRegisterRequest = z.object({
  email: z.string().email(),
  username: z.string().min(3).max(50),
  password: z.string().min(8).max(128),
});

export const ZAuthLoginRequest = z.object({
  identifier: z.string(),
  password: z.string(),
});

export const ZAuthGoogleLoginRequest = z.object({
  idToken: z.string(),
});
