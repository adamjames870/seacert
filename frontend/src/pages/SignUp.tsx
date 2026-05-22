import { useState } from 'react';
import { 
  Typography, 
  Container, 
  Box, 
  TextField, 
  Button, 
  Paper, 
  Link,
  Alert,
  FormControlLabel,
  Checkbox,
  Snackbar
} from '@mui/material';
import { usePostHog } from 'posthog-js/react';
import { Link as RouterLink, useNavigate } from 'react-router-dom';
import { API_BASE_URL } from '../config';
import { supabase } from '../supabaseClient';

const SignUp = () => {
  const posthog = usePostHog();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [emailConsent, setEmailConsent] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    if (password !== confirmPassword) {
      setError("Passwords do not match");
      setLoading(false);
      return;
    }
    
    try {
      const { data, error } = await supabase.auth.signUp({
        email,
        password,
      });

      if (error) throw error;

      if (data.user) {
        posthog.identify(data.user.id, {
          email: email
        });
        posthog.capture('signup completed');
      }

      if (data.session?.access_token) {
        // Record consent
        await fetch(`${API_BASE_URL}/admin/users`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${data.session.access_token}`,
          },
          body: JSON.stringify({
            email_consent: emailConsent,
            email_consent_version: '2026-03-01',
            email_consent_source: 'signup_form'
          }),
        });
      }
      
      setSuccess(true);
      setEmail('');
      setPassword('');
      setConfirmPassword('');

      // Redirect to onboarding if we have a session (means they are logged in automatically)
      // If no session, they probably need to confirm email first.
      if (data.session) {
        navigate('/onboarding/first-certificate', { replace: true });
      } else {
        // If no session, redirect to login with a message that they should check email
        // or just stay here and show the success message.
        // Given the requirement "If that is not possible, direct them to login, then to onboarding"
        // we should probably redirect to login with a state that indicates they just signed up.
        navigate('/login', { state: { signupSuccess: true } });
      }
    } catch (err: any) {
      setError(err.message || 'An error occurred during sign up');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Container maxWidth="xs">
      <Box 
        sx={{ 
          mt: 8, 
          display: 'flex', 
          flexDirection: 'column', 
          alignItems: 'center' 
        }}
      >
        <Paper 
          elevation={0} 
          sx={{ 
            p: { xs: 2, sm: 4 }, 
            width: '100%', 
            border: 1, 
            borderColor: 'divider',
            borderRadius: 2
          }}
        >
          <Typography variant="h5" component="h1" gutterBottom align="center" sx={{ fontWeight: 600 }}>
            Create Account
          </Typography>
          <Typography variant="body2" color="text.secondary" align="center" sx={{ mb: 3 }}>
            Sign up to get started with SeaCert
          </Typography>

          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}

          <Box component="form" onSubmit={handleSubmit} noValidate>
            <TextField
              margin="normal"
              required
              fullWidth
              id="email"
              label="Email Address"
              name="email"
              autoComplete="email"
              autoFocus
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              variant="outlined"
            />
            <TextField
              margin="normal"
              required
              fullWidth
              name="password"
              label="Password"
              type="password"
              id="password"
              autoComplete="new-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              variant="outlined"
            />
            <TextField
              margin="normal"
              required
              fullWidth
              name="confirmPassword"
              label="Confirm Password"
              type="password"
              id="confirmPassword"
              autoComplete="new-password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              variant="outlined"
            />
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
              sx={{ mt: 1, textAlign: 'left' }}
            />
            <Button
              type="submit"
              fullWidth
              variant="contained"
              color="primary"
              disabled={loading}
              sx={{ mt: 3, mb: 2, py: 1.2, fontWeight: 600 }}
            >
              {loading ? 'Creating Account...' : 'Sign Up'}
            </Button>
            <Box sx={{ display: 'flex', justifyContent: 'center', mt: 1 }}>
              <Link component={RouterLink} to="/login" variant="body2" sx={{ textDecoration: 'none' }}>
                Already have an account? Login
              </Link>
            </Box>
          </Box>
        </Paper>
      </Box>

      <Snackbar 
        open={success} 
        autoHideDuration={6000} 
        onClose={() => setSuccess(false)}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      >
        <Alert onClose={() => setSuccess(false)} severity="success" sx={{ width: '100%' }}>
          Registration successful! Please check your email for confirmation.
        </Alert>
      </Snackbar>
    </Container>
  );
};

export default SignUp;
