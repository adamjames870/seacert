import { useEffect, useState, useMemo } from 'react';
import {
  Typography,
  Container,
  Box,
  Paper,
  Grid,
  Button,
  CircularProgress,
  Alert,
  Divider,
  Stack,
  IconButton,
  Tooltip,
} from '@mui/material';
import {
  ShieldCheck,
  Ship,
  History,
  AlertTriangle,
  Plus,
  ArrowRight,
} from 'lucide-react';
import { Link as RouterLink } from 'react-router-dom';
import { supabase } from '../supabaseClient';
import { API_BASE_URL } from '../config';
import { formatDate, calculateDaysInYear } from '../utils/dateUtils';

interface Certificate {
  id: string;
  'cert-type-name': string;
  'expiry-date': string;
  'document-path'?: string;
  'has-successors'?: boolean;
}

interface SeatimeRecord {
  id: string;
  'total-days': number;
  'start-date': string;
  'end-date': string;
  ship: { name: string };
  periods: { 'period-type-name': string; days: number }[];
}

interface Ship {
  id: string;
  status: 'approved' | 'provisional';
}

interface CertType {
  id: string;
  status?: 'approved' | 'provisional';
}

const Dashboard = () => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [certificates, setCertificates] = useState<Certificate[]>([]);
  const [seatime, setSeatime] = useState<SeatimeRecord[]>([]);
  const [ships, setShips] = useState<Ship[]>([]);
  const [certTypes, setCertTypes] = useState<CertType[]>([]);
  const [isAdmin, setIsAdmin] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      setError(null);
      try {
        const { data: { session } } = await supabase.auth.getSession();
        if (!session) {
          setError('Not authenticated');
          setLoading(false);
          return;
        }

        const adminStatus = session.user?.app_metadata?.role === 'admin';
        setIsAdmin(adminStatus);

        const headers = { 'Authorization': `Bearer ${session.access_token}` };

        const requests: Promise<any>[] = [
          fetch(`${API_BASE_URL}/api/certificates`, { headers }).then(r => r.json()),
          fetch(`${API_BASE_URL}/api/seatime`, { headers }).then(r => r.json()),
        ];

        if (adminStatus) {
          requests.push(fetch(`${API_BASE_URL}/api/ships`, { headers }).then(r => r.json()));
          requests.push(fetch(`${API_BASE_URL}/api/cert-types`, { headers }).then(r => r.json()));
        }

        const [certsData, seatimeData, shipsData, certTypesData] = await Promise.all(requests);

        setCertificates(certsData || []);
        setSeatime(seatimeData || []);
        if (adminStatus) {
          setShips(shipsData || []);
          setCertTypes(certTypesData || []);
        }
      } catch (err: any) {
        console.error('Error fetching dashboard data:', err);
        setError('Failed to load dashboard data');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  const stats = useMemo(() => {
    const now = new Date();
    const ninetyDaysFromNow = new Date();
    ninetyDaysFromNow.setDate(now.getDate() + 90);

    const activeCerts = certificates.filter(c => !c['has-successors']);

    const expiringSoon = activeCerts.filter(c => {
      if (!c['expiry-date']) return false;
      const expiry = new Date(c['expiry-date']);
      return expiry > now && expiry <= ninetyDaysFromNow;
    });

    const missingAttachments = activeCerts.filter(c => !c['document-path']);

    const sortedSeatime = [...seatime].sort((a, b) => 
      new Date(b['end-date']).getTime() - new Date(a['end-date']).getTime()
    );
    const lastPeriod = sortedSeatime[0];

    const currentYear = now.getFullYear();
    const lastYear = currentYear - 1;

    const seadaysThisYear = seatime
      .reduce((sum, r) => sum + calculateDaysInYear(r['start-date'], r['end-date'], currentYear), 0);

    const seadaysLastYear = seatime
      .reduce((sum, r) => sum + calculateDaysInYear(r['start-date'], r['end-date'], lastYear), 0);

    const unapprovedShips = ships.filter(s => s.status === 'provisional').length;
    const unapprovedCertTypes = certTypes.filter(ct => ct.status === 'provisional').length;

    return {
      totalCerts: activeCerts.length,
      expiringSoon: expiringSoon.length,
      missingAttachments: missingAttachments.length,
      lastPeriod,
      numPeriods: seatime.length,
      seadaysThisYear,
      seadaysLastYear,
      unapprovedShips,
      unapprovedCertTypes,
    };
  }, [certificates, seatime, ships, certTypes]);

  if (loading) {
    return (
      <Container sx={{ mt: 8, textAlign: 'center' }}>
        <CircularProgress size={60} />
        <Typography sx={{ mt: 2 }}>Loading your dashboard...</Typography>
      </Container>
    );
  }

  if (error) {
    return (
      <Container sx={{ mt: 4 }}>
        <Alert severity="error">{error}</Alert>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 8 }}>
      <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'flex-end' }}>
        <Box>
          <Typography variant="h4" component="h1" sx={{ fontWeight: 800, color: 'primary.main' }}>
            Dashboard
          </Typography>
        </Box>
      </Box>

      <Grid container spacing={4} justifyContent="center" alignItems="stretch" sx={{ px: { xs: 2, md: 0 } }}>
        {/* Certificates Section */}
        <Grid size={{ xs: 12, md: 5 }} sx={{ display: 'flex' }}>
          <Paper sx={{ p: 4, width: '100%', minWidth: 0, borderRadius: 4, boxShadow: '0 4px 25px rgba(0,0,0,0.06)', display: 'flex', flexDirection: 'column' }}>
            <Box sx={{ flexGrow: 1 }}>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 3 }}>
                <Box 
                  component={RouterLink} 
                  to="/certificates" 
                  sx={{ 
                    display: 'flex', 
                    alignItems: 'center', 
                    gap: 1.5, 
                    textDecoration: 'none', 
                    color: 'inherit',
                    '&:hover': {
                      opacity: 0.8,
                      '& .header-icon-box': {
                        bgcolor: 'primary.main',
                        color: 'primary.contrastText',
                      }
                    },
                    transition: 'all 0.2s'
                  }}
                >
                  <Box className="header-icon-box" sx={{ p: 1, bgcolor: 'primary.light', borderRadius: 2, color: 'primary.main', display: 'flex', transition: 'all 0.2s' }}>
                    <ShieldCheck size={26} />
                  </Box>
                  <Typography variant="h6" sx={{ fontWeight: 800 }}>Certificates</Typography>
                </Box>
                <Tooltip title="Smart Add Certificate">
                  <IconButton 
                    component={RouterLink} 
                    to="/certificate-wizard"
                    size="small"
                    sx={{ 
                      bgcolor: 'background.paper', 
                      color: 'primary.main',
                      border: '1px solid',
                      borderColor: 'primary.main',
                      '&:hover': {
                        bgcolor: 'primary.main',
                        color: 'primary.contrastText',
                      }
                    }}
                  >
                    <Plus size={20} />
                  </IconButton>
                </Tooltip>
              </Box>
              
              {stats.totalCerts > 0 ? (
                <Stack spacing={3}>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <Typography variant="body1" color="text.secondary" sx={{ fontWeight: 600 }}>Total Certificates</Typography>
                    <Typography variant="h4" sx={{ fontWeight: 800 }}>{stats.totalCerts}</Typography>
                  </Box>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <Typography variant="body1" color="text.secondary" sx={{ fontWeight: 600 }}>Missing Files</Typography>
                    <Typography variant="h4" sx={{ fontWeight: 800, color: stats.missingAttachments > 0 ? 'warning.main' : 'inherit' }}>
                      {stats.missingAttachments}
                    </Typography>
                  </Box>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <Typography variant="body1" color="text.secondary" sx={{ fontWeight: 600 }}>Expiring Soon</Typography>
                    <Typography variant="h4" sx={{ fontWeight: 800, color: stats.expiringSoon > 0 ? 'error.main' : 'inherit' }}>
                      {stats.expiringSoon}
                    </Typography>
                  </Box>
                </Stack>
              ) : (
                <Box sx={{ py: 4, textAlign: 'center' }}>
                  <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
                    You haven't added any certificates yet.
                  </Typography>
                  <Stack spacing={2} sx={{ alignItems: 'center' }}>
                    <Button 
                      variant="contained" 
                      color="primary"
                      startIcon={<Plus size={20} />} 
                      component={RouterLink} 
                      to="/certificate-wizard"
                      size="large"
                      sx={{ borderRadius: 2, px: 4, width: '100%', maxWidth: 300 }}
                    >
                      Use Smart Add
                    </Button>
                    <Typography variant="caption" color="text.secondary">
                      Upload an image or PDF and we'll extract the data for you!
                    </Typography>
                    <Button 
                      variant="outlined" 
                      color="primary"
                      startIcon={<Plus size={20} />} 
                      component={RouterLink} 
                      to="/add-certificate"
                      sx={{ borderRadius: 2, px: 4, width: '100%', maxWidth: 300 }}
                    >
                      Add Manually
                    </Button>
                  </Stack>
                </Box>
              )}
            </Box>

            {stats.totalCerts > 0 && (
              <>
                <Divider sx={{ my: 3 }} />
                <Box sx={{ display: 'flex', justifyContent: 'center' }}>
                  <Button 
                    component={RouterLink} 
                    to="/certificates" 
                    endIcon={<ArrowRight size={18} />}
                    sx={{ fontWeight: 700 }}
                  >
                    View All Certificates
                  </Button>
                </Box>
              </>
            )}
          </Paper>
        </Grid>

        {/* Seatime Section */}
        <Grid size={{ xs: 12, md: 5 }} sx={{ display: 'flex' }}>
          <Paper sx={{ p: 4, width: '100%', minWidth: 0, borderRadius: 4, boxShadow: '0 4px 25px rgba(0,0,0,0.06)', display: 'flex', flexDirection: 'column' }}>
            <Box sx={{ flexGrow: 1 }}>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 3 }}>
                <Box 
                  component={RouterLink} 
                  to="/seatime" 
                  sx={{ 
                    display: 'flex', 
                    alignItems: 'center', 
                    gap: 1.5, 
                    textDecoration: 'none', 
                    color: 'inherit',
                    '&:hover': {
                      opacity: 0.8,
                      '& .header-icon-box-seatime': {
                        bgcolor: 'secondary.main',
                        color: 'secondary.contrastText',
                      }
                    },
                    transition: 'all 0.2s'
                  }}
                >
                  <Box className="header-icon-box-seatime" sx={{ p: 1, bgcolor: 'secondary.light', borderRadius: 2, color: 'secondary.main', display: 'flex', transition: 'all 0.2s' }}>
                    <Ship size={26} />
                  </Box>
                  <Typography variant="h6" sx={{ fontWeight: 800 }}>Seatime</Typography>
                </Box>
                <Tooltip title="Add Seatime Period">
                  <IconButton 
                    component={RouterLink} 
                    to="/add-seatime"
                    size="small"
                    sx={{ 
                      bgcolor: 'background.paper', 
                      color: 'secondary.main',
                      border: '1px solid',
                      borderColor: 'secondary.main',
                      '&:hover': {
                        bgcolor: 'secondary.main',
                        color: 'secondary.contrastText',
                      }
                    }}
                  >
                    <Plus size={20} />
                  </IconButton>
                </Tooltip>
              </Box>

              {stats.numPeriods > 0 ? (
                <Stack spacing={3}>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <Typography variant="body1" color="text.secondary" sx={{ fontWeight: 600 }}>Days this year ({new Date().getFullYear()})</Typography>
                    <Typography variant="h4" sx={{ fontWeight: 800 }}>{stats.seadaysThisYear}</Typography>
                  </Box>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <Typography variant="body1" color="text.secondary" sx={{ fontWeight: 600 }}>Days last year ({new Date().getFullYear() - 1})</Typography>
                    <Typography variant="h4" sx={{ fontWeight: 800 }}>{stats.seadaysLastYear}</Typography>
                  </Box>
                  
                  <Box>
                    <Typography variant="subtitle2" sx={{ fontWeight: 700, mb: 1, display: 'flex', alignItems: 'center', gap: 1, color: 'text.secondary' }}>
                      <History size={16} /> LAST PERIOD
                    </Typography>
                    <Box sx={{ bgcolor: 'action.hover', p: 2, borderRadius: 2 }}>
                      <Typography variant="body2" sx={{ fontWeight: 700 }}>{stats.lastPeriod.ship?.name}</Typography>
                      <Typography variant="caption" color="text.secondary" sx={{ fontWeight: 600 }}>
                        {formatDate(stats.lastPeriod['start-date'])} - {formatDate(stats.lastPeriod['end-date'])} ({stats.lastPeriod['total-days']} days)
                      </Typography>
                    </Box>
                  </Box>
                </Stack>
              ) : (
                <Box sx={{ py: 4, textAlign: 'center' }}>
                  <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
                    You haven't recorded any seatime yet.
                  </Typography>
                  <Button 
                    variant="contained" 
                    color="primary"
                    startIcon={<Plus size={20} />} 
                    component={RouterLink} 
                    to="/add-seatime"
                    size="large"
                    sx={{ borderRadius: 2, px: 4 }}
                  >
                    Add First Seatime Period
                  </Button>
                </Box>
              )}
            </Box>

            {stats.numPeriods > 0 && (
              <>
                <Divider sx={{ my: 3 }} />
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <Typography variant="body2" color="text.secondary" sx={{ fontWeight: 600 }}>
                    {stats.numPeriods} total periods
                  </Typography>
                  <Button 
                    component={RouterLink} 
                    to="/seatime" 
                    endIcon={<ArrowRight size={18} />}
                    sx={{ fontWeight: 700 }}
                  >
                    Full History
                  </Button>
                </Box>
              </>
            )}
          </Paper>
        </Grid>


        {/* Admin Section */}
        {isAdmin && (
          <Grid size={{ xs: 12, md: 10 }} sx={{ display: 'flex' }}>
            <Paper 
              sx={{ 
                p: 4, 
                width: '100%',
                minWidth: 0,
                borderRadius: 4, 
                border: '1px solid',
                borderColor: 'warning.light',
                bgcolor: 'rgba(255, 152, 0, 0.02)',
                boxShadow: '0 4px 20px rgba(255, 152, 0, 0.05)'
              }}
            >
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 3, gap: 1.5 }}>
                <Box sx={{ p: 1, bgcolor: 'warning.light', borderRadius: 2, color: 'warning.main', display: 'flex' }}>
                  <AlertTriangle size={26} />
                </Box>
                <Typography variant="h6" sx={{ fontWeight: 800 }}>Admin Oversight</Typography>
              </Box>

              <Grid container spacing={6} justifyContent="center">
                <Grid size={{ xs: 12, sm: 5 }} sx={{ textAlign: 'center' }}>
                  <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, alignItems: 'center' }}>
                    <Box>
                      <Typography variant="h3" sx={{ fontWeight: 800, color: 'warning.main' }}>{stats.unapprovedShips}</Typography>
                      <Typography variant="body1" color="text.secondary" sx={{ fontWeight: 600 }}>Unapproved Ships</Typography>
                    </Box>
                    <Box>
                      <Button 
                        variant="contained" 
                        color="warning" 
                        size="medium"
                        component={RouterLink}
                        to="/ships"
                        sx={{ borderRadius: 2, px: 3 }}
                      >
                        Manage Ships
                      </Button>
                    </Box>
                  </Box>
                </Grid>
                <Grid size={{ xs: 12, sm: 5 }} sx={{ textAlign: 'center' }}>
                  <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, alignItems: 'center' }}>
                    <Box>
                      <Typography variant="h3" sx={{ fontWeight: 800, color: 'warning.main' }}>{stats.unapprovedCertTypes}</Typography>
                      <Typography variant="body1" color="text.secondary" sx={{ fontWeight: 600 }}>Unapproved Cert Types</Typography>
                    </Box>
                    <Box>
                      <Button 
                        variant="contained" 
                        color="warning" 
                        size="medium"
                        component={RouterLink}
                        to="/cert-types"
                        sx={{ borderRadius: 2, px: 3 }}
                      >
                        Manage Types
                      </Button>
                    </Box>
                  </Box>
                </Grid>
              </Grid>
            </Paper>
          </Grid>
        )}
      </Grid>
    </Container>
  );
};

export default Dashboard;
