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
  Alert,
  CircularProgress,
  IconButton,
  Stack,
  Breadcrumbs,
  Link,
} from '@mui/material';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import { useParams, useNavigate, Link as RouterLink } from 'react-router-dom';
import { supabase } from '../supabaseClient';
import { API_BASE_URL } from '../config';

interface Certificate {
  id: string;
  'cert-type-name': string;
  'cert-type-short-name': string;
  'issuer-name': string;
  'issuer-country': string;
}

interface User {
  id: string;
  forename: string;
  surname: string;
}

const AdminUserCertificates = () => {
  const { userId } = useParams<{ userId: string }>();
  const [certificates, setCertificates] = useState<Certificate[]>([]);
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        const { data: { session } } = await supabase.auth.getSession();
        if (!session) {
          setError('Not authenticated');
          return;
        }

        // Fetch user details
        const userResponse = await fetch(`${API_BASE_URL}/admin/users`, {
            headers: {
              'Authorization': `Bearer ${session.access_token}`,
            },
        });
        
        if (userResponse.ok) {
            const users = await userResponse.json();
            const foundUser = users.find((u: User) => u.id === userId);
            if (foundUser) setUser(foundUser);
        }

        // Fetch user certificates
        // NOTE: We assume there's an endpoint that allows admin to see a specific user's certificates
        // based on existing patterns, it might be /api/certificates?userId=... or similar
        // If not, we might need a specific admin endpoint.
        const response = await fetch(`${API_BASE_URL}/admin/users/${userId}/certificates`, {
          headers: {
            'Authorization': `Bearer ${session.access_token}`,
          },
        });

        if (!response.ok) {
          throw new Error('Failed to fetch user certificates');
        }

        const data = await response.json();
        setCertificates(data);
      } catch (err: any) {
        setError(err.message || 'An error occurred while fetching data');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [userId]);

  return (
    <Container maxWidth="lg">
      <Box sx={{ mt: 4, mb: 4 }}>
        <Breadcrumbs sx={{ mb: 2 }}>
          <Link component={RouterLink} to="/admin/users" color="inherit" underline="hover">
            User Management
          </Link>
          <Typography color="text.primary">
            {user ? `${user.forename} ${user.surname}'s Certificates` : 'User Certificates'}
          </Typography>
        </Breadcrumbs>

        <Stack
          direction="row"
          spacing={2}
          alignItems="center"
          sx={{ mb: 3 }}
        >
          <IconButton onClick={() => navigate('/admin/users')} color="primary">
            <ArrowBackIcon />
          </IconButton>
          <Typography variant="h4" component="h1" sx={{ fontWeight: 700 }}>
            {user ? `${user.forename} ${user.surname}` : 'User'}'s Certificates
          </Typography>
        </Stack>

        {loading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', mt: 8 }}>
            <CircularProgress />
          </Box>
        ) : error ? (
          <Alert severity="error" sx={{ mb: 3 }}>{error}</Alert>
        ) : (
          <TableContainer component={Paper} elevation={0} sx={{ border: 1, borderColor: 'divider' }}>
            <Table>
              <TableHead sx={{ bgcolor: 'action.hover' }}>
                <TableRow>
                  <TableCell sx={{ fontWeight: 700 }}>Certificate Type</TableCell>
                  <TableCell sx={{ fontWeight: 700 }}>Short Name</TableCell>
                  <TableCell sx={{ fontWeight: 700 }}>Issuer</TableCell>
                  <TableCell sx={{ fontWeight: 700 }}>Country</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {certificates.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={4} align="center" sx={{ py: 4 }}>
                      <Typography variant="body1" color="text.secondary">
                        No certificates found for this user
                      </Typography>
                    </TableCell>
                  </TableRow>
                ) : (
                  certificates.map((cert) => (
                    <TableRow key={cert.id} hover>
                      <TableCell sx={{ fontWeight: 500 }}>{cert['cert-type-name']}</TableCell>
                      <TableCell>{cert['cert-type-short-name'] || 'N/A'}</TableCell>
                      <TableCell>{cert['issuer-name']}</TableCell>
                      <TableCell>{cert['issuer-country']}</TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </TableContainer>
        )}
      </Box>
    </Container>
  );
};

export default AdminUserCertificates;
