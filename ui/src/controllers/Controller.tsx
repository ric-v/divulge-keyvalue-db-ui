import { Container, Grid, Typography } from "@mui/material";
import { useEffect, useState } from "react";
import VerticalLinearStepper from "../components/DBSelectStepper";
import FlexLayoutGrid from "./DatagridList";

const Controller = () => {
  const [status, setStatus] = useState("connected");
  const [dbname, setDbname] = useState("");
  const [dbkey, setDbkey] = useState("");
  const [dbtype, setDbtype] = useState("");
  const [loadView, setLoadView] = useState(false);

  useEffect(() => {}, [dbkey, dbname, status]);
  console.log("loadView:", loadView);

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
                justifyContent: "flex-start",
                alignItems: "center",
              }}
            >
              <Typography variant="h5" mr={"12px"}>
                Database:{" "}
              </Typography>
              <Typography variant="h3">{dbname}</Typography>
            </Grid>
            <FlexLayoutGrid
              dbkey={dbkey}
              dbname={dbname}
              status={status}
              setStatus={setStatus}
              setDbkey={setDbkey}
              setDbname={setDbname}
              setLoadView={setLoadView}
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
