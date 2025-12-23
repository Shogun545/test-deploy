import StatusBadge from "./statusbadge";

type Props = {
  reports: any[];
  onEdit: (id: number) => void;
};

export default function ReportTable({ reports, onEdit }: Props) {
  return (
    <div className="bg-white border rounded-xl shadow-sm">
      <div className="px-4 py-3 font-medium">All Issue</div>

      <table className="w-full text-sm border-collapse">
        <thead className="bg-gray-50 text-gray-500 border-b border-gray-100">

          <tr>
            <th className="px-4 py-3 text-left">ID</th>
            <th className="px-4 py-3 text-left">Topic</th>
            <th className="px-4 py-3 text-left">Status</th>
            <th className="px-4 py-3 text-center">Action</th>
          </tr>
        </thead>

        <tbody>
          {reports.map((r) => (
            <tr
              key={r.id}
              className="border-b border-gray-100 hover:bg-gray-50"
            >
              <td className="px-4 py-3">{r.id}</td>
              <td className="px-4 py-3 font-medium">
                {r.topic?.name}
              </td>
              <td className="px-4 py-3">
                <StatusBadge status={r.status?.name} />
              </td>
              <td className="px-4 py-3 text-center">
                <button
                  onClick={() => onEdit(r.id)}
                  className="px-3 py-1 text-blue-600 border border-blue-500 rounded-lg hover:bg-blue-50"
                >
                  แก้ไข
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
