import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { travelAPI } from '../services/api';
import { TravelPackage } from '../types';
import PackageCard from '../components/PackageCard';
import SuggestionSection from '../components/SuggestionSection';
import '../styles/Home.css';

const Home = () => {
  const [featuredPackages, setFeaturedPackages] = useState<TravelPackage[]>([]);
  const [loading, setLoading] = useState(true);
  
  // Mock user ID - in real app, this would come from authentication
  const currentUserId = 'user123';

  useEffect(() => {
    fetchFeaturedPackages();
  }, []);

  const fetchFeaturedPackages = async () => {
    try {
      const packages = await travelAPI.getAllPackages();
      // Get top-rated packages as featured
      const featured = packages
        .sort((a, b) => b.rating - a.rating)
        .slice(0, 3);
      setFeaturedPackages(featured);
    } catch (error) {
      console.error('Error fetching packages:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="home">
      <section className="hero">
        <div className="hero-content">
          <h1>Discover Your Next Adventure</h1>
          <p>Personalized travel recommendations based on your preferences</p>
          <Link to="/search" className="cta-button">
            Start Exploring
          </Link>
        </div>
        <div className="hero-image">
          <img 
            src="https://images.unsplash.com/photo-1488646953014-85cb44e25828?w=1200" 
            alt="Travel" 
          />
        </div>
      </section>

      <section className="user-message">
        <div className="message-card">
          <h2>Welcome back! 🎉</h2>
          <p>Since you recently booked your trip to <strong>Tokyo</strong>, 
             we've curated some amazing recommendations for you!</p>
        </div>
      </section>

      <SuggestionSection 
        userId={currentUserId} 
        title="Recommended for You"
      />

      <section className="featured-section">
        <h2>Featured Destinations</h2>
        {loading ? (
          <div className="loading">Loading featured packages...</div>
        ) : (
          <div className="featured-grid">
            {featuredPackages.map(pkg => (
              <PackageCard key={pkg.id} package={pkg} />
            ))}
          </div>
        )}
      </section>

      <section className="explore-by-type">
        <h2>Explore by Travel Style</h2>
        <div className="travel-types">
          <Link to="/search?type=beach" className="type-card beach">
            <span className="type-icon">🏖️</span>
            <span className="type-name">Beach</span>
          </Link>
          <Link to="/search?type=adventure" className="type-card adventure">
            <span className="type-icon">🏔️</span>
            <span className="type-name">Adventure</span>
          </Link>
          <Link to="/search?type=city" className="type-card city">
            <span className="type-icon">🏙️</span>
            <span className="type-name">City</span>
          </Link>
          <Link to="/search?type=cultural" className="type-card cultural">
            <span className="type-icon">🏛️</span>
            <span className="type-name">Cultural</span>
          </Link>
          <Link to="/search?type=culinary" className="type-card culinary">
            <span className="type-icon">🍜</span>
            <span className="type-name">Culinary</span>
          </Link>
          <Link to="/search?type=countryside" className="type-card countryside">
            <span className="type-icon">🌾</span>
            <span className="type-name">Countryside</span>
          </Link>
        </div>
      </section>
    </div>
  );
};

export default Home;