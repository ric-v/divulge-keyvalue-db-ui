import { HorizontalRule } from "@mui/icons-material";
import { Container, Grid } from "@mui/material";
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
        <Grid item xs={2} sm={6} md={6}>
          <FileUpload />
        </Grid>
        <Grid item xs={2} sm={6} md={6}>
          <NewFile />
        </Grid>
      </Grid>
    </Container>
  );
};

export default Controller;
