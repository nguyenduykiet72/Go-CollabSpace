import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  createDocument,
  getWorkspaceDocs,
  getDocDetail,
} from "@/lib/api/document";
import type { CreateDocRequest } from "@/lib/types";

export const documentKeys = {
  all: ["documents"] as const,
  workspace: (workspaceId: string) =>
    ["documents", "workspace", workspaceId] as const,
  detail: (docId: string) => ["documents", docId] as const,
};

export function useWorkspaceDocs(workspaceId: string) {
  return useQuery({
    queryKey: documentKeys.workspace(workspaceId),
    queryFn: () => getWorkspaceDocs(workspaceId),
    enabled: !!workspaceId,
    select: (res) => res.data ?? [],
  });
}

export function useDocDetail(docId: string) {
  return useQuery({
    queryKey: documentKeys.detail(docId),
    queryFn: () => getDocDetail(docId),
    enabled: !!docId,
    select: (res) => res.data,
  });
}

export function useCreateDocument() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateDocRequest) => createDocument(data),
    onSuccess: (_, variables) => {
      qc.invalidateQueries({
        queryKey: documentKeys.workspace(variables.workspaceId),
      });
    },
  });
}
