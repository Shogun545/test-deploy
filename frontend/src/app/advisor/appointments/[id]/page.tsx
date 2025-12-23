"use client";

import { use, useEffect, useState } from "react";
import { useRouter } from "next/navigation";

import {
  getAppointmentById,
  approveAppointment as approveAppointmentApi,
  rescheduleAppointment as rescheduleAppointmentApi,
  type AppointmentDetail,
  type AppointmentStatus,
} from "@/src/services/http/approvalservice";


type AppointmentDetailPageProps = {
  params: Promise<{ id: string }>;
};

const getStatusLabel = (s: AppointmentStatus) => {
  switch (s) {
    case "PENDING":
      return "รอพิจารณา";
    case "APPROVED":
      return "อนุมัติแล้ว";
    case "RESCHEDULE":
      return "เสนอเวลาใหม่";
    default:
      return s;
  }
};

export default function AppointmentDetailPage({
  params,
}: AppointmentDetailPageProps) {
  const router = useRouter();

  // ✅ unwrap Promise ด้วย React.use()
  const { id: idParam } = use(params);

  const idNum = Number(idParam);
  const isValidId = Number.isFinite(idNum) && idNum > 0;

  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [data, setData] = useState<AppointmentDetail | null>(null);

  // -----------------------
  // load appointment
  // -----------------------
  const loadAppointment = async () => {
    setLoading(true);
    setError(null);

    try {
      if (!isValidId) throw new Error("invalid id");

      const detail = await getAppointmentById(idNum);
      setData(detail);
    } catch (e: any) {
      const status = e?.response?.status;

      if (status === 401) {
        router.push("/login");
        return;
      }

      setError(e?.message || "เกิดข้อผิดพลาด");
      setData(null);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadAppointment();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [idParam]);

  // -----------------------
  // actions
  // -----------------------
  const approve = async () => {
    if (!data) return;

    setSubmitting(true);
    setError(null);

    try {
      await approveAppointmentApi(data.id, "อนุมัติคำขอ");
      await loadAppointment();
    } catch (e: any) {
      const status = e?.response?.status;

      if (status === 401) {
        router.push("/login");
        return;
      }

      setError(e?.message || "อนุมัติไม่สำเร็จ");
    } finally {
      setSubmitting(false);
    }
  };

  const reschedule = async () => {
    if (!data) return;

    setSubmitting(true);
    setError(null);

    try {
      await rescheduleAppointmentApi(data.id, "เสนอเวลาใหม่");
      await loadAppointment();
    } catch (e: any) {
      const status = e?.response?.status;

      if (status === 401) {
        router.push("/login");
        return;
      }

      setError(e?.message || "เสนอเวลาใหม่ไม่สำเร็จ");
    } finally {
      setSubmitting(false);
    }
  };

  // -----------------------
  // render
  // -----------------------
  if (loading) {
    return <div className="p-6 text-sm">กำลังโหลดข้อมูล...</div>;
  }

  if (!data) {
    return (
      <div className="p-6 text-sm text-red-500">
        ไม่พบข้อมูลคำขอนัดหมาย
        {error ? <div className="mt-2 text-xs text-red-600">{error}</div> : null}
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-[#E3743D]">
          รายละเอียดคำขอนัดหมาย
        </h1>
        <span className="rounded-full bg-[#FFF1E6] px-4 py-1 text-xs font-semibold text-[#E3743D]">
          หมายเลขคำขอ #{data.id}
        </span>
      </div>

      {error && (
        <div className="rounded-lg bg-red-50 px-4 py-2 text-sm text-red-600">
          {error}
        </div>
      )}

      <div className="rounded-2xl border bg-white p-6">
        <div className="mb-4 flex justify-between">
          <span className="text-sm font-semibold text-slate-500">
            สถานะปัจจุบัน:
          </span>
          <span className="rounded-full border px-4 py-1 text-xs font-semibold text-[#E3743D]">
            {getStatusLabel(data.status)} ({data.status})
          </span>
        </div>

        <div className="grid gap-6 text-sm md:grid-cols-2">
          <div>
            <p>
              <b>ชื่อนักศึกษา:</b> {data.studentName}
            </p>
            <p>
              <b>รหัสนักศึกษา:</b> {data.sutId}
            </p>
            <p>
              <b>วันที่ยื่น:</b> {data.submittedAt}
            </p>
          </div>
          <div>
            <p>
              <b>หัวข้อ:</b> {data.topic}
            </p>
            <p>
              <b>รายละเอียด:</b> {data.description}
            </p>
          </div>
        </div>

        <div className="mt-6 flex justify-between border-t pt-4">
          <button
            onClick={() => router.push("/advisor/appointments")}
            className="rounded-full border px-4 py-2 text-xs"
          >
            กลับไปหน้ารายการ
          </button>

          <div className="flex gap-3">
            <button
              disabled={submitting || data.status !== "PENDING"}
              onClick={approve}
              className="rounded-full border border-[#E3743D] px-4 py-2 text-xs text-[#E3743D] disabled:opacity-50"
            >
              อนุมัติคำขอ
            </button>

            <button
              disabled={submitting || data.status !== "PENDING"}
              onClick={reschedule}
              className="rounded-full bg-[#E3743D] px-4 py-2 text-xs text-white disabled:opacity-50"
            >
              เสนอเวลาใหม่
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
