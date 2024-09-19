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
  user_id: z.string().uuid(),
  status: z.enum(["complete", "not_complete"]),
  created_at: z.string(),
  updated_at: z.string(),
});

const LevelSchema = z.object({
  id: z.number(),
  level_up_xp: z.number(),
  cash_reward: z.number(),
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
const LevelResponseSchema = createServerResponseSchema(LevelSchema);

export type User = z.infer<typeof UserSchema>;
export type Goal = z.infer<typeof GoalSchema>;
export type GoalCategory = z.infer<typeof GoalCategorySchema>;
export type Level = z.infer<typeof LevelSchema>;

type ErrorMap = Record<string, string>;
export type ErrorResponse = {
  // we will manually add this field to our errors from the json response
  statusCode?: number;
  // message should always be present
  message: string;
  // in creation requests the server returns an object mapping field names to error messages
  errors?: ErrorMap;
};

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
  LevelResponseSchema: LevelResponseSchema,
} as const;
