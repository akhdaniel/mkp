# Simple Chat Bot

A secure chatbot application with API key protection using a Node.js backend proxy.

## Setup

1. Install dependencies:
```bash
npm install
```

2. Configure environment variables:
   - Copy `.env.example` to `.env`
   - Add your API keys to `.env` file

3. Start the server:
```bash
npm start
```

For development with auto-restart:
```bash
npm run dev
```

4. Open your browser and navigate to:
   - `http://localhost:3000/01-chatbot.html` - Basic chatbot
   - `http://localhost:3000/02-chatbot-tools.html` - Chatbot with weather/currency tools
   - `http://localhost:3000/03-chatbot-trial.html` - Minimal Groq chatbot
   - `http://localhost:3000/04-chatbot-knowledge.html` - Knowledge chatbot

## Security

- API keys are stored in `.env` file (not committed to git)
- Backend proxy handles all API calls
- No API keys exposed in frontend code

## API Services Used

- **DeepSeek**: AI chat completions
- **Groq**: Fast AI inference
- **OpenWeatherMap**: Weather data (optional)
- **Exchange Rate API**: Currency conversion (optional)