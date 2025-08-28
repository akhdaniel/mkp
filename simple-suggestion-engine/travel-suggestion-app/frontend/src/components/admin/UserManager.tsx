import { useState } from 'react';
import '../../styles/admin/UserManager.css';

interface User {
  id: string;
  name: string;
  email: string;
  joinDate: string;
  totalBookings: number;
  totalSpent: number;
  status: 'active' | 'inactive';
}

interface UserManagerProps {
  users: User[];
}

const UserManager = ({ users }: UserManagerProps) => {
  const [searchTerm, setSearchTerm] = useState('');
  const [filterStatus, setFilterStatus] = useState('all');
  const [sortBy, setSortBy] = useState('joinDate');

  const filteredUsers = users
    .filter(user => {
      const matchesSearch = user.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                           user.email.toLowerCase().includes(searchTerm.toLowerCase());
      const matchesStatus = filterStatus === 'all' || user.status === filterStatus;
      return matchesSearch && matchesStatus;
    })
    .sort((a, b) => {
      switch (sortBy) {
        case 'name':
          return a.name.localeCompare(b.name);
        case 'totalSpent':
          return b.totalSpent - a.totalSpent;
        case 'totalBookings':
          return b.totalBookings - a.totalBookings;
        case 'joinDate':
        default:
          return new Date(b.joinDate).getTime() - new Date(a.joinDate).getTime();
      }
    });

  const totalRevenue = users.reduce((sum, user) => sum + user.totalSpent, 0);
  const activeUsers = users.filter(u => u.status === 'active').length;

  return (
    <div className="user-manager">
      <div className="manager-header">
        <h2>User Management</h2>
        <div className="user-stats">
          <span>Total Users: {users.length}</span>
          <span>Active: {activeUsers}</span>
          <span>Total Revenue: ${totalRevenue.toLocaleString()}</span>
        </div>
      </div>

      <div className="manager-controls">
        <input
          type="text"
          placeholder="Search users..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="search-input"
        />
        <select 
          value={filterStatus} 
          onChange={(e) => setFilterStatus(e.target.value)}
          className="filter-select"
        >
          <option value="all">All Status</option>
          <option value="active">Active</option>
          <option value="inactive">Inactive</option>
        </select>
        <select 
          value={sortBy} 
          onChange={(e) => setSortBy(e.target.value)}
          className="sort-select"
        >
          <option value="joinDate">Join Date</option>
          <option value="name">Name</option>
          <option value="totalSpent">Total Spent</option>
          <option value="totalBookings">Total Bookings</option>
        </select>
      </div>

      <div className="users-grid">
        {filteredUsers.map(user => (
          <div key={user.id} className="user-card">
            <div className="user-header">
              <div className="user-avatar">
                {user.name.split(' ').map(n => n[0]).join('')}
              </div>
              <div className="user-info">
                <h3>{user.name}</h3>
                <p>{user.email}</p>
              </div>
              <span className={`status-badge ${user.status}`}>
                {user.status}
              </span>
            </div>
            
            <div className="user-stats-grid">
              <div className="stat">
                <span className="stat-label">Member Since</span>
                <span className="stat-value">
                  {new Date(user.joinDate).toLocaleDateString()}
                </span>
              </div>
              <div className="stat">
                <span className="stat-label">Total Bookings</span>
                <span className="stat-value">{user.totalBookings}</span>
              </div>
              <div className="stat">
                <span className="stat-label">Total Spent</span>
                <span className="stat-value">${user.totalSpent.toLocaleString()}</span>
              </div>
              <div className="stat">
                <span className="stat-label">Avg. Booking</span>
                <span className="stat-value">
                  ${user.totalBookings > 0 ? Math.round(user.totalSpent / user.totalBookings) : 0}
                </span>
              </div>
            </div>

            <div className="user-actions">
              <button className="view-btn">View Details</button>
              <button className="message-btn">Send Message</button>
              {user.status === 'active' ? (
                <button className="suspend-btn">Suspend</button>
              ) : (
                <button className="activate-btn">Activate</button>
              )}
            </div>
          </div>
        ))}
      </div>

      {filteredUsers.length === 0 && (
        <div className="no-users">
          No users found matching your criteria
        </div>
      )}
    </div>
  );
};

export default UserManager;