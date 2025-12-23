import type { ReactNode } from "react";
import DashboardShell from "@/src/components/layout/DashboardShell";

export default function AdvisorLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <DashboardShell>{children}</DashboardShell>;
}
