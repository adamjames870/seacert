import { useEffect, useState } from 'react';
import { 
  Typography, 
  Container, 
  Box, 
  Paper, 
  List, 
  ListItemText, 
  Alert, 
  CircularProgress, 
  FormControl, 
  InputLabel, 
  Select, 
  MenuItem, 
  IconButton, 
  Tooltip, 
  Collapse, 
  Divider, 
  Link, 
  ListItemButton, 
  Button, 
  TextField, 
  InputAdornment, 
  Stack, 
  Dialog, 
  DialogTitle, 
  DialogContent, 
  DialogContentText, 
  DialogActions, 
  Accordion, 
  AccordionSummary, 
  AccordionDetails,
  ButtonGroup,
  Menu,
  ListItemIcon
} from '@mui/material';
import ArrowUpwardIcon from '@mui/icons-material/ArrowUpward';
import ArrowDownwardIcon from '@mui/icons-material/ArrowDownward';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import ExpandLessIcon from '@mui/icons-material/ExpandLess';
import AddIcon from '@mui/icons-material/Add';
import EditIcon from '@mui/icons-material/Edit';
import SearchIcon from '@mui/icons-material/Search';
import DeleteIcon from '@mui/icons-material/Delete';
import ArchiveIcon from '@mui/icons-material/Archive';
import UnarchiveIcon from '@mui/icons-material/Unarchive';
import AutorenewIcon from '@mui/icons-material/Autorenew';
import AttachFileIcon from '@mui/icons-material/AttachFile';
import VisibilityIcon from '@mui/icons-material/Visibility';
import DownloadIcon from '@mui/icons-material/Download';
import DescriptionIcon from '@mui/icons-material/Description';
import CloseIcon from '@mui/icons-material/Close';
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown';
import { supabase } from '../supabaseClient';
import { API_BASE_URL } from '../config';
import { formatDate } from '../utils/dateUtils';
import { getCountryName } from '../utils/countryData';
import ReportPreviewDialog from '../components/ReportPreviewDialog';
import { getErrorMessage } from '../utils/errorUtils';
import { Link as RouterLink, useNavigate } from 'react-router-dom';

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
  'cert-number': string;
  'issuer-id': string;
  'issuer-name': string;
  'issuer-country': string;
  'issuer-website': string;
  'issued-date': string;
  'expiry-date': string;
  'alternative-name': string;
  remarks: string;
  deleted?: boolean;
  predecessors?: Predecessor[];
  'has-successors'?: boolean;
  'document-path'?: string;
  'document-url'?: string;
}

type SortField = 'cert-type-name' | 'issuer-name' | 'issued-date' | 'expiry-date';
type SortOrder = 'asc' | 'desc';

const Certificates = () => {
  const [certificates, setCertificates] = useState<Certificate[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [sortBy, setSortBy] = useState<SortField>('expiry-date');
  const [sortOrder, setSortOrder] = useState<SortOrder>('asc');
  const [expandedId, setExpandedId] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [retireDialogOpen, setRetireDialogOpen] = useState(false);
  const [activateDialogOpen, setActivateDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [selectedCert, setSelectedCert] = useState<Certificate | null>(null);
  const [actionLoading, setActionLoading] = useState(false);
  const [previewOpen, setPreviewOpen] = useState(false);
  const [menuAnchorEl, setMenuAnchorEl] = useState<null | HTMLElement>(null);
  const [activeMenuCertId, setActiveMenuCertId] = useState<string | null>(null);
  const [reportDialogOpen, setReportDialogOpen] = useState(false);
  const [draggedOverCertId, setDraggedOverCertId] = useState<string | null>(null);

  const navigate = useNavigate();

  const handleDragOver = (e: React.DragEvent, certId: string) => {
    e.preventDefault();
    e.stopPropagation();
    setDraggedOverCertId(certId);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setDraggedOverCertId(null);
  };

  const handleDrop = (e: React.DragEvent, certId: string) => {
    e.preventDefault();
    e.stopPropagation();
    setDraggedOverCertId(null);
    
    const file = e.dataTransfer.files?.[0];
    if (file && (file.type === 'application/pdf' || file.type === 'image/jpeg' || file.type === 'image/jpg')) {
      navigate(`/update-certificate/${certId}`, { state: { droppedFile: file } });
    }
  };

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>, certId: string) => {
    event.stopPropagation();
    setMenuAnchorEl(event.currentTarget);
    setActiveMenuCertId(certId);
  };

  const handleMenuClose = () => {
    setMenuAnchorEl(null);
    setActiveMenuCertId(null);
  };

  const handlePreviewOpen = (cert: Certificate) => {
    setSelectedCert(cert);
    setPreviewOpen(true);
    handleMenuClose();
  };

  const handleDownload = (cert: Certificate) => {
    if (!cert['document-url']) return;
    
    const link = document.createElement('a');
    link.href = cert['document-url'];
    // Try to suggest a filename, though R2 pre-signed URLs might override this
    const extension = cert['document-path']?.split('.').pop() || 'file';
    link.download = `${cert['cert-type-short-name']}_${cert['cert-number']}.${extension}`;
    link.target = '_blank';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    handleMenuClose();
  };

  const fetchCertificates = async (retry = true) => {
    try {
      if (retry) setLoading(true);
      const { data: { session }, error: sessionError } = await supabase.auth.getSession();
      
      if (sessionError) throw sessionError;
      if (!session) {
        setError('Not authenticated');
        setLoading(false);
        return;
      }

      const response = await fetch(`${API_BASE_URL}/api/certificates`, {
        headers: {
          'Authorization': `Bearer ${session.access_token}`,
        },
      });

      if (response.status === 401 && retry) {
        const { data: { session: newSession } } = await supabase.auth.refreshSession();
        if (newSession) {
          return fetchCertificates(false);
        }
      }

      if (!response.ok) {
        const message = await getErrorMessage(response);
        throw new Error(message);
      }

      const data = await response.json();
      setCertificates(data);
      setError(null);
    } catch (err: any) {
      setError(err.message || 'Failed to load certificates');
    } finally {
      if (retry) setLoading(false);
    }
  };

  useEffect(() => {
    fetchCertificates();
  }, []);

  const getExpiryStatus = (expiryDate: string) => {
    if (!expiryDate || new Date(expiryDate).getFullYear() <= 1) return 'normal';
    
    const now = new Date();
    const expiry = new Date(expiryDate);
    const twelveMonthsFromNow = new Date();
    twelveMonthsFromNow.setFullYear(now.getFullYear() + 1);

    if (expiry < now) return 'expired';
    if (expiry <= twelveMonthsFromNow) return 'expiring-soon';
    return 'normal';
  };

  const getStatusStyles = (status: string) => {
    switch (status) {
      case 'expired':
        return {
          bgcolor: '#fff1f2', // Rose 50
          borderColor: '#fecdd3', // Rose 200
          textColor: '#9f1239', // Rose 800
          secondaryTextColor: '#be123b', // Rose 700
          labelColor: '#e11d48', // Rose 600
        };
      case 'expiring-soon':
        return {
          bgcolor: '#fffbeb', // Amber 50
          borderColor: '#fde68a', // Amber 200
          textColor: '#92400e', // Amber 800
          secondaryTextColor: '#b45309', // Amber 700
          labelColor: '#d97706', // Amber 600
        };
      default:
        return {
          bgcolor: 'background.paper',
          borderColor: 'divider',
          textColor: 'text.primary',
          secondaryTextColor: 'text.secondary',
          labelColor: 'primary.main',
        };
    }
  };

  const flattenPredecessors = (cert: Certificate): { reason: string, certificate: Certificate }[] => {
    if (!cert.predecessors || cert.predecessors.length === 0) return [];
    
    let result: { reason: string, certificate: Certificate }[] = [];
    for (const pred of cert.predecessors) {
      // Add direct predecessor
      result.push(pred);
      // Recursively add its predecessors
      const subPredecessors = flattenPredecessors(pred.certificate);
      result = [...result, ...subPredecessors];
    }
    return result;
  };

  const handleRetireClick = (cert: Certificate) => {
    setSelectedCert(cert);
    setRetireDialogOpen(true);
  };

  const handleActivateClick = (cert: Certificate) => {
    setSelectedCert(cert);
    setActivateDialogOpen(true);
  };

  const handleDeleteClick = (cert: Certificate) => {
    setSelectedCert(cert);
    setDeleteDialogOpen(true);
  };

  const handleRetireConfirm = async () => {
    if (!selectedCert) return;
    setActionLoading(true);
    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) throw new Error('Not authenticated');

      const response = await fetch(`${API_BASE_URL}/api/certificates`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${session.access_token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          id: selectedCert.id,
          deleted: true,
        }),
      });

      if (!response.ok) throw new Error('Failed to retire certificate');

      setRetireDialogOpen(false);
      fetchCertificates();
    } catch (err: any) {
      setError(err.message);
    } finally {
      setActionLoading(false);
    }
  };

  const handleActivateConfirm = async () => {
    if (!selectedCert) return;
    setActionLoading(true);
    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) throw new Error('Not authenticated');

      const response = await fetch(`${API_BASE_URL}/api/certificates`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${session.access_token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          id: selectedCert.id,
          deleted: false,
        }),
      });

      if (!response.ok) throw new Error('Failed to activate certificate');

      setActivateDialogOpen(false);
      fetchCertificates();
    } catch (err: any) {
      setError(err.message);
    } finally {
      setActionLoading(false);
    }
  };

  const handleDeleteConfirm = async () => {
    if (!selectedCert) return;
    setActionLoading(true);
    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) throw new Error('Not authenticated');

      const response = await fetch(`${API_BASE_URL}/api/certificates?id=${selectedCert.id}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${session.access_token}`,
        },
      });

      if (!response.ok) throw new Error('Failed to delete certificate');

      setDeleteDialogOpen(false);
      fetchCertificates();
    } catch (err: any) {
      setError(err.message);
    } finally {
      setActionLoading(false);
    }
  };

  const filteredCertificates = certificates.filter((cert) => {
    // If a certificate has any successors, it should not be displayed in the main list
    if (cert['has-successors']) {
      return false;
    }

    const query = searchQuery.toLowerCase();
    const matchesSearch = 
      cert['cert-type-name'].toLowerCase().includes(query) ||
      cert['issuer-name'].toLowerCase().includes(query);
    
    return matchesSearch;
  });

  const sortedCertificates = [...filteredCertificates].sort((a, b) => {
    if (sortBy === 'expiry-date') {
      const dateA = new Date(a['expiry-date']);
      const dateB = new Date(b['expiry-date']);
      const isNoExpiryA = dateA.getFullYear() <= 1;
      const isNoExpiryB = dateB.getFullYear() <= 1;

      if (isNoExpiryA && isNoExpiryB) return 0;
      if (isNoExpiryA) return 1; // Always bottom
      if (isNoExpiryB) return -1; // Always bottom

      const comparison = dateA.getTime() - dateB.getTime();
      return sortOrder === 'asc' ? comparison : -comparison;
    }

    let comparison = 0;
    if (sortBy === 'issued-date') {
      const dateA = new Date(a[sortBy]).getTime();
      const dateB = new Date(b[sortBy]).getTime();
      comparison = dateA - dateB;
    } else {
      const valA = a[sortBy].toLowerCase();
      const valB = b[sortBy].toLowerCase();
      comparison = valA.localeCompare(valB);
    }

    return sortOrder === 'asc' ? comparison : -comparison;
  });

  const activeCertificates = sortedCertificates.filter(cert => !cert.deleted);
  const retiredCertificates = sortedCertificates.filter(cert => cert.deleted);

  const toggleSortOrder = () => {
    setSortOrder(prev => prev === 'asc' ? 'desc' : 'asc');
  };

  const handleExpand = (id: string) => {
    setExpandedId(prev => prev === id ? null : id);
  };

  return (
    <Container sx={{ position: 'relative', minHeight: '60vh' }}>
      <Box sx={{ 
        mt: 4, 
        opacity: !loading && !error && certificates.length === 0 ? 0.3 : 1,
        transition: 'opacity 0.3s'
      }}>
        <Stack 
          direction={{ xs: 'column', md: 'row' }} 
          spacing={2} 
          justifyContent="space-between" 
          alignItems={{ xs: 'stretch', md: 'center' }} 
          sx={{ mb: 3 }}
        >
          <Typography variant="h4" component="h1">
            Certificates
          </Typography>
          
          <Stack 
            direction={{ xs: 'column', sm: 'row' }} 
            spacing={2} 
            alignItems={{ xs: 'stretch', sm: 'center' }}
          >
            <TextField
              size="small"
              placeholder="Search certificates..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <SearchIcon fontSize="small" />
                  </InputAdornment>
                ),
              }}
              sx={{ minWidth: { xs: '100%', sm: 250 } }}
            />
            <Button
              variant="outlined"
              color="primary"
              startIcon={<DescriptionIcon />}
              onClick={() => setReportDialogOpen(true)}
              sx={{ whiteSpace: 'nowrap' }}
            >
              Report
            </Button>
            <Button
              variant="contained"
              color="primary"
              startIcon={<AddIcon />}
              component={RouterLink}
              to="/add-certificate"
              sx={{ whiteSpace: 'nowrap' }}
            >
              Add Certificate
            </Button>

            {!loading && !error && certificates.length > 0 && (
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <FormControl size="small" sx={{ minWidth: 120, flexGrow: { xs: 1, sm: 0 } }}>
                  <InputLabel id="sort-by-label">Sort By</InputLabel>
                  <Select
                    labelId="sort-by-label"
                    value={sortBy}
                    label="Sort By"
                    onChange={(e) => setSortBy(e.target.value as SortField)}
                  >
                    <MenuItem value="cert-type-name">Name</MenuItem>
                    <MenuItem value="issuer-name">Issuer</MenuItem>
                    <MenuItem value="issued-date">Issue Date</MenuItem>
                    <MenuItem value="expiry-date">Expiry Date</MenuItem>
                  </Select>
                </FormControl>
                <Tooltip title={sortOrder === 'asc' ? "Sort Descending" : "Sort Ascending"}>
                  <IconButton onClick={toggleSortOrder} color="primary" size="small">
                    {sortOrder === 'asc' ? <ArrowUpwardIcon /> : <ArrowDownwardIcon />}
                  </IconButton>
                </Tooltip>
              </Box>
            )}
          </Stack>
        </Stack>

        {loading && (
          <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
            <CircularProgress />
          </Box>
        )}

        {error && (
          <Alert severity="error" sx={{ mt: 2 }}>
            {error}
          </Alert>
        )}

        {!loading && !error && certificates.length > 0 && filteredCertificates.length === 0 && (
          <Typography variant="body1" sx={{ mt: 2 }}>
            No certificates match your search.
          </Typography>
        )}

        {!loading && !error && filteredCertificates.length > 0 && (
          <>
            {activeCertificates.length > 0 ? (
              <List sx={{ mt: 2 }}>
                {activeCertificates.map((cert) => {
                  const status = getExpiryStatus(cert['expiry-date']);
                  const styles = getStatusStyles(status);
                  const hasAttachment = !!cert['document-path'];
                  
                  return (
                    <Paper 
                      key={cert.id} 
                      elevation={0} 
                      onDragOver={(e) => handleDragOver(e, cert.id)}
                      onDragLeave={handleDragLeave}
                      onDrop={(e) => handleDrop(e, cert.id)}
                      sx={{ 
                        mb: 2, 
                        border: 1, 
                        borderColor: draggedOverCertId === cert.id ? 'primary.main' : styles.borderColor, 
                        bgcolor: draggedOverCertId === cert.id ? 'action.hover' : styles.bgcolor,
                        overflow: 'hidden',
                        position: 'relative',
                        transition: 'all 0.2s ease-in-out',
                        boxShadow: draggedOverCertId === cert.id ? 2 : 'none',
                        '&::before': !hasAttachment ? {
                          content: '""',
                          position: 'absolute',
                          left: 0,
                          top: 0,
                          bottom: 0,
                          width: '4px',
                          bgcolor: 'warning.light',
                          opacity: 0.6
                        } : {}
                      }}
                    >
                      {draggedOverCertId === cert.id && (
                        <Box sx={{ 
                          position: 'absolute', 
                          top: 0, 
                          left: 0, 
                          right: 0, 
                          bottom: 0, 
                          display: 'flex', 
                          alignItems: 'center', 
                          justifyContent: 'center',
                          bgcolor: 'rgba(25, 118, 210, 0.08)',
                          zIndex: 1,
                          pointerEvents: 'none'
                        }}>
                          <Typography variant="subtitle2" sx={{ color: 'primary.main', fontWeight: 'bold' }}>
                            Drop to add/replace attachment
                          </Typography>
                        </Box>
                      )}
                      <ListItemButton 
                        onClick={() => handleExpand(cert.id)}
                        sx={{ 
                          display: 'flex', 
                          flexDirection: { xs: 'column', sm: 'row' },
                          justifyContent: 'space-between', 
                          alignItems: { xs: 'flex-start', sm: 'center' },
                          p: 2,
                          gap: 2
                        }}
                      >
                        <ListItemText 
                          primary={
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                              <Typography variant="subtitle1" sx={{ fontWeight: 600, color: styles.textColor }}>
                                {cert['cert-type-name']}
                              </Typography>
                              {cert['document-path'] && (
                                <Tooltip title="Has attachment">
                                  <AttachFileIcon sx={{ fontSize: '1rem', color: styles.secondaryTextColor, transform: 'rotate(45deg)' }} />
                                </Tooltip>
                              )}
                            </Box>
                          }
                          secondary={
                            <Box component="span" sx={{ color: styles.secondaryTextColor }}>
                              <Typography component="span" variant="body2" sx={{ color: 'inherit', display: 'block' }}>
                                Issuer: {cert['issuer-name']}
                              </Typography>
                              <Typography component="span" variant="caption" sx={{ color: 'inherit' }}>
                                {`No. ${cert['cert-number']} — Issued: ${formatDate(cert['issued-date'])}`}
                              </Typography>
                            </Box>
                          }
                        />
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, width: { xs: '100%', sm: 'auto' }, justifyContent: { xs: 'space-between', sm: 'flex-end' } }}>
                          {!cert['expiry-date'] || new Date(cert['expiry-date']).getFullYear() <= 1 ? (
                            <Box sx={{ textAlign: { xs: 'left', sm: 'right' } }}>
                              <Typography variant="caption" sx={{ display: 'block', textTransform: 'uppercase', fontWeight: 'bold', fontSize: '0.7rem', color: styles.secondaryTextColor }}>
                                Validity
                              </Typography>
                              <Typography variant="body1" sx={{ fontWeight: 'bold', color: styles.labelColor }}>
                                Does not expire
                              </Typography>
                            </Box>
                          ) : (
                            <Box sx={{ textAlign: { xs: 'left', sm: 'right' } }}>
                              <Typography variant="caption" sx={{ display: 'block', textTransform: 'uppercase', fontWeight: 'bold', fontSize: '0.7rem', color: styles.secondaryTextColor }}>
                                Expires
                              </Typography>
                              <Typography variant="body1" sx={{ fontWeight: 'bold', color: styles.labelColor }}>
                                {formatDate(cert['expiry-date'])}
                              </Typography>
                            </Box>
                          )}
                          <IconButton size="small" sx={{ color: styles.textColor }}>
                            {expandedId === cert.id ? <ExpandLessIcon /> : <ExpandMoreIcon />}
                          </IconButton>
                        </Box>
                      </ListItemButton>
                      <Collapse in={expandedId === cert.id} timeout="auto" unmountOnExit>
                        <Divider sx={{ borderColor: styles.borderColor }} />
                        <Box sx={{ p: 2, bgcolor: status === 'normal' ? 'action.hover' : 'rgba(0, 0, 0, 0.02)' }}>
                          <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', sm: '1fr 1fr' }, gap: 2 }}>
                            <Box>
                              <Typography variant="caption" sx={{ color: styles.secondaryTextColor, display: 'block' }}>
                                Short Name
                              </Typography>
                              <Typography variant="body2" sx={{ color: styles.textColor }}>
                                {cert['cert-type-short-name'] || 'N/A'}
                              </Typography>
                            </Box>
                            <Box>
                              <Typography variant="caption" sx={{ color: styles.secondaryTextColor, display: 'block' }}>
                                STCW Reference
                              </Typography>
                              <Typography variant="body2" sx={{ color: styles.textColor }}>
                                {cert['cert-type-stcw-ref'] || 'N/A'}
                              </Typography>
                            </Box>
                            <Box>
                              <Typography variant="caption" sx={{ color: styles.secondaryTextColor, display: 'block' }}>
                                Issuer Country
                              </Typography>
                              <Typography variant="body2" sx={{ color: styles.textColor }}>
                                {getCountryName(cert['issuer-country'])}
                              </Typography>
                            </Box>
                            <Box>
                              <Typography variant="caption" sx={{ color: styles.secondaryTextColor, display: 'block' }}>
                                Issuer Website
                              </Typography>
                              {cert['issuer-website'] ? (
                                <Link href={cert['issuer-website'].startsWith('http') ? cert['issuer-website'] : `https://${cert['issuer-website']}`} target="_blank" rel="noopener" variant="body2" sx={{ color: styles.labelColor }}>
                                  {cert['issuer-website']}
                                </Link>
                              ) : (
                                <Typography variant="body2" sx={{ color: styles.textColor }}>N/A</Typography>
                              )}
                              <Box sx={{ mt: 0.5 }}>
                                <Link 
                                  component={RouterLink} 
                                  to={`/edit-issuer/${cert['issuer-id']}`} 
                                  state={{ from: 'certificates' }}
                                  variant="caption"
                                  sx={{ display: 'inline-flex', alignItems: 'center', gap: 0.5, textDecoration: 'none', color: styles.labelColor }}
                                >
                                  <EditIcon sx={{ fontSize: '0.8rem' }} /> Edit Issuer
                                </Link>
                              </Box>
                            </Box>
                          </Box>
                          {cert.remarks && (
                            <Box sx={{ mt: 2 }}>
                              <Typography variant="caption" sx={{ color: styles.secondaryTextColor, display: 'block' }}>
                                Remarks
                              </Typography>
                              <Typography variant="body2" sx={{ color: styles.textColor }}>
                                {cert.remarks}
                              </Typography>
                            </Box>
                          )}
                          {cert['document-url'] && (
                            <Box sx={{ mt: 2 }}>
                              <ButtonGroup 
                                variant="outlined" 
                                size="small" 
                                sx={{
                                  '& .MuiButton-root': {
                                    color: styles.textColor,
                                    borderColor: styles.borderColor,
                                    '&:hover': {
                                      borderColor: styles.textColor,
                                      bgcolor: 'rgba(0, 0, 0, 0.04)'
                                    }
                                  }
                                }}
                              >
                                <Button
                                  startIcon={<VisibilityIcon />}
                                  onClick={(e) => {
                                    e.stopPropagation();
                                    handlePreviewOpen(cert);
                                  }}
                                >
                                  Preview
                                </Button>
                                <Button
                                  size="small"
                                  onClick={(e) => handleMenuOpen(e, cert.id)}
                                >
                                  <ArrowDropDownIcon />
                                </Button>
                              </ButtonGroup>
                              <Menu
                                anchorEl={menuAnchorEl}
                                open={Boolean(menuAnchorEl) && activeMenuCertId === cert.id}
                                onClose={handleMenuClose}
                                onClick={(e) => e.stopPropagation()}
                              >
                                <MenuItem onClick={() => handlePreviewOpen(cert)}>
                                  <ListItemIcon>
                                    <VisibilityIcon fontSize="small" />
                                  </ListItemIcon>
                                  Preview
                                </MenuItem>
                                <MenuItem onClick={() => handleDownload(cert)}>
                                  <ListItemIcon>
                                    <DownloadIcon fontSize="small" />
                                  </ListItemIcon>
                                  Download
                                </MenuItem>
                              </Menu>
                            </Box>
                          )}
                          {flattenPredecessors(cert).length > 0 && (
                            <Box sx={{ mt: 2 }}>
                              <Typography variant="caption" sx={{ color: styles.secondaryTextColor, display: 'block' }}>
                                Predecessors
                              </Typography>
                              <List sx={{ p: 0 }}>
                                {flattenPredecessors(cert).map((pred, idx) => (
                                  <Box key={`${pred.certificate.id}-${idx}`} sx={{ mb: 0.5 }}>
                                    <Typography variant="body2" sx={{ color: styles.textColor }}>
                                      {pred.certificate['cert-type-name']} ({pred.certificate['cert-number']}) • Issued: {formatDate(pred.certificate['issued-date'])} {pred.reason && `• Reason: ${pred.reason}`}
                                    </Typography>
                                  </Box>
                                ))}
                              </List>
                            </Box>
                          )}
                          <Box sx={{ 
                            mt: 3, 
                            display: 'flex', 
                            justifyContent: 'flex-end', 
                            gap: 1,
                            flexWrap: 'wrap'
                          }}>
                            <Button
                              variant="outlined"
                              size="small"
                              startIcon={<EditIcon />}
                              component={RouterLink}
                              to={`/update-certificate/${cert.id}`}
                              onClick={(e) => e.stopPropagation()}
                              sx={{
                                color: styles.textColor,
                                borderColor: styles.borderColor,
                                '&:hover': {
                                  borderColor: styles.textColor,
                                  bgcolor: 'rgba(0, 0, 0, 0.04)'
                                },
                                flex: { xs: '1 1 auto', sm: 'none' }
                              }}
                            >
                              Edit
                            </Button>
                            <Button
                                variant="outlined"
                                size="small"
                                startIcon={<AddIcon />}
                                component={RouterLink}
                                to="/add-certificate"
                                state={{ 
                                  certTypeId: cert['cert-type-id'],
                                  supersedes: cert.id,
                                  supersedeReason: 'updated'
                                }}
                                onClick={(e) => e.stopPropagation()}
                                sx={{
                                  color: styles.textColor,
                                  borderColor: styles.borderColor,
                                  '&:hover': {
                                    borderColor: styles.textColor,
                                    bgcolor: 'rgba(0, 0, 0, 0.04)'
                                  },
                                  flex: { xs: '1 1 auto', sm: 'none' }
                                }}
                            >
                              Update with New
                            </Button>
                            <Button
                                variant="outlined"
                                size="small"
                                startIcon={<AutorenewIcon />}
                                component={RouterLink}
                                to="/add-certificate"
                                state={{ 
                                  certTypeId: cert['cert-type-id'],
                                  supersedes: cert.id,
                                  supersedeReason: 'replaced'
                                }}
                                onClick={(e) => e.stopPropagation()}
                                sx={{
                                  color: styles.textColor,
                                  borderColor: styles.borderColor,
                                  '&:hover': {
                                    borderColor: styles.textColor,
                                    bgcolor: 'rgba(0, 0, 0, 0.04)'
                                  },
                                  flex: { xs: '1 1 auto', sm: 'none' }
                                }}
                            >
                              Replace with New
                            </Button>
                            <Button
                              variant="outlined"
                              size="small"
                              startIcon={<ArchiveIcon />}
                              onClick={(e) => {
                                e.stopPropagation();
                                handleRetireClick(cert);
                              }}
                              sx={{ 
                                color: styles.textColor, 
                                borderColor: styles.borderColor,
                                '&:hover': {
                                  borderColor: styles.textColor,
                                  bgcolor: 'rgba(0, 0, 0, 0.04)'
                                },
                                flex: { xs: '1 1 auto', sm: 'none' }
                              }}
                            >
                              Retire
                            </Button>
                            <Button
                              variant="outlined"
                              size="small"
                              color="error"
                              startIcon={<DeleteIcon />}
                              onClick={(e) => {
                                e.stopPropagation();
                                handleDeleteClick(cert);
                              }}
                              sx={{ 
                                '&:hover': {
                                  bgcolor: 'error.main',
                                  color: 'white'
                                },
                                flex: { xs: '1 1 auto', sm: 'none' }
                              }}
                            >
                              Delete Permanently
                            </Button>
                          </Box>
                        </Box>
                      </Collapse>
                    </Paper>
                  );
                })}
              </List>
            ) : (
              <Typography variant="body1" sx={{ mt: 2 }}>
                No active certificates found matching your search.
              </Typography>
            )}

            {retiredCertificates.length > 0 && (
              <Accordion 
                sx={{ 
                  mt: 4, 
                  bgcolor: 'transparent', 
                  boxShadow: 'none', 
                  '&:before': { display: 'none' },
                  border: '1px solid',
                  borderColor: 'divider',
                  borderRadius: '8px !important',
                  overflow: 'hidden'
                }}
              >
                <AccordionSummary
                  expandIcon={<ExpandMoreIcon />}
                  sx={{ 
                    bgcolor: 'action.hover',
                    '& .MuiAccordionSummary-content': { 
                      m: 1.5,
                      display: 'flex',
                      alignItems: 'center',
                      gap: 1
                    },
                    '&.Mui-expanded': { 
                      minHeight: 'unset',
                      borderBottom: '1px solid',
                      borderColor: 'divider'
                    },
                    '&:hover': {
                      bgcolor: 'action.selected'
                    }
                  }}
                >
                  <ArchiveIcon sx={{ color: 'text.secondary', fontSize: '1.2rem' }} />
                  <Typography variant="subtitle1" sx={{ fontWeight: 600, color: 'text.secondary' }}>
                    Retired Certificates ({retiredCertificates.length})
                  </Typography>
                </AccordionSummary>
                <AccordionDetails sx={{ p: 0, mt: 2 }}>
                  <List>
                    {retiredCertificates.map((cert) => {
                      const hasAttachment = !!cert['document-path'];
                      return (
                        <Paper 
                          key={cert.id} 
                          elevation={0} 
                          onDragOver={(e) => handleDragOver(e, cert.id)}
                          onDragLeave={handleDragLeave}
                          onDrop={(e) => handleDrop(e, cert.id)}
                          sx={{ 
                            mb: 2, 
                            border: 1, 
                            borderColor: draggedOverCertId === cert.id ? 'primary.main' : 'divider', 
                            bgcolor: draggedOverCertId === cert.id ? 'action.selected' : 'action.hover',
                            opacity: draggedOverCertId === cert.id ? 1 : 0.8,
                            overflow: 'hidden',
                            position: 'relative',
                            transition: 'all 0.2s ease-in-out',
                            boxShadow: draggedOverCertId === cert.id ? 2 : 'none',
                            '&::before': !hasAttachment ? {
                              content: '""',
                              position: 'absolute',
                              left: 0,
                              top: 0,
                              bottom: 0,
                              width: '4px',
                              bgcolor: 'warning.light',
                              opacity: 0.6
                            } : {}
                          }}
                        >
                          {draggedOverCertId === cert.id && (
                            <Box sx={{ 
                              position: 'absolute', 
                              top: 0, 
                              left: 0, 
                              right: 0, 
                              bottom: 0, 
                              display: 'flex', 
                              alignItems: 'center', 
                              justifyContent: 'center',
                              bgcolor: 'rgba(25, 118, 210, 0.08)',
                              zIndex: 1,
                              pointerEvents: 'none'
                            }}>
                              <Typography variant="subtitle2" sx={{ color: 'primary.main', fontWeight: 'bold' }}>
                                Drop to add/replace attachment
                              </Typography>
                            </Box>
                          )}
                          <ListItemButton 
                            onClick={() => handleExpand(cert.id)}
                            sx={{ 
                              display: 'flex', 
                              flexDirection: { xs: 'column', sm: 'row' },
                              justifyContent: 'space-between', 
                              alignItems: { xs: 'flex-start', sm: 'center' },
                              p: 2,
                              gap: 2
                            }}
                          >
                            <ListItemText 
                              primary={
                                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                  <Typography variant="subtitle1" sx={{ fontWeight: 600, color: 'text.secondary' }}>
                                    {cert['cert-type-name']}
                                  </Typography>
                                  {cert['document-path'] && (
                                    <Tooltip title="Has attachment">
                                      <AttachFileIcon sx={{ fontSize: '1rem', color: 'text.disabled', transform: 'rotate(45deg)' }} />
                                    </Tooltip>
                                  )}
                                </Box>
                              }
                              secondary={
                                <Box component="span" sx={{ color: 'text.secondary' }}>
                                  <Typography component="span" variant="body2" sx={{ color: 'inherit', display: 'block' }}>
                                    Issuer: {cert['issuer-name']}
                                  </Typography>
                                  <Typography component="span" variant="caption" sx={{ color: 'inherit' }}>
                                    {`No. ${cert['cert-number']} — Issued: ${formatDate(cert['issued-date'])}`}
                                  </Typography>
                                </Box>
                              }
                            />
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, width: { xs: '100%', sm: 'auto' }, justifyContent: 'flex-end' }}>
                              <IconButton size="small" sx={{ color: 'text.secondary' }}>
                                {expandedId === cert.id ? <ExpandLessIcon /> : <ExpandMoreIcon />}
                              </IconButton>
                            </Box>
                          </ListItemButton>
                          <Collapse in={expandedId === cert.id} timeout="auto" unmountOnExit>
                            <Divider />
                            <Box sx={{ p: 2, bgcolor: 'rgba(0, 0, 0, 0.02)' }}>
                              <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', sm: '1fr 1fr' }, gap: 2 }}>
                                <Box>
                                  <Typography variant="caption" sx={{ color: 'text.secondary', display: 'block' }}>
                                    Short Name
                                  </Typography>
                                  <Typography variant="body2">
                                    {cert['cert-type-short-name'] || 'N/A'}
                                  </Typography>
                                </Box>
                                <Box>
                                  <Typography variant="caption" sx={{ color: 'text.secondary', display: 'block' }}>
                                    STCW Reference
                                  </Typography>
                                  <Typography variant="body2">
                                    {cert['cert-type-stcw-ref'] || 'N/A'}
                                  </Typography>
                                </Box>
                                <Box sx={{ gridColumn: { sm: 'span 2' } }}>
                                  <Typography variant="caption" sx={{ color: 'text.secondary', display: 'block' }}>
                                    Remarks
                                  </Typography>
                                  <Typography variant="body2">
                                    {cert.remarks || 'No remarks provided.'}
                                  </Typography>
                                </Box>
                                {cert['document-url'] && (
                                  <Box sx={{ gridColumn: { sm: 'span 2' }, mt: 1 }}>
                                    <ButtonGroup variant="outlined" size="small">
                                      <Button
                                        startIcon={<VisibilityIcon />}
                                        onClick={(e) => {
                                          e.stopPropagation();
                                          handlePreviewOpen(cert);
                                        }}
                                      >
                                        Preview
                                      </Button>
                                      <Button
                                        size="small"
                                        onClick={(e) => handleMenuOpen(e, cert.id)}
                                      >
                                        <ArrowDropDownIcon />
                                      </Button>
                                    </ButtonGroup>
                                    <Menu
                                      anchorEl={menuAnchorEl}
                                      open={Boolean(menuAnchorEl) && activeMenuCertId === cert.id}
                                      onClose={handleMenuClose}
                                      onClick={(e) => e.stopPropagation()}
                                    >
                                      <MenuItem onClick={() => handlePreviewOpen(cert)}>
                                        <ListItemIcon>
                                          <VisibilityIcon fontSize="small" />
                                        </ListItemIcon>
                                        Preview
                                      </MenuItem>
                                      <MenuItem onClick={() => handleDownload(cert)}>
                                        <ListItemIcon>
                                          <DownloadIcon fontSize="small" />
                                        </ListItemIcon>
                                        Download
                                      </MenuItem>
                                    </Menu>
                                  </Box>
                                )}
                              </Box>
                              {flattenPredecessors(cert).length > 0 && (
                                <Box sx={{ mt: 2 }}>
                                  <Typography variant="caption" sx={{ color: 'text.secondary', display: 'block' }}>
                                    Predecessors
                                  </Typography>
                                  <List sx={{ p: 0 }}>
                                    {flattenPredecessors(cert).map((pred, idx) => (
                                      <Box key={`${pred.certificate.id}-${idx}`} sx={{ mb: 0.5 }}>
                                        <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                                          {pred.certificate['cert-type-name']} ({pred.certificate['cert-number']}) • Issued: {formatDate(pred.certificate['issued-date'])} {pred.reason && `• Reason: ${pred.reason}`}
                                        </Typography>
                                      </Box>
                                    ))}
                                  </List>
                                </Box>
                              )}
                              <Box sx={{ 
                                mt: 3, 
                                display: 'flex', 
                                justifyContent: 'flex-end', 
                                gap: 1,
                                flexWrap: 'wrap'
                              }}>
                                <Button
                                  variant="outlined"
                                  size="small"
                                  startIcon={<UnarchiveIcon />}
                                  onClick={(e) => {
                                    e.stopPropagation();
                                    handleActivateClick(cert);
                                  }}
                                  sx={{
                                    color: 'primary.main',
                                    borderColor: 'primary.main',
                                    '&:hover': {
                                      bgcolor: 'primary.main',
                                      color: 'white'
                                    },
                                    flex: { xs: '1 1 auto', sm: 'none' }
                                  }}
                                >
                                  Make Active
                                </Button>
                                <Button
                                  variant="outlined"
                                  size="small"
                                  color="error"
                                  startIcon={<DeleteIcon />}
                                  onClick={(e) => {
                                    e.stopPropagation();
                                    handleDeleteClick(cert);
                                  }}
                                  sx={{ 
                                    '&:hover': {
                                      bgcolor: 'error.main',
                                      color: 'white'
                                    },
                                    flex: { xs: '1 1 auto', sm: 'none' }
                                  }}
                                >
                                  Delete Permanently
                                </Button>
                                <Button
                                  variant="outlined"
                                  size="small"
                                  startIcon={<EditIcon />}
                                  component={RouterLink}
                                  to={`/update-certificate/${cert.id}`}
                                  onClick={(e) => e.stopPropagation()}
                                  sx={{ flex: { xs: '1 1 auto', sm: 'none' } }}
                                >
                                  Edit
                                </Button>
                                <Button
                                  variant="outlined"
                                  size="small"
                                  startIcon={<AddIcon />}
                                  component={RouterLink}
                                  to="/add-certificate"
                                  state={{ 
                                    certTypeId: cert['cert-type-id'],
                                    supersedes: cert.id,
                                    supersedeReason: 'updated'
                                  }}
                                  onClick={(e) => e.stopPropagation()}
                                  sx={{ flex: { xs: '1 1 auto', sm: 'none' } }}
                                >
                                  Update with New
                                </Button>
                                <Button
                                  variant="outlined"
                                  size="small"
                                  startIcon={<AutorenewIcon />}
                                  component={RouterLink}
                                  to="/add-certificate"
                                  state={{ 
                                    certTypeId: cert['cert-type-id'],
                                    supersedes: cert.id,
                                    supersedeReason: 'replaced'
                                  }}
                                  onClick={(e) => e.stopPropagation()}
                                  sx={{ flex: { xs: '1 1 auto', sm: 'none' } }}
                                >
                                  Replace with New
                                </Button>
                              </Box>
                            </Box>
                          </Collapse>
                        </Paper>
                      );
                    })}
                  </List>
                </AccordionDetails>
              </Accordion>
            )}
          </>
        )}
      </Box>

      {/* Retire Confirmation Dialog */}
      <Dialog
        open={retireDialogOpen}
        onClose={() => !actionLoading && setRetireDialogOpen(false)}
      >
        <DialogTitle>Retire Certificate</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Are you sure you want to retire certificate "{selectedCert?.['cert-type-name']}"?
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setRetireDialogOpen(false)} disabled={actionLoading}>
            Cancel
          </Button>
          <Button onClick={handleRetireConfirm} color="primary" autoFocus disabled={actionLoading}>
            {actionLoading ? <CircularProgress size={24} /> : 'Retire'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Activate Confirmation Dialog */}
      <Dialog
        open={activateDialogOpen}
        onClose={() => !actionLoading && setActivateDialogOpen(false)}
      >
        <DialogTitle>Activate Certificate</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Do you want to restore certificate "{selectedCert?.['cert-type-name']}" to your active list?
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setActivateDialogOpen(false)} disabled={actionLoading}>
            Cancel
          </Button>
          <Button onClick={handleActivateConfirm} color="primary" autoFocus disabled={actionLoading}>
            {actionLoading ? <CircularProgress size={24} /> : 'Restore'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <Dialog
        open={deleteDialogOpen}
        onClose={() => !actionLoading && setDeleteDialogOpen(false)}
      >
        <DialogTitle>Delete Certificate</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Are you sure you want to <strong>permanently delete</strong> certificate "{selectedCert?.['cert-type-name']}"?
            <br /><br />
            <strong>Warning: This action cannot be undone.</strong>
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)} disabled={actionLoading}>
            Cancel
          </Button>
          <Button onClick={handleDeleteConfirm} color="error" autoFocus disabled={actionLoading}>
            {actionLoading ? <CircularProgress size={24} /> : 'Delete'}
          </Button>
        </DialogActions>
      </Dialog>
      {/* Document Preview Dialog */}
      <Dialog
        open={previewOpen}
        onClose={() => setPreviewOpen(false)}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle sx={{ m: 0, p: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          Document Preview - {selectedCert?.['cert-type-short-name']}
          <IconButton onClick={() => setPreviewOpen(false)}>
            <CloseIcon />
          </IconButton>
        </DialogTitle>
        <DialogContent dividers sx={{ p: 0, height: '75vh', overflow: 'hidden' }}>
          {selectedCert?.['document-url'] && (
            selectedCert['document-path']?.toLowerCase().endsWith('.pdf') ? (
              <iframe
                src={selectedCert['document-url']}
                title="PDF Preview"
                width="100%"
                height="100%"
                style={{ border: 'none' }}
              />
            ) : (
              <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100%', p: 2 }}>
                <img
                  src={selectedCert['document-url']}
                  alt="Document Preview"
                  style={{ maxWidth: '100%', maxHeight: '100%', objectFit: 'contain' }}
                />
              </Box>
            )
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setPreviewOpen(false)}>Close</Button>
          <Button 
            variant="contained" 
            startIcon={<DownloadIcon />}
            onClick={() => selectedCert && handleDownload(selectedCert)}
          >
            Download
          </Button>
        </DialogActions>
      </Dialog>
      <ReportPreviewDialog open={reportDialogOpen} onClose={() => setReportDialogOpen(false)} />

      {!loading && !error && certificates.length === 0 && (
        <Box 
          sx={{ 
            position: 'absolute',
            top: '50%',
            left: '50%',
            transform: 'translate(-50%, -50%)',
            display: 'flex', 
            flexDirection: 'column',
            alignItems: 'center', 
            justifyContent: 'center',
            textAlign: 'center',
            gap: 3,
            zIndex: 10,
            width: '100%',
            maxWidth: 500,
            p: 4,
            pointerEvents: 'auto'
          }}
        >
          <Typography variant="h5" color="text.secondary" sx={{ fontWeight: 500 }}>
            Your certificate list is empty
          </Typography>
          <Button
            variant="contained"
            color="primary"
            size="large"
            startIcon={<AddIcon />}
            component={RouterLink}
            to="/add-certificate"
            sx={{ 
              px: 6, 
              py: 2, 
              fontSize: '1.2rem',
              borderRadius: 4,
              boxShadow: 6,
              '&:hover': {
                boxShadow: 10,
                transform: 'translateY(-2px)'
              },
              transition: 'all 0.2s'
            }}
          >
            Add First Certificate
          </Button>
        </Box>
      )}
    </Container>
  );
};

export default Certificates;
