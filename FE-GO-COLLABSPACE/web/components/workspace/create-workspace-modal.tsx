"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useRouter } from "next/navigation";
import { Plus } from "lucide-react";

import { useCreateWorkspace } from "@/hooks/use-workspace";
import { getApiError } from "@/lib/api/client";
import { Modal } from "@/components/ui/modal";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { FormField } from "@/components/ui/form-field";

const schema = z.object({
  name: z
    .string()
    .min(3, "At least 3 characters")
    .max(100, "Max 100 characters"),
  slug: z
    .string()
    .min(3, "At least 3 characters")
    .max(50, "Max 50 characters")
    .regex(/^[a-zA-Z0-9]+$/, "Only letters and numbers"),
});

type FormValues = z.infer<typeof schema>;

export function CreateWorkspaceModal() {
  const [open, setOpen] = useState(false);
  const router = useRouter();
  const mutation = useCreateWorkspace();

  const {
    register,
    handleSubmit,
    reset,
    setValue,
    setError,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
  });

  function handleNameChange(e: React.ChangeEvent<HTMLInputElement>) {
    const v = e.target.value;
    setValue("name", v);
    setValue(
      "slug",
      v
        .toLowerCase()
        .replace(/[^a-z0-9]/g, "")
        .slice(0, 50)
    );
  }

  const onSubmit = async (data: FormValues) => {
    try {
      const res = await mutation.mutateAsync(data);
      reset();
      setOpen(false);
      router.push(`/workspace/${res.data.id}`);
    } catch (err) {
      setError("root", { message: getApiError(err) });
    }
  };

  return (
    <>
      <Button onClick={() => setOpen(true)} size="sm">
        <Plus className="h-4 w-4" />
        New workspace
      </Button>

      <Modal
        open={open}
        onClose={() => {
          setOpen(false);
          reset();
        }}
        title="Create workspace"
        description="Workspaces help you organise your team's documents."
      >
        <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4">
          <FormField label="Name" error={errors.name?.message} required>
            <Input
              placeholder="My team workspace"
              {...register("name")}
              onChange={handleNameChange}
            />
          </FormField>

          <FormField label="Slug (URL)" error={errors.slug?.message} required>
            <div className="flex items-center gap-0 overflow-hidden rounded-md border border-slate-300 bg-white shadow-sm focus-within:border-blue-500 focus-within:ring-2 focus-within:ring-blue-100">
              <span className="shrink-0 border-r border-slate-200 bg-slate-50 px-3 py-2 text-sm text-slate-500">
                /workspace/
              </span>
              <input
                className="min-w-0 flex-1 bg-white px-3 py-2 text-sm text-slate-950 outline-none placeholder:text-slate-400"
                placeholder="my-team"
                {...register("slug")}
              />
            </div>
            {errors.slug && (
              <p className="text-xs text-red-600">{errors.slug.message}</p>
            )}
          </FormField>

          {errors.root && (
            <p className="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">
              {errors.root.message}
            </p>
          )}

          <div className="flex justify-end gap-2 pt-1">
            <Button
              type="button"
              variant="outline"
              onClick={() => {
                setOpen(false);
                reset();
              }}
            >
              Cancel
            </Button>
            <Button type="submit" isLoading={mutation.isPending}>
              Create workspace
            </Button>
          </div>
        </form>
      </Modal>
    </>
  );
}
