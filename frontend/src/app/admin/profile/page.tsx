"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/src/contexts/AuthContext";
import Image from "next/image";
import { message, Select } from "antd";
import { getPrefixes } from "@/src/services/http/masterservice";
import {
  Mail,
  Phone,
  Building2,
  Calendar,
  CheckCircle,
  XCircle,
  PenLine,
  Users,
  LogIn,
  HardHat,
} from "lucide-react";

import type {
  AdminProfileResponse,
  AdminProfileUpdateData,
} from "@/src/interfaces/adminprofile";
import {
  getAdminProfile,
  updateAdminProfile,
} from "@/src/services/http/adminservice";

export default function AdminProfilePage() {
  const { isAuthenticated, isLoading: authLoading } = useAuth();
  const router = useRouter();

  const [profile, setProfile] = useState<AdminProfileResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [editMode, setEditMode] = useState(false);
  const [saving, setSaving] = useState(false);
  const [prefixOptions, setPrefixOptions] = useState<
    { label: string; value: string }[]
  >([]);

  const [formData, setFormData] = useState<AdminProfileUpdateData>({
    prefix: "",
    firstName: "",
    lastName: "",
    phone: "",
    notes: "",
    active: true,
    email: "",
  });

  const fetchAdminProfile = async () => {
    try {
      setIsLoading(true);
      const data = await getAdminProfile();
      setProfile(data);
      setFormData({
        prefix: data.prefix,
        firstName: data.firstName,
        lastName: data.lastName,
        phone: data.phone,
        notes: data.notes,
        active: data.active,
        email: data.email,
      });
    } catch (err) {
      message.error("ไม่สามารถโหลดข้อมูลโปรไฟล์ได้");
    } finally {
      setIsLoading(false);
    }
  };
  useEffect(() => {
    getPrefixes()
      .then((data) =>
        setPrefixOptions(
          data.map((p) => ({
            label: p.prefix,
            value: p.prefix,
          }))
        )
      )
      .catch(() => message.error("โหลดคำนำหน้าไม่สำเร็จ"));
  }, []);

  useEffect(() => {
    if (!authLoading && !isAuthenticated) router.push("/login");
    if (isAuthenticated) fetchAdminProfile();
  }, [authLoading, isAuthenticated, router]);

  const handleSave = async () => {
    if (!formData.prefix.trim()) return message.error("กรุณาเลือกคำนำหน้า");
    if (!formData.firstName.trim()) return message.error("กรุณากรอกชื่อจริง");
    if (!formData.lastName.trim()) return message.error("กรุณากรอกนามสกุล");
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!formData.email || !emailRegex.test(formData.email)) {
      return message.error("รูปแบบอีเมลไม่ถูกต้อง");
    }
    if (!formData.phone || !/^[0-9]{10}$/.test(formData.phone)) {
      return message.error("เบอร์โทรศัพท์ต้องเป็นตัวเลข 10 หลัก");
    }

    try {
      setSaving(true);
      await updateAdminProfile({
        prefix: formData.prefix,
        firstName: formData.firstName,
        lastName: formData.lastName,
        phone: formData.phone,
        notes: formData.notes,
        active: formData.active,
        email: formData.email,
      });
      message.success("บันทึกข้อมูลสำเร็จ");
      setEditMode(false);
      await fetchAdminProfile(); // อัปเดตข้อมูลทันทีโดยไม่รีโหลด
    } catch (err) {
      message.error("บันทึกข้อมูลล้มเหลว");
    } finally {
      setSaving(false);
    }
  };

  if (authLoading || isLoading)
    return <div className="p-10 text-center text-gray-500">กำลังโหลด...</div>;
  if (!profile) return <div className="p-10 text-center">ไม่พบข้อมูล</div>;

  return (
    <div className="pb-6">
      {/* TOP BAR */}
      <div className="flex justify-between items-center mb-6">
        <div className="flex gap-2 ml-auto">
          {!editMode ? (
            <button
              onClick={() => setEditMode(true)}
              className="rounded-full bg-[#F26522] px-6 py-2 text-xs font-medium text-white hover:bg-[#e05a1d] transition-all shadow-sm"
            >
              แก้ไข
            </button>
          ) : (
            <>
              <button
                onClick={() => setEditMode(false)}
                className="rounded-full border border-gray-300 px-6 py-2 text-xs font-medium text-gray-700 bg-white dark:bg-gray-800 dark:text-gray-200"
              >
                ยกเลิก
              </button>
              <button
                onClick={handleSave}
                disabled={saving}
                className="rounded-full bg-[#F26522] px-6 py-2 text-xs font-medium text-white hover:bg-[#e05a1d] transition-all shadow-sm"
              >
                {saving ? "กำลังบันทึก..." : "บันทึก"}
              </button>
            </>
          )}
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        {/* LEFT COLUMN - PROFILE CARD */}
        <div className="lg:col-span-4 space-y-6">
          <div className="bg-white dark:bg-gray-800 rounded-3xl p-8 shadow-sm text-center">
            <div className="relative w-32 h-32 mx-auto mb-4">
              <Image
                src="/profile-placeholder.jpg"
                alt="Profile"
                fill
                className="rounded-full object-cover border-4 border-[#fff2eb] dark:border-gray-700"
              />
              <div
                className={`absolute bottom-1 right-2 w-6 h-6 border-4 border-white dark:border-gray-800 rounded-full ${
                  profile.active ? "bg-green-500" : "bg-red-500"
                }`}
              />
            </div>
            <h2 className="text-base font-bold text-slate-800 dark:text-white">
              {profile.prefix} {profile.firstName} {profile.lastName}
            </h2>
            <p className="text-[#F26522] font-medium text-sm mb-4">
              ผู้ดูแลระบบ (Admin)
            </p>
            <div className="text-sm text-slate-500 dark:text-gray-400 space-y-1">
              <p className="flex items-center justify-center gap-1">
                <Building2 className="w-4 h-4" /> {profile.departmentName}
              </p>
            </div>
          </div>

          {/* QUICK STATS */}
          <div className="bg-white dark:bg-gray-800 rounded-3xl p-6 shadow-sm space-y-4">
            <StatRow
              label="ผู้ใช้ที่ดูแล"
              value={`${profile.totalManaged} คน`}
              Icon={Users}
            />
            <StatRow
              label="เข้าสู่ระบบล่าสุด"
              value={profile.lastLogin || "-"}
              Icon={LogIn}
            />
            <StatRow
              label="วันที่สร้างบัญชี"
              value={profile.createdAt}
              Icon={Calendar}
            />
          </div>
        </div>

        {/* RIGHT COLUMN - DETAILS */}
        <div className="lg:col-span-8 space-y-6">
          <div className="bg-white dark:bg-gray-800 rounded-3xl p-8 shadow-sm">
            <h3 className="text-lg font-bold text-[#F26522] mb-6">
              ข้อมูลส่วนตัวและการติดต่อ
            </h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
              {/* Prefix + First + Last inline row */}
              <div className="md:col-span-2 grid grid-cols-1 md:grid-cols-3 gap-4 border-b border-dashed pb-6">
                {editMode ? (
                  <div>
                    <p className="text-sm uppercase tracking-wide text-gray-400">
                      คำนำหน้า
                    </p>
                    <Select
                      value={formData.prefix || undefined}
                      placeholder="เลือกคำนำหน้า"
                      options={prefixOptions}
                      onChange={(value) =>
                        setFormData({ ...formData, prefix: value })
                      }
                      className="mt-1 w-full"
                    />
                  </div>
                ) : (
                  <DataField label="คำนำหน้า" value={profile.prefix} />
                )}

                {editMode ? (
                  <InputField
                    label="ชื่อจริง"
                    value={formData.firstName}
                    onChange={(v) =>
                      setFormData({ ...formData, firstName: v })
                    }
                  />
                ) : (
                  <DataField label="ชื่อจริง" value={profile.firstName} />
                )}

                {editMode ? (
                  <InputField
                    label="นามสกุล"
                    value={formData.lastName}
                    onChange={(v) => setFormData({ ...formData, lastName: v })}
                  />
                ) : (
                  <DataField label="นามสกุล" value={profile.lastName} />
                )}
              </div>
              {editMode ? (
                <InputField
                  label="อีเมล"
                  value={formData.email}
                  onChange={(v) => setFormData({ ...formData, email: v })}
                />
              ) : (
                <DataField label="อีเมล" value={profile.email} Icon={Mail} />
              )}

              {editMode ? (
                <InputField
                  label="เบอร์โทรศัพท์"
                  value={formData.phone}
                  onChange={(v) => setFormData({ ...formData, phone: v })}
                />
              ) : (
                <DataField
                  label="เบอร์โทรศัพท์"
                  value={profile.phone || "-"}
                  Icon={Phone}
                />
              )}

              <div className="md:col-span-2">
                {editMode ? (
                  <div className="space-y-2">
                    <p className="text-sm font-medium text-gray-400 uppercase tracking-wide">
                      หมายเหตุ / หน้าที่
                    </p>
                    <textarea
                      value={formData.notes}
                      onChange={(e) =>
                        setFormData({ ...formData, notes: e.target.value })
                      }
                      className="w-full bg-slate-50 dark:bg-gray-700 border border-slate-200 dark:border-gray-600 rounded-2xl p-4 text-sm focus:ring-2 focus:ring-[#F26522] outline-none dark:text-white"
                      rows={3}
                    />
                  </div>
                ) : (
                  <DataField
                    label="หมายเหตุ / หน้าที่"
                    value={profile.notes || "-"}
                    Icon={PenLine}
                  />
                )}
              </div>

              <div className="flex flex-col gap-2">
                <p className="text-sm font-medium uppercase tracking-wide text-gray-400">
                  สถานะบัญชี
                </p>
                {editMode ? (
                  <div className="flex items-center gap-3 py-1">
                    <Toggle
                      checked={formData.active}
                      onChange={(v) => setFormData({ ...formData, active: v })}
                    />
                    <span className="text-sm font-medium dark:text-white">
                      {formData.active ? "เปิดการใช้งาน" : "ปิดการใช้งาน"}
                    </span>
                  </div>
                ) : (
                  <span
                    className={`text-sm inline-flex items-center rounded-full px-3 py-1 ${
                      profile.active
                        ? "bg-green-100 text-green-600"
                        : "bg-red-100 text-red-600"
                    }`}
                  >
                    ● {profile.active ? "Active" : "Inactive"}
                  </span>
                )}
              </div>
            </div>
          </div>

          {/* ADMIN LINKS */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <QuickLinkCard
              href="/admin/manageuser"
              title="จัดการรายชื่อผู้ใช้"
              subtitle="Manage all user accounts"
              Icon={Users}
            />
            <QuickLinkCard
              href="/admin/settings"
              title="ตั้งค่าระบบ"
              subtitle="System configuration"
              Icon={PenLine}
            />
          </div>
        </div>
      </div>
    </div>
  );
}

/* --- REUSABLE MINI COMPONENTS --- */

const DataField = ({
  label,
  value,
  Icon,
}: {
  label: string;
  value: string;
  Icon?: any;
}) => (
  <div className="flex flex-col">
    <p className="text-sm uppercase tracking-wide text-gray-400 flex items-center gap-1">
      {Icon && <Icon className="w-3 h-3" />}
      {label}
    </p>
    <p className="mt-1 text-sm font-medium text-gray-800 dark:text-gray-200">
      {value || "-"}
    </p>
  </div>
);

const InputField = ({
  label,
  value,
  onChange,
}: {
  label: string;
  value: string;
  onChange: (v: string) => void;
}) => (
  <div>
    <p className="text-sm uppercase tracking-wide text-gray-400">{label}</p>
    <input
      type="text"
      value={value}
      onChange={(e) => onChange(e.target.value)}
      className="mt-1 w-full rounded-xl border border-gray-300 dark:border-gray-600 px-3 py-2 text-sm text-gray-800 dark:text-white bg-white dark:bg-gray-700 focus:ring-2 focus:ring-[#F26522] outline-none transition-all"
    />
  </div>
);

const SelectField = ({
  label,
  value,
  onChange,
  options,
}: {
  label: string;
  value: string;
  onChange: (v: string) => void;
  options: string[];
}) => (
  <div>
    <p className="text-sm uppercase tracking-wide text-gray-400">{label}</p>
    <select
      value={value}
      onChange={(e) => onChange(e.target.value)}
      className="mt-1 w-full rounded-xl border border-gray-300 dark:border-gray-600 px-3 py-2 text-sm text-gray-800 dark:text-white bg-white dark:bg-gray-700 focus:ring-2 focus:ring-[#F26522] outline-none transition-all"
    >
      <option value="">เลือก</option>
      {options.map((opt) => (
        <option key={opt} value={opt}>
          {opt}
        </option>
      ))}
    </select>
  </div>
);

const StatRow = ({
  label,
  value,
  Icon,
}: {
  label: string;
  value: string;
  Icon: any;
}) => (
  <div className="flex items-center justify-between">
    <div className="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400">
      <Icon className="w-4 h-4" /> <span>{label}</span>
    </div>
    <span className="text-sm font-bold text-slate-800 dark:text-white">
      {value}
    </span>
  </div>
);

const QuickLinkCard = ({
  href,
  title,
  subtitle,
  Icon,
}: {
  href: string;
  title: string;
  subtitle: string;
  Icon: any;
}) => (
  <a
    href={href}
    className="group flex items-center justify-between bg-white dark:bg-gray-800 p-6 rounded-3xl shadow-sm hover:border-[#F26522] border border-transparent transition-all"
  >
    <div className="flex items-center gap-4">
      <div className="p-3 rounded-2xl bg-orange-50 dark:bg-orange-900/20 text-[#F26522]">
        <Icon className="w-5 h-5" />
      </div>
      <div>
        <p className="text-base font-bold text-slate-800 dark:text-white">
          {title}
        </p>
        <p className="text-xs text-slate-400 uppercase tracking-tighter">
          {subtitle}
        </p>
      </div>
    </div>
    <span className="text-[#F26522] opacity-0 group-hover:opacity-100 transition-all">
      →
    </span>
  </a>
);

const Toggle = ({
  checked,
  onChange,
}: {
  checked: boolean;
  onChange: (v: boolean) => void;
}) => (
  <button
    onClick={() => onChange(!checked)}
    className={`w-11 h-6 flex items-center rounded-full p-1 transition-all ${
      checked ? "bg-green-500" : "bg-gray-300"
    }`}
  >
    <div
      className={`bg-white w-4 h-4 rounded-full shadow-md transition-all ${
        checked ? "translate-x-5" : "translate-x-0"
      }`}
    />
  </button>
);
