import { List, ListItem, ListItemIcon, ListItemText, Box } from '@mui/material';
import { CheckCircle2 } from 'lucide-react';

interface ChecklistProps {
  items: string[];
}

const Checklist = ({ items }: ChecklistProps) => {
  return (
    <Box sx={{ my: 3 }}>
      <List sx={{ p: 0 }}>
        {items.map((item, index) => (
          <ListItem key={index} alignItems="flex-start" sx={{ px: 0, py: 0.5 }}>
            <ListItemIcon sx={{ minWidth: 36, mt: 0.5 }}>
              <CheckCircle2 size={20} color="#4A6D8C" />
            </ListItemIcon>
            <ListItemText 
              primary={item} 
              primaryTypographyProps={{ 
                variant: 'body1',
                sx: { fontWeight: 500 }
              }} 
            />
          </ListItem>
        ))}
      </List>
    </Box>
  );
};

export default Checklist;
