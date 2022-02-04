import { Box, Button, Grid, TextField, Typography } from "@mui/material";
import { useSnackbar } from "notistack";
import { useEffect, useState } from "react";
import http from "../services/axios-common";

type Props = {
  dbname: string;
  setDbname: React.Dispatch<React.SetStateAction<string>>;
  errorMessage: string;
  createNewFile: () => void;
};

const NewFileForm = (props: Props) => {
  return (
    <>
      <Typography variant="h3">Create new database file</Typography>
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
                props.dbname.length > 30 ||
                (props.dbname.length > 0 && props.dbname.length < 3)
              }
              helperText={props.errorMessage}
              color="primary"
              size="small"
              label="Database Name"
              value={props.dbname}
              onChange={(e) => {
                e.preventDefault();
                props.setDbname(e.target.value);
              }}
              fullWidth
            />
          </Grid>
          <Grid item xs={1} sm={2} md={2} lg={4}>
            <Button
              variant="contained"
              color="secondary"
              fullWidth
              disabled={props.dbname === ""}
              onClick={() => {
                props.createNewFile();
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

export default NewFileForm;
