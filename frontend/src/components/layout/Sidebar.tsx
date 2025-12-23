"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Menu, ChevronLeft, LogOut } from "lucide-react";
import { LayoutGrid, NotebookPen, Users, User as UserIcon, Settings, CalendarDays,FileText, } from "lucide-react";
import type { MenuItem, Role } from "./menuConfig";
import { MENU_BY_ROLE } from "./menuConfig";
import Image from "next/image";
import { useAuth } from "@/src/contexts/AuthContext";

type SidebarProps = {
  role: Role | string | undefined;
  collapsed: boolean;
  onToggle: () => void;
};

const IconFor = (name?: MenuItem["icon"]) => {
  switch (name) {
    case "dashboard": return <LayoutGrid className="w-5 h-5" />;
    case "appointments": return <NotebookPen className="w-5 h-5" />;
    case "students": return <Users className="w-5 h-5" />;
    case "profile": return <UserIcon className="w-5 h-5" />;
    case "settings": return <Settings className="w-5 h-5" />;
    case "users": return <Users className="w-5 h-5" />;
    case "calendar": return <CalendarDays className="w-5 h-5" />;
     case "document": return <FileText className="w-5 h-5" />;
    default: return null;
  }
};

export default function Sidebar({ role, collapsed, onToggle }: SidebarProps) {
  const pathname = usePathname();
  // ป้องกันกรณี role ผิดพลาด ให้เป็น array ว่าง
  const menus = (role && MENU_BY_ROLE && MENU_BY_ROLE[role as keyof typeof MENU_BY_ROLE]) || [];
  const { logout } = useAuth();

  const width = collapsed ? 90 : 260; // px

  // wrapper สำหรับเรียก logout แล้วจับ error
  const handleLogout = () => {
    try {
      logout();
    } catch (e) {
      console.error("Logout error:", e);
      // fallback: ล้าง localStorage เงียบ ๆ แล้ว reload
      try {
        localStorage.removeItem("token");
        localStorage.removeItem("user");
        window.location.href = "/login";
      } catch {}
    }
  };

  return (
    <div
      className="fixed left-0 top-0 h-screen bg-white shadow-lg transition-all duration-300"
      style={{ width }}
    >
      <button
        onClick={onToggle}
        className="absolute -right-4 top-6 bg-white w-8 h-8 rounded-full shadow flex items-center justify-center"
      >
        {collapsed ? <Menu size={18} /> : <ChevronLeft size={18} />}
      </button>

      <div className="px-6 py-8 flex items-center justify-start">
        {!collapsed ? (
          <div className="flex flex-col">
            <Image
              src="/LOGO.svg"
              alt="ENGI Advisory Logo"
              width={200}
              height={50}
            />
          </div>
        ) : (
            // กรณีถูกหุบ → แสดงโลโก้แบบเล็กกลางแทน
            <Image
            src="/LOGO.svg"
            alt="ENGI Advisory Logo"
            width={40}
            height={40}
            className="mx-auto"
          />
        )}
      </div>

      <nav className="mt-4 px-3 space-y-2">
        {menus.length === 0 ? (
          <div className="px-4 py-2 text-sm text-gray-400">เมนูยังไม่พร้อม</div>
        ) : (
          menus.map((item: MenuItem) => {
            const active = pathname?.startsWith(item.href || "") ?? false;
            return (
              <Link
                key={item.key}
                href={item.href}
                className={`
                  flex items-center gap-3 px-4 py-2 rounded-2xl text-sm transition
                  ${active ? "bg-[#f26522] text-white" : "text-gray-600 hover:bg-gray-100"}
                `}
              >
                <span className="shrink-0">{IconFor(item.icon)}</span>
                {!collapsed && item.label}
              </Link>
            );
          })
        )}
      </nav>

      <div className="absolute bottom-6 w-full px-4 border-t border-gray-200 pt-2">
        <button
          onClick={handleLogout}
          className="flex items-center gap-3 w-full px-4 py-2 rounded-2xl text-sm text-gray-600 hover:bg-gray-100"
        >
          <LogOut className="w-5 h-5" />
          {!collapsed && "ออกจากระบบ"}
        </button>
      </div>
    </div>
  );
}
