import { Typography, Container, Box } from '@mui/material';

const Home = () => {
  return (
    <Container>
      <Box sx={{ mt: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Home Page
        </Typography>
        <Typography variant="body1">
          Welcome to the Home Page (Placeholder).
        </Typography>
      </Box>
    </Container>
  );
};

export default Home;
