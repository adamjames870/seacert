import { useEffect, useState } from 'react';
import { Typography, Container, Box, Paper, List, ListItemText, Alert, CircularProgress, FormControl, InputLabel, Select, MenuItem, IconButton, Tooltip, Collapse, Divider, Link, ListItemButton, Button, TextField, InputAdornment, Stack, Dialog, DialogTitle, DialogContent, DialogContentText, DialogActions, Accordion, AccordionSummary, AccordionDetails } from '@mui/material';
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
import { supabase } from '../supabaseClient';
import { API_BASE_URL } from '../config';
import { formatDate } from '../utils/dateUtils';
import { getCountryName } from '../utils/countryData';
import { Link as RouterLink } from 'react-router-dom';

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
  predecessors?: Certificate[];
  successors?: Certificate[];
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
        throw new Error(`Error fetching certificates: ${response.statusText}`);
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

  const flattenPredecessors = (cert: Certificate): Certificate[] => {
    if (!cert.predecessors || cert.predecessors.length === 0) return [];
    
    let result: Certificate[] = [];
    for (const pred of cert.predecessors) {
      // Add direct predecessor
      result.push(pred);
      // Recursively add its predecessors
      result = [...result, ...flattenPredecessors(pred)];
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
    if (cert.successors && cert.successors.length > 0) {
      return false;
    }

    const query = searchQuery.toLowerCase();
    return (
      cert['cert-type-name'].toLowerCase().includes(query) ||
      cert['issuer-name'].toLowerCase().includes(query)
    );
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
    <Container>
      <Box sx={{ mt: 4 }}>
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

        {!loading && !error && certificates.length === 0 && (
          <Typography variant="body1" sx={{ mt: 2 }}>
            No certificates found.
          </Typography>
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
                  
                  return (
                    <Paper 
                      key={cert.id} 
                      elevation={0} 
                      sx={{ 
                        mb: 2, 
                        border: 1, 
                        borderColor: styles.borderColor, 
                        bgcolor: styles.bgcolor,
                        overflow: 'hidden' 
                      }}
                    >
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
                            <Typography variant="subtitle1" sx={{ fontWeight: 600, color: styles.textColor }}>
                              {cert['cert-type-name']}
                            </Typography>
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
                          {cert['expiry-date'] && new Date(cert['expiry-date']).getFullYear() > 1 && (
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
                          {flattenPredecessors(cert).length > 0 && (
                            <Box sx={{ mt: 2 }}>
                              <Typography variant="caption" sx={{ color: styles.secondaryTextColor, display: 'block' }}>
                                Predecessors
                              </Typography>
                              <List sx={{ p: 0 }}>
                                {flattenPredecessors(cert).map((pred) => (
                                  <Box key={pred.id} sx={{ mb: 0.5 }}>
                                    <Typography variant="body2" sx={{ color: styles.textColor }}>
                                      {pred['cert-type-name']} ({pred['cert-number']}) • Issued: {formatDate(pred['issued-date'])}
                                    </Typography>
                                  </Box>
                                ))}
                              </List>
                            </Box>
                          )}
                          <Box sx={{ mt: 3, display: 'flex', justifyContent: 'flex-end', gap: 1 }}>
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
                                }
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
                                  }
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
                                  }
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
                                }
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
                                }
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
                      return (
                        <Paper 
                          key={cert.id} 
                          elevation={0} 
                          sx={{ 
                            mb: 2, 
                            border: 1, 
                            borderColor: 'divider', 
                            bgcolor: 'action.hover',
                            opacity: 0.8,
                            overflow: 'hidden' 
                          }}
                        >
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
                                <Typography variant="subtitle1" sx={{ fontWeight: 600, color: 'text.secondary' }}>
                                  {cert['cert-type-name']}
                                </Typography>
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
                              </Box>
                              {flattenPredecessors(cert).length > 0 && (
                                <Box sx={{ mt: 2 }}>
                                  <Typography variant="caption" sx={{ color: 'text.secondary', display: 'block' }}>
                                    Predecessors
                                  </Typography>
                                  <List sx={{ p: 0 }}>
                                    {flattenPredecessors(cert).map((pred) => (
                                      <Box key={pred.id} sx={{ mb: 0.5 }}>
                                        <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                                          {pred['cert-type-name']} ({pred['cert-number']}) • Issued: {formatDate(pred['issued-date'])}
                                        </Typography>
                                      </Box>
                                    ))}
                                  </List>
                                </Box>
                              )}
                              <Box sx={{ mt: 3, display: 'flex', justifyContent: 'flex-end', gap: 1 }}>
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
                                    }
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
                                    }
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
    </Container>
  );
};

export default Certificates;
