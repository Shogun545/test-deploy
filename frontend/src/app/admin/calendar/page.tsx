"use client";
import CalendarView from "@/src/components/academiccalendar/pageview"; // ตรวจสอบ Path import ให้ถูกนะครับ

export default function AdminCalendarPage() {
  // ✅ ส่ง canEdit={true}
  return <CalendarView canEdit={true} />;
}