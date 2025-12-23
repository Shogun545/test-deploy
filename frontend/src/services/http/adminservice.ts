import { apiClient } from "./apiclient";
import type {
  AdminProfileResponse,
  AdminProfileUpdateData,
} from "@/src/interfaces/adminprofile";

const ADMIN_API_URL = "/api/admin/me/profile";

/**
 * ดึงข้อมูลโปรไฟล์ผู้ดูแลระบบ
 */
export async function getAdminProfile(): Promise<AdminProfileResponse> {
  try {
    const res = await apiClient.get(ADMIN_API_URL);
    return res.data as AdminProfileResponse;
  } catch (err: any) {
    throw new Error(
      err.response?.data?.message || "โหลดข้อมูลโปรไฟล์ผู้ดูแลระบบล้มเหลว"
    );
  }
}

/**
 * อัปเดตข้อมูลโปรไฟล์ผู้ดูแลระบบ
 * (backend ใช้ PUT)
 */
export async function updateAdminProfile(
  profileData: AdminProfileUpdateData
): Promise<AdminProfileResponse> {
  try {
    const res = await apiClient.put(ADMIN_API_URL, profileData);
    return res.data as AdminProfileResponse;
  } catch (err: any) {
    throw new Error(
      err.response?.data?.message || "อัปเดตข้อมูลโปรไฟล์ผู้ดูแลระบบล้มเหลว"
    );
  }
}
