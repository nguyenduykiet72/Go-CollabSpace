"use client";

import { CreateWorkspaceModal } from "@/components/workspace/create-workspace-modal";
import { BookOpen, FileText, FolderOpen, Search } from "lucide-react";

export default function DocumentsPage() {
  return (
    <div className="mx-auto max-w-6xl px-8 py-8">
      <div className="mb-8 border-b border-stone-200 pb-6">
        <p className="text-xs font-medium uppercase text-emerald-700">
          Documents
        </p>
        <h1 className="mt-2 text-3xl font-semibold tracking-tight text-stone-950">
          Document library
        </h1>
        <p className="mt-2 max-w-2xl text-sm leading-6 text-stone-600">
          Documents belong to a workspace. Open or create a workspace to browse,
          create, and edit shared notes.
        </p>
      </div>

      <div className="grid gap-6 lg:grid-cols-[1fr_320px]">
        <section className="rounded-lg border border-emerald-200 bg-emerald-50/70 p-8 shadow-sm shadow-stone-200/70">
          <div className="flex max-w-xl flex-col gap-4">
            <div className="flex h-11 w-11 items-center justify-center rounded-lg bg-white/75 ring-1 ring-emerald-100">
              <FolderOpen className="h-5 w-5 text-emerald-700" />
            </div>
            <div>
              <h2 className="text-lg font-semibold text-stone-950">
                Open a workspace first
              </h2>
              <p className="mt-2 text-sm leading-6 text-stone-600">
                The backend organizes documents by workspace ID, so the document
                list is available inside each workspace page.
              </p>
            </div>
            <CreateWorkspaceModal />
          </div>
        </section>

        <aside className="rounded-lg border border-stone-200 bg-white p-5 shadow-sm shadow-stone-200/60">
          <div className="flex items-center gap-2">
            <BookOpen className="h-4 w-4 text-amber-700" />
            <h2 className="text-sm font-semibold text-stone-950">
              Library actions
            </h2>
          </div>
          <div className="mt-4 space-y-3">
            {[
              { icon: Search, label: "Search across workspace documents" },
              { icon: FileText, label: "Create notes inside a workspace" },
              { icon: FolderOpen, label: "Keep documents grouped by topic" },
            ].map((item) => (
              <div
                key={item.label}
                className="flex items-center gap-3 rounded-md border border-stone-100 bg-stone-50 px-3 py-2"
              >
                <item.icon className="h-4 w-4 text-stone-500" />
                <p className="text-sm text-stone-700">{item.label}</p>
              </div>
            ))}
          </div>
        </aside>
      </div>
    </div>
  );
}
