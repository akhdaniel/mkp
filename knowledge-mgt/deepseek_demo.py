#!/usr/bin/env python
"""
Demo script showing how to use DeepSeek with the Knowledge Management System
"""
import os
from dotenv import load_dotenv

# Load environment variables
load_dotenv()


def test_deepseek_direct():
    """Test DeepSeek API directly"""
    print("\n" + "="*50)
    print("Testing DeepSeek API Directly")
    print("="*50)
    
    from openai import OpenAI
    
    api_key = os.getenv("DEEPSEEK_API_KEY")
    if not api_key:
        print("❌ DEEPSEEK_API_KEY not found in .env file")
        print("Please add: DEEPSEEK_API_KEY=your_key_here")
        return False
    
    try:
        client = OpenAI(
            api_key=api_key,
            base_url="https://api.deepseek.com/v1"
        )
        
        response = client.chat.completions.create(
            model="deepseek-chat",
            messages=[
                {"role": "system", "content": "You are a helpful assistant."},
                {"role": "user", "content": "Hello! Can you confirm you're working?"}
            ],
            max_tokens=100
        )
        
        print("✅ DeepSeek API is working!")
        print(f"Response: {response.choices[0].message.content}")
        return True
        
    except Exception as e:
        print(f"❌ Error connecting to DeepSeek: {e}")
        return False


def test_deepseek_rag_engine():
    """Test DeepSeek through the RAG engine"""
    print("\n" + "="*50)
    print("Testing DeepSeek RAG Engine")
    print("="*50)
    
    from backend.deepseek_engine import DeepSeekEngine
    
    try:
        engine = DeepSeekEngine()
        print("✅ DeepSeek RAG Engine initialized")
        
        # Test response generation (without actual documents)
        response = engine.generate_response(
            query="What is the purpose of this system?",
            context=[{
                "content": "This is a knowledge management system for internal helpdesk use.",
                "metadata": {"source": "README.md", "page": 1}
            }]
        )
        
        print(f"✅ Response generated: {response[:200]}...")
        return True
        
    except Exception as e:
        print(f"❌ Error with DeepSeek RAG Engine: {e}")
        return False


def test_multi_llm_engine():
    """Test the Multi-LLM RAG Engine with DeepSeek"""
    print("\n" + "="*50)
    print("Testing Multi-LLM RAG Engine with DeepSeek")
    print("="*50)
    
    from backend.multi_llm_rag_engine import MultiLLMRAGEngine
    
    try:
        # Test with DeepSeek
        engine = MultiLLMRAGEngine(provider="deepseek")
        info = engine.get_provider_info()
        print(f"✅ Multi-LLM Engine initialized with {info['provider']}")
        print(f"   Model: {info['model']}")
        
        # You can switch providers
        print("\nSwitching to Anthropic...")
        engine.switch_provider("anthropic")
        info = engine.get_provider_info()
        print(f"✅ Switched to {info['provider']}")
        
        print("\nSwitching back to DeepSeek...")
        engine.switch_provider("deepseek")
        info = engine.get_provider_info()
        print(f"✅ Switched to {info['provider']}")
        
        return True
        
    except Exception as e:
        print(f"❌ Error with Multi-LLM Engine: {e}")
        return False


def setup_instructions():
    """Print setup instructions for using DeepSeek"""
    print("\n" + "="*60)
    print("📚 HOW TO USE DEEPSEEK WITH THIS SYSTEM")
    print("="*60)
    
    print("""
1. GET YOUR DEEPSEEK API KEY:
   - Sign up at https://platform.deepseek.com/
   - Go to API Keys section
   - Create a new API key
   
2. ADD TO YOUR .env FILE:
   ```
   DEEPSEEK_API_KEY=your_deepseek_api_key_here
   LLM_PROVIDER=deepseek  # Set DeepSeek as default
   DEEPSEEK_MODEL=deepseek-chat  # or deepseek-coder
   ```

3. INSTALL OPENAI LIBRARY (if not already installed):
   ```
   pip install openai
   ```

4. USE IN YOUR APPLICATION:

   Option A - Use DeepSeek Engine directly:
   ```python
   from backend.deepseek_engine import DeepSeekEngine
   engine = DeepSeekEngine()
   response = engine.generate_response("Your question here")
   ```

   Option B - Use Multi-LLM Engine (recommended):
   ```python
   from backend.multi_llm_rag_engine import MultiLLMRAGEngine
   
   # Use DeepSeek
   engine = MultiLLMRAGEngine(provider="deepseek")
   
   # Or set LLM_PROVIDER=deepseek in .env and just:
   engine = MultiLLMRAGEngine()
   ```

5. MODIFY app.py TO USE DEEPSEEK:
   Replace the line that initializes RAGEngine with:
   ```python
   from backend.multi_llm_rag_engine import MultiLLMRAGEngine
   rag_engine = MultiLLMRAGEngine()  # Will use LLM_PROVIDER from .env
   ```

DEEPSEEK MODELS:
- deepseek-chat: General purpose chat model (recommended)
- deepseek-coder: Optimized for code-related queries

PRICING:
DeepSeek is typically more cost-effective than Claude or GPT-4
Check current pricing at https://platform.deepseek.com/pricing
""")


def main():
    """Run all tests"""
    print("="*60)
    print("🚀 DEEPSEEK INTEGRATION TEST")
    print("="*60)
    
    # Check if OpenAI library is installed
    try:
        import openai
        print("✅ OpenAI library is installed")
    except ImportError:
        print("❌ OpenAI library not installed")
        print("   Run: pip install openai")
        return
    
    # Run tests
    has_api_key = os.getenv("DEEPSEEK_API_KEY") is not None
    
    if has_api_key:
        test_deepseek_direct()
        test_deepseek_rag_engine()
        test_multi_llm_engine()
    else:
        print("\n⚠️  DEEPSEEK_API_KEY not found in .env file")
        print("   Skipping API tests...")
    
    # Show setup instructions
    setup_instructions()
    
    print("\n" + "="*60)
    if has_api_key:
        print("✅ DeepSeek integration is ready to use!")
    else:
        print("📝 Follow the instructions above to set up DeepSeek")
    print("="*60)


if __name__ == "__main__":
    main()