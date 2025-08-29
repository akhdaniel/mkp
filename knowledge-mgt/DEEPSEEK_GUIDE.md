# DeepSeek Integration Guide

## Overview

Your Knowledge Management System now supports **DeepSeek** as an alternative to Anthropic Claude. DeepSeek offers cost-effective AI models with strong performance, especially for code-related queries.

## Quick Start

### 1. Get Your DeepSeek API Key

1. Sign up at [platform.deepseek.com](https://platform.deepseek.com/)
2. Navigate to API Keys section
3. Create a new API key
4. Copy the key

### 2. Configure Your Environment

Add your DeepSeek API key to the `.env` file:

```bash
# Edit your .env file
DEEPSEEK_API_KEY=your_actual_deepseek_key_here
LLM_PROVIDER=deepseek  # Set DeepSeek as default provider
DEEPSEEK_MODEL=deepseek-chat  # or deepseek-coder for code tasks
```

### 3. Test the Integration

Run the test script to verify everything works:

```bash
source venv/bin/activate
python deepseek_demo.py
```

## Available Models

### deepseek-chat
- **Purpose**: General-purpose conversations and Q&A
- **Best for**: Document search, helpdesk queries, general knowledge
- **Context**: 32K tokens

### deepseek-coder
- **Purpose**: Code generation and analysis
- **Best for**: Technical documentation, code-related queries
- **Context**: 16K tokens

## Usage Examples

### Option 1: Using DeepSeek Engine Directly

```python
from backend.deepseek_engine import DeepSeekEngine

# Initialize the engine
engine = DeepSeekEngine()

# Generate a response
response = engine.generate_response(
    query="What is the vacation policy?",
    context=retrieved_documents
)

print(response)
```

### Option 2: Using Multi-LLM Engine (Recommended)

```python
from backend.multi_llm_rag_engine import MultiLLMRAGEngine

# Initialize with DeepSeek
engine = MultiLLMRAGEngine(provider="deepseek")

# Or automatically use provider from .env
engine = MultiLLMRAGEngine()  # Uses LLM_PROVIDER setting

# Generate response with citations
result = engine.generate_response(
    query="How do I submit an expense report?",
    use_hybrid_search=True
)

print(f"Answer: {result['response']}")
print(f"Sources: {', '.join(result['citations'])}")
print(f"Model used: {result['model']}")
```

### Option 3: Switch Between Providers Dynamically

```python
from backend.multi_llm_rag_engine import MultiLLMRAGEngine

engine = MultiLLMRAGEngine()

# Use DeepSeek for cost-effective queries
engine.switch_provider("deepseek")
response1 = engine.generate_response("Simple question")

# Switch to Claude for complex reasoning
engine.switch_provider("anthropic")
response2 = engine.generate_response("Complex analysis question")
```

## Modifying the Streamlit App

To use DeepSeek in your main application, update `app.py`:

### Find this section:
```python
from backend.rag_engine import RAGEngine
# ...
rag_engine = RAGEngine()
```

### Replace with:
```python
from backend.multi_llm_rag_engine import MultiLLMRAGEngine
# ...
rag_engine = MultiLLMRAGEngine()  # Will use LLM_PROVIDER from .env
```

### Optional: Add Provider Selection to UI

```python
# In the sidebar
provider = st.selectbox(
    "Select AI Model",
    ["anthropic", "deepseek"],
    index=0 if os.getenv("LLM_PROVIDER") == "anthropic" else 1
)

if st.button("Switch Provider"):
    rag_engine.switch_provider(provider)
    st.success(f"Switched to {provider}")
```

## Cost Comparison

| Provider | Model | Input Cost | Output Cost | Context |
|----------|-------|------------|-------------|---------|
| DeepSeek | deepseek-chat | $0.14/1M tokens | $0.28/1M tokens | 32K |
| DeepSeek | deepseek-coder | $0.14/1M tokens | $0.28/1M tokens | 16K |
| Anthropic | Claude 3 Opus | $15/1M tokens | $75/1M tokens | 200K |
| OpenAI | GPT-4 | $30/1M tokens | $60/1M tokens | 128K |

**DeepSeek is ~100x cheaper than Claude 3 Opus!**

## Environment Variables Reference

```bash
# Required for DeepSeek
DEEPSEEK_API_KEY=your_key_here

# Optional DeepSeek Configuration
LLM_PROVIDER=deepseek           # or "anthropic" 
DEEPSEEK_MODEL=deepseek-chat    # or "deepseek-coder"

# Keep your Anthropic key for switching
ANTHROPIC_API_KEY=your_anthropic_key
```

## API Compatibility

DeepSeek uses an OpenAI-compatible API, which means:
- Uses the `openai` Python library
- Compatible with OpenAI SDK patterns
- Easy to switch from GPT models

## Troubleshooting

### Authentication Error
If you see "Authentication Fails", check:
1. Your API key is correctly set in `.env`
2. The key doesn't have extra spaces or quotes
3. You've reloaded the environment: `source venv/bin/activate`

### Module Not Found
If `openai` module is missing:
```bash
pip install openai
```

### Rate Limits
DeepSeek has generous rate limits:
- 500 requests per minute
- 100 concurrent requests

## Advanced Features

### Custom Temperature
```python
response = engine.generate_response(
    query="Your question",
    temperature=0.3  # Lower = more focused, Higher = more creative
)
```

### Streaming Responses
```python
# In deepseek_engine.py, set stream=True
response = client.chat.completions.create(
    model=self.model_name,
    messages=messages,
    stream=True  # Enable streaming
)

for chunk in response:
    print(chunk.choices[0].delta.content, end="")
```

## Best Practices

1. **Use DeepSeek for**:
   - High-volume queries (cost-effective)
   - Code-related questions (use deepseek-coder)
   - Standard helpdesk queries
   
2. **Use Claude for**:
   - Complex reasoning tasks
   - Creative writing
   - When you need larger context (200K tokens)

3. **Hybrid Approach**:
   - Implement query routing based on complexity
   - Use DeepSeek for initial filtering, Claude for detailed analysis

## Support

- DeepSeek Documentation: [platform.deepseek.com/docs](https://platform.deepseek.com/docs)
- API Status: [status.deepseek.com](https://status.deepseek.com)
- Pricing: [platform.deepseek.com/pricing](https://platform.deepseek.com/pricing)

---

Your system is now configured to use either Anthropic Claude or DeepSeek. Simply set `LLM_PROVIDER` in your `.env` file to switch between them!