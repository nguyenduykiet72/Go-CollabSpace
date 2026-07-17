import type { ApiResponse, LoginRequest, TokenResponse, UserResponse } from "./types";

// Temporary mock login for frontend UI work. Delete this file and the imports
// from auth.ts/login-form.tsx when the backend auth flow is ready.
export const MOCK_LOGIN = {
  email: "demo@gocollab.dev",
  password: "password123",
};

export const MOCK_USER: UserResponse = {
  id: "mock-user-1",
  email: MOCK_LOGIN.email,
  fullName: "Demo Learner",
  avatar: "",
  createdAt: new Date().toISOString(),
};

export function isMockLogin(data: LoginRequest) {
  return data.email === MOCK_LOGIN.email && data.password === MOCK_LOGIN.password;
}

export function mockTokenResponse(): ApiResponse<TokenResponse> {
  return {
    data: {
      accessToken: "mock-access-token",
      refreshToken: "mock-refresh-token",
    },
    message: "Mock login successful",
  };
}
