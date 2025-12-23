"use client";
import Image from "next/image";
import { Navbar } from "../components/Navbar/page";
import { GlassButton } from "../components/Button/GlassButton";
import Link from "next/link";
export default function HomePage() {
  return (
    <div className="min-h-screen bg-white flex flex-col">
      <Navbar />
      <div className="relative w-full flex-1">
        <Image
          src="/backgroundlogin.svg"
          alt="Background"
          fill
          className="object-cover"
          priority
        />
        {/* overlay */}
        <div className="absolute inset-0 bg-black/20" />
        <div className="absolute inset-0 flex flex-col items-center justify-center text-center px-4">
          <h1 className="text-[100px] font-medium tracking-tight text-white drop-shadow-2xl">
            <span className="bg-gradient-to-r from-[#F26522] to-[#CFD1D7] bg-clip-text text-transparent">
              ENGI
            </span>{" "}
            ADVISORY
          </h1>
          <p className="-mt-6 text-xl md:text-2xl text-[#F5F5F5] drop-shadow max-w-2xl leading-relaxed">
            Empowering students and advisors through smart connection. Guiding
            every step toward academic success.
          </p>
          <Link href="/login">
            <GlassButton className="mt-6 rounded-full px-10 py-4 text-lg font-normal text-white hover:bg-white/30 cursor-pointer">
              Get Start
            </GlassButton>
          </Link>
        </div>
      </div>
    </div>
  );
}
