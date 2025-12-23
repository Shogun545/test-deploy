import { apiClient } from "./apiclient";
import type { Prefix } from "@/src/interfaces/prefix";

export async function getPrefixes(): Promise<Prefix[]> {
  const res = await apiClient.get("/api/master/prefixes");
  return res.data;
}
