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
  TextField,
  InputAdornment,
  Link,
  Button
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import SearchIcon from '@mui/icons-material/Search';
import AddIcon from '@mui/icons-material/Add';
import { Link as RouterLink } from 'react-router-dom';
import { supabase } from '../supabaseClient';
import { API_BASE_URL } from '../config';
import { getCountryName } from '../utils/countryData';

interface Issuer {
  id: string;
  name: string;
  country: string;
  website: string;
}

const Issuers = () => {
  const [issuers, setIssuers] = useState<Issuer[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');

  useEffect(() => {
    const fetchIssuers = async () => {
      try {
        setLoading(true);
        const { data: { session } } = await supabase.auth.getSession();
        
        if (!session) {
          setError('Not authenticated');
          setLoading(false);
          return;
        }

        const response = await fetch(`${API_BASE_URL}/api/issuers`, {
          headers: {
            'Authorization': `Bearer ${session.access_token}`,
          },
        });

        if (!response.ok) {
          throw new Error(`Error fetching issuers: ${response.statusText}`);
        }

        const data = await response.json();
        setIssuers(data);
        setError(null);
      } catch (err: any) {
        setError(err.message || 'Failed to load issuers');
      } finally {
        setLoading(false);
      }
    };

    fetchIssuers();
  }, []);

  const filteredIssuers = issuers.filter((issuer) => {
    const query = searchQuery.toLowerCase();
    return (
      issuer.name.toLowerCase().includes(query) ||
      issuer.country?.toLowerCase().includes(query) ||
      getCountryName(issuer.country).toLowerCase().includes(query)
    );
  });

  const sortedIssuers = [...filteredIssuers].sort((a, b) => 
    a.name.localeCompare(b.name)
  );

  return (
    <Container>
      <Box sx={{ mt: 4 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h4" component="h1">
            Issuers
          </Typography>
          
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
            <TextField
              size="small"
              placeholder="Search issuers..."
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
              to="/add-issuer"
              state={{ from: 'issuers' }}
            >
              Add Issuer
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

        {!loading && !error && issuers.length === 0 && (
          <Typography variant="body1" sx={{ mt: 2 }}>
            No issuers found.
          </Typography>
        )}

        {!loading && !error && issuers.length > 0 && filteredIssuers.length === 0 && (
          <Typography variant="body1" sx={{ mt: 2 }}>
            No issuers match your search.
          </Typography>
        )}

        {!loading && !error && filteredIssuers.length > 0 && (
          <List sx={{ mt: 2 }}>
            {sortedIssuers.map((issuer) => (
              <Paper key={issuer.id} elevation={0} sx={{ mb: 1, border: 1, borderColor: 'divider', overflow: 'hidden' }}>
                <Box 
                  sx={{ 
                    display: 'flex', 
                    justifyContent: 'space-between', 
                    alignItems: 'center',
                    p: 2
                  }}
                >
                  <ListItemText 
                    primary={issuer.name}
                    secondary={
                      <>
                        <Typography component="span" variant="body2" color="text.secondary">
                          Country: {getCountryName(issuer.country)}
                        </Typography>
                        {issuer.website && (
                          <Box component="span" sx={{ display: 'block' }}>
                            <Link href={issuer.website.startsWith('http') ? issuer.website : `https://${issuer.website}`} target="_blank" rel="noopener" variant="body2">
                              {issuer.website}
                            </Link>
                          </Box>
                        )}
                      </>
                    }
                  />
                  <Tooltip title="Edit Issuer">
                    <IconButton 
                      component={RouterLink} 
                      to={`/edit-issuer/${issuer.id}`}
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

export default Issuers;
