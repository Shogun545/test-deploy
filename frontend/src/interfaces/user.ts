export interface User {
  id: number;
  sut_id: string;   // หรือ sut_id ถ้าอยากใช้ field นี้
  sutId: string;   // หรือ sut_id ถ้าอยากใช้ field นี้
  role: string;       // "Admin" | "Advisor" | "Student"
  profile_image?: string | null;
}

export interface LoginInterface {
  sutId?: string;
  sut_id?: string;
  password: string;
}