import { cn } from "@/lib/utils";
import { getInitials } from "@/lib/utils";

interface AvatarProps {
  name: string;
  src?: string;
  size?: "sm" | "md" | "lg";
  className?: string;
}

export function Avatar({ name, src, size = "md", className }: AvatarProps) {
  const sizes = {
    sm: "h-7 w-7 text-xs",
    md: "h-8 w-8 text-sm",
    lg: "h-10 w-10 text-base",
  };

  if (src) {
    return (
      // eslint-disable-next-line @next/next/no-img-element
      <img
        src={src}
        alt={name}
        className={cn("rounded-full object-cover ring-1 ring-slate-200", sizes[size], className)}
      />
    );
  }

  return (
    <div
      className={cn(
        "flex shrink-0 items-center justify-center rounded-full bg-slate-900 font-medium text-white ring-1 ring-slate-200",
        sizes[size],
        className
      )}
    >
      {getInitials(name)}
    </div>
  );
}
