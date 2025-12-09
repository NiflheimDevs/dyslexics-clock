import api from "./axios";

export const getAlarms = () => api.get("/alarm");

export const createAlarm = (data) => api.post("/alarm", data);

export const updateAlarm = (id, data) => api.patch(`/alarm/${id}`, data);

export const deleteAlarm = (id) => api.delete(`/alarm/${id}`);

export const getDeviceColor = () => api.get("/device/color");

export const getDeviceVolume = () => api.get("/device");

export const updateDeviceColor = (color) =>
  api.patch("/device/color", { color });

export const updateDeviceVolume = (volume) =>
  api.patch("/device/volume", { volume });

export const login = (credentials) => api.post("/login", credentials);
