import { z } from "zod";

export function createServerResponseSchema<TData extends z.ZodTypeAny>(
  schema: TData,
) {
  return z.object({
    message: z.string(),
    data: schema,
  });
}

function createArraySchema<TData extends z.ZodTypeAny>(schema: TData) {
  return z.array(schema);
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

const GoalSchema = z.object({
  id: z.string().uuid(),
  title: z.string(),
  description: z.string(),
  category_id: z.string().uuid(),
  completed: z.boolean(),
  user_id: z.string().uuid(),
  status: z.enum(["complete", "not_complete"]),
  created_at: z.date(),
  updated_at: z.date(),
});

const GoalCategorySchema = z.object({
  id: z.string().uuid(),
  title: z.string(),
  xp_per_goal: z.number(),
  user_id: z.string().uuid(),
  goals: z.array(GoalSchema),
});

const UserResponseSchema = createServerResponseSchema(UserSchema);
const GoalResponseSchema = createServerResponseSchema(GoalSchema);
const GoalCategoryResponseSchema =
  createServerResponseSchema(GoalCategorySchema);
const GoalCategoryArraySchema = createArraySchema(GoalCategorySchema);
const GoalArraySchema = createArraySchema(GoalSchema);
const GoalCategoryResponseArraySchema = createServerResponseSchema(
  GoalCategoryArraySchema,
);
const GoalResponseArraySchema = createServerResponseSchema(GoalArraySchema);

export type User = z.infer<typeof UserSchema>;
export type Goal = z.infer<typeof GoalSchema>;
export type GoalCategory = z.infer<typeof GoalCategorySchema>;

export const Schemas = {
  UserSchema,
  createServerResponseSchema,
  GoalCategorySchema,
  GoalSchema,
  UserResponseSchema,
  GoalResponseSchema,
  GoalCategoryResponseSchema,
  GoalCategoryArraySchema,
  GoalResponseArraySchema,
  GoalCategoryResponseArraySchema,
} as const;
