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
  IconButton, 
  Tooltip,
  Button,
  TextField,
  InputAdornment,
  Stack,
  Chip,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  FormControl,
  InputLabel,
  Select,
  MenuItem
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import SearchIcon from '@mui/icons-material/Search';
import AddIcon from '@mui/icons-material/Add';
import CheckIcon from '@mui/icons-material/Check';
import SwapHorizIcon from '@mui/icons-material/SwapHoriz';
import { Link as RouterLink } from 'react-router-dom';
import { supabase } from '../supabaseClient';
import { API_BASE_URL } from '../config';

interface CertType {
  id: string;
  name: string;
  'short-name': string;
  'stcw-reference': string;
  'normal-validity-months': number;
  status?: 'approved' | 'provisional';
  'created-by'?: string;
}

const CertTypes = () => {
  const [certTypes, setCertTypes] = useState<CertType[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [userRole, setUserRole] = useState<string | null>(null);
  const [resolveDialogOpen, setResolveDialogOpen] = useState(false);
  const [provisionalType, setProvisionalType] = useState<CertType | null>(null);
  const [replacementId, setReplacementId] = useState('');
  const [resolving, setResolving] = useState(false);

  const getMissingFieldsStatus = (type: CertType) => {
    const missing = [];
    if (!type['short-name']) missing.push('short-name');
    if (!type['stcw-reference']) missing.push('stcw-reference');
    
    if (missing.length > 0) return 'incomplete';
    return 'normal';
  };

  const getStatusStyles = (status: string) => {
    switch (status) {
      case 'incomplete':
        return {
          bgcolor: '#fffbeb', // Amber 50
          borderColor: '#fef3c7', // Amber 100
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
    const fetchUserData = async () => {
      const { data: { session } } = await supabase.auth.getSession();
      if (session) {
        // First try to get role from app_metadata (set by Supabase for our admins)
        const role = session.user?.app_metadata?.role;
        if (role) {
          setUserRole(role);
          return;
        }

        // Fallback to fetching if metadata is not available (though it should be)
        if (session.access_token) {
          try {
            const response = await fetch(`${API_BASE_URL}/admin/users`, {
              headers: {
                'Authorization': `Bearer ${session.access_token}`,
              },
            });
            if (response.ok) {
              const data = await response.json();
              const user = Array.isArray(data) ? data.find((u: any) => u.id === session.user.id) : data;
              setUserRole(user?.role || 'user');
            } else {
              // If we can't fetch /admin/users (e.g. we're not an admin), we're likely a regular user
              setUserRole('user');
            }
          } catch (error) {
            console.error('Error fetching user role:', error);
            setUserRole('user');
          }
        }
      }
    };
    fetchUserData();

    const fetchCertTypes = async () => {
      try {
        setLoading(true);
        const { data: { session } } = await supabase.auth.getSession();
        
        if (!session) {
          setError('Not authenticated');
          setLoading(false);
          return;
        }

        const response = await fetch(`${API_BASE_URL}/api/cert-types`, {
          headers: {
            'Authorization': `Bearer ${session.access_token}`,
          },
        });

        if (!response.ok) {
          throw new Error(`Error fetching certificate types: ${response.statusText}`);
        }

        const data = await response.json();
        setCertTypes(data);
        setError(null);
      } catch (err: any) {
        setError(err.message || 'Failed to load certificate types');
      } finally {
        setLoading(false);
      }
    };

    fetchCertTypes();
  }, []);

  const handleApprove = async (type: CertType) => {
    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) return;

      const response = await fetch(`${API_BASE_URL}/api/cert-types?id=${type.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${session.access_token}`,
        },
        body: JSON.stringify({
          id: type.id,
          name: type.name,
          'short-name': type['short-name'] || null,
          'stcw-reference': type['stcw-reference'] || null,
          'normal-validity-months': type['normal-validity-months'] || 0,
          status: 'approved'
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to approve certificate type');
      }

      setCertTypes(prev => prev.map(t => t.id === type.id ? { ...t, status: 'approved' } : t));
    } catch (err: any) {
      setError(err.message);
    }
  };

  const handleOpenResolve = (type: CertType) => {
    setProvisionalType(type);
    setReplacementId('');
    setResolveDialogOpen(true);
  };

  const handleResolve = async () => {
    if (!provisionalType || !replacementId) return;
    
    setResolving(true);
    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) return;

      const response = await fetch(`${API_BASE_URL}/admin/cert-types/resolve`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${session.access_token}`,
        },
        body: JSON.stringify({
          'provisional-id': provisionalType.id,
          'replacement-id': replacementId
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to resolve certificate type');
      }

      setCertTypes(prev => prev.filter(t => t.id !== provisionalType.id));
      setResolveDialogOpen(false);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setResolving(false);
    }
  };

  const filteredCertTypes = certTypes.filter((type) => {
    const query = searchQuery.toLowerCase();
    return (
      type.name.toLowerCase().includes(query) ||
      type['short-name']?.toLowerCase().includes(query) ||
      type['stcw-reference']?.toLowerCase().includes(query)
    );
  });

  const sortedCertTypes = [...filteredCertTypes].sort((a, b) => 
    a.name.localeCompare(b.name)
  );

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
            Certificate Types
          </Typography>
          
          <Stack 
            direction={{ xs: 'column', sm: 'row' }} 
            spacing={2} 
            alignItems={{ xs: 'stretch', sm: 'center' }}
          >
            <TextField
              size="small"
              placeholder="Search types..."
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
              to="/add-cert-type"
              sx={{ whiteSpace: 'nowrap' }}
            >
              Add certificate type
            </Button>
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

        {!loading && !error && certTypes.length === 0 && (
          <Typography variant="body1" sx={{ mt: 2 }}>
            No certificate types found.
          </Typography>
        )}

        {!loading && !error && certTypes.length > 0 && filteredCertTypes.length === 0 && (
          <Typography variant="body1" sx={{ mt: 2 }}>
            No certificate types match your search.
          </Typography>
        )}

        {!loading && !error && filteredCertTypes.length > 0 && (
          <List sx={{ mt: 2 }}>
            {sortedCertTypes.map((type) => {
              const status = getMissingFieldsStatus(type);
              const styles = getStatusStyles(status);

              return (
                <Paper 
                  key={type.id} 
                  elevation={0} 
                  sx={{ 
                    mb: 1, 
                    border: 1, 
                    borderColor: styles.borderColor, 
                    bgcolor: styles.bgcolor,
                    overflow: 'hidden' 
                  }}
                >
                  <Box 
                    sx={{ 
                      display: 'flex', 
                      flexDirection: { xs: 'column', sm: 'row' },
                      justifyContent: 'space-between', 
                      alignItems: { xs: 'stretch', sm: 'center' },
                      p: 2,
                      gap: 1
                    }}
                  >
                    <ListItemText 
                      primary={
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          <Typography variant="subtitle1" sx={{ fontWeight: 600, color: styles.textColor }}>
                            {type.name}
                          </Typography>
                          {type.status === 'provisional' && (
                            <Chip label="Provisional" size="small" color="warning" variant="outlined" />
                          )}
                        </Box>
                      }
                      secondary={
                        <>
                          <Typography component="span" variant="body2" sx={{ color: styles.secondaryTextColor }}>
                            Short Name: {type['short-name'] || (
                              <Box component="span" sx={{ fontStyle: 'italic', fontWeight: 'bold', color: styles.labelColor }}>
                                Missing
                              </Box>
                            )}
                            {' | '}
                            STCW: {type['stcw-reference'] || (
                              <Box component="span" sx={{ fontStyle: 'italic', fontWeight: 'bold', color: styles.labelColor }}>
                                Missing
                              </Box>
                            )}
                          </Typography>
                          {type['normal-validity-months'] ? (
                            <Typography variant="body2" sx={{ color: styles.secondaryTextColor }}>
                              Validity: {type['normal-validity-months']} months
                            </Typography>
                          ) : (
                            <Typography variant="body2" sx={{ color: styles.secondaryTextColor }}>
                              Validity: Does not expire
                            </Typography>
                          )}
                        </>
                      }
                    />
                    <Box sx={{ display: 'flex', justifyContent: 'flex-end', alignItems: 'center', gap: 0.5 }}>
                      {userRole === 'admin' && type.status === 'provisional' && (
                        <>
                          <Tooltip title="Approve">
                            <IconButton 
                              onClick={() => handleApprove(type)}
                              color="success"
                            >
                              <CheckIcon />
                            </IconButton>
                          </Tooltip>
                          <Tooltip title="Resolve (Replace & Delete)">
                            <IconButton 
                              onClick={() => handleOpenResolve(type)}
                              color="warning"
                            >
                              <SwapHorizIcon />
                            </IconButton>
                          </Tooltip>
                        </>
                      )}
                      <Tooltip title="Edit Certificate Type">
                        <IconButton 
                          component={RouterLink} 
                          to={`/edit-cert-type/${type.id}`}
                          sx={{ color: styles.labelColor }}
                        >
                          <EditIcon />
                        </IconButton>
                      </Tooltip>
                    </Box>
                  </Box>
                </Paper>
              );
            })}
          </List>
        )}

        <Dialog open={resolveDialogOpen} onClose={() => !resolving && setResolveDialogOpen(false)} maxWidth="sm" fullWidth>
          <DialogTitle>Resolve Provisional Certificate Type</DialogTitle>
          <DialogContent>
            <DialogContentText sx={{ mb: 2 }}>
              Migrate all existing certificates from <strong>{provisionalType?.name}</strong> to an existing approved type and delete the provisional one.
            </DialogContentText>
            <FormControl fullWidth sx={{ mt: 1 }}>
              <InputLabel id="replacement-type-label">Replacement Approved Type</InputLabel>
              <Select
                labelId="replacement-type-label"
                value={replacementId}
                label="Replacement Approved Type"
                onChange={(e) => setReplacementId(e.target.value)}
                disabled={resolving}
              >
                {certTypes
                  .filter(t => t.status === 'approved' && t.id !== provisionalType?.id)
                  .sort((a, b) => a.name.localeCompare(b.name))
                  .map(t => (
                    <MenuItem key={t.id} value={t.id}>
                      {t.name} ({t['short-name']})
                    </MenuItem>
                  ))
                }
              </Select>
            </FormControl>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setResolveDialogOpen(false)} disabled={resolving}>Cancel</Button>
            <Button 
              onClick={handleResolve} 
              color="warning" 
              variant="contained" 
              disabled={!replacementId || resolving}
              startIcon={resolving ? <CircularProgress size={20} color="inherit" /> : <SwapHorizIcon />}
            >
              Resolve & Delete
            </Button>
          </DialogActions>
        </Dialog>
      </Box>
    </Container>
  );
};

export default CertTypes;
