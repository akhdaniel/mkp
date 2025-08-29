# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an AI-powered internal helpdesk system that enables employees to interact with organizational documents through a chat interface. It uses Retrieval-Augmented Generation (RAG) with vector database storage to provide accurate, context-aware responses from uploaded corporate documents.

## Development Commands

### Environment Setup
```bash
# Create and activate virtual environment
python -m venv venv
source venv/bin/activate  # macOS/Linux
# or
venv\Scripts\activate  # Windows

# Install dependencies
pip install -r requirements.txt

# Set up environment variables
cp .env.example .env
# Edit .env with your API keys
```

### Running the Application
```bash
# Run Streamlit app
streamlit run app.py

# Run with specific port
streamlit run app.py --server.port 8501

# Run in development mode with auto-reload
streamlit run app.py --server.runOnSave true
```

### Testing
```bash
# Run all tests
pytest

# Run specific test file
pytest tests/test_document_processor.py

# Run with coverage
pytest --cov=backend --cov-report=html

# Run specific test
pytest tests/test_rag_engine.py::test_generate_response
```

### Code Quality
```bash
# Format code with black
black backend/ utils/ tests/

# Check linting
pylint backend/ utils/

# Type checking
mypy backend/ utils/
```

## Architecture

### Core Components

1. **Document Processing Layer** (`backend/document_processor.py`)
   - Handles ingestion of PDF, DOCX, PPTX, TXT, XLSX formats
   - Implements recursive text splitting with configurable chunk size (default 512 tokens)
   - Maintains chunk overlap (default 50 tokens) for context preservation

2. **Vector Database Layer** (`backend/vector_store.py`)
   - Manages ChromaDB or Qdrant for vector storage
   - Handles document embeddings storage and retrieval
   - Implements hybrid search (semantic + keyword)

3. **RAG Engine** (`backend/rag_engine.py`)
   - Orchestrates retrieval and generation pipeline
   - Manages context window for Claude API
   - Implements reranking for relevance scoring
   - Handles top-k retrieval (default k=5)

4. **UI Layer** (`app.py`)
   - Streamlit-based interface
   - Session state management for conversation history
   - Document upload and management interface

### Data Flow

1. Document Upload → Text Extraction → Chunking → Embedding → Vector Storage
2. User Query → Query Embedding → Vector Search → Reranking → Context Assembly → LLM Generation → Response with Citations

### Key Design Patterns

- **Singleton Pattern**: Vector database connection management
- **Factory Pattern**: Document parser selection based on file type
- **Strategy Pattern**: Different embedding strategies (OpenAI vs sentence-transformers)
- **Chain of Responsibility**: Document processing pipeline

## Configuration

### Environment Variables
Required in `.env`:
```
ANTHROPIC_API_KEY=your_key
OPENAI_API_KEY=your_key  # Optional, for OpenAI embeddings
VECTOR_DB_PATH=./data/chromadb
UPLOAD_PATH=./data/uploads
MAX_FILE_SIZE=10485760
CHUNK_SIZE=512
CHUNK_OVERLAP=50
TOP_K_RETRIEVAL=5
MODEL_NAME=claude-3-opus-20240229
```

### Vector Database Selection
- Use ChromaDB for development (easier setup)
- Use Qdrant for production (better performance)
- Configuration in `backend/config.py`

## Project Structure

```
knowledge-mgt/
├── app.py                      # Main Streamlit application
├── backend/
│   ├── document_processor.py  # Document ingestion and chunking
│   ├── vector_store.py        # Vector database operations
│   ├── rag_engine.py          # RAG orchestration
│   ├── embeddings.py          # Embedding generation
│   └── config.py              # Configuration management
├── utils/
│   ├── text_splitter.py       # Text chunking utilities
│   ├── file_handlers.py       # File type specific parsers
│   └── validators.py          # Input validation
├── tests/                      # Test files
├── data/
│   ├── uploads/               # Uploaded documents
│   ├── vectordb/              # Vector database storage
│   └── cache/                 # Query cache
└── logs/                      # Application logs
```

## Key Implementation Notes

### Document Processing
- Use `langchain.text_splitter.RecursiveCharacterTextSplitter` for intelligent chunking
- Preserve document metadata throughout the pipeline
- Implement format-specific parsers in `utils/file_handlers.py`

### RAG Pipeline
- Implement hybrid search combining vector similarity and BM25 for better retrieval
- Use reranking with cross-encoder models for improved relevance
- Maintain conversation context across queries
- Always provide source citations with chunk references

### Performance Optimization
- Implement caching layer for frequent queries
- Use batch processing for document uploads
- Implement async operations where possible
- Use connection pooling for database connections

### Security
- Validate all file uploads (type, size, content)
- Sanitize user inputs before database queries
- Implement rate limiting for API calls
- Store API keys in environment variables only

## Dependencies Management

Core dependencies to maintain:
- `streamlit>=1.30.0` - UI framework
- `anthropic>=0.18.0` - Claude API
- `chromadb>=0.4.0` or `qdrant-client>=1.7.0` - Vector database
- `langchain>=0.1.0` - Document processing and RAG utilities
- `sentence-transformers>=2.2.0` - Embeddings (alternative to OpenAI)
- `pypdf2>=3.0.0` - PDF processing
- `python-docx>=1.0.0` - DOCX processing
- `python-pptx>=0.6.0` - PPTX processing

## Common Development Tasks

### Adding New Document Format Support
1. Create parser in `utils/file_handlers.py`
2. Register parser in `backend/document_processor.py`
3. Add format validation in `utils/validators.py`
4. Update tests in `tests/test_file_handlers.py`

### Modifying RAG Parameters
1. Update defaults in `backend/config.py`
2. Expose UI controls in `app.py` sidebar
3. Update environment variables documentation

### Implementing New Embedding Model
1. Add model class in `backend/embeddings.py`
2. Update factory method for model selection
3. Adjust chunk sizes if needed for model context

## Troubleshooting

### Common Issues
- **Import errors**: Ensure virtual environment is activated
- **API errors**: Check API keys in `.env` file
- **Vector DB connection**: Verify path permissions and disk space
- **Memory issues**: Reduce chunk size or implement pagination
- **Slow responses**: Check embedding cache, consider switching to lighter model