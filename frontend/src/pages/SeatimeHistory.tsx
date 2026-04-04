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
  MapPin, 
  Clock,
  ChevronRight,
  Info
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
    'imo-number': string;
    gt: number;
    flag: string;
    'propulsion-power': number;
  };
  'voyage-type-id': string;
  'voyage-type-name': string;
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

const SeatimeHistory = () => {
  const [records, setRecords] = useState<SeatimeRecord[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchSeatime();
  }, []);

  const fetchSeatime = async () => {
    setLoading(true);
    setError(null);
    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) {
        setError('No active session. Please log in.');
        setLoading(false);
        return;
      }

      const response = await fetch(`${API_BASE_URL}/api/seatime`, {
        headers: {
          'Authorization': `Bearer ${session.access_token}`,
        },
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch seatime: ${response.statusText}`);
      }

      const data = await response.json();
      setRecords(data);
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

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4 }}>
        <Typography variant="h4" component="h1" sx={{ fontWeight: 700 }}>
          Seatime History
        </Typography>
        <Button 
          variant="contained" 
          startIcon={<Plus size={20} />}
          component={RouterLink}
          to="/add-seatime"
        >
          Record New Voyage
        </Button>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      {records.length === 0 ? (
        <Paper sx={{ p: 6, textAlign: 'center', bgcolor: 'background.default', border: '2px dashed', borderColor: 'divider' }}>
          <Ship size={48} style={{ marginBottom: 16, opacity: 0.5 }} />
          <Typography variant="h6" color="text.secondary" gutterBottom>
            No seatime recorded yet
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
            Start building your maritime experience profile by logging your first voyage.
          </Typography>
          <Button 
            variant="outlined" 
            component={RouterLink}
            to="/add-seatime"
          >
            Add Your First Entry
          </Button>
        </Paper>
      ) : (
        <Stack spacing={3}>
          {records.map((record) => (
            <Card key={record.id} elevation={2} sx={{ borderRadius: 2 }}>
              <CardContent sx={{ p: 0 }}>
                <Box sx={{ p: 3 }}>
                  <Grid container spacing={2}>
                    <Grid item xs={12} md={4}>
                      <Stack spacing={1}>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          <Ship size={20} className="text-primary" />
                          <Typography variant="h6" sx={{ fontWeight: 700 }}>
                            {record.ship.name}
                          </Typography>
                        </Box>
                        <Typography variant="body2" color="text.secondary">
                          {record.ship['ship-type-name']} • IMO: {record.ship['imo-number']}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          {record.ship.gt} GT • {record.ship.flag}
                        </Typography>
                      </Stack>
                    </Grid>
                    
                    <Grid item xs={12} md={5}>
                      <Stack spacing={1}>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          <Calendar size={18} />
                          <Typography variant="body1">
                            {formatDate(record['start-date'])} to {formatDate(record['end-date'])}
                          </Typography>
                        </Box>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          <MapPin size={18} />
                          <Typography variant="body2">
                            {record['start-location']} → {record['end-location']}
                          </Typography>
                        </Box>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          <Chip 
                            label={record['voyage-type-name']} 
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

                    <Grid item xs={12} md={3} sx={{ display: 'flex', flexDirection: 'column', alignItems: { md: 'flex-end' }, justifyContent: 'center' }}>
                      <Box sx={{ textAlign: { md: 'right' } }}>
                        <Typography variant="h4" sx={{ fontWeight: 800, color: 'primary.main' }}>
                          {record['total-days']}
                        </Typography>
                        <Typography variant="overline" sx={{ fontWeight: 700 }}>
                          Total Days
                        </Typography>
                      </Box>
                      <Button 
                        size="small" 
                        startIcon={<Edit size={16} />}
                        component={RouterLink}
                        to={`/update-seatime/${record.id}`}
                        sx={{ mt: 1 }}
                      >
                        Edit Record
                      </Button>
                    </Grid>
                  </Grid>
                </Box>

                {record.periods && record.periods.length > 0 && (
                  <>
                    <Divider />
                    <Box sx={{ p: 2, bgcolor: 'rgba(0,0,0,0.02)' }}>
                      <Typography variant="subtitle2" sx={{ mb: 1, px: 1, fontWeight: 700, display: 'flex', alignItems: 'center', gap: 1 }}>
                        <Clock size={16} /> Specialized Service Periods
                      </Typography>
                      <Stack direction="row" spacing={1} sx={{ flexWrap: 'wrap', gap: 1 }}>
                        {record.periods.map((period) => (
                          <Tooltip key={period.id} title={`${formatDate(period['start-date'])} - ${formatDate(period['end-date'])}: ${period.remarks}`}>
                            <Chip 
                              label={`${period['period-type-name']}: ${period.days} days`}
                              size="small"
                              variant="filled"
                              sx={{ bgcolor: 'info.light', color: 'info.contrastText' }}
                            />
                          </Tooltip>
                        ))}
                      </Stack>
                    </Box>
                  </>
                )}
                
                <Box sx={{ p: 2, borderTop: 1, borderColor: 'divider', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <Typography variant="caption" color="text.secondary">
                    Capacity: <strong>{record.capacity}</strong> • Company: <strong>{record.company}</strong>
                  </Typography>
                </Box>
              </CardContent>
            </Card>
          ))}
        </Stack>
      )}
    </Container>
  );
};

export default SeatimeHistory;
