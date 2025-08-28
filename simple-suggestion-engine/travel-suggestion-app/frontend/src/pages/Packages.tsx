import { useState, useEffect } from 'react';
import { travelAPI } from '../services/api';
import { TravelPackage, Destination } from '../types';
import PackageCard from '../components/PackageCard';
import '../styles/Packages.css';

const Packages = () => {
  const [packages, setPackages] = useState<TravelPackage[]>([]);
  const [filteredPackages, setFilteredPackages] = useState<TravelPackage[]>([]);
  const [destinations, setDestinations] = useState<Destination[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedCountry, setSelectedCountry] = useState<string>('all');
  const [selectedType, setSelectedType] = useState<string>('all');
  const [priceSort, setPriceSort] = useState<string>('none');

  useEffect(() => {
    fetchData();
  }, []);

  useEffect(() => {
    filterAndSortPackages();
  }, [packages, selectedCountry, selectedType, priceSort]);

  const fetchData = async () => {
    try {
      const [packagesData, destinationsData] = await Promise.all([
        travelAPI.getAllPackages(),
        travelAPI.getDestinations()
      ]);
      setPackages(packagesData);
      setFilteredPackages(packagesData);
      setDestinations(destinationsData);
    } catch (error) {
      console.error('Error fetching data:', error);
    } finally {
      setLoading(false);
    }
  };

  const filterAndSortPackages = () => {
    let filtered = [...packages];

    if (selectedCountry !== 'all') {
      filtered = filtered.filter(pkg => pkg.country === selectedCountry);
    }

    if (selectedType !== 'all') {
      filtered = filtered.filter(pkg => pkg.type === selectedType);
    }

    if (priceSort === 'low') {
      filtered.sort((a, b) => a.price - b.price);
    } else if (priceSort === 'high') {
      filtered.sort((a, b) => b.price - a.price);
    }

    setFilteredPackages(filtered);
  };

  if (loading) {
    return (
      <div className="packages-page">
        <div className="loading">Loading packages...</div>
      </div>
    );
  }

  return (
    <div className="packages-page">
      <div className="page-header">
        <h1>All Travel Packages</h1>
        <p>Explore our complete collection of travel experiences</p>
      </div>

      <div className="filters-section">
        <div className="filter-group">
          <label htmlFor="country-filter">Destination:</label>
          <select 
            id="country-filter"
            value={selectedCountry} 
            onChange={(e) => setSelectedCountry(e.target.value)}
          >
            <option value="all">All Countries</option>
            {destinations.map(dest => (
              <option key={dest.name} value={dest.name}>
                {dest.name} ({dest.count})
              </option>
            ))}
          </select>
        </div>

        <div className="filter-group">
          <label htmlFor="type-filter">Travel Type:</label>
          <select 
            id="type-filter"
            value={selectedType} 
            onChange={(e) => setSelectedType(e.target.value)}
          >
            <option value="all">All Types</option>
            <option value="city">City</option>
            <option value="beach">Beach</option>
            <option value="adventure">Adventure</option>
            <option value="cultural">Cultural</option>
            <option value="culinary">Culinary</option>
            <option value="countryside">Countryside</option>
          </select>
        </div>

        <div className="filter-group">
          <label htmlFor="price-sort">Sort by Price:</label>
          <select 
            id="price-sort"
            value={priceSort} 
            onChange={(e) => setPriceSort(e.target.value)}
          >
            <option value="none">Default</option>
            <option value="low">Low to High</option>
            <option value="high">High to Low</option>
          </select>
        </div>
      </div>

      <div className="results-info">
        Showing {filteredPackages.length} packages
      </div>

      <div className="packages-grid">
        {filteredPackages.map(pkg => (
          <PackageCard key={pkg.id} package={pkg} />
        ))}
      </div>

      {filteredPackages.length === 0 && (
        <div className="no-results">
          No packages found matching your filters. Try adjusting your selection.
        </div>
      )}
    </div>
  );
};

export default Packages;