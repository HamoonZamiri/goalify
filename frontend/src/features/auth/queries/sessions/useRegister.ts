import { useMutation } from "@tanstack/vue-query";
import { toast } from "vue3-toastify";
import { zodFetch } from "@/shared/api";
import { isErrorResponse } from "@/shared/schemas";
import { UserSchema, type RegisterFormData } from "@/features/auth/schemas";
import type { User } from "@/features/auth/schemas";
import { API_BASE, http } from "@/utils/constants";

async function registerQueryDataFn(data: RegisterFormData): Promise<User> {
	const result = await zodFetch(`${API_BASE}/users/signup`, UserSchema, {
		method: http.MethodPost,
		body: JSON.stringify({
			email: data.email,
			password: data.password,
			confirm_password: data.confirmPassword,
		}),
	});

	if (isErrorResponse(result)) {
		throw new Error(result.message);
	}

	return result;
}

export function useRegister() {
	return useMutation({
		mutationFn: registerQueryDataFn,
		onError: (error: Error) => {
			toast.error(`Registration failed: ${error.message}`);
		},
	});
}
