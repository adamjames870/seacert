import { useState, useEffect } from 'react';
import { 
  Snackbar, 
  Button, 
  Typography, 
  Box, 
  Paper
} from '@mui/material';
import CookieIcon from '@mui/icons-material/Cookie';

const CookieConsent = () => {
  const [open, setOpen] = useState(false);

  useEffect(() => {
    // Check if user has already acknowledged the cookie information
    const consent = localStorage.getItem('cookie-consent');
    if (!consent) {
      setOpen(true);
    }
  }, []);

  const handleAccept = () => {
    localStorage.setItem('cookie-consent', 'true');
    setOpen(false);
  };

  return (
    <Snackbar
      open={open}
      anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      sx={{ width: { xs: '100%', sm: 'auto' }, maxWidth: '600px' }}
    >
      <Paper 
        elevation={6} 
        sx={{ 
          p: 3, 
          backgroundColor: '#ffffff', 
          borderTop: '4px solid',
          borderColor: 'primary.main',
          display: 'flex',
          flexDirection: { xs: 'column', sm: 'row' },
          alignItems: 'center',
          gap: 2
        }}
      >
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <CookieIcon color="primary" />
          <Typography variant="body2" color="text.secondary">
            We use cookies to improve your experience and for analytics. By continuing to use our site, you agree to our use of cookies.
          </Typography>
        </Box>
        <Box sx={{ display: 'flex', gap: 1, ml: 'auto', width: { xs: '100%', sm: 'auto' }, justifyContent: 'flex-end' }}>
          <Button 
            variant="contained" 
            onClick={handleAccept}
            size="small"
            sx={{ whiteSpace: 'nowrap' }}
          >
            I understand
          </Button>
        </Box>
      </Paper>
    </Snackbar>
  );
};

export default CookieConsent;
