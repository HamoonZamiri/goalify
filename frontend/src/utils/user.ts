import { type User } from "./schemas";
import Cookies from "js-cookie";

export function setUser(user: User) {
  Cookies.set("user", JSON.stringify(user), { sameSite: "lax", expires: 1 });
}

export function getUser(): User | null {
  const user = Cookies.get("user");
  return user ? JSON.parse(user) : null;
}

export function removeUser() {
  Cookies.remove("user");
}

export function isLoggedIn(): boolean {
  return !!getUser();
}
