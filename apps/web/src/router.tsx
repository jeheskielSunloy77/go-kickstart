import { NotFound } from "./pages/404";
import { RootErrorBoundary } from "./pages/error";
import { RequireAuth } from "@/auth/require-auth";
import { AuthLayout } from "@/pages/auth/auth-layout";
import { ForgotPasswordPage } from "@/pages/auth/forgot-password";
import { LoginPage } from "@/pages/auth/login";
import { MePage } from "@/pages/auth/me";
import { RegisterPage } from "@/pages/auth/register";
import { VerifyEmailPage } from "@/pages/auth/verify-email";
import { HealthDemoPage } from "@/pages/health-demo";
import { HomePage } from "@/pages/home";
import { Navigate, Outlet, createBrowserRouter } from "react-router";

export const router = createBrowserRouter([
  {
    // 🔹 Root layout route
    element: <Outlet />,
    ErrorBoundary: RootErrorBoundary,
    children: [
      {
        path: "/",
        element: <HomePage />,
      },
      {
        path: "/health-demo",
        element: <HealthDemoPage />,
      },
      {
        path: "/auth",
        element: <AuthLayout />,
        children: [
          { index: true, element: <Navigate to="login" replace /> },
          { path: "login", element: <LoginPage /> },
          { path: "register", element: <RegisterPage /> },
          { path: "verify-email", element: <VerifyEmailPage /> },
          { path: "forgot-password", element: <ForgotPasswordPage /> },
          {
            path: "me",
            element: (
              <RequireAuth>
                <MePage />
              </RequireAuth>
            ),
          },
        ],
      },
      {
        path: "*",
        element: <NotFound />,
      },
    ],
  },
]);
