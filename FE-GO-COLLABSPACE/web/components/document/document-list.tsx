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
      <div className="space-y-2 p-4">
        {[...Array(4)].map((_, i) => (
          <Skeleton key={i} className="h-14 w-full" />
        ))}
      </div>
    );
  }

  if (docs.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center px-6 py-16 text-center">
        <div className="mb-3 flex h-11 w-11 items-center justify-center rounded-lg bg-emerald-50">
          <FileText className="h-5 w-5 text-emerald-700" />
        </div>
        <p className="text-sm font-medium text-stone-700">No documents yet</p>
        <p className="mt-1 max-w-sm text-xs leading-5 text-stone-500">
          Create a document for lecture notes, planning, or shared exercises.
        </p>
      </div>
    );
  }

  return (
    <div className="divide-y divide-stone-100">
      {docs.map((doc) => (
        <Link
          key={doc.id}
          href={`/workspace/${workspaceId}/document/${doc.id}`}
          className="group flex items-center gap-3 px-5 py-4 transition-colors hover:bg-emerald-50/45"
        >
          <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-md bg-amber-50 text-base ring-1 ring-inset ring-amber-100">
            {doc.emoji || "📄"}
          </div>
          <div className="flex-1 min-w-0">
            <p className="truncate text-sm font-medium text-stone-950 group-hover:text-teal-800">
              {doc.title}
            </p>
            <p className="mt-0.5 text-xs text-stone-500">
              {formatDate(doc.createdAt)}
            </p>
          </div>
          <ChevronRight className="h-4 w-4 text-stone-300 transition-colors group-hover:text-teal-700" />
        </Link>
      ))}
    </div>
  );
}
