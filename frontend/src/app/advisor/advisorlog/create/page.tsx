"use client";

import { ChangeEvent, useMemo, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { CreateAdvisorLog } from "@/src/services/http/advisorlog";
import type { AdvisorLogInterface } from "@/src/interfaces/advisorlog";

export default function AdvisorLogCreatePage() {
  const router = useRouter();
  const search = useSearchParams();
  const appointmentIdFromQuery = search?.get("appointmentId");

  const [successMsg, setSuccessMsg] = useState<string | null>(null);
  const [formError, setFormError] = useState<string | null>(null); // ✅ เพิ่ม

  const [title, setTitle] = useState("");
  const [body, setBody] = useState("");
  const [requiresReport, setRequiresReport] = useState<boolean | null>(null);

  const [files, setFiles] = useState<File[]>([]);
  const [loading, setLoading] = useState(false);
  const [savingMsg, setSavingMsg] = useState<string | null>(null);
  const [uploadPercent, setUploadPercent] = useState<number>(0);

  const appointmentId = useMemo(() => {
    const n = Number(appointmentIdFromQuery);
    return Number.isFinite(n) && n > 0 ? n : 0;
  }, [appointmentIdFromQuery]);

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    const list = e.target.files;
    if (!list || list.length === 0) return;

    const selected = Array.from(list);
    setFiles((prev) => [...prev, ...selected]);
    e.target.value = "";
  };

  const handleRemoveFile = (index: number) => {
    setFiles((prev) => prev.filter((_, i) => i !== index));
  };

  const goBackToList = () => {
    router.push("/advisor/advisorlog");
  };

  const handleSubmit = async (isDraft = false) => {
    if (loading) return;

    // ✅ ไม่ใช้ alert -> ใช้ข้อความ error ในหน้าแทน
    if (!isDraft) {
      if (!title.trim()) {
        setFormError("กรุณากรอกหัวข้อ");
        return;
      }
      if (!body.trim()) {
        setFormError("กรุณากรอกเนื้อหาสรุป");
        return;
      }
      if (requiresReport === null) {
        setFormError("กรุณาเลือกว่า ต้องส่งรายงานหรือไม่");
        return;
      }
    }

    setFormError(null); 
    setSuccessMsg(null);

    setLoading(true);
    setUploadPercent(0);
    setSavingMsg(isDraft ? "กำลังบันทึกฉบับร่าง..." : "กำลังบันทึก...");

    const payload: AdvisorLogInterface = {
      AppointmentID: appointmentId || 0,
      Title: title.trim(),
      Body: body.trim(),
      RequiresReport: requiresReport === true,
    };

    const res = await CreateAdvisorLog(payload, {
      files,
      isDraft,
      onUploadProgress: (p) => setUploadPercent(p),
    });

    if (res.success) {
      setSavingMsg(null);
      setFormError(null);
      setSuccessMsg(isDraft ? "บันทึกฉบับร่างแล้ว" : "เผยแพร่สำเร็จ");

      setTimeout(() => {
        router.push("/advisor/advisorlog");
      }, 2000);

      return;
    }

    setSavingMsg(null);
    setFormError(res.message || "เกิดข้อผิดพลาดในการบันทึก");
    setLoading(false);
  };

  return (
    <div className="w-full min-h-screen bg-gray-50 flex justify-center py-6 px-4">
      <div className="w-full max-w-5xl bg-white p-10 rounded-2xl shadow-sm border border-gray-200">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">สร้างบันทึกการปรึกษา</h1>
        <p className="text-gray-500 text-sm mb-4 pb-4 border-b border-gray-200">
          AppointmentID: {appointmentId || "-"}
        </p>

        {/* ✅ Error message (ไม่ใช้ alert) */}
        {formError && (
          <div
            className="mb-6 rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700"
            data-testid="form-error"
          >
            {formError}
          </div>
        )}

        {/* ✅ Success message */}
        {successMsg && (
          <div
            className="mb-6 rounded-xl border border-green-200 bg-green-50 px-4 py-3 text-sm text-green-700"
            data-testid="success-message"
          >
            {successMsg}
          </div>
        )}

        {/* saving / uploading */}
        {savingMsg ? (
          <div className="mb-6 rounded-xl border border-orange-200 bg-orange-50 px-4 py-3 text-sm text-orange-700">
            {savingMsg}
            {files.length > 0 ? (
              <span className="ml-2 text-orange-600">({uploadPercent}% อัปโหลด)</span>
            ) : null}
          </div>
        ) : null}

        <div className="space-y-6 max-w-4xl mx-auto">
          {/* Title */}
          <div className="flex flex-col md:flex-row gap-4">
            <label className="w-40 text-left md:text-right pt-2 font-bold text-gray-800 text-lg">
              หัวข้อ
            </label>
            <input
              type="text"
              className="flex-1 border border-gray-300 rounded-lg px-4 py-2 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:border-[#F26522] focus:ring-1 focus:ring-[#F26522]"
              placeholder="ระบุหัวข้อ..."
              value={title}
              onChange={(e) => {
                setTitle(e.target.value);
                if (formError) setFormError(null); // ✅ เคลียร์ error ตอนพิมพ์
              }}
              disabled={loading}
              data-testid="log-title"
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
              placeholder="พิมพ์รายละเอียด..."
              value={body}
              onChange={(e) => {
                setBody(e.target.value);
                if (formError) setFormError(null); // ✅ เคลียร์ error ตอนพิมพ์
              }}
              disabled={loading}
              style={{ whiteSpace: "pre-wrap" }}
              data-testid="log-body"
            />
          </div>

          {/* Files */}
          <div className="flex flex-col md:flex-row gap-4 items-start">
            <label className="w-40 text-left md:text-right font-bold text-gray-800 text-lg pt-2">
              ที่แนบไฟล์
            </label>

            <div className="flex-1 w-full">
              <div className="flex items-center gap-3">
                <label className="cursor-pointer border border-gray-300 text-gray-800 px-4 py-2 rounded-lg hover:bg-gray-50 text-sm transition shadow-sm bg-white">
                  + อัปโหลดไฟล์
                  <input
                    type="file"
                    className="hidden"
                    multiple
                    accept=".pdf,.jpg,.jpeg,.png"
                    onChange={handleFileChange}
                    disabled={loading}
                  />
                </label>

                {files.length > 0 ? (
                  <span className="text-sm text-gray-700">{files.length} ไฟล์</span>
                ) : (
                  <span className="text-sm text-gray-500">ยังไม่ได้เลือกไฟล์</span>
                )}
              </div>

              <div className="mt-3 space-y-2">
                {files.map((file, idx) => (
                  <div
                    key={`${file.name}-${idx}`}
                    className="flex items-center justify-between gap-3 text-sm text-gray-900 bg-gray-50 border border-gray-200 rounded-lg px-3 py-2"
                  >
                    <div className="min-w-0">
                      <div className="truncate font-medium">{file.name}</div>
                      <div className="text-xs text-gray-500">{(file.size / 1024).toFixed(0)} KB</div>
                    </div>

                    <button
                      type="button"
                      className="shrink-0 text-red-600 hover:text-red-800 text-sm font-semibold"
                      onClick={() => handleRemoveFile(idx)}
                      disabled={loading}
                    >
                      ลบ
                    </button>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* RequiresReport */}
          <div className="flex flex-col md:flex-row gap-4 items-center mt-4">
            <label className="w-auto md:pl-44 font-bold text-gray-800 text-md">
              ส่งรายงานความคืบหน้าภายใน 1 สัปดาห์ หรือไม่
            </label>
            <div className="flex gap-2">
              <button
                type="button"
                onClick={() => {
                  setRequiresReport(true);
                  if (formError) setFormError(null); // ✅ เคลียร์ error ตอนเลือก
                }}
                disabled={loading}
                className={`px-4 py-2 rounded-md text-sm font-bold border transition ${
                  requiresReport === true
                    ? "bg-[#F26522] text-white border-[#F26522]"
                    : "bg-white text-gray-600 border-gray-200 hover:border-[#F26522] hover:text-[#F26522]"
                }`}
                data-testid="requires-yes"
              >
                ต้องส่ง
              </button>

              <button
                type="button"
                onClick={() => {
                  setRequiresReport(false);
                  if (formError) setFormError(null); // ✅ เคลียร์ error ตอนเลือก
                }}
                disabled={loading}
                className={`px-4 py-2 rounded-md text-sm font-bold border transition ${
                  requiresReport === false
                    ? "bg-[#F26522] text-white border-[#F26522]"
                    : "bg-white text-gray-600 border-gray-200 hover:border-[#F26522] hover:text-[#F26522]"
                }`}
                data-testid="requires-no"
              >
                ไม่ต้องส่ง
              </button>
            </div>
          </div>

          {/* Buttons */}
          <div className="flex flex-wrap gap-4 mt-10 md:pl-40">
            <button
              disabled={loading}
              onClick={() => handleSubmit(false)}
              className="bg-[#F26522] disabled:opacity-60 hover:bg-[#d1490e] text-white text-lg px-8 py-2.5 rounded-xl font-bold shadow-md transition"
              data-testid="publish"
            >
              {loading ? "กำลังบันทึก..." : "เผยแพร่"}
            </button>

            <button
              disabled={loading}
              onClick={() => handleSubmit(true)}
              className="bg-white disabled:opacity-60 border-2 border-[#F26522] text-[#F26522] text-lg px-6 py-2.5 rounded-xl font-bold hover:bg-orange-50 transition"
              data-testid="save-draft"
            >
              บันทึกฉบับร่าง
            </button>

            <button
              disabled={loading}
              onClick={goBackToList}
              className="text-gray-500 disabled:opacity-60 text-lg px-4 py-2.5 font-medium hover:text-gray-700 ml-auto"
            >
              ยกเลิก
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
