import { LoginForm } from "@/components/auth/login-form";

export const metadata = {
  title: "Sign in — GoCollab",
};

export default function LoginPage() {
  return (
    <div className="min-h-screen bg-zinc-50 flex flex-col items-center justify-center px-4">
      <div className="w-full max-w-sm">
        <div className="mb-8 text-center">
          <h1 className="text-2xl font-bold tracking-tight text-zinc-900">
            GoCollab
          </h1>
          <p className="mt-2 text-sm text-zinc-500">
            Sign in to your workspace
          </p>
        </div>
        <div className="bg-white border border-zinc-200 rounded-xl p-6 shadow-sm">
          <LoginForm />
        </div>
      </div>
    </div>
  );
}
