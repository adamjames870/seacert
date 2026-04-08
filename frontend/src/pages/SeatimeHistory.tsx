import { useEffect, useState } from 'react';
import { 
  Typography, 
  Container, 
  Box, 
  Paper, 
  Button, 
  CircularProgress, 
  Alert,
  IconButton,
  Tooltip,
  Chip,
  Stack,
  Card,
  CardContent,
  Grid,
  Divider
} from '@mui/material';
import { 
  Plus, 
  Edit, 
  Ship, 
  Calendar, 
  MapPin
} from 'lucide-react';
import { Link as RouterLink } from 'react-router-dom';
import { supabase } from '../supabaseClient';
import { API_BASE_URL } from '../config';
import { formatDate } from '../utils/dateUtils';

interface SpecializedPeriod {
  id: string;
  'period-type-id': string;
  'period-type-name': string;
  'start-date': string;
  'end-date': string;
  days: number;
  remarks: string;
}

interface SeatimeRecord {
  id: string;
  'ship-id': string;
  ship: {
    id: string;
    name: string;
    'ship-type-id': string;
    'ship-type-name': string;
    'ship-type-description'?: string;
    'imo-number': string;
    gt: number;
    flag: string;
    'propulsion-power': number;
  };
  'voyage-type-id': string;
  'voyage-type-name': string;
  'voyage-type-description'?: string;
  'seatime-period-type-description'?: string;
  'start-date': string;
  'end-date': string;
  'start-location': string;
  'end-location': string;
  'total-days': number;
  company: string;
  capacity: string;
  'is-watchkeeping': boolean;
  periods: SpecializedPeriod[];
}

interface LookupItem {
  id: string;
  name: string;
  description: string;
}

interface SeatimeLookups {
  'ship-types': LookupItem[];
  'voyage-types': LookupItem[];
  'period-types': LookupItem[];
}

const SeatimeHistory = () => {
  const [records, setRecords] = useState<SeatimeRecord[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    setError(null);
    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) {
        setError('No active session. Please log in.');
        setLoading(false);
        return;
      }

      // Fetch lookups first for descriptions
      const lookupResponse = await fetch(`${API_BASE_URL}/api/seatime/lookups`, {
        headers: {
          'Authorization': `Bearer ${session.access_token}`,
        },
      });

      let lookupData: SeatimeLookups | null = null;
      if (lookupResponse.ok) {
        lookupData = await lookupResponse.json();
      }

      const response = await fetch(`${API_BASE_URL}/api/seatime`, {
        headers: {
          'Authorization': `Bearer ${session.access_token}`,
        },
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch seatime: ${response.statusText}`);
      }

      const data: SeatimeRecord[] = await response.json();
      
      // Map descriptions if lookups are available
      const recordsWithDescriptions = (data || []).map(record => {
        const shipType = lookupData?.['ship-types'].find(t => t.id === record.ship['ship-type-id']);
        const voyageType = lookupData?.['voyage-types'].find(t => t.id === record['voyage-type-id']);
        
        return {
          ...record,
          ship: {
            ...record.ship,
            'ship-type-description': shipType?.description || record.ship['ship-type-name']
          },
          'seatime-period-type-description': voyageType?.description || record['voyage-type-name'],
          periods: (record.periods || []).map(p => {
            const periodType = lookupData?.['period-types'].find(pt => pt.id === p['period-type-id']);
            return {
              ...p,
              'period-type-description': periodType?.description || p['period-type-name']
            };
          })
        };
      });

      setRecords(recordsWithDescriptions);
    } catch (err: any) {
      console.error('Error fetching seatime:', err);
      setError(err.message || 'An unexpected error occurred');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Container sx={{ mt: 4, display: 'flex', justifyContent: 'center' }}>
        <CircularProgress />
      </Container>
    );
  }

  const isEmpty = records.length === 0;

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Box sx={{ 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center', 
        mb: 4,
        opacity: loading || error ? 0.3 : 1,
        transition: 'opacity 0.3s'
      }}>
        <Typography variant="h4" component="h1" sx={{ fontWeight: 700 }}>
          Seatime History
        </Typography>
        {!isEmpty && !loading && !error && (
          <Button 
            variant="contained" 
            startIcon={<Plus size={20} />}
            component={RouterLink}
            to="/add-seatime"
          >
            Record New Seatime Period
          </Button>
        )}
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      {isEmpty && !loading && !error ? (
        <Box sx={{ 
          display: 'flex', 
          flexDirection: 'column', 
          alignItems: 'center', 
          justifyContent: 'center', 
          minHeight: '60vh',
          textAlign: 'center'
        }}>
          <Paper sx={{ 
            p: 8, 
            display: 'flex', 
            flexDirection: 'column', 
            alignItems: 'center', 
            borderRadius: 4,
            bgcolor: 'background.paper',
            boxShadow: '0 8px 40px rgba(0,0,0,0.08)',
            border: '1px solid',
            borderColor: 'divider',
            maxWidth: 600,
            width: '100%'
          }}>
            <Box sx={{ 
              mb: 3, 
              p: 3, 
              borderRadius: '50%', 
              bgcolor: 'primary.light', 
              color: 'primary.main',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center'
            }}>
              <Ship size={64} />
            </Box>
            <Typography variant="h4" gutterBottom sx={{ fontWeight: 700 }}>
              No Seatime Records
            </Typography>
            <Typography variant="body1" color="text.secondary" sx={{ mb: 4, fontSize: '1.1rem' }}>
              You haven't recorded any sea service yet. Start logging your seatime periods to build your professional profile.
            </Typography>
            <Button 
              variant="contained" 
              size="large"
              startIcon={<Plus size={24} />}
              component={RouterLink}
              to="/add-seatime"
              sx={{ px: 4, py: 1.5, borderRadius: 2, fontSize: '1.1rem', fontWeight: 700 }}
            >
              Add First Seatime Record
            </Button>
          </Paper>
        </Box>
      ) : (
        <Stack spacing={3} width="100%">
          {records.map((record) => (
            <Card key={record.id} elevation={2} sx={{ borderRadius: 2, width: '100%' }}>
              <CardContent sx={{ p: 0 }}>
                <Box sx={{ p: 3 }}>
                  <Grid container spacing={3} alignItems="flex-start">
                    <Grid size={{ xs: 12, md: 3 }}>
                      <Stack spacing={1}>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          <Ship size={20} className="text-primary" style={{ color: '#1976d2' }} />
                          <Typography variant="h6" sx={{ fontWeight: 700 }}>
                            {record.ship.name}
                          </Typography>
                        </Box>
                        <Typography variant="body2" color="text.secondary" sx={{ fontWeight: 500 }}>
                          {record.ship['ship-type-description']}
                        </Typography>
                        <Typography variant="caption" color="text.secondary">
                          IMO: {record.ship['imo-number']} • {record.ship.gt} GT • {record.ship.flag}
                        </Typography>
                      </Stack>
                    </Grid>
                    
                    <Grid size={{ xs: 12, md: 3 }}>
                      <Stack spacing={1}>
                        <Typography variant="subtitle1" sx={{ fontWeight: 800, color: 'primary.main' }}>
                          {record.capacity}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          {record.company}
                        </Typography>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, opacity: 0.8 }}>
                          <MapPin size={18} color="#666" />
                          <Typography variant="body2">
                            {record['start-location']} → {record['end-location']}
                          </Typography>
                        </Box>
                        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                          <Chip 
                            label={record['seatime-period-type-description']} 
                            size="small" 
                            color="primary" 
                            variant="outlined" 
                          />
                          {record['is-watchkeeping'] && (
                            <Chip label="Watchkeeping" size="small" color="success" variant="outlined" />
                          )}
                        </Box>
                      </Stack>
                    </Grid>

                    <Grid size={{ xs: 12, md: 4 }} sx={{ display: 'flex', flexDirection: 'column', alignItems: { md: 'flex-end' }, justifyContent: 'center', ml: 'auto' }}>
                      <Stack spacing={0.5} alignItems={{ md: 'flex-end' }}>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          <Calendar size={18} color="#666" />
                          <Typography variant="body1" sx={{ fontWeight: 600 }}>
                            {formatDate(record['start-date'])} — {formatDate(record['end-date'])}
                          </Typography>
                        </Box>
                        <Box sx={{ textAlign: { md: 'right' } }}>
                          <Typography variant="h4" sx={{ fontWeight: 800, color: 'primary.dark', lineHeight: 1 }}>
                            {record['total-days']}
                          </Typography>
                          <Typography variant="overline" sx={{ fontWeight: 700, lineHeight: 1 }}>
                            Total Days
                          </Typography>
                        </Box>
                      </Stack>
                    </Grid>

                    <Grid size={{ xs: 12, md: 1 }} sx={{ display: 'flex', flexDirection: 'column', alignItems: { xs: 'center', md: 'flex-end' }, justifyContent: 'center' }}>
                      <IconButton 
                        size="small" 
                        color="primary"
                        component={RouterLink}
                        to={`/update-seatime/${record.id}`}
                        sx={{ border: '1px solid', borderColor: 'primary.light' }}
                      >
                        <Edit size={18} />
                      </IconButton>
                    </Grid>
                  </Grid>
                </Box>

                {record.periods && record.periods.length > 0 && (
                  <>
                    <Divider />
                    <Box sx={{ p: 2, bgcolor: 'grey.50' }}>
                      <Stack direction="row" spacing={1} sx={{ flexWrap: 'wrap', gap: 1 }}>
                        {record.periods.map((period: any) => (
                          <Tooltip key={period.id} title={`${formatDate(period['start-date'])} - ${formatDate(period['end-date'])}: ${period.remarks}`}>
                            <Chip 
                              label={`${period['period-type-description']}: ${period.days} days`}
                              size="small"
                              variant="filled"
                              sx={{ bgcolor: 'info.light', color: 'info.contrastText', fontWeight: 600 }}
                            />
                          </Tooltip>
                        ))}
                      </Stack>
                    </Box>
                  </>
                )}
              </CardContent>
            </Card>
          ))}
        </Stack>
      )}
    </Container>
  );
};

export default SeatimeHistory;
