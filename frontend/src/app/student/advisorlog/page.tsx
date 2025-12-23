"use client";

import { useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { ListAdvisorLogsByStudent } from "@/src/services/http/advisorlog";

// Type Definition
type LogStatus = "Draft" | "PendingReport" | "Completed" | string;

type AdvisorLogItem = {
  id: number;
  appointmentId?: number;
  title?: string;
  body?: string;
  status?: LogStatus;
  requiresReport?: boolean;
  fileName?: string;
  filePath?: string;
  createdAt?: string;
  updatedAt?: string;
};

// --- Helpers ---

// Format วันที่ภาษาไทย
function formatDate(s?: string) {
  if (!s) return "";
  try {
    return new Date(s).toLocaleString("th-TH", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  } catch {
    return s;
  }
}

// แยกชื่อไฟล์จาก CSV
function splitCsv(s?: string) {
  if (!s) return [];
  return s.split(",").map((x) => x.trim()).filter(Boolean);
}

// --- JWT Helper (Simplified) ---
function decodeJwtPayload(token: string): any | null {
  try {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(window.atob(base64).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));
    return JSON.parse(jsonPayload);
  } catch (e) {
    return null;
  }
}

// Helper เลือกค่าจาก Object (เผื่อ Case Sensitivity)
function pick<T>(obj: any, keys: string[], fallback: T): T {
  for (const k of keys) {
    const v = obj?.[k];
    if (v !== undefined && v !== null) return v as T;
  }
  return fallback;
}

export default function StudentAdvisorLogPage() {
  const router = useRouter();

  const [logs, setLogs] = useState<AdvisorLogItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const totalText = useMemo(
    () => (loading ? "..." : `${logs.length} รายการ`),
    [loading, logs.length]
  );

  useEffect(() => {
    let mounted = true;

    (async () => {
      setLoading(true);
      setError(null);

      try {
        // 1. Get Token
        const token = typeof window !== "undefined" ? localStorage.getItem("token") || localStorage.getItem("access_token") : null;
        
        if (!token) {
          router.push("/login");
          return;
        }

        // 2. Decode JWT to get User ID
        const payload = decodeJwtPayload(token);
        const userId = Number(payload?.user_id ?? payload?.userId ?? payload?.uid ?? 0);

        if (!userId) {
          throw new Error("Invalid Token: user_id not found");
        }

        // 3. Fetch Logs
        const res = await ListAdvisorLogsByStudent(userId);
        if (!res.success) throw new Error(res.message || "โหลดบันทึกไม่สำเร็จ");

        let arr: any[] = res.data?.data ?? res.data ?? [];
        if (!Array.isArray(arr)) arr = arr ? [arr] : [];

        // 4. Map Data
        const mapped: AdvisorLogItem[] = arr.map((it: any) => ({
          id: Number(pick(it, ["ID", "id"], 0)) || 0,
          appointmentId: pick(it, ["AppointmentID", "appointmentId"], undefined),
          title: pick(it, ["Title", "title"], ""),
          body: pick(it, ["Body", "body"], ""),
          status: pick(it, ["Status", "status"], ""),
          requiresReport: Boolean(pick(it, ["RequiresReport", "requiresReport"], false)),
          fileName: pick(it, ["FileName", "fileName"], "") || undefined,
          filePath: pick(it, ["FilePath", "filePath"], "") || undefined,
          createdAt: pick(it, ["CreatedAt", "createdAt"], "") || undefined,
          updatedAt: pick(it, ["UpdatedAt", "updatedAt"], "") || undefined,
        })).filter((x) => x.id > 0 && x.status !== "Draft"); // กรอง Draft ออก (เผื่อ Backend หลุดมา)

        // Sort: ใหม่สุดขึ้นก่อน
        mapped.sort((a, b) => (b.createdAt || "").localeCompare(a.createdAt || ""));

        if (mounted) setLogs(mapped);

      } catch (e: any) {
        if (mounted) {
            // Check 401 Unauthorized
            if (e.response?.status === 401 || e.message?.includes("Invalid Token")) {
                localStorage.removeItem("token");
                router.push("/login");
            } else {
                setError(e?.message ?? "เกิดข้อผิดพลาดในการเชื่อมต่อ");
                setLogs([]);
            }
        }
      } finally {
        if (mounted) setLoading(false);
      }
    })();

    return () => {
      mounted = false;
    };
  }, [router]);

  // UI Helpers
  const badgeClass = (st?: string) => {
    switch (st) {
        case "PendingReport": return "bg-orange-100 text-orange-700 border border-orange-200";
        case "Completed": return "bg-green-100 text-green-700 border border-green-200";
        default: return "bg-gray-100 text-gray-600 border border-gray-200";
    }
  };

  const statusLabel = (st?: string) => {
    switch (st) {
        case "PendingReport": return "รอส่งรายงาน";
        case "Completed": return "เสร็จสิ้น";
        default: return st || "-";
    }
  };

  return (
    <div className="w-full min-h-screen bg-gray-50 flex justify-center py-6 px-4">
      <div className="w-full max-w-6xl bg-white rounded-3xl shadow-sm border border-gray-200 px-6 py-8 sm:px-10">
        
        {/* Header */}
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center gap-3">
            <h1 className="text-2xl font-bold text-gray-900">
              บันทึกการปรึกษาของฉัน
            </h1>
            <span className="rounded-full bg-[#F26522] text-white text-xs font-semibold px-3 py-1">
              {totalText}
            </span>
          </div>
        </div>

        <hr className="border-gray-200 mb-6" />

        {/* State Display */}
        {loading ? (
          <div className="p-10 text-center text-gray-500 animate-pulse">กำลังโหลดข้อมูล...</div>
        ) : error ? (
          <div className="p-4 text-sm text-red-600 bg-red-50 rounded-xl border border-red-100 text-center">
            {error}
          </div>
        ) : logs.length === 0 ? (
          <div className="p-10 text-center text-gray-500 border border-dashed border-gray-300 rounded-xl">
            ยังไม่มีรายการบันทึกการปรึกษา
          </div>
        ) : (
          /* Cards Grid */
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {logs.map((log) => {
              const fileNames = splitCsv(log.fileName);

              return (
                <div
                  key={log.id}
                  className="h-full border border-gray-200 hover:border-orange-300 hover:shadow-md transition rounded-2xl p-5 flex flex-col bg-white"
                >
                  {/* Card Header */}
                  <div className="flex justify-between items-start gap-2 mb-3">
                    <div className="min-w-0">
                      <h3 className="text-base font-bold text-gray-900 truncate" title={log.title}>
                        {log.title || "ไม่มีหัวข้อ"}
                      </h3>
                      <p className="text-xs text-gray-500 mt-1">
                        {formatDate(log.createdAt)}
                      </p>
                    </div>
                    <span className={`shrink-0 rounded-full text-[10px] font-bold px-2 py-1 uppercase tracking-wide ${badgeClass(log.status)}`}>
                      {statusLabel(log.status)}
                    </span>
                  </div>

                  {/* Card Body */}
                  <div className="flex-1">
                    <p className="text-sm text-gray-600 line-clamp-3 mb-3">
                        {log.body || "- ไม่มีรายละเอียด -"}
                    </p>

                    {/* Meta Tags */}
                    <div className="flex flex-wrap gap-2 mb-4">
                        {fileNames.length > 0 && (
                            <span className="inline-flex items-center gap-1 text-xs text-orange-600 bg-orange-50 px-2 py-1 rounded-md">
                                📎 {fileNames.length} ไฟล์แนบ
                            </span>
                        )}
                        {log.requiresReport && (
                            <span className="inline-flex items-center gap-1 text-xs text-blue-600 bg-blue-50 px-2 py-1 rounded-md">
                                📝 ต้องส่งรายงาน
                            </span>
                        )}
                    </div>
                  </div>

                  {/* Card Footer */}
                  <div className="mt-auto pt-3 border-t border-gray-100">
                    <button
                      onClick={() => router.push(`/student/advisorlog/${log.id}`)}
                      className="w-full rounded-xl bg-[#F26522] text-white text-sm font-bold px-4 py-2 hover:bg-[#d1490e] transition flex items-center justify-center gap-2"
                    >
                      ดูรายละเอียด
                      <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                      </svg>
                    </button>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}