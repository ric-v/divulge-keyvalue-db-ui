import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import { SnackbarProvider } from 'notistack';

ReactDOM.render(
  <React.StrictMode>
    <SnackbarProvider
      maxSnack={3}
      anchorOrigin={{
        vertical: 'bottom',
        horizontal: 'right',
      }}
    >
      <App />
    </SnackbarProvider>
  </React.StrictMode>,
  document.getElementById('root')
);
