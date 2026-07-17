import type {
  CreateWorkspaceRequest,
  WorkspaceResponse,
  AddMembersRequest,
  ApiResponse,
} from "../types";
import { apiClient } from "./client";

export async function createWorkspace(
  data: CreateWorkspaceRequest
): Promise<ApiResponse<WorkspaceResponse>> {
  const res = await apiClient.post<ApiResponse<WorkspaceResponse>>(
    "/workspace",
    data
  );
  return res.data;
}

export async function getWorkspaceById(
  workspaceId: string
): Promise<ApiResponse<WorkspaceResponse>> {
  const res = await apiClient.get<ApiResponse<WorkspaceResponse>>(
    `/workspace/${workspaceId}`
  );
  return res.data;
}

export async function addWorkspaceMembers(
  workspaceId: string,
  data: AddMembersRequest
): Promise<ApiResponse<{ count: number }>> {
  const res = await apiClient.post<ApiResponse<{ count: number }>>(
    `/workspace/${workspaceId}/members`,
    data
  );
  return res.data;
}
