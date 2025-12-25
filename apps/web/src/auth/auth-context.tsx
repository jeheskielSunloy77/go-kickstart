import { tsr } from "@/api";
import type { TAuthMeResponse } from "@/api/types";
import { createContext, useContext, type ReactNode } from "react";

const AuthContext = createContext<{
  user?: TAuthMeResponse;
}>(null!);

export function AuthProvider(props: { children: ReactNode }) {
  const user = tsr.auth.me.useQuery({
    queryKey: ["auth", "me"],
    retry: false,
    select: (response) => response.body,
  });

  return (
    <AuthContext.Provider value={{ user: user.data }}>
      {props.children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}

export function useMustUser() {
  return useContext(AuthContext).user!;
}
