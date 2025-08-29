import apiClient from './client';

export const bookingAPI = {
  searchSchedules: async (params) => {
    const response = await apiClient.get('/schedules/search', { params });
    return response.data;
  },

  getSchedule: async (id) => {
    const response = await apiClient.get(`/schedules/${id}`);
    return response.data;
  },

  createBooking: async (data) => {
    const response = await apiClient.post('/bookings', data);
    return response.data;
  },

  getBooking: async (id) => {
    const response = await apiClient.get(`/bookings/${id}`);
    return response.data;
  },

  getBookingByReference: async (reference) => {
    const response = await apiClient.get(`/bookings/reference/${reference}`);
    return response.data;
  },

  cancelBooking: async (id, reason) => {
    const response = await apiClient.post(`/bookings/${id}/cancel`, { reason });
    return response.data;
  },

  getMyBookings: async () => {
    const response = await apiClient.get('/bookings/my-bookings');
    return response.data;
  },

  checkInTicket: async (qrCode) => {
    const response = await apiClient.post('/tickets/check-in', { qr_code: qrCode });
    return response.data;
  },

  getManifest: async (scheduleId) => {
    const response = await apiClient.get(`/schedules/${scheduleId}/manifest`);
    return response.data;
  },
};