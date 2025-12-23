import type { ReactNode } from "react";
import DashboardShell from "@/src/components/layout/DashboardShell";

export default function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <DashboardShell>{children}</DashboardShell>;
}
