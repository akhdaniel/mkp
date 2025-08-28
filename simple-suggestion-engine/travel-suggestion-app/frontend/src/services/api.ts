import axios from 'axios';
import { TravelPackage, Destination, UserPurchase, SearchFilters } from '../types';

const API_BASE_URL = 'http://localhost:5000/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const travelAPI = {
  getAllPackages: async (): Promise<TravelPackage[]> => {
    const response = await api.get('/packages');
    return response.data;
  },

  getPackageById: async (id: number): Promise<TravelPackage> => {
    const response = await api.get(`/packages/${id}`);
    return response.data;
  },

  getDestinations: async (): Promise<Destination[]> => {
    const response = await api.get('/destinations');
    return response.data;
  },

  getPackagesByCountry: async (country: string): Promise<TravelPackage[]> => {
    const response = await api.get(`/destinations/${country}`);
    return response.data;
  },

  searchPackages: async (filters: SearchFilters): Promise<TravelPackage[]> => {
    const response = await api.post('/search', filters);
    return response.data;
  },

  getUserPurchases: async (userId: string): Promise<UserPurchase[]> => {
    const response = await api.get(`/user/${userId}/purchases`);
    return response.data;
  },

  getSuggestions: async (userId: string): Promise<TravelPackage[]> => {
    const response = await api.get(`/suggestions/${userId}`);
    return response.data;
  },

  purchasePackage: async (userId: string, packageId: number): Promise<any> => {
    const response = await api.post(`/user/${userId}/purchase`, { packageId });
    return response.data;
  },
};

export default api;