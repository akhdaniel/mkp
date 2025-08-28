#!/bin/bash

echo "Starting DataInsight AI..."
echo ""

# Backend startup
echo "Starting backend server..."
cd backend
if [ ! -d "venv" ]; then
    echo "Creating virtual environment..."
    python3 -m venv venv
fi

source venv/bin/activate
echo "Installing/updating backend dependencies..."
pip install -q -r requirements.txt

# Start backend in background
uvicorn main:app --reload --port 8000 &
BACKEND_PID=$!
echo "Backend started on http://localhost:8000 (PID: $BACKEND_PID)"

# Frontend startup
echo ""
echo "Starting frontend server..."
cd ../frontend

# Install dependencies if node_modules doesn't exist
if [ ! -d "node_modules" ]; then
    echo "Installing frontend dependencies..."
    npm install
fi

# Start frontend
echo "Frontend starting on http://localhost:3000"
npm start

# Kill backend when frontend is stopped
kill $BACKEND_PID