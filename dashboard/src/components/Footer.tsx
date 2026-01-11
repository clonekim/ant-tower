import { Layout, Button, Space, Popconfirm, message } from "antd";
import { controlPower } from "../services/api";

const { Footer: AntFooter } = Layout;

export const Footer: React.FC = () => {
  const handlePowerControl = async (
    action: "shutdown" | "reboot" | "logoff"
  ) => {
    try {
      await controlPower(action);
      message.success(
        `${
          action === "shutdown"
            ? "종료"
            : action === "reboot"
            ? "재부팅"
            : "로그오프"
        } 명령 전송됨`
      );
    } catch (error) {
      message.error((error as Error).message);
    }
  };

  return (
    <AntFooter
      style={{
        display: "flex",
        justifyContent: "flex-end",
        alignItems: "center",
        position: "fixed",
        bottom: "15px",
        left: 0,
        right: 0,
        backgroundColor: "inherit",
        borderTop: "1px solid #d9d9d9",
        zIndex: 100,
        padding: "16px 24px",
      }}
    >
      <Space>
        <Popconfirm
          title="시스템 로그오프"
          description="정말 로그오프하시겠습니까?"
          onConfirm={() => handlePowerControl("logoff")}
          okText="예"
          cancelText="아니오"
        >
          <Button type="default">Logoff</Button>
        </Popconfirm>
        <Popconfirm
          title="시스템 재부팅"
          description="정말 재부팅하시겠습니까?"
          onConfirm={() => handlePowerControl("reboot")}
          okText="예"
          cancelText="아니오"
        >
          <Button type="default">Reboot</Button>
        </Popconfirm>
        <Popconfirm
          title="시스템 종료"
          description="정말 종료하시겠습니까?"
          onConfirm={() => handlePowerControl("shutdown")}
          okText="예"
          cancelText="아니오"
        >
          <Button type="primary" danger>
            Shutdown
          </Button>
        </Popconfirm>
      </Space>
    </AntFooter>
  );
};
