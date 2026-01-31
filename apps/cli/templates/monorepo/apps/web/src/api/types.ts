import { apiContract } from "@go-kickstart/openapi/contracts";
import type { ClientInferResponses, ServerInferRequest } from "@ts-rest/core";

export type TRequests = ServerInferRequest<typeof apiContract>;

// health
export type THealthGetHealthResponse = ClientInferResponses<
  typeof apiContract.health.getHealth,
  200
>["body"];

// user
export type TUserGetManyResponse = ClientInferResponses<
  typeof apiContract.user.getMany,
  200
>["body"];
export type TUserGetByIdResponse = ClientInferResponses<
  typeof apiContract.user.getById,
  200
>["body"];
export type TUserStoreResponse = ClientInferResponses<
  typeof apiContract.user.store,
  201
>["body"];
export type TUserUpdateResponse = ClientInferResponses<
  typeof apiContract.user.update,
  200
>["body"];
export type TUserDestroyResponse = ClientInferResponses<
  typeof apiContract.user.destroy,
  200
>["body"];
export type TUserRestoreResponse = ClientInferResponses<
  typeof apiContract.user.restore,
  200
>["body"];

// auth
export type TAuthMeResponse = ClientInferResponses<
  typeof apiContract.auth.me,
  200
>["body"];
export type TAuthRegisterResponse = ClientInferResponses<
  typeof apiContract.auth.register,
  201
>["body"];
export type TAuthLoginResponse = ClientInferResponses<
  typeof apiContract.auth.login,
  200
>["body"];
export type TAuthVerifyEmailResponse = ClientInferResponses<
  typeof apiContract.auth.verifyEmail,
  200
>["body"];
export type TAuthRefreshResponse = ClientInferResponses<
  typeof apiContract.auth.refresh,
  200
>["body"];
export type TAuthLogoutResponse = ClientInferResponses<
  typeof apiContract.auth.logout,
  200
>["body"];
export type TAuthResendVerificationResponse = ClientInferResponses<
  typeof apiContract.auth.resendVerification,
  200
>["body"];
export type TAuthLogoutAllResponse = ClientInferResponses<
  typeof apiContract.auth.logoutAll,
  200
>["body"];
