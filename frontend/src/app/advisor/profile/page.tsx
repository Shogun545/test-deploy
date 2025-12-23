"use client";
import { Camera } from "lucide-react";
import { uploadProfileImage } from "@/src/services/http/cloudinary";
import { useEffect, useState, useCallback, useRef } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/src/contexts/AuthContext";
import Image from "next/image";
import { message, Select } from "antd";
import "dayjs/locale/th";
import { getPrefixes } from "@/src/services/http/masterservice";

import type { AdvisorProfileUI } from "@/src/interfaces/advisorprofile";
import {
  getAdvisorProfile,
  updateAdvisorProfile,
} from "@/src/services/http/advisorservice";
import { getAdvisorStudents } from "@/src/services/http/advisorstudents";

export default function AdvisorProfilePage() {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();

  const [profile, setProfile] = useState<AdvisorProfileUI | null>(null);
  const [loadingProfile, setLoadingProfile] = useState(true);
  const [editMode, setEditMode] = useState(false);
  const [saving, setSaving] = useState(false);
  const [realStudentCount, setRealStudentCount] = useState<number | null>(null);
  const [prefixOptions, setPrefixOptions] = useState<
    { label: string; value: string }[]
  >([]);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [profileImage, setProfileImage] = useState<string | null>(null);

  // ข้อมูลที่ใช้ในฟอร์มแก้ไข
  const [formData, setFormData] = useState({
    prefix: "",
    first_name: "",
    last_name: "",
    email: "",
    phone: "",
    officeRoom: "",
    specialties: "",
    isActive: false,
  });

  // ฟังก์ชันดึงข้อมูล (ใช้ซ้ำได้ทั้งตอนโหลดหน้าและหลังบันทึก)
  const fetchInitialData = useCallback(async () => {
    try {
      const [profileData, studentData] = await Promise.all([
        getAdvisorProfile(),
        getAdvisorStudents(),
      ]);

      const mapped: AdvisorProfileUI = {
        fullName: `${profileData.prefix}${profileData.first_name} ${profileData.last_name}`,
        email: profileData.email,
        phone: profileData.phone,
        officeRoom: profileData.office_room,
        majorName: profileData.major_name,
        departmentName: profileData.department_name,
        adviseeCount: profileData.advisee_count,
        specialties: profileData.specialties,
        isActive: profileData.is_active,
        profileImage: profileData.profile_image ?? null,
      };

      setProfile(mapped);
      setRealStudentCount(studentData.students?.length || 0);
      setProfileImage(profileData.profile_image ?? null);

      // อัปเดตข้อมูลในฟอร์ม
      setFormData({
        prefix: profileData.prefix || "",
        first_name: profileData.first_name || "",
        last_name: profileData.last_name || "",
        email: profileData.email || "",
        phone: profileData.phone || "",
        officeRoom: profileData.office_room || "",
        specialties: profileData.specialties || "",
        isActive: profileData.is_active,
      });
    } catch (err: any) {
      console.error("Fetch Error:", err.message);
      message.error("ไม่สามารถโหลดข้อมูลโปรไฟล์ได้");
    } finally {
      setLoadingProfile(false);
    }
  }, []);

  useEffect(() => {
    if (!isLoading && !isAuthenticated) router.push("/login");
    if (isAuthenticated) fetchInitialData();
  }, [isLoading, isAuthenticated, router, fetchInitialData]);

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

    // Validation
    if (!formData.prefix.trim()) return message.error("กรุณาเลือกคำนำหน้า");
    if (!formData.first_name.trim()) return message.error("กรุณากรอกชื่อจริง");
    if (!formData.last_name.trim()) return message.error("กรุณากรอกนามสกุล");
    if (!formData.email.trim()) return message.error("กรุณากรอกอีเมล");

    try {
      setSaving(true);
      await updateAdvisorProfile({
        prefix: formData.prefix,
        first_name: formData.first_name,
        last_name: formData.last_name,
        email: formData.email,
        phone: formData.phone,
        office_room: formData.officeRoom,
        specialties: formData.specialties,
        is_active: formData.isActive,
          profile_image: profileImage, 

      });

      message.success("บันทึกข้อมูลสำเร็จ");
      setEditMode(false);
      await fetchInitialData(); // อัปเดต UI ทันทีโดยไม่ต้อง reload
    } catch (err) {
      message.error("บันทึกข้อมูลล้มเหลว");
    } finally {
      setSaving(false);
    }
  };

  const handleCancel = () => {
    // รีเซ็ตค่ากลับไปเป็นข้อมูลล่าสุดที่โหลดมา
    fetchInitialData();
    setEditMode(false);
  };

  if (isLoading || loadingProfile)
    return <div className="p-10 text-center text-gray-500">กำลังโหลด...</div>;
  if (!profile) return <div className="p-10 text-center">ไม่พบข้อมูล</div>;

  return (
    <div className="pb-6">
      <div className="flex justify-between items-center mb-6">
        <div className="flex gap-2 ml-auto">
          {!editMode ? (
            <button
              onClick={() => setEditMode(true)}
              className="rounded-full bg-[#F26522] px-6 py-2 text-xs font-medium text-white hover:bg-[#e05a1d] transition-all"
            >
              แก้ไขโปรไฟล์
            </button>
          ) : (
            <>
              <button
                onClick={handleCancel}
                className="rounded-full border border-gray-300 px-6 py-2 text-xs font-medium text-gray-700 bg-white"
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
        {/* LEFT COLUMN */}
        <div className="lg:col-span-4 space-y-6">
          <div className="bg-white rounded-3xl p-8 shadow-sm text-center">
            <div className="relative w-32 h-32 mx-auto mb-4">
              <Image
                src={profileImage || "/profile-placeholder.jpg"}
                alt="Profile"
                fill
                className="rounded-full object-cover border-4 border-[#fff2eb]"
              />
              {editMode && (
                <>
                  <input
                    ref={fileInputRef}
                    type="file"
                    accept="image/*"
                    hidden
                    onChange={async (e) => {
                      if (!e.target.files) return;

                      const url = await uploadProfileImage(e.target.files[0]);
                      setProfileImage(url);
                    }}
                  />

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

              <div
                className={`absolute bottom-1 right-2 w-6 h-6 border-4 border-white rounded-full ${
                  profile.isActive ? "bg-green-500" : "bg-slate-300"
                }`}
              ></div>
            </div>
            <h2 className="text-xl font-bold text-slate-800">
              {profile.fullName}
            </h2>
            <p className="text-[#F26522] font-medium text-sm mb-4">
              อาจารย์ที่ปรึกษา
            </p>
            <div className="text-sm text-slate-500 space-y-1">
              <p>{profile.email}</p>
              <p>{profile.phone || "ยังไม่ได้ระบุเบอร์โทรศัพท์"}</p>
            </div>
          </div>

          <div className="bg-white rounded-3xl p-6 shadow-sm">
            <h3 className="font-bold text-slate-800 mb-4 border-b pb-2 text-sm">
              สถิติการดูแล
            </h3>
            <div className="flex justify-between items-center">
              <span className="text-slate-500 text-sm">นักศึกษาในดูแล</span>
              <span className="text-2xl font-bold text-[#F26522]">
                {realStudentCount ?? profile.adviseeCount}{" "}
                <small className="text-xs text-slate-400 font-normal">คน</small>
              </span>
            </div>
          </div>
        </div>

        {/* RIGHT COLUMN */}
        <div className="lg:col-span-8 space-y-6">
          <div className="bg-white rounded-3xl p-8 shadow-sm">
            <h3 className="text-lg font-bold text-[#F26522] mb-6">
              รายละเอียดข้อมูลส่วนตัว
            </h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
              {/* ส่วนแก้ไข ชื่อ-นามสกุล และ คำนำหน้า */}
              <div className="md:col-span-2 grid grid-cols-1 md:grid-cols-3 gap-4 border-b border-dashed pb-6">
                {editMode ? (
                  <div>
                    <p className="text-[11px] uppercase tracking-wide text-gray-400">
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
                  <DataField label="คำนำหน้า" value={formData.prefix} />
                )}

                {editMode ? (
                  <InputField
                    label="ชื่อจริง"
                    value={formData.first_name}
                    onChange={(v) =>
                      setFormData({ ...formData, first_name: v })
                    }
                  />
                ) : (
                  <DataField label="ชื่อจริง" value={formData.first_name} />
                )}

                {editMode ? (
                  <InputField
                    label="นามสกุล"
                    value={formData.last_name}
                    onChange={(v) => setFormData({ ...formData, last_name: v })}
                  />
                ) : (
                  <DataField label="นามสกุล" value={formData.last_name} />
                )}
              </div>

              <DataField label="สาขาวิชา" value={profile.majorName} />

              {editMode ? (
                <InputField
                  label="อีเมล"
                  value={formData.email}
                  onChange={(v) => setFormData({ ...formData, email: v })}
                />
              ) : (
                <DataField label="อีเมล" value={profile.email} />
              )}

              {editMode ? (
                <InputField
                  label="ห้องทำงาน"
                  value={formData.officeRoom}
                  onChange={(v) => setFormData({ ...formData, officeRoom: v })}
                />
              ) : (
                <DataField
                  label="ห้องทำงาน"
                  value={profile.officeRoom || "-"}
                />
              )}

              {editMode ? (
                <InputField
                  label="เบอร์โทรศัพท์"
                  value={formData.phone}
                  onChange={(v) => setFormData({ ...formData, phone: v })}
                />
              ) : (
                <DataField label="เบอร์โทรศัพท์" value={profile.phone || "-"} />
              )}

              <div className="flex flex-col gap-2">
                <p className="text-sm font-medium uppercase tracking-wide text-gray-400">
                  สถานะการรับคำปรึกษา
                </p>
                {editMode ? (
                  <div className="flex items-center gap-3 py-1">
                    <Toggle
                      checked={formData.isActive}
                      onChange={(v) =>
                        setFormData({ ...formData, isActive: v })
                      }
                    />
                    <span className="text-sm font-medium">
                      {formData.isActive
                        ? "พร้อมให้คำปรึกษา"
                        : "ปิดรับชั่วคราว"}
                    </span>
                  </div>
                ) : (
                  <span
                    className={`text-sm inline-flex items-center gap-2 ${
                      profile.isActive ? "text-green-500" : "text-slate-400"
                    }`}
                  >
                    <span
                      className={`w-2 h-2 rounded-full ${
                        profile.isActive ? "bg-green-500" : "bg-slate-400"
                      }`}
                    ></span>
                    {profile.isActive
                      ? "Scheduled (พร้อมรับนัด)"
                      : "Unavailable (ไม่พร้อมรับนัด)"}
                  </span>
                )}
              </div>
            </div>
          </div>

          <div className="bg-white rounded-3xl p-8 shadow-sm">
            <h3 className="text-lg font-bold text-[#F26522] mb-4">
              ความเชี่ยวชาญเฉพาะทาง
            </h3>
            {editMode ? (
              <textarea
                value={formData.specialties}
                onChange={(e) =>
                  setFormData({ ...formData, specialties: e.target.value })
                }
                className="w-full bg-slate-50 border border-slate-100 rounded-2xl p-4 text-sm focus:ring-2 focus:ring-[#F26522] outline-none"
                rows={4}
              />
            ) : (
              <p className="text-slate-600 leading-relaxed bg-[#fffbf9] p-4 rounded-2xl border border-[#fff2eb] text-sm italic">
                {profile.specialties || "ยังไม่ได้ระบุความเชี่ยวชาญพิเศษ"}
              </p>
            )}
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <QuickLinkCard
              href="/advisor/mystudents"
              title="รายชื่อนักศึกษา"
              subtitle="Advisee management"
            />
            <QuickLinkCard
              href="/advisor/appointments"
              title="ตารางนัดหมาย"
              subtitle="Booking schedule"
            />
          </div>
        </div>
      </div>
    </div>
  );
}

/* --- REUSABLE COMPONENTS --- */

const QuickLinkCard = ({
  href,
  title,
  subtitle,
}: {
  href: string;
  title: string;
  subtitle: string;
}) => (
  <a
    href={href}
    className="group flex items-center justify-between bg-white p-6 rounded-3xl shadow-sm hover:border-[#F26522] border border-transparent transition-all"
  >
    <div>
      <p className="text-base font-bold text-slate-800">{title}</p>
      <p className="text-xs text-slate-400 uppercase tracking-tighter">
        {subtitle}
      </p>
    </div>
    <span className="w-10 h-10 rounded-full bg-slate-50 group-hover:bg-[#fff2eb] flex items-center justify-center transition text-[#F26522]">
      →
    </span>
  </a>
);

const DataField = ({ label, value }: { label: string; value: string }) => (
  <div className="flex flex-col">
    <p className="text-sm uppercase tracking-wide text-gray-400">{label}</p>
    <p className="mt-1 text-sm font-medium text-gray-800">{value || "-"}</p>
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
    <p className="text-[11px] uppercase tracking-wide text-gray-400">{label}</p>
    <input
      type="text"
      value={value}
      onChange={(e) => onChange(e.target.value)}
      className="mt-1 w-full rounded-xl border border-gray-300 px-3 py-2 text-sm text-gray-800 bg-white focus:border-[#F26522] outline-none transition-all"
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
    <p className="text-[11px] uppercase tracking-wide text-gray-400">{label}</p>
    <select
      value={value}
      onChange={(e) => onChange(e.target.value)}
      className="mt-1 w-full rounded-xl border border-gray-300 px-3 py-2 text-sm text-gray-800 bg-white focus:border-[#F26522] outline-none transition-all"
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
      checked ? "bg-[#F26522]" : "bg-slate-300"
    }`}
  >
    <div
      className={`bg-white w-4 h-4 rounded-full shadow-md transition-all ${
        checked ? "translate-x-5" : "translate-x-0"
      }`}
    />
  </button>
);
