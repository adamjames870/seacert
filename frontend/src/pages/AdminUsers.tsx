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
  TextField,
  InputAdornment,
  Stack,
  Tooltip,
} from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import VisibilityIcon from '@mui/icons-material/Visibility';
import { useNavigate } from 'react-router-dom';
import { supabase } from '../supabaseClient';
import { API_BASE_URL } from '../config';

interface User {
  id: string;
  forename: string;
  surname: string;
  email: string;
  nationality: string;
  role: string;
  certificate_count?: number;
}

const AdminUsers = () => {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const navigate = useNavigate();

  const fetchUsers = async () => {
    setLoading(true);
    setError(null);
    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) {
        console.error('AdminUsers: No session found');
        setError('Not authenticated');
        return;
      }

      console.log('AdminUsers: Fetching users from', `${API_BASE_URL}/admin/users`);
      const response = await fetch(`${API_BASE_URL}/admin/users`, {
        headers: {
          'Authorization': `Bearer ${session.access_token}`,
        },
      });

      console.log('AdminUsers: Fetch response status:', response.status);
      if (!response.ok) {
        const errorText = await response.text();
        console.error('AdminUsers: Fetch failed', response.status, errorText);
        throw new Error(`Failed to fetch users: ${response.status} ${response.statusText}`);
      }

      const data = await response.json();
      console.log('AdminUsers: Received users data:', data);
      
      let usersArray: User[] = [];
      if (Array.isArray(data)) {
        usersArray = data;
      } else if (data && typeof data === 'object') {
        // Try to find an array in the object properties
        const arrayKey = Object.keys(data).find(key => Array.isArray(data[key]));
        if (arrayKey) {
          usersArray = data[arrayKey];
          console.log(`AdminUsers: Extracted users from key: "${arrayKey}"`);
        } else if (data.id && data.email && data.role) {
          // If it's a single user object (as reported in the issue)
          console.log('AdminUsers: Received a single user object instead of an array. Wrapping in array.');
          usersArray = [data as User];
        } else {
          // If no array found, check for a message or error field
          const message = data.message || data.error || data.details;
          const availableKeys = Object.keys(data).join(', ');
          console.error('AdminUsers: No array found in object:', data);
          
          if (message) {
            throw new Error(`Server message: ${message}`);
          } else {
            throw new Error(`Received invalid data format from server. Expected an array but got object with keys: [${availableKeys}].`);
          }
        }
      } else {
        console.error('AdminUsers: Expected array of users but got:', typeof data, data);
        throw new Error(`Received invalid data format from server. Expected an array but got ${typeof data}.`);
      }
      
      setUsers(usersArray);
    } catch (err: any) {
      console.error('AdminUsers: Error in fetchUsers:', err);
      setError(err.message || 'An error occurred while fetching users');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  const filteredUsers = users.filter((user) => {
    const query = searchQuery.toLowerCase();
    const forename = user.forename || '';
    const surname = user.surname || '';
    const email = user.email || '';
    return (
      forename.toLowerCase().includes(query) ||
      surname.toLowerCase().includes(query) ||
      email.toLowerCase().includes(query)
    );
  });

  const handleViewCertificates = (userId: string) => {
    navigate(`/admin/users/${userId}/certificates`);
  };

  return (
    <Container maxWidth="lg">
      <Box sx={{ mt: 4, mb: 4 }}>
        <Stack
          direction={{ xs: 'column', md: 'row' }}
          spacing={2}
          justifyContent="space-between"
          alignItems={{ xs: 'stretch', md: 'center' }}
          sx={{ mb: 3 }}
        >
          <Typography variant="h4" component="h1" sx={{ fontWeight: 700 }}>
            User Management
          </Typography>

          <TextField
            size="small"
            placeholder="Search users..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <SearchIcon fontSize="small" />
                </InputAdornment>
              ),
            }}
            sx={{ minWidth: { xs: '100%', sm: 300 } }}
          />
        </Stack>

        {loading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', mt: 8 }}>
            <CircularProgress />
          </Box>
        ) : error ? (
          <Alert severity="error">{error}</Alert>
        ) : (
          <TableContainer component={Paper} elevation={0} sx={{ border: 1, borderColor: 'divider' }}>
            <Table>
              <TableHead sx={{ bgcolor: 'action.hover' }}>
                <TableRow>
                  <TableCell sx={{ fontWeight: 700 }}>Name</TableCell>
                  <TableCell sx={{ fontWeight: 700 }}>Email</TableCell>
                  <TableCell sx={{ fontWeight: 700 }}>Nationality</TableCell>
                  <TableCell sx={{ fontWeight: 700 }}>Role</TableCell>
                  <TableCell sx={{ fontWeight: 700 }} align="center">Certificates</TableCell>
                  <TableCell sx={{ fontWeight: 700 }} align="right">Actions</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {filteredUsers.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} align="center" sx={{ py: 4 }}>
                      <Typography variant="body1" color="text.secondary">
                        No users found
                      </Typography>
                    </TableCell>
                  </TableRow>
                ) : (
                  filteredUsers.map((user) => (
                    <TableRow key={user.id} hover>
                      <TableCell>{`${user.forename} ${user.surname}`}</TableCell>
                      <TableCell>{user.email}</TableCell>
                      <TableCell>{user.nationality}</TableCell>
                      <TableCell sx={{ textTransform: 'capitalize' }}>{user.role}</TableCell>
                      <TableCell align="center">{user.certificate_count ?? 'N/A'}</TableCell>
                      <TableCell align="right">
                        <Tooltip title="View Certificates">
                          <IconButton
                            color="primary"
                            onClick={() => handleViewCertificates(user.id)}
                            size="small"
                          >
                            <VisibilityIcon fontSize="small" />
                          </IconButton>
                        </Tooltip>
                      </TableCell>
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

export default AdminUsers;
