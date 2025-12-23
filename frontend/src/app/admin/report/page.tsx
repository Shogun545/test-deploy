"use client";

import { useEffect, useState } from "react";
import {
  getReportSummary,
  getReports,
  getReportById,
} from "@/src/services/http/report";
import ReportTable from "@/src/components/report/reporttable";
import ReportModal from "@/src/components/report/reportmodal-backup";

export default function ReportPage() {
  // summary
  const [summary, setSummary] = useState<any>(null);

  // report list
  const [reports, setReports] = useState<any[]>([]);

  // modal state
  const [open, setOpen] = useState(false);
  const [selected, setSelected] = useState<any>(null);
  const [loading, setLoading] = useState(false);

  // โหลดข้อมูลครั้งแรก
  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    setLoading(true);
    try {
      const [summaryData, reportList] = await Promise.all([
        getReportSummary(),
        getReports(),
      ]);
      setSummary(summaryData);
      setReports(reportList);
    } finally {
      setLoading(false);
    }
  };

  // กดปุ่ม "แก้ไข"
  const handleEdit = async (id: number) => {
    try {
      const data = await getReportById(id);
      setSelected(data);
      setOpen(true);
    } catch (err) {
      alert("ไม่สามารถโหลดข้อมูล Report ได้");
    }
  };

  // หลังอัปเดต status เสร็จ
  const handleUpdated = async () => {
    await loadData();
  };

  return (
    <div className="p-6">
      {/* ===== Title ===== */}
      <h1 className="text-lg font-semibold mb-6">
        Report Dashboard
      </h1>

      {/* ===== Summary ===== */}
      {summary && (
        <div className="grid grid-cols-4 gap-4 mb-6">
          {[
            { label: "Total Issue", value: summary.total },
            { label: "Pending", value: summary.pending },
            { label: "Inprogress", value: summary.inprogress },
            { label: "Resolved", value: summary.resolved },
          ].map((item) => (
            <div
              key={item.label}
              className="bg-white border rounded-xl p-4 text-center shadow-sm"
            >
              <p className="text-sm text-gray-500">
                {item.label}
              </p>
              <p className="text-2xl font-semibold">
                {item.value}
              </p>
            </div>
          ))}
        </div>
      )}

      {/* ===== Table ===== */}
      <div className="bg-white rounded-xl shadow-sm">
        {loading ? (
          <div className="p-6 text-center text-gray-500">
            Loading...
          </div>
        ) : (
    <ReportTable reports={reports} onEdit={handleEdit} />
  )}
</div>


      {/* ===== Modal ===== */}
      <ReportModal
        open={open}
        report={selected}
        onClose={() => setOpen(false)}
        onUpdated={handleUpdated}
      />
    </div>
  );
}
