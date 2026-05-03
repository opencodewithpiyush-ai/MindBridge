import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import './AuthPage.css';
import API_BASE from '../../api';

export default function RegisterPage() {
  const navigate = useNavigate();

  const [form, setForm] = useState({ name: '', username: '', email: '', password: '' });
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

    if (!form.name || !form.username || !form.email || !form.password) {
      setError('All fields are required.');
      setLoading(false);
      return;
    }
    if (form.password.length < 6) {
      setError('Password must be at least 6 characters.');
      setLoading(false);
      return;
    }

    try {
      const res = await fetch(`${API_BASE}/auth/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(form),
      });
      const data = await res.json();

      if (!res.ok || !data.success) {
        throw new Error(data.error || data.message || 'Registration failed. Please try again.');
      }

      setSuccess('Account created! Redirecting to login...');
      setTimeout(() => navigate('/login'), 1800);
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
            "One platform. Every AI model. Zero cost."
          </blockquote>
          <div className="al-models">
            {['GPT-5.5', 'Claude Sonnet 4.6', 'Gemini 3.1 Pro', 'Grok-4', 'DeepSeek R1', 'Kimi K2'].map(m => (
              <span key={m} className="al-chip">{m}</span>
            ))}
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
            <h1 className="auth-title">Create account</h1>
            <p className="auth-sub">Join MindBridge and access 20 AI models for free.</p>
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
            <div className="form-row">
              <div className="form-group">
                <label htmlFor="reg-name" className="form-label">Full Name</label>
                <input
                  id="reg-name"
                  className="form-input"
                  type="text"
                  name="name"
                  placeholder="Piyush Makwana"
                  value={form.name}
                  onChange={handleChange}
                  autoComplete="name"
                  disabled={loading}
                />
              </div>
              <div className="form-group">
                <label htmlFor="reg-username" className="form-label">Username</label>
                <input
                  id="reg-username"
                  className="form-input"
                  type="text"
                  name="username"
                  placeholder="piyush123"
                  value={form.username}
                  onChange={handleChange}
                  autoComplete="username"
                  disabled={loading}
                />
              </div>
            </div>

            <div className="form-group">
              <label htmlFor="reg-email" className="form-label">Email Address</label>
              <input
                id="reg-email"
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
              <label htmlFor="reg-password" className="form-label">Password</label>
              <input
                id="reg-password"
                className="form-input"
                type="password"
                name="password"
                placeholder="Minimum 6 characters"
                value={form.password}
                onChange={handleChange}
                autoComplete="new-password"
                disabled={loading}
              />
            </div>

            <button
              id="register-submit-btn"
              type="submit"
              className="auth-submit"
              disabled={loading}
            >
              {loading ? (
                <><span className="spinner" />Creating account...</>
              ) : (
                <>Create Account — It's Free</>
              )}
            </button>
          </form>

          <p className="auth-switch">
            Already have an account?{' '}
            <Link to="/login" className="auth-link">Log in</Link>
          </p>

          <p className="auth-terms">
            By creating an account, you agree to our{' '}
            <a href="#terms" className="auth-link">Terms of Service</a>{' '}
            and{' '}
            <a href="#privacy" className="auth-link">Privacy Policy</a>.
          </p>
        </div>
      </div>
    </div>
  );
}
