import { useEffect } from 'react';
import { Container, Paper, Box } from '@mui/material';

const Terms = () => {
  useEffect(() => {
    const id = 'termly-jssdk';
    if (document.getElementById(id)) return;
    
    const js = document.createElement('script');
    js.id = id;
    js.src = "https://app.termly.io/embed-policy.min.js";
    
    const tjs = document.getElementsByTagName('script')[0];
    if (tjs && tjs.parentNode) {
      tjs.parentNode.insertBefore(js, tjs);
    } else {
      document.head.appendChild(js);
    }
  }, []);

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Paper elevation={0} sx={{ p: { xs: 2, md: 4 }, borderRadius: 2, border: '1px solid', borderColor: 'divider' }}>
        <Box className="terms-conditions-container">
          <div name="termly-embed" data-id="ade141d5-370d-4936-8e86-5c741af7e75c"></div>
        </Box>
      </Paper>
    </Container>
  );
};

export default Terms;
