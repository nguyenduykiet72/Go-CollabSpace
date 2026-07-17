import axios, { type AxiosError } from "axios";
import { getAccessToken, clearTokens } from "../auth";

// In dev, calls go through Next.js rewrite proxy → no CORS needed.
// In prod, set NEXT_PUBLIC_API_URL to your backend URL.
const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL ||
  (typeof window !== "undefined" ? "/api/v1" : "http://localhost:9999/api/v1");

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

apiClient.interceptors.request.use((config) => {
  const token = getAccessToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

apiClient.interceptors.response.use(
  (res) => res,
  (error: AxiosError) => {
    if (error.response?.status === 401) {
      clearTokens();
      if (typeof window !== "undefined") {
        window.location.href = "/login";
      }
    }
    return Promise.reject(error);
  }
);

export function getApiError(error: unknown): string {
  if (axios.isAxiosError(error)) {
    const data = error.response?.data as
      | { error?: { message?: string } }
      | undefined;
    return (
      data?.error?.message ||
      error.message ||
      "An unexpected error occurred"
    );
  }
  return "An unexpected error occurred";
}
