import * as React from 'react';
import { Box, Button, Container } from "@mui/material";
import NavBar from "./components/NavBar";
import { ThemeProvider, createTheme } from '@mui/material/styles';
import IconButton from '@mui/material/IconButton';
import Brightness4Icon from '@mui/icons-material/Brightness4';
import Brightness7Icon from '@mui/icons-material/Brightness7';

const ColorModeContext = React.createContext({ toggleColorMode: () => { } });

function App() {


  const [mode, setMode] = React.useState<'light' | 'dark'>('dark');
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
      }),
    [mode],
  );

  const newFile = (fileName: string) => {

    // GET the file from API server
    fetch('http://localhost:8080/api/v1/new?file=' + fileName + '&dbtype=boltdb', {
      method: 'POST',
    })
      .then(response => response.json())
      .then(data => {
        console.log('success', data);
        // store the access key to local storage
        localStorage.setItem('access_key', data.dbkey);
        localStorage.setItem('file_name', data.filename);
        localStorage.setItem('db_type', data.dbtype);
        console.log('success', data);
      })
      .catch(err => {
        console.log('error', err);
      });
  }

  return (
    <ColorModeContext.Provider value={colorMode}>
      <ThemeProvider theme={theme}>
        <Box sx={{
          bgcolor: 'background.default',
          color: 'text.primary',
          minHeight: '100vh',
          minWidth: '100%',
        }}>
          <NavBar />
          <Container>
            <div>
              <h1>Hello World</h1>
              <Button variant='contained' color="secondary" onClick={() => newFile('new_file')}>
                click
              </Button>
            </div>
            <IconButton sx={{ ml: 1 }} onClick={colorMode.toggleColorMode} color="inherit">
              {theme.palette.mode === 'dark' ? <Brightness7Icon /> : <Brightness4Icon />}
            </IconButton>
          </Container>
        </Box>
      </ThemeProvider>
    </ColorModeContext.Provider>
  );
}

export default App;
