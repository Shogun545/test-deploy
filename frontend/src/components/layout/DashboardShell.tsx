"use client";

import { useState } from "react";
import { usePathname } from "next/navigation";
import { useAuth } from "@/src/contexts/AuthContext";
import Sidebar from "./Sidebar";
import Topbar from "./Topbar";
import type { Role } from "./menuConfig";
import { MENU_BY_ROLE } from "./menuConfig";

function normalizeRole(roleFromUser?: string | null): Role {
  const r = (roleFromUser || "Student").toLowerCase();
  if (r === "admin") return "admin";
  if (r === "advisor") return "advisor";
  return "student";
}

export default function DashboardShell({
  children,
}: {
  children: React.ReactNode;
}) {
  const { user } = useAuth();
  const pathname = usePathname();
  const [collapsed, setCollapsed] = useState(false);

  const role = normalizeRole(user?.role);
  const menus = MENU_BY_ROLE[role];
  // 🔥 แก้ไข Logic การหาชื่อหัวข้อ (Title) ตรงนี้
  const activeMenu = menus.find((m) => pathname.startsWith(m.href));

  let pageTitle = activeMenu?.label;

  // ✅ เพิ่มเงื่อนไขพิเศษ: ถ้า path มีคำว่า /calendar ให้ตั้งชื่อเป็น Academic Calendar เลย
  // (โดยไม่ต้องสนใจว่าจะมีในเมนู sidebar หรือไม่)
  if (!pageTitle && pathname.includes("/calendar")) {
    pageTitle = "Academic Calendar";
  }
  // Fallback: ถ้าหาไม่เจอจริงๆ ให้ใช้ชื่อเมนูแรก (Dashboard)
  if (!pageTitle) {
    pageTitle = menus[0].label;
  }

  const sidebarWidth = collapsed ? 90 : 260;

  return (
    // เพิ่ม flex เพื่อให้ div ย่อยสามารถใช้ flex-1 ได้
    <div className="h-screen bg-[#f6f7f9] flex overflow-hidden"> 
      {/* SIDEBAR */}
      <Sidebar
        role={role}
        collapsed={collapsed}
        onToggle={() => setCollapsed((prev) => !prev)}
      />

      {/* MAIN AREA */}
      <div
        className="flex flex-col transition-all duration-300 flex-1" // 2. เพิ่ม flex-col และ flex-1
        style={{ marginLeft: sidebarWidth }}
      >
        {/* TOPBAR */}
        <div className="px-8 pt-6 flex-shrink-0"> {/* เพิ่ม flex-shrink-0 */}
          <Topbar />
        </div>

        {/* HEADER + CONTENT */}
        <div className="px-8 pb-8 pt-4 flex flex-col flex-1 min-h-0"> {/* 3. เพิ่ม flex-col, flex-1, min-h-0 */}
          
          {/* หัวข้อหน้า */}
          <div className="mb-3 flex-shrink-0"> 
            <div className="text-[#F26522] font-normal text-lg">
              {/* ✅ ใช้ตัวแปร pageTitle ที่เราคำนวณมาแสดง */}
              {pageTitle}
            </div>
            <div className="mt-1 h-[1px] bg-gray-200" />
          </div>

          <div className="mt-4 rounded-[32px] bg-white shadow-[0_0_25px_rgba(0,0,0,0.03)] p-8 flex-1 overflow-y-auto">
            <div className="flex flex-col h-full">
                {children}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}