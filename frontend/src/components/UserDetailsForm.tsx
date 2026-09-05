import { useState } from 'react';
import { Box, TextField, Autocomplete, Button, CircularProgress } from '@mui/material';
import { countries } from '../utils/countryData';

interface UserDetailsFormProps {
  initialData?: any;
  onSave: (data: any) => Promise<void>;
  onCancel?: () => void;
}

export const UserDetailsForm = ({ initialData, onSave, onCancel }: UserDetailsFormProps) => {
  const [formData, setFormData] = useState(initialData || { forename: '', surname: '', nationality: '' });
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    await onSave(formData);
    setLoading(false);
  };

  return (
    <Box component="form" onSubmit={handleSubmit} sx={{ mt: 2 }}>
      <TextField fullWidth label="Forename" name="forename" value={formData.forename} onChange={(e) => setFormData({...formData, forename: e.target.value})} sx={{ mb: 2 }} required />
      <TextField fullWidth label="Surname" name="surname" value={formData.surname} onChange={(e) => setFormData({...formData, surname: e.target.value})} sx={{ mb: 2 }} required />
      <Autocomplete
        options={countries}
        getOptionLabel={(option) => typeof option === 'string' ? option : option.label}
        value={formData.nationality ? countries.find(c => c.code === formData.nationality) : null}
        onChange={(_, newValue) => setFormData({...formData, nationality: newValue ? newValue.code : ''})}
        renderInput={(params) => <TextField {...params} label="Nationality" sx={{ mb: 2 }} />}
      />
      <Box sx={{ mt: 3, display: 'flex', gap: 2 }}>
        <Button type="submit" variant="contained" disabled={loading}>{loading ? <CircularProgress size={24} /> : 'Save'}</Button>
        {onCancel && <Button onClick={onCancel}>Cancel</Button>}
      </Box>
    </Box>
  );
};
