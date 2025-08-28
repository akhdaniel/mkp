import { Link, useLocation } from 'react-router-dom';
import '../styles/Navbar.css';

const Navbar = () => {
  const location = useLocation();
  
  const isActive = (path: string) => {
    return location.pathname === path ? 'active' : '';
  };

  return (
    <nav className="navbar">
      <div className="nav-container">
        <Link to="/" className="nav-brand">
          <span className="brand-icon">✈️</span>
          TravelWise
        </Link>
        
        <div className="nav-menu">
          <Link to="/" className={`nav-link ${isActive('/')}`}>
            Home
          </Link>
          <Link to="/packages" className={`nav-link ${isActive('/packages')}`}>
            All Packages
          </Link>
          <Link to="/search" className={`nav-link ${isActive('/search')}`}>
            Search
          </Link>
          <Link to="/profile" className={`nav-link ${isActive('/profile')}`}>
            My Profile
          </Link>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;