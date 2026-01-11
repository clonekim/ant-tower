import { Button, Space, Switch, Tag, Typography } from "antd";
import {
  ReloadOutlined,
  WindowsOutlined,
  MoonOutlined,
  SunOutlined,
} from "@ant-design/icons";
import type { WSStatus } from "../types";

const { Title } = Typography;

interface HeaderProps {
  isDarkMode: boolean;
  onToggleTheme: (checked: boolean) => void;
  wsStatus: WSStatus;
  onRefresh: () => void;
  loading: boolean;
}

export const Header: React.FC<HeaderProps> = ({
  isDarkMode,
  onToggleTheme,
  wsStatus,
  onRefresh,
  loading,
}) => {
  return (
    <div
      style={{
        display: "flex",
        alignItems: "center",
        justifyContent: "space-between",
        padding: "16px 20px",
        background: isDarkMode ? undefined : "#ffffff",
        borderBottom: isDarkMode ? undefined : "1px solid #f0f0f0",
      }}
    >
      <div style={{ display: "flex", alignItems: "center", gap: "10px" }}>
        <WindowsOutlined style={{ fontSize: "24px", color: "#1890ff" }} />
        <Title level={4} style={{ margin: 0 }}>
          TWN Monitor
        </Title>
      </div>
      <Space>
        <Switch
          checkedChildren={<MoonOutlined />}
          unCheckedChildren={<SunOutlined />}
          checked={isDarkMode}
          onChange={onToggleTheme}
        />

        <Tag color={wsStatus === "connected" ? "success" : "error"}>
          WS: {wsStatus.toUpperCase()}
        </Tag>

        <Button icon={<ReloadOutlined />} onClick={onRefresh} loading={loading}>
          Refresh
        </Button>
      </Space>
    </div>
  );
};
