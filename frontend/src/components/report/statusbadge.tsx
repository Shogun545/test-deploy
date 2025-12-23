export default function StatusBadge({ status }: { status: string }) {
  const style: Record<string, string> = {
    Resolved: "bg-green-100 text-green-700",
    Pending: "bg-red-100 text-red-700",        // 🔴 unresolved
    Unresolved: "bg-red-100 text-red-700",     // เผื่อมีคำนี้
    Inprogress: "bg-blue-100 text-blue-700",   // 🔵
  };

  return (
    <span
      className={`px-3 py-1 rounded-full text-xs font-medium ${style[status]}`}
    >
      {status}
    </span>
  );
}
