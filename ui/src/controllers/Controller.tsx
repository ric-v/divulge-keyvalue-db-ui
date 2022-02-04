import { Container, Grid, Typography } from "@mui/material";
import FileUpload from "./FileUpload";
import NewFile from "./NewFile";

const Controller = () => {
  return (
    <Container
      sx={{
        marginTop: "1rem",
        flexGrow: 1,
      }}
    >
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
    </Container>
  );
};

export default Controller;
