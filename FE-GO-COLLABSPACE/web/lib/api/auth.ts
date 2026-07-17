import type {
  RegisterRequest,
  LoginRequest,
  TokenResponse,
  UserResponse,
  ApiResponse,
} from "../types";
import { apiClient } from "./client";
import { isMockLogin, mockTokenResponse } from "../mock-auth";

export async function register(
  data: RegisterRequest
): Promise<ApiResponse<TokenResponse>> {
  const res = await apiClient.post<ApiResponse<TokenResponse>>(
    "/auth/register",
    data
  );
  return res.data;
}

export async function login(
  data: LoginRequest
): Promise<ApiResponse<TokenResponse>> {
  if (isMockLogin(data)) {
    return mockTokenResponse();
  }

  const res = await apiClient.post<ApiResponse<TokenResponse>>(
    "/auth/login",
    data
  );
  return res.data;
}

export async function getAllUsers(): Promise<ApiResponse<UserResponse[]>> {
  const res = await apiClient.get<ApiResponse<UserResponse[]>>("/users");
  return res.data;
}
