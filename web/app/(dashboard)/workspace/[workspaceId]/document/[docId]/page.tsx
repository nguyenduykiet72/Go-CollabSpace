"use client";

import { use } from "react";
import Link from "next/link";
import { useDocDetail } from "@/hooks/use-document";
import { Skeleton } from "@/components/ui/skeleton";
import { formatDate } from "@/lib/utils";
import { ChevronLeft, Clock, User } from "lucide-react";

interface Props {
  params: Promise<{ workspaceId: string; docId: string }>;
}

export default function DocumentPage({ params }: Props) {
  const { workspaceId, docId } = use(params);
  const { data: doc, isLoading, isError } = useDocDetail(docId);

  return (
    <div className="px-8 py-8 max-w-3xl">
      {/* Breadcrumb */}
      <Link
        href={`/workspace/${workspaceId}`}
        className="inline-flex items-center gap-1 text-sm text-zinc-500 hover:text-zinc-900 mb-6 transition-colors"
      >
        <ChevronLeft className="h-4 w-4" />
        Back to workspace
      </Link>

      {isLoading ? (
        <div className="space-y-4">
          <Skeleton className="h-10 w-2/3" />
          <Skeleton className="h-4 w-1/3" />
          <div className="mt-8 space-y-3">
            <Skeleton className="h-4 w-full" />
            <Skeleton className="h-4 w-5/6" />
            <Skeleton className="h-4 w-4/6" />
          </div>
        </div>
      ) : isError ? (
        <div className="rounded-lg border border-red-200 bg-red-50 px-4 py-3">
          <p className="text-sm text-red-600">Failed to load document.</p>
        </div>
      ) : doc ? (
        <div>
          {/* Doc header */}
          <div className="mb-8">
            <div className="flex items-center gap-3 mb-2">
              <span className="text-4xl">{doc.emoji || "📄"}</span>
            </div>
            <h1 className="text-3xl font-bold text-zinc-900 tracking-tight mt-3">
              {doc.title}
            </h1>
            <div className="flex items-center gap-4 mt-3">
              <span className="flex items-center gap-1.5 text-xs text-zinc-400">
                <Clock className="h-3.5 w-3.5" />
                {formatDate(doc.createdAt)}
              </span>
            </div>
          </div>

          {/* Placeholder editor area */}
          <div className="rounded-xl border border-zinc-200 bg-white min-h-[400px] p-6">
            <div className="flex flex-col items-center justify-center h-40 text-center">
              <p className="text-sm font-medium text-zinc-500">
                Real-time collaborative editor
              </p>
              <p className="text-xs text-zinc-400 mt-1">
                Yjs + WebSocket integration coming soon
              </p>
              <div className="mt-4 flex gap-2">
                <span className="inline-flex items-center rounded-full bg-zinc-100 px-2.5 py-0.5 text-xs text-zinc-500">
                  CRDT
                </span>
                <span className="inline-flex items-center rounded-full bg-zinc-100 px-2.5 py-0.5 text-xs text-zinc-500">
                  Offline-first
                </span>
                <span className="inline-flex items-center rounded-full bg-zinc-100 px-2.5 py-0.5 text-xs text-zinc-500">
                  Conflict-free
                </span>
              </div>
            </div>
          </div>
        </div>
      ) : null}
    </div>
  );
}
