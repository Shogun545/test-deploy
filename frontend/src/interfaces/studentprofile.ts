export interface StudentProfileResponse {
  sut_id: string;
  prefix: string;
  first_name: string;
  last_name: string;
  email: string;
  phone: string;
  national_id: string;
  birthday: string;
  profile_image: string | null;

  major_name: string;
  year_of_study: number;
  advisor_name: string | null;

  academic_year: string | null;
  semester: number | null;
  term_gpa: number | null;
  cumulative_gpa: number | null;
  academic_status: string | null;

  gpx: number | null;
  gpa_latest: number | null;
  gpa_term_label: string | null;
}

export interface ProfileUI {
  title: string;
  firstName: string;
  lastName: string;
  roleLabel: string;

  studentCode: string;
  advisorName: string;

  email: string;
  phone: string;
  dob: string;
  idCard: string;

  major: string;
  year: string;

  gpx: number;
  gpaLatest: number;
  gpaTermLabel: string;

  status: string;
  statusNote: string;
}
