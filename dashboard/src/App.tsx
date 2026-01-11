import { useEffect, useState, useCallback } from "react";
import { Layout, ConfigProvider, theme, message } from "antd";
import { Header, ProcessGrid, Footer, SystemInfo } from "./components";
import { useWebSocket } from "./hooks/useWebSocket";
import { fetchProcesses, killProcess } from "./services/api";

const { Content } = Layout;

function App() {
  const [isDarkMode, setIsDarkMode] = useState(() => {
    const savedTheme = localStorage.getItem("theme");
    return savedTheme ? savedTheme === "dark" : true;
  });

  const toggleTheme = (checked: boolean) => {
    setIsDarkMode(checked);
    localStorage.setItem("theme", checked ? "dark" : "light");
  };

  const [loading, setLoading] = useState(false);
  const { wsStatus, processes, setProcesses, uptime } = useWebSocket(
    (process) => {
      //message.info(`프로세스 시작: ${process.name}`);
    },
    (process) => {
      // message.warning(`프로세스 종료: ${process.name}`);
    }
  );

  // 초기 스냅샷 로딩
  const fetchSnapshot = useCallback(async () => {
    setLoading(true);
    try {
      const data = await fetchProcesses();
      setProcesses(data);
    } catch (error) {
      message.error("데이터 로딩 실패");
    } finally {
      setLoading(false);
    }
  }, [setProcesses]);

  // 프로세스 종료 요청
  const handleKillProcess = async (pid: number) => {
    try {
      await killProcess(pid);
      message.success(`PID ${pid} 종료 명령 전송됨`);
    } catch (error) {
      message.error(error.message);
    }
  };

  // 초기 로딩
  useEffect(() => {
    fetchSnapshot();
  }, [fetchSnapshot]);

  return (
    <ConfigProvider
      theme={{
        algorithm: isDarkMode ? theme.darkAlgorithm : theme.defaultAlgorithm,
        token: {
          colorPrimary: "#1890ff",
        },
      }}
    >
      <Layout
        data-ag-theme-mode={isDarkMode ? "dark" : "light"}
        style={{ minHeight: "100vh", display: "flex", flexDirection: "column" }}
      >
        <Header
          isDarkMode={isDarkMode}
          onToggleTheme={toggleTheme}
          wsStatus={wsStatus}
          onRefresh={fetchSnapshot}
          loading={loading}
        />

        <Content
          style={{
            padding: "24px",
            flex: 1,
            display: "flex",
            flexDirection: "column",
          }}
        >
          <SystemInfo
            processCount={processes?.length}
            uptime={uptime}
            loading={loading}
          />
          <ProcessGrid processes={processes} onKill={handleKillProcess} />
        </Content>

        <Footer />
      </Layout>
    </ConfigProvider>
  );
}

export default App;
