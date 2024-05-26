type User = {
  id: string;
  email: string;
  xp: number;
  level_id: number;
  cash_available: number;
  access_token: string;
  refresh_token: string;
};

type Goal = {
  title: string;
  description: string;
  status: "complete" | "not_complete";
  id: string;
  user_id: string;
  category_id: string;
};

type GoalCategory = {
  title: string;
  goals: Goal[];
  xp_per_goal: number;
  id: string;
  user_id: string;
};
