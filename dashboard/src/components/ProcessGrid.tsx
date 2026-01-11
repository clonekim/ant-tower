import { Button, Tag } from "antd";
import { PoweroffOutlined } from "@ant-design/icons";
import { AgGridReact } from "ag-grid-react";
import { themeQuartz, type ColDef } from "ag-grid-community";
import { format } from "date-fns";
import type { Process } from "../types";

interface ProcessGridProps {
  processes: Process[];
  onKill: (pid: number) => void;
}

export const ProcessGrid: React.FC<ProcessGridProps> = ({
  processes,
  onKill,
}) => {
  const columns: ColDef<Process>[] = [
    {
      field: "pid",
      headerName: "PID",
      width: 100,
      sortable: true,
      cellRenderer: (params) => <Tag color="blue">{params.value}</Tag>,
    },
    {
      field: "name",
      headerName: "Name",
      flex: 1,
      filter: true,
      minWidth: 150,
      sortable: true,
    },
    {
      field: "start_time",
      headerName: "Started At",
      width: 200,
      sortable: true,
      valueFormatter: (params) => {
        if (params.value) {
          return format(new Date(params.value), "yyyy-MM-dd HH:mm:ss");
        }
        return "";
      },
    },
    {
      headerName: "Action",
      width: 120,
      sortable: false,
      cellRenderer: (params) => (
        <Button
          danger
          size="small"
          icon={<PoweroffOutlined />}
          onClick={() => onKill(params.data.pid)}
        >
          Kill
        </Button>
      ),
    },
  ];

  return (
    <div style={{ height: "60vh", width: "100%" }}>
      <AgGridReact
        theme={themeQuartz}
        rowData={processes}
        columnDefs={columns}
        pagination={true}
        paginationPageSize={15}
        getRowId={(params) => String(params.data.pid)}
      />
    </div>
  );
};
