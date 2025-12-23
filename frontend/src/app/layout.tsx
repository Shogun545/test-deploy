// src/app/layout.tsx
import "./globals.css";
import type { ReactNode } from "react";
import type { Metadata } from "next";
import { Kanit } from "next/font/google";
import { AuthProvider } from "../contexts/AuthContext";

const fontKanit = Kanit({
  subsets: ["thai", "latin"],
  weight: ["100", "200", "300", "400", "500", "600", "700", "800", "900"],
  variable: "--font-kanit",
});

export const metadata: Metadata = {
  title: "ENGI Advisory",
  description: "SUT Advisory Appointment System",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="th">
      <body className={`${fontKanit.variable} antialiased h-screen w-screen`}>
        <AuthProvider>
          {children}
        </AuthProvider>
      </body>
    </html>
  );
}
