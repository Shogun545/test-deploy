"use client";

import { useEffect, useMemo, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { GetAdvisorLogByID } from "@/src/services/http/advisorlog";
import { apiClient } from "@/src/services/http/apiclient";

// Helper: Format วันที่ไทย
function formatDate(s?: string) {
  if (!s) return "-";
  try {
    return new Date(s).toLocaleString("th-TH", {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  } catch {
    return s;
  }
}

// Helper: แยก CSV ชื่อไฟล์
function splitCsv(s?: string) {
  if (!s) return [];
  return s.split(",").map((x) => x.trim()).filter(Boolean);
}

// Helper: สี Badge สถานะ
function getBadgeClass(st?: string) {
  switch (st) {
    case "Draft": return "bg-gray-100 text-gray-600 border border-gray-200";
    case "PendingReport": return "bg-orange-100 text-orange-700 border border-orange-200";
    case "Completed": return "bg-green-100 text-green-700 border border-green-200";
    default: return "bg-gray-100 text-gray-600 border border-gray-200";
  }
}

// Helper: ข้อความสถานะภาษาไทย
function getStatusLabel(st?: string) {
  switch (st) {
    case "Draft": return "ฉบับร่าง";
    case "PendingReport": return "รอส่งรายงาน";
    case "Completed": return "เสร็จสิ้น";
    default: return st || "-";
  }
}

// Type Definition (ใช้ camelCase เพื่อความสม่ำเสมอใน JS)
type AdvisorLogItem = {
  id: number;
  appointmentId?: number;
  title?: string;
  body?: string;
  status?: string;
  requiresReport?: boolean;
  fileName?: string;
  filePath?: string;
  createdAt?: string;
  updatedAt?: string;
};

export default function StudentAdvisorLogDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params?.id;

  const [loading, setLoading] = useState(true);
  const [log, setLog] = useState<AdvisorLogItem | null>(null);
  const [error, setError] = useState<string | null>(null);
  
  // State กันกดปุ่มเปิดไฟล์ซ้ำ
  const [openingIndex, setOpeningIndex] = useState<number | null>(null);

  const numericId = useMemo(() => {
    const v = Array.isArray(id) ? id[0] : id;
    const n = Number(v);
    return Number.isFinite(n) && n > 0 ? n : null;
  }, [id]);

  useEffect(() => {
    let mounted = true;

    const load = async () => {
      if (!numericId) {
        setError("ไม่พบรหัสเอกสาร");
        setLoading(false);
        return;
      }

      setLoading(true);
      setError(null);

      try {
        const res = await GetAdvisorLogByID(numericId);
        if (!mounted) return;

        if (!res.success) {
          setError(res.message || "โหลดข้อมูลไม่สำเร็จ");
          setLog(null);
        } else {
          const d = res.data?.data ?? res.data ?? {};
          
          // Map Data เป็น camelCase
          setLog({
            id: Number(d.ID ?? d.id ?? numericId),
            appointmentId: d.AppointmentID ?? d.appointmentId,
            title: d.Title ?? d.title ?? "",
            body: d.Body ?? d.body ?? "",
            status: d.Status ?? d.status ?? "",
            requiresReport: !!(d.RequiresReport ?? d.requiresReport),
            fileName: d.FileName ?? d.fileName ?? "",
            filePath: d.FilePath ?? d.filePath ?? "",
            createdAt: d.CreatedAt ?? d.createdAt,
            updatedAt: d.UpdatedAt ?? d.updatedAt,
          });
        }
      } catch (err: any) {
        if (mounted) setError(err.message || "เกิดข้อผิดพลาดในการเชื่อมต่อ");
      } finally {
        if (mounted) setLoading(false);
      }
    };

    load();
    return () => {
      mounted = false;
    };
  }, [numericId]);

  const fileNames = splitCsv(log?.fileName);
  // เราใช้ index ในการอ้างอิงไฟล์ตาม Logic Backend

  // ✅ ฟังก์ชันเปิดไฟล์แบบเสถียร (ไม่โดน Block Popup + คืน Memory เร็ว)
  const openFile = async (index: number, fileName?: string) => {
    if (!log?.id || openingIndex !== null) return;
    setOpeningIndex(index);

    try {
      const resp = await apiClient.get(`/api/advisor_logs/${log.id}/files/${index}`, {
        responseType: "blob",
      });

      // สร้าง Blob URL
      const blobUrl = URL.createObjectURL(resp.data);

      // สร้าง Link จำลองเพื่อสั่งเปิด/ดาวน์โหลด
      const link = document.createElement("a");
      link.href = blobUrl;
      link.target = "_blank"; // พยายามเปิดแท็บใหม่
      link.download = fileName || `file-${log.id}-${index}`; // ช่วย Browser ตั้งชื่อไฟล์
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);

      // คืน Memory ทันทีหลังจาก Browser รับคำสั่งไปแล้วเล็กน้อย
      setTimeout(() => URL.revokeObjectURL(blobUrl), 100);

    } catch (e: any) {
      const msg = e?.response?.data?.error || e?.message || "เปิดไฟล์ไม่สำเร็จ";
      alert(msg);
    } finally {
      setOpeningIndex(null);
    }
  };

  return (
    <div className="w-full min-h-screen bg-gray-50 flex justify-center py-6 px-4">
      <div className="w-full max-w-4xl bg-white rounded-3xl shadow-sm border border-gray-200 px-6 py-8 sm:px-10">
        
        {/* Header */}
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-2xl font-bold text-gray-900">รายละเอียดบันทึก</h1>
          <button
            onClick={() => router.push("/student/advisorlog")}
            className="rounded-full border border-gray-300 text-gray-700 text-sm font-semibold px-5 py-2 hover:bg-gray-50 transition"
          >
            ← กลับ
          </button>
        </div>

        <hr className="border-gray-200 mb-6" />

        {/* Content */}
        {loading ? (
          <div className="p-10 text-center text-gray-500 animate-pulse">กำลังโหลดข้อมูล...</div>
        ) : error ? (
          <div className="p-4 text-sm text-red-600 bg-red-50 rounded-xl border border-red-100 text-center">
            {error}
          </div>
        ) : !log ? (
          <div className="p-10 text-center text-gray-500">ไม่พบบันทึก</div>
        ) : (
          <div className="space-y-8">
            
            {/* Meta Info */}
            <div className="bg-gray-50 rounded-2xl p-5 border border-gray-100 grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                    <div className="flex items-center gap-2 mb-2">
                         <span className={`px-3 py-1 rounded-full text-xs font-bold ${getBadgeClass(log.status)}`}>
                            {getStatusLabel(log.status)}
                         </span>
                         {log.requiresReport && (
                            <span className="px-3 py-1 rounded-full text-xs font-bold bg-blue-100 text-blue-700 border border-blue-200">
                                ต้องส่งรายงาน
                            </span>
                         )}
                    </div>
                    <div className="text-sm text-gray-500">Log ID: {log.id}</div>
                    <div className="text-sm text-gray-500">Appt ID: {log.appointmentId ?? "-"}</div>
                </div>
                <div className="md:text-right">
                    <p className="text-xs text-gray-500 uppercase tracking-wide font-semibold mb-1">หัวข้อ</p>
                    <p className="text-lg font-bold text-gray-900 mb-2">{log.title}</p>
                    <p className="text-xs text-gray-500">สร้างเมื่อ: {formatDate(log.createdAt)}</p>
                    {log.updatedAt && log.updatedAt !== log.createdAt && (
                        <p className="text-xs text-gray-400">อัปเดต: {formatDate(log.updatedAt)}</p>
                    )}
                </div>
            </div>

            {/* Body */}
            <div>
              <h3 className="text-sm font-bold text-gray-900 mb-3 uppercase tracking-wider">รายละเอียด</h3>
              <div className="bg-white border border-gray-200 rounded-2xl p-6 shadow-sm">
                 <pre className="whitespace-pre-wrap break-words text-sm text-gray-800 font-sans leading-relaxed">
                    {log.body || "- ไม่มีรายละเอียด -"}
                 </pre>
              </div>
            </div>

            {/* Files */}
            <div>
              <h3 className="text-sm font-bold text-gray-900 mb-3 uppercase tracking-wider flex items-center gap-2">
                ไฟล์แนบ <span className="text-gray-400 font-normal">({fileNames.length})</span>
              </h3>

              {fileNames.length === 0 ? (
                <div className="p-6 border border-dashed border-gray-300 rounded-xl text-center text-gray-400 text-sm italic bg-gray-50">
                  ไม่มีไฟล์แนบ
                </div>
              ) : (
                <ul className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                  {fileNames.map((name, idx) => (
                    <li key={`${name}-${idx}`} className="flex items-center justify-between p-3 border border-gray-200 rounded-xl hover:border-orange-300 hover:shadow-sm transition bg-white">
                       <div className="flex items-center gap-3 overflow-hidden">
                          <span className="text-xl">📄</span>
                          <span className="truncate text-sm text-gray-700 font-medium" title={name}>
                            {name}
                          </span>
                       </div>
                       
                       <button
                         onClick={() => openFile(idx, name)}
                         disabled={openingIndex !== null}
                         className="shrink-0 text-[#F26522] text-sm font-semibold px-3 py-1.5 rounded-lg hover:bg-orange-50 disabled:opacity-50 transition"
                       >
                         {openingIndex === idx ? "กำลังเปิด..." : "เปิดไฟล์"}
                       </button>
                    </li>
                  ))}
                </ul>
              )}
            </div>

          </div>
        )}
      </div>
    </div>
  );
}