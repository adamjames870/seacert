import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Box,
  CircularProgress,
  Typography,
  IconButton,
  Alert
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import DownloadIcon from '@mui/icons-material/Download';
import { API_BASE_URL } from '../config';
import { supabase } from '../supabaseClient';

interface ReportPreviewDialogProps {
  open: boolean;
  onClose: () => void;
}

const ReportPreviewDialog: React.FC<ReportPreviewDialogProps> = ({ open, onClose }) => {
  const [loading, setLoading] = useState(true);
  const [pdfUrl, setPdfUrl] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (open) {
      fetchReport();
    } else {
      // Clean up the URL when dialog closes
      if (pdfUrl) {
        URL.revokeObjectURL(pdfUrl);
        setPdfUrl(null);
      }
      setError(null);
      setLoading(true);
    }
  }, [open]);

  const fetchReport = async () => {
    setLoading(true);
    setError(null);
    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) {
        throw new Error('You must be logged in to view the report.');
      }

      const response = await fetch(`${API_BASE_URL}/api/certificates/report`, {
        headers: {
          'Authorization': `Bearer ${session.access_token}`,
        },
      });

      if (!response.ok) {
        throw new Error('Failed to generate report.');
      }

      const blob = await response.blob();
      const url = URL.createObjectURL(blob);
      setPdfUrl(url);
    } catch (err: any) {
      console.error('Error fetching report:', err);
      setError(err.message || 'An error occurred while generating the report.');
    } finally {
      setLoading(false);
    }
  };

  const handleDownload = () => {
    if (pdfUrl) {
      const link = document.createElement('a');
      link.href = pdfUrl;
      link.download = 'seacert-report.pdf';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    }
  };

  return (
    <Dialog
      open={open}
      onClose={onClose}
      maxWidth="lg"
      fullWidth
      PaperProps={{
        sx: { height: '90vh' }
      }}
    >
      <DialogTitle sx={{ m: 0, p: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        Certificate Report Preview
        <IconButton
          aria-label="close"
          onClick={onClose}
          sx={{
            color: (theme) => theme.palette.grey[500],
          }}
        >
          <CloseIcon />
        </IconButton>
      </DialogTitle>
      <DialogContent dividers sx={{ p: 0, display: 'flex', flexDirection: 'column', bgcolor: 'grey.100' }}>
        {loading ? (
          <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', flexGrow: 1, gap: 2 }}>
            <CircularProgress />
            <Typography>Generating your report...</Typography>
          </Box>
        ) : error ? (
          <Box sx={{ p: 3 }}>
            <Alert severity="error">{error}</Alert>
          </Box>
        ) : (
          <Box sx={{ flexGrow: 1, width: '100%', height: '100%' }}>
            <iframe
              src={`${pdfUrl}#view=FitH`}
              title="Certificate Report"
              width="100%"
              height="100%"
              style={{ border: 'none' }}
            />
          </Box>
        )}
      </DialogContent>
      <DialogActions sx={{ p: 2 }}>
        <Button onClick={onClose} color="inherit">
          Close
        </Button>
        <Button
          onClick={handleDownload}
          variant="contained"
          startIcon={<DownloadIcon />}
          disabled={loading || !!error}
        >
          Download PDF
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default ReportPreviewDialog;
