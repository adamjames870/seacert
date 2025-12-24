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
  Autocomplete
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { supabase } from '../supabaseClient';
import { API_BASE_URL } from '../config';
import { countries } from '../utils/countryData';

const EditAccount = () => {
  const [formData, setFormData] = useState({
    forename: '',
    surname: '',
    nationality: ''
  });
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        setLoading(true);
        const { data: { session } } = await supabase.auth.getSession();
        
        if (!session) {
          setError('Not authenticated');
          setLoading(false);
          return;
        }

        const response = await fetch(`${API_BASE_URL}/admin/users`, {
          headers: {
            'Authorization': `Bearer ${session.access_token}`,
          },
        });

        if (!response.ok) {
          throw new Error('Failed to fetch user details');
        }

        const data = await response.json();
        // Based on App.tsx, the API might return an array
        const user = Array.isArray(data) ? data.find((u: any) => u.id === session.user.id) : data;

        if (user) {
          setFormData({
            forename: user.forename || '',
            surname: user.surname || '',
            nationality: user.nationality || ''
          });
        }
      } catch (err: any) {
        setError(err.message || 'An error occurred while fetching user data');
      } finally {
        setLoading(false);
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

      const response = await fetch(`${API_BASE_URL}/admin/users`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${session.access_token}`,
        },
        body: JSON.stringify({
          forename: formData.forename,
          surname: formData.surname,
          nationality: formData.nationality || null
        }),
      });

      if (!response.ok) {
        const responseData = await response.json().catch(() => ({}));
        const errorMessage = responseData.message || 
                           responseData.error || 
                           (responseData.errors && typeof responseData.errors === 'object' ? JSON.stringify(responseData.errors) : null) ||
                           'Failed to update user details';
        throw new Error(errorMessage);
      }

      navigate('/certificates');
    } catch (err: any) {
      setError(err.message || 'An error occurred during update');
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', mt: 8 }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Container maxWidth="md">
      <Box sx={{ mt: 4, mb: 4 }}>
        <Paper elevation={0} sx={{ p: { xs: 2, sm: 4 }, border: 1, borderColor: 'divider', borderRadius: 2 }}>
          <Typography variant="h5" gutterBottom sx={{ fontWeight: 600, mb: 3 }}>
            Edit Account Details
          </Typography>

          {error && (
            <Alert severity="error" sx={{ mb: 3 }}>
              {error}
            </Alert>
          )}

          <Box component="form" onSubmit={handleSubmit} noValidate>
            <Grid container spacing={3}>
              <Grid size={{ xs: 12, sm: 6 }}>
                <TextField
                  required
                  fullWidth
                  id="forename"
                  label="Forename"
                  name="forename"
                  value={formData.forename}
                  onChange={handleChange}
                  autoFocus
                />
              </Grid>
              <Grid size={{ xs: 12, sm: 6 }}>
                <TextField
                  required
                  fullWidth
                  id="surname"
                  label="Surname"
                  name="surname"
                  value={formData.surname}
                  onChange={handleChange}
                />
              </Grid>
              <Grid size={{ xs: 12 }}>
                <Autocomplete
                  id="nationality-select"
                  options={countries}
                  autoHighlight
                  getOptionLabel={(option) => option.label}
                  value={countries.find(c => c.code === formData.nationality) || null}
                  onChange={(_event, newValue) => {
                    setFormData(prev => ({ ...prev, nationality: newValue ? newValue.code : '' }));
                  }}
                  renderInput={(params) => (
                    <TextField
                      {...params}
                      label="Nationality"
                      inputProps={{
                        ...params.inputProps,
                        autoComplete: 'new-password',
                      }}
                    />
                  )}
                />
              </Grid>

              <Grid size={{ xs: 12 }} sx={{ display: 'flex', gap: 2, justifyContent: 'flex-end', mt: 2 }}>
                <Button 
                  variant="outlined" 
                  onClick={() => navigate(-1)}
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
                  {submitting ? 'Saving...' : 'Save Changes'}
                </Button>
              </Grid>
            </Grid>
          </Box>
        </Paper>
      </Box>
    </Container>
  );
};

export default EditAccount;
