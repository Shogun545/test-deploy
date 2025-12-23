export interface ManagedUserEntry {
    sutId: string;
    prefix: string; 
    firstName: string; 
    lastName: string; 
    fullName: string;
    email: string;
    role: string;
    majorName: string;      // สาขาวิชา
    active: boolean;
    createdAt?: string; 
    
    // Field เสริมที่อาจถูกเพิ่มใน Frontend เพื่อใช้เป็น Key หรือสำหรับการแสดงผล
    id?: string; 
}


// AdminManagedUsersResponse (อิงตาม Go Backend DTO)
export interface AdminManagedUsersResponse {
    totalUsers: number;
    users: ManagedUserEntry[];
}

// UserFilters (ใช้สำหรับส่ง Query Params ใน Request List)
export interface UserFilters {
    role?: string;
    status?: string;
    major?: string;
    search?: string; // สำหรับค้นหา (SutID, Name, Email)
}


// AdminUserDetailResponse ใช้สำหรับแสดงรายละเอียดผู้ใช้ใน Modal
export interface AdminUserDetailResponse {
    sutId: string;
    fullName: string;
    email: string;
    phone: string;
    role: string;
    departmentName: string;
    majorName: string;
    active: boolean;

    // ข้อมูลเฉพาะอาจารย์ (optional)
    officeRoom?: string;   
    specialties?: string;  

    // ข้อมูลเฉพาะนักศึกษา (optional)
    yearOfStudy?: number; 
    gpaLatest?: number;   
}


export interface MajorEntry {
    major: string;
}

export interface UpdateUserStatusRequest {
    status: string; // เช่น "Active", "Inactive", "Suspended"
}

export interface UserCreatedDateResponse {
  sut_id: string
  created_at: string 
}