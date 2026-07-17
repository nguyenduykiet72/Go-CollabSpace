import { type ReactNode } from "react";
import { cn } from "@/lib/utils";

interface CardProps {
  className?: string;
  children: ReactNode;
}

export function Card({ className, children }: CardProps) {
  return (
    <div
      className={cn(
        "rounded-lg border border-slate-200 bg-white shadow-sm shadow-slate-200/50",
        className
      )}
    >
      {children}
    </div>
  );
}

export function CardHeader({ className, children }: CardProps) {
  return (
    <div className={cn("border-b border-slate-100 px-6 py-4", className)}>
      {children}
    </div>
  );
}

export function CardTitle({ className, children }: CardProps) {
  return (
    <h2 className={cn("text-lg font-semibold text-slate-950", className)}>
      {children}
    </h2>
  );
}

export function CardDescription({ className, children }: CardProps) {
  return (
    <p className={cn("mt-0.5 text-sm text-slate-500", className)}>{children}</p>
  );
}

export function CardBody({ className, children }: CardProps) {
  return <div className={cn("px-6 py-4", className)}>{children}</div>;
}

export function CardFooter({ className, children }: CardProps) {
  return (
    <div
      className={cn(
        "rounded-b-lg border-t border-slate-100 bg-slate-50/80 px-6 py-4",
        className
      )}
    >
      {children}
    </div>
  );
}
