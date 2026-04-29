import { createContext, useContext, useState, useEffect, useCallback } from "react";
import { apiFetch, setSessionExpiredHandler, verifySession } from "../lib/api";

interface User {
  id: string;
  name: string;
  username: string;
  email: string;
}

interface AuthContextType {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  signup: (
    name: string,
    username: string,
    email: string,
    password: string
  ) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  const clearSession = useCallback(() => {
    setToken(null);
    setUser(null);
    localStorage.removeItem("token");
    localStorage.removeItem("user");
  }, []);

  useEffect(() => {
    setSessionExpiredHandler(clearSession);
  }, [clearSession]);

  useEffect(() => {
    const storedToken = localStorage.getItem("token");
    const storedUser = localStorage.getItem("user");

    if (storedToken && storedUser) {
      verifySession(storedToken).then((valid) => {
        if (valid) {
          setToken(storedToken);
          setUser(JSON.parse(storedUser));
        } else {
          clearSession();
        }
        setLoading(false);
      });
    } else {
      setLoading(false);
    }
  }, [clearSession]);

  const login = async (email: string, password: string) => {
    const response = await apiFetch<{
      data: { token: string; user: User };
    }>("/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });

    setToken(response.data.token);
    setUser(response.data.user);
    localStorage.setItem("token", response.data.token);
    localStorage.setItem("user", JSON.stringify(response.data.user));
  };

  const signup = async (
    name: string,
    username: string,
    email: string,
    password: string
  ) => {
    const response = await apiFetch<{
      data: { token: string; user: User };
    }>("/auth/register", {
      method: "POST",
      body: JSON.stringify({ name, username, email, password }),
    });

    setToken(response.data.token);
    setUser(response.data.user);
    localStorage.setItem("token", response.data.token);
    localStorage.setItem("user", JSON.stringify(response.data.user));
  };

  const logout = async () => {
    const storedToken = localStorage.getItem("token");
    if (storedToken) {
      try {
        await apiFetch(
          "/auth/logout",
          {
            method: "POST",
            headers: {
              Authorization: `Bearer ${storedToken}`,
            },
          },
          true
        );
      } catch {
        // Ignore errors on logout
      }
    }
    clearSession();
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        isAuthenticated: !!token,
        loading,
        login,
        signup,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return context;
}
