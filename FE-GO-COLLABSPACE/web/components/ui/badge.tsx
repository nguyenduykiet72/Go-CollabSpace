import { type ReactNode } from "react";
import { cn } from "@/lib/utils";

interface BadgeProps {
  children: ReactNode;
  variant?: "default" | "muted" | "success" | "warning";
  className?: string;
}

export function Badge({
  children,
  variant = "default",
  className,
}: BadgeProps) {
  const variants = {
    default: "bg-slate-950 text-white",
    muted: "bg-slate-100 text-slate-600 ring-1 ring-inset ring-slate-200",
    success: "bg-teal-50 text-teal-700 ring-1 ring-inset ring-teal-200",
    warning: "bg-amber-50 text-amber-700 ring-1 ring-inset ring-amber-200",
  };

  return (
    <span
      className={cn(
        "inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium",
        variants[variant],
        className
      )}
    >
      {children}
    </span>
  );
}
