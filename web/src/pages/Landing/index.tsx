import { Link } from "react-router-dom";
import styles from "./Landing.module.css";

function Landing() {
  return (
    <div className={styles.landing}>
      <div className={styles["hero-glow"]} />
      <div className={styles["grid-bg"]} />

      <nav className={styles.navbar}>
        <div className={styles.logo}>
          <div className={styles["logo-icon"]}>M</div>
          <span className={styles["logo-text"]}>MindBridge</span>
        </div>
        <div className={styles["nav-actions"]}>
          <Link to="/login" className={`${styles.btn} ${styles["btn-ghost"]}`}>
            Log in
          </Link>
          <Link to="/signup" className={`${styles.btn} ${styles["btn-primary"]}`}>
            Sign up
          </Link>
        </div>
      </nav>

      <main className={styles.hero}>
        <div className={styles.badge}>
          <span className={styles["badge-dot"]} />
          Powered by 17+ AI models
        </div>

        <h1 className={styles["hero-title"]}>
          Chat with any AI,
          <br />
          <span className={styles["gradient-text"]}>one gateway.</span>
        </h1>

        <p className={styles["hero-description"]}>
          Access GPT-5, Claude, Gemini, Grok and more through a single clean interface.
          Sign up free and start chatting in seconds.
        </p>

        <div className={styles["hero-actions"]}>
          <Link to="/signup" className={`${styles.btn} ${styles["btn-primary"]} ${styles["btn-large"]}`}>
            Get Started Free
          </Link>
          <Link to="/login" className={`${styles.btn} ${styles["btn-outline"]} ${styles["btn-large"]}`}>
            Log in
          </Link>
        </div>
      </main>

      <section className={styles.features}>
        <div className={styles["feature-card"]}>
          <div className={styles["feature-icon"]}>⚡</div>
          <h3 className={styles["feature-title"]}>Lightning Fast</h3>
          <p className={styles["feature-description"]}>
            Real-time streaming responses with zero lag across all models.
          </p>
        </div>

        <div className={styles["feature-card"]}>
          <div className={styles["feature-icon"]}>🛡️</div>
          <h3 className={styles["feature-title"]}>Secure & Private</h3>
          <p className={styles["feature-description"]}>
            JWT authentication and encrypted sessions keep your data safe.
          </p>
        </div>

        <div className={styles["feature-card"]}>
          <div className={styles["feature-icon"]}>🎯</div>
          <h3 className={styles["feature-title"]}>17+ Models</h3>
          <p className={styles["feature-description"]}>
            Switch between GPT, Claude, Gemini, Grok and more instantly.
          </p>
        </div>
      </section>
    </div>
  );
}

export default Landing;
