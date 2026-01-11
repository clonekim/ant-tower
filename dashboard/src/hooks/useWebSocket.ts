import { useEffect, useRef, useState } from "react";
import type { Process, WebSocketMessage, WSStatus } from "../types";

interface UseWebSocketReturn {
  wsStatus: WSStatus;
  processes: Process[];
  setProcesses: React.Dispatch<React.SetStateAction<Process[]>>;
  uptime: number;
  setUptime: React.Dispatch<React.SetStateAction<number>>;
}

export const useWebSocket = (
  onProcessStart?: (process: Process) => void,
  onProcessEnd?: (process: Process) => void
): UseWebSocketReturn => {
  const [processes, setProcesses] = useState<Process[]>([]);
  const [wsStatus, setWsStatus] = useState<WSStatus>("disconnected");
  const [uptime, setUptime] = useState<number>(0);
  const ws = useRef<WebSocket | null>(null);

  useEffect(() => {
    const protocol = window.location.protocol === "https:" ? "wss" : "ws";
    const socket = new WebSocket(`${protocol}://${window.location.host}/ws`);

    socket.onopen = () => setWsStatus("connected");
    socket.onclose = () => setWsStatus("disconnected");

    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data) as WebSocketMessage;
      const data = msg.data;

      // Update uptime if provided
      if (msg.uptime !== undefined) {
        setUptime(msg.uptime);
      }

      if (msg.type === "START") {
        setProcesses((prev) => [...prev, data]);
        onProcessStart?.(data);
      } else if (msg.type === "END") {
        setProcesses((prev) => prev.filter((p) => p.pid !== data.pid));
        onProcessEnd?.(data);
      }
    };

    ws.current = socket;
    return () => socket.close();
  }, []);

  return { wsStatus, processes, setProcesses, uptime, setUptime };
};
