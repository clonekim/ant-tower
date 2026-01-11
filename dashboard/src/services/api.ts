import axios from "axios";
import type { Process } from "../types";

export const fetchProcesses = async (): Promise<Process[]> => {
  try {
    const res = await axios.get("/api/process");
    return res.data.processes;
  } catch (error) {
    const message = (error as any).response?.data?.message || "데이터 로딩 실패";
    throw new Error(message);
  }
};

export const killProcess = async (pid: number): Promise<void> => {
  try {
    await axios.post("/api/process/kill", { pid });
  } catch (error) {
    const message = (error as any).response?.data?.message || "종료 실패";
    throw new Error(message);
  }
};

export const fetchUptime = async (): Promise<number> => {
  try {
    const res = await axios.get("/api/uptime");
    return res.data.uptime_seconds;
  } catch (error) {
    const message = (error as any).response?.data?.message || "업타임 로딩 실패";
    throw new Error(message);
  }
};

export const fetchCurrentUser = async (): Promise<string> => {
  try {
    const res = await axios.get("/api/login");
    return res.data.username;
  } catch (error) {
    const message = (error as any).response?.data?.message || "사용자 정보 로딩 실패";
    throw new Error(message);
  }
};

export const controlPower = async (action: "shutdown" | "reboot" | "logoff"): Promise<void> => {
  try {
    await axios.post("/api/power", { action });
  } catch (error) {
    const message = (error as any).response?.data?.message || "전원 제어 실패";
    throw new Error(message);
  }
};
