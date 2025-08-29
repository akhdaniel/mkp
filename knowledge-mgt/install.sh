#!/bin/bash

# Installation script for macOS to avoid compilation issues

echo "Installing AI Helpdesk Chat dependencies for macOS..."

# Upgrade pip and install wheel first
pip install --upgrade pip wheel setuptools

# Install core dependencies that don't require compilation
echo "Installing core dependencies..."
pip install --no-cache-dir streamlit==1.31.0
pip install --no-cache-dir anthropic==0.18.1
pip install --no-cache-dir python-dotenv==1.0.1
pip install --no-cache-dir pydantic==2.5.3
pip install --no-cache-dir pydantic-settings==2.1.0

# Install document processing libraries
echo "Installing document processing libraries..."
pip install --no-cache-dir pypdf2==3.0.1
pip install --no-cache-dir python-docx==1.1.0
pip install --no-cache-dir python-pptx==0.6.23
pip install --no-cache-dir openpyxl==3.1.2

# Install pandas with pre-built wheels
echo "Installing pandas..."
pip install --no-cache-dir --only-binary :all: pandas==2.2.3

# Install numpy with pre-built wheels
echo "Installing numpy..."
pip install --no-cache-dir --only-binary :all: numpy==2.1.3

# Install ChromaDB
echo "Installing ChromaDB..."
pip install --no-cache-dir chromadb==0.4.22

# Install embeddings
echo "Installing sentence-transformers..."
pip install --no-cache-dir sentence-transformers==2.3.1
pip install --no-cache-dir tiktoken==0.5.2

# Install LangChain (may take a moment)
echo "Installing LangChain components..."
pip install --no-cache-dir langchain==0.1.7
pip install --no-cache-dir langchain-community==0.0.20
pip install --no-cache-dir langchain-anthropic==0.1.1

# Optional: Install OpenAI if needed
# pip install --no-cache-dir openai==1.12.0

echo "Installation complete!"
echo "To run the application: streamlit run app.py"