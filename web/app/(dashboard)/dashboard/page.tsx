"use client";

import { CreateWorkspaceModal } from "@/components/workspace/create-workspace-modal";
import { LayoutDashboard } from "lucide-react";

export default function DashboardPage() {
  return (
    <div className="px-8 py-8 max-w-2xl">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-zinc-900 tracking-tight">
          Dashboard
        </h1>
        <p className="mt-1 text-sm text-zinc-500">
          Create or open a workspace to get started.
        </p>
      </div>

      <div className="rounded-xl border border-dashed border-zinc-300 bg-white p-10 flex flex-col items-center justify-center text-center gap-4">
        <div className="h-12 w-12 rounded-full bg-zinc-100 flex items-center justify-center">
          <LayoutDashboard className="h-6 w-6 text-zinc-400" />
        </div>
        <div>
          <p className="text-sm font-medium text-zinc-900">
            No workspace open
          </p>
          <p className="text-xs text-zinc-500 mt-1">
            Create a new workspace or navigate to one using its ID.
          </p>
        </div>
        <CreateWorkspaceModal />
      </div>
    </div>
  );
}
