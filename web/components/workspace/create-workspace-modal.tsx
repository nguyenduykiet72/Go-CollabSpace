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
    watch,
    setValue,
    setError,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
  });

  // Auto-generate slug from name
  const nameValue = watch("name", "");

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
            <div className="flex items-center gap-0 rounded-md border border-zinc-300 focus-within:ring-2 focus-within:ring-zinc-900 overflow-hidden">
              <span className="px-3 py-2 text-sm text-zinc-400 bg-zinc-50 border-r border-zinc-200 shrink-0">
                /workspace/
              </span>
              <input
                className="flex-1 px-3 py-2 text-sm outline-none bg-white text-zinc-900"
                placeholder="my-team"
                {...register("slug")}
              />
            </div>
            {errors.slug && (
              <p className="text-xs text-red-500">{errors.slug.message}</p>
            )}
          </FormField>

          {errors.root && (
            <p className="text-sm text-red-500 bg-red-50 border border-red-200 rounded-md px-3 py-2">
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
