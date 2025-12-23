"use client";

import { useEffect, useMemo, useState } from "react";
import { useAuth } from "@/src/contexts/AuthContext";
import { useRouter } from "next/navigation";
import {
  getAdvisorStudents,
  getAdvisorStudentDetail,
} from "@/src/services/http/advisorstudents";
import type { StudentForAdvisor } from "@/src/interfaces/advisorstudent";
import { Select, ConfigProvider } from "antd"; // นำเข้า ConfigProvider

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

export default function MyStudentsPage() {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();

  const [students, setStudents] = useState<StudentForAdvisor[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [search, setSearch] = useState("");
  const [selectedMajor, setSelectedMajor] = useState<string>("");
  const [majors, setMajors] = useState<string[]>([]);

  // pagination
  const [page, setPage] = useState(1);
  const pageSize = 8;

  // detail modal
  const [detailSutId, setDetailSutId] = useState<string | null>(null);
  const [detailData, setDetailData] = useState<StudentForAdvisor | null>(null);
  const [loadingDetail, setLoadingDetail] = useState(false);

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push("/login");
    }
  }, [isLoading, isAuthenticated, router]);

  // fetch list
  useEffect(() => {
    if (isLoading || !isAuthenticated) return; // รอการยืนยันตัวตนก่อน
    const fetchList = async () => {
      setLoading(true);
      setError(null);
      try {
        const data = await getAdvisorStudents(
          selectedMajor ? { major: selectedMajor } : undefined
        );
        setStudents(data.students || []);

        const uniqueMajors = Array.from(
          new Set(
            (data.students || []).map((s) => s.major_name).filter(Boolean)
          )
        ) as string[];
        setMajors(uniqueMajors);
        setPage(1);
      } catch (err: any) {
        setError(
          err?.response?.data?.message ||
            err.message ||
            "โหลดรายชื่อนักศึกษาไม่สำเร็จ"
        );
      } finally {
        setLoading(false);
      }
    };

    fetchList();
  }, [selectedMajor, isLoading, isAuthenticated]);

  // fetch detail
  useEffect(() => {
    if (!detailSutId) {
      setDetailData(null);
      return;
    }

    const fetchDetail = async () => {
      setLoadingDetail(true);
      try {
        const d = await getAdvisorStudentDetail(detailSutId);
        setDetailData(d);
      } catch (err: any) {
        setDetailData({
          first_name: "",
          last_name: "",
          sut_id: detailSutId,
        } as StudentForAdvisor);
        setError(
          err?.response?.data?.message || err.message || "โหลดข้อมูลไม่สำเร็จ"
        );
      } finally {
        setLoadingDetail(false);
      }
    };

    fetchDetail();
  }, [detailSutId]);

  const filtered = useMemo(() => {
    const q = search.trim().toLowerCase();
    return students.filter((s) => {
      if (selectedMajor && s.major_name !== selectedMajor) return false;
      if (q === "") return true;
      return [s.first_name, s.last_name, s.sut_id, s.major_name]
        .filter(Boolean)
        .some((field) => field!.toString().toLowerCase().includes(q));
    });
  }, [students, search, selectedMajor]);

  const totalPages = Math.max(1, Math.ceil(filtered.length / pageSize));
  const pageData = filtered.slice((page - 1) * pageSize, page * pageSize);

  const handleMajorChange = (value: string) => {
    setSelectedMajor(value);
  };


  return (
    <ConfigProvider
      theme={{
        token: {
          colorPrimary: '#F26522', 
        },
      }}
    >
      {/* เพิ่ม pb-4 เพื่อไม่ให้ Pagination ติดขอบ และใช้ h-full สำหรับ Layout ที่ยืดหยุ่น */}
      <div className="flex flex-col h-full space-y-4 pb-4">
        {/* ส่วนหัว / ตัวกรอง (ใช้ flex-shrink-0) */}
        <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-3 flex-shrink-0">
          <div className="flex items-center gap-3">
            <label className="text-sm text-gray-600">สาขาวิชา</label>
            
            <Select
              value={selectedMajor}
              style={{ width: 150 }}
              onChange={handleMajorChange}
              options={[
                { value: "", label: "ทุกสาขาวิชา" },
                ...majors.map((m) => ({ value: m, label: m })),
              ]}
              // ⬇️ [1] เพิ่ม showSearch เพื่อบังคับใช้ Custom UI (แก้ปัญหาฟอนต์)
              showSearch={true} 
              filterOption={false}
              className="text-gray-700" 
            />
            
            {/* <button
              onClick={async () => {
                setLoading(true);
                try {
                  const data = await getAdvisorStudents(
                    selectedMajor ? { major: selectedMajor } : undefined
                  );
                  setStudents(data.students || []);
                  setPage(1);
                } catch (err: any) {
                  setError(
                    err?.response?.data?.message ||
                      err.message ||
                      "โหลดรายชื่อนักศึกษาไม่สำเร็จ"
                  );
                } finally {
                  setLoading(false);
                }
              }}
              className="rounded-full bg-[#F26522] px-4 py-1 text-white text-sm hover:bg-[#e05a1d]"
            >
              ดึงข้อมูล
            </button> */}
          </div>

          <div className="text-sm text-gray-500">
            แสดง {filtered.length} ผลลัพธ์
          </div>
        </div>

        {/* ตารางนักศึกษา (Scrollable Container) */}
        <div className="overflow-hidden rounded-xl border border-slate-100 bg-white flex flex-col flex-1 min-h-0">
          {/* ส่วนหัวตาราง: คงที่ (flex-shrink-0) */}
          <div className="flex-shrink-0 overflow-x-auto">
            {/* เพิ่ม table-fixed และกำหนดความกว้างคอลัมน์เพื่อ Alignment ที่แม่นยำ */}
            <table className="w-full min-w-full border-collapse text-sm table-fixed">
              <thead className="bg-slate-50 sticky top-0 z-10">
                <tr>
                  <th className="w-2/5 px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-500">
                    นักศึกษา
                  </th>
                  <th className="w-1/5 px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-500">
                    สาขา
                  </th>
                  <th className="w-[10%] px-4 py-3 text-center text-xs font-semibold uppercase tracking-wide text-slate-500">
                    ชั้นปี
                  </th>
                  <th className="w-[10%] px-4 py-3 text-center text-xs font-semibold uppercase tracking-wide text-slate-500">
                    GPA
                  </th>
                  <th className="w-1/5 px-4 py-3 text-right text-xs font-semibold uppercase tracking-wide text-slate-500">
                    การดำเนินการ
                  </th>
                </tr>
              </thead>
            </table>
          </div>

          {/* ส่วนเนื้อหาตาราง: เลื่อนได้และกินพื้นที่เหลือทั้งหมด */}
          <div className="overflow-y-auto flex-1 overflow-x-auto">
            {/* เพิ่ม table-fixed และกำหนดความกว้างคอลัมน์เพื่อ Alignment ที่แม่นยำ */}
            <table className="w-full min-w-full border-collapse text-sm table-fixed">
              <tbody>
                {loading ? (
                  <tr>
                    <td
                      colSpan={5}
                      className="px-4 py-6 text-center text-gray-500"
                    >
                      กำลังโหลด...
                    </td>
                  </tr>
                ) : error ? (
                  <tr>
                    <td
                      colSpan={5}
                      className="px-4 py-6 text-center text-red-500"
                    >
                      {error}
                    </td>
                  </tr>
                ) : pageData.length === 0 ? (
                  <tr>
                    <td
                      colSpan={5}
                      className="px-4 py-6 text-center text-gray-500"
                    >
                      ไม่มีนักศึกษาที่ตรงกับการค้นหา
                    </td>
                  </tr>
                ) : (
                  pageData.map((s) => (
                    <tr
                      key={s.sut_id}
                      className="border-t border-slate-100 hover:bg-slate-50/60"
                    >
                      {/* คอลัมน์ 1: นักศึกษา */}
                      <td className="w-2/5 px-4 py-3 align-middle">
                        <div className="flex items-center gap-3">
                          <Avatar name={`${s.first_name} ${s.last_name}`} />
                          <div className="flex flex-col">
                            <span className="font-medium text-slate-900">
                              {s.first_name} {s.last_name}
                            </span>
                            <span className="text-xs text-slate-500">
                              {s.sut_id}
                            </span>
                          </div>
                        </div>
                      </td>
                      {/* คอลัมน์ 2: สาขา */}
                      <td className="w-1/5 px-4 py-3 align-middle text-slate-700">
                        {s.major_name ?? "-"}
                      </td>
                      {/* คอลัมน์ 3: ชั้นปี */}
                      <td className="w-[10%] px-4 py-3 align-middle text-center text-slate-700">
                        {s.year_of_study ?? "-"}
                      </td>
                      {/* คอลัมน์ 4: GPA */}
                      <td className="w-[10%] px-4 py-3 align-middle text-center text-slate-700">
                        {s.gpa_latest != null ? s.gpa_latest.toFixed(2) : "-"}
                      </td>
                      {/* คอลัมน์ 5: การดำเนินการ (ใช้ text-right) */}
                      <td className="w-1/5 px-4 py-3 align-middle text-right">
                        <div className="flex items-center justify-end gap-2">
                          <button
                            onClick={() => setDetailSutId(s.sut_id)}
                            className="text-sm rounded-lg border px-3 py-1 text-[#F26522] border-[#F26522] hover:bg-orange-50"
                          >
                            ดู
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>

        {/* การแบ่งหน้า (flex-shrink-0) */}
        <div className="flex items-center justify-between flex-shrink-0">
          {/* แก้ไขสีข้อความสำหรับ Pagination Status */}
          <div className="text-sm text-gray-700">
            หน้า {page} จาก {totalPages}
          </div>
          <div className="flex items-center gap-2">
            {/* แก้ไขสีข้อความและพื้นหลังสำหรับปุ่ม Pagination */}
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

        {/* โมดัลรายละเอียด (บังคับ Light Mode เสมอ) */}
        {detailSutId != null && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
            {/* ⬇️ [2] ลบคลาส Dark Mode ออกจาก Modal หลัก (ใช้ Light Mode เสมอ) ⬇️ */}
            <div className="w-full max-w-2xl rounded-lg bg-white p-6 shadow-lg text-slate-900">
              <div className="flex items-start justify-between gap-4">
                <h2 className="text-lg font-semibold">รายละเอียดนักศึกษา</h2>
                {/* ⬇️ [3] ลบคลาส Dark Mode ออกจากปุ่มปิด ⬇️ */}
                <button
                  onClick={() => setDetailSutId(null)}
                  className="text-sm text-gray-500"
                >
                  ปิด
                </button>
              </div>
              <div className="mt-4">
                {loadingDetail ? (
                  <div>กำลังโหลดข้อมูล...</div>
                ) : detailData ? (
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                      <div className="flex items-center gap-3">
                        {/* ⬇️ [4] ลบคลาส Dark Mode ออกจาก Avatar ⬇️ */}
                        <div className="h-16 w-16 rounded-full bg-gray-200 flex items-center justify-center text-xl text-gray-600">
                          {detailData.first_name?.[0] ?? ""}
                          {detailData.last_name?.[0] ?? ""}
                        </div>
                        <div>
                          <div className="text-lg font-medium">
                            {detailData.first_name} {detailData.last_name}
                          </div>
                          {/* ⬇️ [5] ลบคลาส Dark Mode ออกจากรหัสนักศึกษา ⬇️ */}
                          <div className="text-sm text-gray-500">
                            {detailData.sut_id}
                          </div>
                        </div>
                      </div>

                      <div className="mt-4 text-sm">
                        {/* รายละเอียดอื่นๆ ใช้สีหลัก (สีดำ) */}
                        <p>
                          <span className="font-medium">สาขาวิชา:</span>{" "}
                          {detailData.major_name ?? "-"}
                        </p>
                        <p className="mt-1">
                          <span className="font-medium">ชั้นปี:</span>{" "}
                          {detailData.year_of_study ?? "-"}
                        </p>
                        <p className="mt-1">
                          <span className="font-medium">อีเมล:</span>{" "}
                          {detailData.email ?? "-"}
                        </p>
                        <p className="mt-1">
                          <span className="font-medium">เบอร์โทร:</span>{" "}
                          {detailData.phone ?? "-"}
                        </p>
                      </div>
                    </div>

                    <div>
                      <p className="text-sm">
                        <span className="font-medium">GPA (ล่าสุด):</span>{" "}
                        {detailData.gpa_latest != null
                          ? detailData.gpa_latest.toFixed(2)
                          : "-"}
                      </p>
                      <p className="mt-1 text-sm">
                        <span className="font-medium">วันเกิด:</span>{" "}
                        {detailData.birthday ?? "-"}
                      </p>
                      

                      {detailData.advisor_note && (
                        <div className="mt-3">
                          <p className="font-medium text-sm">หมายเหตุที่ปรึกษา</p>
                          {/* ⬇️ [6] ลบคลาส Dark Mode ออกจาก Note Text ⬇️ */}
                          <p className="mt-1 text-sm text-gray-700">
                            {detailData.advisor_note}
                          </p>
                        </div>
                      )}
                    </div>
                  </div>
                ) : (
                  <div>ไม่มีข้อมูล</div>
                )}
              </div>

              <div className="mt-6 flex justify-end gap-3">
                {/* ⬇️ [7] ลบคลาส Dark Mode ออกจากปุ่มปิด ⬇️ */}
                <button
                  onClick={() => setDetailSutId(null)}
                  className="rounded-md border px-4 py-2 text-sm text-gray-700 bg-white"
                >
                  ปิด
                </button>
                <button
                  onClick={() => {
                    window.location.href = `/advisor/students/${detailSutId}`;
                  }}
                  className="rounded-md bg-[#F26522] px-4 py-2 text-sm text-white"
                >
                  ไปที่หน้ารายละเอียด
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </ConfigProvider>
  );
}