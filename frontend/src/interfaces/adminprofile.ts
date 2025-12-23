export interface AdminProfileResponse {
  sutId: string;
  firstName: string;
  lastName: string;
  prefix: string;
  email: string;
  phone: string;
  role: string;
  departmentName: string;
  notes: string;
  active: boolean;
  createdAt: string;
  managedUsers?: number;
  lastLogin?: string;
  totalManaged?: number;
}


// ข้อมูลที่ Component จะส่งไปเพื่ออัปเดต (PATCH payload)
export type AdminProfileUpdateData = {
  prefix: string;
  firstName: string;
  lastName: string; 
  phone: string;
  notes: string;
  active: boolean;
  email: string;
};
