import { useState, useEffect, Dispatch, SetStateAction } from "react";
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

  useEffect(() => {
    http
      .get(
        "/api/v1/db/?accesskey=" +
          props.dbkey +
          "&filename=" +
          props.dbname +
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

  const closedbConnection = (
    dbkey: string,
    setDbname: Dispatch<SetStateAction<string>>,
    setDbkey: Dispatch<SetStateAction<string>>
  ) => {
    http
      .delete("/api/v1/clear?accesskey=" + dbkey)
      .then((resp) => {
        enqueueSnackbar("Database connection closed", {
          key: "success",
          variant: "success",
          draggable: true,
          onClick: () => {
            closeSnackbar("success");
          },
        });
        props.setStatus("disconnected");
        setDbname("");
        setDbkey("");
        props.setLoadView(false);
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
  };

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
      <Button
        color="primary"
        variant="contained"
        onClick={() =>
          props.setStatus((current: string) => {
            closedbConnection(props.dbkey, props.setDbname, props.setDbkey);
            return current === "connected" ? "disconnected" : "connected";
          })
        }
      >
        {props.status === "connected" ? "Disconnect" : "Connect"}
      </Button>
    </Box>
  );
}
