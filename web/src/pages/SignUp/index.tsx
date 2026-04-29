import { useState, type FormEvent } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../../context/AuthContext";
import styles from "./Auth.module.css";

interface ValidationError {
  field: string;
  message: string;
}

function SignUp() {
  const { signup, isAuthenticated } = useAuth();
  const navigate = useNavigate();

  if (isAuthenticated) {
    navigate("/chat");
    return null;
  }
  const [name, setName] = useState("");
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [errors, setErrors] = useState<ValidationError[]>([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setErrors([]);
    setError("");
    setLoading(true);

    try {
      await signup(name, username, email, password);
      navigate("/chat");
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : "Signup failed";
      try {
        const parsed = JSON.parse(message);
        if (Array.isArray(parsed.errors)) {
          setErrors(parsed.errors);
        } else if (parsed.error) {
          setError(parsed.error);
        } else {
          setError(message);
        }
      } catch {
        setError(message);
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.auth}>
      <div className={styles["auth-glow"]} />
      <div className={styles["grid-bg"]} />

      <nav className={styles.navbar}>
        <Link to="/" className={styles.logo}>
          <div className={styles["logo-icon"]}>M</div>
          <span className={styles["logo-text"]}>MindBridge</span>
        </Link>
      </nav>

      <div className={styles.container}>
        <div className={styles.card}>
          <h1 className={styles.title}>Create account</h1>
          <p className={styles.subtitle}>Get started with MindBridge for free</p>

          {error && <div className={styles["error-message"]}>{error}</div>}

          {errors.length > 0 && (
            <div className={styles["error-message"]}>
              {errors.map((e) => (
                <p key={e.field}>
                  <strong>{e.field}:</strong> {e.message}
                </p>
              ))}
            </div>
          )}

          <form onSubmit={handleSubmit} className={styles.form}>
            <div className={styles["form-group"]}>
              <label htmlFor="name" className={styles.label}>
                Full Name
              </label>
              <input
                id="name"
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                className={styles.input}
                placeholder="John Doe"
                required
              />
            </div>

            <div className={styles["form-group"]}>
              <label htmlFor="username" className={styles.label}>
                Username
              </label>
              <input
                id="username"
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className={styles.input}
                placeholder="johndoe"
                required
              />
            </div>

            <div className={styles["form-group"]}>
              <label htmlFor="email" className={styles.label}>
                Email
              </label>
              <input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className={styles.input}
                placeholder="you@example.com"
                required
              />
            </div>

            <div className={styles["form-group"]}>
              <label htmlFor="password" className={styles.label}>
                Password
              </label>
              <input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className={styles.input}
                placeholder="Min 8 characters"
                required
              />
            </div>

            <button
              type="submit"
              className={styles["btn-submit"]}
              disabled={loading}
            >
              {loading ? "Creating account..." : "Sign up"}
            </button>
          </form>

          <p className={styles.footer}>
            Already have an account?{" "}
            <Link to="/login" className={styles.link}>
              Sign in
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}

export default SignUp;
