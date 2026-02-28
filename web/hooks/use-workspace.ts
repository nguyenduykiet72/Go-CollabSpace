import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  createWorkspace,
  getWorkspaceById,
  addWorkspaceMembers,
} from "@/lib/api/workspace";
import type { CreateWorkspaceRequest, AddMembersRequest } from "@/lib/types";

export const workspaceKeys = {
  all: ["workspaces"] as const,
  detail: (id: string) => ["workspaces", id] as const,
};

export function useWorkspace(workspaceId: string) {
  return useQuery({
    queryKey: workspaceKeys.detail(workspaceId),
    queryFn: () => getWorkspaceById(workspaceId),
    enabled: !!workspaceId,
    select: (res) => res.data,
  });
}

export function useCreateWorkspace() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateWorkspaceRequest) => createWorkspace(data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: workspaceKeys.all });
    },
  });
}

export function useAddMembers(workspaceId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: AddMembersRequest) =>
      addWorkspaceMembers(workspaceId, data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: workspaceKeys.detail(workspaceId) });
    },
  });
}
