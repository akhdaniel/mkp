import { useState, useEffect } from 'react';
import { adminAPI, travelAPI } from '../services/api';
import { TravelPackage } from '../types';
import StatsCard from '../components/admin/StatsCard';
import RevenueChart from '../components/admin/RevenueChart';
import PackageManager from '../components/admin/PackageManager';
import UserManager from '../components/admin/UserManager';
import RecentBookings from '../components/admin/RecentBookings';
import '../styles/AdminDashboard.css';

const AdminDashboard = () => {
  const [activeTab, setActiveTab] = useState<'overview' | 'packages' | 'users'>('overview');
  const [stats, setStats] = useState<any>(null);
  const [packages, setPackages] = useState<TravelPackage[]>([]);
  const [users, setUsers] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      const [statsData, packagesData, usersData] = await Promise.all([
        adminAPI.getStats(),
        travelAPI.getAllPackages(),
        adminAPI.getUsers()
      ]);
      setStats(statsData);
      setPackages(packagesData);
      setUsers(usersData);
    } catch (error) {
      console.error('Error fetching admin data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handlePackageUpdate = async (id: number, data: any) => {
    try {
      await adminAPI.updatePackage(id, data);
      await fetchData(); // Refresh data
      alert('Package updated successfully!');
    } catch (error) {
      console.error('Error updating package:', error);
      alert('Failed to update package');
    }
  };

  const handlePackageDelete = async (id: number) => {
    if (window.confirm('Are you sure you want to delete this package?')) {
      try {
        await adminAPI.deletePackage(id);
        await fetchData(); // Refresh data
        alert('Package deleted successfully!');
      } catch (error) {
        console.error('Error deleting package:', error);
        alert('Failed to delete package');
      }
    }
  };

  if (loading) {
    return (
      <div className="admin-dashboard">
        <div className="loading">Loading admin dashboard...</div>
      </div>
    );
  }

  return (
    <div className="admin-dashboard">
      <div className="dashboard-header">
        <h1>Admin Dashboard</h1>
        <p>Manage your travel business</p>
      </div>

      <div className="dashboard-tabs">
        <button 
          className={`tab ${activeTab === 'overview' ? 'active' : ''}`}
          onClick={() => setActiveTab('overview')}
        >
          📊 Overview
        </button>
        <button 
          className={`tab ${activeTab === 'packages' ? 'active' : ''}`}
          onClick={() => setActiveTab('packages')}
        >
          📦 Packages
        </button>
        <button 
          className={`tab ${activeTab === 'users' ? 'active' : ''}`}
          onClick={() => setActiveTab('users')}
        >
          👥 Users
        </button>
      </div>

      <div className="dashboard-content">
        {activeTab === 'overview' && stats && (
          <div className="overview-section">
            <div className="stats-grid">
              <StatsCard
                title="Total Revenue"
                value={`$${stats.overview.revenue.toLocaleString()}`}
                icon="💰"
                trend="+12.5%"
                trendUp={true}
              />
              <StatsCard
                title="Total Bookings"
                value={stats.overview.totalBookings.toString()}
                icon="🎫"
                trend="+8.2%"
                trendUp={true}
              />
              <StatsCard
                title="Active Users"
                value={stats.overview.totalUsers.toString()}
                icon="👥"
                trend="+15.3%"
                trendUp={true}
              />
              <StatsCard
                title="Available Packages"
                value={stats.overview.totalPackages.toString()}
                icon="📦"
                trend="+2"
                trendUp={true}
              />
            </div>

            <div className="charts-section">
              <div className="chart-container">
                <h2>Revenue Trend</h2>
                <RevenueChart data={stats.monthlyRevenue} />
              </div>

              <div className="top-destinations">
                <h2>Popular Destinations</h2>
                <div className="destination-list">
                  {Object.entries(stats.destinationBookings).map(([country, bookings]: [string, any]) => (
                    <div key={country} className="destination-item">
                      <span className="destination-name">{country}</span>
                      <div className="destination-bar">
                        <div 
                          className="bar-fill" 
                          style={{ width: `${(bookings / 3) * 100}%` }}
                        ></div>
                      </div>
                      <span className="destination-count">{bookings} bookings</span>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            <div className="performance-section">
              <div className="top-packages">
                <h2>Top Performing Packages</h2>
                <div className="package-performance-list">
                  {stats.topPackages.map((pkg: any) => (
                    <div key={pkg.id} className="performance-item">
                      <img src={pkg.image} alt={pkg.title} className="performance-image" />
                      <div className="performance-info">
                        <h4>{pkg.title}</h4>
                        <p>{pkg.destination}, {pkg.country}</p>
                      </div>
                      <div className="performance-stats">
                        <span className="bookings-count">{pkg.bookings} bookings</span>
                        <span className="revenue">${(pkg.price * pkg.bookings).toLocaleString()}</span>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              <RecentBookings bookings={stats.recentBookings} />
            </div>
          </div>
        )}

        {activeTab === 'packages' && (
          <PackageManager 
            packages={packages}
            onUpdate={handlePackageUpdate}
            onDelete={handlePackageDelete}
          />
        )}

        {activeTab === 'users' && (
          <UserManager users={users} />
        )}
      </div>
    </div>
  );
};

export default AdminDashboard;