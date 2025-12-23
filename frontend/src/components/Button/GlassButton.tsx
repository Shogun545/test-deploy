// src/components/GlassButton.tsx
import React from "react";

type GlassButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  children: React.ReactNode;
};

export function GlassButton({ children, className = "", ...props }: GlassButtonProps) {
  return (
    <button
      className={
        `
        relative flex items-center justify-center
        px-6 py-3
        rounded-2xl
        text-base md:text-lg font-medium
        text-white
        bg-white/15 dark:bg-white/10
        backdrop-blur-xl
        border border-white/30
        shadow-[0_8px_32px_rgba(15,23,42,0.35)]
        hover:bg-white/25 hover:shadow-[0_12px_40px_rgba(15,23,42,0.45)]
        active:scale-[0.97]
        transition-all duration-200

        ` + " " + className
      }
      {...props}
    >
      {children}
    </button>
  );
}
