import React, { useState } from 'react';
import { Button } from '@mui/material';
import { Link, useLocation } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import './ModernNavbar.css';

const navLinks = [
  { label: 'VERIFY', to: '/' },
  { label: 'FORGOT PASSWORD', to: '/forgot-password' },
  { label: 'RESEND VERIFICATION', to: '/resend-verification' },
];

export default function AnimatedNavbar() {
  const [showNav, setShowNav] = useState(false);
  const location = useLocation();

  return (
    <div className="navbar-center-area">
      <div
        className="site-title-centered"
        onMouseEnter={() => setShowNav(true)}
        onMouseLeave={() => setShowNav(false)}
        tabIndex={0}
        onFocus={() => setShowNav(true)}
        onBlur={() => setShowNav(false)}
      >
        E-verify
        <AnimatePresence>
          {showNav && (
            <motion.nav
              className="modern-navbar-float-centered"
              initial={{ opacity: 0, y: -10 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -10 }}
              transition={{ duration: 0.25, ease: 'easeOut' }}
            >
              {navLinks.map((link, idx) => {
                const isActive = location.pathname === link.to || (link.to === '/' && location.pathname === '');
                return (
                  <div key={link.to} className="modern-nav-item-wrapper">
                    <div className={`nav-dot${isActive ? ' nav-dot-active' : ''}`}></div>
                    <Button
                      className={`modern-nav-btn${isActive ? ' nav-btn-active' : ''}`}
                      component={Link}
                      to={link.to}
                      sx={{ textTransform: 'uppercase', letterSpacing: '0.12em', fontWeight: 700 }}
                    >
                      {link.label}
                      {isActive && (
                        <motion.div
                          className="nav-btn-highlight"
                          layoutId="nav-btn-highlight"
                          transition={{ type: 'spring', stiffness: 500, damping: 30 }}
                        />
                      )}
                    </Button>
                  </div>
                );
              })}
            </motion.nav>
          )}
        </AnimatePresence>
      </div>
    </div>
  );
} 