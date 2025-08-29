#!/usr/bin/env python
"""
Test script to verify the Knowledge Management System components work
"""
import os
import sys

def test_imports():
    """Test that all core imports work"""
    print("Testing imports...")
    try:
        import streamlit as st
        print("✓ Streamlit imported successfully")
        
        import anthropic
        print(f"✓ Anthropic imported successfully (version {anthropic.__version__})")
        
        import chromadb
        print("✓ ChromaDB imported successfully")
        
        import langchain
        print(f"✓ Langchain imported successfully (version {langchain.__version__})")
        
        from backend.config import Settings
        print("✓ Backend config imported successfully")
        
        from backend.document_processor import DocumentProcessor
        print("✓ Document processor imported successfully")
        
        from backend.vector_store import VectorStore
        print("✓ Vector store imported successfully")
        
        from backend.rag_engine import RAGEngine
        print("✓ RAG engine imported successfully")
        
        return True
    except ImportError as e:
        print(f"✗ Import failed: {e}")
        return False

def test_config():
    """Test configuration loading"""
    print("\nTesting configuration...")
    try:
        from backend.config import Settings
        settings = Settings()
        print(f"✓ Settings loaded successfully")
        print(f"  - Vector DB: {settings.vector_db_type}")
        print(f"  - Chunk size: {settings.chunk_size}")
        print(f"  - Model: {settings.model_name}")
        return True
    except Exception as e:
        print(f"✗ Configuration failed: {e}")
        return False

def test_api_key():
    """Test API key configuration"""
    print("\nTesting API keys...")
    api_key = os.getenv("ANTHROPIC_API_KEY", "").strip()
    if not api_key or api_key == "test-api-key" or api_key == "your_anthropic_api_key_here":
        print("⚠ Warning: ANTHROPIC_API_KEY not configured")
        print("  Please set your actual API key in the .env file")
        return False
    else:
        print("✓ ANTHROPIC_API_KEY is configured")
        return True

def main():
    """Run all tests"""
    print("=" * 50)
    print("Knowledge Management System - Component Test")
    print("=" * 50)
    
    results = []
    results.append(test_imports())
    results.append(test_config())
    results.append(test_api_key())
    
    print("\n" + "=" * 50)
    if all(results):
        print("✅ All tests passed! The application is ready to run.")
        print("\nTo start the application, run:")
        print("  streamlit run app.py")
    else:
        print("⚠️ Some tests failed. Please check the errors above.")
        if not test_api_key():
            print("\n📝 Next step: Add your Anthropic API key to the .env file")
            print("  Edit .env and replace 'your_anthropic_api_key_here' with your actual key")
    print("=" * 50)

if __name__ == "__main__":
    main()