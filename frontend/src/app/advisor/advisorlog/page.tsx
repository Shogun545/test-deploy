"use client";

import { useEffect, useMemo, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";

import { apiClient } from "@/src/services/http/apiclient";
import { ListAdvisorLogs, ListAdvisorLogsByStudent } from "@/src/services/http/advisorlog";

type LogStatus = "Draft" | "PendingReport" | "Completed";
type BadgeTone = "warn" | "gray" | "success";

type AdvisorLogNorm = {
  id: number;
  appointmentId: number;
  status: LogStatus;
  requiresReport: boolean;
  fileName?: string;
  filePath?: string;
  createdAt?: string;
};

type AppointmentNorm = {
  id: number; // appointment id
  topic?: string;
  startAt?: string;
  endAt?: string;
  student: { sutId?: string; fullName?: string };
};

type Card = {
  appointment: AppointmentNorm;
  advisorLog: AdvisorLogNorm | null;
  ui: {
    badgeText: "รอการบันทึก" | "ฉบับร่าง" | "บันทึกแล้ว";
    badgeTone: BadgeTone;
    primaryAction: { label: string; href: string };
    secondaryAction?: { label: string; href: string };
  };
};

/* =========================
   Small utils
========================= */
const toNum = (v: any) => {
  const n = Number(v);
  return Number.isFinite(n) ? n : 0;
};

const pick = <T,>(obj: any, keys: string[], fallback: T): T => {
  for (const k of keys) {
    const v = obj?.[k];
    if (v !== undefined && v !== null) return v as T;
  }
  return fallback;
};

function formatDateMaybe(s?: string) {
  if (!s) return "";
  const d = new Date(s);
  if (Number.isNaN(d.getTime())) return s;
  return d.toLocaleString();
}

function badgeStyle(tone: BadgeTone) {
  if (tone === "warn") return "bg-[#FFE6D6] text-[#F26522]";
  if (tone === "success") return "bg-[#E0F9EA] text-[#16A34A]";
  return "bg-[#E5E7EB] text-[#4B5563]";
}

/* =========================
   JWT helper (get user_id)
   - ใช้เพื่อส่ง advisor_id ให้ endpoint /api/appointments/done
========================= */
function getAdvisorIdFromToken(): number | null {
  if (typeof window === "undefined") return null;

  const token =
    window.localStorage.getItem("token") ||
    window.localStorage.getItem("access_token") ||
    window.localStorage.getItem("accessToken") ||
    "";

  if (!token) return null;

  try {
    const payload = JSON.parse(atob(token.split(".")[1]));
    const id = Number(payload?.user_id);
    return Number.isFinite(id) && id > 0 ? id : null;
  } catch {
    return null;
  }
}

/* =========================
   Normalizers
========================= */
function normalizeLog(raw: any): AdvisorLogNorm {
  const id = toNum(pick(raw, ["ID", "id"], 0));
  const appointmentId = toNum(pick(raw, ["AppointmentID", "appointmentId", "appointment_id"], 0));
  const status = pick<string>(raw, ["Status", "status"], "Draft") as LogStatus;

  const requiresReport = Boolean(
    pick(raw, ["RequiresReport", "requiresReport", "requires_report"], false)
  );

  const fileName = pick<string>(raw, ["FileName", "fileName", "file_name"], "");
  const filePath = pick<string>(raw, ["FilePath", "filePath", "file_path"], "");
  const createdAt = pick<string>(raw, ["CreatedAt", "createdAt", "created_at"], "");

  return {
    id,
    appointmentId,
    status,
    requiresReport,
    fileName: fileName || undefined,
    filePath: filePath || undefined,
    createdAt: createdAt || undefined,
  };
}

function normalizeAppointment(raw: any): AppointmentNorm {
  const id = toNum(pick(raw, ["ID", "id", "appointmentId", "appointment_id"], 0));

  const topic =
    String(
      pick(raw, ["TopicName", "topicName", "topic"], "") ||
        pick(raw?.Topic, ["Name", "name"], "") ||
        pick(raw?.topic, ["name"], "")
    ) || undefined;

  const startAt =
    String(pick(raw, ["StartAt", "startAt", "start_at", "appointment_time"], "")) || undefined;

  const endAt =
    String(pick(raw, ["EndAt", "endAt", "end_at"], "")) || undefined;

  const sutId =
    String(
      pick(raw, ["sutId", "sut_id", "StudentSutID", "studentSutId"], "") ||
        pick(raw?.StudentUser, ["SutID", "sutId", "sut_id"], "") ||
        pick(raw?.student, ["sutId", "sut_id"], "")
    ) || undefined;

  const fullName =
    String(
      pick(raw, ["studentName", "StudentName", "student_full_name", "fullName"], "") ||
        pick(raw?.StudentUser, ["FullName", "fullName", "Name", "name"], "") ||
        pick(raw?.student, ["fullName", "name"], "")
    ) || undefined;

  return {
    id,
    topic,
    startAt,
    endAt,
    student: { sutId, fullName },
  };
}

/* =========================
   Fetch appointments (DONE)
   - ใช้ /api/appointments/done เป็นหลัก (ตาม network ของคุณ)
   - ส่ง advisor_id ถ้ามี
========================= */
async function fetchDoneAppointments(studentId?: string) {
  const advisorId = getAdvisorIdFromToken();

  const candidates = [
    "/api/appointments/done",
    "/api/advisor/me/appointments/done",
    "/api/advisor/me/appointments", // เผื่อไม่มี done แยก
    "/api/appointments",
  ];

  const params: Record<string, any> = {};
  if (advisorId) params.advisor_id = advisorId;
  if (studentId) params.studentId = studentId;

  let lastErr: any = null;

  for (const url of candidates) {
    try {
      const resp = await apiClient.get(url, { params });
      const payload = resp.data?.data ?? resp.data ?? [];
      if (Array.isArray(payload)) return payload;
      return payload ? [payload] : [];
    } catch (e: any) {
      lastErr = e;
      // ลองตัวถัดไป
    }
  }

  const msg =
    lastErr?.response?.data?.error ||
    lastErr?.message ||
    "ไม่พบ endpoint สำหรับดึงนัดหมายที่เสร็จสิ้น";
  throw new Error(msg);
}

/* =========================
   Page
========================= */
export default function AdvisorLogPage() {
  const router = useRouter();
  const search = useSearchParams();

  const studentId = search?.get("studentId") ?? "";

  const [cards, setCards] = useState<Card[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const totalText = useMemo(
    () => (loading ? "..." : `${cards.length} รายการ`),
    [loading, cards.length]
  );

  useEffect(() => {
    let mounted = true;

    (async () => {
      setLoading(true);
      setError(null);

      try {
        // 1) appointments (done)
        const apptsRaw = await fetchDoneAppointments(studentId || undefined);
        let apptsArr = Array.isArray(apptsRaw) ? apptsRaw : (apptsRaw ? [apptsRaw] : []);
        const apptsNorm = apptsArr.map(normalizeAppointment).filter((a) => a.id > 0);

        // 2) logs
        let logsRaw: any[] = [];
        if (studentId) {
          const res = await ListAdvisorLogsByStudent(studentId);
          if (!res.success) throw new Error(res.message || "โหลดบันทึกไม่สำเร็จ");
          logsRaw = res.data?.data ?? res.data ?? [];
        } else {
          const res = await ListAdvisorLogs();
          if (!res.success) throw new Error(res.message || "โหลดบันทึกไม่สำเร็จ");
          logsRaw = res.data?.data ?? res.data ?? [];
        }
        if (!Array.isArray(logsRaw)) logsRaw = logsRaw ? [logsRaw] : [];

        const logByAppt = new Map<number, AdvisorLogNorm>();
        for (const l of logsRaw) {
          const nl = normalizeLog(l);
          if (nl.appointmentId) logByAppt.set(nl.appointmentId, nl);
        }

        // 3) build cards
        const built: Card[] = apptsNorm.map((appt) => {
          const log = logByAppt.get(appt.id) ?? null;

          // UI rules ตามแบบ
          let badgeText: Card["ui"]["badgeText"] = "รอการบันทึก";
          let badgeTone: BadgeTone = "warn";
          let primaryAction: Card["ui"]["primaryAction"] = {
            label: "สร้างบันทึกการปรึกษา",
            href: `/advisor/advisorlog/create?appointmentId=${appt.id}`,
          };
          let secondaryAction: Card["ui"]["secondaryAction"] | undefined;

          if (log) {
            secondaryAction = { label: "ดูรายละเอียด", href: `/advisor/advisorlog/${log.id}` };

            if (log.status === "Draft") {
              badgeText = "ฉบับร่าง";
              badgeTone = "gray";
              primaryAction = { label: "แก้ไขฉบับร่าง", href: `/advisor/advisorlog/edit/${log.id}` };
            } else if (log.status === "Completed") {
              badgeText = "บันทึกแล้ว";
              badgeTone = "success";
              primaryAction = { label: "ดูรายละเอียด", href: `/advisor/advisorlog/${log.id}` };
              secondaryAction = undefined; 
            } else {
              badgeText = "บันทึกแล้ว";
              badgeTone = "success";
              primaryAction = { label: "ดูรายละเอียด", href: `/advisor/advisorlog/${log.id}` };
              secondaryAction = undefined; 
            }
          }

          return {
            appointment: appt,
            advisorLog: log,
            ui: { badgeText, badgeTone, primaryAction, secondaryAction },
          };
        });

        // sort ล่าสุดก่อน (ใช้ startAt ถ้ามี)
        built.sort((a, b) => (b.appointment.startAt || "").localeCompare(a.appointment.startAt || ""));

        if (!mounted) return;
        setCards(built);
      } catch (e: any) {
        if (!mounted) return;
        setError(e?.message ?? "เกิดข้อผิดพลาด");
        setCards([]);
      } finally {
        if (mounted) setLoading(false);
      }
    })();

    return () => {
      mounted = false;
    };
  }, [studentId]);

  return (
  <div className="w-full min-h-screen bg-gray-50 flex justify-center py-6 px-4">
    <div className="w-full max-w-6xl bg-white rounded-3xl shadow-sm border border-gray-200 px-10 py-8">
      <div className="flex items-center gap-3 mb-6">
        <h1 className="text-2xl font-bold text-gray-900">รายการนัดหมายที่เสร็จสิ้น</h1>
        <span className="rounded-full bg-[#F26522] text-white text-xs font-semibold px-3 py-1">
          {totalText}
        </span>
      </div>

      <hr className="border-gray-200 mb-6" />

      {loading ? (
        <div className="p-6 text-center text-gray-500">กำลังโหลดข้อมูล...</div>
      ) : error ? (
        <div className="p-4 text-sm text-orange-600 bg-orange-50 rounded mb-4">
          {error}
          <div className="mt-2 text-xs text-orange-700">
            ถ้า error คือ “advisor_id is required” ให้เช็คว่า token มี{" "}
            <span className="font-mono">user_id</span> ไหม (เพราะหน้านี้ส่ง{" "}
            <span className="font-mono">advisor_id</span> อัตโนมัติจาก token)
          </div>
        </div>
      ) : cards.length === 0 ? (
        <div className="p-6 text-center text-gray-500">ยังไม่มีนัดหมายที่เสร็จสิ้น</div>
      ) : null}

      {/* ✅ GRID: 1 / 2 / 3 columns */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        {cards.map((c) => {
          const appt = c.appointment;
          const log = c.advisorLog;

          return (
            <div
              key={appt.id}
              className="h-full border border-gray-300 rounded-2xl px-4 py-4 flex flex-col bg-white"
            >
              {/* header */}
              <div className="flex justify-between items-start gap-2">
                <div className="min-w-0">
                  <p className="text-sm font-semibold text-gray-900 truncate">
                    {appt.student.fullName || "ไม่ทราบชื่อ"}{" "}
                    {appt.student.sutId ? `(${appt.student.sutId})` : ""}
                  </p>
                  <p className="text-xs text-gray-500 line-clamp-1">
                    ขอปรึกษาเรื่อง {appt.topic || "-"}
                  </p>
                </div>

                <span
                  className={`shrink-0 rounded-full text-xs font-semibold px-3 py-1 ${badgeStyle(
                    c.ui.badgeTone
                  )}`}
                >
                  {c.ui.badgeText}
                </span>
              </div>

              {/* body */}
              <div className="mt-3 text-xs text-gray-700 space-y-1">
                {appt.startAt ? (
                  <p className="truncate">🕒 ปรึกษาเมื่อ: {formatDateMaybe(appt.startAt)}</p>
                ) : null}
                {appt.endAt ? (
                  <p className="truncate">✅ เสร็จเมื่อ: {formatDateMaybe(appt.endAt)}</p>
                ) : null}

                {log?.fileName ? (
                  <p className="text-xs text-[#F26522] truncate">ไฟล์: {log.fileName}</p>
                ) : null}
              </div>

              {/* footer */}
              <div className="mt-auto pt-4 flex gap-2">
                <button
                  onClick={() => router.push(c.ui.primaryAction.href)}
                  className="rounded-full bg-[#F26522] text-white text-xs font-semibold px-4 py-1.5 hover:bg-[#d1490e] transition"
                >
                  {c.ui.primaryAction.label}
                </button>

                {c.ui.secondaryAction ? (
                  <button
                    onClick={() => router.push(c.ui.secondaryAction!.href)}
                    className="rounded-full border border-[#F26522] text-[#F26522] text-xs font-semibold px-4 py-1.5 hover:bg-orange-50 transition"
                  >
                    {c.ui.secondaryAction.label}
                  </button>
                ) : null}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  </div>
);
}
