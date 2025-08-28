import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { travelAPI } from '../services/api';
import { UserPurchase, TravelPackage } from '../types';
import '../styles/UserProfile.css';

const UserProfile = () => {
  const navigate = useNavigate();
  const [purchases, setPurchases] = useState<UserPurchase[]>([]);
  const [suggestions, setSuggestions] = useState<TravelPackage[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'purchases' | 'suggestions'>('purchases');
  
  // Mock user data
  const currentUserId = 'user123';
  const userName = 'John Doe';

  useEffect(() => {
    fetchUserData();
  }, []);

  const fetchUserData = async () => {
    try {
      const [purchasesData, suggestionsData] = await Promise.all([
        travelAPI.getUserPurchases(currentUserId),
        travelAPI.getSuggestions(currentUserId)
      ]);
      setPurchases(purchasesData);
      setSuggestions(suggestionsData);
    } catch (error) {
      console.error('Error fetching user data:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="user-profile">
        <div className="loading">Loading profile...</div>
      </div>
    );
  }

  return (
    <div className="user-profile">
      <div className="profile-header">
        <div className="profile-info">
          <div className="avatar">
            {userName.split(' ').map(n => n[0]).join('')}
          </div>
          <div className="user-details">
            <h1>{userName}</h1>
            <p>Travel Enthusiast</p>
            <p className="member-since">Member since 2024</p>
          </div>
        </div>
        
        <div className="profile-stats">
          <div className="stat">
            <span className="stat-value">{purchases.length}</span>
            <span className="stat-label">Trips Booked</span>
          </div>
          <div className="stat">
            <span className="stat-value">
              {purchases.filter(p => p.package).map(p => p.package!.country)
                .filter((v, i, a) => a.indexOf(v) === i).length}
            </span>
            <span className="stat-label">Countries Visited</span>
          </div>
        </div>
      </div>

      <div className="profile-tabs">
        <button 
          className={`tab ${activeTab === 'purchases' ? 'active' : ''}`}
          onClick={() => setActiveTab('purchases')}
        >
          My Bookings
        </button>
        <button 
          className={`tab ${activeTab === 'suggestions' ? 'active' : ''}`}
          onClick={() => setActiveTab('suggestions')}
        >
          Personalized Suggestions
        </button>
      </div>

      <div className="profile-content">
        {activeTab === 'purchases' ? (
          <div className="purchases-section">
            <h2>Your Travel History</h2>
            {purchases.length > 0 ? (
              <div className="purchases-list">
                {purchases.map((purchase) => (
                  <div key={purchase.packageId} className="purchase-card">
                    {purchase.package && (
                      <>
                        <img 
                          src={purchase.package.image} 
                          alt={purchase.package.title}
                          className="purchase-image"
                        />
                        <div className="purchase-info">
                          <h3>{purchase.package.title}</h3>
                          <p className="purchase-destination">
                            📍 {purchase.package.destination}, {purchase.package.country}
                          </p>
                          <p className="purchase-details">
                            {purchase.package.duration} days • ${purchase.package.price}
                          </p>
                          <p className="purchase-date">
                            Booked on: {new Date(purchase.purchaseDate).toLocaleDateString()}
                          </p>
                          <span className={`status ${purchase.status}`}>
                            {purchase.status}
                          </span>
                        </div>
                        <button 
                          onClick={() => navigate(`/packages/${purchase.package!.id}`)}
                          className="view-button"
                        >
                          View Details
                        </button>
                      </>
                    )}
                  </div>
                ))}
              </div>
            ) : (
              <div className="empty-state">
                <p>No bookings yet. Start exploring our packages!</p>
                <button onClick={() => navigate('/packages')} className="explore-button">
                  Explore Packages
                </button>
              </div>
            )}
          </div>
        ) : (
          <div className="suggestions-section">
            <h2>Recommended for You</h2>
            <p className="suggestions-intro">
              Based on your travel history and preferences, we think you'll love these:
            </p>
            <div className="suggestions-grid">
              {suggestions.slice(0, 6).map(pkg => (
                <div key={pkg.id} className="suggestion-card">
                  {pkg.score && (
                    <div className="match-score">{Math.round(pkg.score)}% Match</div>
                  )}
                  <img src={pkg.image} alt={pkg.title} />
                  <div className="suggestion-content">
                    <h3>{pkg.title}</h3>
                    <p className="suggestion-location">
                      {pkg.destination}, {pkg.country}
                    </p>
                    <p className="suggestion-details">
                      {pkg.duration} days • ${pkg.price}
                    </p>
                    <div className="suggestion-tags">
                      {pkg.tags.slice(0, 2).map(tag => (
                        <span key={tag} className="tag">#{tag}</span>
                      ))}
                    </div>
                    <button 
                      onClick={() => navigate(`/packages/${pkg.id}`)}
                      className="view-suggestion"
                    >
                      Learn More
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default UserProfile;