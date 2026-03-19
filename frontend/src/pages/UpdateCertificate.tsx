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
  Autocomplete,
  Link
} from '@mui/material';
import { useNavigate, useParams, Link as RouterLink, useLocation } from 'react-router-dom';
import { supabase } from '../supabaseClient';
import { API_BASE_URL } from '../config';

interface Issuer {
  id: string;
  name: string;
}

const UpdateCertificate = () => {
  const { id } = useParams<{ id: string }>();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [issuers, setIssuers] = useState<Issuer[]>([]);

  const [formData, setFormData] = useState({
    certNumber: '',
    issuedDate: '',
    remarks: '',
    issuerId: ''
  });
  const [certTypeName, setCertTypeName] = useState('');

  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    const fetchCertificate = async () => {
      try {
        setLoading(true);
        const { data: { session } } = await supabase.auth.getSession();
        
        if (!session) {
          setError('Not authenticated');
          setLoading(false);
          return;
        }

        const [certResponse, issuersResponse] = await Promise.all([
          fetch(`${API_BASE_URL}/api/certificates`, {
            headers: {
              'Authorization': `Bearer ${session.access_token}`,
            },
          }),
          fetch(`${API_BASE_URL}/api/issuers`, {
            headers: {
              'Authorization': `Bearer ${session.access_token}`,
            },
          })
        ]);

        if (!certResponse.ok || !issuersResponse.ok) {
          throw new Error('Failed to fetch certificate details or issuers');
        }

        const [certsData, issuersData] = await Promise.all([
          certResponse.json(),
          issuersResponse.json()
        ]);

        const certs = Array.isArray(certsData) ? certsData : (certsData.certificates || []);
        const cert = certs.find((c: any) => c.id === id);

        if (!cert) {
          throw new Error('Certificate not found');
        }

        const sortedIssuers = (issuersData as Issuer[]).sort((a, b) => 
          a.name.localeCompare(b.name)
        );
        setIssuers(sortedIssuers);

        setFormData({
          certNumber: cert['cert-number'] || '',
          issuedDate: cert['issued-date'] ? cert['issued-date'].split('T')[0] : '',
          remarks: cert.remarks || '',
          issuerId: cert['issuer-id'] || ''
        });
        setCertTypeName(cert['cert-type-name'] || '');
      } catch (err: any) {
        setError(err.message || 'An error occurred while fetching data');
      } finally {
        setLoading(false);
      }
    };

    fetchCertificate();
  }, [id]);

  useEffect(() => {
    if (location.state?.newIssuerId) {
      setFormData(prev => ({ ...prev, issuerId: location.state.newIssuerId }));
      // Clear state
      navigate(location.pathname, { replace: true, state: {} });
    }
  }, [location.state, location.pathname, navigate]);

  const handleChange = (e: any) => {
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

      const payload = {
        id: id,
        'issuer-id': formData.issuerId,
        'cert-number': formData.certNumber,
        'issued-date': formData.issuedDate ? new Date(formData.issuedDate).toISOString() : null,
        remarks: formData.remarks
      };
      console.log('Sending update payload:', payload);

      const response = await fetch(`${API_BASE_URL}/api/certificates`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${session.access_token}`,
        },
        body: JSON.stringify(payload),
      });

      const responseData = await response.json().catch(() => ({}));
      console.log('Update response:', responseData);

      if (!response.ok) {
        // Extract as much error information as possible
        const errorMessage = responseData.message || 
                           responseData.error || 
                           (responseData.errors && typeof responseData.errors === 'object' ? JSON.stringify(responseData.errors) : null) ||
                           'Failed to update certificate';
        throw new Error(errorMessage);
      }

      // Check if the returned certificate matches expectations (e.g. has updated fields)
      console.log('Successfully updated certificate:', responseData);
      
      // We could optionally verify if the returned data matches formData here,
      // but if the response is OK, we assume the backend handled it or returned the new state.

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
          <Typography variant="h5" gutterBottom sx={{ fontWeight: 600, mb: 1 }}>
            Update Certificate
          </Typography>
          <Typography variant="subtitle1" color="text.secondary" sx={{ mb: 3 }}>
            {certTypeName}
          </Typography>

          {error && (
            <Alert severity="error" sx={{ mb: 3 }}>
              {error}
            </Alert>
          )}

          <Box component="form" onSubmit={handleSubmit} noValidate>
            <Grid container spacing={3}>
              <Grid size={{ xs: 12 }}>
                <Autocomplete
                  id="issuerId"
                  options={issuers}
                  getOptionLabel={(option) => option.name}
                  autoHighlight
                  autoSelect
                  value={issuers.find(issuer => issuer.id === formData.issuerId) || null}
                  onChange={(_event, newValue) => {
                    setFormData(prev => ({ ...prev, issuerId: newValue ? newValue.id : '' }));
                  }}
                  renderInput={(params) => (
                    <TextField
                      {...params}
                      label="Issuer"
                      required
                      error={!formData.issuerId && submitting}
                    />
                  )}
                />
                <Box sx={{ mt: 1 }}>
                  <Typography variant="caption">
                    Can't find your issuer?{' '}
                    <Link component={RouterLink} to="/add-issuer" state={{ from: 'update-certificate', id }} sx={{ textDecoration: 'none' }}>
                      Add New Issuer
                    </Link>
                  </Typography>
                </Box>
              </Grid>

              <Grid size={{ xs: 12, sm: 6 }}>
                <TextField
                  required
                  fullWidth
                  id="certNumber"
                  label="Certificate Number"
                  name="certNumber"
                  value={formData.certNumber}
                  onChange={handleChange}
                />
              </Grid>

              <Grid size={{ xs: 12, sm: 6 }}>
                <TextField
                  required
                  fullWidth
                  id="issuedDate"
                  label="Issued Date"
                  name="issuedDate"
                  type="date"
                  value={formData.issuedDate}
                  onChange={handleChange}
                  InputLabelProps={{ shrink: true }}
                />
              </Grid>

              <Grid size={{ xs: 12 }}>
                <TextField
                  fullWidth
                  id="remarks"
                  label="Remarks"
                  name="remarks"
                  multiline
                  rows={3}
                  value={formData.remarks}
                  onChange={handleChange}
                />
              </Grid>

              <Grid size={{ xs: 12 }} sx={{ display: 'flex', gap: 2, justifyContent: 'flex-end', mt: 2 }}>
                <Button 
                  variant="outlined" 
                  onClick={() => navigate('/certificates')}
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
                  {submitting ? 'Updating...' : 'Update Certificate'}
                </Button>
              </Grid>
            </Grid>
          </Box>
        </Paper>
      </Box>
    </Container>
  );
};

export default UpdateCertificate;
