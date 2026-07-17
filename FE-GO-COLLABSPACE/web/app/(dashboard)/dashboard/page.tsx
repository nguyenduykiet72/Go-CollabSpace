"use client";

import { CreateWorkspaceModal } from "@/components/workspace/create-workspace-modal";
import {
  BookOpen,
  FileText,
  LayoutDashboard,
  LockKeyhole,
  MessageSquareText,
  Users,
} from "lucide-react";

export default function DashboardPage() {
  return (
    <div className="mx-auto max-w-6xl px-8 py-8">
      <div className="mb-8 flex flex-col gap-4 border-b border-stone-200 pb-6 md:flex-row md:items-end md:justify-between">
        <div>
          <p className="text-xs font-medium uppercase text-teal-700">
            Workspace hub
          </p>
          <h1 className="mt-2 text-3xl font-semibold tracking-tight text-stone-950">
            Dashboard
          </h1>
          <p className="mt-2 max-w-2xl text-sm leading-6 text-stone-600">
            Create a shared space for notes, study plans, and live documents.
          </p>
        </div>
        <CreateWorkspaceModal />
      </div>

      <div className="grid gap-4 md:grid-cols-3">
        {[
          {
            label: "Live documents",
            value: "Yjs",
            icon: FileText,
            className: "border-emerald-200 bg-emerald-50/80",
            iconClassName: "text-emerald-700",
          },
          {
            label: "Role controls",
            value: "RBAC",
            icon: LockKeyhole,
            className: "border-amber-200 bg-amber-50/80",
            iconClassName: "text-amber-700",
          },
          {
            label: "Study rooms",
            value: "Shared",
            icon: Users,
            className: "border-rose-200 bg-rose-50/80",
            iconClassName: "text-rose-700",
          },
        ].map((item) => (
          <div
            key={item.label}
            className={`rounded-lg border p-4 shadow-sm shadow-stone-200/60 ${item.className}`}
          >
            <div className="flex items-center justify-between">
              <p className="text-sm font-medium text-stone-700">{item.label}</p>
              <item.icon className={`h-4 w-4 ${item.iconClassName}`} />
            </div>
            <p className="mt-4 text-2xl font-semibold text-stone-950">
              {item.value}
            </p>
          </div>
        ))}
      </div>

      <div className="mt-6 grid gap-6 lg:grid-cols-[1.15fr_0.85fr]">
        <div className="rounded-lg border border-dashed border-teal-300 bg-white p-8 shadow-sm shadow-stone-200/70">
          <div className="flex max-w-xl flex-col gap-4">
            <div className="flex h-11 w-11 items-center justify-center rounded-lg bg-teal-50">
              <LayoutDashboard className="h-5 w-5 text-teal-700" />
            </div>
            <div>
              <h2 className="text-lg font-semibold text-stone-950">
                No workspace selected
              </h2>
              <p className="mt-2 text-sm leading-6 text-stone-600">
                Start with a workspace for a class, cohort, or research group.
                Documents created inside it can later sync through the Go
                backend over WebSocket.
              </p>
            </div>
            <CreateWorkspaceModal />
          </div>
        </div>

        <div className="rounded-lg border border-stone-200 bg-[#fffaf0] p-5 shadow-sm shadow-stone-200/60">
          <div className="flex items-center gap-2">
            <BookOpen className="h-4 w-4 text-amber-700" />
            <h2 className="text-sm font-semibold text-stone-950">
              Good workspace names
            </h2>
          </div>
          <div className="mt-4 space-y-3">
            {[
              "Distributed Systems",
              "IELTS Speaking Circle",
              "ML Paper Club",
            ].map((name) => (
              <div
                key={name}
                className="rounded-md border border-amber-100 bg-white/70 px-3 py-2"
              >
                <p className="text-sm font-medium text-stone-800">{name}</p>
                <p className="text-xs text-stone-500">
                  {name.toLowerCase().replace(/[^a-z0-9]/g, "")}
                </p>
              </div>
            ))}
          </div>
        </div>
      </div>

      <div className="mt-6 rounded-lg border border-indigo-100 bg-indigo-50/70 p-5">
        <div className="flex items-center gap-2">
          <MessageSquareText className="h-4 w-4 text-indigo-700" />
          <h2 className="text-sm font-semibold text-stone-950">
            Collaboration rhythm
          </h2>
        </div>
        <p className="mt-2 max-w-3xl text-sm leading-6 text-stone-600">
          Use one workspace per learning group, then split documents by topic,
          meeting, or assignment so shared notes stay easy to scan.
        </p>
      </div>
    </div>
  );
}
