"use client";

import { useEffect, useState, useRef } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/src/contexts/AuthContext";
import Image from "next/image";
import { DatePicker, message, Select } from "antd";
import { Camera } from "lucide-react";
import dayjs from "dayjs";
import "dayjs/locale/th";
import { getPrefixes } from "@/src/services/http/masterservice";
import { uploadProfileImage } from "@/src/services/http/cloudinary";
import type { ProfileUI } from "@/src/interfaces/studentprofile";
import {
  getStudentProfile,
  updateStudentProfile,
} from "@/src/services/http/studentservice";

dayjs.locale("th");

export default function StudentDashboardPage() {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [profile, setProfile] = useState<ProfileUI | null>(null);
  const [loadingProfile, setLoadingProfile] = useState(true);
  const [prefixOptions, setPrefixOptions] = useState<
    { label: string; value: string }[]
  >([]);

  const [editMode, setEditMode] = useState(false);
  const [saving, setSaving] = useState(false);
  const [formData, setFormData] = useState({
    prefix: "",
    first_name: "",
    last_name: "",
    phone: "",
    birthday: "",
    profile_image: "",
    email: "",
    year_of_study: 0,
    term_gpa: 0,
    cumulative_gpa: 0,
    academic_year: "", // เพิ่ม
    semester: 1, // เพิ่ม
  });

  // Move fetchProfile to component scope so it can be reused
  const fetchProfile = async () => {
    setLoadingProfile(true);
    try {
      const data = await getStudentProfile();
      const gpxVal = data.cumulative_gpa ?? 0;
      const mapped: ProfileUI = {
        title: data.prefix, // ใช้จาก backend
        firstName: data.first_name,
        lastName: data.last_name,
        roleLabel: "นักศึกษา",
        studentCode: data.sut_id,
        advisorName: data.advisor_name ?? "ยังไม่มีข้อมูล",
        email: data.email,
        phone: data.phone,
        dob: data.birthday,
        idCard: data.national_id,
        major: data.major_name,
        year: String(data.year_of_study),
        gpx: data.cumulative_gpa ?? 0,
        gpaLatest: data.term_gpa ?? 0,
        gpaTermLabel: data.gpa_term_label ?? "-",
        status:
          gpxVal === 0 ? "ไม่พบข้อมูลเกรด" : data.academic_status ?? "ปกติ",
        statusNote:
          gpxVal === 0
            ? "กรุณากรอกเกรดเฉลี่ยสะสมเพื่อดูสถานะการศึกษา"
            : "สถานะถูกคำนวณจากเกรดเฉลี่ยสะสมล่าสุด",
      };

      setProfile(mapped);
      setFormData({
        prefix: mapped.title ?? "",
        first_name: mapped.firstName ?? "",
        last_name: mapped.lastName ?? "",
        phone: mapped.phone ?? "",
        birthday: mapped.dob ?? "",
        profile_image: data.profile_image ?? "",
        email: mapped.email ?? "",
        year_of_study: data.year_of_study,
        academic_year: data.academic_year ?? "",
        semester: data.semester ?? 1,
        term_gpa: data.term_gpa ?? 0,
        cumulative_gpa: gpxVal ?? 0,
      });
    } catch (err: any) {
      console.error("โหลดข้อมูลโปรไฟล์ล้มเหลว:", err.message);
    } finally {
      setLoadingProfile(false);
    }
  };

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push("/login");
    }
  }, [isLoading, isAuthenticated]);

  useEffect(() => {
    if (!isAuthenticated) return;
    fetchProfile();
  }, [isAuthenticated]);

  useEffect(() => {
    getPrefixes()
      .then((data) => {
        setPrefixOptions(
          data.map((p) => ({
            label: p.prefix,
            value: p.prefix,
          }))
        );
      })
      .catch(() => {
        message.error("โหลดคำนำหน้าไม่สำเร็จ");
      });
  }, []);

  const handleSave = async () => {
    if (!profile) return;
    if (!formData.prefix.trim()) return message.error("กรุณาเลือกคำนำหน้า");
    if (!formData.first_name.trim()) return message.error("กรุณากรอกชื่อจริง");
    if (!formData.last_name.trim()) return message.error("กรุณากรอกนามสกุล");
    if (!formData.birthday) return message.error("กรุณากรอกวันเกิด");
    if (!formData.phone) return message.error("กรุณากรอกเบอร์โทรศัพท์");
    if (!/^[0-9]{10}$/.test(formData.phone))
      return message.error("เบอร์โทรศัพท์ต้องเป็นตัวเลข 10 หลัก");

    const email = formData.email.trim();
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!email) return message.error("กรุณากรอกอีเมล");
    if (!emailRegex.test(email)) return message.error("อีเมลไม่ถูกต้อง");

    const year = Number(formData.year_of_study);
    if (!Number.isInteger(year) || year <= 0)
      return message.error("กรุณากรอกชั้นปีเป็นตัวเลขที่ถูกต้อง");

    const termGPA = Number(formData.term_gpa);
    const cumulativeGPA = Number(formData.cumulative_gpa);
    const inRange = (v: number) => v >= 0 && v <= 4.0;
    if (!inRange(termGPA) || !inRange(cumulativeGPA))
      return message.error("เกรดต้องอยู่ระหว่าง 0.00 ถึง 4.00");

    try {
      setSaving(true);
      await updateStudentProfile({
        prefix: formData.prefix,
        first_name: formData.first_name,
        last_name: formData.last_name,
        phone: formData.phone,
        birthday: formData.birthday,
        profile_image: formData.profile_image,
        email: formData.email,
        year_of_study: year,
        term_gpa: termGPA,
        cumulative_gpa: cumulativeGPA,
        academic_year: formData.academic_year,
        semester: formData.semester,
      });
      message.success("บันทึกข้อมูลสำเร็จ");
      setEditMode(false);
      await fetchProfile();
    } catch (err: any) {
      message.error(err?.response?.data?.error || "อัปเดตข้อมูลล้มเหลว");
    } finally {
      setSaving(false);
    }
  };

  if (isLoading || loadingProfile)
    return (
      <div className="p-10 text-center text-gray-500">กำลังโหลดข้อมูล...</div>
    );
  if (!profile)
    return <div className="p-10 text-center">ไม่พบข้อมูลโปรไฟล์</div>;

  const gpxPercent = (profile.gpx / 4) * 100;

  return (
    <div className="pb-6">
      {/* TOP BAR / ACTION BUTTONS */}
      <div className="flex justify-between items-center mb-6">
        <div className="flex gap-2 ml-auto">
          {!editMode ? (
            <button
              onClick={() => setEditMode(true)}
              className="rounded-full bg-[#F26522] px-6 py-2 text-xs font-medium text-white hover:bg-[#e05a1d] transition-all"
            >
              แก้ไข
            </button>
          ) : (
            <>
              <button
                onClick={() => setEditMode(false)}
                className="rounded-full border border-gray-300 px-4 py-2 text-xs font-medium text-gray-700 bg-white dark:bg-gray-800 dark:text-gray-200"
              >
                ยกเลิก
              </button>
              <button
                onClick={handleSave}
                disabled={saving}
                className="rounded-full bg-[#F26522] px-6 py-2 text-xs font-medium text-white hover:bg-[#e05a1d] transition-all"
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
                src={formData.profile_image || "/profile-placeholder.jpg"}
                alt="Profile"
                fill
                className="rounded-full object-cover border-4"
              />

              {editMode && (
                <>
                  {/* hidden file input */}
                  <input
                    ref={fileInputRef}
                    type="file"
                    accept="image/*"
                    hidden
                    onChange={async (e) => {
                      if (!e.target.files) return;

                      const url = await uploadProfileImage(e.target.files[0]);
                      setFormData((prev) => ({
                        ...prev,
                        profile_image: url,
                      }));
                    }}
                  />

                  {/* camera button */}
                  <button
                    type="button"
                    onClick={() => fileInputRef.current?.click()}
                    className="absolute bottom-1 right-1 flex h-9 w-9 items-center justify-center rounded-full bg-[#F26522] text-white shadow hover:bg-[#e05a1d] transition"
                    title="เปลี่ยนรูปโปรไฟล์"
                  >
                    <Camera size={18} />
                  </button>
                </>
              )}
            </div>
            <h2 className="text-sm font-bold text-slate-800 dark:text-white">
              {profile.title} {profile.firstName} {profile.lastName}
            </h2>
            <p className="text-[#F26522] font-medium text-sm mb-4">
              {profile.roleLabel}
            </p>
            <div className="text-sm text-slate-500 dark:text-gray-400 space-y-1">
              <p>{profile.studentCode}</p>
              <p className="text-xs">อาจารย์ที่ปรึกษา: {profile.advisorName}</p>
            </div>

            {/* Academic Status merged into profile card */}
            <div className="mt-6 border-t dark:border-gray-700 pt-4 text-left">
              <h3 className="font-bold text-slate-800 dark:text-white mb-3 text-sm">
                สถานะการศึกษา
              </h3>
              <div className="flex justify-between items-center">
                <span className="text-slate-500 text-sm">สถานะปัจจุบัน</span>
                <span className="inline-flex items-center rounded-full bg-orange-100 dark:bg-orange-900/30 px-3 py-1 text-xs font-bold text-[#F26522]">
                  ● {profile.status}
                </span>
              </div>
              <p className="text-[11px] text-slate-400 text-right mt-2">
                {profile.statusNote}
              </p>
            </div>
          </div>
        </div>

        {/* RIGHT COLUMN - DETAILS & GRADES */}
        <div className="lg:col-span-8 space-y-6">
          {/* CARD: PERSONAL INFO */}
          <div className="bg-white dark:bg-gray-800 rounded-3xl p-8 shadow-sm">
            <h3 className="text-lg font-bold text-[#F26522] mb-6">
              ข้อมูลส่วนตัว
            </h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
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
                <DataField label="คำนำหน้า" value={profile.title} />
              )}

              {editMode ? (
                <InputField
                  label="ชื่อจริง"
                  value={formData.first_name}
                  onChange={(val) =>
                    setFormData({ ...formData, first_name: val })
                  }
                  inputClassName="text-base"
                  labelClassName="text-sm"
                />
              ) : (
                <DataField label="ชื่อจริง" value={profile.firstName} />
              )}

              {editMode ? (
                <InputField
                  label="นามสกุล"
                  value={formData.last_name}
                  onChange={(val) =>
                    setFormData({ ...formData, last_name: val })
                  }
                  inputClassName="text-base"
                  labelClassName="text-sm"
                />
              ) : (
                <DataField label="นามสกุล" value={profile.lastName} />
              )}

              {editMode ? (
                <FieldDatePicker
                  label="วันเดือนปีเกิด"
                  currentDateString={formData.birthday}
                  onDateChange={(val) =>
                    setFormData({ ...formData, birthday: val })
                  }
                />
              ) : (
                <DataField label="วันเดือนปีเกิด" value={profile.dob} />
              )}

              {editMode ? (
                <InputField
                  label="อีเมล"
                  value={formData.email}
                  onChange={(val) => setFormData({ ...formData, email: val })}
                  inputClassName="text-base"
                  labelClassName="text-sm"
                />
              ) : (
                <DataField
                  label="อีเมล"
                  value={profile.email}
                  valueClassName="text-base"
                />
              )}

              {editMode ? (
                <InputField
                  label="เบอร์โทรศัพท์"
                  value={formData.phone}
                  onChange={(val) => setFormData({ ...formData, phone: val })}
                />
              ) : (
                <DataField label="เบอร์โทรศัพท์" value={profile.phone} />
              )}
            </div>
          </div>

          {/* CARD: EDUCATION INFO */}
          <div className="bg-white dark:bg-gray-800 rounded-3xl p-8 shadow-sm">
            <h3 className="text-lg font-bold text-[#F26522] mb-6">
              ข้อมูลการศึกษา
            </h3>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mb-8">
              <DataField label="รหัสนักศึกษา" value={profile.studentCode} />
              <DataField label="สาขาวิชา" value={profile.major} />
              {editMode ? (
                <InputField
                  label="ชั้นปี"
                  value={String(formData.year_of_study)}
                  onChange={(val) =>
                    setFormData({
                      ...formData,
                      year_of_study: Number(val) || 0,
                    })
                  }
                  inputClassName="text-base"
                  labelClassName="text-sm"
                />
              ) : (
                <DataField
                  label="ชั้นปี"
                  value={profile.year}
                  valueClassName="text-base"
                />
              )}
            </div>
            {/* --- เพิ่มส่วนนี้เข้าไปตรงนี้ครับ --- */}
            {editMode && (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-8 pb-6 border-b border-dashed border-gray-200">
                <InputField
                  label="ปีการศึกษา (เช่น 2568)"
                  value={formData.academic_year}
                  onChange={(val) =>
                    setFormData({ ...formData, academic_year: val })
                  }
                />
                <SelectField
                  label="ภาคการศึกษา"
                  value={String(formData.semester)}
                  onChange={(val) =>
                    setFormData({ ...formData, semester: Number(val) })
                  }
                  options={["1", "2", "3"]}
                />
              </div>
            )}
            {/* ---------------------------------- */}

            {/* GRADES GRID */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="bg-[#fafafa] dark:bg-gray-700/50 p-5 rounded-2xl border border-gray-100 dark:border-gray-700">
                <p className="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1">
                  GPX (สะสม)
                </p>
                {editMode ? (
                  <input
                    type="number"
                    step="0.01"
                    min={0}
                    max={4}
                    value={formData.cumulative_gpa}
                    onChange={(e) =>
                      setFormData({
                        ...formData,
                        cumulative_gpa: Number(e.target.value),
                      })
                    }
                    className="text-4xl font-bold text-[#F26522] w-24 bg-transparent outline-none"
                  />
                ) : (
                  <p className="text-4xl font-bold text-[#F26522]">
                    {profile.gpx.toFixed(2)}
                  </p>
                )}
                <div className="mt-3 h-2 w-full rounded-full bg-gray-200 dark:bg-gray-600 overflow-hidden">
                  <div
                    className="h-full bg-[#F26522] transition-all duration-500"
                    style={{ width: `${gpxPercent}%` }}
                  />
                </div>
              </div>

              <div className="bg-[#fafafa] dark:bg-gray-700/50 p-5 rounded-2xl border border-gray-100 dark:border-gray-700">
                <p className="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1">
                  GPA เทอมล่าสุด
                </p>
                {editMode ? (
                  <input
                    type="number"
                    step="0.01"
                    min={0}
                    max={4}
                    value={formData.term_gpa}
                    onChange={(e) =>
                      setFormData({
                        ...formData,
                        term_gpa: Number(e.target.value),
                      })
                    }
                    className="text-4xl font-bold text-[#F26522] w-24 bg-transparent outline-none"
                  />
                ) : (
                  <p className="text-4xl font-bold text-[#F26522]">
                    {profile.gpaLatest.toFixed(2)}
                  </p>
                )}
                <p className="mt-1 text-xs text-slate-400 font-medium">
                  {profile.gpaTermLabel}
                </p>
              </div>
            </div>
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
  labelClassName,
  valueClassName,
}: {
  label: string;
  value: string;
  labelClassName?: string;
  valueClassName?: string;
}) => (
  <div className="flex flex-col">
    <p
      className={`text-sm uppercase tracking-wide text-gray-400 ${
        labelClassName || ""
      }`}
    >
      {label}
    </p>
    <p
      className={`mt-1 text-sm font-medium text-gray-800 dark:text-gray-200 ${
        valueClassName || ""
      }`}
    >
      {value || "-"}
    </p>
  </div>
);
const InputField = ({
  label,
  value,
  onChange,
  labelClassName,
  inputClassName,
}: {
  label: string;
  value: string;
  onChange: (v: string) => void;
  labelClassName?: string;
  inputClassName?: string;
}) => (
  <div>
    <p
      className={`text-sm uppercase tracking-wide text-gray-400 ${
        labelClassName || ""
      }`}
    >
      {label}
    </p>
    <input
      type="text"
      value={value}
      onChange={(e) => onChange(e.target.value)}
      className={`mt-1 w-full rounded-md border border-gray-300 px-2 py-1 text-[13px] text-gray-800 bg-white dark:text-gray-200 dark:bg-gray-700 dark:border-gray-600 ${
        inputClassName || ""
      }`}
    />
  </div>
);
const SelectField = ({
  label,
  value,
  onChange,
  options,
  labelClassName,
  selectClassName,
}: {
  label: string;
  value: string;
  onChange: (v: string) => void;
  options: string[];
  labelClassName?: string;
  selectClassName?: string;
}) => (
  <div>
    <p
      className={`text-sm uppercase tracking-wide text-gray-400 ${
        labelClassName || ""
      }`}
    >
      {label}
    </p>
    <select
      value={value}
      onChange={(e) => onChange(e.target.value)}
      className={`mt-1 w-full rounded-md border border-gray-300 px-2 py-1 text-[13px] text-gray-800 bg-white dark:text-gray-200 dark:bg-gray-700 dark:border-gray-600 ${
        selectClassName || ""
      }`}
    >
      <option value="">-- กรุณาเลือก --</option>
      {options.map((opt) => (
        <option key={opt} value={opt}>
          {opt}
        </option>
      ))}
    </select>
  </div>
);
const FieldDatePicker = ({
  label,
  currentDateString,
  onDateChange,
}: {
  label: string;
  currentDateString: string;
  onDateChange: (v: string) => void;
}) => {
  const dayjsValue = currentDateString ? dayjs(currentDateString) : null;
  return (
    <div className="flex flex-col">
      <p className="text-sm uppercase tracking-wide text-gray-400 mb-1">
        {label}
      </p>
      <DatePicker
        allowClear={false}
        value={dayjsValue}
        onChange={(_, dateString) =>
          onDateChange(Array.isArray(dateString) ? dateString[0] : dateString)
        }
        format="YYYY-MM-DD"
        className="w-full rounded-xl py-2 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
      />
    </div>
  );
};
