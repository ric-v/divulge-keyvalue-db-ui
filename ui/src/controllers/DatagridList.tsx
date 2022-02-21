import { useState, useEffect } from "react";
import { DataGrid, GridSelectionModel } from "@mui/x-data-grid";
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
  status: string;
  dbtype: string;
  bucket: string;
};

export default function FixedSizeGrid({
  dbkey,
  status,
  dbtype,
  bucket,
}: Props) {
  const data = {
    columns: [],
    rows: [],
    initialState: {},
  };

  const [dataGrid, setDataGrid] = useState(data);
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();
  const [updated, setUpdated] = useState(false);
  const [selectionModel, setSelectionModel] = useState<GridSelectionModel>([]);
  const [keysToDelete, setKeysToDelete] = useState<any[]>([]);
  const [showDelete, setShowDelete] = useState(false);

  useEffect(() => {
    if (dbtype === "boltdb" && bucket === "") {
      enqueueSnackbar("Select a default bucket to load data from DB", {
        key: "warning",
        variant: "warning",
        draggable: true,
        onClick: () => {
          closeSnackbar("warning");
        },
      });
      return;
    }
    http(dbkey)
      .get("/api/v1/db/")
      .then((resp) => {
        console.log("setting data", resp.data);
        setDataGrid(resp.data.data);
        setUpdated(false);
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
    setSelectionModel([]);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [updated, dbtype, bucket]);

  const handleUpdate = (e: any) => {
    const newValue = e.value;
    const key: any = dataGrid.rows[+e.id - 1];
    console.log("newValue", newValue, "key", key.key);
    http(dbkey)
      .put("/api/v1/db/" + key.key, {
        value: newValue,
        key: key.key,
      })
      .then((resp) => {
        enqueueSnackbar(resp.data.message, {
          key: "updated",
          variant: "success",
          onClick: () => {
            closeSnackbar("updated");
          },
        });
        setUpdated(true);
      })
      .catch((err) => {
        enqueueSnackbar(err.response.data.message, {
          key: "error",
          variant: "error",
          onClick: () => {
            closeSnackbar("error");
          },
        });
      });
  };

  const handleSelection = (newSelectionModel: GridSelectionModel) => {
    setSelectionModel(newSelectionModel);
    // replace id with key
    const keys = newSelectionModel.map((item) => {
      const gridItem: any = dataGrid.rows[+item - 1];
      return gridItem.key;
    });
    setKeysToDelete(keys);
    console.log("keys", keys);

    if (keys.length > 0) {
      setShowDelete(true);
    } else {
      setShowDelete(false);
    }
  };

  return (
    <Box sx={{ width: 1 }}>
      <Box sx={{ height: 720, width: 1, mb: 2 }}>
        <DataGrid
          density={"compact"}
          disableColumnMenu={true}
          components={{
            Toolbar: CustomToolbar,
            NoRowsOverlay: CustomNoRowsOverlay,
            Footer: CustomFooterStatusComponent,
          }}
          componentsProps={{
            footer: {
              status: status,
              setUpdated: setUpdated,
              dbkey: dbkey,
              showDelete: showDelete,
              keys: keysToDelete,
            },
          }}
          pageSize={15}
          disableSelectionOnClick
          checkboxSelection
          onSelectionModelChange={(newSelectionModel) =>
            handleSelection(newSelectionModel)
          }
          selectionModel={selectionModel}
          onCellEditCommit={(e) => handleUpdate(e)}
          {...dataGrid}
        />
      </Box>
    </Box>
  );
}
