import { apiClient } from "./apiclient";
import type { StudentProfileResponse } from "@/src/interfaces/studentprofile";

export async function getStudentProfile() {
  try {
    const res = await apiClient.get("/api/student/me/profile");
    return res.data as StudentProfileResponse;
  } catch (err: any) {
    throw new Error(err.response?.data?.message || "โหลดข้อมูลโปรไฟล์ล้มเหลว");
  }
}

export async function updateStudentProfile(data: {
  prefix?: string;
  first_name?: string;
  last_name?: string;
  phone?: string;
  birthday?: string;
  profile_image?: string;
  email?: string;
  year_of_study?: number;
  term_gpa?: number;
  cumulative_gpa?: number;
  academic_year?: string; 
  semester?: number;      
}) {
  try {
    const res = await apiClient.put("/api/student/me/profile", data);
    return res.data as StudentProfileResponse;
  } catch (err: any) {
    throw new Error(err.response?.data?.message || "อัปเดตโปรไฟล์ไม่สำเร็จ");
  }
}
