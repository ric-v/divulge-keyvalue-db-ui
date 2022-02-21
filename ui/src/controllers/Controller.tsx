import { Box, Container, Grid, Typography, Button } from "@mui/material";
import { Dispatch, SetStateAction, useEffect, useState } from "react";
import VerticalLinearStepper from "../components/DBSelectStepper";
import FlexLayoutGrid from "./DatagridList";
import FiberManualRecordIcon from "@mui/icons-material/FiberManualRecord";
import CloseIcon from "@material-ui/icons/Close";
import GetAppIcon from "@material-ui/icons/GetApp";
import http from "../services/axios-common";
import { useSnackbar } from "notistack";
import ManageBucketsModal from "../components/ManageBucketsModal";

const Controller = () => {
  const [status, setStatus] = useState("connected");
  const [dbname, setDbname] = useState("");
  const [dbkey, setDbkey] = useState("");
  const [dbtype, setDbtype] = useState("");
  const [loadView, setLoadView] = useState(false);
  const [bucket, setBucket] = useState("");
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();

  useEffect(() => {
    // check if db access key is available in local storage
    if (localStorage.getItem("dbkey") !== null && dbkey === "") {
      const dbkey = localStorage.getItem("dbkey");
      console.log("laoding dbkey:", dbkey);
      if (dbkey) {
        setDbkey(dbkey);
        setStatus("connected");

        http(dbkey)
          .post("/api/v1/load")
          .then((resp) => {
            enqueueSnackbar("File loaded successfully", {
              key: "success",
              variant: "success",
              draggable: true,
              onClick: () => {
                closeSnackbar("success");
              },
            });
            console.log(resp.data);
            setDbname(resp.data.filename);
            setDbkey(resp.data.dbkey);
            setDbtype(resp.data.dbtype);
            setStatus("connected");
            setLoadView(true);
            localStorage.setItem("dbkey", resp.data.dbkey);
          })
          .catch((err) => {
            if (err.response.status !== 400) console.log(err);
          });
      } else {
        setStatus("disconnected");
        setLoadView(false);
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [dbkey, bucket, loadView]);

  const downloadFile = (dbkey: string, dbname: string) => {
    http(dbkey)
      .get("/api/v1/download")
      .then((resp) => {
        const url = window.URL.createObjectURL(
          new Blob([resp.data], { type: "application/file" })
        );
        const link = document.createElement("a");
        link.href = url;
        link.setAttribute("download", dbname); //or any other extension
        document.body.appendChild(link);
        link.click();
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

  const closedbConnection = (
    dbkey: string,
    setDbname: Dispatch<SetStateAction<string>>,
    setDbkey: Dispatch<SetStateAction<string>>
  ) => {
    http(dbkey)
      .delete("/api/v1/clear")
      .then((_resp) => {
        enqueueSnackbar("Database connection closed", {
          key: "success",
          variant: "success",
          draggable: true,
          onClick: () => {
            closeSnackbar("success");
          },
        });
        setStatus("disconnected");
        setDbname("");
        setDbkey("");
        setLoadView(false);
        localStorage.removeItem("dbkey");
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
    <Container
      sx={{
        marginTop: "1rem",
        flexGrow: 1,
      }}
    >
      {loadView ? (
        // Datagrid view after data is loaded
        <Grid container columns={{ xs: 2, sm: 12, md: 12 }} spacing={3}>
          <Grid item flex={1}>
            <Grid
              sx={{
                display: "flex",
                flexDirection: "row",
                justifyContent: "space-between",
                alignItems: "center",
              }}
            >
              <Box
                sx={{
                  display: "flex",
                  flexDirection: "row",
                  alignItems: "center",
                }}
              >
                <Typography variant="h5" mr={"12px"}>
                  Database:{" "}
                </Typography>
                <Typography variant="h3">{dbname}</Typography>
                {dbtype === "boltdb" ? (
                  <Typography variant="caption" ml={"12px"}>
                    {bucket}
                  </Typography>
                ) : (
                  ""
                )}

                {dbtype === "boltdb" ? (
                  <ManageBucketsModal
                    defBucket={bucket}
                    setBucket={setBucket}
                  ></ManageBucketsModal>
                ) : (
                  ""
                )}
              </Box>
              <Box
                sx={{
                  display: "flex",
                  flexDirection: "row",
                  alignItems: "center",
                }}
              >
                Connection Status:
                <FiberManualRecordIcon
                  fontSize="small"
                  sx={{
                    mx: 2,
                    color: status === "connected" ? "#4caf50" : "#d9182e",
                  }}
                />
                <Button
                  aria-label="download"
                  color="primary"
                  onClick={() => downloadFile(dbkey, dbname)}
                >
                  <GetAppIcon />
                  download
                </Button>
                <Button
                  aria-label="close"
                  color="error"
                  onClick={() => closedbConnection(dbkey, setDbname, setDbkey)}
                >
                  <CloseIcon />
                  Close DB
                </Button>
              </Box>
            </Grid>
            <FlexLayoutGrid
              dbkey={dbkey}
              status={status}
              dbtype={dbtype}
              bucket={bucket}
            />
          </Grid>
        </Grid>
      ) : (
        // Stepper view for s
        <Grid container columns={{ xs: 2, sm: 12, md: 12 }}>
          <VerticalLinearStepper
            dbKey={dbkey}
            setDbkey={setDbkey}
            dbName={dbname}
            setDbname={setDbname}
            status={status}
            setStatus={setStatus}
            dbtype={dbtype}
            setDbtype={setDbtype}
            loadView={loadView}
            setLoadView={setLoadView}
          />
        </Grid>
      )}
    </Container>
  );
};

export default Controller;
