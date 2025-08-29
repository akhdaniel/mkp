import apiClient from './client';

export const portsAPI = {
  getAllPorts: async (params = {}) => {
    const response = await apiClient.get('/ports', { params });
    return response.data;
  },

  getPort: async (id) => {
    const response = await apiClient.get(`/ports/${id}`);
    return response.data;
  },

  searchPorts: async (city, country) => {
    const response = await apiClient.get('/ports/search', {
      params: { city, country },
    });
    return response.data;
  },

  getRoutes: async (departurePortId, arrivalPortId) => {
    const response = await apiClient.get('/routes/search', {
      params: {
        departure_port_id: departurePortId,
        arrival_port_id: arrivalPortId,
      },
    });
    return response.data;
  },
};