# Travel Suggestion Engine

A personalized travel recommendation system inspired by Trip.com's suggestion engine. This application provides tailored travel package recommendations based on user preferences and purchase history.

## Features

- **Smart Recommendation Engine**: Suggests travel packages based on user's previous purchases and preferences
- **Search & Filter**: Advanced search with filters for price, duration, and travel type
- **User Profiles**: Track bookings and get personalized suggestions
- **Detailed Package Information**: Comprehensive details about each travel package
- **Responsive Design**: Works seamlessly on desktop and mobile devices

## Tech Stack

### Backend
- Node.js
- Express.js
- CORS for cross-origin requests
- Mock data simulating a real database

### Frontend
- React with TypeScript
- React Router for navigation
- Axios for API calls
- Vite for fast development
- Custom CSS for styling

## Installation & Setup

### Prerequisites
- Node.js (v14 or higher)
- npm or yarn

### Backend Setup

1. Navigate to the backend directory:
```bash
cd travel-suggestion-app/backend
```

2. Install dependencies:
```bash
npm install
```

3. Start the backend server:
```bash
npm start
```

The backend server will run on `http://localhost:5000`

### Frontend Setup

1. Open a new terminal and navigate to the frontend directory:
```bash
cd travel-suggestion-app/frontend
```

2. Install dependencies:
```bash
npm install
```

3. Start the development server:
```bash
npm run dev
```

The frontend will run on `http://localhost:3000`

## How the Suggestion Engine Works

The recommendation algorithm considers multiple factors:

1. **Purchase History**: Analyzes previously purchased packages
2. **Country Preferences**: Suggests more packages in visited countries
3. **Travel Type Preferences**: Matches user's preferred travel styles (beach, adventure, cultural, etc.)
4. **Price Range**: Recommends packages within similar price ranges
5. **Tag Matching**: Finds packages with similar tags (food, culture, nature, etc.)
6. **Diversification**: Also suggests new destinations for variety

Each suggestion gets a score based on these factors, and the top-scoring packages are displayed.

## API Endpoints

### Packages
- `GET /api/packages` - Get all travel packages
- `GET /api/packages/:id` - Get specific package details
- `GET /api/destinations` - Get all destinations
- `GET /api/destinations/:country` - Get packages for a specific country

### Search & Suggestions
- `POST /api/search` - Search packages with filters
- `GET /api/suggestions/:userId` - Get personalized suggestions for a user

### User
- `GET /api/user/:userId/purchases` - Get user's purchase history
- `POST /api/user/:userId/purchase` - Record a new purchase

## Mock Users

The application includes mock data for testing:
- User ID: `user123` - Has purchased Tokyo and Bangkok packages
- User ID: `user456` - Has purchased Paris package

## Features Demo

1. **Homepage**: Shows personalized recommendations based on user's purchase history
2. **All Packages**: Browse all available travel packages with filters
3. **Search**: Advanced search with multiple filter options
4. **Package Details**: View comprehensive information about each package
5. **User Profile**: See purchase history and personalized suggestions

## Future Enhancements

- Real database integration (MongoDB/PostgreSQL)
- User authentication and registration
- Payment gateway integration
- Reviews and ratings system
- Real-time availability checking
- Email notifications
- Admin panel for package management
- Machine learning for improved recommendations

## License

MIT