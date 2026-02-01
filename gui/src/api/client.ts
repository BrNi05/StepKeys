import axios from 'axios';

export const api = axios.create({
  baseURL: '/api',
  timeout: 5000,
});

// Settings
export const getBoot = () => api.get('/boot');
export const toggleBoot = () => api.post('/boot');

export const getEnabled = () => api.get('/enabled');
export const toggleEnabled = () => api.post('/enabled');

// Lifecycle
export const getUpdate = (force = false) => api.get(`/update?force=${force}`);
export const quitApp = () => api.post('/quit');

// Additional
export const getSerial = () => api.get('/serial');
export const getLogs = () => api.get('/logs');
export const getValidKeys = () => api.get('/valid-keys');

// Pedals
export const getPedals = () => api.get('/pedals');
export async function setPedals(pedals: any) {
  try {
    const res = await api.post('/pedals', pedals);
    return { ok: true, data: res.data };
  } catch (err: any) {
    return {
      ok: false,
      // Order of precedence: server error message, axios error message, fallback message
      message: err?.response?.data?.errorMessage || err?.message || 'Unknown API call error',
    };
  }
}
