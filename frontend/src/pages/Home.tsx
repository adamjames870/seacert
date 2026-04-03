import { Typography, Container, Box, Button, Grid, Stack, Card, CardContent, CardMedia } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';
import { Anchor, ShieldCheck, FileText, DownloadCloud, AlertCircle } from 'lucide-react';

const Home = () => {
  return (
    <Box>
      {/* Hero Section */}
      <Box 
        component="section"
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
      <Box component="section" sx={{ bgcolor: 'white', py: 10 }}>
        <Container>
          <Typography variant="h4" component="h2" align="center" gutterBottom sx={{ fontWeight: 700, mb: 8 }}>
            Why Use SeaCert?
          </Typography>
          <Grid container spacing={4} justifyContent="center">
            {[
              {
                icon: <ShieldCheck size={48} />,
                title: "Complete Records",
                description: "Maintain comprehensive records of all your certifications and STCW training in one secure location.",
                image: "https://images.unsplash.com/photo-1550751827-4bd374c3f58b?auto=format&fit=crop&q=80&w=2070"
              },
              {
                icon: <AlertCircle size={48} />,
                title: "Expiry Reminder",
                description: "Easily see when your certificates need renewing. Plan your renewal before they expire so you're always ready for your next contract.",
                image: "https://images.unsplash.com/photo-1508962914676-134849a727f0?auto=format&fit=crop&q=80&w=2070"
              },
              {
                icon: <DownloadCloud size={48} />,
                title: "Easy Access Anywhere",
                description: "Store secure copies of your certificates for instant access. Download them on ship or at home, whenever you need them.",
                image: "https://images.unsplash.com/photo-1451187580459-43490279c0fa?auto=format&fit=crop&q=80&w=2072"
              },
              {
                icon: <FileText size={48} />,
                title: "Professional Reports",
                description: "Generate a comprehensive report of all your certifications in seconds. Perfect for sending to crewing agents and employers.",
                image: "https://images.unsplash.com/photo-1450101499163-c8848c66ca85?auto=format&fit=crop&q=80&w=2070"
              }
            ].map((feature, index) => (
              <Grid size={{ xs: 12, sm: 6, md: 3 }} key={index} sx={{ display: 'flex' }}>
                <Card sx={{ 
                  width: '100%', 
                  display: 'flex', 
                  flexDirection: 'column',
                  transition: 'transform 0.2s',
                  '&:hover': { transform: 'translateY(-8px)' },
                  boxShadow: 2
                }}>
                  <CardMedia
                    component="img"
                    height="140"
                    image={feature.image}
                    alt={feature.title}
                  />
                  <CardContent sx={{ flexGrow: 1, textAlign: 'center', pt: 3 }}>
                    <Box 
                      sx={{ 
                        color: 'primary.main', 
                        bgcolor: 'rgba(74, 109, 140, 0.1)', 
                        p: 2, 
                        borderRadius: '50%',
                        display: 'inline-flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        mb: 2
                      }}
                    >
                      {feature.icon}
                    </Box>
                    <Typography variant="h6" component="h3" sx={{ fontWeight: 600, mb: 1 }}>
                      {feature.title}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {feature.description}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
            ))}
          </Grid>
        </Container>
      </Box>

      {/* CTA Section */}
      <Box component="section" sx={{ bgcolor: 'background.paper', py: 10, borderTop: 1, borderColor: 'divider' }}>
        <Container maxWidth="sm" sx={{ textAlign: 'center' }}>
          <Typography variant="h4" component="h2" gutterBottom sx={{ fontWeight: 700 }}>
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

    </Box>
  );
};

export default Home;
