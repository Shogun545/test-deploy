import { apiClient } from "./apiclient";
import type { AdvisorProfileResponse } from "@/src/interfaces/advisorprofile";

export async function getAdvisorProfile() {
  try {
    const res = await apiClient.get("/api/advisor/me/profile");
    return res.data as AdvisorProfileResponse;
  } catch (err: any) {
    throw new Error(err.response?.data?.message || "โหลดข้อมูลโปรไฟล์ล้มเหลว");
  }
}

export async function updateAdvisorProfile(
  profileData: Partial<AdvisorProfileResponse>
) {
  try {
    const res = await apiClient.put("/api/advisor/me/profile", profileData);
    return res.data as AdvisorProfileResponse;
  } catch (err: any) {
    throw new Error(err.response?.data?.message || "อัปเดตข้อมูลโปรไฟล์ล้มเหลว");
  }
}

// export async function uploadAdvisorProfilePicture(
//   file: File
// ): Promise<string> {
//   try {
//     const formData = new FormData();
//     formData.append("profilePicture", file);

//     const res = await apiClient.post(
//       "/api/advisor/me/profile/picture",
//       formData,
//       {
//         headers: {
//           "Content-Type": "multipart/form-data",
//         },
//       }
//     );

//     return res.data.imageUrl as string;
//   } catch (err: any) {
//     throw new Error(
//       err.response?.data?.message || "อัปโหลดรูปโปรไฟล์ล้มเหลว"
//     );
//   }
// }
