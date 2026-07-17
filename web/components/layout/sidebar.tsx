"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  Files,
  LayoutDashboard,
  LogOut,
  Search,
  Settings,
  Sparkles,
  Users,
} from "lucide-react";
import { useAuth } from "@/components/providers/auth-provider";
import { Avatar } from "@/components/ui/avatar";
import { cn } from "@/lib/utils";

interface NavItem {
  label: string;
  href: string;
  icon: React.ElementType;
  color?: string;
}

const NAV: NavItem[] = [
  { label: "Dashboard", href: "/dashboard", icon: LayoutDashboard },
  {
    label: "Documents",
    href: "/documents",
    icon: Files,
    color: "text-emerald-700",
  },
  { label: "Members", href: "/members", icon: Users, color: "text-amber-700" },
  {
    label: "Settings",
    href: "/settings",
    icon: Settings,
    color: "text-rose-700",
  },
];

export function Sidebar() {
  const pathname = usePathname();
  const { user, logout } = useAuth();

  return (
    <aside className="flex h-screen w-64 shrink-0 flex-col border-r border-stone-200 bg-[#fbfaf7]">
      <div className="border-b border-stone-200 px-5 py-4">
        <div className="flex items-center gap-3">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-[#24322f] text-sm font-semibold text-amber-100">
            G
          </div>
          <div className="min-w-0">
            <p className="text-sm font-semibold text-stone-950">GoCollab</p>
            <p className="text-xs text-stone-500">Shared learning spaces</p>
          </div>
        </div>
      </div>

      <div className="px-3 py-3">
        <div className="flex h-9 w-full items-center gap-2 rounded-md border border-stone-300 bg-white px-3 shadow-sm shadow-stone-200/60 focus-within:border-teal-600 focus-within:ring-2 focus-within:ring-teal-100">
          <Search className="h-4 w-4 text-stone-400" />
          <input
            type="search"
            placeholder="Search documents"
            className="min-w-0 flex-1 bg-transparent text-sm text-stone-900 outline-none placeholder:text-stone-400"
          />
        </div>
      </div>

      <nav className="flex-1 space-y-1 overflow-y-auto px-3 pb-3">
        {NAV.map((item) => (
          <Link
            key={item.href}
            href={item.href}
            className={cn(
              "flex items-center gap-2.5 rounded-md px-3 py-2 text-sm transition-colors",
              pathname === item.href
                ? "bg-teal-50 font-medium text-teal-800 ring-1 ring-inset ring-teal-100"
                : "text-stone-600 hover:bg-stone-100 hover:text-stone-950"
            )}
          >
            <item.icon className={cn("h-4 w-4 shrink-0", item.color)} />
            {item.label}
          </Link>
        ))}

        <div className="mx-3 mt-5 rounded-lg border border-amber-200 bg-amber-50 p-3">
          <div className="flex items-center gap-2">
            <Sparkles className="h-4 w-4 text-amber-700" />
            <p className="text-xs font-semibold text-amber-950">Study mode</p>
          </div>
          <p className="mt-2 text-xs leading-5 text-amber-800">
            These pages are frontend shells until workspace-specific endpoints
            are wired in.
          </p>
        </div>
      </nav>

      <div className="border-t border-stone-200 px-3 py-3">
        {user && (
          <div className="flex items-center gap-3 px-2 py-2">
            <Avatar name={user.fullName} src={user.avatar} size="sm" />
            <div className="flex-1 min-w-0">
              <p className="truncate text-sm font-medium text-stone-950">
                {user.fullName}
              </p>
              <p className="truncate text-xs text-stone-500">{user.email}</p>
            </div>
          </div>
        )}
        <button
          onClick={logout}
          className="mt-1 flex w-full items-center gap-2.5 rounded-md px-3 py-2 text-sm text-stone-600 transition-colors hover:bg-stone-100 hover:text-stone-950"
        >
          <LogOut className="h-4 w-4" />
          Sign out
        </button>
      </div>
    </aside>
  );
}
