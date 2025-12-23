"use client";

import { useEffect, useMemo, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { GetAdvisorLogByID } from "@/src/services/http/advisorlog";
import { apiClient } from "@/src/services/http/apiclient";

// กำหนด Type ให้ตรงกับข้อมูลที่เราจะใช้งานในหน้าเว็บ
type AdvisorLogItem = {
  id: number;
  appointmentId: number;
  title: string;
  body: string;
  status: string;
  requiresReport: boolean;
  fileName?: string;
  filePath?: string;
  createdAt?: string;
  updatedAt?: string;
};

// Helper: แยก CSV (สำหรับชื่อไฟล์และ Path)
function splitCsv(s?: string) {
  if (!s) return [];
  return s.split(",").map((x) => x.trim()).filter(Boolean);
}

// Helper: จัด Format วันที่สวยๆ (ภาษาไทย)
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

// Helper: สีของ Badge สถานะ
function getStatusColor(status: string) {
  switch (status) {
    case "Draft":
      return "bg-gray-100 text-gray-600 border-gray-200";
    case "PendingReport":
      return "bg-orange-100 text-orange-700 border-orange-200";
    case "Completed":
      return "bg-green-100 text-green-700 border-green-200";
    default:
      return "bg-gray-100 text-gray-600 border-gray-200";
  }
}

// Helper: ข้อความสถานะภาษาไทย
function getStatusLabel(status: string) {
  switch (status) {
    case "Draft": return "ฉบับร่าง";
    case "PendingReport": return "รอส่งรายงาน";
    case "Completed": return "บันทึกแล้ว";
    default: return status;
  }
}

export default function AdvisorLogDetailPage() {
  const router = useRouter();
  const params = useParams();
  const id = params?.id;

  const [loading, setLoading] = useState(true);
  const [log, setLog] = useState<AdvisorLogItem | null>(null);
  const [error, setError] = useState<string | null>(null);

  // State สำหรับ Loading ตอนกดโหลดไฟล์ (UX) เพื่อกัน User กดย้ำ
  const [downloadingIndex, setDownloadingIndex] = useState<number | null>(null);

  const numericId = useMemo(() => {
    const v = Array.isArray(id) ? id[0] : id;
    const n = Number(v);
    return Number.isFinite(n) && n > 0 ? n : null;
  }, [id]);

  useEffect(() => {
    let mounted = true;

    const load = async () => {
      if (!numericId) {
        setError("ไม่พบรหัสเอกสาร (Invalid ID)");
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
          // Data Mapping: รองรับทั้ง camelCase และ PascalCase เพื่อความชัวร์
          const d = res.data?.data ?? res.data ?? {};
          setLog({
            id: Number(d.id ?? d.ID),
            appointmentId: Number(d.appointmentId ?? d.AppointmentID),
            title: d.title ?? d.Title ?? "",
            body: d.body ?? d.Body ?? "",
            status: d.status ?? d.Status ?? "Draft",
            requiresReport: !!(d.requiresReport ?? d.RequiresReport),
            fileName: d.fileName ?? d.FileName,
            filePath: d.filePath ?? d.FilePath,
            createdAt: d.createdAt ?? d.CreatedAt,
            updatedAt: d.updatedAt ?? d.UpdatedAt,
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

  // ฟังก์ชันเปิดไฟล์ (แก้ไข Memory Leak เรียบร้อย)
  async function openFile(logId: number, index: number, fileName?: string) {
    if (downloadingIndex !== null) return; // กันกดซ้ำ
    setDownloadingIndex(index);

    try {
      const resp = await apiClient.get(`/api/advisor_logs/${logId}/files/${index}`, {
        responseType: "blob", // สำคัญมาก
      });

      // ✅ สร้าง Blob URL
      const blobUrl = URL.createObjectURL(resp.data);

      // ✅ สร้าง Link จำลองเพื่อสั่ง Download/Open
      const link = document.createElement("a");
      link.href = blobUrl;
      link.target = "_blank"; 
      link.download = fileName || `file-${logId}-${index}`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);

      // ✅ Cleanup Memory: คืน RAM ให้ Browser
      setTimeout(() => {
        URL.revokeObjectURL(blobUrl);
      }, 100);

    } catch (err) {
      alert("ไม่สามารถเปิดไฟล์ได้");
      console.error(err);
    } finally {
      setDownloadingIndex(null);
    }
  }

  return (
    <div className="w-full min-h-screen bg-gray-50 flex justify-center py-6 px-4">
      <div className="w-full max-w-4xl bg-white rounded-3xl shadow-sm border border-gray-200 px-6 py-8 sm:px-10">
        
        {/* Header Section */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
          <div className="flex flex-wrap items-center gap-3">
            <h1 className="text-2xl font-bold text-gray-900">
              รายละเอียดบันทึก
            </h1>
            {log && (
              <span className={`px-3 py-1 rounded-full text-xs font-bold border ${getStatusColor(log.status)}`}>
                {getStatusLabel(log.status)}
              </span>
            )}
          </div>

          <button
            onClick={() => router.push("/advisor/advisorlog")}
            className="self-start sm:self-auto rounded-full border border-gray-300 text-gray-700 text-sm font-semibold px-5 py-2 hover:bg-gray-50 transition"
          >
            ← กลับไปหน้า My Logs
          </button>
        </div>

        <hr className="border-gray-200 mb-6" />

        {/* Content Section */}
        {loading ? (
          <div className="p-10 text-center text-gray-500 animate-pulse">กำลังโหลดข้อมูล...</div>
        ) : error ? (
          <div className="p-4 text-sm text-red-600 bg-red-50 rounded-xl border border-red-100 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
            </svg>
            {error}
          </div>
        ) : !log ? (
          <div className="p-10 text-center text-gray-500">ไม่พบข้อมูลบันทึก</div>
        ) : (
          <div className="space-y-8">
            
            {/* Meta Info Grid */}
            <div className="bg-gray-50 rounded-2xl p-5 border border-gray-100">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <p className="text-xs text-gray-500 uppercase tracking-wide font-semibold mb-1">หัวข้อการปรึกษา</p>
                  <p className="text-lg font-bold text-gray-900">{log.title}</p>
                </div>
                <div className="md:text-right">
                  <p className="text-xs text-gray-500 uppercase tracking-wide font-semibold mb-1">วันที่บันทึก</p>
                  <p className="text-gray-900 font-medium">{formatDate(log.createdAt)}</p>
                  <p className="text-xs text-gray-400 mt-1 font-mono">Log ID: {log.id} | Appt ID: {log.appointmentId}</p>
                </div>
              </div>
            </div>

            {/* Body Description */}
            <div>
              <h3 className="text-sm font-bold text-gray-900 mb-3 uppercase tracking-wider flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 text-gray-500" viewBox="0 0 20 20" fill="currentColor">
                  <path fillRule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4z" clipRule="evenodd" />
                </svg>
                รายละเอียด
              </h3>
              <div className="bg-white border border-gray-200 rounded-2xl p-6 shadow-sm">
                <pre className="whitespace-pre-wrap break-words text-sm text-gray-800 font-sans leading-relaxed">
                  {log.body || "- ไม่มีการระบุรายละเอียด -"}
                </pre>
              </div>
            </div>

            {/* Attachments */}
            <div>
              <h3 className="text-sm font-bold text-gray-900 mb-3 uppercase tracking-wider flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 text-gray-500" viewBox="0 0 20 20" fill="currentColor">
                  <path fillRule="evenodd" d="M8 4a3 3 0 00-3 3v4a5 5 0 0010 0V7a1 1 0 112 0v4a7 7 0 11-14 0V7a5 5 0 0110 0v4a3 3 0 11-6 0V7a1 1 0 012 0v4a1 1 0 102 0V7a3 3 0 00-3-3z" clipRule="evenodd" />
                </svg>
                ไฟล์แนบ <span className="text-gray-400 font-normal">({fileNames.length})</span>
              </h3>

              {fileNames.length === 0 ? (
                <div className="p-6 border border-dashed border-gray-300 rounded-xl text-center text-gray-400 text-sm italic bg-gray-50">
                  ไม่มีไฟล์แนบ
                </div>
              ) : (
                <ul className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                  {fileNames.map((name, idx) => (
                    <li
                      key={`${idx}-${name}`}
                      className="group flex items-center justify-between p-3 border border-gray-200 rounded-xl hover:border-orange-300 hover:shadow-md transition bg-white"
                    >
                      <div className="flex items-center gap-3 overflow-hidden min-w-0">
                        <div className="bg-orange-50 p-2 rounded-lg text-orange-500">
                            {/* Icon File */}
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                                <path fillRule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4z" clipRule="evenodd" />
                            </svg>
                        </div>
                        <span className="truncate text-sm text-gray-700 font-medium" title={name}>
                          {name}
                        </span>
                      </div>

                      <button
                        disabled={downloadingIndex !== null}
                        onClick={() => openFile(log.id, idx, name)}
                        className="shrink-0 text-[#F26522] text-sm font-semibold px-4 py-2 rounded-lg hover:bg-orange-50 disabled:opacity-50 transition flex items-center gap-1"
                      >
                        {downloadingIndex === idx ? (
                           <span className="animate-pulse">Loading...</span>
                        ) : "เปิดดู"}
                      </button>
                    </li>
                  ))}
                </ul>
              )}
            </div>

            {/* Footer Action (สำหรับ Draft เท่านั้น) */}
            {log.status === "Draft" && (
              <div className="pt-6 mt-8 border-t border-gray-100 flex justify-end">
                <button
                  onClick={() => router.push(`/advisor/advisorlog/edit/${log.id}`)}
                  className="rounded-full bg-[#F26522] text-white text-sm font-bold px-8 py-3 hover:bg-[#d1490e] shadow-lg shadow-orange-200 transition transform hover:-translate-y-0.5 flex items-center gap-2"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                    <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
                  </svg>
                  แก้ไขฉบับร่าง
                </button>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}