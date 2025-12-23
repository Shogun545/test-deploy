// src/services/http/advisorservice.ts
import { apiClient } from "./apiclient";
import type { AdvisorProfileResponse } from "@/src/interfaces/advisorprofile";

// ==============================
// Advisor Profile APIs
// ==============================
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

// ==============================
// Appointment Types + Helpers (ใช้ร่วมกันทั้ง list + detail)
// ==============================
export type AppointmentStatus = "PENDING" | "APPROVED" | "RESCHEDULE";
export type ViewMode = "PENDING" | "DONE" | "ALL"; // ✅ เพิ่มให้หน้า list import ได้

export type AppointmentListRow = {
  id: number;
  studentName: string;
  sutId: string;
  topic: string;
  submittedAt: string;
  status: AppointmentStatus;
};

export type AppointmentDetail = {
  id: number;
  studentName: string;
  sutId: string;
  topic: string;
  submittedAt: string;
  description: string;
  status: AppointmentStatus;
};

// helper: map backend status_code -> UI status
const normalizeAppointmentStatus = (raw: unknown): AppointmentStatus => {
  const s = String(raw || "").toUpperCase();
  if (s === "APPROVED") return "APPROVED";
  if (s === "RESCHEDULE" || s === "REJECTED") return "RESCHEDULE"; // ✅ กันกรณี backend ยังส่ง REJECTED มา
  return "PENDING";
};

// helper: safe date
const formatThaiDate = (raw: unknown) => {
  if (!raw) return "-";
  const d = new Date(String(raw));
  if (Number.isNaN(d.getTime())) return "-";
  return d.toLocaleDateString("th-TH");
};

// ==============================
// Appointment List APIs (pending / done / all)
// ==============================

/**
 * GET /api/appointments/pending
 */
export async function listAppointmentsPending(): Promise<AppointmentListRow[]> {
  try {
    const res = await apiClient.get("/api/appointments/pending");
    const data = (res.data || []) as Array<{
      id: number;
      studentName: string;
      sutId: string;
      topic: string;
      submittedAt: string;
      status: string; // PENDING | APPROVED | RESCHEDULE
    }>;

    return data.map((x) => ({
      id: Number(x.id),
      studentName: x.studentName || "-",
      sutId: x.sutId || "-",
      topic: x.topic || "-",
      submittedAt: x.submittedAt || "-",
      status: normalizeAppointmentStatus(x.status),
    }));
  } catch (err: any) {
    throw new Error(
      err.response?.data?.message || "โหลดรายการคำขอที่รอพิจารณาล้มเหลว"
    );
  }
}

/**
 * GET /api/appointments/done (APPROVED + RESCHEDULE)
 */
export async function listAppointmentsDone(): Promise<AppointmentListRow[]> {
  try {
    const res = await apiClient.get("/api/appointments/done");
    const data = (res.data || []) as Array<{
      id: number;
      studentName: string;
      sutId: string;
      topic: string;
      submittedAt: string;
      status: string;
    }>;

    return data.map((x) => ({
      id: Number(x.id),
      studentName: x.studentName || "-",
      sutId: x.sutId || "-",
      topic: x.topic || "-",
      submittedAt: x.submittedAt || "-",
      status: normalizeAppointmentStatus(x.status),
    }));
  } catch (err: any) {
    throw new Error(
      err.response?.data?.message || "โหลดรายการคำขอที่พิจารณาแล้วล้มเหลว"
    );
  }
}

/**
 * GET /api/appointments (all)
 */
export async function listAppointmentsAll(): Promise<AppointmentListRow[]> {
  try {
    const res = await apiClient.get("/api/appointments");
    const data = (res.data || []) as Array<{
      id: number;
      studentName: string;
      sutId: string;
      topic: string;
      submittedAt: string;
      status: string;
    }>;

    return data.map((x) => ({
      id: Number(x.id),
      studentName: x.studentName || "-",
      sutId: x.sutId || "-",
      topic: x.topic || "-",
      submittedAt: x.submittedAt || "-",
      status: normalizeAppointmentStatus(x.status),
    }));
  } catch (err: any) {
    throw new Error(
      err.response?.data?.message || "โหลดรายการคำขอทั้งหมดล้มเหลว"
    );
  }
}

/**
 * helper รวม: เรียกตามแท็บ (PENDING/DONE/ALL)
 * ใช้ในหน้า app/advisor/appointments/page.tsx ได้เลย
 */
export async function fetchAdvisorAppointments(
  mode: ViewMode
): Promise<AppointmentListRow[]> {
  if (mode === "PENDING") return listAppointmentsPending();
  if (mode === "DONE") return listAppointmentsDone();
  return listAppointmentsAll();
}

// ==============================
// Appointment Detail + Actions APIs ([id])
// ==============================

/**
 * GET /api/appointments/:id
 */
export async function getAppointmentById(
  id: number
): Promise<AppointmentDetail> {
  try {
    const res = await apiClient.get(`/api/appointments/${id}`);
    const appt = res.data;

    const apptId = Number(appt?.id ?? appt?.ID ?? 0);

    const firstName =
      appt?.student_user?.first_name ?? appt?.student_user?.FirstName ?? "";
    const lastName =
      appt?.student_user?.last_name ?? appt?.student_user?.LastName ?? "";
    const studentName = `${firstName} ${lastName}`.trim() || "-";

    const sutId =
      appt?.student_user?.sut_id ?? appt?.student_user?.SutID ?? "-";

    const description = appt?.description ?? appt?.Description ?? "-";
    const submittedAt = formatThaiDate(appt?.created_at ?? appt?.CreatedAt);

    const rawStatus =
      appt?.appointment_status?.status_code ??
      appt?.appointment_status?.StatusCode ??
      appt?.appointment_status?.status ??
      appt?.appointment_status?.Status ??
      appt?.AppointmentStatus?.status_code ??
      appt?.AppointmentStatus?.StatusCode;

    return {
      id: apptId,
      studentName,
      sutId,
      topic: description,
      description,
      submittedAt,
      status: normalizeAppointmentStatus(rawStatus),
    };
  } catch (err: any) {
    throw new Error(
      err.response?.data?.message || "โหลดรายละเอียดคำขอนัดหมายล้มเหลว"
    );
  }
}

/**
 * PUT /api/appointments/:id/approve
 */
export async function approveAppointment(
  id: number,
  description = "อนุมัติคำขอ"
) {
  try {
    const res = await apiClient.put(`/api/appointments/${id}/approve`, {
      description,
    });
    return res.data;
  } catch (err: any) {
    throw new Error(err.response?.data?.message || "อนุมัติคำขอล้มเหลว");
  }
}

/**
 * PUT /api/appointments/:id/reschedule
 */
export async function rescheduleAppointment(
  id: number,
  description = "เสนอเวลาใหม่"
) {
  try {
    const res = await apiClient.put(`/api/appointments/${id}/reschedule`, {
      description,
    });
    return res.data;
  } catch (err: any) {
    throw new Error(err.response?.data?.message || "เสนอเวลาใหม่ล้มเหลว");
  }
}
