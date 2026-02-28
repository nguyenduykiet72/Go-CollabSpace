"use client";

import { use } from "react";
import { useWorkspace } from "@/hooks/use-workspace";
import { useWorkspaceDocs } from "@/hooks/use-document";
import { DocumentList } from "@/components/document/document-list";
import { CreateDocumentModal } from "@/components/document/create-document-modal";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
import { formatDate } from "@/lib/utils";
import { Users, Hash } from "lucide-react";

interface Props {
  params: Promise<{ workspaceId: string }>;
}

export default function WorkspacePage({ params }: Props) {
  const { workspaceId } = use(params);
  const { data: workspace, isLoading: wsLoading } = useWorkspace(workspaceId);
  const { data: docs = [], isLoading: docsLoading } =
    useWorkspaceDocs(workspaceId);

  return (
    <div className="px-8 py-8 max-w-4xl">
      {/* Workspace header */}
      <div className="mb-8">
        {wsLoading ? (
          <div className="space-y-2">
            <Skeleton className="h-8 w-48" />
            <Skeleton className="h-4 w-32" />
          </div>
        ) : workspace ? (
          <>
            <div className="flex items-start justify-between">
              <div>
                <h1 className="text-2xl font-bold text-zinc-900 tracking-tight">
                  {workspace.name}
                </h1>
                <div className="flex items-center gap-3 mt-2">
                  <span className="flex items-center gap-1 text-xs text-zinc-500">
                    <Hash className="h-3 w-3" />
                    {workspace.slug}
                  </span>
                  <span className="text-xs text-zinc-400">·</span>
                  <span className="text-xs text-zinc-500">
                    Created {formatDate(workspace.createdAt)}
                  </span>
                </div>
              </div>
              <CreateDocumentModal workspaceId={workspaceId} />
            </div>
          </>
        ) : (
          <p className="text-sm text-red-500">Workspace not found.</p>
        )}
      </div>

      {/* Documents */}
      <div className="rounded-xl border border-zinc-200 bg-white overflow-hidden">
        <div className="flex items-center justify-between px-4 py-3 border-b border-zinc-100">
          <div className="flex items-center gap-2">
            <h2 className="text-sm font-semibold text-zinc-900">Documents</h2>
            {!docsLoading && (
              <Badge variant="muted">{docs.length}</Badge>
            )}
          </div>
        </div>
        <DocumentList
          docs={docs}
          workspaceId={workspaceId}
          isLoading={docsLoading}
        />
      </div>
    </div>
  );
}
