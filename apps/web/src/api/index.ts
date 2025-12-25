import { API_URL } from "@/config/env";
import { apiContract } from "@go-kickstart/openapi/contracts";
import {
  initClient,
  type ApiFetcherArgs,
  type ClientArgs,
} from "@ts-rest/core";
import { initTsrReactQuery } from "@ts-rest/react-query/v5";
import axios, {
  AxiosError,
  isAxiosError,
  type AxiosResponse,
  type Method,
} from "axios";

type Headers = Awaited<
  ReturnType<NonNullable<Parameters<typeof initClient>[1]["api"]>>
>["headers"];

export type TApiClient = ReturnType<typeof useApiClient>;

const getToken = async ({
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  template,
}: {
  template: "custom" | "service";
}): Promise<string | null> => {
  // Implement your token retrieval logic here.
  // This is a placeholder implementation.
  return null;
};

const createApiFetcher =
  ({ isBlob = false }: { isBlob?: boolean } = {}) =>
  async ({ path, method, headers, body, fetchOptions }: ApiFetcherArgs) => {
    const token = await getToken({ template: "custom" });

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const makeRequest = async (retryCount = 0): Promise<any> => {
      try {
        const result = await axios.request({
          method: method as Method,
          url: path,
          headers: {
            ...headers,
            ...(token ? { Authorization: `Bearer ${token}` } : {}),
          },
          data: body,
          signal: fetchOptions?.signal ?? undefined,
          ...(isBlob ? { responseType: "blob" } : {}),
        });
        return {
          status: result.status,
          body: result.data,
          headers: result.headers as unknown as Headers,
        };
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
      } catch (e: Error | AxiosError | any) {
        if (isAxiosError(e)) {
          const error = e as AxiosError;
          const response = error.response as AxiosResponse;

          // If unauthorized and we haven't retried yet, retry
          if (response?.status === 401 && retryCount < 2) {
            return makeRequest(retryCount + 1);
          }

          return {
            status: response?.status || 500,
            body: response?.data || { message: "Internal server error" },
            headers: (response?.headers as unknown as Headers) || {},
          };
        }
        throw e;
      }
    };

    return makeRequest();
  };

const createClientArgs = ({
  isBlob = false,
}: {
  isBlob?: boolean;
} = {}): ClientArgs => ({
  baseUrl: API_URL,
  baseHeaders: {
    "Content-Type": "application/json",
  },
  api: createApiFetcher({ isBlob }),
});

export const useApiClient = ({ isBlob = false }: { isBlob?: boolean } = {}) =>
  initClient(apiContract, createClientArgs({ isBlob }));

export const tsr = initTsrReactQuery(apiContract, createClientArgs());
