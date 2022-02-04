import { Container, Grid, Typography } from "@mui/material";
import FlexLayoutGrid from "./DatagridList";
import FileUpload from "./FileUpload";
import NewFile from "./NewFile";

const Controller = () => {
  const dbname = localStorage.getItem("dbname");
  const accesskey = localStorage.getItem("accesskey");

  return (
    <Container
      sx={{
        marginTop: "1rem",
        flexGrow: 1,
      }}
    >
      {dbname && accesskey ? (
        <Grid container columns={{ xs: 2, sm: 12, md: 12 }} spacing={3}>
          <Grid item flex={1}>
            <Typography variant="h4">
              Welcome to the BoltDB Database Manager
            </Typography>
            <FlexLayoutGrid accesskey={accesskey} filename={dbname} />
          </Grid>
        </Grid>
      ) : (
        <Grid container columns={{ xs: 2, sm: 12, md: 12 }}>
          <Grid item xs={2} sm={5.5} md={5.5} px={5} py={5} borderRadius={5}>
            <FileUpload />
          </Grid>
          <Grid
            item
            xs={2}
            sm={1}
            md={1}
            px={5}
            py={5}
            borderRadius={5}
            textAlign="center"
          >
            <Typography variant="h3" color="secondary.dark">
              OR
            </Typography>
          </Grid>
          <Grid item xs={2} sm={5.5} md={5.5} px={5} py={5} borderRadius={5}>
            <NewFile />
          </Grid>
        </Grid>
      )}
    </Container>
  );
};

export default Controller;
