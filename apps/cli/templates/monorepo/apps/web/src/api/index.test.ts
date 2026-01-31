import axios, { AxiosError, type AxiosResponse } from "axios";
import { describe, expect, it, vi, afterEach } from "vitest";
import { __testing } from "@/api";

type AxiosResponseLike = AxiosResponse;

const apiBase =
  (import.meta.env.VITE_API_URL as string | undefined) ||
  "http://localhost:3001";

function makeAxiosError(status: number, data: unknown) {
  const response = {
    status,
    data,
    headers: {},
    config: {},
  } as AxiosResponseLike;
  return new AxiosError("Request failed", "ERR_BAD_REQUEST", {}, {}, response);
}

afterEach(() => {
  vi.restoreAllMocks();
});

describe("api fetcher", () => {
  it("retries once after refresh succeeds", async () => {
    const requestSpy = vi.spyOn(axios, "request");
    const postSpy = vi.spyOn(axios, "post");

    const unauthorizedError = makeAxiosError(401, { message: "Unauthorized" });

    requestSpy
      .mockRejectedValueOnce(unauthorizedError)
      .mockResolvedValueOnce({
        status: 200,
        data: { ok: true },
        headers: { "x-test": "1" },
      });

    postSpy.mockResolvedValue({ status: 200, data: {} });

    const fetcher = __testing.createApiFetcher();

    const result = await fetcher({
      path: `${apiBase}/api/v1/auth/me`,
      method: "GET",
      headers: {},
      body: undefined,
      fetchOptions: undefined,
    });

    expect(result.status).toBe(200);
    expect(result.body).toEqual({ ok: true });
    expect(requestSpy).toHaveBeenCalledTimes(2);
    expect(postSpy).toHaveBeenCalledTimes(1);
  });

  it("returns the error response when refresh fails", async () => {
    const requestSpy = vi.spyOn(axios, "request");
    const postSpy = vi.spyOn(axios, "post");

    const unauthorizedError = makeAxiosError(401, { message: "Unauthorized" });

    requestSpy.mockRejectedValueOnce(unauthorizedError);
    postSpy.mockRejectedValue(new Error("Refresh failed"));

    const fetcher = __testing.createApiFetcher();

    const result = await fetcher({
      path: `${apiBase}/api/v1/auth/me`,
      method: "GET",
      headers: {},
      body: undefined,
      fetchOptions: undefined,
    });

    expect(result.status).toBe(401);
    expect(result.body).toEqual({ message: "Unauthorized" });
    expect(postSpy).toHaveBeenCalledTimes(1);
  });
});
