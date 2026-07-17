"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import Link from "next/link";
import { useRouter } from "next/navigation";

import { useRegister } from "@/hooks/use-auth";
import { useAuth } from "@/components/providers/auth-provider";
import { getApiError } from "@/lib/api/client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { FormField } from "@/components/ui/form-field";

const registerSchema = z.object({
  fullName: z.string().min(1, "Full name is required"),
  email: z.string().email("Invalid email address"),
  password: z
    .string()
    .min(6, "Password must be at least 6 characters"),
});

type RegisterForm = z.infer<typeof registerSchema>;

export function RegisterForm() {
  const router = useRouter();
  const { login: authLogin } = useAuth();
  const registerMutation = useRegister();

  const {
    register,
    handleSubmit,
    setError,
    formState: { errors },
  } = useForm<RegisterForm>({
    resolver: zodResolver(registerSchema),
  });

  const onSubmit = async (data: RegisterForm) => {
    try {
      const res = await registerMutation.mutateAsync(data);
      authLogin(res.data.accessToken, res.data.refreshToken);
      router.push("/dashboard");
    } catch (err) {
      setError("root", { message: getApiError(err) });
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4">
      <FormField label="Full name" error={errors.fullName?.message} required>
        <Input
          type="text"
          placeholder="John Doe"
          autoComplete="name"
          {...register("fullName")}
        />
      </FormField>

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
          autoComplete="new-password"
          {...register("password")}
        />
      </FormField>

      {errors.root && (
        <p className="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">
          {errors.root.message}
        </p>
      )}

      <Button
        type="submit"
        isLoading={registerMutation.isPending}
        className="mt-1 w-full"
        size="lg"
      >
        Create account
      </Button>

      <p className="text-center text-sm text-slate-500">
        Already have an account?{" "}
        <Link
          href="/login"
          className="font-medium text-slate-950 hover:text-blue-700"
        >
          Sign in
        </Link>
      </p>
    </form>
  );
}
