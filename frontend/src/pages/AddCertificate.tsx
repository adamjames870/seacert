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
  LinearProgress,
  Snackbar,
  Dialog,
  DialogContent,
  DialogTitle,
  IconButton,
  Grid,
  Autocomplete,
  Chip,
  CircularProgress
} from '@mui/material';
import { useNavigate, Link as RouterLink, useLocation } from 'react-router-dom';
import DeleteIcon from '@mui/icons-material/Delete';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';
import CloseIcon from '@mui/icons-material/Close';
import VisibilityIcon from '@mui/icons-material/Visibility';
import { supabase } from '../supabaseClient';
import { API_BASE_URL } from '../config';
import { calculateExpiryDate, formatDate } from '../utils/dateUtils';
import { useFileUpload } from '../hooks/useFileUpload';

interface CertType {
  id: string;
  name: string;
  'short-name': string;
  'stcw-reference': string;
  'normal-validity-months'?: number;
  status?: 'approved' | 'provisional';
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
  const [isPdfSizeError, setIsPdfSizeError] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [documentPath, setDocumentPath] = useState<string | null>(null);
  const [fileName, setFileName] = useState<string | null>(null);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [showSuccess, setShowSuccess] = useState(false);
  const [previewOpen, setPreviewOpen] = useState(false);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [fileType, setFileType] = useState<string | null>(null);
  const [isDragging, setIsDragging] = useState(false);

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
    certTypeId: '',
    issuerId: '',
    certNumber: '',
    issuedDate: new Date().toISOString().split('T')[0],
    remarks: '',
    supersedes: '',
    supersedeReason: '',
    manualExpiry: ''
  });
  const [showManualExpiry, setShowManualExpiry] = useState(false);

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

        // Sort cert types alphabetically by name
        const sortedCertTypes = (certTypesData as CertType[]).sort((a, b) => 
          a.name.localeCompare(b.name)
        );

        // Sort issuers alphabetically by name
        const sortedIssuers = (issuersData as Issuer[]).sort((a, b) => 
          a.name.localeCompare(b.name)
        );

        setCertTypes(sortedCertTypes);
        setIssuers(sortedIssuers);

        // Pre-fill from location state (Update/Replace from Certificates page)
        if (location.state?.certTypeId || location.state?.supersedes || location.state?.supersedeReason) {
          setFormData(prev => ({
            ...prev,
            certTypeId: location.state.certTypeId || prev.certTypeId,
            supersedes: location.state.supersedes || prev.supersedes,
            supersedeReason: location.state.supersedeReason || prev.supersedeReason
          }));
        }

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

  const handleFileUpload = async (file: File) => {
    if (!file) return;

    setError(null);
    setIsPdfSizeError(false);
    setFileName(file.name);
    setFileType(file.type);

    try {
      const { fileKey } = await uploadFile(file);
      setDocumentPath(fileKey);
      
      // Create local preview URL
      const url = URL.createObjectURL(file);
      setPreviewUrl(url);
    } catch (err: any) {
      // Error is handled by useFileUpload and synced via useEffect
      setFileName(null);
      setPreviewUrl(null);
      setFileType(null);
    }
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
    
    const file = e.dataTransfer.files?.[0];
    if (file && (file.type === 'application/pdf' || file.type === 'image/jpeg' || file.type === 'image/jpg')) {
      handleFileUpload(file);
    } else if (file) {
      setError('Only PDF and JPG files are allowed.');
    }
  };

  const handleRemoveFile = () => {
    setDocumentPath(null);
    setFileName(null);
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
          'issued-date': formData.issuedDate ? new Date(formData.issuedDate).toISOString() : null,
          'manual-expiry': formData.manualExpiry ? new Date(formData.manualExpiry).toISOString() : null,
          remarks: formData.remarks,
          'supersedes': formData.supersedes || undefined,
          'supersede-reason': formData.supersedeReason || undefined,
          'document-path': documentPath
        }),
      });

      const responseData = await response.json().catch(() => ({}));
      console.log('Add response:', responseData);

      if (!response.ok) {
        const errorMessage = responseData.error || 
                           responseData.message || 
                           (responseData.errors && typeof responseData.errors === 'object' ? JSON.stringify(responseData.errors) : null) ||
                           'Failed to add certificate';
        throw new Error(errorMessage);
      }

      // setShowSuccess(true);
      navigate('/certificates');
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
          <Typography variant="h5" gutterBottom sx={{ fontWeight: 600, mb: 1 }}>
            {formData.supersedes 
              ? `${formData.supersedeReason === 'updated' ? 'Update' : 'Replace'} Certificate` 
              : 'Add New Certificate'}
          </Typography>
          {formData.supersedes && (
            <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
              This will {formData.supersedeReason} the existing certificate.
            </Typography>
          )}

          <Box component="form" onSubmit={handleSubmit} noValidate>
            <Grid container spacing={3}>
              <Grid size={{ xs: 12, sm: 6 }}>
                <Autocomplete
                  id="certTypeId"
                  options={certTypes}
                  getOptionLabel={(option) => option.name}
                  filterOptions={(options, { inputValue }) => {
                    const query = inputValue.toLowerCase();
                    return options.filter(option => 
                      option.name.toLowerCase().includes(query) || 
                      option['short-name']?.toLowerCase().includes(query)
                    );
                  }}
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
                  renderOption={(props, option) => {
                    const { key, ...optionProps } = props as any;
                    return (
                      <Box component="li" key={key} {...optionProps} sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', width: '100%' }}>
                        <Box>
                          <Typography variant="body1">{option.name}</Typography>
                          {option['short-name'] && (
                            <Typography variant="caption" color="text.secondary">
                              {option['short-name']}
                            </Typography>
                          )}
                        </Box>
                        {option.status === 'provisional' && (
                          <Chip label="Provisional" size="small" color="warning" variant="outlined" sx={{ ml: 1 }} />
                        )}
                      </Box>
                    );
                  }}
                  noOptionsText={
                    <Box sx={{ p: 1 }}>
                      <Typography variant="body2" sx={{ mb: 1 }}>No results found</Typography>
                      <Button
                        size="small"
                        color="primary"
                        component={RouterLink}
                        to="/add-cert-type"
                        state={{ from: 'add-certificate' }}
                        fullWidth
                        sx={{ justifyContent: 'flex-start' }}
                      >
                        Add certificate type
                      </Button>
                    </Box>
                  }
                />
                <Box sx={{ mt: 1 }}>
                  <Typography variant="caption">
                    Can't find the certificate type?{' '}
                    <Link component={RouterLink} to="/add-cert-type" state={{ from: 'add-certificate' }} sx={{ textDecoration: 'none' }}>
                      Add certificate type
                    </Link>
                  </Typography>
                </Box>
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
                  noOptionsText={
                    <Box sx={{ p: 1 }}>
                      <Typography variant="body2" sx={{ mb: 1 }}>No results found</Typography>
                      <Button
                        size="small"
                        color="primary"
                        component={RouterLink}
                        to="/add-issuer"
                        state={{ from: 'add-certificate' }}
                        fullWidth
                        sx={{ justifyContent: 'flex-start' }}
                      >
                        Add New Issuer
                      </Button>
                    </Box>
                  }
                />
                <Box sx={{ mt: 1 }}>
                  <Typography variant="caption">
                    Can't find your issuer?{' '}
                    <Link component={RouterLink} to="/add-issuer" state={{ from: 'add-certificate' }} sx={{ textDecoration: 'none' }}>
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
                  const selectedCertType = certTypes.find(type => type.id === formData.certTypeId);
                  if (!selectedCertType || !formData.issuedDate) return null;
                  
                  const validityMonths = selectedCertType['normal-validity-months'];
                  const expiryDate = calculateExpiryDate(formData.issuedDate, validityMonths ?? null);
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
                  <Box 
                    onDragOver={handleDragOver}
                    onDragLeave={handleDragLeave}
                    onDrop={handleDrop}
                    sx={{ 
                      border: '2px dashed',
                      borderColor: isDragging ? 'primary.main' : 'divider',
                      borderRadius: 2,
                      p: 3,
                      textAlign: 'center',
                      bgcolor: isDragging ? 'action.hover' : 'background.paper',
                      transition: 'all 0.2s ease-in-out',
                      mb: 1
                    }}
                  >
                    <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 1 }}>
                      <CloudUploadIcon sx={{ fontSize: 40, color: isDragging ? 'primary.main' : 'text.secondary', mb: 1 }} />
                      <Typography variant="body2" color="text.secondary">
                        Drag and drop your certificate here, or
                      </Typography>
                      <Button
                        variant="outlined"
                        component="label"
                        disabled={uploadingFile || submitting}
                        size="small"
                        sx={{ mt: 1 }}
                      >
                        {documentPath ? 'Change File' : 'Select File'}
                        <input
                          type="file"
                          hidden
                          accept="application/pdf,image/jpeg,image/jpg"
                          onChange={(e) => {
                            const file = e.target.files?.[0];
                            if (file) handleFileUpload(file);
                          }}
                        />
                      </Button>
                      <Typography variant="caption" color="text.secondary" sx={{ mt: 1 }}>
                        Maximum size: 3MB (PDF or JPG)
                      </Typography>
                    </Box>
                  </Box>
                  {fileName && !uploadingFile && (
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Chip
                        label={fileName}
                        onDelete={handleRemoveFile}
                        deleteIcon={<DeleteIcon />}
                        color="primary"
                        variant="outlined"
                      />
                      {previewUrl && (
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

              <Grid size={{ xs: 12 }}>
                {error && (
                  <Alert severity="error" sx={{ mb: 2 }}>
                    {error}
                    {isPdfSizeError && (
                      <Box component="span" sx={{ display: 'block', mt: 1 }}>
                        Try using <Link href="https://www.ilovepdf.com/compress_pdf" target="_blank" rel="noopener">iLovePDF</Link> to reduce your file size.
                      </Box>
                    )}
                  </Alert>
                )}
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
                  {submitting ? 'Saving...' : (formData.supersedes ? (formData.supersedeReason === 'updated' ? 'Update Certificate' : 'Replace Certificate') : 'Add Certificate')}
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
          {previewUrl && (
            fileType === 'application/pdf' ? (
              <iframe
                src={previewUrl}
                title="PDF Preview"
                width="100%"
                height="100%"
                style={{ border: 'none' }}
              />
            ) : (
              <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100%', p: 2 }}>
                <img
                  src={previewUrl}
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
        message="Certificate added successfully!"
      />
    </Container>
  );
};

export default AddCertificate;
