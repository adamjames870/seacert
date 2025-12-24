import { useState } from 'react';
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
import { useNavigate, useLocation } from 'react-router-dom';
import { supabase } from '../supabaseClient';
import { API_BASE_URL } from '../config';
import { countries } from '../utils/countryData';

const AddIssuer = () => {
  const [formData, setFormData] = useState({
    name: '',
    country: '',
    website: ''
  });
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();

  const from = location.state?.from || 'add-certificate';

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

      const response = await fetch(`${API_BASE_URL}/api/issuers`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${session.access_token}`,
        },
        body: JSON.stringify({
          name: formData.name,
          country: formData.country || null,
          website: formData.website || null
        }),
      });

      const responseData = await response.json().catch(() => ({}));

      if (!response.ok) {
        const errorMessage = responseData.message || 
                           responseData.error || 
                           (responseData.errors && typeof responseData.errors === 'object' ? JSON.stringify(responseData.errors) : null) ||
                           'Failed to add issuer';
        throw new Error(errorMessage);
      }

      // Navigate back to the caller
      if (from === 'issuers') {
        navigate('/issuers');
      } else {
        navigate('/add-certificate', { state: { newIssuerId: responseData.id } });
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
          <Typography variant="h5" gutterBottom sx={{ fontWeight: 600, mb: 3 }}>
            Add New Issuer
          </Typography>

          {error && (
            <Alert severity="error" sx={{ mb: 3 }}>
              {error}
            </Alert>
          )}

          <Box component="form" onSubmit={handleSubmit} noValidate>
            <Grid container spacing={3}>
              <Grid size={{ xs: 12 }}>
                <TextField
                  required
                  fullWidth
                  id="name"
                  label="Issuer Name"
                  name="name"
                  value={formData.name}
                  onChange={handleChange}
                  autoFocus
                />
              </Grid>

              <Grid size={{ xs: 12, sm: 6 }}>
                <Autocomplete
                  id="country-select"
                  options={countries}
                  autoHighlight
                  getOptionLabel={(option) => option.label}
                  value={countries.find(c => c.code === formData.country) || null}
                  onChange={(_event, newValue) => {
                    setFormData(prev => ({ ...prev, country: newValue ? newValue.code : '' }));
                  }}
                  renderInput={(params) => (
                    <TextField
                      {...params}
                      label="Country"
                      inputProps={{
                        ...params.inputProps,
                        autoComplete: 'new-password', // disable autocomplete and autofill
                      }}
                    />
                  )}
                />
              </Grid>

              <Grid size={{ xs: 12, sm: 6 }}>
                <TextField
                  fullWidth
                  id="website"
                  label="Website"
                  name="website"
                  value={formData.website}
                  onChange={handleChange}
                  placeholder="https://example.com"
                />
              </Grid>

              <Grid size={{ xs: 12 }} sx={{ display: 'flex', gap: 2, justifyContent: 'flex-end', mt: 2 }}>
                <Button 
                  variant="outlined" 
                  onClick={() => navigate(from === 'issuers' ? '/issuers' : '/add-certificate')}
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
                  {submitting ? <CircularProgress size={24} color="inherit" /> : 'Add Issuer'}
                </Button>
              </Grid>
            </Grid>
          </Box>
        </Paper>
      </Box>
    </Container>
  );
};

export default AddIssuer;
