import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { travelAPI } from '../services/api';
import { TravelPackage } from '../types';
import '../styles/PackageDetail.css';

const PackageDetail = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [packageData, setPackageData] = useState<TravelPackage | null>(null);
  const [loading, setLoading] = useState(true);
  const [relatedPackages, setRelatedPackages] = useState<TravelPackage[]>([]);
  
  // Mock user ID
  const currentUserId = 'user123';

  useEffect(() => {
    if (id) {
      fetchPackageDetails(parseInt(id));
    }
  }, [id]);

  const fetchPackageDetails = async (packageId: number) => {
    try {
      setLoading(true);
      const data = await travelAPI.getPackageById(packageId);
      setPackageData(data);
      
      // Fetch related packages from same country
      const countryPackages = await travelAPI.getPackagesByCountry(data.country);
      setRelatedPackages(countryPackages.filter(p => p.id !== packageId).slice(0, 3));
    } catch (error) {
      console.error('Error fetching package details:', error);
      navigate('/packages');
    } finally {
      setLoading(false);
    }
  };

  const handlePurchase = async () => {
    if (!packageData) return;
    
    try {
      await travelAPI.purchasePackage(currentUserId, packageData.id);
      alert('Package booked successfully! 🎉');
      navigate('/profile');
    } catch (error) {
      console.error('Error purchasing package:', error);
      alert('Failed to book package. Please try again.');
    }
  };

  if (loading) {
    return (
      <div className="package-detail">
        <div className="loading">Loading package details...</div>
      </div>
    );
  }

  if (!packageData) {
    return (
      <div className="package-detail">
        <div className="error">Package not found</div>
      </div>
    );
  }

  return (
    <div className="package-detail">
      <button onClick={() => navigate(-1)} className="back-button">
        ← Back
      </button>

      <div className="detail-hero">
        <img src={packageData.image} alt={packageData.title} className="hero-image" />
        <div className="hero-overlay">
          <h1>{packageData.title}</h1>
          <p className="destination">
            📍 {packageData.destination}, {packageData.country}
          </p>
        </div>
      </div>

      <div className="detail-content">
        <div className="main-info">
          <div className="info-header">
            <div className="badges">
              <span className="badge type">{packageData.type}</span>
              <span className="badge duration">{packageData.duration} days</span>
              <span className="badge rating">⭐ {packageData.rating}</span>
            </div>
          </div>

          <div className="description-section">
            <h2>About This Trip</h2>
            <p>{packageData.description}</p>
          </div>

          <div className="highlights-section">
            <h2>Trip Highlights</h2>
            <ul className="highlights-list">
              {packageData.highlights.map((highlight, index) => (
                <li key={index}>
                  <span className="highlight-icon">✓</span>
                  {highlight}
                </li>
              ))}
            </ul>
          </div>

          <div className="included-section">
            <h2>What's Included</h2>
            <div className="included-grid">
              {packageData.included.map((item, index) => (
                <div key={index} className="included-item">
                  <span className="included-icon">✅</span>
                  <span>{item}</span>
                </div>
              ))}
            </div>
          </div>

          <div className="tags-section">
            <h3>Perfect for:</h3>
            <div className="tags">
              {packageData.tags.map(tag => (
                <span key={tag} className="tag">#{tag}</span>
              ))}
            </div>
          </div>
        </div>

        <div className="booking-sidebar">
          <div className="price-card">
            <div className="price-header">
              <span className="price-label">Starting from</span>
              <div className="price-amount">${packageData.price}</div>
              <span className="price-per">per person</span>
            </div>
            
            <button onClick={handlePurchase} className="book-button">
              Book This Package
            </button>
            
            <div className="booking-info">
              <p>✈️ Instant confirmation</p>
              <p>💳 Secure payment</p>
              <p>🔄 Free cancellation up to 24h</p>
            </div>
          </div>
        </div>
      </div>

      {relatedPackages.length > 0 && (
        <div className="related-section">
          <h2>More in {packageData.country}</h2>
          <div className="related-grid">
            {relatedPackages.map(pkg => (
              <div key={pkg.id} className="related-card" onClick={() => navigate(`/packages/${pkg.id}`)}>
                <img src={pkg.image} alt={pkg.title} />
                <div className="related-info">
                  <h4>{pkg.title}</h4>
                  <p>{pkg.duration} days • ${pkg.price}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default PackageDetail;