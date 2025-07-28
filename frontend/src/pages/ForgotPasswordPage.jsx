import React, { useState, useEffect } from 'react';
import { Box, Button, TextField, Typography, Paper, InputAdornment } from '@mui/material';
import { motion, AnimatePresence } from 'framer-motion';
import Lottie from 'react-lottie-player';
import EmailIcon from '@mui/icons-material/Email';

// Helper to fetch Lottie JSON from public folder
function useLottie(url) {
  const [data, setData] = useState(null);
  useEffect(() => {
    fetch(url)
      .then(res => res.json())
      .then(setData);
  }, [url]);
  return data;
}

const ForgotPasswordPage = () => {
  const [email, setEmail] = useState('');
  const [status, setStatus] = useState('idle'); // idle | loading | success | error
  const [message, setMessage] = useState('');

  const successAnim = useLottie(process.env.PUBLIC_URL + '/animations/greencheck.json');
  const errorAnim = useLottie(process.env.PUBLIC_URL + '/animations/redcross.json');
  const loadingAnim = useLottie(process.env.PUBLIC_URL + '/animations/loading.json');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setStatus('loading');
    setMessage('');
    try {
      const res = await fetch('http://localhost:8080/forgot-password', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email })
      });
      const data = await res.json();
      if (data.status === 'success' || data.message?.toLowerCase().includes('reset link')) {
        setStatus('success');
        setMessage(data.message);
      } else {
        setStatus('error');
        setMessage(data.message || 'An error occurred.');
      }
    } catch (err) {
      setStatus('error');
      setMessage('Server error. Please try again.');
    }
  };

  return (
    <Box minHeight="100vh" display="flex" alignItems="center" justifyContent="center" sx={{ background: 'radial-gradient(circle at 30% 50%, #a259ff 0%, #181c24 70%)', fontFamily: 'Inter, Arial, sans-serif' }}>
      <Paper elevation={6} sx={{ borderRadius: 4, p: 5, maxWidth: 410, width: '100%' }}>
        <motion.div initial={{ opacity: 0, y: 30 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.6, ease: 'easeOut' }}>
          <Typography variant="h4" align="center" fontWeight={700} mb={4}>
            Forgot Password
          </Typography>
          <form onSubmit={handleSubmit} autoComplete="off">
            <TextField
              label="Email"
              type="email"
              fullWidth
              required
              margin="normal"
              variant="outlined"
              size="medium"
              value={email}
              onChange={e => setEmail(e.target.value)}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <EmailIcon color="primary" />
                  </InputAdornment>
                ),
              }}
            />
            <motion.div initial={{ scale: 0.95, opacity: 0 }} animate={{ scale: 1, opacity: 1 }} transition={{ duration: 0.4, delay: 0.2 }}>
              <Button
                type="submit"
                variant="contained"
                fullWidth
                sx={{
                  mt: 2,
                  py: 1.5,
                  fontWeight: 700,
                  fontSize: '1.1rem',
                  background: 'linear-gradient(90deg, #4f8cff 0%, #6a5cff 100%)',
                  boxShadow: '0 2px 8px 0 #4f8cff22',
                  letterSpacing: 0.5,
                  borderRadius: 2,
                  '&:hover': {
                    background: 'linear-gradient(90deg, #6a5cff 0%, #4f8cff 100%)',
                    boxShadow: '0 4px 16px 0 #4f8cff44',
                  },
                }}
              >
                SEND RESET LINK
              </Button>
            </motion.div>
          </form>
        </motion.div>
        <Box mt={3} textAlign="center" minHeight={32}>
          <AnimatePresence>
            {status === 'loading' && (
              <motion.div key="loading" initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }}>
                <Lottie loop play animationData={loadingAnim} style={{ width: 80, height: 80, margin: '0 auto' }} />
                <Typography color="text.secondary">Processing...</Typography>
              </motion.div>
            )}
            {status === 'success' && (
              <motion.div key="success" initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }}>
                <Lottie loop play animationData={successAnim} style={{ width: 80, height: 80, margin: '0 auto' }} />
                <Typography color="success.main">{message}</Typography>
              </motion.div>
            )}
            {status === 'error' && (
              <motion.div key="error" initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }}>
                <Lottie loop play animationData={errorAnim} style={{ width: 80, height: 80, margin: '0 auto' }} />
                <Typography color="error.main">{message}</Typography>
              </motion.div>
            )}
          </AnimatePresence>
        </Box>
      </Paper>
    </Box>
  );
};

export default ForgotPasswordPage; 