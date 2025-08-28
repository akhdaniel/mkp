import { useState } from 'react';
import { TravelPackage } from '../../types';
import '../../styles/admin/PackageManager.css';

interface PackageManagerProps {
  packages: TravelPackage[];
  onUpdate: (id: number, data: any) => void;
  onDelete: (id: number) => void;
}

const PackageManager = ({ packages, onUpdate, onDelete }: PackageManagerProps) => {
  const [editingId, setEditingId] = useState<number | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterType, setFilterType] = useState('all');
  const [showAddForm, setShowAddForm] = useState(false);
  const [editForm, setEditForm] = useState<any>({});

  const filteredPackages = packages.filter(pkg => {
    const matchesSearch = pkg.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         pkg.country.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesType = filterType === 'all' || pkg.type === filterType;
    return matchesSearch && matchesType;
  });

  const handleEdit = (pkg: TravelPackage) => {
    setEditingId(pkg.id);
    setEditForm({
      title: pkg.title,
      price: pkg.price,
      duration: pkg.duration,
      description: pkg.description
    });
  };

  const handleSave = () => {
    if (editingId) {
      onUpdate(editingId, editForm);
      setEditingId(null);
      setEditForm({});
    }
  };

  const handleCancel = () => {
    setEditingId(null);
    setEditForm({});
  };

  return (
    <div className="package-manager">
      <div className="manager-header">
        <h2>Package Management</h2>
        <button className="add-package-btn" onClick={() => setShowAddForm(true)}>
          + Add New Package
        </button>
      </div>

      <div className="manager-filters">
        <input
          type="text"
          placeholder="Search packages..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="search-input"
        />
        <select 
          value={filterType} 
          onChange={(e) => setFilterType(e.target.value)}
          className="filter-select"
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

      <div className="packages-table-container">
        <table className="packages-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Image</th>
              <th>Title</th>
              <th>Destination</th>
              <th>Type</th>
              <th>Duration</th>
              <th>Price</th>
              <th>Rating</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {filteredPackages.map(pkg => (
              <tr key={pkg.id}>
                <td>{pkg.id}</td>
                <td>
                  <img src={pkg.image} alt={pkg.title} className="table-image" />
                </td>
                <td>
                  {editingId === pkg.id ? (
                    <input
                      type="text"
                      value={editForm.title}
                      onChange={(e) => setEditForm({...editForm, title: e.target.value})}
                      className="edit-input"
                    />
                  ) : (
                    pkg.title
                  )}
                </td>
                <td>{pkg.destination}, {pkg.country}</td>
                <td>
                  <span className={`type-badge ${pkg.type}`}>{pkg.type}</span>
                </td>
                <td>
                  {editingId === pkg.id ? (
                    <input
                      type="number"
                      value={editForm.duration}
                      onChange={(e) => setEditForm({...editForm, duration: parseInt(e.target.value)})}
                      className="edit-input small"
                    />
                  ) : (
                    `${pkg.duration} days`
                  )}
                </td>
                <td>
                  {editingId === pkg.id ? (
                    <input
                      type="number"
                      value={editForm.price}
                      onChange={(e) => setEditForm({...editForm, price: parseInt(e.target.value)})}
                      className="edit-input small"
                    />
                  ) : (
                    `$${pkg.price}`
                  )}
                </td>
                <td>⭐ {pkg.rating}</td>
                <td>
                  {editingId === pkg.id ? (
                    <div className="action-buttons">
                      <button onClick={handleSave} className="save-btn">Save</button>
                      <button onClick={handleCancel} className="cancel-btn">Cancel</button>
                    </div>
                  ) : (
                    <div className="action-buttons">
                      <button onClick={() => handleEdit(pkg)} className="edit-btn">Edit</button>
                      <button onClick={() => onDelete(pkg.id)} className="delete-btn">Delete</button>
                    </div>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {showAddForm && (
        <div className="modal-overlay">
          <div className="add-package-modal">
            <h3>Add New Package</h3>
            <form onSubmit={(e) => {
              e.preventDefault();
              // Handle form submission
              setShowAddForm(false);
            }}>
              <div className="form-group">
                <label>Title</label>
                <input type="text" required />
              </div>
              <div className="form-group">
                <label>Destination</label>
                <input type="text" required />
              </div>
              <div className="form-group">
                <label>Country</label>
                <input type="text" required />
              </div>
              <div className="form-group">
                <label>Type</label>
                <select required>
                  <option value="city">City</option>
                  <option value="beach">Beach</option>
                  <option value="adventure">Adventure</option>
                  <option value="cultural">Cultural</option>
                  <option value="culinary">Culinary</option>
                  <option value="countryside">Countryside</option>
                </select>
              </div>
              <div className="form-group">
                <label>Duration (days)</label>
                <input type="number" required />
              </div>
              <div className="form-group">
                <label>Price</label>
                <input type="number" required />
              </div>
              <div className="form-group">
                <label>Description</label>
                <textarea required></textarea>
              </div>
              <div className="form-actions">
                <button type="submit" className="submit-btn">Add Package</button>
                <button type="button" onClick={() => setShowAddForm(false)} className="cancel-btn">
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default PackageManager;