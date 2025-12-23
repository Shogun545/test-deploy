import { apiClient } from "./apiclient";
import type { AdvisorStudentsResponse, StudentForAdvisor } from "@/src/interfaces/advisorstudent";

// fetch list (optionally by major)
export async function getAdvisorStudents(params?: { major?: string }): Promise<AdvisorStudentsResponse> {
  const res = await apiClient.get<AdvisorStudentsResponse>("/api/advisor/me/students", { params });
  return res.data;
}

// fetch detail by sut_id (backend ต้องมี route เช่น /api/advisor/students/:sutId)
export async function getAdvisorStudentDetail(sutId: string): Promise<StudentForAdvisor> {
  const res = await apiClient.get<StudentForAdvisor>(`/api/advisor/me/students/${encodeURIComponent(sutId)}`);
  return res.data;
}
