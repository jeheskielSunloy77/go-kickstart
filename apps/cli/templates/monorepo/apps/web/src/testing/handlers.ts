import { http, HttpResponse } from "msw";
import { defaultHealthResponse, makeUser } from "@/testing/fixtures";

const apiBaseUrl =
  (import.meta.env.VITE_API_URL as string | undefined) ||
  "http://localhost:3001";
const api = apiBaseUrl.replace(/\/$/, "");

export const handlers = [
  http.get(`${api}/api/v1/auth/me`, () => {
    return HttpResponse.json(makeUser());
  }),
  http.post(`${api}/api/v1/auth/login`, () => {
    return HttpResponse.json(makeUser());
  }),
  http.post(`${api}/api/v1/auth/register`, () => {
    return HttpResponse.json(makeUser({ emailVerifiedAt: undefined }));
  }),
  http.post(`${api}/api/v1/auth/verify-email`, () => {
    return HttpResponse.json(makeUser());
  }),
  http.post(`${api}/api/v1/auth/refresh`, () => {
    return HttpResponse.json(makeUser());
  }),
  http.post(`${api}/api/v1/auth/resend-verification`, () => {
    return HttpResponse.json({ status: 200, message: "Sent", success: true });
  }),
  http.post(`${api}/api/v1/auth/logout`, () => {
    return HttpResponse.json({ status: 200, message: "Logged out", success: true });
  }),
  http.post(`${api}/api/v1/auth/logout-all`, () => {
    return HttpResponse.json({
      status: 200,
      message: "Logged out everywhere",
      success: true,
    });
  }),
  http.get(`${api}/health`, () => {
    return HttpResponse.json(defaultHealthResponse);
  }),
];
