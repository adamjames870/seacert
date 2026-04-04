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
  Clock,
  Save,
  ArrowLeft
} from 'lucide-react';
import { useNavigate, useParams, Link as RouterLink } from 'react-router-dom';
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

interface SpecializedPeriod {
  id?: string;
  'period-type-id': string;
  'start-date': string;
  'end-date': string;
  days: number;
  remarks: string;
}

const UpdateSeatime = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  // Lookups
  const [shipTypes, setShipTypes] = useState<ShipType[]>([]);
  const [seatimePeriodTypes, setSeatimePeriodTypes] = useState<SeatimePeriodType[]>([]);
  const [periodTypes, setPeriodTypes] = useState<PeriodType[]>([]);
  const [ships, setShips] = useState<ShipRecord[]>([]);
  
  // Form State
  const [selectedShip, setSelectedShip] = useState<ShipRecord | null>(null);
  
  const [seatimePeriodTypeId, setSeatimePeriodTypeId] = useState('');
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [startLocation, setStartLocation] = useState('');
  const [endLocation, setEndLocation] = useState('');
  const [totalDays, setTotalDays] = useState<number>(0);
  const [company, setCompany] = useState('');
  const [capacity, setCapacity] = useState('');
  const [isWatchkeeping, setIsWatchkeeping] = useState(true);
  const [periods, setPeriods] = useState<SpecializedPeriod[]>([]);

  useEffect(() => {
    fetchLookupsAndRecord();
  }, [id]);

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

  const fetchLookupsAndRecord = async () => {
    setLoading(true);
    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) return;

      // Fetch Lookups
      const lookupResp = await fetch(`${API_BASE_URL}/api/seatime/lookups`, {
        headers: { 'Authorization': `Bearer ${session.access_token}` },
      });
      if (!lookupResp.ok) throw new Error('Failed to fetch lookups');
      const lookups = await lookupResp.json();
      setShipTypes(lookups['ship-types']);
      setSeatimePeriodTypes(lookups['voyage-types']);
      setPeriodTypes(lookups['period-types']);

      // Fetch Ships for autocomplete
      const shipsDataResponse = await fetch(`${API_BASE_URL}/api/ships`, {
        headers: { 'Authorization': `Bearer ${session.access_token}` },
      });
      let fetchedShips: ShipRecord[] = [];
      if (shipsDataResponse.ok) {
        fetchedShips = await shipsDataResponse.json();
        setShips(fetchedShips);
      }

      // Fetch Record - Note: Backend might provide specific endpoint or use list
      const recordResp = await fetch(`${API_BASE_URL}/api/seatime/${id}`, {
        headers: { 'Authorization': `Bearer ${session.access_token}` },
      });
      
      let record;
      if (recordResp.ok) {
        record = await recordResp.json();
      } else {
        // Fallback to list search if single fetch fails
        const listResp = await fetch(`${API_BASE_URL}/api/seatime`, {
          headers: { 'Authorization': `Bearer ${session.access_token}` },
        });
        if (!listResp.ok) throw new Error('Failed to fetch seatime records');
        const records = await listResp.json();
        record = records.find((r: any) => r.id === id);
      }
      
      if (!record) throw new Error('Record not found');

      console.log('Record found:', record);

      // Populate form
      if (record.ship) {
        setSelectedShip(record.ship);
      } else if (record['ship-id']) {
        // Find ship in list if only ID provided
        const ship = fetchedShips.find(s => s.id === record['ship-id']);
        if (ship) setSelectedShip(ship);
      } else {
        setSelectedShip(null);
      }
      
      setSeatimePeriodTypeId(record['voyage-type-id'] || '');
      setStartDate(record['start-date'] ? record['start-date'].split('T')[0] : '');
      setEndDate(record['end-date'] ? record['end-date'].split('T')[0] : '');
      setStartLocation(record['start-location'] || '');
      setEndLocation(record['end-location'] || '');
      setTotalDays(record['total-days'] || 0);
      setCompany(record.company || '');
      setCapacity(record.capacity || '');
      setIsWatchkeeping(!!record['is-watchkeeping']);
      setPeriods((record.periods || []).map((p: any) => ({
        id: p.id,
        'period-type-id': p['period-type-id'],
        'start-date': p['start-date'] ? p['start-date'].split('T')[0] : '',
        'end-date': p['end-date'] ? p['end-date'].split('T')[0] : '',
        days: p.days,
        remarks: p.remarks
      })));

    } catch (err: any) {
      console.error('Error fetching data:', err);
      setError(err.message || 'Could not load record.');
    } finally {
      setLoading(false);
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

  const getShipTypeDescription = (ship: ShipRecord) => {
    const type = shipTypes.find(t => t.id === ship['ship-type-id']);
    return type?.description || ship['ship-type-name'] || 'Unknown Type';
  };

  const handlePeriodChange = (index: number, field: keyof SpecializedPeriod, value: any) => {
    const newPeriods = [...periods];
    newPeriods[index] = { ...newPeriods[index], [field]: value };
    
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

      const body = {
        id,
        'ship-id': selectedShip?.id,
        'voyage-type-id': seatimePeriodTypeId,
        'start-date': startDate,
        'end-date': endDate,
        'start-location': startLocation,
        'end-location': endLocation,
        'total-days': Number(totalDays),
        company,
        capacity,
        'is-watchkeeping': isWatchkeeping,
        periods: periods
          .filter(p => p['period-type-id'] !== '')
          .map(p => {
            const { id: _ignoredId, ...rest } = p;
            return {
              ...rest,
              days: Number(rest.days)
            };
          })
      };

      console.log('Update body:', body);

      // Note: Assuming PUT /api/seatime/:id exists as per proposal
      const response = await fetch(`${API_BASE_URL}/api/seatime/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${session.access_token}`,
        },
        body: JSON.stringify(body),
      });

      if (!response.ok) {
        const responseData = await response.json().catch(() => ({}));
        console.error('Update error response:', responseData);
        throw new Error(responseData.message || responseData.error || 'Failed to update seatime');
      }

      navigate('/seatime');
    } catch (err: any) {
      console.error('Error updating seatime:', err);
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
        Update Seatime Period
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
                  options={ships}
                  getOptionLabel={(option) => `${option.name} (${option['imo-number']})`}
                  value={selectedShip}
                  onChange={(_, newValue) => setSelectedShip(newValue)}
                  fullWidth
                  renderInput={(params) => (
                    <TextField
                      {...params}
                      label="Select Ship"
                      required
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
                />
              </Grid>

              {selectedShip && (
                <Grid size={12}>
                  <Box sx={{ mt: 2, p: 2, bgcolor: 'background.default', borderRadius: 1 }}>
                    <Typography variant="body2" color="text.secondary">
                      Type: {getShipTypeDescription(selectedShip)} | Flag: {selectedShip.flag} | GT: {selectedShip.gt}
                    </Typography>
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
                    <Switch 
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
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
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
              {submitting ? 'Updating...' : 'Update Record'}
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

export default UpdateSeatime;
