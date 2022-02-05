import { useState, useEffect } from "react";
import { DataGrid } from "@mui/x-data-grid";
import { useSnackbar } from "notistack";
import http from "../services/axios-common";
import {
  CustomToolbar,
  CustomNoRowsOverlay,
  CustomFooterStatusComponent,
} from "../components/DatagridComponents";
import { Box, Button } from "@mui/material";

type Props = {
  accesskey: string;
  filename: string;
};

export default function FixedSizeGrid(props: Props) {
  const data = {
    columns: [
      {
        field: "id",
        hide: true,
      },
      {
        field: "key",
        headerName: "KEY",
        flex: 1,
      },
      {
        field: "value",
        headerName: "VALUE",
        flex: 3,
        editable: true,
      },
    ],
    rows: [],
    initialState: {
      columns: {
        columnVisibilityModel: {
          id: false,
        },
      },
    },
  };

  const [dataGrid, setDataGrid] = useState(data);
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();
  const [status, setStatus] = useState("connected");

  useEffect(() => {
    http
      .get(
        "/api/v1/db/?accesskey=" +
          props.accesskey +
          "&filename=" +
          props.filename +
          "&dbtype=buntdb"
      )
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
      <Box sx={{ height: 350, width: 1, mb: 2 }}>
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
            footer: { status: status },
          }}
          pageSize={15}
          {...dataGrid}
        />
      </Box>
      <Button
        color="primary"
        variant="contained"
        onClick={() =>
          setStatus((current) =>
            current === "connected" ? "disconnected" : "connected"
          )
        }
      >
        {status === "connected" ? "Disconnect" : "Connect"}
      </Button>
    </Box>
  );
}
