import { useState } from "react";
import { updateReportStatus } from "@/src/services/http/report";

type Props = {
  open: boolean;
  onClose: () => void;
  report: any;
  onUpdated: () => void; // ให้ parent refresh table
};

export default function ReportModal({
  open,
  onClose,
  report,
  onUpdated,
}: Props) {

console.log("MODAL REPORT =", report);


  const [statusId, setStatusId] = useState<number>(
    report?.status?.id
  );
  const [loading, setLoading] = useState(false);

  if (!open || !report) return null;

const handleSave = async () => {
  try {
    setLoading(true);
    await updateReportStatus(report.ID, statusId); // ✅ แก้ตรงนี้
    onUpdated();
    onClose();
  } catch (err) {
    alert("ไม่สามารถอัปเดตสถานะได้");
  } finally {
    setLoading(false);
  }
};


  return (
    <div className="fixed inset-0 bg-black/30 flex items-center justify-center z-50">
      <div className="bg-white w-[560px] rounded-xl p-6">
        <h2 className="text-lg font-semibold mb-4">
          Report Details
        </h2>

        {/* Info */}
        <div className="text-sm space-y-1 mb-4">
          <p>
            <b>Issue ID:</b> {report.ID} &nbsp;
            <b>Topic:</b> {report.topic?.reporttopic_name}
          </p>
          <p>
            <b>Report By:</b>{" "}
            {report.user?.sut_id} {report.user?.first_name}{" "}
            {report.user?.last_name}
            
          </p>
          <p>
            <b>Email:</b> {report.user?.email}
          </p>
        </div>

        {/* Description */}
        <textarea
          className="w-full border rounded-lg p-3 text-sm mb-4"
          rows={4}
          readOnly
          value={report.description}
        />

        {/* Action row */}
        <div className="flex items-center justify-between">
          {/* Attachment (ยังไม่ทำ) */}
          <button
            disabled
            className="px-3 py-1 border rounded-lg text-gray-400"
          >
            Download
          </button>

          {/* Status */}
          <div className="flex items-center gap-2">
            <span className="text-sm">Status</span>
            <select
              className="border rounded-lg px-3 py-1 text-sm"
              value={statusId}
              onChange={(e) => setStatusId(Number(e.target.value))}
            >
              <option value={1}>Pending</option>
              <option value={2}>Inprogress</option>
              <option value={3}>Resolved</option>
            </select>

            <button
              onClick={handleSave}
              disabled={loading}
              className="px-4 py-1 bg-orange-500 text-white rounded-lg"
            >
              {loading ? "Saving..." : "Save"}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
  
}
