import type {
  CreateDocRequest,
  DocumentResponse,
  ApiResponse,
} from "../types";
import { apiClient } from "./client";

export async function createDocument(
  data: CreateDocRequest
): Promise<ApiResponse<DocumentResponse>> {
  const res = await apiClient.post<ApiResponse<DocumentResponse>>(
    "/document",
    data
  );
  return res.data;
}

export async function getWorkspaceDocs(
  workspaceId: string
): Promise<ApiResponse<DocumentResponse[]>> {
  const res = await apiClient.get<ApiResponse<DocumentResponse[]>>(
    `/document/doc/${workspaceId}`
  );
  return res.data;
}

export async function getDocDetail(
  docId: string
): Promise<ApiResponse<DocumentResponse>> {
  const res = await apiClient.get<ApiResponse<DocumentResponse>>(
    `/document/${docId}`
  );
  return res.data;
}
