"use client";

import { ChangeEvent, useEffect, useMemo, useState } from "react";
import { useParams, useRouter } from "next/navigation";

import {
  GetAdvisorLogByID,
  UpdateAdvisorLog,
  UpdateAdvisorLogStatus,
} from "@/src/services/http/advisorlog";

type LogStatus = "Draft" | "PendingReport" | "Completed" | string;

export default function EditAdvisorLogPage() {
  const router = useRouter();
  const params = useParams();
  const id = params?.id;

  const numericId = useMemo(() => {
    const v = Array.isArray(id) ? id[0] : id;
    const n = Number(v);
    return Number.isFinite(n) && n > 0 ? n : null;
  }, [id]);

  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  const [error, setError] = useState<string | null>(null);
  const [formError, setFormError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);

  // form state
  const [title, setTitle] = useState("");
  const [body, setBody] = useState("");
  const [requiresReport, setRequiresReport] = useState<boolean | null>(null);
  const [status, setStatus] = useState<LogStatus>("");

  // existing files (from backend string "a.pdf,b.jpg")
  const [existingFileName, setExistingFileName] = useState<string>("");
  // const [existingFilePath, setExistingFilePath] = useState<string>(""); // ไม่ได้ใช้ path ในหน้านี้

  // new files to upload
  const [files, setFiles] = useState<File[]>([]);

  const goBack = () => router.push("/advisor/advisorlog");

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    const list = e.target.files;
    if (!list || list.length === 0) return;

    const selected = Array.from(list);
    setFiles((prev) => [...prev, ...selected]);
    e.target.value = "";
  };

  const removeNewFile = (idx: number) => {
    setFiles((prev) => prev.filter((_, i) => i !== idx));
  };

  useEffect(() => {
    let mounted = true;

    (async () => {
      if (!numericId) {
        setError("ไม่พบรหัสเอกสาร (Invalid ID)");
        setLoading(false);
        return;
      }

      setLoading(true);
      setError(null);

      const res = await GetAdvisorLogByID(numericId);
      if (!mounted) return;

      if (!res.success) {
        setError(res.message || "โหลดข้อมูลไม่สำเร็จ");
        setLoading(false);
        return;
      }

      // Mapping Data (รองรับทั้งตัวเล็ก/ตัวใหญ่)
      const payload = res.data?.data ?? res.data ?? {};
      
      const t = payload?.Title ?? payload?.title ?? "";
      const b = payload?.Body ?? payload?.body ?? "";
      const rr = payload?.RequiresReport ?? payload?.requiresReport ?? payload?.requires_report;
      const st = payload?.Status ?? payload?.status ?? "";
      const fn = payload?.FileName ?? payload?.fileName ?? payload?.file_name ?? "";
      // const fp = payload?.FilePath ?? payload?.filePath ?? payload?.file_path ?? "";

      setTitle(String(t));
      setBody(String(b));
      // เช็คให้ชัวร์ว่าเป็น boolean
      setRequiresReport(rr === true || rr === "true"); 
      setStatus(st);
      setExistingFileName(fn);
      // setExistingFilePath(fp);

      setLoading(false);
    })();

    return () => {
      mounted = false;
    };
  }, [numericId]);

  // Validate Function
  const validateForm = () => {
    if (!title.trim()) return "กรุณากรอกหัวข้อ";
    if (!body.trim()) return "กรุณากรอกเนื้อหาสรุป";
    if (requiresReport === null) return "กรุณาเลือก ต้องส่ง / ไม่ต้องส่ง รายงาน";
    return null;
  };

  const saveContent = async () => {
    if (!numericId || saving) return;

    const err = validateForm();
    if (err) {
      setFormError(err);
      return;
    }

    setFormError(null);
    setSaving(true);

    // Call API Update Content
    const res = await UpdateAdvisorLog(
      numericId,
      {
        Title: title.trim(),
        Body: body.trim(),
        RequiresReport: requiresReport === true,
      } as any, // Cast as any เพื่อความง่าย (หรือจะแก้ Interface ให้ตรงเป๊ะก็ได้)
      { files }
    );

    if (res.success) {
      setSuccessMsg("บันทึกการแก้ไขแล้ว");
      setFiles([]); // เคลียร์ไฟล์ใหม่ เพราะมันถูกเซฟไปแล้ว

      setTimeout(() => {
        // Refresh หน้าจอ หรือ Redirect
        // router.refresh(); 
        router.push("/advisor/advisorlog");
      }, 1500);
      return;
    }

    setFormError(res.message || "บันทึกไม่สำเร็จ");
    setSaving(false);
  };

  const publish = async () => {
    if (!numericId || saving) return;

    const err = validateForm();
    if (err) {
      setFormError(err);
      return;
    }

    setFormError(null);
    setSaving(true);

    // ✅ STEP 1: Save Content ก่อน
    const saveRes = await UpdateAdvisorLog(
      numericId,
      {
        Title: title.trim(),
        Body: body.trim(),
        RequiresReport: requiresReport === true,
      } as any,
      { files }
    );

    if (!saveRes.success) {
      setFormError(saveRes.message || "บันทึกข้อมูลไม่สำเร็จ (ขั้นตอนที่ 1/2)");
      setSaving(false);
      return;
    }

    // ✅ STEP 2: Update Status
    const nextStatus = requiresReport ? "PendingReport" : "Completed";
    const statusRes = await UpdateAdvisorLogStatus(numericId, nextStatus);

    if (statusRes.success) {
      setSuccessMsg("เผยแพร่สำเร็จ");
      setTimeout(() => router.push("/advisor/advisorlog"), 1500);
      return;
    }

    setFormError(statusRes.message || "เปลี่ยนสถานะไม่สำเร็จ (ขั้นตอนที่ 2/2)");
    setSaving(false);
  };

  const fileNames = useMemo(
    () => (existingFileName ? existingFileName.split(",").map((s: string) => s.trim()).filter(Boolean) : []),
    [existingFileName]
  );

  return (
    <div className="w-full min-h-screen bg-gray-50 flex justify-center py-6 px-4">
      <div className="w-full max-w-5xl bg-white p-10 rounded-2xl shadow-sm border border-gray-200">
        
        {/* Header */}
        <div className="flex items-center gap-3 mb-6">
          <h1 className="text-2xl font-bold text-gray-900">แก้ไขบันทึก</h1>
          <button
            onClick={goBack}
            className="ml-auto rounded-full border border-gray-300 text-gray-700 text-sm font-semibold px-5 py-2 hover:bg-gray-50 transition"
            disabled={saving}
          >
            ยกเลิก / กลับ
          </button>
        </div>

        <hr className="border-gray-200 mb-6" />

        {/* Alerts */}
        {formError && (
          <div className="mb-6 rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
            {formError}
          </div>
        )}

        {successMsg && (
          <div className="mb-6 rounded-xl border border-green-200 bg-green-50 px-4 py-3 text-sm text-green-700">
            {successMsg}
          </div>
        )}

        {/* Content */}
        {loading ? (
          <div className="p-10 text-center text-gray-500 animate-pulse">กำลังโหลดข้อมูล...</div>
        ) : error ? (
          <div className="p-4 text-sm text-orange-600 bg-orange-50 rounded text-center">{error}</div>
        ) : (
          <div className="space-y-8 max-w-4xl mx-auto">
            
            {/* Status Info */}
            <div className="bg-blue-50 border border-blue-100 rounded-lg p-3 text-sm text-blue-800 flex items-center gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" />
              </svg>
              <span>สถานะปัจจุบัน: <strong>{status}</strong></span>
            </div>

            {/* Title */}
            <div className="flex flex-col md:flex-row gap-4">
              <label className="w-40 text-left md:text-right pt-2 font-bold text-gray-800 text-lg">
                หัวข้อ
              </label>
              <input
                type="text"
                className="flex-1 border border-gray-300 rounded-lg px-4 py-2 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:border-[#F26522] focus:ring-1 focus:ring-[#F26522]"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                disabled={saving}
              />
            </div>

            {/* Body */}
            <div className="flex flex-col md:flex-row gap-4">
              <label className="w-40 text-left md:text-right pt-2 font-bold text-gray-800 text-lg leading-tight">
                พิมพ์ข้อสรุป<br className="hidden md:block" />/สิ่งที่ต้องทำ
              </label>
              <textarea
                rows={8}
                className="flex-1 border border-[#1E90FF] rounded-lg px-4 py-2 text-sm text-black placeholder-gray-400 focus:outline-none focus:border-[#1E90FF] focus:ring-1 focus:ring-[#1E90FF]"
                value={body}
                onChange={(e) => setBody(e.target.value)}
                disabled={saving}
                style={{ whiteSpace: "pre-wrap" }}
              />
            </div>

            {/* Existing Files */}
            <div className="flex flex-col md:flex-row gap-4 items-start">
              <label className="w-40 text-left md:text-right font-bold text-gray-800 text-lg pt-2">
                ไฟล์เดิม
              </label>
              <div className="flex-1 w-full bg-gray-50 rounded-lg border border-gray-200 p-4">
                {fileNames.length === 0 ? (
                  <div className="text-sm text-gray-500 italic">ไม่มีไฟล์แนบเดิม</div>
                ) : (
                  <ul className="space-y-2">
                    {fileNames.map((n, idx) => (
                      <li key={idx} className="flex items-center gap-2 text-sm text-gray-700">
                        <span className="text-gray-400">📄</span> {n}
                      </li>
                    ))}
                  </ul>
                )}
                
                {/* ⚠️ คำเตือนเรื่องการแทนที่ไฟล์ */}
                <div className="mt-3 text-xs text-red-500 bg-red-50 p-2 rounded border border-red-100 flex items-start gap-1">
                   <span>⚠️</span>
                   <span>หากคุณอัปโหลดไฟล์ใหม่ในช่องด้านล่าง <strong>ระบบจะลบไฟล์เดิมทั้งหมด</strong> และใช้ไฟล์ใหม่แทน</span>
                </div>
              </div>
            </div>

            {/* New Files Upload */}
            <div className="flex flex-col md:flex-row gap-4 items-start">
              <label className="w-40 text-left md:text-right font-bold text-gray-800 text-lg pt-2">
                อัปโหลดไฟล์ใหม่
              </label>
              <div className="flex-1 w-full">
                <div className="flex items-center gap-3">
                  <label className="cursor-pointer border border-gray-300 text-gray-800 px-4 py-2 rounded-lg hover:bg-gray-50 text-sm transition shadow-sm bg-white">
                    + เลือกไฟล์ (แทนที่เดิม)
                    <input
                      type="file"
                      className="hidden"
                      multiple
                      accept=".pdf,.jpg,.jpeg,.png"
                      onChange={handleFileChange}
                      disabled={saving}
                    />
                  </label>
                  <span className="text-sm text-gray-600">
                    {files.length > 0 ? `${files.length} ไฟล์ใหม่` : "ยังไม่ได้เลือกไฟล์ใหม่"}
                  </span>
                </div>

                {files.length > 0 && (
                  <div className="mt-3 space-y-2">
                    {files.map((f, idx) => (
                      <div key={idx} className="flex items-center justify-between gap-3 text-sm text-gray-900 bg-blue-50 border border-blue-100 rounded-lg px-3 py-2">
                        <span className="truncate">{f.name} <span className="text-gray-400 text-xs">({(f.size/1024).toFixed(0)}KB)</span></span>
                        <button
                          type="button"
                          className="text-red-600 hover:text-red-800 text-xs font-bold"
                          onClick={() => removeNewFile(idx)}
                          disabled={saving}
                        >
                          ลบ
                        </button>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>

            {/* Requires Report */}
            <div className="flex flex-col md:flex-row gap-4 items-start mt-4">
              <label className="w-40 text-left md:text-right pt-2 font-bold text-gray-800 text-lg leading-tight">
                ส่งรายงาน<br className="hidden md:block"/>ความคืบหน้า?
              </label>
              <div className="flex gap-2 pt-1">
                <button
                  type="button"
                  onClick={() => setRequiresReport(true)}
                  disabled={saving}
                  className={`px-4 py-2 rounded-md text-sm font-bold border transition ${
                    requiresReport === true
                      ? "bg-[#F26522] text-white border-[#F26522]"
                      : "bg-white text-gray-600 border-gray-200 hover:border-[#F26522] hover:text-[#F26522]"
                  }`}
                >
                  ต้องส่ง
                </button>
                <button
                  type="button"
                  onClick={() => setRequiresReport(false)}
                  disabled={saving}
                  className={`px-4 py-2 rounded-md text-sm font-bold border transition ${
                    requiresReport === false
                      ? "bg-[#F26522] text-white border-[#F26522]"
                      : "bg-white text-gray-600 border-gray-200 hover:border-[#F26522] hover:text-[#F26522]"
                  }`}
                >
                  ไม่ต้องส่ง
                </button>
              </div>
            </div>

            {/* Action Buttons */}
            <div className="flex flex-wrap gap-4 mt-10 md:pl-44 border-t border-gray-100 pt-6">
              <button
                disabled={saving}
                onClick={saveContent}
                className="bg-[#1E90FF] disabled:opacity-60 hover:bg-[#1877d6] text-white text-lg px-6 py-2.5 rounded-xl font-bold shadow-md transition"
              >
                {saving ? "กำลังทำงาน..." : "บันทึกการแก้ไข"}
              </button>

              <button
                disabled={saving}
                onClick={publish}
                className="bg-[#F26522] disabled:opacity-60 hover:bg-[#d1490e] text-white text-lg px-6 py-2.5 rounded-xl font-bold shadow-md transition"
              >
                {saving ? "กำลังทำงาน..." : "เผยแพร่ทันที"}
              </button>
            </div>

          </div>
        )}
      </div>
    </div>
  );
}