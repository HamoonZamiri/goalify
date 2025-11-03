import Cookies from "js-cookie";
import { ref } from "vue";
import { UserSchema, type User } from "@/features/auth/schemas";

const authState = ref<User>();

function useAuth() {
	function setUser(user: User) {
		Cookies.set("user", JSON.stringify(user), { sameSite: "lax", expires: 1 });
		authState.value = user;
	}

	function logout() {
		Cookies.remove("user");
		authState.value = undefined;
	}

	function isLoggedIn(): boolean {
		return !!getUser();
	}

	function getUser(): User | undefined {
		if (!authState.value) {
			const userCookie = Cookies.get("user");
			if (!userCookie) return;
			const parsed = JSON.parse(userCookie);
			const parseResult = UserSchema.safeParse(parsed);
			if (!parseResult.success) return;
			authState.value = parseResult.data;
		}
		return authState.value;
	}

	return {
		authState,
		getUser,
		setUser,
		logout,
		isLoggedIn,
	};
}

export default useAuth;
