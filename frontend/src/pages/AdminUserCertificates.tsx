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
            const userData = await userResponse.json();
            let usersArray: User[] = [];
            
            if (Array.isArray(userData)) {
              usersArray = userData;
            } else if (userData && typeof userData === 'object') {
              const arrayKey = Object.keys(userData).find(key => Array.isArray(userData[key]));
              if (arrayKey) {
                usersArray = userData[arrayKey];
              } else if (userData.id && userData.email) {
                // Handle single user object
                usersArray = [userData as User];
              }
            }

            if (usersArray.length > 0) {
            const foundUser = usersArray.find((u: User) => u.id === userId);
            if (foundUser) setUser(foundUser);
          }
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
          const errorText = await response.text();
          console.error('AdminUserCertificates: Fetch certificates failed', response.status, errorText);
          throw new Error(`Failed to fetch user certificates: ${response.status} ${response.statusText}`);
        }

        const data = await response.json();
        console.log('AdminUserCertificates: Received certificates data:', data);
        
        let certsArray: Certificate[] = [];
        if (Array.isArray(data)) {
          certsArray = data;
        } else if (data && typeof data === 'object') {
          // Try to find an array in the object properties
          const arrayKey = Object.keys(data).find(key => Array.isArray(data[key]));
          if (arrayKey) {
            certsArray = data[arrayKey];
            console.log(`AdminUserCertificates: Extracted certificates from key: "${arrayKey}"`);
          } else if (data.id && (data['cert-type-name'] || data.certificate_type)) {
            // Handle single certificate object if returned
            console.log('AdminUserCertificates: Received a single certificate object instead of an array. Wrapping in array.');
            certsArray = [data as Certificate];
          } else {
            // If no array found, check for a message or error field
            const message = data.message || data.error || data.details;
            const availableKeys = Object.keys(data).join(', ');
            console.error('AdminUserCertificates: No array found in object:', data);
            
            if (message) {
              throw new Error(`Server message: ${message}`);
            } else {
              throw new Error(`Received invalid data format for certificates. Expected an array but got object with keys: [${availableKeys}].`);
            }
          }
        } else {
          console.error('AdminUserCertificates: Expected array of certificates but got:', typeof data, data);
          throw new Error(`Received invalid data format for certificates. Expected an array but got ${typeof data}.`);
        }
        
        setCertificates(certsArray);
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
