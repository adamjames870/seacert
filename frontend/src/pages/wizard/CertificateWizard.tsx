import React, { useState, useEffect } from 'react';
import {
  Typography,
  Container,
  Box,
  Paper,
  Stepper,
  Step,
  StepLabel,
  Button,
  CircularProgress,
  Alert,
  TextField,
  Autocomplete,
  Grid,
  IconButton,
  Dialog,
  DialogContent,
  DialogTitle,
  DialogActions,
  Link
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';
import CloseIcon from '@mui/icons-material/Close';
import VisibilityIcon from '@mui/icons-material/Visibility';
import { supabase } from '../../supabaseClient';
import { API_BASE_URL } from '../../config';
import { useFileUpload } from '../../hooks/useFileUpload';
import { countries } from '../../utils/countryData';

interface CertType {
  id: string;
  name: string;
  'short-name': string;
  'stcw-reference': string;
  'normal-validity-months'?: number;
}

interface Issuer {
  id: string;
  name: string;
  country: string;
  website: string;
}

interface ExtractedData {
  'cert-type-name': string;
  'cert-number': string;
  'issuer-name': string;
  'issued-date': string;
  'cert-type-id'?: string;
  'issuer-id'?: string;
}

const steps = [
  'Upload Certificate',
  'Certificate Type',
  'Issuer Information',
  'Confirm Details',
  'Final Review'
];

const CertificateWizard = () => {
  const [activeStep, setActiveStep] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  // Data from API
  const [certTypes, setCertTypes] = useState<CertType[]>([]);
  const [issuers, setIssuers] = useState<Issuer[]>([]);
  
  // Wizard State
  const [file, setFile] = useState<File | null>(null);
  const [documentPath, setDocumentPath] = useState<string | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [previewOpen, setPreviewOpen] = useState(false);
  const [isDragging, setIsDragging] = useState(false);
  
  const [extractedData, setExtractedData] = useState<ExtractedData | null>(null);
  
  // Form State for each step
  const [certDetails, setCertDetails] = useState({
    certNumber: '',
    issuedDate: ''
  });
  
  const [selectedIssuer, setSelectedIssuer] = useState<Issuer | null>(null);
  const [newIssuer, setNewIssuer] = useState({ name: '', country: '', website: '' });
  const [isCreatingNewIssuer, setIsCreatingNewIssuer] = useState(false);
  
  const [selectedCertType, setSelectedCertType] = useState<CertType | null>(null);
  const [newCertType, setNewCertType] = useState({ 
    name: '', 
    shortName: '', 
    stcwReference: '', 
    normalValidityMonths: '',
    status: 'provisional' 
  });
  const [isCreatingNewCertType, setIsCreatingNewCertType] = useState(false);
  const [privacyAgreed, setPrivacyAgreed] = useState(false);

  const navigate = useNavigate();
  const { uploadFile, uploading: uploadingFile } = useFileUpload();

  useEffect(() => {
    fetchInitialData();
  }, []);

  const fetchInitialData = async () => {
    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) return;

      const headers = { 'Authorization': `Bearer ${session.access_token}` };
      const [certTypesRes, issuersRes] = await Promise.all([
        fetch(`${API_BASE_URL}/api/cert-types`, { headers }),
        fetch(`${API_BASE_URL}/api/issuers`, { headers })
      ]);

      if (certTypesRes.ok && issuersRes.ok) {
        const [certTypesData, issuersData] = await Promise.all([
          certTypesRes.json(),
          issuersRes.json()
        ]);
        setCertTypes((certTypesData as CertType[]).sort((a, b) => a.name.localeCompare(b.name)));
        setIssuers((issuersData as Issuer[]).sort((a, b) => a.name.localeCompare(b.name)));
      }
    } catch (err) {
      console.error('Error fetching data:', err);
    }
  };

  const handleFileUpload = async (selectedFile: File) => {
    setFile(selectedFile);
    setPreviewUrl(URL.createObjectURL(selectedFile));
    setError(null);
  };

  const handleExtract = async () => {
    if (!file) return;
    
    if (!privacyAgreed) return;

    setLoading(true);
    setError(null);

    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) throw new Error('Not authenticated');

      // 1. Upload the file first to get the document path (used for final submission)
      const { fileKey } = await uploadFile(file);
      setDocumentPath(fileKey);

      // 2. Call the extraction API
      const formData = new FormData();
      formData.append('certificate', file);

      const response = await fetch(`${API_BASE_URL}/api/certificates/extract`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${session.access_token}`,
        },
        body: formData,
      });

      if (!response.ok) {
        throw new Error('Failed to extract data from certificate');
      }

      const data: ExtractedData = await response.json();
      setExtractedData(data);
      
      // Initialize form states with extracted data
      setCertDetails({
        certNumber: data['cert-number'] || '',
        issuedDate: data['issued-date'] || ''
      });

      // Handle Issuer
      if (data['issuer-id']) {
        const issuer = issuers.find(i => i.id === data['issuer-id']);
        if (issuer) {
          setSelectedIssuer(issuer);
          setIsCreatingNewIssuer(false);
        } else {
          // If ID provided but not in our list (rare), fallback to name
          setNewIssuer({ ...newIssuer, name: data['issuer-name'] || '' });
          setIsCreatingNewIssuer(true);
        }
      } else {
        setNewIssuer({ ...newIssuer, name: data['issuer-name'] || '' });
        setIsCreatingNewIssuer(true);
      }

      // Handle Cert Type
      if (data['cert-type-id']) {
        const certType = certTypes.find(ct => ct.id === data['cert-type-id']);
        if (certType) {
          setSelectedCertType(certType);
          setIsCreatingNewCertType(false);
        } else {
          setNewCertType({ ...newCertType, name: data['cert-type-name'] || '' });
          setIsCreatingNewCertType(true);
        }
      } else {
        setNewCertType({ ...newCertType, name: data['cert-type-name'] || '' });
        setIsCreatingNewCertType(true);
      }

      setActiveStep(1);
    } catch (err: any) {
      setError(err.message || 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  const handleNext = () => {
    setActiveStep((prev) => prev + 1);
  };

  const handleBack = () => {
    setActiveStep((prev) => prev - 1);
  };

  const handleSubmit = async () => {
    setLoading(true);
    setError(null);

    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) throw new Error('Not authenticated');

      let finalIssuerId = selectedIssuer?.id;
      let finalCertTypeId = selectedCertType?.id;

      // Create Issuer if needed
      if (isCreatingNewIssuer) {
        const issuerRes = await fetch(`${API_BASE_URL}/api/issuers`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${session.access_token}`,
          },
          body: JSON.stringify({
            name: newIssuer.name,
            country: newIssuer.country || null,
            website: newIssuer.website || null
          }),
        });
        if (!issuerRes.ok) throw new Error('Failed to create issuer');
        const issuerData = await issuerRes.json();
        finalIssuerId = issuerData.id;
      }

      // Create Cert Type if needed
      if (isCreatingNewCertType) {
        const certTypeRes = await fetch(`${API_BASE_URL}/api/cert-types`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${session.access_token}`,
          },
          body: JSON.stringify({
            name: newCertType.name,
            'short-name': newCertType.shortName || null,
            'stcw-reference': newCertType.stcwReference || null,
            'normal-validity-months': newCertType.normalValidityMonths ? parseInt(newCertType.normalValidityMonths) : 0,
            status: newCertType.status
          }),
        });
        if (!certTypeRes.ok) throw new Error('Failed to create certificate type');
        const certTypeData = await certTypeRes.json();
        finalCertTypeId = certTypeData.id;
      }

      // Finally create the certificate
      const response = await fetch(`${API_BASE_URL}/api/certificates`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${session.access_token}`,
        },
        body: JSON.stringify({
          'cert-type-id': finalCertTypeId,
          'issuer-id': finalIssuerId,
          'cert-number': certDetails.certNumber,
          'issued-date': certDetails.issuedDate,
          'document-path': documentPath
        }),
      });

      if (!response.ok) throw new Error('Failed to save certificate');

      navigate('/certificates', { state: { showSuccess: true } });
    } catch (err: any) {
      setError(err.message || 'An error occurred during submission');
    } finally {
      setLoading(false);
    }
  };

  const renderStepContent = (step: number) => {
    switch (step) {
      case 0: // Upload
        return (
          <Box sx={{ mt: 4, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
            <Paper
              variant="outlined"
              onDragOver={(e) => { e.preventDefault(); if (!file) setIsDragging(true); }}
              onDragLeave={() => setIsDragging(false)}
              onDrop={(e) => {
                e.preventDefault();
                setIsDragging(false);
                if (!file) {
                  const droppedFile = e.dataTransfer.files[0];
                  if (droppedFile) handleFileUpload(droppedFile);
                }
              }}
              sx={{
                p: 4,
                width: '100%',
                textAlign: 'center',
                backgroundColor: isDragging ? 'action.hover' : 'background.paper',
                borderStyle: 'dashed',
                cursor: !file ? 'pointer' : 'default',
                position: 'relative',
                minHeight: 200,
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'center',
                alignItems: 'center',
                overflow: 'hidden'
              }}
              onClick={() => !file && document.getElementById('file-upload')?.click()}
            >
              <input
                type="file"
                id="file-upload"
                hidden
                accept="application/pdf,image/jpeg,image/jpg"
                onChange={(e) => e.target.files?.[0] && handleFileUpload(e.target.files[0])}
              />
              
              {!file ? (
                <>
                  <CloudUploadIcon sx={{ fontSize: 48, color: 'primary.main', mb: 2 }} />
                  <Typography variant="h6">
                    Drag and drop or click to upload certificate
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Supports PDF and JPG (Max 3MB)
                  </Typography>
                </>
              ) : (
                <>
                  {/* File Info */}
                  <Box sx={{ zIndex: 1 }}>
                    <Typography variant="h6" sx={{ mb: 1 }}>
                      {file.name}
                    </Typography>
                    <Box sx={{ display: 'flex', gap: 1, justifyContent: 'center' }}>
                      <Button 
                        size="small"
                        startIcon={<VisibilityIcon />} 
                        onClick={(e) => { e.stopPropagation(); setPreviewOpen(true); }}
                      >
                        Preview
                      </Button>
                      <Button 
                        size="small"
                        startIcon={<CloseIcon />} 
                        color="error" 
                        onClick={(e) => { 
                          e.stopPropagation(); 
                          setFile(null); 
                          setPreviewUrl(null); 
                          setPrivacyAgreed(false);
                        }}
                      >
                        Remove
                      </Button>
                    </Box>
                  </Box>

                  {/* Privacy Overlay if not agreed */}
                  {!privacyAgreed && (
                    <Box 
                      sx={{ 
                        position: 'absolute', 
                        top: 0, 
                        left: 0, 
                        right: 0, 
                        bottom: 0, 
                        bgcolor: 'background.paper',
                        zIndex: 2,
                        p: 3,
                        display: 'flex',
                        flexDirection: 'column',
                        justifyContent: 'center',
                        alignItems: 'center'
                      }}
                    >
                      <Typography variant="subtitle2" sx={{ display: 'flex', alignItems: 'center', mb: 1, fontWeight: 600 }}>
                        🔒 Privacy Notice
                      </Typography>
                      <Typography variant="caption" color="text.secondary" align="center" sx={{ mb: 1, maxWidth: '90%' }}>
                        To process your request, the file will be sent securely to Google’s Gemini AI. 
                        Google may use this data to improve its models.
                      </Typography>
                      <Typography variant="caption" color="warning.main" align="center" sx={{ fontWeight: 600, mb: 1 }}>
                        ⚠ Do not include sensitive personal information.
                      </Typography>
                      <Link 
                        href="https://policies.google.com/privacy" 
                        target="_blank" 
                        rel="noopener" 
                        variant="caption"
                        sx={{ mb: 2 }}
                      >
                        View Google Privacy Policy
                      </Link>
                      <Box sx={{ display: 'flex', gap: 2 }}>
                        <Button 
                          variant="outlined" 
                          size="small" 
                          onClick={(e) => { 
                            e.stopPropagation(); 
                            setFile(null); 
                            setPreviewUrl(null); 
                          }}
                        >
                          Cancel
                        </Button>
                        <Button 
                          variant="contained" 
                          size="small" 
                          onClick={(e) => { e.stopPropagation(); setPrivacyAgreed(true); }}
                        >
                          I Agree
                        </Button>
                      </Box>
                    </Box>
                  )}
                </>
              )}
            </Paper>

            <Button
              variant="contained"
              disabled={!file || !privacyAgreed || loading || uploadingFile}
              onClick={handleExtract}
              sx={{ mt: 4 }}
              size="large"
            >
              {loading ? <CircularProgress size={24} /> : 'Process Certificate'}
            </Button>
          </Box>
        );
      case 1: // Cert Type
        return (
          <Box sx={{ mt: 2 }}>
            <Typography variant="subtitle1" gutterBottom>
              Extracted Certificate Type: <strong>{extractedData?.['cert-type-name']}</strong>
            </Typography>
            <Box sx={{ mb: 3 }}>
              <Button 
                variant={isCreatingNewCertType ? "outlined" : "contained"} 
                onClick={() => setIsCreatingNewCertType(false)}
                sx={{ mr: 2 }}
              >
                Select Existing
              </Button>
              <Button 
                variant={isCreatingNewCertType ? "contained" : "outlined"} 
                onClick={() => setIsCreatingNewCertType(true)}
              >
                Create New
              </Button>
            </Box>

            {isCreatingNewCertType ? (
              <Grid container spacing={3}>
                <Grid item xs={12}>
                  <TextField
                    fullWidth
                    required
                    label="Name"
                    value={newCertType.name}
                    onChange={(e) => setNewCertType({ ...newCertType, name: e.target.value })}
                    autoFocus
                  />
                </Grid>
                <Grid item xs={12} sm={6}>
                  <TextField
                    fullWidth
                    label="Short Name"
                    value={newCertType.shortName}
                    onChange={(e) => setNewCertType({ ...newCertType, shortName: e.target.value })}
                  />
                </Grid>
                <Grid item xs={12} sm={6}>
                  <TextField
                    fullWidth
                    label="STCW Reference"
                    value={newCertType.stcwReference}
                    onChange={(e) => setNewCertType({ ...newCertType, stcwReference: e.target.value })}
                  />
                </Grid>
                <Grid item xs={12} sm={6}>
                  <TextField
                    fullWidth
                    type="number"
                    label="Normal Validity (Months)"
                    value={newCertType.normalValidityMonths}
                    onChange={(e) => setNewCertType({ ...newCertType, normalValidityMonths: e.target.value })}
                    helperText="Leave blank if certificate does not expire"
                  />
                </Grid>
              </Grid>
            ) : (
              <Autocomplete
                options={certTypes}
                getOptionLabel={(option) => option.name}
                value={selectedCertType}
                onChange={(_, newValue) => setSelectedCertType(newValue)}
                renderInput={(params) => <TextField {...params} label="Search Certificate Type" fullWidth />}
              />
            )}
          </Box>
        );
      case 2: // Issuer
        return (
          <Box sx={{ mt: 2 }}>
            <Typography variant="subtitle1" gutterBottom>
              Extracted Issuer: <strong>{extractedData?.['issuer-name']}</strong>
            </Typography>
            <Box sx={{ mb: 3 }}>
              <Button 
                variant={isCreatingNewIssuer ? "outlined" : "contained"} 
                onClick={() => setIsCreatingNewIssuer(false)}
                sx={{ mr: 2 }}
              >
                Select Existing
              </Button>
              <Button 
                variant={isCreatingNewIssuer ? "contained" : "outlined"} 
                onClick={() => setIsCreatingNewIssuer(true)}
              >
                Create New
              </Button>
            </Box>

            {isCreatingNewIssuer ? (
              <Grid container spacing={3}>
                <Grid item xs={12}>
                  <TextField
                    fullWidth
                    required
                    label="Issuer Name"
                    value={newIssuer.name}
                    onChange={(e) => setNewIssuer({ ...newIssuer, name: e.target.value })}
                    autoFocus
                  />
                </Grid>
                <Grid item xs={12} sm={6}>
                  <Autocomplete
                    options={countries}
                    autoHighlight
                    getOptionLabel={(option) => option.label}
                    value={countries.find(c => c.code === newIssuer.country) || null}
                    onChange={(_, newValue) => setNewIssuer({ ...newIssuer, country: newValue ? newValue.code : '' })}
                    renderInput={(params) => (
                      <TextField 
                        {...params} 
                        label="Country" 
                        fullWidth 
                        inputProps={{
                          ...params.inputProps,
                          autoComplete: 'new-password',
                        }}
                      />
                    )}
                  />
                </Grid>
                <Grid item xs={12} sm={6}>
                  <TextField
                    fullWidth
                    label="Website"
                    placeholder="https://example.com"
                    value={newIssuer.website}
                    onChange={(e) => setNewIssuer({ ...newIssuer, website: e.target.value })}
                  />
                </Grid>
              </Grid>
            ) : (
              <Autocomplete
                options={issuers}
                getOptionLabel={(option) => option.name}
                value={selectedIssuer}
                onChange={(_, newValue) => setSelectedIssuer(newValue)}
                renderInput={(params) => <TextField {...params} label="Search Issuer" fullWidth />}
              />
            )}
          </Box>
        );
      case 3: // Confirm Details
        return (
          <Grid container spacing={3} sx={{ mt: 2 }}>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Certificate Number"
                value={certDetails.certNumber}
                onChange={(e) => setCertDetails({ ...certDetails, certNumber: e.target.value })}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                type="date"
                label="Issued Date"
                InputLabelProps={{ shrink: true }}
                value={certDetails.issuedDate}
                onChange={(e) => setCertDetails({ ...certDetails, issuedDate: e.target.value })}
              />
            </Grid>
          </Grid>
        );
      case 4: // Final Review
        return (
          <Box sx={{ mt: 2 }}>
            <Typography variant="h6" gutterBottom>Summary</Typography>
            <Grid container spacing={2}>
              <Grid item xs={4}><Typography variant="body2" color="text.secondary">Number:</Typography></Grid>
              <Grid item xs={8}><Typography variant="body1">{certDetails.certNumber}</Typography></Grid>
              
              <Grid item xs={4}><Typography variant="body2" color="text.secondary">Issued Date:</Typography></Grid>
              <Grid item xs={8}><Typography variant="body1">{certDetails.issuedDate}</Typography></Grid>
              
              <Grid item xs={4}><Typography variant="body2" color="text.secondary">Issuer:</Typography></Grid>
              <Grid item xs={8}>
                <Typography variant="body1">
                  {isCreatingNewIssuer ? `${newIssuer.name} (New)` : selectedIssuer?.name}
                </Typography>
              </Grid>
              
              <Grid item xs={4}><Typography variant="body2" color="text.secondary">Type:</Typography></Grid>
              <Grid item xs={8}>
                <Typography variant="body1">
                  {isCreatingNewCertType ? `${newCertType.name} (New)` : selectedCertType?.name}
                </Typography>
              </Grid>
            </Grid>
          </Box>
        );
      default:
        return null;
    }
  };

  return (
    <Container maxWidth="md">
      <Box sx={{ mt: 4, mb: 4 }}>
        <Paper elevation={0} sx={{ p: { xs: 2, sm: 4 }, border: 1, borderColor: 'divider', borderRadius: 2 }}>
          <Typography variant="h5" gutterBottom sx={{ fontWeight: 600, mb: 3 }}>
            Smart Add Certificate
          </Typography>

          <Stepper activeStep={activeStep} alternativeLabel sx={{ mb: 4 }}>
            {steps.map((label) => (
              <Step key={label}>
                <StepLabel>{label}</StepLabel>
              </Step>
            ))}
          </Stepper>

          {error && <Alert severity="error" sx={{ mb: 3 }}>{error}</Alert>}

          {renderStepContent(activeStep)}

          <Box sx={{ display: 'flex', justifyContent: 'flex-end', mt: 4, pt: 2, borderTop: 1, borderColor: 'divider' }}>
            {activeStep > 0 && activeStep < steps.length && (
              <Button onClick={handleBack} sx={{ mr: 1 }} disabled={loading}>
                Back
              </Button>
            )}
            {activeStep > 0 && activeStep < steps.length - 1 && (
              <Button variant="contained" onClick={handleNext} disabled={loading}>
                Next
              </Button>
            )}
            {activeStep === steps.length - 1 && (
              <Button variant="contained" onClick={handleSubmit} disabled={loading}>
                {loading ? <CircularProgress size={24} /> : 'Save Certificate'}
              </Button>
            )}
          </Box>
        </Paper>
      </Box>

      {/* Preview Dialog */}
      <Dialog open={previewOpen} onClose={() => setPreviewOpen(false)} maxWidth="lg" fullWidth>
        <DialogTitle sx={{ m: 0, p: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          Preview: {file?.name}
          <IconButton onClick={() => setPreviewOpen(false)}><CloseIcon /></IconButton>
        </DialogTitle>
        <DialogContent dividers sx={{ p: 0, height: '80vh' }}>
          {previewUrl && (
            file?.type === 'application/pdf' ? (
              <iframe src={previewUrl} width="100%" height="100%" title="Preview" />
            ) : (
              <Box sx={{ display: 'flex', justifyContent: 'center', p: 2 }}>
                <img src={previewUrl} alt="Preview" style={{ maxWidth: '100%', maxHeight: '100%', objectFit: 'contain' }} />
              </Box>
            )
          )}
        </DialogContent>
      </Dialog>
    </Container>
  );
};

export default CertificateWizard;
