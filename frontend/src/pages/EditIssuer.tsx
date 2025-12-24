import { useState, useEffect } from 'react';
import { 
  Typography, 
  Container, 
  Box, 
  TextField, 
  Button, 
  Paper, 
  Alert, 
  CircularProgress,
  Grid,
  Autocomplete
} from '@mui/material';
import { useNavigate, useParams, useLocation } from 'react-router-dom';
import { supabase } from '../supabaseClient';
import { API_BASE_URL } from '../config';
import { countries } from '../utils/countryData';

const EditIssuer = () => {
  const { id } = useParams<{ id: string }>();
  const location = useLocation();
  const [formData, setFormData] = useState({
    name: '',
    country: '',
    website: ''
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const navigate = useNavigate();

  const fromCertificates = location.state?.from === 'certificates';

  useEffect(() => {
    const fetchIssuer = async () => {
      try {
        setLoading(true);
        const { data: { session } } = await supabase.auth.getSession();
        if (!session) throw new Error('Not authenticated');

        const response = await fetch(`${API_BASE_URL}/api/issuers`, {
          headers: {
            'Authorization': `Bearer ${session.access_token}`,
          },
        });

        if (!response.ok) throw new Error('Failed to fetch issuer');

        const data = await response.json();
        const issuer = data.find((i: any) => i.id === id);

        if (!issuer) throw new Error('Issuer not found');

        setFormData({
          name: issuer.name || '',
          country: issuer.country || '',
          website: issuer.website || ''
        });
      } catch (err: any) {
        setError(err.message || 'An error occurred');
      } finally {
        setLoading(false);
      }
    };

    fetchIssuer();
  }, [id]);

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

      const body: any = {
        id: id,
        name: formData.name,
        country: formData.country || null,
        website: formData.website || null
      };

      const response = await fetch(`${API_BASE_URL}/api/issuers?id=${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${session.access_token}`,
        },
        body: JSON.stringify(body),
      });

      if (!response.ok) {
        const responseData = await response.json().catch(() => ({}));
        throw new Error(responseData.message || 'Failed to update issuer');
      }

      if (fromCertificates) {
        navigate('/certificates');
      } else {
        navigate('/issuers');
      }
    } catch (err: any) {
      setError(err.message || 'An error occurred during submission');
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
            Edit Issuer
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
                        autoComplete: 'new-password',
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
                  onClick={() => fromCertificates ? navigate('/certificates') : navigate('/issuers')}
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
                  {submitting ? <CircularProgress size={24} color="inherit" /> : 'Update Issuer'}
                </Button>
              </Grid>
            </Grid>
          </Box>
        </Paper>
      </Box>
    </Container>
  );
};

export default EditIssuer;
