import { 
  Typography, 
  Container, 
  Box, 
  Button, 
  Paper, 
  Stack 
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { usePostHog } from 'posthog-js/react';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';

const FirstCertificate = () => {
  const navigate = useNavigate();
  const posthog = usePostHog();

  const handleUploadClick = () => {
    posthog.capture('onboarding upload clicked');
    navigate('/certificate-wizard');
  };

  const handleSkipClick = () => {
    posthog.capture('onboarding skipped');
    navigate('/');
  };

  return (
    <Container maxWidth="sm">
      <Box 
        sx={{ 
          mt: 8, 
          display: 'flex', 
          flexDirection: 'column', 
          alignItems: 'center',
          textAlign: 'center'
        }}
      >
        <Paper 
          elevation={0} 
          sx={{ 
            p: { xs: 4, sm: 6 }, 
            width: '100%', 
            border: 1, 
            borderColor: 'divider',
            borderRadius: 4,
            bgcolor: 'background.paper'
          }}
        >
          <Typography variant="h4" component="h1" gutterBottom sx={{ fontWeight: 700, mb: 2 }}>
            Add your first certificate
          </Typography>
          <Typography variant="body1" color="text.secondary" sx={{ mb: 6, fontSize: '1.1rem' }}>
            Upload a photo or PDF — we’ll fill everything automatically
          </Typography>

          <Stack spacing={2} alignItems="center">
            <Button
              variant="contained"
              color="primary"
              size="large"
              fullWidth
              startIcon={<CloudUploadIcon />}
              onClick={handleUploadClick}
              sx={{ 
                py: 2, 
                fontSize: '1.1rem', 
                fontWeight: 600,
                borderRadius: 2,
                textTransform: 'none'
              }}
            >
              Upload Certificate
            </Button>
            <Button
              variant="text"
              color="inherit"
              size="small"
              onClick={handleSkipClick}
              sx={{ 
                textTransform: 'none',
                color: 'text.secondary',
                '&:hover': {
                  bgcolor: 'transparent',
                  textDecoration: 'underline'
                }
              }}
            >
              Skip for now
            </Button>
          </Stack>
        </Paper>
      </Box>
    </Container>
  );
};

export default FirstCertificate;
