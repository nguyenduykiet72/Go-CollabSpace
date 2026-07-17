"use client";

import { type ReactNode } from "react";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { Sidebar } from "@/components/layout/sidebar";
import { isAuthenticated } from "@/lib/auth";

export default function DashboardLayout({ children }: { children: ReactNode }) {
  const router = useRouter();
  const [canRender] = useState(() => isAuthenticated());

  useEffect(() => {
    if (!canRender) {
      router.replace("/login");
    }
  }, [canRender, router]);

  if (!canRender) {
    return (
      <div className="flex h-screen items-center justify-center bg-slate-50">
        <div className="h-2 w-32 overflow-hidden rounded-full bg-slate-200">
          <div className="h-full w-1/2 rounded-full bg-slate-400" />
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-screen overflow-hidden bg-slate-50">
      <Sidebar />
      <main className="flex-1 overflow-y-auto">{children}</main>
    </div>
  );
}
