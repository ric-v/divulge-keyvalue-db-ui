import { useState, useEffect } from "react";
import { DataGrid } from "@mui/x-data-grid";
import { useSnackbar } from "notistack";
import http from "../services/axios-common";
import {
  CustomToolbar,
  CustomNoRowsOverlay,
  CustomFooterStatusComponent,
} from "../components/DatagridComponents";
import { Box } from "@mui/material";

type Props = {
  dbkey: string;
  dbname: string;
  status: string;
  setStatus: React.Dispatch<React.SetStateAction<string>>;
  setDbname: React.Dispatch<React.SetStateAction<string>>;
  setDbkey: React.Dispatch<React.SetStateAction<string>>;
  setLoadView: React.Dispatch<React.SetStateAction<boolean>>;
};

export default function FixedSizeGrid(props: Props) {
  const data = {
    columns: [],
    rows: [],
    initialState: {},
  };

  const [dataGrid, setDataGrid] = useState(data);
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();

  useEffect(() => {
    http
      .get("/api/v1/db/?dbkey=" + props.dbkey)
      .then((resp) => {
        enqueueSnackbar("Updating tables", {
          key: "load",
          variant: "info",
          onClick: () => {
            closeSnackbar("load");
          },
        });
        console.log("setting data", resp.data);
        setDataGrid(resp.data.data);
      })
      .catch((err) => {
        enqueueSnackbar(err.response.data.message, {
          key: "error",
          variant: "error",
          draggable: true,
          onClick: () => {
            closeSnackbar("error");
          },
        });
      });
  }, []);

  return (
    <Box sx={{ width: 1 }}>
      <Box sx={{ height: 720, width: 1, mb: 2 }}>
        {/* @ts-ignore */}
        <DataGrid
          density={"compact"}
          disableColumnMenu={true}
          components={{
            Toolbar: CustomToolbar,
            NoRowsOverlay: CustomNoRowsOverlay,
            Footer: CustomFooterStatusComponent,
          }}
          componentsProps={{
            footer: { status: props.status },
          }}
          pageSize={15}
          {...dataGrid}
        />
      </Box>
    </Box>
  );
}
