"use client";

import {
  createContext,
  useContext,
  useState,
  useEffect,
  type ReactNode,
} from "react";
import type { UserResponse } from "@/lib/types";
import {
  getStoredUser,
  setStoredUser,
  setTokens,
  clearTokens,
  isAuthenticated,
} from "@/lib/auth";

interface AuthContextValue {
  user: UserResponse | null;
  isLoggedIn: boolean;
  login: (
    accessToken: string,
    refreshToken: string,
    user?: UserResponse
  ) => void;
  logout: () => void;
  setUser: (user: UserResponse) => void;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUserState] = useState<UserResponse | null>(null);
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  useEffect(() => {
    if (isAuthenticated()) {
      const stored = getStoredUser();
      if (stored) {
        setUserState(stored);
        setIsLoggedIn(true);
      } else {
        setIsLoggedIn(true);
      }
    }
  }, []);

  function login(
    accessToken: string,
    refreshToken: string,
    user?: UserResponse
  ) {
    setTokens(accessToken, refreshToken);
    if (user) {
      setStoredUser(user);
      setUserState(user);
    }
    setIsLoggedIn(true);
  }

  function logout() {
    clearTokens();
    setUserState(null);
    setIsLoggedIn(false);
    window.location.href = "/login";
  }

  function setUser(u: UserResponse) {
    setStoredUser(u);
    setUserState(u);
  }

  return (
    <AuthContext.Provider value={{ user, isLoggedIn, login, logout, setUser }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used inside AuthProvider");
  return ctx;
}
