export const API_URL = 'http://localhost:8080';

export async function apiFetch(endpoint, options = {}) {
  const url = `${API_URL}${endpoint}`;
  
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  if (options.body instanceof FormData) {
    delete headers['Content-Type'];
  }

  const config = {
    ...options,
    headers,
    credentials: 'include',
  };

  const response = await fetch(url, config);

  if (response.status === 401) {
    if (typeof window !== 'undefined') {
      window.dispatchEvent(new Event('auth:unauthorized'));
    }
  }

  return response;
}

export function toMediaUrl(path) {
  if (!path) return null;
  if (typeof path !== 'string') return null;
  const p = path.trim();
  if (p.startsWith('http://') || p.startsWith('https://')) return p;
  if (p.startsWith('/uploads')) return `${API_URL}${p}`;
  if (p.startsWith('uploads')) return `${API_URL}/${p}`;
  if (p.startsWith('/')) return `${API_URL}${p}`;
  return `${API_URL}/uploads/${p}`;
}
