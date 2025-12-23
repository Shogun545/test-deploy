import { API_BASE_URL } from "@/src/config/api";

export function ensureAbsoluteUrl(
  url?: string | null
): string | null {
  if (!url) return null;
  if (/^https?:\/\//i.test(url)) return url;

  const prefix = API_BASE_URL.replace(/\/?$/, "");
  const path = url.startsWith("/") ? url : `/${url}`;
  return `${prefix}${path}`;
}
