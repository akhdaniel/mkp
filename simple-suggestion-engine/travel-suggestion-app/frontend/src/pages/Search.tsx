import { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import { travelAPI } from '../services/api';
import { TravelPackage, SearchFilters } from '../types';
import PackageCard from '../components/PackageCard';
import '../styles/Search.css';

const Search = () => {
  const [searchParams] = useSearchParams();
  const [packages, setPackages] = useState<TravelPackage[]>([]);
  const [loading, setLoading] = useState(false);
  const [filters, setFilters] = useState<SearchFilters>({
    query: '',
    priceRange: [0, 5000],
    duration: undefined,
    type: undefined,
  });

  useEffect(() => {
    // Check for URL parameters
    const type = searchParams.get('type');
    if (type) {
      setFilters(prev => ({ ...prev, type }));
      performSearch({ ...filters, type });
    }
  }, [searchParams]);

  const performSearch = async (searchFilters: SearchFilters) => {
    setLoading(true);
    try {
      const results = await travelAPI.searchPackages(searchFilters);
      setPackages(results);
    } catch (error) {
      console.error('Error searching packages:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    performSearch(filters);
  };

  const handleFilterChange = (key: keyof SearchFilters, value: any) => {
    setFilters(prev => ({ ...prev, [key]: value }));
  };

  return (
    <div className="search-page">
      <div className="search-header">
        <h1>Find Your Perfect Trip</h1>
        <p>Search and filter through our extensive collection of travel packages</p>
      </div>

      <div className="search-container">
        <form className="search-form" onSubmit={handleSearch}>
          <div className="search-bar">
            <input
              type="text"
              placeholder="Search destinations, countries, or keywords..."
              value={filters.query}
              onChange={(e) => handleFilterChange('query', e.target.value)}
              className="search-input"
            />
            <button type="submit" className="search-button">
              Search
            </button>
          </div>

          <div className="filters">
            <div className="filter-group">
              <label>Price Range</label>
              <div className="price-inputs">
                <input
                  type="number"
                  placeholder="Min"
                  min="0"
                  value={filters.priceRange?.[0] || 0}
                  onChange={(e) => handleFilterChange('priceRange', [
                    parseInt(e.target.value) || 0,
                    filters.priceRange?.[1] || 5000
                  ])}
                  className="price-input"
                />
                <span>-</span>
                <input
                  type="number"
                  placeholder="Max"
                  min="0"
                  value={filters.priceRange?.[1] || 5000}
                  onChange={(e) => handleFilterChange('priceRange', [
                    filters.priceRange?.[0] || 0,
                    parseInt(e.target.value) || 5000
                  ])}
                  className="price-input"
                />
              </div>
            </div>

            <div className="filter-group">
              <label htmlFor="duration">Duration (days)</label>
              <select
                id="duration"
                value={filters.duration || ''}
                onChange={(e) => handleFilterChange('duration', e.target.value ? parseInt(e.target.value) : undefined)}
                className="filter-select"
              >
                <option value="">Any duration</option>
                <option value="3">3 days</option>
                <option value="4">4 days</option>
                <option value="5">5 days</option>
                <option value="6">6 days</option>
                <option value="7">7 days</option>
              </select>
            </div>

            <div className="filter-group">
              <label htmlFor="type">Travel Type</label>
              <select
                id="type"
                value={filters.type || ''}
                onChange={(e) => handleFilterChange('type', e.target.value || undefined)}
                className="filter-select"
              >
                <option value="">All types</option>
                <option value="city">City</option>
                <option value="beach">Beach</option>
                <option value="adventure">Adventure</option>
                <option value="cultural">Cultural</option>
                <option value="culinary">Culinary</option>
                <option value="countryside">Countryside</option>
              </select>
            </div>
          </div>
        </form>
      </div>

      <div className="search-results">
        {loading ? (
          <div className="loading">Searching...</div>
        ) : packages.length > 0 ? (
          <>
            <div className="results-count">
              Found {packages.length} package{packages.length !== 1 ? 's' : ''}
            </div>
            <div className="results-grid">
              {packages.map(pkg => (
                <PackageCard key={pkg.id} package={pkg} />
              ))}
            </div>
          </>
        ) : (
          <div className="no-results">
            <p>No packages found matching your criteria.</p>
            <p>Try adjusting your filters or search terms.</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default Search;