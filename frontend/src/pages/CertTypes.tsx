import { useEffect, useState } from 'react';
import { 
  Typography, 
  Container, 
  Box, 
  Paper, 
  List, 
  ListItemText, 
  Alert, 
  CircularProgress, 
  IconButton, 
  Tooltip,
  ListItemButton,
  Button,
  TextField,
  InputAdornment
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import SearchIcon from '@mui/icons-material/Search';
import AddIcon from '@mui/icons-material/Add';
import { Link as RouterLink } from 'react-router-dom';
import { supabase } from '../supabaseClient';
import { API_BASE_URL } from '../config';

interface CertType {
  id: string;
  name: string;
  'short-name': string;
  'stcw-ref': string;
  'normal-validity-months': number;
}

const CertTypes = () => {
  const [certTypes, setCertTypes] = useState<CertType[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');

  useEffect(() => {
    const fetchCertTypes = async () => {
      try {
        setLoading(true);
        const { data: { session } } = await supabase.auth.getSession();
        
        if (!session) {
          setError('Not authenticated');
          setLoading(false);
          return;
        }

        const response = await fetch(`${API_BASE_URL}/api/cert-types`, {
          headers: {
            'Authorization': `Bearer ${session.access_token}`,
          },
        });

        if (!response.ok) {
          throw new Error(`Error fetching certificate types: ${response.statusText}`);
        }

        const data = await response.json();
        setCertTypes(data);
        setError(null);
      } catch (err: any) {
        setError(err.message || 'Failed to load certificate types');
      } finally {
        setLoading(false);
      }
    };

    fetchCertTypes();
  }, []);

  const filteredCertTypes = certTypes.filter((type) => {
    const query = searchQuery.toLowerCase();
    return (
      type.name.toLowerCase().includes(query) ||
      type['short-name']?.toLowerCase().includes(query) ||
      type['stcw-ref']?.toLowerCase().includes(query)
    );
  });

  const sortedCertTypes = [...filteredCertTypes].sort((a, b) => 
    a.name.localeCompare(b.name)
  );

  return (
    <Container>
      <Box sx={{ mt: 4 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h4" component="h1">
            Certificate Types
          </Typography>
          
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
            <TextField
              size="small"
              placeholder="Search types..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <SearchIcon fontSize="small" />
                  </InputAdornment>
                ),
              }}
              sx={{ width: 250 }}
            />
            <Button
              variant="contained"
              color="primary"
              startIcon={<AddIcon />}
              component={RouterLink}
              to="/add-cert-type"
            >
              Add Certificate Type
            </Button>
          </Box>
        </Box>

        {loading && (
          <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
            <CircularProgress />
          </Box>
        )}

        {error && (
          <Alert severity="error" sx={{ mt: 2 }}>
            {error}
          </Alert>
        )}

        {!loading && !error && certTypes.length === 0 && (
          <Typography variant="body1" sx={{ mt: 2 }}>
            No certificate types found.
          </Typography>
        )}

        {!loading && !error && certTypes.length > 0 && filteredCertTypes.length === 0 && (
          <Typography variant="body1" sx={{ mt: 2 }}>
            No certificate types match your search.
          </Typography>
        )}

        {!loading && !error && filteredCertTypes.length > 0 && (
          <List sx={{ mt: 2 }}>
            {sortedCertTypes.map((type) => (
              <Paper key={type.id} elevation={0} sx={{ mb: 1, border: 1, borderColor: 'divider', overflow: 'hidden' }}>
                <Box 
                  sx={{ 
                    display: 'flex', 
                    justifyContent: 'space-between', 
                    alignItems: 'center',
                    p: 2
                  }}
                >
                  <ListItemText 
                    primary={type.name}
                    secondary={
                      <>
                        <Typography component="span" variant="body2" color="text.secondary">
                          Short Name: {type['short-name'] || 'N/A'} | STCW: {type['stcw-ref'] || 'N/A'}
                        </Typography>
                        {type['normal-validity-months'] && (
                          <Typography variant="body2" color="text.secondary">
                            Validity: {type['normal-validity-months']} months
                          </Typography>
                        )}
                      </>
                    }
                  />
                  <Tooltip title="Edit Certificate Type">
                    <IconButton 
                      component={RouterLink} 
                      to={`/edit-cert-type/${type.id}`}
                      color="primary"
                    >
                      <EditIcon />
                    </IconButton>
                  </Tooltip>
                </Box>
              </Paper>
            ))}
          </List>
        )}
      </Box>
    </Container>
  );
};

export default CertTypes;
