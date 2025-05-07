import type { Goal, GoalCategory, Level, User } from "@/utils/schemas";

export const goal: Goal = {
	id: "026a715f-a023-4e6b-973e-4bb0e96562ae",
	title: "Test Goal",
	description: "Test Goal Description",
	category_id: "026a715f-a023-4e6b-973e-4bb0e96562ae",
	user_id: "026a715f-a023-4e6b-973e-4bb0e96562ae",
	status: "not_complete",
	created_at: new Date().toString(),
	updated_at: new Date().toString(),
};

export const goalCategory: GoalCategory = {
	title: "new Category testing events 2",
	goals: [],
	xp_per_goal: 50,
	id: "026a715f-a023-4e6b-973e-4bb0e96562ae",
	user_id: "1c0736f6-2c67-4e4b-a34c-afe3b78d1ab1",
};

export const user: User = {
	id: "026a715f-a023-4e6b-973e-4bb0e96562ae",
	access_token: "somerandomaccesstoken",
	cash_available: 200,
	email: "test@mail.com",
	level_id: 1,
	refresh_token: "026a715f-a023-4e6b-973e-4bb0e96562ae",
	xp: 50,
};

export const levelOne: Level = {
	cash_reward: 100,
	id: 1,
	level_up_xp: 100,
};
