"use client";

import Image from "next/image";
import { useState, type FormEvent } from "react";
import { Navbar } from "../../components/Navbar/page";
import { GlassButton } from "../../components/Button/GlassButton";
import { User, Lock } from "lucide-react";
import { useAuth } from "@/src/contexts/AuthContext";
import { message } from "antd";
import { useRouter } from "next/navigation";

export default function LoginPage() {
  const { login } = useAuth();
  const router = useRouter();
  const [sutId, setSutId] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);

  const [messageApi, contextHolder] = message.useMessage();

  async function handleSubmit(e?: FormEvent<HTMLFormElement>) {
    if (e) e.preventDefault();
    if (!sutId.trim() || !password) {
      messageApi.error("กรุณากรอกข้อมูลให้ครบทุกช่อง");
      return;
    }
    const sutIdRegex = /^[TBA][0-9]{7}$/;
    if (!sutIdRegex.test(sutId)) {
      messageApi.error("usernameต้องมี 8 ตัว และขึ้นต้นด้วย T, B หรือ A");
      return;
    }

    const passwordRegex = /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,16}$/;
    

    setLoading(true);

    try {
      const loggedInUser = await login({ sut_id: sutId, password });
      const role = loggedInUser?.role || "Student";

      if (role === "Admin") router.push("/admin/dashboard");
      else if (role === "Advisor") router.push("/advisor/dashboard");
      else router.push("/student/dashboard");
    } catch (err: any) {
      if (err?.response) {
        const status = err.response.status;

        if (status === 401) {
          messageApi.error("รหัสผู้ใช้หรือรหัสผ่านไม่ถูกต้อง");
        } else {
          messageApi.error("เข้าสู่ระบบไม่สำเร็จ");
        }
      } else {
        messageApi.error("ไม่สามารถเชื่อมต่อเซิร์ฟเวอร์ได้");
      }
    } finally {
      setLoading(false);
    }
  }

  return (
    <>
      {contextHolder}

      {/* เพิ่ม overflow-hidden เพื่อป้องกัน Scrollbar หลัก */}
      <div className="min-h-screen bg-white flex flex-col overflow-hidden">
        <Navbar />

        <div className="relative w-full flex-1">
          <Image
            src="/backgroundlogin.svg"
            alt="Background"
            fill
            className="object-cover"
            priority
          />
          <div className="absolute inset-0 bg-black/20" />

          {/* จัดกึ่งกลางด้วย flex-col items-center justify-center */}
          {/* เพิ่ม py-4/py-8 เพื่อให้มี padding ด้านบน/ล่าง ป้องกันการติดขอบ */}
          <div className="absolute inset-0 flex flex-col items-center justify-center px-4 py-8">
            {/* HEADER: อนุญาตให้ถูกบีบอัด (flex-shrink: 1 โดย default) */}
            <div className="text-center mb-6 flex-shrink">
              <h1 className="text-[64px] md:text-[88px] font-semibold tracking-tight text-white drop-shadow-2xl">
                <span className="bg-gradient-to-r from-[#F26522] to-[#CFD1D7] bg-clip-text text-transparent">
                  ENGI
                </span>{" "}
                ADVISORY
              </h1>
              <p className="-mt-4 text-lg md:text-xl text-[#F5F5F5] drop-shadow max-w-2xl leading-relaxed mx-auto">
                Empowering students and advisors through smart connection.
                Guiding every step toward academic success.
              </p>
            </div>

            {/* LOGIN CARD: ใช้ flex-shrink-0 เพื่อคงความสูงไว้ (ห้ามบีบอัด) */}
            <div
              className="
                glass-surface glass-surface--fallback
                w-[450px]
                rounded-[40px]
                p-10 md:p-12
                bg-white/40
                backdrop-blur-2xl
                border border-white/40
                shadow-[0_18px_45px_rgba(0,0,0,0.25)]
                flex-shrink-0
              "
            >
              <form
                noValidate
                className="glass-surface__content flex flex-col gap-6"
                onSubmit={handleSubmit}
              >
                <h2 className="text-4xl md:text-4xl font-medium text-white text-center">
                  USER LOGIN
                </h2>

                {/* INPUTS */}
                <div className="space-y-4 mt-10">
                  {/* USERNAME */}
                  <div
                    className="
                      glass-surface glass-surface--fallback
                      flex items-center gap-3
                      px-5 py-3
                      rounded-full
                      bg-white/30
                      backdrop-blur-2xl
                      border-[0.3px] border-white/40
                      shadow-[0_10px_25px_rgba(0,0,0,0.15)]
                      focus-within:border-[#F26522]
                      focus-within:ring-2
                      focus-within:ring-[#F26522]/40
                    "
                  >
                    <User className="text-gray-400 text-lg" />
                    <input
                      value={sutId}
                      onChange={(e) => setSutId(e.target.value)}
                      name="sut_id"
                      type="text"
                      placeholder="Username"
                      data-testid="sut-id-input"
                      className="w-full bg-transparent outline-none text-gray-700 placeholder-gray-400"
                    />
                  </div>

                  {/* PASSWORD */}
                  <div
                    className="
                      glass-surface glass-surface--fallback
                      flex items-center gap-3
                      px-5 py-3
                      rounded-full
                      bg-white/30
                      backdrop-blur-2xl
                      border-[0.3px] border-white/40
                      shadow-[0_10px_25px_rgba(0,0,0,0.15)]
                      focus-within:border-[#F26522]
                      focus-within:ring-2
                      focus-within:ring-[#F26522]/40
                    "
                  >
                    <Lock className="text-gray-400 text-lg" />
                    <input
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      name="password"
                      type="password"
                      placeholder="Password"
                      data-testid="password-input"
                      className="w-full bg-transparent outline-none text-gray-700 placeholder-gray-400"
                    />
                  </div>
                  {/* REMEMBER & FORGOT */}
                  <div className="flex items-center justify-end px-2 mt-2 text-xs md:text-sm text-gray-600 w-full">
                    <button
                      type="button"
                      className="text-gray-500 hover:text-[#F26522] hover:underline underline-offset-2"
                    >
                      Forgot password?
                    </button>
                  </div>
                </div>

                {/* LOGIN BUTTON */}
                <GlassButton
                  type="submit"
                  data-testid="login-btn"
                  className="mt-1 w-[150px] mx-auto rounded-full px-6 py-3 text-base font-medium !text-white hover:bg-white/30 cursor-pointer"
                  disabled={loading}
                >
                  {loading ? "Loading..." : "LOGIN"}
                </GlassButton>
              </form>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
