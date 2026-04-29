export const API_BASE_URL =
  import.meta.env.VITE_API_URL || "https://test-mindbridge-v1-1.onrender.com";
export const FILE_BASE_URL =
  import.meta.env.VITE_FILE_BASE_URL || "https://files.use.ai";

type SessionExpiredHandler = () => void;
let onSessionExpired: SessionExpiredHandler | null = null;

export function setSessionExpiredHandler(handler: SessionExpiredHandler) {
  onSessionExpired = handler;
}

export function triggerSessionExpired() {
  onSessionExpired?.();
}

export async function apiFetch<T>(
  path: string,
  options: RequestInit = {},
  skipSessionCheck = false
): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
  });

  if (response.status === 401 && !skipSessionCheck) {
    onSessionExpired?.();
    throw new Error("Session expired");
  }

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error || data.message || "Something went wrong");
  }

  return data;
}

export async function uploadFile(
  file: File,
  token: string
): Promise<{ key: string; url: string; preview: string; name: string; type: string }> {
  const formData = new FormData();
  formData.append("file", file);
  formData.append("name", file.name);
  formData.append("type", file.type);

  const response = await fetch(`${API_BASE_URL}/upload`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token}`,
    },
    body: formData,
  });

  if (response.status === 401) {
    onSessionExpired?.();
    throw new Error("Session expired");
  }

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error || "Upload failed");
  }

  const key: string = data.key;
  const rawUrl: string = data.url;
  const fullUrl = `${FILE_BASE_URL}${rawUrl}`;
  const previewUrl = `https://use.ai/_next/image?url=${encodeURIComponent(fullUrl)}&w=1024&q=75`;

  return { key, url: fullUrl, preview: previewUrl, name: file.name, type: file.type };
}

export async function getFileURL(
  key: string,
  token: string
): Promise<string> {
  const response = await fetch(`${API_BASE_URL}/file/${key}`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  if (response.status === 401) {
    onSessionExpired?.();
    throw new Error("Session expired");
  }

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error || "Failed to get file URL");
  }

  return data.url;
}

export async function verifySession(token: string): Promise<boolean> {
  try {
    const response = await fetch(`${API_BASE_URL}/models`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.status !== 401;
  } catch {
    return false;
  }
}
