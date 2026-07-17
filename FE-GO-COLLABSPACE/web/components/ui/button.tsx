import { forwardRef, type ButtonHTMLAttributes } from "react";
import { cn } from "@/lib/utils";
import { Loader2 } from "lucide-react";

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "default" | "outline" | "ghost" | "destructive";
  size?: "sm" | "md" | "lg";
  isLoading?: boolean;
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  (
    {
      className,
      variant = "default",
      size = "md",
      isLoading,
      disabled,
      children,
      ...props
    },
    ref
  ) => {
    const base =
      "inline-flex items-center justify-center gap-2 rounded-md font-medium tracking-normal transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50";

    const variants = {
      default: "bg-slate-950 text-white shadow-sm hover:bg-slate-800",
      outline:
        "border border-slate-300 bg-white text-slate-900 shadow-sm hover:bg-slate-50",
      ghost: "text-slate-700 hover:bg-slate-100 hover:text-slate-950",
      destructive: "bg-red-600 text-white hover:bg-red-700",
    };

    const sizes = {
      sm: "text-sm px-3 py-1.5 h-8",
      md: "text-sm px-4 py-2 h-9",
      lg: "text-base px-5 py-2.5 h-11",
    };

    return (
      <button
        ref={ref}
        className={cn(base, variants[variant], sizes[size], className)}
        disabled={disabled || isLoading}
        {...props}
      >
        {isLoading && <Loader2 className="h-4 w-4 animate-spin" />}
        {children}
      </button>
    );
  }
);

Button.displayName = "Button";
