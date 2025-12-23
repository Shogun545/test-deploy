import axios from "axios";
import { API_BASE_URL } from "@/src/config/api";

//const BASE_URL = "http://localhost:8080"
export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

// ใส่ token ให้ทุก request อัตโนมัติ
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});
