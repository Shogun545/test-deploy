"use client";

import { useState, useEffect, ChangeEvent } from "react";

import { ReporttopicInterface } from "@/src/interfaces/ireporttopic";

import { ReportInterface } from "@/src/interfaces/ireport";

export interface IssueReportFormData {
  title: string;
  topic: number | "";
  description: string;
  file?: File | null;
}

export default function IssueReportForm() {
  const [topics, setTopics] = useState<ReporttopicInterface[]>([]);

  const [form, setForm] = useState<IssueReportFormData>({
    title: "",
    topic: "",
    description: "",
    file: null,
  });

  const handleChange = (
    e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>
  ) => {
    const { name, value } = e.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  };

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (!['image/png', 'image/jpeg'].includes(file.type)) {
      alert('รองรับเฉพาะไฟล์ PNG หรือ JPG เท่านั้น');
      e.target.value = '';
      return;
    }

    setForm((prev) => ({ ...prev, file }));
  };

  useEffect(() => {
    // TODO: replace with real API call
    const mockTopics: ReporttopicInterface[] = [
      { ID: 1, reporttopic_name: "บัญชีผู้ใช้" },
      { ID: 2, reporttopic_name: "การเข้าสู่ระบบ" },
      { ID: 3, reporttopic_name: "ระบบ" },
    ];

    setTopics(mockTopics);
  }, []);

  const handleSubmit = () => {
    const payload: ReportInterface = {
      description: form.description,
      report_topic: form.topic || undefined,
      // report_by: userId,      // ดึงจาก auth context
      // report_status: 1,        // default status
      // report_image: imageId,   // ได้หลัง upload
    };

    console.log(payload);
    // TODO: submit payload + image to backend
  };

  return (
    <div className="w-full rounded-2xl bg-white p-6 shadow-sm border border-gray-200">
      <h2 className="text-lg font-semibold text-gray-800 mb-4">
        การรายงานปัญหา
      </h2>

      {/* Title */}
      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700 mb-1">
          หัวเรื่อง
        </label>
        <input
          name="title"
          value={form.title}
          onChange={handleChange}
          placeholder="กรุณากรอกข้อมูล"
          className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm text-gray-900 bg-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-orange-400"
        />
      </div>

      {/* Topic */}
      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700 mb-1">
          หัวข้อปัญหา
        </label>
        <select
          name="topic"
          value={form.topic}
          onChange={handleChange}
          className="w-48 rounded-lg border border-gray-300 px-3 py-2 text-sm text-gray-900 bg-white focus:outline-none focus:ring-2 focus:ring-orange-400"
        >
          <option value="">เลือกหัวข้อ</option>
          {topics.map((t) => (
            <option key={t.ID} value={t.ID}>
              {t.reporttopic_name}
            </option>
          ))}
        </select>
      </div>

      {/* Description */}
      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700 mb-1">
          ข้อมูลปัญหา
        </label>
        <textarea
          name="description"
          value={form.description}
          onChange={handleChange}
          rows={6}
          className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm resize-none focus:outline-none focus:ring-2 focus:ring-orange-400"
        />
      </div>

      {/* Attach file */}
      <div className="mb-6 flex items-center gap-4">
        <label className="flex items-center gap-2 cursor-pointer rounded-lg border border-gray-300 px-4 py-2 text-sm text-gray-700 hover:bg-gray-50">
          <input type="file" className="hidden" onChange={handleFileChange} />
          Attach File
        </label>
        {form.file && (
          <span className="text-sm text-gray-600 truncate max-w-xs">
            {form.file.name}
          </span>
        )}
      </div>

      {/* Submit */}
      <div className="flex justify-end">
        <button
          onClick={handleSubmit}
          className="rounded-full bg-orange-500 px-6 py-2 text-sm font-medium text-white hover:bg-orange-600"
        >
          Submit
        </button>
      </div>
    </div>
  );
}
