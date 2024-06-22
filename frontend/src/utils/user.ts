import type { UserDTO } from "@/utils/types";
import Cookies from "js-cookie";

export function setUser(user: UserDTO) {
  Cookies.set("user", JSON.stringify(user), { expires: 1 });
}

export function getUser(): UserDTO | null {
  const user = Cookies.get("user");
  return user ? JSON.parse(user) : null;
}

export function removeUser() {
  Cookies.remove("user");
}
