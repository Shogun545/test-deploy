"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import { useRouter } from "next/navigation";

import {
  listAppointmentsPending,
  listAppointmentsDone,
  listAppointmentsAll,
  type AppointmentListRow,
  type AppointmentStatus,
} from "@/src/services/http/approvalservice";

type ViewMode = "PENDING" | "DONE" | "ALL";

// ---------- helper functions ----------
const getStatusLabel = (status: AppointmentStatus) => {
  switch (status) {
    case "PENDING":
      return "รอพิจารณา";
    case "APPROVED":
      return "อนุมัติแล้ว";
    case "RESCHEDULE":
      return "เสนอเวลาใหม่";
    default:
      return status;
  }
};

const getStatusClassName = (status: AppointmentStatus) => {
  switch (status) {
    case "PENDING":
      return "bg-white text-[#E3743D] border border-[#E3743D]";
    case "APPROVED":
      return "bg-white text-[#E3743D] border border-[#E3743D]";
    case "RESCHEDULE":
      return "bg-[#E3743D] text-white border border-white";
    default:
      return "bg-slate-100 text-slate-700";
  }
};

export default function AdvisorRequestsPage() {
  const router = useRouter();

  const [viewMode, setViewMode] = useState<ViewMode>("PENDING");
  const [requests, setRequests] = useState<AppointmentListRow[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  // กัน race ตอนสลับแท็บเร็ว ๆ
  const abortRef = useRef<AbortController | null>(null);

  const load = async (mode: ViewMode) => {
    setLoading(true);

    // abort request เก่า (กันสลับแท็บเร็ว)
    abortRef.current?.abort();
    const controller = new AbortController();
    abortRef.current = controller;

    try {
      // NOTE: axios (apiClient) ไม่รองรับ AbortController signal แบบ fetch ตรง ๆ
      // เราใช้ abort เพื่อกัน state update หลังเปลี่ยนแท็บเป็นหลัก (เช็คด้านล่าง)

      let rows: AppointmentListRow[] = [];
      if (mode === "PENDING") rows = await listAppointmentsPending();
      else if (mode === "DONE") rows = await listAppointmentsDone();
      else rows = await listAppointmentsAll();

      // ถ้า request นี้ถูก abort ไปแล้ว ไม่ต้อง setState
      if (controller.signal.aborted) return;

      setRequests(rows);
    } catch (e: any) {
      if (controller.signal.aborted) return;

      // ถ้า apiClient ของส้มโยน error แบบ axios: err.response.status
      const status = e?.response?.status;

      if (status === 401) {
        setRequests([]);
        setLoading(false);
        router.push("/login");
        return;
      }

      console.error("❌ load appointments failed:", e);
      setRequests([]);
    } finally {
      if (!controller.signal.aborted) setLoading(false);
    }
  };

  useEffect(() => {
    if (typeof window === "undefined") return;
    load(viewMode);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [viewMode]);

  useEffect(() => {
    if (typeof window === "undefined") return;
    const onFocus = () => load(viewMode);
    window.addEventListener("focus", onFocus);
    return () => window.removeEventListener("focus", onFocus);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [viewMode]);

  const filteredRequests = useMemo(() => requests, [requests]);

  const headerTitle =
    viewMode === "PENDING"
      ? "รายการคำขอที่รอพิจารณา"
      : viewMode === "DONE"
      ? "รายการคำขอที่พิจารณาแล้ว"
      : "รายการคำขอทั้งหมด";

  return (
    <div className="space-y-6">
      <div className="mb-2 flex items-center justify-between gap-4">
        <div className="flex items-baseline gap-2">
          <h2 className="text-xl font-bold text-[#E3743D]">{headerTitle}</h2>
          <span className="rounded-full bg-amber-50 px-3 py-1 text-xs font-semibold text-amber-700">
            {filteredRequests.length} รายการ
          </span>
        </div>
      </div>

      {/* ปุ่มสลับมุมมอง */}
      <div className="flex flex-wrap items-center gap-2">
        {(["PENDING", "DONE", "ALL"] as ViewMode[]).map((m) => (
          <button
            key={m}
            onClick={() => setViewMode(m)}
            className={`rounded-full px-4 py-2 text-xs font-semibold border ${
              viewMode === m
                ? "bg-[#E3743D] text-white border-[#E3743D]"
                : "bg-white text-[#E3743D] border-[#E3743D]"
            }`}
          >
            {m === "PENDING" ? "รอพิจารณา" : m === "DONE" ? "พิจารณาแล้ว" : "ทั้งหมด"}
          </button>
        ))}
      </div>

      {loading && (
        <div className="rounded-xl border border-slate-100 bg-white p-4 text-sm text-slate-600">
          กำลังโหลดข้อมูล...
        </div>
      )}

      <div className="overflow-x-auto rounded-xl border border-slate-100 bg-white">
        <table className="min-w-full border-collapse text-sm">
          <thead className="bg-slate-50">
            <tr>
              <th className="px-4 py-3 text-left text-xs font-semibold text-slate-500">
                ลำดับ
              </th>
              <th className="px-4 py-3 text-left text-xs font-semibold text-slate-500">
                ชื่อ - รหัสนักศึกษา
              </th>
              <th className="px-4 py-3 text-left text-xs font-semibold text-slate-500">
                หัวข้อคำขอ
              </th>
              <th className="px-4 py-3 text-left text-xs font-semibold text-slate-500">
                วันที่ยื่นคำขอ
              </th>
              <th className="px-4 py-3 text-left text-xs font-semibold text-slate-500">
                สถานะ
              </th>
              <th className="px-4 py-3 text-right text-xs font-semibold text-slate-500">
                การดำเนินการ
              </th>
            </tr>
          </thead>

          <tbody>
            {filteredRequests.map((req, index) => (
              <tr
                key={req.id}
                className="border-t border-slate-100 hover:bg-slate-50/60"
              >
                <td className="px-4 py-3">{index + 1}</td>

                <td className="px-4 py-3">
                  <div className="flex flex-col">
                    <span className="font-medium">{req.studentName}</span>
                    <span className="text-xs text-slate-500">{req.sutId}</span>
                  </div>
                </td>

                <td className="px-4 py-3">{req.topic}</td>
                <td className="px-4 py-3">{req.submittedAt}</td>

                <td className="px-4 py-3">
                  <span
                    className={
                      "inline-flex rounded-full px-3 py-1 text-xs font-semibold " +
                      getStatusClassName(req.status)
                    }
                  >
                    {getStatusLabel(req.status)}
                  </span>
                </td>

                <td className="px-4 py-3 text-right">
                  <button
                    onClick={() => router.push(`/advisor/appointments/${req.id}`)}
                    className="rounded-full bg-[#E3743D] px-4 py-2 text-xs font-semibold text-white hover:bg-[#c65d2e]"
                  >
                    ดูรายละเอียด
                  </button>
                </td>
              </tr>
            ))}

            {!loading && filteredRequests.length === 0 && (
              <tr>
                <td
                  colSpan={6}
                  className="px-4 py-6 text-center text-sm text-slate-500"
                >
                  {viewMode === "PENDING"
                    ? "ยังไม่มีคำขอที่รอพิจารณา"
                    : viewMode === "DONE"
                    ? "ยังไม่มีคำขอที่พิจารณาแล้ว"
                    : "ยังไม่มีคำขอ"}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
