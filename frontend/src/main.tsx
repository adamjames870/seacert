import { StrictMode, useMemo } from 'react'
import { createRoot } from 'react-dom/client'
import { CssBaseline, ThemeProvider, createTheme } from '@mui/material'
import { BrowserRouter } from 'react-router-dom'
import './index.css'
import App from './App.tsx'

function Main() {
  const theme = useMemo(
    () =>
      createTheme({
        palette: {
          mode: 'light',
          primary: {
            main: '#4a6d8c', // Muted blues
          },
          secondary: {
            main: '#6b8ca6',
          },
          background: {
            default: '#f4f7f9',
          },
        },
      }),
    [],
  );

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <App />
    </ThemeProvider>
  );
}

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <Main />
    </BrowserRouter>
  </StrictMode>,
)
