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

export type User = z.infer<typeof UserSchema>;
