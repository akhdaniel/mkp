const express = require('express');
const cors = require('cors');
const bodyParser = require('body-parser');
const { travelPackages, destinations, getUserPurchases, getSuggestions, getAdminStats, getAllUsers } = require('./data/mockData');

const app = express();
const PORT = 5000;

app.use(cors());
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));

// Routes
app.get('/api/packages', (req, res) => {
  res.json(travelPackages);
});

app.get('/api/destinations', (req, res) => {
  res.json(destinations);
});

app.get('/api/packages/:id', (req, res) => {
  const package = travelPackages.find(p => p.id === parseInt(req.params.id));
  if (package) {
    res.json(package);
  } else {
    res.status(404).json({ error: 'Package not found' });
  }
});

app.get('/api/destinations/:country', (req, res) => {
  const countryPackages = travelPackages.filter(
    p => p.country.toLowerCase() === req.params.country.toLowerCase()
  );
  res.json(countryPackages);
});

app.post('/api/search', (req, res) => {
  const { query, priceRange, duration, type } = req.body;
  
  let results = travelPackages;
  
  if (query) {
    results = results.filter(p => 
      p.title.toLowerCase().includes(query.toLowerCase()) ||
      p.destination.toLowerCase().includes(query.toLowerCase()) ||
      p.country.toLowerCase().includes(query.toLowerCase())
    );
  }
  
  if (priceRange) {
    const [min, max] = priceRange;
    results = results.filter(p => p.price >= min && p.price <= max);
  }
  
  if (duration) {
    results = results.filter(p => p.duration === duration);
  }
  
  if (type) {
    results = results.filter(p => p.type === type);
  }
  
  res.json(results);
});

app.get('/api/user/:userId/purchases', (req, res) => {
  const purchases = getUserPurchases(req.params.userId);
  res.json(purchases);
});

app.get('/api/suggestions/:userId', (req, res) => {
  const suggestions = getSuggestions(req.params.userId);
  res.json(suggestions);
});

app.post('/api/user/:userId/purchase', (req, res) => {
  const { packageId } = req.body;
  // In a real app, this would update a database
  res.json({ success: true, message: 'Purchase recorded', packageId });
});

// Admin endpoints
app.get('/api/admin/stats', (req, res) => {
  res.json(getAdminStats());
});

app.get('/api/admin/users', (req, res) => {
  res.json(getAllUsers());
});

app.post('/api/admin/packages', (req, res) => {
  // In a real app, this would add to database
  res.json({ success: true, message: 'Package created', package: req.body });
});

app.put('/api/admin/packages/:id', (req, res) => {
  // In a real app, this would update database
  res.json({ success: true, message: 'Package updated', id: req.params.id });
});

app.delete('/api/admin/packages/:id', (req, res) => {
  // In a real app, this would delete from database
  res.json({ success: true, message: 'Package deleted', id: req.params.id });
});

app.listen(PORT, () => {
  console.log(`Server running on http://localhost:${PORT}`);
});