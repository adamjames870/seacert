import { useState, useEffect } from 'react';
import { 
  Typography, 
  Container, 
  Box, 
  TextField, 
  Button, 
  Paper, 
  Alert, 
  Link,
  CircularProgress,
  Autocomplete,
  Grid
} from '@mui/material';
import { useNavigate, Link as RouterLink, useLocation } from 'react-router-dom';
import { supabase } from '../supabaseClient';

interface CertType {
  id: string;
  name: string;
  'short-name': string;
  'stcw-ref': string;
}

interface Issuer {
  id: string;
  name: string;
  country: string;
  website: string;
}

const AddCertificate = () => {
  const [certTypes, setCertTypes] = useState<CertType[]>([]);
  const [issuers, setIssuers] = useState<Issuer[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const [formData, setFormData] = useState({
    certTypeId: '',
    issuerId: '',
    certNumber: '',
    issuedDate: '',
    remarks: ''
  });

  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const { data: { session } } = await supabase.auth.getSession();
        
        if (!session) {
          setError('Not authenticated');
          setLoading(false);
          return;
        }

        const headers = {
          'Authorization': `Bearer ${session.access_token}`,
        };

        const [certTypesRes, issuersRes] = await Promise.all([
          fetch(`${API_BASE_URL}/api/cert-types`, { headers }),
          fetch(`${API_BASE_URL}/api/issuers`, { headers })
        ]);

        if (!certTypesRes.ok || !issuersRes.ok) {
          throw new Error('Failed to fetch required data');
        }

        const [certTypesData, issuersData] = await Promise.all([
          certTypesRes.json(),
          issuersRes.json()
        ]);

        setCertTypes(certTypesData);
        setIssuers(issuersData);

        // If we came back from Add Issuer with a new issuer ID, select it automatically
        if (location.state?.newIssuerId) {
          setFormData(prev => ({ ...prev, issuerId: location.state.newIssuerId }));
          // Clear the state so it doesn't persist if they navigate away and back
          navigate(location.pathname, { replace: true, state: {} });
        }
      } catch (err: any) {
        setError(err.message || 'An error occurred while fetching data');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
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

      const response = await fetch(`${API_BASE_URL}/api/certificates`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${session.access_token}`,
        },
        body: JSON.stringify({
          'cert-type-id': formData.certTypeId,
          'issuer-id': formData.issuerId,
          'cert-number': formData.certNumber,
          'issued-date': formData.issuedDate,
          remarks: formData.remarks
        }),
      });

      const responseData = await response.json().catch(() => ({}));
      console.log('Add response:', responseData);

      if (!response.ok) {
        const errorMessage = responseData.message || 
                           responseData.error || 
                           (responseData.errors && typeof responseData.errors === 'object' ? JSON.stringify(responseData.errors) : null) ||
                           'Failed to add certificate';
        throw new Error(errorMessage);
      }

      navigate('/dashboard');
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
        <Paper elevation={0} sx={{ p: 4, border: 1, borderColor: 'divider', borderRadius: 2 }}>
          <Typography variant="h5" gutterBottom sx={{ fontWeight: 600, mb: 3 }}>
            Add New Certificate
          </Typography>

          {error && (
            <Alert severity="error" sx={{ mb: 3 }}>
              {error}
            </Alert>
          )}

          <Box component="form" onSubmit={handleSubmit} noValidate>
            <Grid container spacing={3}>
              <Grid size={{ xs: 12, sm: 6 }}>
                <Autocomplete
                  id="certTypeId"
                  options={certTypes}
                  getOptionLabel={(option) => option.name}
                  autoHighlight
                  autoSelect
                  value={certTypes.find(type => type.id === formData.certTypeId) || null}
                  onChange={(_event, newValue) => {
                    setFormData(prev => ({ ...prev, certTypeId: newValue ? newValue.id : '' }));
                  }}
                  renderInput={(params) => (
                    <TextField
                      {...params}
                      label="Certificate Type"
                      required
                      error={!formData.certTypeId && submitting}
                    />
                  )}
                />
              </Grid>

              <Grid size={{ xs: 12, sm: 6 }}>
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
                    <Link component={RouterLink} to="/add-issuer" sx={{ textDecoration: 'none' }}>
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
                  onClick={() => navigate('/dashboard')}
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
                  {submitting ? 'Saving...' : 'Add Certificate'}
                </Button>
              </Grid>
            </Grid>
          </Box>
        </Paper>
      </Box>
    </Container>
  );
};

export default AddCertificate;
