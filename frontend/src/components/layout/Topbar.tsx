"use client";

import { Bell, CalendarDays, Search } from "lucide-react";
import { useAuth } from "@/src/contexts/AuthContext";
import { useState, useEffect } from "react";
import { ensureAbsoluteUrl } from "@/src/utils/url";


import { getStudentProfile } from "@/src/services/http/studentservice";
import { getAdvisorProfile } from "@/src/services/http/advisorservice";
import { getAdminProfile } from "@/src/services/http/adminservice";
import Link from "next/link";

export default function Topbar() {
  const { user } = useAuth();
  const displayId = user ? user.sut_id ?? user.sutId ?? "ไม่พบรหัส" : "Guest";
  const displayRole = user?.role ?? "Unknown";
  const [profileImage, setProfileImage] = useState<string | null>(null);

  useEffect(() => {
    if (!user) return;

    const currentUser = user; 

    async function loadProfileImage() {
      try {
        const role = (currentUser.role || "").toLowerCase();

        if (currentUser.profile_image) {
          setProfileImage(ensureAbsoluteUrl(currentUser.profile_image));
          return;
        }

        if (role === "student") {
          const profile = await getStudentProfile();
          setProfileImage(ensureAbsoluteUrl(profile.profile_image));
          return;
        }

        if (role === "advisor") {
          const profile = await getAdvisorProfile();
          setProfileImage(ensureAbsoluteUrl(profile.profile_image));
          return;
        }

        if (role === "admin") {
          setProfileImage(null);
          return;
        }
      } catch (err) {
        console.error("Failed to load profile image", err);
        setProfileImage(null);
      }
    }

    loadProfileImage();
  }, [user]);

  // ✅ Logic คำนวณ Path ปฏิทินตาม Role
  const role = user?.role?.toLowerCase() || "";
  let calendarPath = "/student/calendar"; // default เป็น student
  if (role === "admin") {
      calendarPath = "/admin/calendar";
  } else if (role === "advisor") {
      calendarPath = "/advisor/calendar";
  }

  return (
    <div className="w-full h-[70px] rounded-[25px] bg-[#e96a26] px-8 flex items-center text-white shadow-md">
      {/* LEFT: Search + Title */}
      <div className="flex items-center gap-6">
        <div className="flex items-center gap-3 w-[320px]">
          <div className="flex items-center bg-white/95 rounded-full px-4 py-2 flex-1">
            <input
              type="text"
              placeholder="ค้นหาอาจารย์"
              className="bg-transparent outline-none text-sm text-gray-700 flex-1"
            />
            <Search className="w-4 h-4 text-gray-400" />
          </div>
        </div>
        <div className="text-sm whitespace-nowrap">
          Welcome to ENGI Advisory
        </div>
      </div>

      {/* RIGHT */}
      <div className="flex items-center gap-5 ml-auto">
        <Link href={calendarPath}>
          <CalendarDays className="w-5 h-5" />
        </Link>
        <Bell className="w-5 h-5" />
        <div className="flex items-center gap-3 bg-white/10 rounded-full px-3 py-1">
          <div className="text-xs leading-tight">
            <div className="font-semibold">{displayId}</div>
            <div className="text-[10px] opacity-80">{displayRole}</div>
          </div>
          <div className="w-9 h-9 rounded-full overflow-hidden bg-gray-200">
            {profileImage ? (
              <img
                src={profileImage}
                alt="Profile"
                className="w-full h-full object-cover"
                referrerPolicy="no-referrer"
              />
            ) : (
              <div className="w-full h-full bg-gray-300" />
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
