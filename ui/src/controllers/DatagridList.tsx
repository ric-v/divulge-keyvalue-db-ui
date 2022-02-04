import * as React from "react";
import { DataGrid } from "@mui/x-data-grid";
import { useSnackbar } from "notistack";
import http from "../services/axios-common";

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
        width: 110,
      },
      {
        field: "value",
        headerName: "VALUE",
        width: 180,
        editable: true,
      },
    ],
    rows: [
      {
        id: "93b20245-b180-503f-b45e-317b3e2a0ee1",
        key: "D-8966",
        value: "Sugar No.11",
      },
      {
        id: "d6ef22ed-f6f8-53a2-8341-975df7e29279",
        key: "D-3840",
        value: "Rapeseed",
      },
      {
        id: "98eb7594-1010-55b2-9dc3-9af24529a856",
        key: "D-1614",
        value: "Wheat",
      },
    ],
    initialState: {
      columns: {
        columnVisibilityModel: {
          id: false,
        },
      },
    },
  };

  const [dataGrid, setDataGrid] = React.useState(data);
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();

  const loadData = () => {
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
        console.log('setting data',resp.data);
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
  };
  loadData();
  console.log("-->", data);

  return (
    <div style={{ height: 400, width: "100%" }}>
      <div style={{ height: 350, width: "100%" }}>
        {/* @ts-ignore */}
        <DataGrid {...dataGrid} />
      </div>
    </div>
  );
}
