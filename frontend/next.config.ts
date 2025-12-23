import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // 1. จำเป็นมากสำหรับ Dockerfile ที่ให้ไปก่อนหน้านี้
  output: "standalone",

  // 2. แนะนำให้ลบ unoptimized: true ทิ้ง เพื่อให้เว็บโหลดรูปเร็ว
  // แต่ถ้า Server สเปคต่ำมาก (RAM < 500MB) ค่อยเปิดไว้ครับ
  images: {
    // unoptimized: true, // <-- Comment ทิ้ง หรือลบออก ถ้าอยากได้ Performance ดีๆ
    remotePatterns: [
      {
        protocol: "https",
        hostname: "**", // อนุญาตให้ดึงรูปจากภายนอกได้ (ถ้ามี)
      },
    ],
  },
};

export default nextConfig;