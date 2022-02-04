import * as React from "react";
import { Box, Grid } from "@mui/material";
import NavBar from "./components/NavBar";
import { ThemeProvider, createTheme } from "@mui/material/styles";
import { SnackbarProvider } from "notistack";
import "@fontsource/fira-sans";
import Controller from "./controllers/Controller";

const ColorModeContext = React.createContext({ toggleColorMode: () => {} });

function App() {
  const [mode, setMode] = React.useState<"light" | "dark">("light");
  const colorMode = React.useMemo(
    () => ({
      toggleColorMode: () => {
        setMode((prevMode) => (prevMode === "light" ? "dark" : "light"));
      },
    }),
    []
  );

  const theme = React.useMemo(
    () =>
      createTheme({
        palette: {
          mode: mode,
          primary: {
            main: "#594f72",
            light: "#867ba1",
            dark: "#2f2746",
          },
          secondary: {
            light: "#fff990",
            main: "#e2c660",
            dark: "#ae9631",
          },
          contrastThreshold: 3,
          tonalOffset: 0.2,
        },
        typography: {
          fontFamily: ["Fira Sans", "sans-serif"].join(","),
          h1: {
            fontSize: "1.7rem",
            fontWeight: "bold",
            fontVariant: "small-caps",
            fontVariantLigatures: "common-ligatures",
            "@media (min-width:800px)": {
              fontSize: "2.5rem",
            },
            "@media (min-width:1600px)": {
              fontSize: "3.5rem",
            },
          },
          h2: {
            fontSize: "1.5rem",
            fontWeight: "bold",
            fontVariantLigatures: "common-ligatures",
            fontVariant: "small-caps",
            "@media (min-width:800px)": {
              fontSize: "2.2rem",
            },
            "@media (min-width:1600px)": {
              fontSize: "3rem",
            },
          },
          h3: {
            fontSize: "1.4rem",
            fontWeight: "bold",
            fontVariantLigatures: "common-ligatures",
            fontVariant: "small-caps",
            "@media (min-width:800px)": {
              fontSize: "2rem",
            },
            "@media (min-width:1600px)": {
              fontSize: "2.8rem",
            },
          },
          h4: {
            fontSize: "1.2rem",
            fontWeight: "bold",
            fontVariantLigatures: "common-ligatures",
            fontVariant: "small-caps",
            "@media (min-width:800px)": {
              fontSize: "1.8rem",
            },
            "@media (min-width:1600px)": {
              fontSize: "2.5rem",
            },
          },
        },
        components: {
          MuiCssBaseline: {
            styleOverrides: `
              @font-face {
                font-family: 'Fira Sans';
              }
            `,
          },
        },
      }),
    [mode]
  );

  return (
    <SnackbarProvider
      maxSnack={3}
      anchorOrigin={{
        vertical: "top",
        horizontal: "right",
      }}
      autoHideDuration={3000}
    >
      <ColorModeContext.Provider value={colorMode}>
        <ThemeProvider theme={theme}>
          <Grid
            container
            spacing={0}
            direction="column"
            alignItems="stretch"
            justifyContent="center"
            sx={{
              my: "auto",
              bgcolor: "background.default",
              color: "text.primary",
              minHeight: "100vh",
              py: [5, 8],
            }}
          >
            <Grid item>
              <NavBar
                toggleColorMode={colorMode.toggleColorMode}
                theme={theme}
              />
              {/* margin */}
              <Controller />
            </Grid>
          </Grid>
        </ThemeProvider>
      </ColorModeContext.Provider>
    </SnackbarProvider>
  );
}

export default App;
