"use client"

import Image from "next/image"
import Link from "next/link"

export function Navbar() {
  return (
    <header className="flex items-center justify-between px-10 py-4">
      <div className="flex items-center gap-3">
        <Image
          src="/LOGOSUT.svg"
          alt="SUT Logo"
          width={150}
          height={150}
          className="object-contain"
        />
      </div>

      {/* MENU */}
      <nav className="flex items-center gap-6 text-base text-[#6D6E71] font-normal ">
        <Link href="/public/login?role=teacher" className="hover:text-[#800020]">
          Teacher
        </Link>
        <Link href="/public/login?role=student" className="hover:text-[#800020]">
          Student
        </Link>
      </nav>
    </header>
  )
}
