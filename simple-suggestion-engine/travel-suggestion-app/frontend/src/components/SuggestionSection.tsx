import { useState, useEffect } from 'react';
import { travelAPI } from '../services/api';
import { TravelPackage } from '../types';
import PackageCard from './PackageCard';
import '../styles/SuggestionSection.css';

interface SuggestionSectionProps {
  userId: string;
  title: string;
}

const SuggestionSection = ({ userId, title }: SuggestionSectionProps) => {
  const [suggestions, setSuggestions] = useState<TravelPackage[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchSuggestions();
  }, [userId]);

  const fetchSuggestions = async () => {
    try {
      setLoading(true);
      const data = await travelAPI.getSuggestions(userId);
      setSuggestions(data);
      setError(null);
    } catch (err) {
      console.error('Error fetching suggestions:', err);
      setError('Failed to load suggestions');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="suggestion-section">
        <h2>{title}</h2>
        <div className="loading">Loading personalized suggestions...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="suggestion-section">
        <h2>{title}</h2>
        <div className="error">{error}</div>
      </div>
    );
  }

  return (
    <div className="suggestion-section">
      <h2>{title}</h2>
      <div className="suggestion-grid">
        {suggestions.slice(0, 6).map(pkg => (
          <PackageCard key={pkg.id} package={pkg} showScore={true} />
        ))}
      </div>
    </div>
  );
};

export default SuggestionSection;