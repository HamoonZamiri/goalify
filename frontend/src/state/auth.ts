import { reactive } from "vue";
import { type User } from "../utils/schemas";
import Cookies from "js-cookie";

function setUserCookie(user: User) {
  Cookies.set("user", JSON.stringify(user), { sameSite: "lax", expires: 1 });
}

function getUserFromCookie(): User | null {
  const user = Cookies.get("user");
  return user ? JSON.parse(user) : null;
}

function removeUserFromCookie() {
  Cookies.remove("user");
}

function userCookieExists(): boolean {
  return !!getUserFromCookie();
}

const authState = reactive({
  user: getUserFromCookie(),
  isLoggedIn() {
    return userCookieExists() && !!this.user;
  },
  setUser(user: User) {
    setUserCookie(user);
    this.user = user;
  },
  logout() {
    removeUserFromCookie();
    this.user = null;
  },
});

export default authState;
