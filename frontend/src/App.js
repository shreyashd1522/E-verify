import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import { Button, Box } from '@mui/material';
import VerifyPage from './pages/VerifyPage';
import ForgotPasswordPage from './pages/ForgotPasswordPage';
import ResetPasswordPage from './pages/ResetPasswordPage';
import ResendVerificationPage from './pages/ResendVerificationPage';
import AnimatedNavbar from './AnimatedNavbar';
import './ModernNavbar.css';

function App() {
  return (
    <Router>
      <AnimatedNavbar />
      <Routes>
        <Route path="/" element={<VerifyPage />} />
        <Route path="/forgot-password" element={<ForgotPasswordPage />} />
        <Route path="/reset-password" element={<ResetPasswordPage />} />
        <Route path="/resend-verification" element={<ResendVerificationPage />} />
      </Routes>
    </Router>
  );
}

export default App;
