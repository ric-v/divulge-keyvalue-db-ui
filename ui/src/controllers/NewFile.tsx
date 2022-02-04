import { Box, Button, Grid } from "@mui/material";
import ProgressBar from "../components/ProgressBar";
import React, { ReactChildren } from "react";

const NewFile = () => {
  return (
    <Box>
      <form>
        <h1>Select your database file</h1>
        <Button
          variant="contained"
          color="secondary"
          type="submit"
          // disabled={!file || uploading}
        >
          Upload
        </Button>
      </form>
    </Box>
  );
};

export default NewFile;
