import * as React from 'react';
import { Box } from "@mui/material";
import NavBar from "./components/NavBar";
import { ThemeProvider, createTheme } from '@mui/material/styles';
import { SnackbarProvider } from 'notistack';
import '@fontsource/fira-sans';
import Controller from './components/Controller';

const ColorModeContext = React.createContext({ toggleColorMode: () => { } });

function App() {

  const [mode, setMode] = React.useState<'light' | 'dark'>('light');
  const colorMode = React.useMemo(
    () => ({
      toggleColorMode: () => {
        setMode((prevMode) => (prevMode === 'light' ? 'dark' : 'light'));
      },
    }),
    [],
  );

  const theme = React.useMemo(
    () =>
      createTheme({
        palette: {
          mode: mode,
          primary: {
            main: '#594f72',
            light: '#867ba1',
            dark: '#2f2746',
          },
          secondary: {
            light: '#fff990',
            main: '#e2c660',
            dark: '#ae9631',
          },
          contrastThreshold: 3,
          tonalOffset: 0.2,
        },
        typography: {
          fontFamily: ['Fira Sans', 'sans-serif'].join(','),
        },
        components: {
          MuiCssBaseline: {
            styleOverrides: `
              @font-face {
                font-family: 'Fira Sans';
              }
            `
          }
        }
      }),
    [mode],
  );

  return (
    <SnackbarProvider
      maxSnack={3}
      anchorOrigin={{
        vertical: 'top',
        horizontal: 'right',
      }}
    >
      <ColorModeContext.Provider value={colorMode}>
        <ThemeProvider theme={theme}>
          <Box sx={{
            bgcolor: 'background.default',
            color: 'text.primary',
            minHeight: '100vh',
            minWidth: '100%',
            py: [5, 8],
          }}>
            <NavBar toggleColorMode={colorMode.toggleColorMode} theme={theme} />
            {/* margin */}
            <Controller />
          </Box>
        </ThemeProvider>
      </ColorModeContext.Provider>
    </SnackbarProvider>
  );
}

export default App;
