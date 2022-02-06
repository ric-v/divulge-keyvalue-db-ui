import { useSnackbar } from "notistack";
import { useEffect, useState } from "react";
import http from "../services/axios-common";
import { Box, Button, Grid, TextField } from "@mui/material";

type FileUploadProps = {
  dbName: string;
  setDbname: React.Dispatch<React.SetStateAction<string>>;
  setDbkey: React.Dispatch<React.SetStateAction<string>>;
  setStatus: React.Dispatch<React.SetStateAction<string>>;
};

const NewFile = ({
  dbName,
  setDbname,
  setDbkey,
  setStatus,
}: FileUploadProps) => {
  const [inputDbName, setInputDbName] = useState("");
  const [errorMessage, setErrorMessage] = useState("");
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();

  // use effect to do db name validation
  useEffect(() => {
    if (inputDbName.length > 0 && inputDbName.length < 3) {
      setErrorMessage("Database name must be at least 3 characters long");
    } else if (inputDbName.length > 30) {
      setErrorMessage("Database name is too long");
    } else if (inputDbName.match(/[^a-zA-Z0-9_]/)) {
      setErrorMessage("File name should not contain special characters");
    } else {
      setErrorMessage("");
    }
  }, [inputDbName, errorMessage]);

  const createNewFile = () => {
    http
      .post("/api/v1/new?file=" + dbName + ".db&dbtype=buntdb")
      .then((resp) => {
        enqueueSnackbar("File created successfully", {
          key: "success",
          variant: "success",
          draggable: true,
          onClick: () => {
            closeSnackbar("success");
          },
        });
        console.log(resp);
        setInputDbName("");
        setDbname(resp.data.filename);
        setDbkey(resp.data.dbkey);
        setStatus("connected");
        localStorage.setItem("dbkey", resp.data.dbkey);
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
    <>
      <Box
        component="form"
        sx={{
          "& > :not(style)": { m: 1 },
          display: "flex",
          flexDirection: "column",
        }}
        alignItems={["center"]}
      >
        <Grid
          container
          columns={{ xs: 3, sm: 6, md: 8, lg: 12 }}
          rowSpacing={["1rem"]}
          spacing={["0rem", "1rem", "0rem", "1rem"]}
        >
          <Grid item xs={2} sm={4} md={6} lg={8}>
            <TextField
              error={
                dbName.length > 30 || (dbName.length > 0 && dbName.length < 3)
              }
              helperText={errorMessage}
              color="primary"
              size="small"
              label="Database Name"
              value={dbName}
              onChange={(e) => {
                e.preventDefault();
                setDbname(e.target.value);
              }}
              fullWidth
            />
          </Grid>
          <Grid item xs={1} sm={2} md={2} lg={4}>
            <Button
              variant="contained"
              color="secondary"
              fullWidth
              disabled={dbName === ""}
              onClick={() => {
                createNewFile();
              }}
            >
              Create
            </Button>
          </Grid>
        </Grid>
      </Box>
    </>
  );
};

export default NewFile;
