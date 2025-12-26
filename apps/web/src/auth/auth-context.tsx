import type { TAuthMeResponse } from "@/api/types";
import { createContext, useContext, type ReactNode } from "react";

const AuthContext = createContext<{
  user: TAuthMeResponse;
}>(null!);

export function AuthProvider(props: {
  children: ReactNode;
  user: TAuthMeResponse;
}) {
  return (
    <AuthContext.Provider value={{ user: props.user }}>
      {props.children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}
