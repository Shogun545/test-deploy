"use client";

import {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
} from "react";
import { useRouter } from "next/navigation";
import { loginApi } from "@/src/services/http/auth";
import type { User, LoginInterface } from "@/src/interfaces/user";

type AuthContextType = {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (credentials: LoginInterface) => Promise<User>;
  logout: () => void;
  setUser: (user: User | null) => void;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const router = useRouter();

  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  // 🔥 สำคัญที่สุด — ป้องกัน UI render ก่อนโหลด token
  const [isLoading, setIsLoading] = useState(true);

  // โหลด user + token จาก localStorage
  useEffect(() => {
    try {
      const savedUser = localStorage.getItem("user");
      const savedToken = localStorage.getItem("token");

      if (savedUser && savedToken && savedUser !== "undefined" && savedToken !== "undefined") {
        const parsedUser = JSON.parse(savedUser);

        setUser(parsedUser);
        setToken(savedToken);
        setIsAuthenticated(true);
      } else {
        localStorage.removeItem("user");
        localStorage.removeItem("token");
      }
    } catch (e) {
      console.warn("Invalid auth data, clearing...", e);
      localStorage.removeItem("user");
      localStorage.removeItem("token");
    } finally {
      setIsLoading(false);
    }
  }, []);

  // --------------------------
  // LOGIN
  // --------------------------
  const login = async (credentials: LoginInterface): Promise<User> => {
    const data = await loginApi(credentials);

    if (!data.token) throw new Error("Server ไม่ได้ส่ง token มา");

    localStorage.setItem("token", data.token);
    localStorage.setItem("user", JSON.stringify(data.user));

    setUser(data.user);
    setToken(data.token);
    setIsAuthenticated(true);

    return data.user;
  };

  // --------------------------
  // LOGOUT
  // --------------------------
  const logout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("user");

    setUser(null);
    setToken(null);
    setIsAuthenticated(false);

    router.replace("/login");
  };

  // -------------------------------------------
  // ⛔ ป้องกันการ render UI ก่อนโหลด user เสร็จ
  // -------------------------------------------
  if (isLoading) {
    return (
      <div className="w-full h-screen flex items-center justify-center text-gray-500">
        Loading...
      </div>
    );
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        isAuthenticated,
        isLoading,
        login,
        logout,
        setUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth ต้องใช้ภายใน AuthProvider เท่านั้น");
  return ctx;
};
