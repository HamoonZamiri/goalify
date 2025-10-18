import { useMutation } from "@tanstack/vue-query";
import { toast } from "vue3-toastify";
import { zodFetch } from "@/shared/api";
import { isErrorResponse } from "@/shared/schemas";
import { UserSchema, type LoginFormData } from "@/features/auth/schemas";
import type { User } from "@/features/auth/schemas";
import { API_BASE, http } from "@/utils/constants";

async function loginQueryDataFn(data: LoginFormData): Promise<User> {
	const result = await zodFetch(`${API_BASE}/users/login`, UserSchema, {
		method: http.MethodPost,
		body: JSON.stringify(data),
	});

	if (isErrorResponse(result)) {
		throw new Error(result.message);
	}

	return result;
}

export function useLogin() {
	return useMutation({
		mutationFn: loginQueryDataFn,
		onError: (error: Error) => {
			toast.error(`Login failed: ${error.message}`);
		},
	});
}
