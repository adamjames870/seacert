import { useEffect, useState } from 'react';
import { 
  Typography, 
  Container, 
  Box, 
  Paper, 
  Grid, 
  Button, 
  CircularProgress, 
  Alert, 
  Divider, 
  Link,
  Card,
  CardContent,
  Stack,
  Chip
} from '@mui/material';
import { useParams, Link as RouterLink } from 'react-router-dom';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import FilePresentIcon from '@mui/icons-material/FilePresent';
import VisibilityIcon from '@mui/icons-material/Visibility';
import { supabase } from '../supabaseClient';
import { API_BASE_URL } from '../config';
import { formatDate } from '../utils/dateUtils';
import { getCountryName } from '../utils/countryData';

interface Predecessor {
  reason: string;
  certificate: Certificate;
}

interface Certificate {
  id: string;
  'created-at': string;
  'updated-at': string;
  'cert-type-id': string;
  'cert-type-name': string;
  'cert-type-short-name': string;
  'cert-type-stcw-ref': string;
  'cert-type-normal-validity-months': number;
  'cert-number': string;
  'issuer-id': string;
  'issuer-name': string;
  'issuer-country': string;
  'issuer-website': string;
  'issued-date': string;
  'expiry-date': string;
  'alternative-name': string;
  remarks: string;
  deleted: boolean;
  'has-successors': boolean;
  'manual-expiry'?: string;
  'document-path'?: string;
  'document-url'?: string;
  predecessors?: Predecessor[];
}

const CertificateDetails = () => {
  const { id } = useParams<{ id: string }>();
  const [certificate, setCertificate] = useState<Certificate | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showPreview, setShowPreview] = useState(false);

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

        // The issue description says: GET api/certificates?id=<uuid>
        const response = await fetch(`${API_BASE_URL}/api/certificates?id=${id}`, {
          headers: {
            'Authorization': `Bearer ${session.access_token}`,
          },
        });

        if (!response.ok) {
          throw new Error('Failed to fetch certificate details');
        }

        const data = await response.json();
        
        // Handle both single object or array depending on API behavior
        let cert: Certificate | null = null;
        if (Array.isArray(data)) {
            cert = data.find((c: any) => c.id === id) || null;
        } else if (data.certificates && Array.isArray(data.certificates)) {
            cert = data.certificates.find((c: any) => c.id === id) || null;
        } else if (data && data.id === id) {
            cert = data;
        }

        if (!cert) {
          throw new Error('Certificate not found');
        }

        setCertificate(cert);
      } catch (err: any) {
        setError(err.message || 'An error occurred while fetching data');
      } finally {
        setLoading(false);
      }
    };

    if (id) {
      fetchCertificate();
    }
  }, [id]);

  if (loading) {
    return (
      <Container sx={{ mt: 4, display: 'flex', justifyContent: 'center' }}>
        <CircularProgress />
      </Container>
    );
  }

  if (error || !certificate) {
    return (
      <Container sx={{ mt: 4 }}>
        <Alert severity="error">{error || 'Certificate not found'}</Alert>
        <Button 
          startIcon={<ArrowBackIcon />} 
          component={RouterLink} 
          to="/certificates" 
          sx={{ mt: 2 }}
        >
          Back to certificate list
        </Button>
      </Container>
    );
  }

  const isPdf = certificate['document-url']?.toLowerCase().includes('.pdf') || certificate['document-path']?.toLowerCase().endsWith('.pdf');
  const isImage = certificate['document-url']?.toLowerCase().match(/\.(jpg|jpeg|png|gif)$/) || certificate['document-path']?.toLowerCase().match(/\.(jpg|jpeg|png|gif)$/);

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Box sx={{ mb: 3, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h4" component="h1">
          Certificate Details
        </Typography>
        <Button 
          startIcon={<ArrowBackIcon />} 
          component={RouterLink} 
          to="/certificates"
          variant="outlined"
        >
          Back to List
        </Button>
      </Box>

      <Grid container spacing={3}>
        <Grid item xs={12} md={8}>
          <Paper sx={{ p: 3, mb: 3 }}>
            <Typography variant="h6" gutterBottom color="primary">
              {certificate['cert-type-name']}
            </Typography>
            <Typography variant="body2" color="text.secondary" gutterBottom>
              {certificate['cert-type-short-name']} {certificate['cert-type-stcw-ref'] && `• ${certificate['cert-type-stcw-ref']}`}
            </Typography>
            
            <Divider sx={{ my: 2 }} />
            
            <Grid container spacing={2}>
              <Grid item xs={12} sm={6}>
                <Typography variant="caption" color="text.secondary">Certificate Number</Typography>
                <Typography variant="body1" sx={{ fontWeight: 500 }}>{certificate['cert-number'] || 'N/A'}</Typography>
              </Grid>
              <Grid item xs={12} sm={6}>
                <Typography variant="caption" color="text.secondary">Alternative Name</Typography>
                <Typography variant="body1">{certificate['alternative-name'] || 'None'}</Typography>
              </Grid>
              
              <Grid item xs={12} sm={6}>
                <Typography variant="caption" color="text.secondary">Issued Date</Typography>
                <Typography variant="body1">{formatDate(certificate['issued-date'])}</Typography>
              </Grid>
              <Grid item xs={12} sm={6}>
                <Typography variant="caption" color="text.secondary">Expiry Date</Typography>
                <Typography variant="body1" sx={{ 
                  fontWeight: 500,
                  color: new Date(certificate['expiry-date']) < new Date() ? 'error.main' : 'inherit'
                }}>
                  {formatDate(certificate['expiry-date'])}
                  {certificate['manual-expiry'] && ' (Manual)'}
                </Typography>
              </Grid>

              <Grid item xs={12}>
                <Typography variant="caption" color="text.secondary">Remarks</Typography>
                <Typography variant="body1" sx={{ whiteSpace: 'pre-wrap' }}>{certificate.remarks || 'No remarks'}</Typography>
              </Grid>
            </Grid>
          </Paper>

          <Paper sx={{ p: 3 }}>
            <Typography variant="h6" gutterBottom>
              Issuer Information
            </Typography>
            <Grid container spacing={2}>
              <Grid item xs={12} sm={6}>
                <Typography variant="caption" color="text.secondary">Name</Typography>
                <Typography variant="body1">{certificate['issuer-name']}</Typography>
              </Grid>
              <Grid item xs={12} sm={6}>
                <Typography variant="caption" color="text.secondary">Country</Typography>
                <Typography variant="body1">{getCountryName(certificate['issuer-country'])}</Typography>
              </Grid>
              {certificate['issuer-website'] && (
                <Grid item xs={12}>
                  <Typography variant="caption" color="text.secondary">Website</Typography>
                  <Box>
                    <Link href={certificate['issuer-website']} target="_blank" rel="noopener noreferrer">
                      {certificate['issuer-website']}
                    </Link>
                  </Box>
                </Grid>
              )}
            </Grid>
          </Paper>
          
          {certificate.predecessors && certificate.predecessors.length > 0 && (
            <Paper sx={{ p: 3, mt: 3 }}>
              <Typography variant="h6" gutterBottom>
                Predecessors
              </Typography>
              <Stack spacing={1}>
                {certificate.predecessors.map((p, index) => (
                  <Card key={index} variant="outlined">
                    <CardContent sx={{ py: 1, '&:last-child': { pb: 1 } }}>
                      <Typography variant="body2">
                        <strong>{p.certificate['cert-type-name']}</strong> ({p.certificate['cert-number']})
                      </Typography>
                      <Typography variant="caption" display="block">
                        Reason: {p.reason}
                      </Typography>
                    </CardContent>
                  </Card>
                ))}
              </Stack>
            </Paper>
          )}
        </Grid>

        <Grid item xs={12} md={4}>
          <Paper sx={{ p: 3, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
            <Typography variant="h6" gutterBottom alignSelf="flex-start">
              Document
            </Typography>
            
            {certificate['document-url'] ? (
              <>
                <FilePresentIcon sx={{ fontSize: 60, color: 'primary.main', mb: 2 }} />
                <Typography variant="body2" gutterBottom align="center">
                  Certificate file is attached
                </Typography>
                <Stack spacing={2} sx={{ width: '100%', mt: 1 }}>
                  <Button 
                    variant="contained" 
                    startIcon={<VisibilityIcon />}
                    onClick={() => setShowPreview(!showPreview)}
                  >
                    {showPreview ? 'Hide Preview' : 'Show Preview'}
                  </Button>
                  <Button 
                    variant="outlined" 
                    component={Link} 
                    href={certificate['document-url']} 
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    Download / Open in New Tab
                  </Button>
                </Stack>
              </>
            ) : (
              <Box sx={{ py: 4, textAlign: 'center' }}>
                <Typography variant="body2" color="text.secondary">
                  No document attached
                </Typography>
              </Box>
            )}
          </Paper>

          {certificate['document-url'] && showPreview && (
            <Paper sx={{ p: 1, mt: 2, height: 500, overflow: 'hidden' }}>
              {isPdf ? (
                <iframe
                  src={`${certificate['document-url']}#toolbar=0`}
                  width="100%"
                  height="100%"
                  style={{ border: 'none' }}
                  title="Certificate Preview"
                />
              ) : isImage ? (
                <Box 
                  component="img"
                  src={certificate['document-url']}
                  alt="Certificate Preview"
                  sx={{ width: '100%', height: '100%', objectFit: 'contain' }}
                />
              ) : (
                <Box sx={{ p: 4, textAlign: 'center' }}>
                  <Typography>Preview not available for this file type. Please download to view.</Typography>
                </Box>
              )}
            </Paper>
          )}
          
          <Paper sx={{ p: 3, mt: 3 }}>
            <Typography variant="h6" gutterBottom>
              Status
            </Typography>
            <Stack spacing={1}>
              <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                <Typography variant="body2">Active</Typography>
                <Chip 
                  size="small" 
                  label={certificate.deleted ? 'Archived' : 'Active'} 
                  color={certificate.deleted ? 'default' : 'success'} 
                />
              </Box>
              <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                <Typography variant="body2">Has Successors</Typography>
                <Typography variant="body2">{certificate['has-successors'] ? 'Yes' : 'No'}</Typography>
              </Box>
              <Divider sx={{ my: 1 }} />
              <Typography variant="caption" color="text.secondary">
                Added: {formatDate(certificate['created-at'])}
              </Typography>
              <Typography variant="caption" color="text.secondary">
                Last Updated: {formatDate(certificate['updated-at'])}
              </Typography>
            </Stack>
          </Paper>
        </Grid>
      </Grid>
    </Container>
  );
};

export default CertificateDetails;
