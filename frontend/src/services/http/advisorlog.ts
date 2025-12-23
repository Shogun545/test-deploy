import { apiClient } from "./apiclient";
import { AdvisorLogInterface } from "@/src/interfaces/advisorlog";
import { AxiosProgressEvent } from "axios";

/* ---------- Common Result Type ---------- */
export type ApiResult<T = any> =
  | { success: true; data: T }
  | { success: false; status?: number; message?: string };

/* ---------- Upload Options ---------- */
type UploadOptions = {
  files?: File[];
  isDraft?: boolean;
  onUploadProgress?: (percent: number) => void;
};

/* ---------- Base URL ---------- */
const BASE_URL = "/api/advisor_logs";

/* =========================================================
   CREATE Advisor Log (multipart + files)
   POST /api/advisor_logs
========================================================= */
export async function CreateAdvisorLog(
  data: AdvisorLogInterface,
  options: UploadOptions = {}
): Promise<ApiResult> {
  const { files = [], isDraft = false, onUploadProgress } = options;

  try {
    const formData = new FormData();

    // ✅ match backend dto: appointmentId, title, body, requiresReport, status
    formData.append("appointmentId", String(data.AppointmentID ?? 0));
    formData.append("title", data.Title ?? "");
    formData.append("body", data.Body ?? "");
    formData.append("requiresReport", data.RequiresReport ? "true" : "false");

    // draft ควรส่งเป็นตัวเล็ก (backend เช็ค strings.ToLower == "draft")
    if (isDraft) formData.append("status", "draft");
    else if (data.Status) formData.append("status", data.Status);

    files.forEach((file) => formData.append("files", file));

    const resp = await apiClient.post(BASE_URL, formData, {
      headers: { "Content-Type": "multipart/form-data" },
      onUploadProgress: (ev: AxiosProgressEvent) => {
        if (onUploadProgress && ev.total) {
          const percent = Math.round((ev.loaded * 100) / ev.total);
          onUploadProgress(percent);
        }
      },
    });

    return { success: true, data: resp.data };
  } catch (err: any) {
    return {
      success: false,
      status: err.response?.status,
      message: err.response?.data?.error || err.message,
    };
  }
}

/* =========================================================
   GET Advisor Log by ID
   GET /api/advisor_logs/:id
========================================================= */
export async function GetAdvisorLogByID(id: number | string): Promise<ApiResult> {
  try {
    const resp = await apiClient.get(`${BASE_URL}/${id}`);
    return { success: true, data: resp.data };
  } catch (err: any) {
    return {
      success: false,
      status: err.response?.status,
      message: err.response?.data?.error || err.message,
    };
  }
}

/* =========================================================
   LIST Logs by Student
   GET /api/advisor_logs/student/:studentId
========================================================= */
export async function ListAdvisorLogsByStudent(
  studentId: number | string
): Promise<ApiResult> {
  try {
    const resp = await apiClient.get(`${BASE_URL}/student/${studentId}`);
    return { success: true, data: resp.data };
  } catch (err: any) {
    return {
      success: false,
      status: err.response?.status,
      message: err.response?.data?.error || err.message,
    };
  }
}

/* =========================================================
   LIST Logs (Advisor sees all his logs)
   GET /api/advisor_logs
========================================================= */
export async function ListAdvisorLogs(): Promise<ApiResult> {
  try {
    const resp = await apiClient.get(BASE_URL);
    return { success: true, data: resp.data };
  } catch (err: any) {
    return {
      success: false,
      status: err.response?.status,
      message: err.response?.data?.error || err.message,
    };
  }
}

/* =========================================================
   UPDATE STATUS ONLY
   PATCH /api/advisor_logs/:id
   ✅ backend รับค่า: "Draft" | "PendingReport" | "Completed"
========================================================= */
export async function UpdateAdvisorLogStatus(
  id: number | string,
  status: string
): Promise<ApiResult> {
  try {
    const resp = await apiClient.patch(`${BASE_URL}/${id}`, { status });
    return { success: true, data: resp.data };
  } catch (err: any) {
    return {
      success: false,
      status: err.response?.status,
      message: err.response?.data?.error || err.message,
    };
  }
}

/* =========================================================
   UPDATE FULL LOG (Title / Body / RequiresReport / Files)
   PATCH /api/advisor_logs/:id/edit (multipart)
========================================================= */
export async function UpdateAdvisorLog(
  id: number | string,
  data: Partial<AdvisorLogInterface>,
  options: UploadOptions = {}
): Promise<ApiResult> {
  const { files = [], onUploadProgress } = options;

  try {
    const formData = new FormData();

    // ✅ แก้ไข: เปลี่ยน Key เป็น lowercase ให้ตรงกับ Backend DTO (`form:"title"`)
    if (data.Title !== undefined) formData.append("title", data.Title ?? ""); // แก้ Title -> title
    if (data.Body !== undefined) formData.append("body", data.Body ?? "");    // แก้ Body -> body
    if (data.RequiresReport !== undefined) {
      formData.append("requiresReport", data.RequiresReport ? "true" : "false"); // แก้ RequiresReport -> requiresReport
    }

    files.forEach((file) => formData.append("files", file));

    const resp = await apiClient.patch(`${BASE_URL}/${id}/edit`, formData, {
      headers: { "Content-Type": "multipart/form-data" },
      onUploadProgress: (ev: AxiosProgressEvent) => {
        if (onUploadProgress && ev.total) {
          const percent = Math.round((ev.loaded * 100) / ev.total);
          onUploadProgress(percent);
        }
      },
    });

    return { success: true, data: resp.data };
  } catch (err: any) {
    return {
      success: false,
      status: err.response?.status,
      message: err.response?.data?.error || err.message,
    };
  }
}
