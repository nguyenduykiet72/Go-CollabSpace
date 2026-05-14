"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useRouter } from "next/navigation";
import { Plus } from "lucide-react";

import { useCreateDocument } from "@/hooks/use-document";
import { getApiError } from "@/lib/api/client";
import { Modal } from "@/components/ui/modal";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { FormField } from "@/components/ui/form-field";

const EMOJIS = ["📄", "📝", "📊", "🗂️", "💡", "🚀", "✅", "🔖"];

const schema = z.object({
  title: z
    .string()
    .min(1, "Title is required")
    .max(200, "Max 200 characters"),
  emoji: z.string().optional(),
});

type FormValues = z.infer<typeof schema>;

interface CreateDocumentModalProps {
  workspaceId: string;
}

export function CreateDocumentModal({ workspaceId }: CreateDocumentModalProps) {
  const [open, setOpen] = useState(false);
  const [selectedEmoji, setSelectedEmoji] = useState("📄");
  const router = useRouter();
  const mutation = useCreateDocument();

  const {
    register,
    handleSubmit,
    reset,
    setError,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { emoji: "📄" },
  });

  const onSubmit = async (data: FormValues) => {
    try {
      const res = await mutation.mutateAsync({
        ...data,
        emoji: selectedEmoji,
        workspaceId,
      });
      reset();
      setOpen(false);
      router.push(`/workspace/${workspaceId}/document/${res.data.id}`);
    } catch (err) {
      setError("root", { message: getApiError(err) });
    }
  };

  return (
    <>
      <Button onClick={() => setOpen(true)} size="sm">
        <Plus className="h-4 w-4" />
        New document
      </Button>

      <Modal
        open={open}
        onClose={() => {
          setOpen(false);
          reset();
        }}
        title="New document"
        description="Create a new document in this workspace."
      >
        <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4">
          {/* Emoji picker */}
          <div>
            <p className="text-sm font-medium text-zinc-700 mb-2">Icon</p>
            <div className="flex gap-2 flex-wrap">
              {EMOJIS.map((e) => (
                <button
                  key={e}
                  type="button"
                  onClick={() => setSelectedEmoji(e)}
                  className={`h-9 w-9 rounded-md text-xl flex items-center justify-center transition-colors ${
                    selectedEmoji === e
                      ? "bg-zinc-900 ring-2 ring-zinc-900 ring-offset-1"
                      : "bg-zinc-100 hover:bg-zinc-200"
                  }`}
                >
                  {e}
                </button>
              ))}
            </div>
          </div>

          <FormField label="Title" error={errors.title?.message} required>
            <Input
              placeholder="Untitled document"
              autoFocus
              {...register("title")}
            />
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
              Create document
            </Button>
          </div>
        </form>
      </Modal>
    </>
  );
}
