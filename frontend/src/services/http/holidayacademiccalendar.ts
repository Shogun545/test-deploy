// src/services/http/holidayacademic.ts

export interface Holiday {
  date: string;
  localName: string;
  name: string;
  countryCode: string;
  types: string[];
  // field อื่นๆ ที่อาจจะไม่ได้ใช้แล้วสามารถตัดออกหรือทำให้เป็น optional ได้
  fixed?: boolean;
  global?: boolean;
  counties?: string[] | null;
  launchYear?: number | null;
}

export const fetchThaiHolidays = async (year: number): Promise<Holiday[]> => {
  if (!year) {
    console.error("Year is undefined!");
    return [];
  }

  // ** เปลี่ยน URL ตรงนี้ ** // ชี้ไปที่ Go Backend ของคุณ (เช่น http://localhost:8080/api/holidays)
  // หรือถ้าตั้ง Proxy ไว้ใน next.config.js แล้ว ให้ใช้ path เดิมได้เลย แต่มั่นใจว่า route.ts เดิมถูกลบหรือเปลี่ยนชื่อแล้ว
  const apiUrl = `${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/holidays?year=${year}`;

  try {
    const response = await fetch(apiUrl, {
      method: 'GET',
      headers: { 
        'Content-Type': 'application/json',
        // 'Authorization': `Bearer ${Token}` // ถ้า Backend ต้องการ Token อย่าลืมใส่
      },
    });

    if (!response.ok) {
      console.error(`Error status: ${response.status}`);
      const errorBody = await response.text();
      console.error(`Error body: ${errorBody}`);
      throw new Error('Network response was not ok');
    }

    return await response.json();
  } catch (error) {
    console.error("Failed to fetch holidays:", error);
    return [];
  }
};