import { useState, useEffect } from 'react';
import { 
  Typography, 
  Container, 
  Box, 
  TextField, 
  Button, 
  Grid, 
  Paper, 
  FormControl, 
  InputLabel, 
  Select, 
  MenuItem, 
  Stack,
  Alert,
  CircularProgress
} from '@mui/material';
import { 
  Ship, 
  Save,
  ArrowLeft
} from 'lucide-react';
import { useNavigate, Link as RouterLink, useParams } from 'react-router-dom';
import { supabase } from '../supabaseClient';
import { API_BASE_URL } from '../config';

interface ShipType {
  id: string;
  name: string;
  description: string;
}

const ShipForm = () => {
  const { id } = useParams();
  const isEdit = Boolean(id);
  const navigate = useNavigate();
  const [loading, setLoading] = useState(isEdit);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [shipTypes, setShipTypes] = useState<ShipType[]>([]);
  const [status, setStatus] = useState<'approved' | 'provisional'>('provisional');

  // Form State
  const [name, setName] = useState('');
  const [shipTypeId, setShipTypeId] = useState('');
  const [imoNumber, setImoNumber] = useState('');
  const [gt, setGt] = useState<number | ''>('');
  const [flag, setFlag] = useState('');
  const [propulsionPower, setPropulsionPower] = useState<number | ''>('');

  useEffect(() => {
    fetchShipTypes();
    if (isEdit) {
      fetchShip();
    }
  }, [id]);

  const fetchShipTypes = async () => {
    try {
      const { data: { session } } = await supabase.auth.getSession();
      const response = await fetch(`${API_BASE_URL}/api/seatime/lookups`, {
        headers: { 'Authorization': `Bearer ${session?.access_token}` },
      });
      if (response.ok) {
        const data = await response.json();
        setShipTypes(data['ship-types']);
      }
    } catch (err) {
      console.error('Error fetching ship types:', err);
    }
  };

  const fetchShip = async () => {
    try {
      const { data: { session } } = await supabase.auth.getSession();
      const response = await fetch(`${API_BASE_URL}/api/ships`, {
        headers: { 'Authorization': `Bearer ${session?.access_token}` },
      });
      if (response.ok) {
        const data = await response.json();
        const ship = data.find((s: any) => s.id === id);
        if (ship) {
          setName(ship.name);
          setShipTypeId(ship['ship-type-id']);
          setImoNumber(ship['imo-number']);
          setGt(ship.gt);
          setFlag(ship.flag);
          setPropulsionPower(ship['propulsion-power'] || '');
          setStatus(ship.status);
        } else {
          setError('Ship not found');
        }
      }
    } catch (err) {
      setError('Failed to load ship details');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setError(null);

    try {
      const { data: { session } } = await supabase.auth.getSession();
      const body = {
        id,
        name,
        'ship-type-id': shipTypeId,
        'imo-number': imoNumber,
        gt: Number(gt),
        flag,
        'propulsion-power': propulsionPower ? Number(propulsionPower) : null
      };

      const response = await fetch(`${API_BASE_URL}/api/ships`, {
        method: isEdit ? 'PATCH' : 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${session?.access_token}`,
        },
        body: JSON.stringify(body),
      });

      if (!response.ok) {
        const errData = await response.json();
        throw new Error(errData.message || 'Failed to save ship');
      }

      navigate('/ships');
    } catch (err: any) {
      setError(err.message);
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return <Container sx={{ mt: 4, textAlign: 'center' }}><CircularProgress /></Container>;
  }

  return (
    <Container maxWidth="sm" sx={{ mt: 4, mb: 8 }}>
      <Button 
        component={RouterLink} 
        to="/ships" 
        startIcon={<ArrowLeft size={18} />} 
        sx={{ mb: 3 }}
      >
        Back to Ships
      </Button>

      <Box sx={{ mb: 4, display: 'flex', alignItems: 'center', gap: 2 }}>
        <Typography variant="h4" component="h1" sx={{ fontWeight: 700 }}>
          {isEdit ? 'Edit Ship' : 'Add New Ship'}
        </Typography>
        {isEdit && (
          <Chip 
            label={status.toUpperCase()} 
            color={status === 'approved' ? 'success' : 'warning'}
            variant="outlined"
            size="small"
            sx={{ fontWeight: 700 }}
          />
        )}
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      <form onSubmit={handleSubmit}>
        <Paper elevation={2} sx={{ p: 3, borderRadius: 2 }}>
          <Stack spacing={3}>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1 }}>
              <Ship size={24} className="text-primary" />
              <Typography variant="h6" sx={{ fontWeight: 600 }}>Ship Details</Typography>
            </Box>

            <TextField
              required
              fullWidth
              label="Ship Name"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />

            <FormControl fullWidth required>
              <InputLabel>Ship Type</InputLabel>
              <Select
                value={shipTypeId}
                label="Ship Type"
                onChange={(e) => setShipTypeId(e.target.value)}
              >
                {shipTypes.map((type) => (
                  <MenuItem key={type.id} value={type.id}>{type.description}</MenuItem>
                ))}
              </Select>
            </FormControl>

            <TextField
              required
              fullWidth
              label="IMO Number"
              placeholder="IMO1234567"
              value={imoNumber}
              onChange={(e) => setImoNumber(e.target.value)}
            />

            <Grid container spacing={2}>
              <Grid item xs={6}>
                <TextField
                  required
                  fullWidth
                  type="number"
                  label="Gross Tonnage (GT)"
                  value={gt}
                  onChange={(e) => setGt(e.target.value === '' ? '' : Number(e.target.value))}
                />
              </Grid>
              <Grid item xs={6}>
                <TextField
                  required
                  fullWidth
                  label="Flag State"
                  value={flag}
                  onChange={(e) => setFlag(e.target.value)}
                />
              </Grid>
            </Grid>

            <TextField
              fullWidth
              type="number"
              label="Propulsion Power (kW)"
              value={propulsionPower}
              onChange={(e) => setPropulsionPower(e.target.value === '' ? '' : Number(e.target.value))}
            />

            <Button
              type="submit"
              variant="contained"
              size="large"
              disabled={submitting}
              startIcon={submitting ? <CircularProgress size={20} /> : <Save size={20} />}
              sx={{ mt: 2, py: 1.5 }}
            >
              {submitting ? 'Saving...' : isEdit ? 'Update Ship' : 'Create Ship'}
            </Button>
          </Stack>
        </Paper>
      </form>
    </Container>
  );
};

export default ShipForm;
