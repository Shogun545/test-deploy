"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/src/contexts/AuthContext";

export default function StudentDashboardPage() {
  const { isAuthenticated, isLoading, user } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push("/login");
    }
  }, [isLoading, isAuthenticated, router]);

  if (isLoading) return <div>กำลังโหลด...</div>;
  if (!isAuthenticated) return <div>กำลังเปลี่ยนหน้า...</div>;

  return (
    <div>
      <p className="mb-1">ยินดีต้อนรับ {user?.sut_id}</p>
      <p className="mb-6">Role: {user?.role}</p>
      {/* ตรงนี้เพื่อนใส่การ์ด / ตาราง / อะไรก็ได้ */}
    </div>
  );
}
