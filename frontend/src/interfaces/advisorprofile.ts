export interface AdvisorProfileResponse {
  prefix: string;
  first_name: string;
  last_name: string;
  email: string;
  phone: string;
  office_room: string;
  major_name: string;
  department_name: string;
  advisee_count: number;
  specialties: string;
  is_active: boolean;
  profile_image?: string | null;
}

export interface AdvisorProfileUI {
  fullName: string;
  email: string;
  phone: string;
  officeRoom: string;

  majorName: string;
  departmentName: string;
  adviseeCount: number;

  specialties: string;
  isActive: boolean;
   profileImage?: string | null;
}

