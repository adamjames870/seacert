import { useState } from 'react';
import { 
  Dialog, 
  DialogTitle, 
  DialogContent, 
  DialogActions, 
  Button, 
  Typography, 
  FormControlLabel, 
  Checkbox,
  Box,
  Link
} from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';
import { API_BASE_URL } from '../config';
import { supabase } from '../supabaseClient';

interface EmailConsentModalProps {
  open: boolean;
  onClose: (updatedUserData?: any) => void;
}

const EmailConsentModal = ({ open, onClose }: EmailConsentModalProps) => {
  const [emailConsent, setEmailConsent] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async () => {
    setLoading(true);
    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) return;

      const response = await fetch(`${API_BASE_URL}/admin/users`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${session.access_token}`,
        },
        body: JSON.stringify({
          email_consent: emailConsent,
          email_consent_version: '2026-03-01',
          email_consent_source: 'consent_banner'
        }),
      });

      if (response.ok) {
        const updatedUser = await response.json();
        onClose(updatedUser);
      } else {
        onClose();
      }
    } catch (error) {
      console.error('Error updating consent:', error);
      onClose();
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} maxWidth="sm" fullWidth>
      <DialogTitle sx={{ fontWeight: 600 }}>We value your privacy</DialogTitle>
      <DialogContent>
        <Typography variant="body1" paragraph>
          Before you continue, please review our updated{' '}
          <Link component={RouterLink} to="/privacy" target="_blank">Privacy Policy</Link> and{' '}
          <Link component={RouterLink} to="/terms" target="_blank">Terms & Conditions</Link>.
        </Typography>
        <Typography variant="body2" color="text.secondary" paragraph>
          We would also like to keep you updated with important notifications and news about your seafarer certificates.
        </Typography>
        <Box sx={{ mt: 2 }}>
          <FormControlLabel
            control={
              <Checkbox 
                checked={emailConsent} 
                onChange={(e) => setEmailConsent(e.target.checked)} 
                color="primary" 
              />
            }
            label={
              <Typography variant="body2">
                I agree to receive email updates and notifications from SeaCert.
              </Typography>
            }
          />
        </Box>
      </DialogContent>
      <DialogActions sx={{ p: 3 }}>
        <Button 
          variant="contained" 
          onClick={handleSubmit} 
          disabled={loading}
          fullWidth
          size="large"
        >
          {loading ? 'Saving...' : 'Accept & Continue'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default EmailConsentModal;
