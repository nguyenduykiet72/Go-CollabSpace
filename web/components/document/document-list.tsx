"use client";

import { formatDate } from "@/lib/utils";
import type { DocumentResponse } from "@/lib/types";
import { FileText, ChevronRight } from "lucide-react";
import Link from "next/link";
import { Skeleton } from "@/components/ui/skeleton";

interface DocumentListProps {
  docs: DocumentResponse[];
  workspaceId: string;
  isLoading?: boolean;
}

export function DocumentList({
  docs,
  workspaceId,
  isLoading,
}: DocumentListProps) {
  if (isLoading) {
    return (
      <div className="space-y-2">
        {[...Array(4)].map((_, i) => (
          <Skeleton key={i} className="h-14 w-full" />
        ))}
      </div>
    );
  }

  if (docs.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-16 text-center">
        <FileText className="h-10 w-10 text-zinc-300 mb-3" />
        <p className="text-sm font-medium text-zinc-600">No documents yet</p>
        <p className="text-xs text-zinc-400 mt-1">
          Create your first document to get started
        </p>
      </div>
    );
  }

  return (
    <div className="divide-y divide-zinc-100">
      {docs.map((doc) => (
        <Link
          key={doc.id}
          href={`/workspace/${workspaceId}/document/${doc.id}`}
          className="flex items-center gap-3 px-4 py-3 hover:bg-zinc-50 transition-colors group"
        >
          <div className="flex h-8 w-8 items-center justify-center rounded-md bg-zinc-100 text-base shrink-0">
            {doc.emoji || "📄"}
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-sm font-medium text-zinc-900 truncate group-hover:text-zinc-700">
              {doc.title}
            </p>
            <p className="text-xs text-zinc-400 mt-0.5">
              {formatDate(doc.createdAt)}
            </p>
          </div>
          <ChevronRight className="h-4 w-4 text-zinc-300 group-hover:text-zinc-500 transition-colors" />
        </Link>
      ))}
    </div>
  );
}
