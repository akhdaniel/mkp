export interface TravelPackage {
  id: number;
  title: string;
  destination: string;
  country: string;
  type: 'city' | 'beach' | 'adventure' | 'cultural' | 'culinary' | 'countryside';
  duration: number;
  price: number;
  rating: number;
  image: string;
  description: string;
  highlights: string[];
  included: string[];
  tags: string[];
  score?: number;
}

export interface Destination {
  name: string;
  count: number;
  popular: boolean;
}

export interface UserPurchase {
  packageId: number;
  purchaseDate: string;
  status: 'completed' | 'upcoming' | 'cancelled';
  package?: TravelPackage;
}

export interface SearchFilters {
  query?: string;
  priceRange?: [number, number];
  duration?: number;
  type?: string;
}