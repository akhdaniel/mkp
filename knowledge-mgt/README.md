# AI Helpdesk Chat System

An intelligent document chat application that enables employees to interact with organizational documents through an AI-powered interface using Retrieval-Augmented Generation (RAG).

## Features

- **Document Processing**: Support for PDF, DOCX, PPTX, TXT, and XLSX files
- **Intelligent Chunking**: Smart document splitting with context preservation
- **Hybrid Search**: Combines semantic and keyword search for better results
- **Conversational AI**: Claude-powered responses with context awareness
- **Source Citations**: All responses include document sources
- **Confidence Scoring**: Transparency in response reliability
- **Session Management**: Persistent chat history within sessions

## Quick Start

### Prerequisites

- Python 3.10+
- Anthropic API key (required)
- OpenAI API key (optional, for OpenAI embeddings)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd knowledge-mgt
```

2. Create a virtual environment:
```bash
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
```

3. Install dependencies:
```bash
pip install -r requirements.txt
```

4. Set up environment variables:
```bash
cp .env.example .env
# Edit .env and add your API keys
```

### Running the Application

```bash
streamlit run app.py
```

The application will open in your browser at `http://localhost:8501`

## Usage

1. **Upload Documents**: Use the sidebar to upload your documents (PDF, DOCX, etc.)
2. **Ask Questions**: Type questions in the chat interface
3. **View Sources**: Expand the sources section to see which documents were used
4. **Export Chat**: Download conversation history as markdown

## Configuration

Key settings in `.env`:

- `ANTHROPIC_API_KEY`: Your Anthropic API key (required)
- `CHUNK_SIZE`: Document chunk size (default: 512 tokens)
- `TOP_K_RETRIEVAL`: Number of chunks to retrieve (default: 5)
- `MODEL_NAME`: Claude model to use (default: claude-3-opus-20240229)

## Project Structure

```
knowledge-mgt/
├── app.py                  # Main Streamlit application
├── backend/
│   ├── config.py          # Configuration management
│   ├── document_processor.py  # Document processing
│   ├── embeddings.py      # Embedding generation
│   ├── vector_store.py    # Vector database operations
│   └── rag_engine.py      # RAG orchestration
├── utils/
│   ├── file_handlers.py  # File type handlers
│   └── text_splitter.py  # Text chunking utilities
├── data/
│   ├── uploads/          # Uploaded documents
│   ├── vectordb/         # Vector database storage
│   └── cache/            # Query cache
└── requirements.txt      # Python dependencies
```

## Development

### Running Tests
```bash
pytest
```

### Code Formatting
```bash
black backend/ utils/
```

### Type Checking
```bash
mypy backend/ utils/
```

## Troubleshooting

**Import errors**: Ensure virtual environment is activated

**API errors**: Check API keys in `.env` file

**Memory issues**: Reduce `CHUNK_SIZE` or implement pagination

**Slow responses**: Check embedding cache, consider using a lighter model

## License

This project is for internal use only.

## Support

For issues or questions, please contact the IT helpdesk team.