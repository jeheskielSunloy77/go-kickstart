import type { TAuthMeResponse } from "@/api/types";

const now = new Date("2024-01-01T00:00:00.000Z").toISOString();

export function makeUser(overrides: Partial<TAuthMeResponse> = {}): TAuthMeResponse {
  return {
    id: "11111111-1111-1111-1111-111111111111",
    email: "jane.doe@example.com",
    username: "janedoe",
    isAdmin: false,
    createdAt: now,
    updatedAt: now,
    emailVerifiedAt: now,
    ...overrides,
  };
}

export const defaultHealthResponse = {
  status: 200,
  message: "OK",
  success: true,
  data: {
    status: "healthy",
    timestamp: "2024-01-01T00:00:00.000Z",
    environment: "test",
    checks: {
      database: {
        status: "healthy",
        response_time: "12ms",
      },
    },
  },
};
