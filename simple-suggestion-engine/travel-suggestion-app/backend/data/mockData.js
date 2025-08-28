const travelPackages = [
  // Japan Packages
  {
    id: 1,
    title: "Tokyo Explorer",
    destination: "Tokyo",
    country: "Japan",
    type: "city",
    duration: 5,
    price: 1299,
    rating: 4.8,
    image: "https://images.unsplash.com/photo-1540959733332-eab4deabeeaf?w=800",
    description: "Experience the vibrant culture of Tokyo with visits to Shibuya, Harajuku, and traditional temples.",
    highlights: ["Tokyo Tower", "Senso-ji Temple", "Shibuya Crossing", "Tsukiji Market"],
    included: ["Hotel", "Breakfast", "City Tours", "Airport Transfer"],
    tags: ["culture", "city", "shopping", "food"]
  },
  {
    id: 2,
    title: "Kyoto Heritage Tour",
    destination: "Kyoto",
    country: "Japan",
    type: "cultural",
    duration: 4,
    price: 999,
    rating: 4.9,
    image: "https://images.unsplash.com/photo-1493976040374-85c8e12f0c0e?w=800",
    description: "Discover ancient temples, traditional gardens, and geisha districts in historic Kyoto.",
    highlights: ["Fushimi Inari", "Kinkaku-ji", "Arashiyama Bamboo Grove", "Gion District"],
    included: ["Ryokan Stay", "Kaiseki Dinner", "Temple Tours", "Tea Ceremony"],
    tags: ["culture", "history", "temples", "traditional"]
  },
  {
    id: 3,
    title: "Mount Fuji Adventure",
    destination: "Mount Fuji",
    country: "Japan",
    type: "adventure",
    duration: 3,
    price: 799,
    rating: 4.7,
    image: "https://images.unsplash.com/photo-1570459027562-4a916cc6113f?w=800",
    description: "Climb Japan's iconic Mount Fuji and enjoy breathtaking views of the surrounding landscape.",
    highlights: ["Mount Fuji Climb", "Lake Kawaguchi", "Fuji Five Lakes", "Hakone Hot Springs"],
    included: ["Mountain Guide", "Equipment", "Mountain Hut Stay", "Meals"],
    tags: ["adventure", "hiking", "nature", "mountains"]
  },

  // Thailand Packages
  {
    id: 4,
    title: "Bangkok Street Food Paradise",
    destination: "Bangkok",
    country: "Thailand",
    type: "culinary",
    duration: 4,
    price: 699,
    rating: 4.6,
    image: "https://images.unsplash.com/photo-1508009603885-50cf7c579365?w=800",
    description: "Explore Bangkok's legendary street food scene and vibrant markets.",
    highlights: ["Chatuchak Market", "Grand Palace", "Wat Pho", "Street Food Tours"],
    included: ["Hotel", "Food Tours", "Cooking Class", "River Cruise"],
    tags: ["food", "city", "culture", "shopping"]
  },
  {
    id: 5,
    title: "Phuket Beach Escape",
    destination: "Phuket",
    country: "Thailand",
    type: "beach",
    duration: 7,
    price: 1199,
    rating: 4.5,
    image: "https://images.unsplash.com/photo-1537956965359-7573183d1f57?w=800",
    description: "Relax on pristine beaches and enjoy crystal-clear waters in tropical Phuket.",
    highlights: ["Patong Beach", "Phi Phi Islands", "Big Buddha", "Old Town"],
    included: ["Beach Resort", "Island Hopping", "Snorkeling", "Spa Treatment"],
    tags: ["beach", "relaxation", "islands", "water sports"]
  },

  // France Packages
  {
    id: 6,
    title: "Paris Romance",
    destination: "Paris",
    country: "France",
    type: "city",
    duration: 5,
    price: 1599,
    rating: 4.8,
    image: "https://images.unsplash.com/photo-1502602898657-3e91760cbb34?w=800",
    description: "Fall in love with the City of Light, from the Eiffel Tower to charming cafés.",
    highlights: ["Eiffel Tower", "Louvre Museum", "Notre-Dame", "Versailles"],
    included: ["Boutique Hotel", "Museum Passes", "Seine Cruise", "Wine Tasting"],
    tags: ["romance", "culture", "art", "cuisine"]
  },
  {
    id: 7,
    title: "Provence Lavender Fields",
    destination: "Provence",
    country: "France",
    type: "countryside",
    duration: 6,
    price: 1399,
    rating: 4.7,
    image: "https://images.unsplash.com/photo-1562079225-bf9413703aa1?w=800",
    description: "Wander through lavender fields and explore charming villages in Provence.",
    highlights: ["Lavender Fields", "Avignon", "Gordes", "Wine Vineyards"],
    included: ["Country Inn", "Wine Tours", "Market Visits", "Cooking Class"],
    tags: ["nature", "wine", "relaxation", "countryside"]
  },

  // USA Packages
  {
    id: 8,
    title: "New York City Lights",
    destination: "New York",
    country: "USA",
    type: "city",
    duration: 5,
    price: 1799,
    rating: 4.6,
    image: "https://images.unsplash.com/photo-1538970272646-f61fabb3a8a2?w=800",
    description: "Experience the energy of NYC from Times Square to Central Park.",
    highlights: ["Statue of Liberty", "Central Park", "Broadway Show", "MoMA"],
    included: ["Manhattan Hotel", "City Pass", "Broadway Ticket", "Walking Tours"],
    tags: ["city", "culture", "entertainment", "shopping"]
  },
  {
    id: 9,
    title: "Grand Canyon Adventure",
    destination: "Arizona",
    country: "USA",
    type: "adventure",
    duration: 4,
    price: 999,
    rating: 4.9,
    image: "https://images.unsplash.com/photo-1474044159687-1ee9f3a51722?w=800",
    description: "Witness the majesty of the Grand Canyon with hiking and helicopter tours.",
    highlights: ["South Rim", "Helicopter Tour", "Hiking Trails", "Sunset Views"],
    included: ["Lodge Stay", "Park Passes", "Guided Tours", "Helicopter Ride"],
    tags: ["nature", "adventure", "hiking", "photography"]
  },

  // Italy Packages
  {
    id: 10,
    title: "Rome Ancient Wonders",
    destination: "Rome",
    country: "Italy",
    type: "cultural",
    duration: 5,
    price: 1399,
    rating: 4.8,
    image: "https://images.unsplash.com/photo-1552832230-c0197dd311b5?w=800",
    description: "Walk through history in the Eternal City from the Colosseum to Vatican City.",
    highlights: ["Colosseum", "Vatican Museums", "Trevi Fountain", "Roman Forum"],
    included: ["Central Hotel", "Skip-the-line Tickets", "Food Tour", "Vatican Tour"],
    tags: ["history", "culture", "art", "cuisine"]
  },
  {
    id: 11,
    title: "Tuscany Wine Country",
    destination: "Tuscany",
    country: "Italy",
    type: "culinary",
    duration: 6,
    price: 1599,
    rating: 4.9,
    image: "https://images.unsplash.com/photo-1523978591478-c753949ff840?w=800",
    description: "Savor world-class wines and cuisine in the rolling hills of Tuscany.",
    highlights: ["Chianti Vineyards", "Florence", "Siena", "San Gimignano"],
    included: ["Villa Stay", "Wine Tastings", "Cooking Class", "Truffle Hunting"],
    tags: ["wine", "food", "countryside", "relaxation"]
  },

  // Indonesia Packages
  {
    id: 12,
    title: "Bali Paradise",
    destination: "Bali",
    country: "Indonesia",
    type: "beach",
    duration: 7,
    price: 1099,
    rating: 4.7,
    image: "https://images.unsplash.com/photo-1537996194471-e657df975ab4?w=800",
    description: "Find your zen in Bali with beaches, temples, and rice terraces.",
    highlights: ["Uluwatu Temple", "Rice Terraces", "Ubud", "Beach Clubs"],
    included: ["Villa with Pool", "Spa Treatments", "Temple Tours", "Surfing Lesson"],
    tags: ["beach", "culture", "wellness", "spirituality"]
  },

  // Spain Packages
  {
    id: 13,
    title: "Barcelona Art & Architecture",
    destination: "Barcelona",
    country: "Spain",
    type: "city",
    duration: 5,
    price: 1299,
    rating: 4.8,
    image: "https://images.unsplash.com/photo-1583422409516-2895a77efded?w=800",
    description: "Discover Gaudí's masterpieces and vibrant Catalan culture.",
    highlights: ["Sagrada Familia", "Park Güell", "Las Ramblas", "Gothic Quarter"],
    included: ["Boutique Hotel", "Museum Passes", "Tapas Tour", "Flamenco Show"],
    tags: ["art", "architecture", "culture", "nightlife"]
  }
];

const destinations = [
  { name: "Japan", count: 3, popular: true },
  { name: "Thailand", count: 2, popular: true },
  { name: "France", count: 2, popular: true },
  { name: "USA", count: 2, popular: false },
  { name: "Italy", count: 2, popular: true },
  { name: "Indonesia", count: 1, popular: true },
  { name: "Spain", count: 1, popular: false }
];

// Mock user purchase history
const userPurchases = {
  "user123": [
    { packageId: 1, purchaseDate: "2024-01-15", status: "completed" },
    { packageId: 4, purchaseDate: "2024-02-20", status: "completed" }
  ],
  "user456": [
    { packageId: 6, purchaseDate: "2024-03-10", status: "upcoming" }
  ]
};

function getUserPurchases(userId) {
  const purchases = userPurchases[userId] || [];
  return purchases.map(purchase => {
    const packageDetails = travelPackages.find(p => p.id === purchase.packageId);
    return { ...purchase, package: packageDetails };
  });
}

function getSuggestions(userId) {
  const purchases = userPurchases[userId] || [];
  
  if (purchases.length === 0) {
    // Return popular packages for new users
    return travelPackages.filter(p => p.rating >= 4.7).slice(0, 6);
  }
  
  // Get purchased packages details
  const purchasedPackages = purchases.map(p => 
    travelPackages.find(pkg => pkg.id === p.packageId)
  );
  
  // Extract preferences from purchase history
  const purchasedCountries = purchasedPackages.map(p => p.country);
  const purchasedTypes = purchasedPackages.map(p => p.type);
  const purchasedTags = purchasedPackages.flatMap(p => p.tags);
  
  // Calculate average price
  const avgPrice = purchasedPackages.reduce((sum, p) => sum + p.price, 0) / purchasedPackages.length;
  
  // Score each package based on similarity
  const scoredPackages = travelPackages
    .filter(p => !purchases.some(purchase => purchase.packageId === p.id)) // Exclude already purchased
    .map(pkg => {
      let score = 0;
      
      // Same country (for exploring more in visited countries)
      if (purchasedCountries.includes(pkg.country)) {
        score += 30;
      }
      
      // Similar type preference
      if (purchasedTypes.includes(pkg.type)) {
        score += 25;
      }
      
      // Matching tags
      const matchingTags = pkg.tags.filter(tag => purchasedTags.includes(tag));
      score += matchingTags.length * 10;
      
      // Price similarity (within 30% range)
      const priceDiff = Math.abs(pkg.price - avgPrice) / avgPrice;
      if (priceDiff <= 0.3) {
        score += 20;
      }
      
      // Rating bonus
      score += pkg.rating * 5;
      
      // Diversification bonus (different countries for variety)
      if (!purchasedCountries.includes(pkg.country) && pkg.rating >= 4.5) {
        score += 15;
      }
      
      return { ...pkg, score };
    })
    .sort((a, b) => b.score - a.score);
  
  // Return top suggestions
  return scoredPackages.slice(0, 8);
}

module.exports = {
  travelPackages,
  destinations,
  getUserPurchases,
  getSuggestions
};