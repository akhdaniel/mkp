require('dotenv').config();
const express = require('express');
const cors = require('cors');
const axios = require('axios');
const path = require('path');

const app = express();
const PORT = process.env.PORT || 3000;

// Middleware
app.use(cors());
app.use(express.json());
app.use(express.static('.'));

// API proxy for DeepSeek
app.post('/api/deepseek', async (req, res) => {
  try {
    const response = await axios.post('https://api.deepseek.com/v1/chat/completions', req.body, {
      headers: {
        'Authorization': `Bearer ${process.env.DEEPSEEK_API_KEY}`,
        'Content-Type': 'application/json'
      }
    });
    res.json(response.data);
  } catch (error) {
    console.error('DeepSeek API Error:', error.response?.data || error.message);
    res.status(error.response?.status || 500).json({ 
      error: error.response?.data || { message: 'API request failed' }
    });
  }
});

// API proxy for Groq
app.post('/api/groq', async (req, res) => {
  try {
    const response = await axios.post('https://api.groq.com/openai/v1/chat/completions', req.body, {
      headers: {
        'Authorization': `Bearer ${process.env.GROQ_API_KEY}`,
        'Content-Type': 'application/json'
      }
    });
    res.json(response.data);
  } catch (error) {
    console.error('Groq API Error:', error.response?.data || error.message);
    res.status(error.response?.status || 500).json({ 
      error: error.response?.data || { message: 'API request failed' }
    });
  }
});

// Health check endpoint
app.get('/api/health', (req, res) => {
  res.json({ status: 'ok', message: 'Server is running' });
});

app.listen(PORT, () => {
  console.log(`Server is running on http://localhost:${PORT}`);
  console.log('API endpoints available:');
  console.log('  - POST /api/deepseek');
  console.log('  - POST /api/groq');
  console.log('  - GET /api/health');
});