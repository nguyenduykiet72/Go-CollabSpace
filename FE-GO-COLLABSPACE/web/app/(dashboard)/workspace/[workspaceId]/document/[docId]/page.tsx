"use client";

import { use } from "react";
import Link from "next/link";
import { useDocDetail } from "@/hooks/use-document";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
import { formatDate } from "@/lib/utils";
import {
  Bold,
  ChevronLeft,
  Clock,
  Italic,
  List,
  MessageSquareText,
  Share2,
} from "lucide-react";

interface Props {
  params: Promise<{ workspaceId: string; docId: string }>;
}

export default function DocumentPage({ params }: Props) {
  const { workspaceId, docId } = use(params);
  const { data: doc, isLoading, isError } = useDocDetail(docId);

  return (
    <div className="mx-auto max-w-6xl px-8 py-8">
      <Link
        href={`/workspace/${workspaceId}`}
        className="mb-6 inline-flex items-center gap-1 text-sm text-stone-500 transition-colors hover:text-teal-800"
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
          <p className="text-sm text-red-700">Failed to load document.</p>
        </div>
      ) : doc ? (
        <div className="grid gap-6 lg:grid-cols-[1fr_280px]">
          <section className="min-w-0">
            <div className="mb-5 rounded-lg border border-stone-200 bg-white shadow-sm shadow-stone-200/60">
              <div className="flex flex-wrap items-center justify-between gap-3 border-b border-stone-100 bg-[#fffaf0] px-4 py-3">
                <div className="flex items-center gap-2">
                  {[
                    { icon: Bold, label: "Bold" },
                    { icon: Italic, label: "Italic" },
                    { icon: List, label: "List" },
                  ].map((item) => (
                    <button
                      key={item.label}
                      className="flex h-8 w-8 items-center justify-center rounded-md text-stone-500 transition-colors hover:bg-amber-100 hover:text-stone-950"
                      aria-label={item.label}
                    >
                      <item.icon className="h-4 w-4" />
                    </button>
                  ))}
                </div>
                <div className="flex items-center gap-2">
                  <Badge variant="warning">Editor pending</Badge>
                  <button className="inline-flex h-8 items-center gap-2 rounded-md border border-stone-300 bg-white/80 px-3 text-sm font-medium text-stone-700 hover:bg-amber-50">
                    <Share2 className="h-4 w-4" />
                    Share
                  </button>
                </div>
              </div>

              <div className="px-8 py-8">
                <div className="mb-6">
                  <span className="text-4xl">{doc.emoji || "📄"}</span>
                  <h1 className="mt-4 text-4xl font-semibold tracking-tight text-stone-950">
                    {doc.title}
                  </h1>
                  <div className="mt-3 flex flex-wrap items-center gap-4">
                    <span className="flex items-center gap-1.5 text-xs text-stone-500">
                      <Clock className="h-3.5 w-3.5" />
                      Created {formatDate(doc.createdAt)}
                    </span>
                  </div>
                </div>

                <div className="min-h-[440px] rounded-lg border border-stone-200 bg-[#fffdf8] px-8 py-7">
                  <div className="mx-auto max-w-2xl space-y-4">
                    <div className="h-4 w-10/12 rounded bg-stone-100" />
                    <div className="h-4 w-11/12 rounded bg-stone-100" />
                    <div className="h-4 w-7/12 rounded bg-stone-100" />
                    <div className="py-4">
                      <div className="rounded-md border border-dashed border-teal-300 bg-teal-50/60 px-4 py-4">
                        <p className="text-sm font-medium text-teal-900">
                          Real-time collaborative editor
                        </p>
                        <p className="mt-1 text-sm leading-6 text-stone-600">
                          The document shell is ready. Connect the Yjs client
                          and WebSocket provider here when the frontend moves to
                          the next implementation pass.
                        </p>
                      </div>
                    </div>
                    <div className="h-4 w-9/12 rounded bg-stone-100" />
                    <div className="h-4 w-8/12 rounded bg-stone-100" />
                  </div>
                </div>
              </div>
            </div>
          </section>

          <aside className="space-y-4">
            <div className="rounded-lg border border-teal-200 bg-teal-50/60 p-5 shadow-sm shadow-stone-200/60">
              <div className="flex items-center gap-2">
                <MessageSquareText className="h-4 w-4 text-teal-700" />
                <h2 className="text-sm font-semibold text-stone-950">
                  Session
                </h2>
              </div>
              <div className="mt-4 space-y-3">
                <div className="rounded-md border border-teal-100 bg-white/75 px-3 py-2">
                  <p className="text-xs font-medium uppercase text-teal-700">
                    Sync
                  </p>
                  <p className="mt-1 text-sm font-medium text-stone-800">
                    WebSocket not connected
                  </p>
                </div>
                <div className="rounded-md border border-teal-100 bg-white/75 px-3 py-2">
                  <p className="text-xs font-medium uppercase text-teal-700">
                    Document ID
                  </p>
                  <p className="mt-1 truncate text-sm font-medium text-stone-800">
                    {docId}
                  </p>
                </div>
              </div>
            </div>
            <div className="rounded-lg border border-rose-200 bg-rose-50/60 p-5 shadow-sm shadow-stone-200/60">
              <h2 className="text-sm font-semibold text-stone-950">
                Capabilities
              </h2>
              <div className="mt-3 flex flex-wrap gap-2">
                <Badge variant="muted">CRDT</Badge>
                <Badge variant="muted">Yjs</Badge>
                <Badge variant="muted">Role aware</Badge>
              </div>
            </div>
          </aside>
        </div>
      ) : null}
    </div>
  );
}
