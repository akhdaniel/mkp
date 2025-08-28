import { Link } from 'react-router-dom';
import { TravelPackage } from '../types';
import '../styles/PackageCard.css';

interface PackageCardProps {
  package: TravelPackage;
  showScore?: boolean;
}

const PackageCard = ({ package: pkg, showScore = false }: PackageCardProps) => {
  return (
    <div className="package-card">
      {showScore && pkg.score && (
        <div className="recommendation-score">
          {Math.round(pkg.score)}% Match
        </div>
      )}
      
      <div className="package-image-container">
        <img src={pkg.image} alt={pkg.title} className="package-image" />
        <div className="package-overlay">
          <span className="package-type">{pkg.type}</span>
          <span className="package-duration">{pkg.duration} days</span>
        </div>
      </div>
      
      <div className="package-content">
        <div className="package-header">
          <h3 className="package-title">{pkg.title}</h3>
          <div className="package-rating">
            ⭐ {pkg.rating}
          </div>
        </div>
        
        <p className="package-location">
          📍 {pkg.destination}, {pkg.country}
        </p>
        
        <p className="package-description">
          {pkg.description.substring(0, 100)}...
        </p>
        
        <div className="package-tags">
          {pkg.tags.slice(0, 3).map(tag => (
            <span key={tag} className="tag">#{tag}</span>
          ))}
        </div>
        
        <div className="package-footer">
          <div className="package-price">
            <span className="price-label">From</span>
            <span className="price-amount">${pkg.price}</span>
          </div>
          
          <Link to={`/packages/${pkg.id}`} className="view-details-btn">
            View Details
          </Link>
        </div>
      </div>
    </div>
  );
};

export default PackageCard;