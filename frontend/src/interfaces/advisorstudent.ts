export type StudentForAdvisor = {
  sut_id: string;
  first_name: string;
  last_name: string;
  year_of_study?: number | null;
  email?: string | null;
  phone?: string | null;
  major_name?: string | null;
  gpa_latest?: number | null;
  birthday?: string | null;
  national_id?: string | null;
  advisor_note?: string | null;
  // ถ้ามี id ใน DB ให้เพิ่ม field id: number;
};

export type AdvisorStudentsResponse = {
  advisor_sut_id: string;
  advisor_name: string;
  students: StudentForAdvisor[];
};