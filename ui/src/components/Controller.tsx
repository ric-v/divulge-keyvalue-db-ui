import { Container, Grid } from "@mui/material";
import Upload from "./Upload";

const Controller = () => {

  return (
    <Container sx={{
      marginTop: "1rem",
      flexGrow: 1,
    }}>
      <Grid container columns={{ xs: 2, sm: 12, md: 12 }}>
        <Grid item xs={2} sm={6} md={6}>
          <Upload />
        </Grid>
      </Grid>
    </Container>
  );
};

export default Controller;
