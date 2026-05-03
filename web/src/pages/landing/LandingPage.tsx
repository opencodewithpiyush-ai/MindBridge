import { useNavigate } from 'react-router-dom';
import './LandingPage.css';

const models = [
  { name: 'GPT-5.5', provider: 'OpenAI', badge: 'Latest', color: '#10a37f' },
  { name: 'GPT-5.4', provider: 'OpenAI', badge: 'Fast', color: '#10a37f' },
  { name: 'GPT-5', provider: 'OpenAI', badge: 'Powerful', color: '#10a37f' },
  { name: 'GPT-4o', provider: 'OpenAI', badge: 'Vision', color: '#10a37f' },
  { name: 'Claude Sonnet 4.6', provider: 'Anthropic', badge: 'Smart', color: '#d97706' },
  { name: 'Claude Opus 4.5', provider: 'Anthropic', badge: 'Deep', color: '#d97706' },
  { name: 'Gemini 3.1 Pro', provider: 'Google', badge: 'Multi', color: '#4285f4' },
  { name: 'Gemini 2.5 Flash', provider: 'Google', badge: 'Quick', color: '#4285f4' },
  { name: 'Grok-4', provider: 'xAI', badge: 'Witty', color: '#1d9bf0' },
  { name: 'DeepSeek V4 Pro', provider: 'DeepSeek', badge: 'Reason', color: '#6366f1' },
  { name: 'DeepSeek R1', provider: 'DeepSeek', badge: 'Think', color: '#6366f1' },
  { name: 'Llama 3.3 70B', provider: 'Meta', badge: 'Open', color: '#0866ff' },
  { name: 'Kimi K2', provider: 'Moonshot', badge: 'New', color: '#7c3aed' },
  { name: 'Qwen 3 Max', provider: 'Alibaba', badge: 'Max', color: '#f59e0b' },
];

const features = [
  {
    icon: '🧠',
    title: '20 AI Models',
    desc: 'GPT-5.5, Claude, Gemini, Grok, DeepSeek — all under one roof, completely free.',
  },
  {
    icon: '🖼️',
    title: 'Image Editing',
    desc: 'Upload your photos and let AI enhance, retouch, or transform them instantly.',
  },
  {
    icon: '📄',
    title: 'Document Summarizer',
    desc: 'Drop any long PDF or doc — get a crisp summary in seconds, no reading needed.',
  },
  {
    icon: '⚡',
    title: 'Real-time Streaming',
    desc: 'Responses stream live as the AI thinks — no waiting for the full reply.',
  },
  {
    icon: '🔒',
    title: 'Private & Secure',
    desc: 'Your conversations are yours. We never store or share your data with anyone.',
  },
  {
    icon: '🌐',
    title: 'Always Free',
    desc: 'No subscriptions, no credit card, no limits. Just pure AI power, always.',
  },
];

export default function LandingPage() {
  const navigate = useNavigate();

  return (
    <div className="landing">
      {/* ── Navbar ── */}
      <nav className="nav">
        <div className="nav-inner">
          <div className="nav-brand">
            <span className="nav-logo">🧠</span>
            <span className="nav-name">MindBridge</span>
          </div>
          <div className="nav-actions">
            <button className="btn-ghost" onClick={() => navigate('/login')}>
              Log in
            </button>
            <button className="btn-primary" onClick={() => navigate('/register')}>
              Get Started
            </button>
          </div>
        </div>
      </nav>

      {/* ── Hero ── */}
      <section className="hero-section">
        {/* Decorative blobs */}
        <div className="blob blob-1" aria-hidden="true" />
        <div className="blob blob-2" aria-hidden="true" />
        <div className="blob blob-3" aria-hidden="true" />

        <div className="hero-content">
          <div className="hero-badge">
            <span className="badge-dot" />
            20 World-Class AI Models — Free Forever
          </div>

          <h1 className="hero-title">
            One platform.<br />
            <span className="hero-gradient">Every AI model.</span>
          </h1>

          <p className="hero-subtitle">
            Chat, create, summarize, and edit images — powered by GPT-5.5,
            Claude, Gemini, Grok and 16 more top AI models. No subscriptions. No limits.
          </p>

          <div className="hero-actions">
            <button
              id="get-started-btn"
              className="btn-hero-primary"
              onClick={() => navigate('/register')}
            >
              Get Started — It's Free
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
                <path d="M5 12h14M12 5l7 7-7 7" />
              </svg>
            </button>
            <button
              id="continue-email-btn"
              className="btn-hero-secondary"
              onClick={() => navigate('/login')}
            >
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <rect x="2" y="4" width="20" height="16" rx="2" />
                <path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7" />
              </svg>
              Continue with Email
            </button>
          </div>

          <p className="hero-note">No credit card required · Instant access · Always free</p>
        </div>

        {/* Floating model chips */}
        <div className="hero-chips" aria-hidden="true">
          {['GPT-5.5', 'Claude Sonnet', 'Gemini 3.1', 'Grok-4', 'DeepSeek R1'].map((m, i) => (
            <div key={m} className="chip" style={{ '--delay': `${i * 0.3}s` } as React.CSSProperties}>
              {m}
            </div>
          ))}
        </div>
      </section>

      {/* ── Models Scroll ── */}
      <section className="models-section">
        <div className="section-label">Powered By</div>
        <div className="models-track-wrapper">
          <div className="models-track">
            {[...models, ...models].map((m, i) => (
              <div key={i} className="model-card">
                <span className="model-dot" style={{ background: m.color }} />
                <div className="model-info">
                  <span className="model-name">{m.name}</span>
                  <span className="model-provider">{m.provider}</span>
                </div>
                <span className="model-badge">{m.badge}</span>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Features ── */}
      <section className="features-section">
        <div className="section-inner">
          <div className="section-header">
            <div className="section-label">What You Can Do</div>
            <h2 className="section-title">Everything you need,<br /> nothing you don't.</h2>
            <p className="section-desc">
              MindBridge brings the world's best AI models to you — for free, forever.
            </p>
          </div>

          <div className="features-grid">
            {features.map((f) => (
              <div key={f.title} className="feature-card">
                <div className="feature-icon">{f.icon}</div>
                <h3 className="feature-title">{f.title}</h3>
                <p className="feature-desc">{f.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Use Cases ── */}
      <section className="usecases-section">
        <div className="section-inner">
          <div className="section-header">
            <div className="section-label">Real World Uses</div>
            <h2 className="section-title">Built for everyone.</h2>
          </div>
          <div className="usecases-grid">
            <div className="usecase-card usecase-large">
              <div className="usecase-emoji">✍️</div>
              <h3>Write & Edit</h3>
              <p>Draft emails, essays, code, and stories — with AI that actually understands context.</p>
            </div>
            <div className="usecase-card">
              <div className="usecase-emoji">🔍</div>
              <h3>Research Fast</h3>
              <p>Ask complex questions and get structured, sourced answers in seconds.</p>
            </div>
            <div className="usecase-card">
              <div className="usecase-emoji">📊</div>
              <h3>Analyze Data</h3>
              <p>Upload spreadsheets or docs and get instant insights, trends, and summaries.</p>
            </div>
            <div className="usecase-card usecase-large">
              <div className="usecase-emoji">🎨</div>
              <h3>Edit Images</h3>
              <p>Upload your photo, describe what you want changed — AI handles the rest with stunning precision.</p>
            </div>
          </div>
        </div>
      </section>

      {/* ── CTA Banner ── */}
      <section className="cta-section">
        <div className="cta-inner">
          <div className="blob cta-blob-1" aria-hidden="true" />
          <div className="blob cta-blob-2" aria-hidden="true" />
          <h2 className="cta-title">Ready to experience<br />the future of AI?</h2>
          <p className="cta-sub">
            Join thousands of users already using MindBridge — for free.
          </p>
          <div className="cta-actions">
            <button className="btn-hero-primary" onClick={() => navigate('/register')}>
              Create Free Account
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
                <path d="M5 12h14M12 5l7 7-7 7" />
              </svg>
            </button>
          </div>
        </div>
      </section>

      {/* ── Footer ── */}
      <footer className="footer">
        <div className="footer-inner">
          <div className="footer-brand">
            <span className="nav-logo">🧠</span>
            <span className="nav-name">MindBridge</span>
          </div>
          <p className="footer-copy">© 2026 MindBridge · Built by Piyush Makwana · MIT License</p>
          <div className="footer-links">
            <a href="https://github.com/piyushmakwana" target="_blank" rel="noreferrer">GitHub</a>
            <a href="https://test-mindbridge-v1-1.onrender.com" target="_blank" rel="noreferrer">API Docs</a>
          </div>
        </div>
      </footer>
    </div>
  );
}
