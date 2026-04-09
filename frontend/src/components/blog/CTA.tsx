import { Box, Typography, Button, Paper } from '@mui/material';

interface CTAProps {
  title: string;
  description: string;
  buttonText: string;
  href: string;
}

const CTA = ({ title, description, buttonText, href }: CTAProps) => {
  const isExternal = href.startsWith('http');
  
  return (
    <Paper 
      elevation={0} 
      sx={{ 
        p: 4, 
        my: 4, 
        textAlign: 'center', 
        bgcolor: 'primary.main', 
        color: 'white',
        borderRadius: 2
      }}
    >
      <Typography variant="h5" sx={{ fontWeight: 700, mb: 2 }}>
        {title}
      </Typography>
      <Typography variant="body1" sx={{ mb: 3, opacity: 0.9 }}>
        {description}
      </Typography>
      <Button 
        variant="contained" 
        color="secondary" 
        size="large"
        href={href}
        target={isExternal ? "_blank" : "_self"}
        rel={isExternal ? "noopener noreferrer" : ""}
        sx={{ 
          px: 4, 
          py: 1.5, 
          fontWeight: 700,
          bgcolor: 'white',
          color: 'primary.main',
          '&:hover': {
            bgcolor: '#f5f5f5'
          }
        }}
      >
        {buttonText}
      </Button>
    </Paper>
  );
};

export default CTA;
