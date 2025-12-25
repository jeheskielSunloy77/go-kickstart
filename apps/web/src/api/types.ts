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
