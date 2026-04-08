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

interface Ship {
  id: string;
  'created-at': string;
  'updated-at': string;
  name: string;
  'ship-type-id': string;
  'ship-type-name': string;
  'imo-number': string;
  gt: number;
  flag: string;
  'propulsion-power': number;
  status: 'approved' | 'provisional';
  'created-by': string;
}

interface ShipType {
  id: string;
  name: string;
  description: string;
}

const Ships = () => {
  const [ships, setShips] = useState<Ship[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [userRole, setUserRole] = useState<string | null>(null);
  const [currentUserId, setCurrentUserId] = useState<string | null>(null);
  const [resolveDialogOpen, setResolveDialogOpen] = useState(false);
  const [provisionalShip, setProvisionalShip] = useState<Ship | null>(null);
  const [replacementId, setReplacementId] = useState('');
  const [resolving, setResolving] = useState(false);
  const [shipTypes, setShipTypes] = useState<ShipType[]>([]);

  useEffect(() => {
    const fetchUserData = async () => {
      const { data: { session } } = await supabase.auth.getSession();
      if (session) {
        setUserRole(session.user?.app_metadata?.role || 'user');
        setCurrentUserId(session.user.id);
      }
    };
    fetchUserData();

    const fetchData = async () => {
      try {
        setLoading(true);
        const { data: { session } } = await supabase.auth.getSession();
        
        if (!session) {
          setError('Not authenticated');
          setLoading(false);
          return;
        }

        // Fetch lookups for ship type descriptions
        const lookupResponse = await fetch(`${API_BASE_URL}/api/seatime/lookups`, {
          headers: {
            'Authorization': `Bearer ${session.access_token}`,
          },
        });

        let lookupData: ShipType[] = [];
        if (lookupResponse.ok) {
          const data = await lookupResponse.json();
          lookupData = data['ship-types'];
          setShipTypes(lookupData);
        }

        const response = await fetch(`${API_BASE_URL}/api/ships`, {
          headers: {
            'Authorization': `Bearer ${session.access_token}`,
          },
        });

        if (!response.ok) {
          throw new Error(`Error fetching ships: ${response.statusText}`);
        }

        const data: Ship[] = await response.json();
        setShips(data);
        setError(null);
      } catch (err: any) {
        setError(err.message || 'Failed to load ships');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  const getShipTypeDescription = (ship: Ship) => {
    const type = shipTypes.find(t => t.id === ship['ship-type-id']);
    return type ? type.description : ship['ship-type-name'];
  };

  const handleApprove = async (ship: Ship) => {
    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) return;

      const response = await fetch(`${API_BASE_URL}/api/admin/ships/approve/${ship.id}`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${session.access_token}`,
        },
      });

      if (!response.ok) {
        throw new Error('Failed to approve ship');
      }

      setShips(prev => prev.map(s => s.id === ship.id ? { ...s, status: 'approved' } : s));
    } catch (err: any) {
      setError(err.message);
    }
  };

  const handleOpenResolve = (ship: Ship) => {
    setProvisionalShip(ship);
    setReplacementId('');
    setResolveDialogOpen(true);
  };

  const handleResolve = async () => {
    if (!provisionalShip || !replacementId) return;
    
    setResolving(true);
    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) return;

      const response = await fetch(`${API_BASE_URL}/api/admin/ships/resolve`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${session.access_token}`,
        },
        body: JSON.stringify({
          'provisional-id': provisionalShip.id,
          'replacement-id': replacementId
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to resolve ship');
      }

      setShips(prev => prev.filter(s => s.id !== provisionalShip.id));
      setResolveDialogOpen(false);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setResolving(false);
    }
  };

  const filteredShips = ships.filter((ship) => {
    const query = searchQuery.toLowerCase();
    return (
      ship.name.toLowerCase().includes(query) ||
      ship['imo-number']?.toLowerCase().includes(query) ||
      ship['ship-type-name']?.toLowerCase().includes(query) ||
      ship.flag?.toLowerCase().includes(query)
    );
  });

  const sortedShips = [...filteredShips].sort((a, b) => 
    a.name.localeCompare(b.name)
  );

  const canEdit = (ship: Ship) => {
    if (userRole === 'admin') return true;
    return ship.status === 'provisional' && ship['created-by'] === currentUserId;
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
            Ships
          </Typography>
          
          <Stack 
            direction={{ xs: 'column', sm: 'row' }} 
            spacing={2} 
            alignItems={{ xs: 'stretch', sm: 'center' }}
          >
            <TextField
              size="small"
              placeholder="Search ships..."
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
              to="/add-ship"
              sx={{ whiteSpace: 'nowrap' }}
            >
              Add ship
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

        {!loading && !error && ships.length === 0 && (
          <Typography variant="body1" sx={{ mt: 2 }}>
            No ships found.
          </Typography>
        )}

        {!loading && !error && ships.length > 0 && filteredShips.length === 0 && (
          <Typography variant="body1" sx={{ mt: 2 }}>
            No ships match your search.
          </Typography>
        )}

        {!loading && !error && filteredShips.length > 0 && (
          <List sx={{ mt: 2 }}>
            {sortedShips.map((ship) => (
              <Paper 
                key={ship.id} 
                elevation={0} 
                sx={{ 
                  mb: 1, 
                  border: 1, 
                  borderColor: ship.status === 'provisional' ? '#fef3c7' : 'divider', 
                  bgcolor: ship.status === 'provisional' ? '#fffbeb' : 'background.paper',
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
                        <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>
                          {ship.name}
                        </Typography>
                        {ship.status === 'provisional' && (
                          <Chip label="Provisional" size="small" color="warning" variant="outlined" />
                        )}
                      </Box>
                    }
                    secondary={
                      <>
                        <Typography component="span" variant="body2" color="text.secondary">
                          IMO: {ship['imo-number']} | Type: {getShipTypeDescription(ship)} | Flag: {ship.flag}
                        </Typography>
                        <br />
                        <Typography component="span" variant="body2" color="text.secondary">
                          GT: {ship.gt} | Power: {ship['propulsion-power'] || 'N/A'} kW
                        </Typography>
                      </>
                    }
                  />
                  <Box sx={{ display: 'flex', justifyContent: 'flex-end', alignItems: 'center', gap: 0.5 }}>
                    {userRole === 'admin' && ship.status === 'provisional' && (
                      <>
                        <Tooltip title="Approve">
                          <IconButton 
                            onClick={() => handleApprove(ship)}
                            color="success"
                          >
                            <CheckIcon />
                          </IconButton>
                        </Tooltip>
                        <Tooltip title="Resolve (Merge & Delete)">
                          <IconButton 
                            onClick={() => handleOpenResolve(ship)}
                            color="warning"
                          >
                            <SwapHorizIcon />
                          </IconButton>
                        </Tooltip>
                      </>
                    )}
                    {canEdit(ship) && (
                      <Tooltip title="Edit Ship">
                        <IconButton 
                          component={RouterLink} 
                          to={`/edit-ship/${ship.id}`}
                        >
                          <EditIcon />
                        </IconButton>
                      </Tooltip>
                    )}
                  </Box>
                </Box>
              </Paper>
            ))}
          </List>
        )}

        <Dialog open={resolveDialogOpen} onClose={() => !resolving && setResolveDialogOpen(false)} maxWidth="sm" fullWidth>
          <DialogTitle>Resolve Provisional Ship</DialogTitle>
          <DialogContent>
            <DialogContentText sx={{ mb: 2 }}>
              Migrate all seatime records from <strong>{provisionalShip?.name}</strong> to an existing approved ship and delete the provisional one.
            </DialogContentText>
            <FormControl fullWidth sx={{ mt: 1 }}>
              <InputLabel id="replacement-ship-label">Replacement Approved Ship</InputLabel>
              <Select
                labelId="replacement-ship-label"
                value={replacementId}
                label="Replacement Approved Ship"
                onChange={(e) => setReplacementId(e.target.value)}
                disabled={resolving}
              >
                {ships
                  .filter(s => s.status === 'approved' && s.id !== provisionalShip?.id)
                  .sort((a, b) => a.name.localeCompare(b.name))
                  .map(s => (
                    <MenuItem key={s.id} value={s.id}>
                      {s.name} ({s['imo-number']})
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

export default Ships;
