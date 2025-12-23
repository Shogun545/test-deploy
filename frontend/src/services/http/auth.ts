import { apiClient } from "./apiclient";
import type { LoginInterface } from "../../interfaces/user";

export const loginApi = async (credentials: LoginInterface) => {
  const res = await apiClient.post("/api/auth/login", credentials);
  return res.data; // axios แปลง JSON ให้แล้ว
};
