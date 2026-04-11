import { Paper, Box, Typography } from '@mui/material';
import { AlertCircle } from 'lucide-react';
import type { ReactNode } from 'react';

interface WarningBoxProps {
  children: ReactNode;
}

const WarningBox = ({ children }: WarningBoxProps) => {
  return (
    <Paper 
      elevation={0} 
      sx={{ 
        p: 3, 
        my: 4, 
        bgcolor: 'rgba(239, 68, 68, 0.1)', 
        borderLeft: '4px solid #ef4444',
        borderRadius: 1,
        display: 'flex',
        gap: 2
      }}
    >
      <Box sx={{ color: '#ef4444', mt: 0.5 }}>
        <AlertCircle size={24} />
      </Box>
      <Box sx={{ flexGrow: 1 }}>
        <Typography variant="body1" sx={{ color: 'text.primary', fontWeight: 500 }}>
          {children}
        </Typography>
      </Box>
    </Paper>
  );
};

export default WarningBox;
