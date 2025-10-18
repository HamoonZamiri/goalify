import { z } from "zod";

const SYMBOLS = "!@#$%^";
const PASSWORD_MIN_LEN = 8;

/**
 * Login form validation schema
 */
export const loginFormSchema = z.object({
	email: z.string().min(1, "Email is required").email("Email is invalid"),
	password: z
		.string()
		.min(1, "Password is required")
		.min(
			PASSWORD_MIN_LEN,
			`Password must be at least ${PASSWORD_MIN_LEN} characters`,
		)
		.refine((val) => !val.includes(" "), "Password cannot contain spaces")
		.refine((val) => /\d/.test(val), "Password must contain at least one digit")
		.refine(
			(val) =>
				new RegExp(`[${SYMBOLS.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")}]`).test(
					val,
				),
			`Password must contain at least one symbol (${SYMBOLS})`,
		),
});
export type LoginFormData = z.infer<typeof loginFormSchema>;

/**
 * Register form validation schema
 */
export const registerFormSchema = z
	.object({
		email: z.string().min(1, "Email is required").email("Email is invalid"),
		password: z
			.string()
			.min(1, "Password is required")
			.min(
				PASSWORD_MIN_LEN,
				`Password must be at least ${PASSWORD_MIN_LEN} characters`,
			)
			.refine((val) => !val.includes(" "), "Password cannot contain spaces")
			.refine(
				(val) => /\d/.test(val),
				"Password must contain at least one digit",
			)
			.refine(
				(val) =>
					new RegExp(
						`[${SYMBOLS.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")}]`,
					).test(val),
				`Password must contain at least one symbol (${SYMBOLS})`,
			),
		confirmPassword: z.string().min(1, "Confirm password is required"),
	})
	.refine((data) => data.password === data.confirmPassword, {
		message: "Passwords do not match",
		path: ["confirmPassword"],
	});
export type RegisterFormData = z.infer<typeof registerFormSchema>;
