export type User = {
  id: string;
  email: string;
  xp: number;
  level_id: number;
  cash_available: number;
  access_token: string;
  refresh_token: string;
};

export type TGoal = {
  title: string;
  description: string;
  status: "complete" | "not_complete";
  id: string;
  user_id: string;
  category_id: string;
  updated_at: Date;
  created_at: Date;
};

export type TGoalCategory = {
  title: string;
  goals: TGoal[];
  xp_per_goal: number;
  id: string;
  user_id: string;
  updated_at: Date;
  created_at: Date;
};

export type UserDTO = {
  email: string;
  access_token: string;
  refresh_token: string;
  xp: number;
  level_id: number;
  cash_available: number;
  id: string;
};

export const mockGoal: TGoal = {
  title: "Complete 10 leetcode questions with a score of 80% or higher",
  description: "Test description",
  status: "not_complete",
  id: "1",
  user_id: "1",
  category_id: "1",
  updated_at: new Date(),
  created_at: new Date(),
};

export const mockGoalCategory: TGoalCategory = {
  title: "Daily",
  goals: [mockGoal, mockGoal, mockGoal],
  xp_per_goal: 10,
  id: "1",
  user_id: "1",
  updated_at: new Date(),
  created_at: new Date(),
};
