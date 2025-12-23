import { useEffect, useState } from 'react';
import { Typography, Container, Box, Paper, List, ListItemText, Alert, CircularProgress, FormControl, InputLabel, Select, MenuItem, IconButton, Tooltip, Collapse, Divider, Link, ListItemButton, Button, TextField, InputAdornment } from '@mui/material';
import ArrowUpwardIcon from '@mui/icons-material/ArrowUpward';
import ArrowDownwardIcon from '@mui/icons-material/ArrowDownward';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import ExpandLessIcon from '@mui/icons-material/ExpandLess';
import AddIcon from '@mui/icons-material/Add';
import EditIcon from '@mui/icons-material/Edit';
import SearchIcon from '@mui/icons-material/Search';
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

  useEffect(() => {
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
          // Token might be expired, attempt to refresh by calling getSession again
          // Supabase getSession handles refreshing if it can.
          // We can also try refreshSession explicitly if needed, but getSession is usually enough.
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

    fetchCertificates();
  }, []);

  const filteredCertificates = certificates.filter((cert) => {
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

  const toggleSortOrder = () => {
    setSortOrder(prev => prev === 'asc' ? 'desc' : 'asc');
  };

  const handleExpand = (id: string) => {
    setExpandedId(prev => prev === id ? null : id);
  };

  return (
    <Container>
      <Box sx={{ mt: 4 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h4" component="h1">
            Certificates
          </Typography>
          
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
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
              sx={{ width: 250 }}
            />
            <Button
              variant="contained"
              color="primary"
              startIcon={<AddIcon />}
              component={RouterLink}
              to="/add-certificate"
            >
              Add Certificate
            </Button>

            {!loading && !error && certificates.length > 0 && (
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                <FormControl size="small" sx={{ minWidth: 150 }}>
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
                  <IconButton onClick={toggleSortOrder} color="primary">
                    {sortOrder === 'asc' ? <ArrowUpwardIcon /> : <ArrowDownwardIcon />}
                  </IconButton>
                </Tooltip>
              </Box>
            )}
          </Box>
        </Box>

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
          <List sx={{ mt: 2 }}>
            {sortedCertificates.map((cert) => {
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
                      justifyContent: 'space-between', 
                      alignItems: 'center',
                      p: 2
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
                          <Typography component="span" variant="body2" sx={{ color: 'inherit' }}>
                            Issuer: {cert['issuer-name']}
                          </Typography>
                          {` — No. ${cert['cert-number']} — Issued on: ${formatDate(cert['issued-date'])}`}
                        </Box>
                      }
                    />
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                      {cert['expiry-date'] && new Date(cert['expiry-date']).getFullYear() > 1 && (
                        <Box sx={{ textAlign: 'right', mr: 2 }}>
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
                      <Box sx={{ mt: 3, display: 'flex', justifyContent: 'flex-end' }}>
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
                          Update
                        </Button>
                      </Box>
                    </Box>
                  </Collapse>
                </Paper>
              );
            })}
          </List>
        )}
      </Box>
    </Container>
  );
};

export default Certificates;
