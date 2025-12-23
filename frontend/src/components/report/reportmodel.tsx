import { useState, useEffect } from "react";
import { updateReportStatus } from "@/src/services/http/report";

type Props = {
  open: boolean;
  onClose: () => void;
  report: any;
  onUpdated: () => void;
};


export default function ReportModal({
  open,
  onClose,
  report,
  onUpdated,
}: Props) {

  const [statusId, setStatusId] = useState<number>(0);
  const [loading, setLoading] = useState(false);

  // ✅ สำคัญมาก
  useEffect(() => {
    if (report?.status?.id) {
      setStatusId(report.status.id);
    }
  }, [report]);

  if (!open || !report) return null;

  const handleSave = async () => {
    console.log("SEND report.ID =", report.ID);
    console.log("SEND statusId =", statusId);

    try {
      setLoading(true);
      await updateReportStatus(report.ID, statusId);
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

        <textarea
          className="w-full border rounded-lg p-3 text-sm mb-4"
          rows={4}
          readOnly
          value={report.description}
        />

        <div className="flex items-center justify-between">
          <button
            disabled
            className="px-3 py-1 border rounded-lg text-gray-400"
          >
            Download
          </button>

          <div className="flex items-center gap-2">
            <span className="text-sm">Status</span>

            <select
              className="border rounded-lg px-3 py-1 text-sm"
              value={statusId}
              onChange={(e) => setStatusId(Number(e.target.value))}
            >
              <option value={1}>Pending</option>
              <option value={2}>Resolved</option>
              <option value={3}>Inprogress</option>
            </select>

            <button
              onClick={handleSave}
              disabled={loading || statusId === 0}
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
