import { RegisterForm } from "@/components/auth/register-form";
import { FileText, ShieldCheck, Users } from "lucide-react";

export const metadata = {
  title: "Create account — GoCollab",
};

export default function RegisterPage() {
  return (
    <div
      className="grid min-h-screen bg-cover bg-center lg:grid-cols-[1fr_460px]"
      style={{ backgroundImage: "url('/img-1.png')" }}
    >
      <div className="hidden bg-slate-950/25 px-10 py-10 text-white backdrop-blur-[1px] lg:flex lg:flex-col">
        <div className="flex items-center gap-3 drop-shadow-sm">
          <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-white/95 text-sm font-semibold text-slate-950">
            G
          </div>
          <div>
            <p className="text-sm font-semibold">GoCollab</p>
            <p className="text-xs text-white/75">Collaborative learning</p>
          </div>
        </div>

        <div className="flex flex-1 items-center">
          <div className="max-w-xl">
            <p className="text-xs font-medium uppercase text-amber-100">
              Start a learning space
            </p>
            <h1 className="mt-4 text-4xl font-semibold tracking-tight text-white drop-shadow-sm">
              Build a workspace where every note has a place.
            </h1>
            <p className="mt-4 text-sm leading-6 text-white/82 drop-shadow-sm">
              Create documents, organize them by group, and prepare for live
              collaborative editing through the Go backend.
            </p>
            <div className="mt-8 grid gap-4 select-none">
              {[
                { icon: FileText, label: "Nested document spaces" },
                { icon: ShieldCheck, label: "Workspace-level access" },
                { icon: Users, label: "Groups built for learning" },
              ].map((item) => (
                <div
                  key={item.label}
                  className="flex items-center gap-3 text-white/90"
                >
                  <span className="flex h-8 w-8 items-center justify-center rounded-full bg-white/18 ring-1 ring-white/25">
                    <item.icon className="h-4 w-4 text-amber-100" />
                  </span>
                  <span className="text-sm font-medium drop-shadow-sm">
                    {item.label}
                  </span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>

      <div className="flex items-center justify-center bg-white/88 px-4 py-10 shadow-2xl shadow-slate-950/10 backdrop-blur-md">
        <div className="w-full max-w-sm">
          <div className="mb-8 lg:hidden">
            <div className="flex items-center gap-3">
              <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-slate-950 text-sm font-semibold text-white">
                G
              </div>
              <div>
                <h1 className="text-base font-semibold text-slate-950">
                  GoCollab
                </h1>
                <p className="text-sm text-slate-500">Create your account</p>
              </div>
            </div>
          </div>
          <div className="mb-6">
            <h2 className="text-2xl font-semibold tracking-tight text-slate-950">
              Create account
            </h2>
            <p className="mt-2 text-sm text-slate-500">
              Set up your profile and create a workspace when you get in.
            </p>
          </div>
          <div className="rounded-lg border border-white/80 bg-white/90 p-6 shadow-xl shadow-slate-900/10">
            <RegisterForm />
          </div>
        </div>
      </div>
    </div>
  );
}
