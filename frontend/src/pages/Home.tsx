import { Typography, Container, Box, Button, Grid, Stack } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';
import { Anchor, ShieldCheck, Bell, Smartphone, Clock } from 'lucide-react';

const Home = () => {
  return (
    <Box>
      {/* Hero Section */}
      <Box 
        sx={{ 
          bgcolor: 'primary.main', 
          color: 'white', 
          pt: 12, 
          pb: 10,
          backgroundImage: 'linear-gradient(rgba(74, 109, 140, 0.8), rgba(74, 109, 140, 0.9)), url("https://images.unsplash.com/photo-1517048676732-d65bc937f952?auto=format&fit=crop&q=80&w=2070")',
          backgroundSize: 'cover',
          backgroundPosition: 'center',
        }}
      >
        <Container maxWidth="md">
          <Stack spacing={4} alignItems="center" textAlign="center">
            <Anchor size={80} strokeWidth={1.5} />
            <Box>
              <Typography variant="h2" component="h1" gutterBottom sx={{ fontWeight: 800 }}>
                Manage Your Sea Career
              </Typography>
              <Typography variant="h5" sx={{ opacity: 0.9, mb: 4 }}>
                Keep all your maritime certificates in one place. Never miss an expiry date again.
              </Typography>
            </Box>
            <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
              <Button 
                variant="contained" 
                color="secondary" 
                size="large" 
                component={RouterLink} 
                to="/signup"
                sx={{ px: 4, py: 1.5, fontSize: '1.1rem', fontWeight: 600, bgcolor: 'white', color: 'primary.main', '&:hover': { bgcolor: '#f0f0f0' } }}
              >
                Start for Free
              </Button>
              <Button 
                variant="outlined" 
                color="inherit" 
                size="large" 
                component={RouterLink} 
                to="/login"
                sx={{ px: 4, py: 1.5, fontSize: '1.1rem', fontWeight: 600, borderWidth: 2, '&:hover': { borderWidth: 2 } }}
              >
                Sign In
              </Button>
            </Stack>
          </Stack>
        </Container>
      </Box>

      {/* Features Section */}
      <Box sx={{ bgcolor: 'white', py: 10 }}>
        <Container>
          <Typography variant="h4" component="h2" align="center" gutterBottom sx={{ fontWeight: 700, mb: 8 }}>
            Why Use SeaCert?
          </Typography>
          <Grid container spacing={6} justifyContent="center">
            {[
              {
                icon: <ShieldCheck size={48} />,
                title: "Compliance Ready",
                description: "Easily track all your STCW and mandatory certificates required for your next contract."
              },
              {
                icon: <Bell size={48} />,
                title: "Expiry Alerts",
                description: "Get notified before your certificates expire, giving you plenty of time for renewal courses."
              },
              {
                icon: <Smartphone size={48} />,
                title: "Always Accessible",
                description: "Access your certificate details anytime, anywhere. On the ship or at home."
              },
              {
                icon: <Clock size={48} />,
                title: "Time Saving",
                description: "No more digging through folders. Find issuer details and certificate numbers in seconds."
              }
            ].map((feature, index) => (
              <Grid size={{ xs: 12, sm: 6, md: 3 }} key={index} sx={{ display: 'flex' }}>
                <Stack 
                  spacing={2} 
                  alignItems="center" 
                  textAlign="center"
                  sx={{ width: '100%' }}
                >
                  <Box 
                    sx={{ 
                      color: 'primary.main', 
                      bgcolor: 'rgba(74, 109, 140, 0.1)', 
                      p: 2, 
                      borderRadius: '50%',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      mb: 1
                    }}
                  >
                    {feature.icon}
                  </Box>
                  <Typography variant="h6" component="h3" sx={{ fontWeight: 600 }}>
                    {feature.title}
                  </Typography>
                  <Typography variant="body1" color="text.secondary" sx={{ maxWidth: 280 }}>
                    {feature.description}
                  </Typography>
                </Stack>
              </Grid>
            ))}
          </Grid>
        </Container>
      </Box>

      {/* CTA Section */}
      <Box sx={{ bgcolor: 'background.paper', py: 10, borderTop: 1, borderColor: 'divider' }}>
        <Container maxWidth="sm" sx={{ textAlign: 'center' }}>
          <Typography variant="h4" gutterBottom sx={{ fontWeight: 700 }}>
            Ready to get organized?
          </Typography>
          <Typography variant="body1" color="text.secondary" sx={{ mb: 4 }}>
            Join other seafarers who manage their career with SeaCert.
          </Typography>
          <Button 
            variant="contained" 
            size="large" 
            component={RouterLink} 
            to="/signup"
            sx={{ px: 6, py: 2, borderRadius: 2, fontWeight: 600 }}
          >
            Create Your Account
          </Button>
        </Container>
      </Box>

      {/* Footer */}
      <Box component="footer" sx={{ py: 6, borderTop: 1, borderColor: 'divider' }}>
        <Container>
          <Typography variant="body2" color="text.secondary" align="center">
            © {new Date().getFullYear()} SeaCert. All rights reserved.
          </Typography>
        </Container>
      </Box>
    </Box>
  );
};

export default Home;
