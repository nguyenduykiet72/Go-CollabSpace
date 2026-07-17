"use client";

import { Badge } from "@/components/ui/badge";
import { Bell, Database, KeyRound, Settings } from "lucide-react";

const SETTINGS = [
  {
    label: "Workspace profile",
    detail: "Name, slug, and default document behavior",
    icon: Settings,
    className: "border-rose-200 bg-rose-50",
  },
  {
    label: "Authentication",
    detail: "JWT access, refresh token rotation, and OAuth settings",
    icon: KeyRound,
    className: "border-amber-200 bg-amber-50",
  },
  {
    label: "Storage",
    detail: "S3 uploads and presigned file access",
    icon: Database,
    className: "border-emerald-200 bg-emerald-50",
  },
  {
    label: "Notifications",
    detail: "Email and async worker delivery preferences",
    icon: Bell,
    className: "border-indigo-200 bg-indigo-50",
  },
];

export default function SettingsPage() {
  return (
    <div className="mx-auto max-w-6xl px-8 py-8">
      <div className="mb-8 border-b border-stone-200 pb-6">
        <p className="text-xs font-medium uppercase text-rose-700">
          Settings
        </p>
        <h1 className="mt-2 text-3xl font-semibold tracking-tight text-stone-950">
          Settings
        </h1>
        <p className="mt-2 max-w-2xl text-sm leading-6 text-stone-600">
          Settings are shown as frontend placeholders until workspace-specific
          update endpoints are connected.
        </p>
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        {SETTINGS.map((item) => (
          <section
            key={item.label}
            className={`rounded-lg border p-5 shadow-sm shadow-stone-200/60 ${item.className}`}
          >
            <div className="flex items-center gap-3">
              <div className="flex h-9 w-9 items-center justify-center rounded-md bg-white/75 ring-1 ring-inset ring-white/70">
                <item.icon className="h-4 w-4 text-stone-700" />
              </div>
              <div>
                <h2 className="text-sm font-semibold text-stone-950">
                  {item.label}
                </h2>
                <p className="mt-1 text-sm leading-5 text-stone-600">
                  {item.detail}
                </p>
              </div>
            </div>
          </section>
        ))}
      </div>

      <section className="mt-6 rounded-lg border border-stone-200 bg-white p-5 shadow-sm shadow-stone-200/60">
        <div className="flex items-center gap-2">
          <Settings className="h-4 w-4 text-rose-700" />
          <h2 className="text-sm font-semibold text-stone-950">
            Configuration status
          </h2>
          <Badge variant="muted">Read-only</Badge>
        </div>
        <p className="mt-3 max-w-3xl text-sm leading-6 text-stone-600">
          This keeps the sidebar complete without pretending settings are
          already editable. The controls can be enabled as soon as the frontend
          has the matching API contracts.
        </p>
      </section>
    </div>
  );
}
