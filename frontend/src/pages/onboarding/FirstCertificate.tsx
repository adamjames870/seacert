import { useState } from 'react';
import { 
  Typography, 
  Container, 
  Box, 
  Button, 
  Paper, 
  Stack,
  Dialog,
  DialogTitle,
  DialogContent
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { usePostHog } from 'posthog-js/react';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';
import { UserDetailsForm } from '../../components/UserDetailsForm';
import { API_BASE_URL } from '../../config';
import { supabase } from '../../supabaseClient';

const FirstCertificate = () => {
  const navigate = useNavigate();
  const posthog = usePostHog();
  const [isDialogOpen, setIsDialogOpen] = useState(true);

  const sendWelcomeEmail = async () => {
    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) return;

      const response = await fetch(`${API_BASE_URL}/notify/welcome-email`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${session.access_token}`
        },
      });

      if (!response.ok) {
        throw new Error(`Failed to send welcome email: ${response.statusText}`);
      }
    } catch (e) {
      console.error('Error sending welcome email', e);
    }
  };

  const handleUploadClick = () => {
    posthog.capture('onboarding upload clicked');
    setIsDialogOpen(true);
  };

  const handleSaveUserDetails = async (data: any) => {
    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) return;

      // Fetch existing user data
      const response = await fetch(`${API_BASE_URL}/admin/users`, {
        headers: {
          'Authorization': `Bearer ${session.access_token}`,
        },
      });

      if (!response.ok) {
        throw new Error('Failed to fetch user data');
      }

      const userData = await response.json();
      const currentUser = Array.isArray(userData) ? userData.find((u: any) => u.id === session.user.id) : userData;

      // Merge changes
      const updatedUser = {
        ...currentUser,
        forename: data.forename,
        surname: data.surname,
        nationality: data.nationality || null,
      };

      const putResponse = await fetch(`${API_BASE_URL}/admin/users`, {
        method: 'PUT',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${session.access_token}`
        },
        body: JSON.stringify(updatedUser),
      });

      if (!putResponse.ok) {
        throw new Error('Failed to update user details');
      }

      await sendWelcomeEmail();
      setIsDialogOpen(false);
      navigate('/certificate-wizard');
    } catch (e) {
      console.error('Error saving user details', e);
      // Optionally notify user about error
    }
  };

  const handleCloseDialog = async () => {
    await sendWelcomeEmail();
    setIsDialogOpen(false);
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

      <Dialog open={isDialogOpen} onClose={handleCloseDialog} fullWidth maxWidth="sm">
        <DialogTitle>Complete your profile</DialogTitle>
        <DialogContent>
          <UserDetailsForm onSave={handleSaveUserDetails} onCancel={handleCloseDialog} />
        </DialogContent>
      </Dialog>
    </Container>
  );
};

export default FirstCertificate;
