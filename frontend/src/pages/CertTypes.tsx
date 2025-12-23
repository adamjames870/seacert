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
  'stcw-reference': string;
  'normal-validity-months': number;
}

const CertTypes = () => {
  const [certTypes, setCertTypes] = useState<CertType[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');

  const getMissingFieldsStatus = (type: CertType) => {
    const missing = [];
    if (!type['short-name']) missing.push('short-name');
    if (!type['stcw-reference']) missing.push('stcw-reference');
    if (!type['normal-validity-months']) missing.push('normal-validity-months');
    
    if (missing.length > 0) return 'incomplete';
    return 'normal';
  };

  const getStatusStyles = (status: string) => {
    switch (status) {
      case 'incomplete':
        return {
          bgcolor: '#fffbeb', // Amber 50
          borderColor: '#fef3c7', // Amber 100
          textColor: '#92400e', // Amber 800
          secondaryTextColor: '#b45309', // Amber 700
          labelColor: '#d97706', // Amber 600
        };
      default:
        return {
          bgcolor: 'background.paper',
          borderColor: 'divider',
          textColor: 'text.primary',
          secondaryTextColor: 'text.secondary',
          labelColor: 'primary.main',
        };
    }
  };

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
      type['stcw-reference']?.toLowerCase().includes(query)
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
            {sortedCertTypes.map((type) => {
              const status = getMissingFieldsStatus(type);
              const styles = getStatusStyles(status);

              return (
                <Paper 
                  key={type.id} 
                  elevation={0} 
                  sx={{ 
                    mb: 1, 
                    border: 1, 
                    borderColor: styles.borderColor, 
                    bgcolor: styles.bgcolor,
                    overflow: 'hidden' 
                  }}
                >
                  <Box 
                    sx={{ 
                      display: 'flex', 
                      justifyContent: 'space-between', 
                      alignItems: 'center',
                      p: 2
                    }}
                  >
                    <ListItemText 
                      primary={
                        <Typography variant="subtitle1" sx={{ fontWeight: 600, color: styles.textColor }}>
                          {type.name}
                        </Typography>
                      }
                      secondary={
                        <>
                          <Typography component="span" variant="body2" sx={{ color: styles.secondaryTextColor }}>
                            Short Name: {type['short-name'] || (
                              <Box component="span" sx={{ fontStyle: 'italic', fontWeight: 'bold', color: styles.labelColor }}>
                                Missing
                              </Box>
                            )}
                            {' | '}
                            STCW: {type['stcw-reference'] || (
                              <Box component="span" sx={{ fontStyle: 'italic', fontWeight: 'bold', color: styles.labelColor }}>
                                Missing
                              </Box>
                            )}
                          </Typography>
                          {type['normal-validity-months'] ? (
                            <Typography variant="body2" sx={{ color: styles.secondaryTextColor }}>
                              Validity: {type['normal-validity-months']} months
                            </Typography>
                          ) : (
                            <Typography variant="body2" sx={{ fontStyle: 'italic', fontWeight: 'bold', color: styles.labelColor }}>
                              Validity: Missing
                            </Typography>
                          )}
                        </>
                      }
                    />
                    <Tooltip title="Edit Certificate Type">
                      <IconButton 
                        component={RouterLink} 
                        to={`/edit-cert-type/${type.id}`}
                        sx={{ color: styles.labelColor }}
                      >
                        <EditIcon />
                      </IconButton>
                    </Tooltip>
                  </Box>
                </Paper>
              );
            })}
          </List>
        )}
      </Box>
    </Container>
  );
};

export default CertTypes;
