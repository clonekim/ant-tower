import { useEffect, useState } from "react";
import { Card, Statistic, Space, message } from "antd";
import {
  UserOutlined,
  DesktopOutlined,
  ClockCircleOutlined,
} from "@ant-design/icons";
import { fetchCurrentUser } from "../services/api";
import { formatUptime } from "../utils/time";

interface SystemInfoProps {
  processCount: number;
  uptime: number;
  loading: boolean;
}

export const SystemInfo: React.FC<SystemInfoProps> = ({
  processCount,
  uptime,
  loading,
}) => {
  const [username, setUsername] = useState<string>("");

  // CurrentUser: 초기 1번만 로드
  useEffect(() => {
    const loadUsername = async () => {
      try {
        const user = await fetchCurrentUser();
        setUsername(user);
      } catch {
        message.error("사용자 정보 로딩 실패");
      }
    };
    loadUsername();
  }, []);

  return (
    <Card style={{ marginBottom: "24px" }} loading={loading}>
      <div
        style={{
          display: "flex",
          width: "100%",
          justifyContent: "space-between",
          alignItems: "flex-start",
        }}
      >
        <Space size="large">
          <Statistic
            title="로그인 사용자"
            value={username}
            prefix={<UserOutlined />}
          />
          <div style={{ borderLeft: "1px solid #d9d9d9", height: "60px" }} />
          <Statistic
            title="프로세스"
            value={processCount}
            suffix="개"
            prefix={<DesktopOutlined />}
          />
        </Space>
        <Statistic
          title="시스템 업타임"
          value={formatUptime(uptime, true)}
          prefix={<ClockCircleOutlined />}
        />
      </div>
    </Card>
  );
};
