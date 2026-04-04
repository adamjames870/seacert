import { useEffect, useState } from 'react';
import { 
  Typography, 
  Container, 
  Box, 
  Paper, 
  Table, 
  TableBody, 
  TableCell, 
  TableContainer, 
  TableHead, 
  TableRow, 
  Button, 
  CircularProgress, 
  Alert,
  IconButton,
  Tabs,
  Tab,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Stack
} from '@mui/material';
import { Edit, Trash2, Plus, ArrowLeft } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { supabase } from '../supabaseClient';
import { API_BASE_URL } from '../config';

interface LookupItem {
  id: string;
  name: string;
  description: string;
}

const ManageSeatimeLookups = () => {
  const [tabValue, setTabValue] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [lookups, setLookups] = useState<{
    'ship-types': LookupItem[];
    'voyage-types': LookupItem[];
    'period-types': LookupItem[];
  }>({
    'ship-types': [],
    'voyage-types': [],
    'period-types': [],
  });

  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingItem, setEditingItem] = useState<LookupItem | null>(null);
  const [formData, setFormData] = useState({ name: '', description: '' });

  useEffect(() => {
    fetchLookups();
  }, []);

  const fetchLookups = async () => {
    setLoading(true);
    try {
      const { data: { session } } = await supabase.auth.getSession();
      const response = await fetch(`${API_BASE_URL}/api/seatime/lookups`, {
        headers: { 'Authorization': `Bearer ${session?.access_token}` },
      });
      if (!response.ok) throw new Error('Failed to fetch lookups');
      const data = await response.json();
      setLookups(data);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleTabChange = (_: any, newValue: number) => {
    setTabValue(newValue);
  };

  const getCurrentType = () => {
    if (tabValue === 0) return 'ship-types';
    if (tabValue === 1) return 'voyage-types';
    return 'period-types';
  };

  const handleOpenDialog = (item: LookupItem | null = null) => {
    setEditingItem(item);
    setFormData(item ? { name: item.name, description: item.description } : { name: '', description: '' });
    setDialogOpen(true);
  };

  const handleCloseDialog = () => {
    setDialogOpen(false);
    setEditingItem(null);
  };

  const handleSave = async () => {
    try {
      const { data: { session } } = await supabase.auth.getSession();
      const type = getCurrentType();
      const baseUrl = `${API_BASE_URL}/api/admin/seatime/${type}`;
      
      const response = await fetch(baseUrl, {
        method: editingItem ? 'PUT' : 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${session?.access_token}` 
        },
        body: JSON.stringify(editingItem ? { ...formData, id: editingItem.id } : formData),
      });

      if (!response.ok) throw new Error('Failed to save lookup item');
      
      fetchLookups();
      handleCloseDialog();
    } catch (err: any) {
      alert(err.message);
    }
  };

  const handleDelete = async (id: string) => {
    if (!window.confirm('Are you sure you want to delete this item?')) return;
    
    try {
      const { data: { session } } = await supabase.auth.getSession();
      const type = getCurrentType();
      const url = `${API_BASE_URL}/api/admin/seatime/${type}/${id}`;
      
      const response = await fetch(url, {
        method: 'DELETE',
        headers: { 
          'Authorization': `Bearer ${session?.access_token}` 
        },
      });

      if (!response.ok) throw new Error('Failed to delete lookup item');
      
      fetchLookups();
    } catch (err: any) {
      alert(err.message);
    }
  };

  if (loading) return <Container sx={{ mt: 4, textAlign: 'center' }}><CircularProgress /></Container>;

  const currentList = lookups[getCurrentType()];

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Typography variant="h4" gutterBottom sx={{ fontWeight: 700 }}>
        Manage Seatime Lookups
      </Typography>

      {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

      <Paper sx={{ width: '100%', mb: 2 }}>
        <Tabs value={tabValue} onChange={handleTabChange} sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tab label="Ship Types" />
          <Tab label="Voyage Types" />
          <Tab label="Period Types" />
        </Tabs>
        
        <Box sx={{ p: 2, display: 'flex', justifyContent: 'flex-end' }}>
          <Button startIcon={<Plus size={20} />} variant="contained" onClick={() => handleOpenDialog()}>
            Add New {getCurrentType().replace('-types', '').replace('-', ' ')}
          </Button>
        </Box>

        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Name</TableCell>
                <TableCell>Description</TableCell>
                <TableCell align="right">Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {currentList.map((item) => (
                <TableRow key={item.id}>
                  <TableCell sx={{ fontWeight: 600 }}>{item.name}</TableCell>
                  <TableCell>{item.description}</TableCell>
                  <TableCell align="right">
                    <IconButton size="small" onClick={() => handleOpenDialog(item)}>
                      <Edit size={18} />
                    </IconButton>
                    <IconButton size="small" color="error" onClick={() => handleDelete(item.id)}>
                      <Trash2 size={18} />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>

      <Dialog open={dialogOpen} onClose={handleCloseDialog} fullWidth maxWidth="sm">
        <DialogTitle>{editingItem ? 'Edit' : 'Add'} {getCurrentType().replace('-types', '')}</DialogTitle>
        <DialogContent>
          <Stack spacing={3} sx={{ mt: 1 }}>
            <TextField
              fullWidth
              label="Name"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            />
            <TextField
              fullWidth
              label="Description"
              multiline
              rows={3}
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            />
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseDialog}>Cancel</Button>
          <Button onClick={handleSave} variant="contained">Save</Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
};

export default ManageSeatimeLookups;
