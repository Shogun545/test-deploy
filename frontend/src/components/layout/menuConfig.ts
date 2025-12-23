export type Role = "admin" | "advisor" | "student";

export type MenuItem = {
  key: string;
  label: string;
  href: string;
  icon?: "dashboard" | "users" | "settings" | "appointments" | "profile" | "students" | "calendar" | "document";
};

export const MENU_BY_ROLE: Record<Role, MenuItem[]> = {
  admin: [
    { key: "admin-dashboard", label: "Dashboard", href: "/admin/dashboard", icon: "dashboard" },
    { key:  "admin-academiccalendar", label: "Academic Calendar", href: "/admin/calendar", icon: "calendar" },
    { key: "admin-users", label: "Manage Users", href: "/admin/manageuser", icon: "users" },
    { key: "admin-profile", label: "My Profile", href: "/admin/profile", icon: "profile" },
    { key: "admin-settings", label: "Settings", href: "/admin/settings", icon: "settings" },
    { key: "admin-report",label: "Report",href: "/admin/report", icon: "profile"},
    { key: "admin-status",label: "report-status",href: "/admin/report-status", icon: "profile"},

  ],
  advisor: [
    { key: "advisor-dashboard", label: "Dashboard", href: "/advisor/dashboard", icon: "dashboard" },
    { key: "advisor-appointments",label: "My Appointments",href: "/advisor/appointments", icon: "appointments"},
    { key: "advisor-advisorlog",label: "My Logs",href: "/advisor/advisorlog", icon: "document"},
    { key: "advisor-students",label: "My Students",href: "/advisor/mystudents", icon: "users"},
    { key: "advisor-profile",label: "My Profile",href: "/advisor/profile", icon: "profile"},
    { key: "advisor-report",label: "Report",href: "/advisor/report", icon: "profile"},
  ],
  student: [
    {key: "student-dashboard",label: "My Dashboard",href: "/student/dashboard", icon: "dashboard"},
    {key: "student-appointments",label: "My Appointments",href: "/student/appointments", icon: "appointments"},
    {key: "student-profile",label: "My Profile",href: "/student/profile", icon: "profile"},
    {key: "student-advisorlog",label: "Advisor Logs",href: "/student/advisorlog",icon: "document"},
  ],
};
