import { useState, useEffect } from 'react';
import { 
  Typography, 
  Container, 
  Box, 
  TextField, 
  Button, 
  Paper, 
  Alert, 
  LinearProgress,
  Snackbar,
  Dialog,
  DialogContent,
  DialogTitle,
  IconButton,
  Chip,
  Grid,
  Autocomplete,
  CircularProgress,
  Link
} from '@mui/material';
import { useNavigate, useParams, Link as RouterLink, useLocation } from 'react-router-dom';
import DeleteIcon from '@mui/icons-material/Delete';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';
import VisibilityIcon from '@mui/icons-material/Visibility';
import CloseIcon from '@mui/icons-material/Close';
import { supabase } from '../supabaseClient';
import { API_BASE_URL } from '../config';
import { calculateExpiryDate, formatDate } from '../utils/dateUtils';
import { useFileUpload } from '../hooks/useFileUpload';

interface Issuer {
  id: string;
  name: string;
}

const UpdateCertificate = () => {
  const { id } = useParams<{ id: string }>();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isPdfSizeError, setIsPdfSizeError] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [documentPath, setDocumentPath] = useState<string | null>(null);
  const [documentUrl, setDocumentUrl] = useState<string | null>(null);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [showSuccess, setShowSuccess] = useState(false);
  const [previewOpen, setPreviewOpen] = useState(false);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [fileType, setFileType] = useState<string | null>(null);
  const [issuers, setIssuers] = useState<Issuer[]>([]);

  const { uploadFile, uploading: uploadingFile, progress, error: uploadError, isPdfSizeError: uploadPdfError } = useFileUpload();

  useEffect(() => {
    if (uploadError) {
      setError(uploadError);
      setIsPdfSizeError(uploadPdfError);
    }
  }, [uploadError, uploadPdfError]);

  useEffect(() => {
    setUploadProgress(progress);
  }, [progress]);

  const [formData, setFormData] = useState({
    certNumber: '',
    issuedDate: '',
    remarks: '',
    issuerId: '',
    manualExpiry: ''
  });
  const [showManualExpiry, setShowManualExpiry] = useState(false);
  const [certTypeName, setCertTypeName] = useState('');
  const [validityMonths, setValidityMonths] = useState<number | null>(null);

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
          issuerId: cert['issuer-id'] || '',
          manualExpiry: cert['manual-expiry'] ? cert['manual-expiry'].split('T')[0] : ''
        });
        if (cert['manual-expiry']) {
          setShowManualExpiry(true);
        }
        setCertTypeName(cert['cert-type-name'] || '');
        setValidityMonths(cert['cert-type-normal-validity-months'] ?? null);
        setDocumentPath(cert['document-path'] || null);
        setDocumentUrl(cert['document-url'] || null);
        if (cert['document-path']?.toLowerCase().endsWith('.pdf')) {
          setFileType('application/pdf');
        } else if (cert['document-path']?.toLowerCase().endsWith('.jpg') || cert['document-path']?.toLowerCase().endsWith('.jpeg')) {
          setFileType('image/jpeg');
        }
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

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setError(null);
    setIsPdfSizeError(false);
    setFileType(file.type);

    try {
      const { fileKey } = await uploadFile(file);
      setDocumentPath(fileKey);
      setDocumentUrl(null); // Clear old preview if any
      
      // Create local preview URL
      const url = URL.createObjectURL(file);
      setPreviewUrl(url);
    } catch (err: any) {
      // Error handled by useFileUpload
      setPreviewUrl(null);
      setFileType(null);
    }
  };

  const handleRemoveFile = () => {
    setDocumentPath(null);
    setDocumentUrl(null);
    setPreviewUrl(null);
    setFileType(null);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setError(null);
    setIsPdfSizeError(false);

    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) throw new Error('Not authenticated');

      const payload = {
        id: id,
        'issuer-id': formData.issuerId,
        'cert-number': formData.certNumber,
        'issued-date': formData.issuedDate ? new Date(formData.issuedDate).toISOString() : null,
        'manual-expiry': formData.manualExpiry ? new Date(formData.manualExpiry).toISOString() : null,
        remarks: formData.remarks,
        'document-path': documentPath
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

      // setShowSuccess(true);
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
              {isPdfSizeError && (
                <Box component="span" sx={{ display: 'block', mt: 1 }}>
                  Try using <Link href="https://www.ilovepdf.com/compress_pdf" target="_blank" rel="noopener">iLovePDF</Link> to reduce your file size.
                </Box>
              )}
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
                {(() => {
                  if (!formData.issuedDate) return null;

                  const expiryDate = calculateExpiryDate(formData.issuedDate, validityMonths);
                  return (
                    <Box sx={{ mt: 1 }}>
                      <Typography variant="caption" color="primary" sx={{ display: 'inline', fontWeight: 500 }}>
                        {validityMonths === null || validityMonths === 0 ? 'Does not expire' : (expiryDate ? `Expiry Date: ${formatDate(expiryDate)}` : 'Does not expire')}
                      </Typography>
                      {!showManualExpiry && (
                        <Link
                          component="button"
                          type="button"
                          variant="caption"
                          onClick={() => setShowManualExpiry(true)}
                          sx={{ ml: 1, textDecoration: 'none', verticalAlign: 'baseline' }}
                        >
                          (Set manual expiry)
                        </Link>
                      )}
                      {showManualExpiry && (
                        <Box sx={{ mt: 2 }}>
                          <TextField
                            fullWidth
                            id="manualExpiry"
                            label="Manual Expiry Date"
                            name="manualExpiry"
                            type="date"
                            value={formData.manualExpiry}
                            onChange={handleChange}
                            InputLabelProps={{ shrink: true }}
                            helperText="Override the calculated expiry date"
                          />
                          <Link
                            component="button"
                            type="button"
                            variant="caption"
                            color="error"
                            onClick={() => {
                              setShowManualExpiry(false);
                              setFormData(prev => ({ ...prev, manualExpiry: '' }));
                            }}
                            sx={{ mt: 0.5, textDecoration: 'none' }}
                          >
                            Remove manual expiry
                          </Link>
                        </Box>
                      )}
                    </Box>
                  );
                })()}
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

              <Grid size={{ xs: 12 }}>
                <Box sx={{ mt: 1 }}>
                  <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 500 }}>
                    Certificate Attachment (PDF/JPG)
                  </Typography>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, flexWrap: 'wrap' }}>
                    <Button
                      variant="outlined"
                      component="label"
                      startIcon={<CloudUploadIcon />}
                      disabled={uploadingFile || submitting}
                    >
                      {documentPath ? 'Change File' : 'Upload File'}
                      <input
                        type="file"
                        hidden
                        accept="application/pdf,image/jpeg,image/jpg"
                        onChange={handleFileUpload}
                      />
                    </Button>
                    
                    {documentPath && !uploadingFile && (
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                        <Chip
                          label={documentPath.split('/').pop()}
                          onDelete={handleRemoveFile}
                          deleteIcon={<DeleteIcon />}
                          color="primary"
                          variant="outlined"
                        />
                        {(documentUrl || previewUrl) && (
                          <IconButton 
                            color="primary" 
                            onClick={() => setPreviewOpen(true)}
                            size="small"
                            title="Preview document"
                          >
                            <VisibilityIcon />
                          </IconButton>
                        )}
                      </Box>
                    )}
                  </Box>

                  {uploadingFile && (
                    <Box sx={{ width: '100%', mt: 1, mb: 1 }}>
                      <LinearProgress variant="determinate" value={uploadProgress} />
                      <Typography variant="caption" color="text.secondary">
                        Uploading: {uploadProgress}%
                      </Typography>
                    </Box>
                  )}

                  <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mt: 1 }}>
                    Max size: 3MB. Allowed formats: PDF, JPG.
                  </Typography>
                </Box>
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

      <Dialog
        open={previewOpen}
        onClose={() => setPreviewOpen(false)}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle sx={{ m: 0, p: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          Document Preview
          <IconButton onClick={() => setPreviewOpen(false)}>
            <CloseIcon />
          </IconButton>
        </DialogTitle>
        <DialogContent dividers sx={{ p: 0, height: '70vh', overflow: 'hidden' }}>
          {(previewUrl || documentUrl) && (
            (fileType === 'application/pdf' || documentPath?.toLowerCase().endsWith('.pdf')) ? (
              <iframe
                src={previewUrl || documentUrl || ''}
                title="PDF Preview"
                width="100%"
                height="100%"
                style={{ border: 'none' }}
              />
            ) : (
              <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100%', p: 2 }}>
                <img
                  src={previewUrl || documentUrl || ''}
                  alt="Document Preview"
                  style={{ maxWidth: '100%', maxHeight: '100%', objectFit: 'contain' }}
                />
              </Box>
            )
          )}
        </DialogContent>
      </Dialog>

      <Snackbar
        open={showSuccess}
        autoHideDuration={3000}
        onClose={() => setShowSuccess(false)}
        message="Certificate updated successfully!"
      />
    </Container>
  );
};

export default UpdateCertificate;
