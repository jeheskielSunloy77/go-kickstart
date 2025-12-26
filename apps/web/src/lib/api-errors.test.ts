import { describe, expect, it, vi } from "vitest";

vi.mock("@ts-rest/react-query/v5", async () => {
  const actual = await vi.importActual<typeof import("@ts-rest/react-query/v5")>(
    "@ts-rest/react-query/v5",
  );
  return {
    ...actual,
    isFetchError: (error: unknown) =>
      Boolean((error as { __fetchError?: boolean }).__fetchError),
  };
});

import {
  applyFieldErrors,
  getApiErrorMessage,
  isUnauthorizedError,
} from "@/lib/api-errors";

describe("api error helpers", () => {
  it("returns a safe message when no error is provided", () => {
    expect(getApiErrorMessage(null)).toBe("Something went wrong.");
  });

  it("returns a fetch error message for network failures", () => {
    const message = getApiErrorMessage({ __fetchError: true });
    expect(message).toBe(
      "We could not reach the server. Check your connection and try again.",
    );
  });

  it("returns the server-provided message when available", () => {
    const message = getApiErrorMessage({
      body: { message: "Bad credentials" },
    });
    expect(message).toBe("Bad credentials");
  });

  it("applies field errors from the API response", () => {
    const setError = vi.fn();

    applyFieldErrors(
      {
        body: {
          errors: [
            { field: "email", error: "Invalid email" },
            { field: "password", error: "Too short" },
          ],
        },
      },
      setError,
    );

    expect(setError).toHaveBeenCalledWith("email", {
      type: "server",
      message: "Invalid email",
    });
    expect(setError).toHaveBeenCalledWith("password", {
      type: "server",
      message: "Too short",
    });
  });

  it("detects unauthorized errors", () => {
    expect(isUnauthorizedError({ status: 401 })).toBe(true);
    expect(isUnauthorizedError({ status: 403 })).toBe(false);
    expect(isUnauthorizedError({ __fetchError: true })).toBe(false);
  });
});
