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
  Checkbox,
  IconButton, 
  Divider, 
  Stack,
  Alert,
  CircularProgress,
  Autocomplete,
  Chip
} from '@mui/material';
import { 
  Trash2, 
  Plus, 
  Ship, 
  Calendar, 
  Save,
  ArrowLeft,
  Clock
} from 'lucide-react';
import { useNavigate, Link as RouterLink } from 'react-router-dom';
import { supabase } from '../supabaseClient';
import { API_BASE_URL } from '../config';

interface ShipType {
  id: string;
  name: string;
  description: string;
}

interface SeatimePeriodType {
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

interface ShipRecord {
  id: string;
  name: string;
  'ship-type-id': string;
  'ship-type-name': string;
  'imo-number': string;
  gt: number;
  flag: string;
  'propulsion-power': number;
  status: 'approved' | 'provisional';
}

interface ShipType {
  id: string;
  name: string;
  description: string;
}

const AddSeatime = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  // Lookups
  const [shipTypes, setShipTypes] = useState<ShipType[]>([]);
  const [seatimePeriodTypes, setSeatimePeriodTypes] = useState<SeatimePeriodType[]>([]);
  const [periodTypes, setPeriodTypes] = useState<PeriodType[]>([]);
  const [ships, setShips] = useState<ShipRecord[]>([]);

  const getShipTypeDescription = (ship: ShipRecord) => {
    const type = shipTypes.find(t => t.id === ship['ship-type-id']);
    return type ? type.description : ship['ship-type-name'];
  };
  
  // Form State
  const [selectedShip, setSelectedShip] = useState<ShipRecord | null>(null);
  const [shipName, setShipName] = useState('');
  const [shipTypeId, setShipTypeId] = useState('');
  const [imoNumber, setImoNumber] = useState('');
  const [gt, setGt] = useState<number | ''>('');
  const [flag, setFlag] = useState('');
  const [propulsionPower, setPropulsionPower] = useState<number | ''>('');
  
  const [seatimePeriodTypeId, setSeatimePeriodTypeId] = useState('');
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [startLocation, setStartLocation] = useState('');
  const [endLocation, setEndLocation] = useState('');
  const [totalDays, setTotalDays] = useState<number>(0);
  const [company, setCompany] = useState('');
  const [capacity, setCapacity] = useState('');
  const [isWatchkeeping, setIsWatchkeeping] = useState(false);
  const [periods, setPeriods] = useState<SpecializedPeriod[]>([]);
  const [showShipForm, setShowShipForm] = useState(false);

  useEffect(() => {
    fetchLookups();
    fetchShips();
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
      setSeatimePeriodTypes(data['voyage-types']);
      setPeriodTypes(data['period-types']);
    } catch (err: any) {
      console.error('Error fetching lookups:', err);
      setError('Could not load lookup data. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const fetchShips = async () => {
    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) return;

      const response = await fetch(`${API_BASE_URL}/api/ships`, {
        headers: {
          'Authorization': `Bearer ${session.access_token}`,
        },
      });

      if (!response.ok) throw new Error('Failed to fetch ships');

      const data = await response.json();
      setShips(Array.isArray(data) ? data : []);
    } catch (err) {
      console.error('Error fetching ships:', err);
      setShips([]);
    }
  };

  const handleShipSelect = (_event: any, newValue: ShipRecord | null) => {
    setSelectedShip(newValue);
    if (newValue) {
      setShowShipForm(false);
      // Optional: prepopulate if needed, but the API expects ship-id or ship object
    } else {
      setShowShipForm(true);
    }
  };

  const handleAddPeriod = () => {
    let initialDays = 0;
    if (startDate && endDate) {
      const pStart = new Date(startDate);
      const pEnd = new Date(endDate);
      if (pEnd >= pStart && !isNaN(pStart.getTime()) && !isNaN(pEnd.getTime())) {
        const diffTime = Math.abs(pEnd.getTime() - pStart.getTime());
        initialDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24)) + 1;
      }
    }

    setPeriods([
      ...periods,
      {
        'period-type-id': '',
        'start-date': startDate,
        'end-date': endDate,
        days: initialDays,
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
        'voyage-type-id': seatimePeriodTypeId,
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

      if (selectedShip) {
        body['ship-id'] = selectedShip.id;
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
        Record New Seatime Period
      </Typography>
      <Typography variant="body1" color="text.secondary" sx={{ mb: 4 }}>
        Fill in the details of your seatime period to update your seatime records.
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
            
            <Grid container spacing={2}>
              <Grid size={12}>
                <Autocomplete
                  options={ships || []}
                  getOptionLabel={(option) => option ? `${option.name} (${option['imo-number']})` : ''}
                  value={selectedShip}
                  onChange={handleShipSelect}
                  noOptionsText="No ships found. Add new ship details below."
                  fullWidth
                  renderInput={(params) => (
                    <TextField
                      {...params}
                      label="Select Existing Ship or Search IMO"
                      placeholder={ships.length === 0 ? "No existing ships found - add new details below" : "Start typing ship name or IMO..."}
                      variant="outlined"
                      fullWidth
                      InputProps={{
                        ...params.InputProps,
                        sx: { 
                          fontSize: '1.2rem', 
                          py: 1,
                          '& .MuiOutlinedInput-input': {
                            padding: '10px 14px',
                          }
                        }
                      }}
                    />
                  )}
                  renderOption={(props, option) => {
                    if (!option) return null;
                    const { key, ...optionProps } = props as any;
                    return (
                      <li key={key} {...optionProps}>
                        <Box sx={{ display: 'flex', flexDirection: 'column' }}>
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                            <Typography variant="body1" sx={{ fontWeight: 500 }}>{option.name}</Typography>
                            {option.status === 'provisional' && (
                              <Chip label="Provisional" size="small" color="warning" variant="outlined" />
                            )}
                          </Box>
                          <Typography variant="caption" color="text.secondary">
                            IMO: {option['imo-number']} | {getShipTypeDescription(option)}
                          </Typography>
                        </Box>
                      </li>
                    );
                  }}
                  filterOptions={(options, params) => {
                    const filtered = options.filter(o => 
                      o && (
                        o.name.toLowerCase().includes(params.inputValue.toLowerCase()) ||
                        o['imo-number'].toLowerCase().includes(params.inputValue.toLowerCase())
                      )
                    );
                    return filtered;
                  }}
                />
              </Grid>

              {!selectedShip && (
                <Grid size={12}>
                  <Box sx={{ display: 'flex', alignItems: 'center', my: 1 }}>
                    <Divider sx={{ flexGrow: 1 }} />
                    <Typography variant="body2" color="text.secondary" sx={{ mx: 2, fontWeight: 700 }}>OR</Typography>
                    <Divider sx={{ flexGrow: 1 }} />
                  </Box>
                  <Box sx={{ display: 'flex', justifyContent: 'center' }}>
                    <Button 
                      variant="outlined" 
                      onClick={() => setShowShipForm(!showShipForm)}
                      startIcon={showShipForm ? <Trash2 size={16} /> : <Plus size={16} />}
                      color={showShipForm ? "error" : "primary"}
                      size="small"
                      sx={{ 
                        px: 2, 
                        py: 0.5, 
                        fontSize: '0.8rem',
                        borderRadius: 2
                      }}
                    >
                      {showShipForm ? 'Cancel New Ship' : 'Add New Ship Details'}
                    </Button>
                  </Box>
                </Grid>
              )}

              {showShipForm && !selectedShip && (
                <Grid size={12}>
                  <Box sx={{ mt: 2, p: 1, borderTop: '1px solid', borderColor: 'divider' }}>
                    <Typography variant="subtitle2" gutterBottom sx={{ mb: 3, fontWeight: 600, color: 'primary.main', textTransform: 'uppercase', letterSpacing: 1 }}>New Ship Details</Typography>
                    <Grid container spacing={3}>
                      <Grid size={{ xs: 12, md: 6 }}>
                        <TextField
                          required
                          fullWidth
                          label="Ship Name"
                          placeholder="Enter vessel name"
                          value={shipName}
                          onChange={(e) => setShipName(e.target.value)}
                        />
                      </Grid>
                      <Grid size={{ xs: 12, md: 6 }}>
                        <TextField
                          required
                          fullWidth
                          label="IMO Number"
                          placeholder="e.g. IMO1234567"
                          value={imoNumber}
                          onChange={(e) => setImoNumber(e.target.value)}
                        />
                      </Grid>
                      <Grid size={{ xs: 12, md: 6 }}>
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
                      </Grid>
                      <Grid size={{ xs: 12, md: 2 }}>
                        <TextField
                          required
                          fullWidth
                          type="number"
                          label="GT"
                          placeholder="Gross Tonnage"
                          value={gt}
                          onChange={(e) => setGt(e.target.value === '' ? '' : Number(e.target.value))}
                        />
                      </Grid>
                      <Grid size={{ xs: 12, md: 2 }}>
                        <TextField
                          required
                          fullWidth
                          label="Flag"
                          placeholder="e.g. UK"
                          value={flag}
                          onChange={(e) => setFlag(e.target.value)}
                        />
                      </Grid>
                      <Grid size={{ xs: 12, md: 2 }}>
                        <TextField
                          fullWidth
                          type="number"
                          label="kW"
                          placeholder="Propulsion Power"
                          value={propulsionPower}
                          onChange={(e) => setPropulsionPower(e.target.value === '' ? '' : Number(e.target.value))}
                        />
                      </Grid>
                    </Grid>
                  </Box>
                </Grid>
              )}
            </Grid>
          </Paper>

          {/* Seatime Period Details */}
          <Paper elevation={2} sx={{ p: 3, borderRadius: 2 }}>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 3 }}>
              <Calendar size={24} className="text-primary" />
              <Typography variant="h6" sx={{ fontWeight: 600 }}>Seatime Period Details</Typography>
            </Box>

            <Grid container spacing={3} sx={{ mt: 1 }}>
              <Grid size={{ xs: 12, md: 9 }}>
                <FormControl fullWidth required>
                  <InputLabel>Seatime Period Type</InputLabel>
                  <Select
                    value={seatimePeriodTypeId}
                    label="Seatime Period Type"
                    onChange={(e) => setSeatimePeriodTypeId(e.target.value)}
                  >
                    {seatimePeriodTypes.map((type) => (
                      <MenuItem key={type.id} value={type.id}>{type.description}</MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>
              <Grid size={{ xs: 12, md: 3 }} sx={{ display: 'flex', alignItems: 'center', justifyContent: { xs: 'flex-start', md: 'flex-end' } }}>
                <FormControlLabel
                  control={
                    <Checkbox 
                      checked={isWatchkeeping} 
                      onChange={(e) => setIsWatchkeeping(e.target.checked)} 
                    />
                  }
                  label="Watchkeeping"
                />
              </Grid>

              <Grid size={{ xs: 12, md: 6 }}>
                <TextField
                  required
                  fullWidth
                  label="Capacity / Rank"
                  placeholder="e.g. Chief Officer"
                  value={capacity}
                  onChange={(e) => setCapacity(e.target.value)}
                />
              </Grid>
              <Grid size={{ xs: 12, md: 6 }}>
                <TextField
                  required
                  fullWidth
                  label="Company / Employer"
                  placeholder="e.g. Global Shipping Ltd"
                  value={company}
                  onChange={(e) => setCompany(e.target.value)}
                />
              </Grid>

              <Grid size={{ xs: 12, md: 6 }}>
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
              <Grid size={{ xs: 12, md: 6 }}>
                <TextField
                  required
                  fullWidth
                  label="Start Location"
                  placeholder="e.g. Southampton"
                  value={startLocation}
                  onChange={(e) => setStartLocation(e.target.value)}
                />
              </Grid>

              <Grid size={{ xs: 12, md: 6 }}>
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
              <Grid size={{ xs: 12, md: 6 }}>
                <TextField
                  required
                  fullWidth
                  label="End Location"
                  placeholder="e.g. New York"
                  value={endLocation}
                  onChange={(e) => setEndLocation(e.target.value)}
                />
              </Grid>

              <Grid size={12}>
                <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'flex-end', p: 1.5, bgcolor: 'primary.light', borderRadius: 1, color: 'primary.contrastText' }}>
                  <Typography variant="subtitle1" sx={{ fontWeight: 700 }}>
                    Total Calculated Days: {totalDays}
                  </Typography>
                </Box>
              </Grid>
            </Grid>
          </Paper>

          {/* Specialized Periods */}
          <Paper elevation={2} sx={{ p: 3, borderRadius: 2 }}>
            <Box sx={{ 
              display: 'flex', 
              flexDirection: { xs: 'column', sm: 'row' },
              justifyContent: 'space-between', 
              alignItems: { xs: 'flex-start', sm: 'center' }, 
              gap: { xs: 2, sm: 0 },
              mb: 3 
            }}>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <Clock size={24} className="text-primary" />
                <Typography variant="h6" sx={{ fontWeight: 600 }}>Specialized Service Periods</Typography>
              </Box>
              <Button 
                startIcon={<Plus size={18} />} 
                onClick={handleAddPeriod}
                variant="outlined"
                size="small"
                sx={{ width: { xs: '100%', sm: 'auto' } }}
              >
                Add Period
              </Button>
            </Box>

            {periods.length === 0 ? (
              <Typography variant="body2" color="text.secondary" sx={{ fontStyle: 'italic' }}>
                No specialized periods (Polar, DP, etc.) added to this seatime period.
              </Typography>
            ) : (
              <Stack spacing={3}>
                {periods.map((period, index) => (
                  <Box key={index} sx={{ p: 2, border: '1px solid', borderColor: 'divider', borderRadius: 1, position: 'relative' }}>
                    <Grid container spacing={3} alignItems="center">
                      <Grid size={{ xs: 12, md: 4 }}>
                        <FormControl fullWidth required size="small">
                          <InputLabel>Period Type</InputLabel>
                          <Select
                            value={period['period-type-id']}
                            label="Period Type"
                            onChange={(e) => handlePeriodChange(index, 'period-type-id', e.target.value)}
                          >
                            {periodTypes.map((type) => (
                              <MenuItem key={type.id} value={type.id}>{type.description}</MenuItem>
                            ))}
                          </Select>
                        </FormControl>
                      </Grid>
                      <Grid size={{ xs: 12, md: 3 }}>
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
                      <Grid size={{ xs: 12, md: 3 }}>
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
                      <Grid size={{ xs: 12, md: 2 }} sx={{ display: 'flex', alignItems: 'center', justifyContent: 'flex-end' }}>
                        <Typography variant="subtitle2" sx={{ fontWeight: 700, mr: 2 }}>
                          {period.days} Days
                        </Typography>
                        <IconButton 
                          size="small" 
                          color="error" 
                          onClick={() => handleRemovePeriod(index)}
                        >
                          <Trash2 size={18} />
                        </IconButton>
                      </Grid>
                      <Grid size={12}>
                        <TextField
                          fullWidth
                          size="small"
                          label="Remarks"
                          placeholder="Optional remarks"
                          value={period.remarks}
                          onChange={(e) => handlePeriodChange(index, 'remarks', e.target.value)}
                        />
                      </Grid>
                    </Grid>
                  </Box>
                ))}
              </Stack>
            )}
          </Paper>

          <Box sx={{ mt: 4, display: 'flex', gap: 2 }}>
            <Button
              type="submit"
              variant="contained"
              size="large"
              disabled={submitting}
              startIcon={submitting ? <CircularProgress size={20} /> : <Save size={20} />}
              sx={{ flexGrow: 1, py: 1.5 }}
            >
              {submitting ? 'Saving...' : 'Record Seatime Period'}
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
