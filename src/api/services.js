import api from './axios';

export const getAlarms = () => api.get('/alarm');

export const createAlarm = (data) => api.post('/alarm', data);

export const updateAlarm = (id, data) => api.put(`/alarm/${id}`, data);

export const deleteAlarm = (id) => api.delete(`/alarm/${id}`);

export const getDeviceColor = () => api.get('/device/color');

export const updateDeviceColor = (color) => api.put('/device/color', { color });

export const login = (credentials) => api.post('/login', credentials);