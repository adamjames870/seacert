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
  FormControlLabel, 
  Switch, 
  IconButton, 
  Divider, 
  Stack,
  Alert,
  CircularProgress,
  Autocomplete,
  Card,
  CardContent
} from '@mui/material';
import { 
  Trash2, 
  Plus, 
  Ship, 
  Calendar, 
  MapPin, 
  Anchor,
  Save,
  ArrowLeft,
  Search,
  CheckCircle2
} from 'lucide-react';
import { useNavigate, Link as RouterLink } from 'react-router-dom';
import { supabase } from '../supabaseClient';
import { API_BASE_URL } from '../config';

interface ShipType {
  id: string;
  name: string;
  description: string;
}

interface VoyageType {
  id: string;
  name: string;
  description: string;
}

interface PeriodType {
  id: string;
  name: string;
  description: string;
}

interface SpecializedPeriod {
  'period-type-id': string;
  'start-date': string;
  'end-date': string;
  days: number;
  remarks: string;
}

const AddSeatime = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  // Lookups
  const [shipTypes, setShipTypes] = useState<ShipType[]>([]);
  const [voyageTypes, setVoyageTypes] = useState<VoyageType[]>([]);
  const [periodTypes, setPeriodTypes] = useState<PeriodType[]>([]);
  
  // Form State
  const [shipId, setShipId] = useState<string | null>(null);
  const [shipName, setShipName] = useState('');
  const [shipTypeId, setShipTypeId] = useState('');
  const [imoNumber, setImoNumber] = useState('');
  const [gt, setGt] = useState<number | ''>('');
  const [flag, setFlag] = useState('');
  const [propulsionPower, setPropulsionPower] = useState<number | ''>('');
  
  const [voyageTypeId, setVoyageTypeId] = useState('');
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [startLocation, setStartLocation] = useState('');
  const [endLocation, setEndLocation] = useState('');
  const [totalDays, setTotalDays] = useState<number>(0);
  const [company, setCompany] = useState('');
  const [capacity, setCapacity] = useState('');
  const [isWatchkeeping, setIsWatchkeeping] = useState(true);
  const [periods, setPeriods] = useState<SpecializedPeriod[]>([]);

  const [showShipForm, setShowShipForm] = useState(true);

  useEffect(() => {
    fetchLookups();
  }, []);

  useEffect(() => {
    if (startDate && endDate) {
      const start = new Date(startDate);
      const end = new Date(endDate);
      if (end >= start) {
        const diffTime = Math.abs(end.getTime() - start.getTime());
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24)) + 1;
        setTotalDays(diffDays);
      } else {
        setTotalDays(0);
      }
    } else {
      setTotalDays(0);
    }
  }, [startDate, endDate]);

  const fetchLookups = async () => {
    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) return;

      const response = await fetch(`${API_BASE_URL}/api/seatime/lookups`, {
        headers: {
          'Authorization': `Bearer ${session.access_token}`,
        },
      });

      if (!response.ok) throw new Error('Failed to fetch lookups');

      const data = await response.json();
      setShipTypes(data['ship-types']);
      setVoyageTypes(data['voyage-types']);
      setPeriodTypes(data['period-types']);
    } catch (err: any) {
      console.error('Error fetching lookups:', err);
      setError('Could not load lookup data. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleAddPeriod = () => {
    setPeriods([
      ...periods,
      {
        'period-type-id': '',
        'start-date': startDate,
        'end-date': endDate,
        days: 0,
        remarks: ''
      }
    ]);
  };

  const handleRemovePeriod = (index: number) => {
    const newPeriods = [...periods];
    newPeriods.splice(index, 1);
    setPeriods(newPeriods);
  };

  const handlePeriodChange = (index: number, field: keyof SpecializedPeriod, value: any) => {
    const newPeriods = [...periods];
    newPeriods[index] = { ...newPeriods[index], [field]: value };
    
    // Auto-calculate period days if dates change
    if (field === 'start-date' || field === 'end-date') {
      const pStart = new Date(field === 'start-date' ? value : newPeriods[index]['start-date']);
      const pEnd = new Date(field === 'end-date' ? value : newPeriods[index]['end-date']);
      if (pEnd >= pStart && !isNaN(pStart.getTime()) && !isNaN(pEnd.getTime())) {
        const diffTime = Math.abs(pEnd.getTime() - pStart.getTime());
        newPeriods[index].days = Math.ceil(diffTime / (1000 * 60 * 60 * 24)) + 1;
      } else {
        newPeriods[index].days = 0;
      }
    }
    
    setPeriods(newPeriods);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setError(null);

    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) throw new Error('No active session');

      const body: any = {
        'voyage-type-id': voyageTypeId,
        'start-date': startDate,
        'end-date': endDate,
        'start-location': startLocation,
        'end-location': endLocation,
        'total-days': totalDays,
        company,
        capacity,
        'is-watchkeeping': isWatchkeeping,
        periods: periods.filter(p => p['period-type-id'] !== '')
      };

      if (shipId) {
        body['ship-id'] = shipId;
      } else {
        body.ship = {
          name: shipName,
          'ship-type-id': shipTypeId,
          'imo-number': imoNumber,
          gt: Number(gt),
          flag,
          'propulsion-power': propulsionPower ? Number(propulsionPower) : undefined
        };
      }

      const response = await fetch(`${API_BASE_URL}/api/seatime`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${session.access_token}`,
        },
        body: JSON.stringify(body),
      });

      if (!response.ok) {
        const errData = await response.json();
        throw new Error(errData.message || 'Failed to record seatime');
      }

      navigate('/seatime');
    } catch (err: any) {
      console.error('Error submitting seatime:', err);
      setError(err.message || 'An error occurred while saving.');
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <Container sx={{ mt: 4, display: 'flex', justifyContent: 'center' }}>
        <CircularProgress />
      </Container>
    );
  }

  return (
    <Container maxWidth="md" sx={{ mt: 4, mb: 8 }}>
      <Button 
        component={RouterLink} 
        to="/seatime" 
        startIcon={<ArrowLeft size={18} />} 
        sx={{ mb: 3 }}
      >
        Back to History
      </Button>

      <Typography variant="h4" component="h1" gutterBottom sx={{ fontWeight: 700 }}>
        Record New Voyage
      </Typography>
      <Typography variant="body1" color="text.secondary" sx={{ mb: 4 }}>
        Fill in the details of your voyage to update your seatime records.
      </Typography>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      <form onSubmit={handleSubmit}>
        <Stack spacing={4}>
          {/* Ship Details */}
          <Paper elevation={2} sx={{ p: 3, borderRadius: 2 }}>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 3 }}>
              <Ship size={24} className="text-primary" />
              <Typography variant="h6" sx={{ fontWeight: 600 }}>Ship Details</Typography>
            </Box>
            
            <Grid container spacing={3}>
              <Grid item xs={12}>
                <TextField
                  fullWidth
                  label="Search by IMO Number (Optional)"
                  placeholder="IMO1234567"
                  value={imoNumber}
                  onChange={(e) => setImoNumber(e.target.value)}
                  helperText="Enter IMO to check if ship exists in our database"
                  InputProps={{
                    endAdornment: (
                      <IconButton size="small">
                        <Search size={18} />
                      </IconButton>
                    )
                  }}
                />
              </Grid>

              {showShipForm && (
                <>
                  <Grid item xs={12} md={8}>
                    <TextField
                      required
                      fullWidth
                      label="Ship Name"
                      value={shipName}
                      onChange={(e) => setShipName(e.target.value)}
                    />
                  </Grid>
                  <Grid item xs={12} md={4}>
                    <FormControl fullWidth required>
                      <InputLabel>Ship Type</InputLabel>
                      <Select
                        value={shipTypeId}
                        label="Ship Type"
                        onChange={(e) => setShipTypeId(e.target.value)}
                      >
                        {shipTypes.map((type) => (
                          <MenuItem key={type.id} value={type.id}>{type.name}</MenuItem>
                        ))}
                      </Select>
                    </FormControl>
                  </Grid>
                  <Grid item xs={12} md={4}>
                    <TextField
                      required
                      fullWidth
                      type="number"
                      label="Gross Tonnage (GT)"
                      value={gt}
                      onChange={(e) => setGt(e.target.value === '' ? '' : Number(e.target.value))}
                    />
                  </Grid>
                  <Grid item xs={12} md={4}>
                    <TextField
                      required
                      fullWidth
                      label="Flag State"
                      value={flag}
                      onChange={(e) => setFlag(e.target.value)}
                    />
                  </Grid>
                  <Grid item xs={12} md={4}>
                    <TextField
                      fullWidth
                      type="number"
                      label="Propulsion Power (kW)"
                      value={propulsionPower}
                      onChange={(e) => setPropulsionPower(e.target.value === '' ? '' : Number(e.target.value))}
                    />
                  </Grid>
                </>
              )}
            </Grid>
          </Paper>

          {/* Voyage Details */}
          <Paper elevation={2} sx={{ p: 3, borderRadius: 2 }}>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 3 }}>
              <Calendar size={24} className="text-primary" />
              <Typography variant="h6" sx={{ fontWeight: 600 }}>Voyage Details</Typography>
            </Box>

            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <FormControl fullWidth required>
                  <InputLabel>Voyage Type</InputLabel>
                  <Select
                    value={voyageTypeId}
                    label="Voyage Type"
                    onChange={(e) => setVoyageTypeId(e.target.value)}
                  >
                    {voyageTypes.map((type) => (
                      <MenuItem key={type.id} value={type.id}>{type.description}</MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>
              <Grid item xs={12} md={6}>
                <FormControlLabel
                  control={
                    <Switch 
                      checked={isWatchkeeping} 
                      onChange={(e) => setIsWatchkeeping(e.target.checked)} 
                    />
                  }
                  label="Watchkeeping Service"
                />
              </Grid>
              
              <Grid item xs={12} md={6}>
                <TextField
                  required
                  fullWidth
                  type="date"
                  label="Start Date"
                  value={startDate}
                  onChange={(e) => setStartDate(e.target.value)}
                  InputLabelProps={{ shrink: true }}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  required
                  fullWidth
                  type="date"
                  label="End Date"
                  value={endDate}
                  onChange={(e) => setEndDate(e.target.value)}
                  InputLabelProps={{ shrink: true }}
                />
              </Grid>
              
              <Grid item xs={12} md={6}>
                <TextField
                  required
                  fullWidth
                  label="Start Location"
                  placeholder="e.g. GB SOU"
                  value={startLocation}
                  onChange={(e) => setStartLocation(e.target.value)}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  required
                  fullWidth
                  label="End Location"
                  placeholder="e.g. US NYC"
                  value={endLocation}
                  onChange={(e) => setEndLocation(e.target.value)}
                />
              </Grid>

              <Grid item xs={12} md={6}>
                <TextField
                  required
                  fullWidth
                  label="Company / Employer"
                  value={company}
                  onChange={(e) => setCompany(e.target.value)}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  required
                  fullWidth
                  label="Capacity / Rank"
                  value={capacity}
                  onChange={(e) => setCapacity(e.target.value)}
                />
              </Grid>

              <Grid item xs={12}>
                <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'flex-end', p: 2, bgcolor: 'primary.light', borderRadius: 1, color: 'primary.contrastText' }}>
                  <Typography variant="h6" sx={{ fontWeight: 700 }}>
                    Total Calculated Days: {totalDays}
                  </Typography>
                </Box>
              </Grid>
            </Grid>
          </Paper>

          {/* Specialized Periods */}
          <Box>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <Clock size={24} className="text-primary" />
                <Typography variant="h6" sx={{ fontWeight: 600 }}>Specialized Service Periods</Typography>
              </Box>
              <Button 
                startIcon={<Plus size={18} />} 
                onClick={handleAddPeriod}
                variant="outlined"
                size="small"
              >
                Add Period
              </Button>
            </Box>

            {periods.length === 0 ? (
              <Typography variant="body2" color="text.secondary" sx={{ fontStyle: 'italic' }}>
                No specialized periods (Polar, DP, etc.) added to this voyage.
              </Typography>
            ) : (
              <Stack spacing={2}>
                {periods.map((period, index) => (
                  <Card key={index} variant="outlined">
                    <CardContent sx={{ position: 'relative', pt: 4 }}>
                      <IconButton 
                        size="small" 
                        color="error" 
                        onClick={() => handleRemovePeriod(index)}
                        sx={{ position: 'absolute', top: 8, right: 8 }}
                      >
                        <Trash2 size={18} />
                      </IconButton>
                      
                      <Grid container spacing={2}>
                        <Grid item xs={12} md={4}>
                          <FormControl fullWidth required size="small">
                            <InputLabel>Period Type</InputLabel>
                            <Select
                              value={period['period-type-id']}
                              label="Period Type"
                              onChange={(e) => handlePeriodChange(index, 'period-type-id', e.target.value)}
                            >
                              {periodTypes.map((type) => (
                                <MenuItem key={type.id} value={type.id}>{type.name}</MenuItem>
                              ))}
                            </Select>
                          </FormControl>
                        </Grid>
                        <Grid item xs={12} md={4}>
                          <TextField
                            required
                            fullWidth
                            size="small"
                            type="date"
                            label="Start Date"
                            value={period['start-date']}
                            onChange={(e) => handlePeriodChange(index, 'start-date', e.target.value)}
                            InputLabelProps={{ shrink: true }}
                          />
                        </Grid>
                        <Grid item xs={12} md={4}>
                          <TextField
                            required
                            fullWidth
                            size="small"
                            type="date"
                            label="End Date"
                            value={period['end-date']}
                            onChange={(e) => handlePeriodChange(index, 'end-date', e.target.value)}
                            InputLabelProps={{ shrink: true }}
                          />
                        </Grid>
                        <Grid item xs={12} md={9}>
                          <TextField
                            fullWidth
                            size="small"
                            label="Remarks"
                            value={period.remarks}
                            onChange={(e) => handlePeriodChange(index, 'remarks', e.target.value)}
                          />
                        </Grid>
                        <Grid item xs={12} md={3} sx={{ display: 'flex', alignItems: 'center' }}>
                          <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>
                            Days: {period.days}
                          </Typography>
                        </Grid>
                      </Grid>
                    </CardContent>
                  </Card>
                ))}
              </Stack>
            )}
          </Box>

          <Box sx={{ mt: 4, display: 'flex', gap: 2 }}>
            <Button
              type="submit"
              variant="contained"
              size="large"
              disabled={submitting}
              startIcon={submitting ? <CircularProgress size={20} /> : <Save size={20} />}
              sx={{ flexGrow: 1, py: 1.5 }}
            >
              {submitting ? 'Saving...' : 'Record Voyage'}
            </Button>
            <Button
              variant="outlined"
              size="large"
              onClick={() => navigate('/seatime')}
              sx={{ py: 1.5 }}
            >
              Cancel
            </Button>
          </Box>
        </Stack>
      </form>
    </Container>
  );
};

export default AddSeatime;
