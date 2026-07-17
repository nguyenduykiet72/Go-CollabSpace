"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import Link from "next/link";
import { useRouter } from "next/navigation";

import { useLogin } from "@/hooks/use-auth";
import { useAuth } from "@/components/providers/auth-provider";
import { getApiError } from "@/lib/api/client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { FormField } from "@/components/ui/form-field";
import { isMockLogin, MOCK_LOGIN, MOCK_USER } from "@/lib/mock-auth";

const loginSchema = z.object({
  email: z.string().email("Invalid email address"),
  password: z.string().min(1, "Password is required"),
});

type LoginForm = z.infer<typeof loginSchema>;

export function LoginForm() {
  const router = useRouter();
  const { login: authLogin } = useAuth();
  const loginMutation = useLogin();

  const {
    register,
    handleSubmit,
    setValue,
    setError,
    formState: { errors },
  } = useForm<LoginForm>({
    resolver: zodResolver(loginSchema),
  });

  const onSubmit = async (data: LoginForm) => {
    try {
      const res = await loginMutation.mutateAsync(data);
      authLogin(
        res.data.accessToken,
        res.data.refreshToken,
        isMockLogin(data) ? MOCK_USER : undefined
      );
      router.push("/dashboard");
    } catch (err) {
      setError("root", { message: getApiError(err) });
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4">
      <FormField label="Email" error={errors.email?.message} required>
        <Input
          type="email"
          placeholder="you@example.com"
          autoComplete="email"
          {...register("email")}
        />
      </FormField>

      <FormField label="Password" error={errors.password?.message} required>
        <Input
          type="password"
          placeholder="••••••••"
          autoComplete="current-password"
          {...register("password")}
        />
      </FormField>

      <div className="rounded-md border border-amber-200 bg-amber-50 px-3 py-2">
        <div className="flex items-start justify-between gap-3">
          <div>
            <p className="text-xs font-medium text-amber-900">
              Temporary demo login
            </p>
            <p className="mt-1 text-xs leading-5 text-amber-800">
              {MOCK_LOGIN.email} / {MOCK_LOGIN.password}
            </p>
          </div>
          <button
            type="button"
            onClick={() => {
              setValue("email", MOCK_LOGIN.email, { shouldValidate: true });
              setValue("password", MOCK_LOGIN.password, {
                shouldValidate: true,
              });
            }}
            className="shrink-0 text-xs font-medium text-amber-950 hover:text-slate-950"
          >
            Fill
          </button>
        </div>
      </div>

      {errors.root && (
        <p className="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">
          {errors.root.message}
        </p>
      )}

      <Button
        type="submit"
        isLoading={loginMutation.isPending}
        className="mt-1 w-full"
        size="lg"
      >
        Sign in
      </Button>

      <p className="text-center text-sm text-slate-500">
        Don&apos;t have an account?{" "}
        <Link
          href="/register"
          className="font-medium text-slate-950 hover:text-blue-700"
        >
          Create one
        </Link>
      </p>
    </form>
  );
}
