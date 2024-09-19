import { type User } from "@/utils/schemas";
import Cookies from "js-cookie";
import { ref } from "vue";

const authState = ref<User | null>(null);

function useAuth() {
  function setUser(user: User) {
    Cookies.set("user", JSON.stringify(user), { sameSite: "lax", expires: 1 });
    authState.value = user;
  }

  function getUser(): User | null {
    if (authState.value) return authState.value;
    const user = Cookies.get("user");
    if (user) setUser(JSON.parse(user));
    return user ? JSON.parse(user) : null;
  }

  function logout() {
    Cookies.remove("user");
    authState.value = null;
  }

  function isLoggedIn(): boolean {
    return !!getUser();
  }

  return {
    authState,
    setUser,
    getUser,
    logout,
    isLoggedIn,
  };
}

export default useAuth;
