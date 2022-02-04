import { Box, Button, Grid, TextField, Typography } from "@mui/material";
import { useSnackbar } from "notistack";
import { useEffect, useState } from "react";
import NewFileForm from "../components/NewFileForm";
import http from "../services/axios-common";

const NewFile = () => {
  const [dbname, setDbname] = useState("");
  const [errorMessage, setErrorMessage] = useState("");
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();

  // use effect to do db name validation
  useEffect(() => {
    if (dbname.length > 0 && dbname.length < 3) {
      setErrorMessage("Database name must be at least 3 characters long");
    } else if (dbname.length > 30) {
      setErrorMessage("Database name is too long");
    } else if (dbname.match(/[^a-zA-Z0-9_]/)) {
      setErrorMessage("File name should not contain special characters");
    } else {
      setErrorMessage("");
    }
  }, [dbname, errorMessage]);

  const createNewFile = () => {
    http
      .post("/api/v1/new?file=" + dbname + ".db&dbtype=boltdb")
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
        setDbname("");
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
    <NewFileForm
      dbname={dbname}
      setDbname={setDbname}
      errorMessage={errorMessage}
      createNewFile={createNewFile}
    />
  );
};

export default NewFile;
