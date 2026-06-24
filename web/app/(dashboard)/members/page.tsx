"use client";

import { Badge } from "@/components/ui/badge";
import { ShieldCheck, UserPlus, Users } from "lucide-react";

const ROLES = [
  {
    name: "Owner",
    detail: "Full workspace control",
    className: "border-teal-200 bg-teal-50",
  },
  {
    name: "Admin",
    detail: "Manage people and documents",
    className: "border-amber-200 bg-amber-50",
  },
  {
    name: "Editor",
    detail: "Create and update documents",
    className: "border-emerald-200 bg-emerald-50",
  },
  {
    name: "Viewer",
    detail: "Read-only workspace access",
    className: "border-rose-200 bg-rose-50",
  },
];

export default function MembersPage() {
  return (
    <div className="mx-auto max-w-6xl px-8 py-8">
      <div className="mb-8 flex flex-col gap-4 border-b border-stone-200 pb-6 md:flex-row md:items-end md:justify-between">
        <div>
          <p className="text-xs font-medium uppercase text-amber-700">
            Members
          </p>
          <h1 className="mt-2 text-3xl font-semibold tracking-tight text-stone-950">
            Workspace members
          </h1>
          <p className="mt-2 max-w-2xl text-sm leading-6 text-stone-600">
            Member management is workspace-specific. This page shows the access
            model while the member endpoints are connected to the frontend.
          </p>
        </div>
        <button
          disabled
          className="inline-flex h-9 cursor-not-allowed items-center gap-2 rounded-md border border-stone-300 bg-white px-3 text-sm font-medium text-stone-400"
        >
          <UserPlus className="h-4 w-4" />
          Invite member
        </button>
      </div>

      <div className="grid gap-4 md:grid-cols-4">
        {ROLES.map((role) => (
          <div
            key={role.name}
            className={`rounded-lg border p-4 shadow-sm shadow-stone-200/60 ${role.className}`}
          >
            <div className="flex items-center justify-between">
              <h2 className="text-sm font-semibold text-stone-950">
                {role.name}
              </h2>
              <ShieldCheck className="h-4 w-4 text-stone-500" />
            </div>
            <p className="mt-3 text-sm leading-5 text-stone-600">
              {role.detail}
            </p>
          </div>
        ))}
      </div>

      <section className="mt-6 rounded-lg border border-stone-200 bg-white p-5 shadow-sm shadow-stone-200/60">
        <div className="flex items-center gap-2">
          <Users className="h-4 w-4 text-amber-700" />
          <h2 className="text-sm font-semibold text-stone-950">
            Current state
          </h2>
          <Badge variant="warning">Frontend shell</Badge>
        </div>
        <p className="mt-3 max-w-3xl text-sm leading-6 text-stone-600">
          The backend README describes role-based workspace access. Once the
          member list endpoint is wired here, this page can show invites,
          current roles, and permission changes.
        </p>
      </section>
    </div>
  );
}
