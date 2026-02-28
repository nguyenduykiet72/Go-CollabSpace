import { useMutation } from "@tanstack/react-query";
import { login as loginApi, register as registerApi } from "@/lib/api/auth";
import { getApiError } from "@/lib/api/client";
import type { LoginRequest, RegisterRequest } from "@/lib/types";

export function useLogin() {
  return useMutation({
    mutationFn: (data: LoginRequest) => loginApi(data),
    onError: (err) => {
      console.error("Login error:", getApiError(err));
    },
  });
}

export function useRegister() {
  return useMutation({
    mutationFn: (data: RegisterRequest) => registerApi(data),
    onError: (err) => {
      console.error("Register error:", getApiError(err));
    },
  });
}
