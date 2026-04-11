import { useState, useEffect } from 'react';
import { 
  Typography, 
  Container, 
  Box, 
  TextField, 
  Button, 
  Paper, 
  Alert, 
  Grid,
  CircularProgress,
  ToggleButton,
  ToggleButtonGroup
} from '@mui/material';
import { useNavigate, useLocation } from 'react-router-dom';
import { supabase } from '../supabaseClient';
import { API_BASE_URL } from '../config';

const AddCertType = () => {
  const [formData, setFormData] = useState({
    name: '',
    shortName: '',
    stcwReference: '',
    normalValidityMonths: '',
    status: 'provisional'
  });
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [userRole, setUserRole] = useState<string | null>(null);
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    const fetchUserData = async () => {
      const { data: { session } } = await supabase.auth.getSession();
      if (session) {
        // First try to get role from app_metadata
        const role = session.user?.app_metadata?.role;
        if (role) {
          setUserRole(role);
          if (role === 'admin') {
            setFormData(prev => ({ ...prev, status: 'approved' }));
          }
          return;
        }

        // Fallback to fetching
        if (session.access_token) {
          try {
            const response = await fetch(`${API_BASE_URL}/admin/users`, {
              headers: {
                'Authorization': `Bearer ${session.access_token}`,
              },
            });
            if (response.ok) {
              const data = await response.json();
              const user = Array.isArray(data) ? data.find((u: any) => u.id === session.user.id) : data;
              const role = user?.role || 'user';
              setUserRole(role);
              if (role === 'admin') {
                setFormData(prev => ({ ...prev, status: 'approved' }));
              }
            } else {
              setUserRole('user');
            }
          } catch (error) {
            console.error('Error fetching user role:', error);
            setUserRole('user');
          }
        }
      }
    };
    fetchUserData();
  }, []);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setError(null);

    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) throw new Error('Not authenticated');

      const response = await fetch(`${API_BASE_URL}/api/cert-types`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${session.access_token}`,
        },
        body: JSON.stringify({
          name: formData.name,
          'short-name': formData.shortName || null,
          'stcw-reference': formData.stcwReference || null,
          'normal-validity-months': formData.normalValidityMonths ? parseInt(formData.normalValidityMonths) : 0,
          status: formData.status
        }),
      });

      const responseData = await response.json().catch(() => ({}));

      if (!response.ok) {
        const errorMessage = responseData.error || 
                           responseData.message || 
                           (responseData.errors && typeof responseData.errors === 'object' ? JSON.stringify(responseData.errors) : null) ||
                           'Failed to add certificate type';
        throw new Error(errorMessage);
      }

      // Navigate back to CertTypes list or where we came from
      if (location.state?.from === 'add-certificate') {
        navigate('/add-certificate', { 
          state: { 
            certTypeId: responseData.id,
            // Preserve other state if needed, though AddCertificate might not have much else in state
          },
          replace: true 
        });
      } else {
        navigate('/cert-types');
      }
    } catch (err: any) {
      setError(err.message || 'An error occurred during submission');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Container maxWidth="md">
      <Box sx={{ mt: 4, mb: 4 }}>
        <Paper elevation={0} sx={{ p: { xs: 2, sm: 4 }, border: 1, borderColor: 'divider', borderRadius: 2 }}>
          <Typography variant="h5" gutterBottom sx={{ fontWeight: 600, mb: 1 }}>
            Add certificate type
          </Typography>
          <Typography variant="body2" color="textSecondary" sx={{ mb: 3 }}>
            Can't find the certificate type you're looking for? Add it here. 
            Once submitted, it will be available for you to use immediately. 
            An administrator will review it for global approval.
          </Typography>

          <Box component="form" onSubmit={handleSubmit} noValidate>
            <Grid container spacing={3}>
              <Grid size={{ xs: 12 }}>
                <TextField
                  required
                  fullWidth
                  id="name"
                  label="Name"
                  name="name"
                  value={formData.name}
                  onChange={handleChange}
                  autoFocus
                />
              </Grid>

              <Grid size={{ xs: 12, sm: 6 }}>
                <TextField
                  fullWidth
                  id="shortName"
                  label="Short Name"
                  name="shortName"
                  value={formData.shortName}
                  onChange={handleChange}
                />
              </Grid>

              <Grid size={{ xs: 12, sm: 6 }}>
                <TextField
                  fullWidth
                  id="stcwReference"
                  label="STCW Reference"
                  name="stcwReference"
                  value={formData.stcwReference}
                  onChange={handleChange}
                />
              </Grid>

              <Grid size={{ xs: 12, sm: 6 }}>
                <TextField
                  fullWidth
                  id="normalValidityMonths"
                  label="Normal Validity (Months)"
                  name="normalValidityMonths"
                  type="number"
                  value={formData.normalValidityMonths}
                  onChange={handleChange}
                  helperText="Leave blank if certificate does not expire"
                />
              </Grid>

              {userRole === 'admin' && (
                <Grid size={{ xs: 12, sm: 6 }}>
                  <Typography variant="body2" color="textSecondary" gutterBottom>
                    Status
                  </Typography>
                  <ToggleButtonGroup
                    color="primary"
                    value={formData.status}
                    exclusive
                    onChange={(_e, value) => value && setFormData(prev => ({ ...prev, status: value }))}
                    fullWidth
                    size="small"
                  >
                    <ToggleButton value="provisional">Provisional</ToggleButton>
                    <ToggleButton value="approved">Approved</ToggleButton>
                  </ToggleButtonGroup>
                </Grid>
              )}

              <Grid size={{ xs: 12 }}>
                {error && (
                  <Alert severity="error" sx={{ mb: 2 }}>
                    {error}
                  </Alert>
                )}
              </Grid>

              <Grid size={{ xs: 12 }} sx={{ display: 'flex', gap: 2, justifyContent: 'flex-end', mt: 2 }}>
                <Button 
                  variant="outlined" 
                  onClick={() => {
                    if (location.state?.from === 'add-certificate') {
                      navigate('/add-certificate', { replace: true });
                    } else {
                      navigate('/cert-types');
                    }
                  }}
                  disabled={submitting}
                >
                  Cancel
                </Button>
                <Button 
                  type="submit" 
                  variant="contained" 
                  color="primary"
                  disabled={submitting}
                >
                  {submitting ? <CircularProgress size={24} color="inherit" /> : 'Add certificate type'}
                </Button>
              </Grid>
            </Grid>
          </Box>
        </Paper>
      </Box>
    </Container>
  );
};

export default AddCertType;
