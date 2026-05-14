"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { LayoutDashboard, LogOut, Settings, ChevronRight } from "lucide-react";
import { useAuth } from "@/components/providers/auth-provider";
import { Avatar } from "@/components/ui/avatar";
import { cn } from "@/lib/utils";

interface NavItem {
  label: string;
  href: string;
  icon: React.ElementType;
}

const NAV: NavItem[] = [
  { label: "Dashboard", href: "/dashboard", icon: LayoutDashboard },
];

export function Sidebar() {
  const pathname = usePathname();
  const { user, logout } = useAuth();

  return (
    <aside className="flex h-screen w-56 flex-col border-r border-zinc-200 bg-white">
      {/* Brand */}
      <div className="flex items-center gap-2 px-5 py-4 border-b border-zinc-100">
        <div className="h-6 w-6 rounded-md bg-zinc-900 flex items-center justify-center">
          <span className="text-white text-xs font-bold">G</span>
        </div>
        <span className="text-sm font-semibold text-zinc-900">GoCollab</span>
      </div>

      {/* Nav */}
      <nav className="flex-1 px-3 py-3 space-y-0.5 overflow-y-auto">
        {NAV.map((item) => (
          <Link
            key={item.href}
            href={item.href}
            className={cn(
              "flex items-center gap-2.5 px-3 py-2 rounded-md text-sm transition-colors",
              pathname === item.href
                ? "bg-zinc-100 text-zinc-900 font-medium"
                : "text-zinc-600 hover:bg-zinc-50 hover:text-zinc-900"
            )}
          >
            <item.icon className="h-4 w-4 shrink-0" />
            {item.label}
          </Link>
        ))}
      </nav>

      {/* User */}
      <div className="border-t border-zinc-100 px-3 py-3">
        {user && (
          <div className="flex items-center gap-3 px-2 py-2">
            <Avatar name={user.fullName} src={user.avatar} size="sm" />
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-zinc-900 truncate">
                {user.fullName}
              </p>
              <p className="text-xs text-zinc-500 truncate">{user.email}</p>
            </div>
          </div>
        )}
        <button
          onClick={logout}
          className="mt-1 flex w-full items-center gap-2.5 px-3 py-2 rounded-md text-sm text-zinc-600 hover:bg-zinc-50 hover:text-zinc-900 transition-colors"
        >
          <LogOut className="h-4 w-4" />
          Sign out
        </button>
      </div>
    </aside>
  );
}
