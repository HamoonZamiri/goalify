import { z } from "zod";

export function createServerResponseSchema<TData extends z.ZodTypeAny>(
  schema: TData,
) {
  return z.object({
    message: z.string(),
    data: schema,
  });
}

export const UserSchema = z.object({
  id: z.string().uuid(),
  email: z.string(),
  xp: z.number(),
  level_id: z.number(),
  cash_available: z.number(),
  access_token: z.string(),
  refresh_token: z.string().uuid(),
});

const UserResponseSchema = createServerResponseSchema(UserSchema);

export type User = z.infer<typeof UserSchema>;

export const Schemas = {
  UserSchema,
  createServerResponseSchema,
  UserResponseSchema,
} as const;
