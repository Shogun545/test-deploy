import { apiClient } from "./apiclient";
import type {
  AdminManagedUsersResponse,
  AdminUserDetailResponse,
  UserFilters,MajorEntry
  ,UserCreatedDateResponse
} from "@/src/interfaces/adminusers"; 

const USERS_API_URL = "/api/admin/users";

// // **[เพิ่มใหม่]** กำหนด URL ฐานของ Go GIN Backend (สมมติพอร์ต 8080)
// // ควรย้ายค่านี้ไปไว้ใน Environment Variables (.env)
// const BASE_BACKEND_URL = "http://localhost:8080"; 

/**
 ดึงรายการผู้ใช้งานทั้งหมด (GET /api/admin/users)
 * @param filters ตัวกรองสำหรับส่งเป็น Query Parameters
 * @returns Promise<AdminManagedUsersResponse>
 */
export async function getManagedUsers(
  filters: UserFilters = {}
): Promise<AdminManagedUsersResponse> {
  const params = new URLSearchParams();
  
  // แปลง filters เป็น Query Params
  Object.entries(filters).forEach(([key, value]) => {
    if (value) {
      params.append(key, value);
    }
  });
  
  const queryString = params.toString();
  
  try {
    const res = await apiClient.get(`${USERS_API_URL}?${queryString}`);
    // Go Controller (GetManagedUsers) ส่ง DTO ตรงๆ (AdminManagedUsersResponse)
    return res.data as AdminManagedUsersResponse;
  } catch (err: any) {
    const errorMessage = err.response?.data?.error || "โหลดรายการผู้ใช้งานล้มเหลว";
    throw new Error(errorMessage);
  }
}

/**
 * ดึงรายละเอียดผู้ใช้งานคนใดคนหนึ่ง (GET /api/admin/users/:sut_id)
 * @param sutId SutID ของผู้ใช้งานที่ต้องการดูรายละเอียด
 * @returns Promise<AdminUserDetailResponse>
 */
export async function getUserDetailBySutId(
  sutId: string
): Promise<AdminUserDetailResponse> {
  try {
    const res = await apiClient.get(`${USERS_API_URL}/${sutId}`);
    // Go Controller (GetUserBySutID) ส่ง DTO ตรงๆ (AdminUserDetailResponse)
    return res.data as AdminUserDetailResponse;
  } catch (err: any) {
    const errorMessage = err.response?.data?.error || "โหลดรายละเอียดผู้ใช้งานล้มเหลว";
    throw new Error(errorMessage);
  }
}

export async function getMajors(): Promise<MajorEntry[]> {
  try {
    const res = await apiClient.get("/api/admin/majors"); // ตามที่ตั้งใน Controller
    return res.data as MajorEntry[];
  } catch (err: any) {
    throw new Error(
      err.response?.data?.error || "โหลดรายการสาขาวิชาล้มเหลว"
    );
  }
}

export async function updateUserStatus(
  sutId: string,
  newStatus: string
): Promise<void> {
  try {
    await apiClient.put(`${USERS_API_URL}/${sutId}/status`, {
      status: newStatus,
    });
  } catch (err: any) {
    const errorMessage = err.response?.data?.error || "อัปเดตสถานะผู้ใช้งานล้มเหลว";
    throw new Error(errorMessage);
  }
}

export const getUserCreatedDate = async (sutId: string): Promise<UserCreatedDateResponse> => {
  try {
    // **[แก้ไข] ใช้ apiClient.get**
    // ใช้ Path สัมพัทธ์ เพราะ apiClient น่าจะมี baseURL ถูกตั้งค่าแล้ว
    const res = await apiClient.get(`${USERS_API_URL}/${sutId}/created-date`); 

    // เนื่องจากใช้ apiClient (Axios) ข้อมูลจะอยู่ใน res.data 
    return res.data as UserCreatedDateResponse;

  } catch (err: any) {
    // ใช้ Error Handling ของ apiClient (Axios)
    const errorMessage = err.response?.data?.error || "โหลดวันที่สร้างผู้ใช้งานล้มเหลว (อาจไม่มีสิทธิ์)";
    throw new Error(errorMessage);
  }
};

export async function updateManagedUser(
  sutId: string,
  payload: {
    phone?: string;
    active?: boolean;
  }
) {
  await apiClient.put(`/api/admin/users/${sutId}`, payload);
}
