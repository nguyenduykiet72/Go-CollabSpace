"use client";

import { use } from "react";
import { useWorkspace } from "@/hooks/use-workspace";
import { useWorkspaceDocs } from "@/hooks/use-document";
import { DocumentList } from "@/components/document/document-list";
import { CreateDocumentModal } from "@/components/document/create-document-modal";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
import { formatDate } from "@/lib/utils";
import { CalendarDays, FileText, Hash, Users } from "lucide-react";

interface Props {
  params: Promise<{ workspaceId: string }>;
}

export default function WorkspacePage({ params }: Props) {
  const { workspaceId } = use(params);
  const { data: workspace, isLoading: wsLoading } = useWorkspace(workspaceId);
  const { data: docs = [], isLoading: docsLoading } =
    useWorkspaceDocs(workspaceId);

  return (
    <div className="mx-auto max-w-6xl px-8 py-8">
      <div className="mb-8">
        {wsLoading ? (
          <div className="space-y-2">
            <Skeleton className="h-8 w-48" />
            <Skeleton className="h-4 w-32" />
          </div>
        ) : workspace ? (
          <>
            <div className="flex flex-col gap-5 border-b border-stone-200 pb-6 lg:flex-row lg:items-end lg:justify-between">
              <div>
                <div className="mb-3 flex items-center gap-2">
                  <Badge variant="success">Workspace</Badge>
                  <span className="flex items-center gap-1 text-xs text-stone-500">
                    <Hash className="h-3.5 w-3.5" />
                    {workspace.slug}
                  </span>
                </div>
                <h1 className="text-3xl font-semibold tracking-tight text-stone-950">
                  {workspace.name}
                </h1>
                <div className="mt-3 flex flex-wrap items-center gap-4">
                  <span className="flex items-center gap-1.5 text-xs text-stone-500">
                    <CalendarDays className="h-3.5 w-3.5" />
                    Created {formatDate(workspace.createdAt)}
                  </span>
                  <span className="flex items-center gap-1.5 text-xs text-stone-500">
                    <Users className="h-3.5 w-3.5" />
                    Members managed by backend roles
                  </span>
                </div>
              </div>
              <CreateDocumentModal workspaceId={workspaceId} />
            </div>
          </>
        ) : (
          <p className="text-sm text-red-600">Workspace not found.</p>
        )}
      </div>

      <div className="grid gap-6 lg:grid-cols-[1fr_280px]">
        <div className="overflow-hidden rounded-lg border border-stone-200 bg-white shadow-sm shadow-stone-200/60">
          <div className="flex items-center justify-between border-b border-stone-100 bg-teal-50/50 px-5 py-4">
            <div className="flex items-center gap-2">
              <FileText className="h-4 w-4 text-teal-700" />
              <h2 className="text-sm font-semibold text-stone-950">
                Documents
              </h2>
              {!docsLoading && <Badge variant="muted">{docs.length}</Badge>}
            </div>
            <span className="text-xs text-stone-500">
              {docsLoading ? "Loading" : "Sorted by newest"}
            </span>
          </div>
          <DocumentList
            docs={docs}
            workspaceId={workspaceId}
            isLoading={docsLoading}
          />
        </div>

        <aside className="rounded-lg border border-amber-200 bg-amber-50/70 p-5 shadow-sm shadow-stone-200/60">
          <div className="flex items-center gap-2">
            <Users className="h-4 w-4 text-amber-700" />
            <h2 className="text-sm font-semibold text-stone-950">
              Collaboration
            </h2>
          </div>
          <div className="mt-4 space-y-3 text-sm text-stone-600">
            <p className="leading-6">
              Workspace access is enforced by the Go API. Add member controls
              here when those endpoints are wired into the frontend.
            </p>
            <div className="rounded-md border border-amber-100 bg-white/70 px-3 py-2">
              <p className="text-xs font-medium uppercase text-amber-700">
                Roles
              </p>
              <p className="mt-1 text-sm font-medium text-stone-800">
                Owner, Admin, Editor, Viewer
              </p>
            </div>
          </div>
        </aside>
      </div>
    </div>
  );
}
