import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import './AuthPage.css';

const API_BASE = 'https://test-mindbridge-v1-1.onrender.com';

export default function LoginPage() {
  const navigate = useNavigate();

  const [form, setForm] = useState({ email: '', password: '' });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm(prev => ({ ...prev, [e.target.name]: e.target.value }));
    setError('');
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setSuccess('');

    if (!form.email || !form.password) {
      setError('Email and password are required.');
      setLoading(false);
      return;
    }

    try {
      const res = await fetch(`${API_BASE}/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(form),
      });
      const data = await res.json();

      if (!res.ok || !data.success) {
        throw new Error(data.error || data.message || 'Login failed. Check your credentials.');
      }

      // Store token
      if (data.token) {
        localStorage.setItem('mb_token', data.token);
      }
      if (data.data?.token) {
        localStorage.setItem('mb_token', data.data.token);
      }

      setSuccess('Welcome back! Logging you in...');
      setTimeout(() => navigate('/'), 1500);
    } catch (err: unknown) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('Something went wrong. Please try again.');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-page">
      {/* Left panel */}
      <div className="auth-left" aria-hidden="true">
        <div className="auth-left-inner">
          <div className="al-brand">
            <span className="al-logo">🧠</span>
            <span className="al-name">MindBridge</span>
          </div>
          <blockquote className="al-quote">
            "Your AI, your way. Every model. Zero cost."
          </blockquote>
          <div className="al-stats">
            <div className="al-stat">
              <span className="al-stat-num">20</span>
              <span className="al-stat-label">AI Models</span>
            </div>
            <div className="al-stat-divider" />
            <div className="al-stat">
              <span className="al-stat-num">∞</span>
              <span className="al-stat-label">Free Chats</span>
            </div>
            <div className="al-stat-divider" />
            <div className="al-stat">
              <span className="al-stat-num">0₹</span>
              <span className="al-stat-label">Cost</span>
            </div>
          </div>
        </div>
        <div className="al-blob-1" />
        <div className="al-blob-2" />
      </div>

      {/* Right panel */}
      <div className="auth-right">
        <div className="auth-card">
          {/* Back */}
          <Link to="/" className="auth-back">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
              <path d="M19 12H5M12 19l-7-7 7-7" />
            </svg>
            Back to Home
          </Link>

          <div className="auth-header">
            <h1 className="auth-title">Welcome back</h1>
            <p className="auth-sub">Log in to your MindBridge account.</p>
          </div>

          {error && (
            <div className="auth-alert auth-alert-error" role="alert">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <circle cx="12" cy="12" r="10" /><line x1="12" y1="8" x2="12" y2="12" /><line x1="12" y1="16" x2="12.01" y2="16" />
              </svg>
              {error}
            </div>
          )}

          {success && (
            <div className="auth-alert auth-alert-success" role="status">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <circle cx="12" cy="12" r="10" /><path d="m9 12 2 2 4-4" />
              </svg>
              {success}
            </div>
          )}

          <form className="auth-form" onSubmit={handleSubmit} noValidate>
            <div className="form-group">
              <label htmlFor="login-email" className="form-label">Email Address</label>
              <input
                id="login-email"
                className="form-input"
                type="email"
                name="email"
                placeholder="you@example.com"
                value={form.email}
                onChange={handleChange}
                autoComplete="email"
                disabled={loading}
              />
            </div>

            <div className="form-group">
              <div className="form-label-row">
                <label htmlFor="login-password" className="form-label">Password</label>
                <a href="#forgot" className="auth-link form-label-right">Forgot password?</a>
              </div>
              <input
                id="login-password"
                className="form-input"
                type="password"
                name="password"
                placeholder="Your password"
                value={form.password}
                onChange={handleChange}
                autoComplete="current-password"
                disabled={loading}
              />
            </div>

            <button
              id="login-submit-btn"
              type="submit"
              className="auth-submit"
              disabled={loading}
            >
              {loading ? (
                <><span className="spinner" />Logging in...</>
              ) : (
                <>Continue with Email</>
              )}
            </button>
          </form>

          <p className="auth-switch">
            Don't have an account?{' '}
            <Link to="/register" className="auth-link">Create one free</Link>
          </p>
        </div>
      </div>
    </div>
  );
}
