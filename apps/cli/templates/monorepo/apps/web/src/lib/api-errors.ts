import { isFetchError } from "@ts-rest/react-query/v5";
import type { UseFormSetError } from "react-hook-form";

type ApiFieldError = {
  field: string;
  error: string;
};

type ApiErrorBody = {
  message?: string;
  errors?: ApiFieldError[];
};

type ApiError = {
  status?: number;
  body?: ApiErrorBody;
};

export function getApiErrorMessage(error: unknown) {
  if (!error) return "Something went wrong.";

  if (isFetchError(error)) {
    return "We could not reach the server. Check your connection and try again.";
  }

  const err = error as ApiError;
  return err.body?.message || "Something went wrong.";
}

export function applyFieldErrors<T>(
  error: unknown,
  setError: UseFormSetError<T>,
) {
  const err = error as ApiError;
  const fieldErrors = err.body?.errors || [];

  fieldErrors.forEach((fieldError) => {
    setError(fieldError.field as keyof T, {
      type: "server",
      message: fieldError.error,
    });
  });
}

export function isUnauthorizedError(error: unknown) {
  if (!error) return false;
  if (isFetchError(error)) return false;
  const err = error as ApiError;
  return err.status === 401;
}
