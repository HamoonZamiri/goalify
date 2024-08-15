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
  getUser: getUserFromCookie,
  isLoggedIn() {
    return userCookieExists() && !!this.getUser;
  },
  setUser(user: User) {
    setUserCookie(user);
  },
  logout() {
    removeUserFromCookie();
  },
});

export default authState;
