"use client";

import { useEffect, useMemo, useState, useCallback } from "react";
import { useAuth } from "@/src/contexts/AuthContext";
import { useRouter } from "next/navigation";
import { Select, ConfigProvider } from "antd";
import {
  getManagedUsers,
  getUserDetailBySutId,
  getMajors,
  updateUserStatus,
  getUserCreatedDate,
  updateManagedUser,
} from "@/src/services/http/adminusers";
import type {
  ManagedUserEntry as BaseManagedUserEntry,
  UserFilters,
  AdminUserDetailResponse,
  MajorEntry,
} from "@/src/interfaces/adminusers";

// ใน Go DTO เราใช้ SutID เป็น string ดังนั้นเราจะใช้ string ใน Frontend ทั้งหมด
export interface ManagedUserEntry extends BaseManagedUserEntry {
  id: string; // บังคับให้ใช้ string (sutId) เป็น key
}
type UserData = ManagedUserEntry & {
  id: string; // ใช้ sutId เป็น key หลัก
  is_active: boolean; // Field เสริมสำหรับ Toggle
  createdAt: string;
};

// ข้อมูลสำหรับ Dropdown (อาจดึงจาก API แยกในอนาคต) mock ก่อนเดี๋ยวกลับมาแก้
const mockRoles = ["Admin", "Advisor", "Student"];

// --- Components ย่อย  ---
function Avatar({ name }: { name: string }) {
  const parts = name.trim().split(/\s+/);
  const initials = parts
    .slice(0, 2)
    .map((p) => p[0] ?? "")
    .join("");
  return (
    <div className="h-10 w-10 rounded-full bg-gray-100 flex items-center justify-center text-sm font-medium text-gray-700">
      {initials || "?"}
    </div>
  );
}

function RoleBadge({ role }: { role: string }) {
  let color = "bg-gray-100 text-gray-800";
  if (role === "Admin") color = "bg-red-100 text-red-800";
  if (role === "Advisor") color = "bg-blue-100 text-blue-800";
  if (role === "Student") color = "bg-green-100 text-green-800";

  return (
    <span
      className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${color}`}
    >
      {role}
    </span>
  );
}
function StatusToggle({
  userId,
  isActive,
  onToggle,
}: {
  userId: string;
  isActive: boolean;
  onToggle: (sutId: string, status: boolean) => void;
}) {
  return (
    <label className="relative inline-flex items-center cursor-pointer">
      <input
        type="checkbox"
        checked={isActive}
        onChange={() => onToggle(userId, !isActive)}
        className="sr-only peer"
      />
      <div className="w-9 h-5 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-[#F26522]/30 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-[#F26522]"></div>
      <span className="ml-3 text-xs font-medium text-gray-900">
        {isActive ? "ใช้งานอยู่" : "ปิดการใช้งาน"}
      </span>
    </label>
  );
}

function formatDate(dateString: string) {
  if (!dateString) return "-";
  return new Date(dateString).toLocaleDateString("th-TH", {
    day: "numeric",
    month: "short",
    year: "numeric",
  });
}
// --- สิ้นสุด Components ย่อย ---

// --- Main Page Component ---
export default function ManageUsersPage() {
  const { isAuthenticated, isLoading, user } = useAuth();
  const router = useRouter();

  // States ใช้ Type ใหม่ (ManagedUserEntry)
  const [users, setUsers] = useState<UserData[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [search, setSearch] = useState("");
  const [filters, setFilters] = useState<UserFilters>({
    role: "",
    status: "",
    major: "",
  });

  const [page, setPage] = useState(1);
  const pageSize = 10;

  const [editingUserId, setEditingUserId] = useState<string | null>(null);
  const [detailData, setDetailData] = useState<AdminUserDetailResponse | null>(
    null
  );
  const [loadingDetail, setLoadingDetail] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [majorOptions, setMajorOptions] = useState<MajorEntry[]>([]);

  // --- Handlers ---

  const handleFilterChange = (name: keyof UserFilters, value: string) => {
    setFilters((prev) => ({ ...prev, [name]: value }));
    setPage(1);
  };

  /**
   * จัดการการสลับสถานะผู้ใช้ด้วย Optimistic Update และ Rollback หาก API ล้มเหลว
   */
  const handleStatusToggle = useCallback(
    async (sutId: string, newActiveState: boolean) => {
      // 1. กำหนดสถานะใหม่เป็น String ตามที่ Backend คาดหวัง
      const newStatusString = newActiveState ? "active" : "inactive";

      // เก็บสถานะเดิมไว้เผื่อต้อง Rollback
      const userToToggle = users.find((u) => u.sutId === sutId);
      const oldActiveState = userToToggle?.active;

      if (oldActiveState === undefined) return;

      // 2. Optimistic Update (อัปเดต UI ก่อนเรียก API)
      setUsers((prevUsers) =>
        prevUsers.map((u) =>
          u.sutId === sutId
            ? { ...u, active: newActiveState, is_active: newActiveState }
            : u
        )
      );

      try {
        // 3. API Call จริง
        // ใช้ Type Assertion เพื่อให้ TypeScript ยอมรับว่าค่าจะถูกต้อง
        await updateUserStatus(
          sutId,
          newStatusString as "active" | "inactive" | "suspended"
        );

        console.log(
          `User ${sutId} status successfully updated to: ${newStatusString}`
        );
      } catch (e: any) {
        // 4. Rollback State หาก API Call ล้มเหลว
        console.error("Failed to update status via API:", e);

        // Rollback state กลับไปเป็นสถานะเดิมที่เก็บไว้
        setUsers((prevUsers) =>
          prevUsers.map((u) =>
            u.sutId === sutId
              ? { ...u, active: oldActiveState, is_active: oldActiveState } // เปลี่ยนกลับไปใช้สถานะเก่า
              : u
          )
        );
        alert(`การอัปเดตสถานะล้มเหลว: ${e.message}`);
      }
    },
    [users]
  );

  const handleResetPassword = async () => {
    if (!detailData) return;

    if (!confirm(`ยืนยันการ Reset รหัสผ่านสำหรับ ${detailData.fullName}?`)) {
      return;
    }

    try {
      // API Call: await resetUserPassword(detailData.sutId);
      alert("รหัสผ่านถูก Reset เรียบร้อยแล้ว (รหัสผ่านชั่วคราว)");
      setEditingUserId(null);
    } catch (e) {
      alert("Reset รหัสผ่านล้มเหลว");
    }
  };

  const handleSaveChanges = async (formData: AdminUserDetailResponse) => {
    try {
      await updateManagedUser(formData.sutId, {
        phone: formData.phone,
      });

      // อัปเดตเฉพาะ user ที่แก้
      setUsers((prev) =>
        prev.map((u) =>
          u.sutId === formData.sutId ? { ...u, phone: formData.phone } : u
        )
      );

      alert("บันทึกข้อมูลเรียบร้อยแล้ว");
      setEditingUserId(null);
    } catch (e: any) {
      alert("บันทึกข้อมูลล้มเหลว: " + e.message);
    }
  };

  // --- Effects ---

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push("/login");
    }
  }, [isLoading, isAuthenticated, router]);

  useEffect(() => {
    if (isLoading || !isAuthenticated) return;
    const fetchMajorsList = async () => {
      try {
        const data = await getMajors();
        setMajorOptions(data);
      } catch (e) {
        console.error("Failed to load majors:", e);
      }
    };
    fetchMajorsList();
  }, [isLoading, isAuthenticated]);

  // Fetch User List (ใช้ API จริง: getManagedUsers)
  useEffect(() => {
    if (isLoading || !isAuthenticated) return;
    const fetchList = async () => {
      setLoading(true);
      setError(null);
      try {
        const response = await getManagedUsers(filters);
        const usersWithCreationDatePromises = response.users.map(async (u) => {
          try {
            const createdDateResponse = await getUserCreatedDate(u.sutId);
            return {
              ...u,
              id: u.sutId, // ใช้ sutId เป็น id (string)
              is_active: u.active, // แมป active is_active
              createdAt: createdDateResponse.created_at,
            } as UserData;
          } catch (e) {
            console.error(`Failed to get created date for ${u.sutId}:`, e);
            return {
              ...u,
              id: u.sutId,
              is_active: u.active,
              createdAt: "",
            } as UserData;
          }
        });
        const usersData = await Promise.all(usersWithCreationDatePromises);
        setUsers(usersData);
      } catch (err: any) {
        if (err.message.includes("Unauthorized")) {
          router.push("/login");
          return;
        }
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };
    fetchList();
  }, [filters, isLoading, isAuthenticated, router]);

  // useEffect: Fetch User Detail (ใช้ API จริง: getUserDetailBySutId)
  useEffect(() => {
    if (!editingUserId) {
      setDetailData(null);
      setIsEditing(false);
      return;
    }

    const fetchDetail = async () => {
      setLoadingDetail(true);
      setError(null);
      try {
        const d = await getUserDetailBySutId(editingUserId);
        setDetailData(d);
        setIsEditing(false);
      } catch (err: any) {
        setDetailData(null);
        setError(err.message);
      } finally {
        setLoadingDetail(false);
      }
    };
    fetchDetail();
  }, [editingUserId]);

  // --- Filtering & Pagination Logic ---

  const filtered = useMemo(() => {
    const q = search.trim().toLowerCase();
    return users.filter((u) => {
      if (q === "") return true;
      return [u.fullName, u.sutId, u.email]
        .filter(Boolean)
        .some((field) => field!.toString().toLowerCase().includes(q));
    });
  }, [users, search]);

  const totalPages = Math.max(1, Math.ceil(filtered.length / pageSize));
  const pageData = filtered.slice((page - 1) * pageSize, page * pageSize);

  // --- Render Components ---

  const UserDetailModal = () => {
    if (!editingUserId || !detailData) return null;

    const [formData, setFormData] =
      useState<AdminUserDetailResponse>(detailData);
    const isSuperAdmin = user?.role === "Admin";

    const getFullName = () => detailData.fullName || "ผู้ใช้งาน";

    const handleFormChange = (
      e: React.ChangeEvent<
        HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
      >
    ) => {
      const { name, value } = e.target;
      setFormData(
        (prev) => ({ ...prev, [name]: value } as AdminUserDetailResponse)
      );
    };

    return (
      <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
        <div className="w-full max-w-2xl rounded-lg bg-white p-6 shadow-lg text-slate-900">
          <div className="flex items-start justify-between gap-4">
            <h2 className="text-lg font-semibold">
              {isEditing ? "แก้ไขข้อมูลผู้ใช้" : "รายละเอียดผู้ใช้"}
            </h2>
            <button
              onClick={() => setEditingUserId(null)}
              className="text-sm text-gray-500"
            >
              ปิด
            </button>
          </div>

          <div className="mt-4">
            {loadingDetail ? (
              <div className="text-center text-gray-500">กำลังโหลด...</div>
            ) : (
              <>
                {/* --------------------- ส่วนแสดง/แก้ไขข้อมูล --------------------- */}
                <form
                  onSubmit={(e) => {
                    e.preventDefault();
                    handleSaveChanges(formData);
                  }}
                >
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {/* ชื่อ */}
                    <div className="md:col-span-2 flex items-center gap-3">
                      <div className="h-16 w-16 rounded-full bg-gray-200 flex items-center justify-center text-xl text-gray-600">
                        {getFullName()
                          .trim()
                          .split(/\s+/)
                          .map((p) => p[0] ?? "")
                          .slice(0, 2)
                          .join("") || "?"}
                      </div>
                      <div className="flex-1">
                        <div className="text-lg font-medium">
                          {getFullName()}
                        </div>
                        <div className="text-sm text-gray-500">
                          {formData.sutId} / {formData.email}
                        </div>
                      </div>
                      {/* <div className="flex items-center gap-3">
                        <RoleBadge role={formData.role} />
                        <StatusToggle
                          userId={formData.sutId}
                          isActive={formData.active}
                          onToggle={handleStatusToggle}
                        />
                      </div> */}
                    </div>

                    {/* Field: Phone */}
                    <div>
                      <label className="block text-sm font-medium text-gray-700">
                        เบอร์โทร
                      </label>
                      <input
                        type="text"
                        name="phone"
                        value={formData.phone || ""}
                        onChange={handleFormChange}
                        disabled={!isEditing}
                        className="mt-1 w-full rounded-xl border border-gray-300 px-3 py-2 text-sm text-gray-800 bg-white focus:ring-2 focus:ring-[#F26522] outline-none transition-all disabled:bg-gray-50"
                      />
                    </div>

                    {/* Field: /Major */}
                    <div>
                      <label className="block text-sm font-medium text-gray-700">
                        ภาควิชา/สาขา
                      </label>
                      <p className="mt-1 w-full rounded-xl border border-gray-300 px-3 py-2 text-sm text-gray-800 bg-white focus:ring-2 focus:ring-[#F26522] outline-none transition-all disabled:bg-gray-50">
                        {formData.majorName || "-"}
                      </p>
                    </div>

                    {/* --------------------- ส่วนจัดการ Role & Advisor --------------------- */}
                    <div className="md:col-span-2 mt-2">
                      <div className="grid grid-cols-2 gap-4">
                        {/* Field: Advisor (สำหรับ Student เท่านั้น) */}
                        {formData.role === "Student" && (
                          <div>{/* ... (Assign Advisor logic) ... */}</div>
                        )}

                        {/* Field: Student Count (สำหรับ Advisor เท่านั้น) */}
                        {/* ... */}
                      </div>
                    </div>
                  </div>

                  {/* --------------------- ปุ่ม Action --------------------- */}
                  <div className="mt-6 flex justify-end gap-3">
                    {isEditing ? (
                      <>
                        <button
                          type="button"
                          onClick={() => {
                            setIsEditing(false);
                            setFormData(detailData);
                          }}
                          className="rounded-md border px-4 py-2 text-sm text-gray-700 bg-white"
                        >
                          ยกเลิก
                        </button>
                        <button
                          type="submit"
                          className="rounded-md bg-[#F26522] px-4 py-2 text-sm text-white"
                        >
                          บันทึก
                        </button>
                      </>
                    ) : (
                      <>
                        <button
                          type="button"
                          onClick={() => setEditingUserId(null)}
                          className="rounded-md border px-4 py-2 text-sm text-gray-700 bg-white"
                        >
                          ปิด
                        </button>
                        <button
                          type="button"
                          onClick={() => setIsEditing(true)}
                          className="rounded-md bg-[#F26522] px-4 py-2 text-sm text-white"
                        >
                          แก้ไขข้อมูล
                        </button>
                      </>
                    )}
                  </div>
                </form>
              </>
            )}
          </div>
        </div>
      </div>
    );
  };

  // --- Main Render ---

  if (isLoading) return <div className="p-8">กำลังตรวจสอบสิทธิ์...</div>;
  if (!isAuthenticated || user?.role !== "Admin")
    return <div className="p-8 text-red-500">คุณไม่มีสิทธิ์เข้าถึงหน้านี้</div>;

  return (
    <ConfigProvider
      theme={{
        token: {
          colorPrimary: "#F26522",
        },
      }}
    >
      <div className="space-y-6 pd-6">
        {/* Header and Filters */}
        <div className="flex items-center gap-3 flex-wrap w-full md:w-auto">
          {/* Search Input: w-60 */}
          <div className="relative w-full sm:w-60 flex-shrink-0">
            {/* Search Icon & Input */}
            <input
              type="text"
              placeholder="ค้นหา (ชื่อ, ID, Email)"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="pl-10 pr-4 py-2 w-full rounded-xl border border-gray-300 text-sm 
                             transition duration-150"
            />
          </div>

          {/* Filter: Role: w-36 */}
          <div className="w-full sm:w-36 flex-shrink-0">
            <Select
              value={filters.role}
              style={{ width: "100%" }}
              placeholder="ทุก Role"
              onChange={(value) => handleFilterChange("role", value)}
              options={[
                { value: "", label: "ทุก Role" },
                ...mockRoles.map((r) => ({ value: r, label: r })),
              ]}
              className="custom-select"
              styles={{
                popup: {
                  root: {
                    borderRadius: "8px",
                  },
                },
              }}
            />
          </div>

          {/* Filter: Status: w-36 */}
          <div className="w-full sm:w-36 flex-shrink-0">
            <Select
              value={filters.status}
              style={{ width: "100%" }}
              placeholder="ทุกสถานะ"
              onChange={(value) => handleFilterChange("status", value)}
              options={[
                { value: "", label: "ทุกสถานะ" },
                { value: "active", label: "Active" },
                { value: "inactive", label: "Inactive" },
              ]}
              className="custom-select"
              styles={{
                popup: {
                  root: {
                    borderRadius: "8px",
                  },
                },
              }}
            />
          </div>

          {/* Filter: (Major): w-44 */}
          <div className="w-full sm:w-44 flex-shrink-0">
            <Select
              value={filters.major}
              style={{ width: "100%" }}
              placeholder="ทุกสาขาวิชา"
              onChange={(value) => handleFilterChange("major", value)}
              options={[
                { value: "", label: "ทุกสาขาวิชา" },
                ...majorOptions.map((m) => ({
                  value: m.major,
                  label: m.major,
                })),
              ]}
              className="custom-select"
              styles={{
                popup: {
                  root: {
                    borderRadius: "8px",
                  },
                },
              }}
            />
          </div>
          {/* แสดงผลรวมผู้ใช้ทั้งหมด */}
          <div className="flex-shrink-0 text-sm font-medium text-gray-600 px-4 py-2 rounded-xl  ml-auto flex items-center h-11">
            แสดง {filtered.length} ผู้ใช้
          </div>
        </div>

        {/* 3. ตารางผู้ใช้ (Scrollable Container) */}
        <div className="overflow-hidden rounded-xl border border-slate-100 bg-white flex flex-col flex-1 min-h-0 ">
          <div className="overflow-y-auto flex-1 overflow-x-auto">
            <table className="w-full min-w-full border-collapse text-sm">
              <thead className="bg-slate-50 sticky top-0 z-10">
                <tr>
                  <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-500 w-2/6">
                    ชื่อ – ID / Email
                  </th>
                  <th className="px-4 py-3 text-center text-xs font-semibold uppercase tracking-wide text-slate-500 w-[15%]">
                    Role
                  </th>
                  <th className="px-4 py-3 text-center text-xs font-semibold uppercase tracking-wide text-slate-500 w-[15%]">
                    สถานะ
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-500 w-[15%]">
                    วันที่สร้างบัญชี
                  </th>
                  <th className="px-4 py-3 text-right text-xs font-semibold uppercase tracking-wide text-slate-500 w-[15%]">
                    Action
                  </th>
                </tr>
              </thead>
              <tbody>
                {loading ? (
                  <tr>
                    <td
                      colSpan={5}
                      className="px-4 py-12 text-center text-gray-500"
                    >
                      กำลังโหลด...
                    </td>
                  </tr>
                ) : error ? (
                  <tr>
                    <td
                      colSpan={5}
                      className="px-4 py-12 text-center text-red-500"
                    >
                      {error}
                    </td>
                  </tr>
                ) : pageData.length === 0 ? (
                  <tr>
                    <td
                      colSpan={5}
                      className="px-4 py-12 text-center text-gray-500"
                    >
                      ไม่พบผู้ใช้ที่ตรงตามเงื่อนไข
                    </td>
                  </tr>
                ) : (
                  pageData.map((u) => (
                    <tr
                      key={u.sutId}
                      className="border-t border-slate-100 hover:bg-slate-50/60"
                    >
                      {/* คอลัมน์ 1: ชื่อ – ID / Email */}
                      <td className="px-4 py-3 align-middle w-2/6">
                        <div className="flex items-center gap-3">
                          <Avatar name={u.fullName} />
                          <div className="flex flex-col">
                            <span className="font-medium text-slate-900">
                              {u.fullName}
                            </span>
                            <span className="text-xs text-slate-500">
                              {u.sutId ?? u.email}
                            </span>
                          </div>
                        </div>
                      </td>
                      {/* คอลัมน์ 2: Role */}
                      <td className="px-4 py-3 align-middle w-[15%] text-center">
                        <div className="inline-flex justify-center">
                          <RoleBadge role={u.role} />
                        </div>
                      </td>
                      {/* คอลัมน์ 3: สถานะ */}
                      <td className="px-4 py-3 align-middle w-[15%] text-center">
                        <div className="inline-flex justify-center">
                          <StatusToggle
                            userId={u.sutId}
                            isActive={u.active}
                            onToggle={handleStatusToggle}
                          />
                        </div>
                      </td>
                      {/* คอลัมน์ 4: วันที่สร้างบัญชี */}
                      <td className="px-4 py-3 align-middle text-slate-700 w-[15%]">
                        {u.createdAt ? formatDate(u.createdAt) : "-"}
                      </td>
                      {/* คอลัมน์ 5: Action */}
                      <td className="px-4 py-3 align-middle text-right w-[15%]">
                        <button
                          onClick={() => setEditingUserId(u.sutId)}
                          className="text-sm rounded-lg border px-3 py-1 text-blue-600 border-blue-600 hover:bg-blue-50 transition"
                        >
                          จัดการ
                        </button>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>

        {/* 4. การแบ่งหน้า (Pagination) */}
        <div className="flex items-center justify-between flex-shrink-0 mt-4 px-4 py-3">
          <div className="text-sm text-gray-700">
            หน้า {page} จาก {totalPages} ({filtered.length} รายการ)
          </div>
          <div className="flex items-center gap-2">
            <button
              disabled={page === 1}
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              className="px-3 py-1 rounded-md border disabled:opacity-50 text-gray-700 bg-white hover:bg-gray-50"
            >
              ก่อนหน้า
            </button>
            <button
              disabled={page === totalPages}
              onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
              className="px-3 py-1 rounded-md border disabled:opacity-50 text-gray-700 bg-white hover:bg-gray-50"
            >
              ถัดไป
            </button>
          </div>
        </div>

        {/* Modal: รายละเอียด/แก้ไขผู้ใช้ */}
        <UserDetailModal />
      </div>
    </ConfigProvider>
  );
}
