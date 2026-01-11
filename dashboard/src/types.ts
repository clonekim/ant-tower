// 데이터 타입 정의 (Go 백엔드와 일치)
export interface Process {
  pid: number;
  name: string;
  start_time: string; // ISO string
  end_time: string | null;
  duration?: string;
}

export type WebSocketMessageType = "START" | "END";

export interface WebSocketMessage {
  type: WebSocketMessageType;
  data: Process;
  uptime?: number; // 시스템 업타임 (초)
}

export type WSStatus = "connected" | "disconnected";
