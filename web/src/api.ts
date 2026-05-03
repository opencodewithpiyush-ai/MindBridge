/**
 * api.ts — MindBridge API client
 *
 * In local development (Vite dev server), requests go through the Vite proxy
 * at /api, which forwards them to the Render backend — bypassing browser CORS.
 *
 * In production, VITE_API_URL should be set in .env to the real backend URL.
 * If unset, it falls back to the Render URL directly (backend has CORS headers).
 */
const API_BASE: string =
  import.meta.env.VITE_API_URL ??
  (import.meta.env.DEV
    ? '/api'
    : 'http://127.0.0.1:5000');

export default API_BASE;
