# Knowledge Management System - Technical Documentation

## Table of Contents
1. [System Architecture](#system-architecture)
2. [Technology Stack](#technology-stack)
3. [Installation Guide](#installation-guide)
4. [Configuration](#configuration)
5. [API Documentation](#api-documentation)
6. [Database Schema](#database-schema)
7. [Security](#security)
8. [Performance Optimization](#performance-optimization)
9. [Monitoring & Logging](#monitoring--logging)
10. [Deployment](#deployment)
11. [Maintenance](#maintenance)
12. [Troubleshooting](#troubleshooting)

---

## 1. System Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         Users                                │
└─────────────────┬───────────────────────────────────────────┘
                  │ HTTPS
                  ▼
┌─────────────────────────────────────────────────────────────┐
│                   Streamlit Web Application                  │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                  UI Components                        │   │
│  │  • Chat Interface  • File Upload  • Settings        │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│                    Backend Services                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   Document   │  │     RAG      │  │   Vector     │     │
│  │  Processor   │  │    Engine    │  │    Store     │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│                    External Services                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │  Anthropic   │  │   DeepSeek   │  │   ChromaDB   │     │
│  │   Claude     │  │     API      │  │   Qdrant     │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
```

### Component Architecture

```python
# Core Components Structure
knowledge-mgt/
├── app.py                          # Main Streamlit application
├── backend/
│   ├── __init__.py
│   ├── config.py                   # Configuration management
│   ├── document_processor.py       # Document ingestion & chunking
│   ├── vector_store.py            # Vector database operations
│   ├── rag_engine.py              # RAG orchestration
│   ├── embeddings.py              # Embedding generation
│   ├── multi_llm_rag_engine.py   # Multi-LLM support
│   └── deepseek_engine.py        # DeepSeek integration
├── utils/
│   ├── __init__.py
│   ├── file_handlers.py          # File type handlers
│   ├── text_splitter.py          # Text chunking utilities
│   └── validators.py              # Input validation
└── data/
    ├── uploads/                   # Uploaded documents
    ├── vectordb/                  # Vector database storage
    └── cache/                     # Query cache
```

### Data Flow

1. **Document Ingestion Flow**
```
File Upload → Validation → Save to Disk → Text Extraction 
→ Chunking → Embedding Generation → Vector Storage
```

2. **Query Processing Flow**
```
User Query → Query Embedding → Vector Search → Context Retrieval 
→ Reranking → LLM Processing → Response Generation → Citation Formatting
```

3. **Session Management Flow**
```
Session Init → State Management → Message History 
→ Context Preservation → Export/Clear
```

---

## 2. Technology Stack

### Core Technologies

| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| **Framework** | Streamlit | 1.31.0 | Web UI framework |
| **Language** | Python | 3.11+ | Primary language |
| **Vector DB** | ChromaDB | 0.4.22 | Default vector storage |
| **Vector DB Alt** | Qdrant | 1.7.0 | Alternative vector storage |
| **LLM Primary** | Anthropic Claude | 0.64.0 | Primary AI model |
| **LLM Secondary** | DeepSeek | Latest | Cost-effective AI model |
| **Doc Processing** | LangChain | 0.3.27 | Document processing |
| **Embeddings** | Sentence-Transformers | 2.3.1 | Text embeddings |

### Dependencies

#### Production Dependencies
```python
# Core
streamlit==1.31.0
anthropic==0.64.0
openai==1.102.0  # For DeepSeek
langchain==0.3.27
langchain-community==0.3.29
langchain-anthropic==0.3.19

# Vector Storage
chromadb==0.4.22
qdrant-client==1.7.0  # Optional

# Document Processing
pypdf2==3.0.1
python-docx==1.1.0
python-pptx==0.6.23
pandas==2.2.3
openpyxl==3.1.2

# ML/AI
numpy==1.26.4
scikit-learn==1.5.2
sentence-transformers==2.3.1
tiktoken==0.5.2

# Utilities
python-dotenv==1.0.1
pydantic==2.11.7
pydantic-settings==2.10.1
watchdog==3.0.0
```

#### Development Dependencies
```python
# Testing
pytest==8.0.0
pytest-cov==4.1.0
pytest-asyncio==0.21.0

# Code Quality
black==24.1.1
pylint==3.0.3
mypy==1.8.0
flake8==6.1.0

# Documentation
mkdocs==1.5.3
mkdocs-material==9.5.0
```

---

## 3. Installation Guide

### Prerequisites

#### System Requirements
- **OS**: Ubuntu 20.04+, macOS 11+, Windows 10+
- **Python**: 3.11 or 3.12 (3.13 has compatibility issues)
- **RAM**: Minimum 4GB, Recommended 8GB
- **Storage**: 10GB free space
- **Network**: Stable internet connection

#### Software Requirements
- Git
- Python 3.11+
- pip or conda
- Virtual environment tool (venv/virtualenv)

### Installation Steps

#### 1. Clone Repository
```bash
git clone https://github.com/yourorg/knowledge-mgt.git
cd knowledge-mgt
```

#### 2. Create Virtual Environment
```bash
# Using venv
python -m venv venv

# Activate on Linux/Mac
source venv/bin/activate

# Activate on Windows
venv\Scripts\activate
```

#### 3. Install Dependencies
```bash
# Upgrade pip
pip install --upgrade pip

# Install requirements
pip install -r requirements.txt

# For development
pip install -r requirements-dev.txt
```

#### 4. Environment Configuration
```bash
# Copy example environment file
cp .env.example .env

# Edit .env with your configuration
nano .env  # or use any text editor
```

Required environment variables:
```env
# API Keys (at least one required)
ANTHROPIC_API_KEY=sk-ant-api03-xxxxx
DEEPSEEK_API_KEY=sk-xxxxx

# Optional
OPENAI_API_KEY=sk-xxxxx  # For OpenAI embeddings

# Configuration
LLM_PROVIDER=anthropic  # or deepseek
VECTOR_DB_TYPE=chromadb  # or qdrant
```

#### 5. Initialize Directories
```bash
# Create required directories
mkdir -p data/uploads data/vectordb data/cache logs

# Set permissions (Linux/Mac)
chmod 755 data/uploads data/vectordb data/cache logs
```

#### 6. Run Application
```bash
# Start the application
streamlit run app.py

# With specific port
streamlit run app.py --server.port 8501

# With custom configuration
streamlit run app.py --server.maxUploadSize 20
```

### Docker Installation

#### Dockerfile
```dockerfile
FROM python:3.11-slim

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application
COPY . .

# Create data directories
RUN mkdir -p data/uploads data/vectordb data/cache logs

# Expose port
EXPOSE 8501

# Health check
HEALTHCHECK CMD curl --fail http://localhost:8501/_stcore/health || exit 1

# Run application
CMD ["streamlit", "run", "app.py", "--server.port=8501", "--server.address=0.0.0.0"]
```

#### Docker Compose
```yaml
version: '3.8'

services:
  kms:
    build: .
    ports:
      - "8501:8501"
    environment:
      - ANTHROPIC_API_KEY=${ANTHROPIC_API_KEY}
      - DEEPSEEK_API_KEY=${DEEPSEEK_API_KEY}
      - LLM_PROVIDER=${LLM_PROVIDER:-anthropic}
    volumes:
      - ./data:/app/data
      - ./logs:/app/logs
    restart: unless-stopped

  # Optional: Qdrant vector database
  qdrant:
    image: qdrant/qdrant:latest
    ports:
      - "6333:6333"
    volumes:
      - ./qdrant_storage:/qdrant/storage
```

---

## 4. Configuration

### Configuration File Structure

#### backend/config.py
```python
from pydantic_settings import BaseSettings
from pydantic import Field
from typing import Optional
from pathlib import Path

class Settings(BaseSettings):
    # API Keys
    anthropic_api_key: Optional[str] = Field(default=None, env="ANTHROPIC_API_KEY")
    openai_api_key: Optional[str] = Field(default=None, env="OPENAI_API_KEY")
    deepseek_api_key: Optional[str] = Field(default=None, env="DEEPSEEK_API_KEY")
    
    # Database Configuration
    vector_db_path: str = Field(default="./data/vectordb", env="VECTOR_DB_PATH")
    vector_db_type: str = Field(default="chromadb", env="VECTOR_DB_TYPE")
    
    # File Upload Configuration
    upload_path: str = Field(default="./data/uploads", env="UPLOAD_PATH")
    max_file_size: int = Field(default=10485760, env="MAX_FILE_SIZE")  # 10MB
    
    # RAG Configuration
    chunk_size: int = Field(default=512, env="CHUNK_SIZE")
    chunk_overlap: int = Field(default=50, env="CHUNK_OVERLAP")
    top_k_retrieval: int = Field(default=5, env="TOP_K_RETRIEVAL")
    
    # Model Configuration
    model_name: str = Field(default="claude-3-opus-20240229", env="MODEL_NAME")
    embedding_model: str = Field(default="all-MiniLM-L6-v2", env="EMBEDDING_MODEL")
    llm_provider: str = Field(default="anthropic", env="LLM_PROVIDER")
    deepseek_model: str = Field(default="deepseek-chat", env="DEEPSEEK_MODEL")
    
    # Application Configuration
    log_level: str = Field(default="INFO", env="LOG_LEVEL")
    cache_path: str = Field(default="./data/cache", env="CACHE_PATH")
    session_timeout: int = Field(default=3600, env="SESSION_TIMEOUT")
    
    class Config:
        env_file = ".env"
        env_file_encoding = "utf-8"
```

### Environment Variables

#### Required Variables
| Variable | Description | Example |
|----------|-------------|---------|
| `ANTHROPIC_API_KEY` | Claude API key | sk-ant-api03-xxxxx |
| `DEEPSEEK_API_KEY` | DeepSeek API key | sk-xxxxx |

#### Optional Variables
| Variable | Default | Description |
|----------|---------|-------------|
| `LLM_PROVIDER` | anthropic | LLM provider (anthropic/deepseek) |
| `VECTOR_DB_TYPE` | chromadb | Vector database type |
| `CHUNK_SIZE` | 512 | Document chunk size in tokens |
| `CHUNK_OVERLAP` | 50 | Overlap between chunks |
| `TOP_K_RETRIEVAL` | 5 | Number of chunks to retrieve |
| `MAX_FILE_SIZE` | 10485760 | Maximum file size in bytes |
| `LOG_LEVEL` | INFO | Logging level |

### Model Configuration

#### Anthropic Claude Models
```python
CLAUDE_MODELS = {
    "claude-3-opus-20240229": {
        "context_window": 200000,
        "max_output": 4096,
        "cost_per_1k_input": 0.015,
        "cost_per_1k_output": 0.075
    },
    "claude-3-sonnet-20240229": {
        "context_window": 200000,
        "max_output": 4096,
        "cost_per_1k_input": 0.003,
        "cost_per_1k_output": 0.015
    }
}
```

#### DeepSeek Models
```python
DEEPSEEK_MODELS = {
    "deepseek-chat": {
        "context_window": 32000,
        "max_output": 4096,
        "cost_per_1m_input": 0.14,
        "cost_per_1m_output": 0.28
    },
    "deepseek-coder": {
        "context_window": 16000,
        "max_output": 4096,
        "cost_per_1m_input": 0.14,
        "cost_per_1m_output": 0.28
    }
}
```

---

## 5. API Documentation

### Internal APIs

#### Document Processor API

```python
class DocumentProcessor:
    def process_document(
        self,
        file_path: Path,
        metadata: Optional[Dict[str, Any]] = None
    ) -> Dict[str, Any]:
        """
        Process a single document.
        
        Args:
            file_path: Path to document file
            metadata: Optional metadata to attach
            
        Returns:
            {
                "id": "document_hash",
                "title": "Document Title",
                "content": "Full text content",
                "chunks": [
                    {
                        "content": "Chunk text",
                        "metadata": {
                            "chunk_index": 0,
                            "source": "file.pdf",
                            "page": 1
                        }
                    }
                ],
                "metadata": {...},
                "num_chunks": 10,
                "total_tokens": 5000
            }
        """
```

#### Vector Store API

```python
class VectorStore:
    def add_documents(
        self, 
        documents: List[Dict[str, Any]]
    ) -> Dict[str, Any]:
        """
        Add documents to vector store.
        
        Args:
            documents: List of processed documents
            
        Returns:
            {
                "total_documents": 1,
                "total_chunks": 10,
                "failed_chunks": 0,
                "success": true
            }
        """
    
    def search(
        self,
        query: str,
        top_k: int = 5,
        filter_metadata: Optional[Dict] = None
    ) -> List[Dict[str, Any]]:
        """
        Search for relevant documents.
        
        Args:
            query: Search query
            top_k: Number of results
            filter_metadata: Metadata filters
            
        Returns:
            [
                {
                    "content": "Relevant text",
                    "metadata": {...},
                    "score": 0.95
                }
            ]
        """
```

#### RAG Engine API

```python
class MultiLLMRAGEngine:
    def generate_response(
        self,
        query: str,
        conversation_history: Optional[List[Dict]] = None,
        filter_metadata: Optional[Dict] = None,
        use_hybrid_search: bool = True
    ) -> Dict[str, Any]:
        """
        Generate AI response.
        
        Args:
            query: User question
            conversation_history: Previous messages
            filter_metadata: Search filters
            use_hybrid_search: Use hybrid search
            
        Returns:
            {
                "response": "AI generated answer",
                "citations": ["doc1.pdf, p.5"],
                "documents_used": 3,
                "provider": "anthropic",
                "model": "claude-3-opus",
                "response_time": 1.5
            }
        """
```

### REST API Endpoints (Future)

#### Authentication
```http
POST /api/auth/login
Content-Type: application/json

{
    "username": "user@example.com",
    "password": "password123"
}

Response:
{
    "token": "jwt_token_here",
    "expires_in": 3600
}
```

#### Document Management
```http
# Upload document
POST /api/documents/upload
Authorization: Bearer {token}
Content-Type: multipart/form-data

file: binary_data
metadata: {"department": "HR"}

# List documents
GET /api/documents
Authorization: Bearer {token}

Response:
{
    "documents": [
        {
            "id": "doc123",
            "filename": "policy.pdf",
            "chunks": 10,
            "upload_date": "2024-01-01T12:00:00Z"
        }
    ],
    "total": 15
}

# Delete document
DELETE /api/documents/{document_id}
Authorization: Bearer {token}
```

#### Query Processing
```http
POST /api/query
Authorization: Bearer {token}
Content-Type: application/json

{
    "query": "What is the vacation policy?",
    "context_size": 5,
    "include_citations": true
}

Response:
{
    "response": "The vacation policy states...",
    "citations": [
        {
            "source": "HR_Policy.pdf",
            "page": 12,
            "relevance": 0.95
        }
    ],
    "processing_time": 1.2
}
```

---

## 6. Database Schema

### Vector Database Structure

#### ChromaDB Collections
```python
# Collection: documents
{
    "name": "documents",
    "metadata": {
        "created_at": "2024-01-01T00:00:00Z",
        "embedding_function": "all-MiniLM-L6-v2",
        "dimension": 384
    }
}

# Document Entry
{
    "id": "doc123_chunk_0",
    "embedding": [0.1, 0.2, ...],  # 384-dimensional vector
    "document": "This is the chunk text content...",
    "metadata": {
        "document_id": "doc123",
        "document_title": "Company Policy",
        "source": "/data/uploads/policy.pdf",
        "chunk_index": 0,
        "page": 1,
        "upload_date": "2024-01-01T12:00:00Z",
        "upload_user": "admin",
        "department": "HR"
    }
}
```

#### Qdrant Schema (Alternative)
```python
# Collection Configuration
{
    "collection_name": "documents",
    "vectors": {
        "size": 384,
        "distance": "Cosine"
    },
    "optimizers_config": {
        "default_segment_number": 2
    },
    "replication_factor": 2
}

# Point Structure
{
    "id": "uuid-v4",
    "vector": [0.1, 0.2, ...],
    "payload": {
        "document_id": "doc123",
        "content": "Chunk text",
        "metadata": {...}
    }
}
```

### File Storage Structure
```
data/
├── uploads/
│   ├── {timestamp}_{filename}
│   └── .metadata.json
├── vectordb/
│   ├── chroma.sqlite3
│   └── indexes/
└── cache/
    └── {query_hash}.json
```

### Metadata Schema
```json
{
    "document_id": "sha256_hash",
    "filename": "original_name.pdf",
    "file_path": "/data/uploads/123456_file.pdf",
    "file_size": 1048576,
    "file_type": "application/pdf",
    "upload_date": "2024-01-01T12:00:00Z",
    "upload_user": "admin",
    "department": "HR",
    "tags": ["policy", "vacation", "2024"],
    "processing": {
        "chunks": 10,
        "tokens": 5000,
        "processing_time": 2.5,
        "embedding_model": "all-MiniLM-L6-v2"
    }
}
```

---

## 7. Security

### Authentication & Authorization

#### Session Management
```python
# Session configuration
SESSION_CONFIG = {
    "timeout": 3600,  # 1 hour
    "refresh_token": True,
    "secure_cookie": True,
    "httponly": True,
    "samesite": "Strict"
}
```

#### API Key Management
- Store in environment variables
- Never commit to version control
- Rotate regularly (90 days)
- Use separate keys for dev/prod

### Data Security

#### Encryption
```python
# At-rest encryption
ENCRYPTION_CONFIG = {
    "algorithm": "AES-256-GCM",
    "key_derivation": "PBKDF2",
    "iterations": 100000
}

# In-transit encryption
TLS_CONFIG = {
    "min_version": "TLS 1.2",
    "ciphers": "ECDHE+AESGCM:ECDHE+CHACHA20"
}
```

#### Input Validation
```python
def validate_file_upload(file):
    # File type validation
    allowed_types = {'.pdf', '.docx', '.pptx', '.txt', '.xlsx'}
    if not file.suffix in allowed_types:
        raise ValueError("Invalid file type")
    
    # File size validation
    if file.size > MAX_FILE_SIZE:
        raise ValueError("File too large")
    
    # Content validation
    if has_malicious_content(file):
        raise ValueError("Malicious content detected")
    
    # Filename sanitization
    safe_name = sanitize_filename(file.name)
    return safe_name
```

### Privacy & Compliance

#### Data Retention
```python
RETENTION_POLICY = {
    "documents": 90,  # days
    "chat_history": 30,  # days
    "logs": 180,  # days
    "cache": 7  # days
}
```

#### GDPR Compliance
- Right to access: Export user data
- Right to deletion: Clear user documents
- Data minimization: Only collect necessary data
- Purpose limitation: Use only for helpdesk

#### Audit Logging
```python
def audit_log(action, user, details):
    log_entry = {
        "timestamp": datetime.now().isoformat(),
        "action": action,
        "user": user,
        "ip_address": request.remote_addr,
        "details": details
    }
    audit_logger.info(json.dumps(log_entry))
```

---

## 8. Performance Optimization

### Caching Strategy

#### Query Cache
```python
class QueryCache:
    def __init__(self, ttl=3600):
        self.cache = {}
        self.ttl = ttl
    
    def get_cached_response(self, query_hash):
        if query_hash in self.cache:
            entry = self.cache[query_hash]
            if time.time() - entry['timestamp'] < self.ttl:
                return entry['response']
        return None
    
    def cache_response(self, query_hash, response):
        self.cache[query_hash] = {
            'response': response,
            'timestamp': time.time()
        }
```

#### Embedding Cache
```python
EMBEDDING_CACHE = {
    "max_size": 10000,  # Maximum cached embeddings
    "ttl": 86400,  # 24 hours
    "eviction": "LRU"  # Least Recently Used
}
```

### Database Optimization

#### ChromaDB Settings
```python
CHROMA_CONFIG = {
    "persist_directory": "./data/vectordb",
    "anonymized_telemetry": False,
    "cache_size": 1000,
    "batch_size": 100
}
```

#### Indexing Strategy
```python
def optimize_indexes():
    # Create indexes for common queries
    collection.create_index(
        field="metadata.department",
        index_type="hash"
    )
    collection.create_index(
        field="metadata.upload_date",
        index_type="btree"
    )
```

### Resource Management

#### Connection Pooling
```python
CONNECTION_POOL = {
    "min_connections": 2,
    "max_connections": 10,
    "connection_timeout": 30,
    "idle_timeout": 300
}
```

#### Memory Management
```python
# Chunk processing in batches
BATCH_CONFIG = {
    "chunk_batch_size": 50,
    "embedding_batch_size": 100,
    "max_memory_usage": "2GB"
}
```

### Performance Metrics

#### Key Metrics to Monitor
```python
PERFORMANCE_METRICS = {
    "response_time_p50": 1.0,  # seconds
    "response_time_p95": 3.0,
    "response_time_p99": 5.0,
    "throughput": 100,  # requests/minute
    "error_rate": 0.01,  # 1%
    "cache_hit_rate": 0.7  # 70%
}
```

---

## 9. Monitoring & Logging

### Logging Configuration

#### Logger Setup
```python
import logging
from logging.handlers import RotatingFileHandler

def setup_logging():
    logger = logging.getLogger('kms')
    logger.setLevel(logging.INFO)
    
    # File handler
    file_handler = RotatingFileHandler(
        'logs/kms.log',
        maxBytes=10485760,  # 10MB
        backupCount=10
    )
    file_handler.setFormatter(
        logging.Formatter(
            '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
        )
    )
    
    # Console handler
    console_handler = logging.StreamHandler()
    console_handler.setLevel(logging.WARNING)
    
    logger.addHandler(file_handler)
    logger.addHandler(console_handler)
    
    return logger
```

#### Log Levels
```python
LOG_LEVELS = {
    "DEBUG": "Detailed debugging information",
    "INFO": "General informational messages",
    "WARNING": "Warning messages for potential issues",
    "ERROR": "Error messages for failures",
    "CRITICAL": "Critical failures requiring immediate attention"
}
```

### Metrics Collection

#### Application Metrics
```python
from prometheus_client import Counter, Histogram, Gauge

# Define metrics
request_count = Counter('kms_requests_total', 'Total requests')
request_duration = Histogram('kms_request_duration_seconds', 'Request duration')
active_users = Gauge('kms_active_users', 'Active users')
document_count = Gauge('kms_documents_total', 'Total documents')

# Collect metrics
@track_metrics
def process_query(query):
    with request_duration.time():
        request_count.inc()
        result = rag_engine.generate_response(query)
    return result
```

#### System Metrics
```python
SYSTEM_METRICS = {
    "cpu_usage": psutil.cpu_percent(),
    "memory_usage": psutil.virtual_memory().percent,
    "disk_usage": psutil.disk_usage('/').percent,
    "network_io": psutil.net_io_counters(),
    "open_connections": len(psutil.net_connections())
}
```

### Health Checks

#### Endpoint Health
```python
@app.route('/health')
def health_check():
    checks = {
        "status": "healthy",
        "timestamp": datetime.now().isoformat(),
        "checks": {
            "database": check_database(),
            "vector_store": check_vector_store(),
            "llm_api": check_llm_api(),
            "storage": check_storage()
        }
    }
    
    # Determine overall health
    if all(check["status"] == "ok" for check in checks["checks"].values()):
        return jsonify(checks), 200
    else:
        checks["status"] = "degraded"
        return jsonify(checks), 503
```

### Alerting

#### Alert Configuration
```python
ALERT_RULES = [
    {
        "name": "High Error Rate",
        "condition": "error_rate > 0.05",
        "severity": "critical",
        "notification": ["email", "slack"]
    },
    {
        "name": "Slow Response Time",
        "condition": "response_time_p95 > 5",
        "severity": "warning",
        "notification": ["email"]
    },
    {
        "name": "Low Disk Space",
        "condition": "disk_usage > 90",
        "severity": "critical",
        "notification": ["email", "pagerduty"]
    }
]
```

---

## 10. Deployment

### Production Deployment

#### Server Requirements
```yaml
# Minimum Requirements
CPU: 4 cores
RAM: 8 GB
Storage: 50 GB SSD
Network: 100 Mbps

# Recommended Requirements
CPU: 8 cores
RAM: 16 GB
Storage: 100 GB SSD
Network: 1 Gbps
```

#### Deployment Checklist
- [ ] Environment variables configured
- [ ] SSL certificates installed
- [ ] Firewall rules configured
- [ ] Backup strategy implemented
- [ ] Monitoring enabled
- [ ] Log rotation configured
- [ ] Health checks passing
- [ ] Load testing completed
- [ ] Security scan passed
- [ ] Documentation updated

### Cloud Deployment

#### AWS Deployment
```yaml
# EC2 Instance
Type: t3.large
AMI: Ubuntu 22.04 LTS
Security Group:
  - Port 443 (HTTPS)
  - Port 22 (SSH, restricted)

# RDS (if using PostgreSQL)
Engine: PostgreSQL 14
Instance: db.t3.medium
Storage: 100 GB SSD
Backup: 7 days retention

# S3 (for document storage)
Bucket: kms-documents
Versioning: Enabled
Encryption: AES-256
```

#### Docker Deployment
```bash
# Build image
docker build -t kms:latest .

# Run container
docker run -d \
  --name kms \
  -p 8501:8501 \
  -v ./data:/app/data \
  -v ./logs:/app/logs \
  --env-file .env \
  --restart unless-stopped \
  kms:latest
```

#### Kubernetes Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kms
spec:
  replicas: 3
  selector:
    matchLabels:
      app: kms
  template:
    metadata:
      labels:
        app: kms
    spec:
      containers:
      - name: kms
        image: kms:latest
        ports:
        - containerPort: 8501
        env:
        - name: ANTHROPIC_API_KEY
          valueFrom:
            secretKeyRef:
              name: kms-secrets
              key: anthropic-api-key
        resources:
          requests:
            memory: "2Gi"
            cpu: "1"
          limits:
            memory: "4Gi"
            cpu: "2"
        volumeMounts:
        - name: data
          mountPath: /app/data
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: kms-data-pvc
```

### CI/CD Pipeline

#### GitHub Actions
```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Set up Python
      uses: actions/setup-python@v2
      with:
        python-version: '3.11'
    - name: Install dependencies
      run: |
        pip install -r requirements.txt
        pip install -r requirements-dev.txt
    - name: Run tests
      run: pytest --cov=backend --cov-report=xml
    - name: Run linting
      run: |
        black --check .
        pylint backend/
        mypy backend/
    
  deploy:
    needs: test
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
    - name: Deploy to production
      run: |
        # Deployment script
        ssh ${{ secrets.PROD_SERVER }} 'cd /app && git pull && docker-compose up -d'
```

---

## 11. Maintenance

### Routine Maintenance

#### Daily Tasks
- [ ] Check system health endpoints
- [ ] Review error logs
- [ ] Monitor disk usage
- [ ] Verify backup completion

#### Weekly Tasks
- [ ] Review performance metrics
- [ ] Clean up old cache files
- [ ] Update document indexes
- [ ] Check for security updates

#### Monthly Tasks
- [ ] Rotate API keys
- [ ] Update dependencies
- [ ] Performance optimization review
- [ ] Security audit
- [ ] Backup restoration test

### Database Maintenance

#### Vector Store Optimization
```python
def optimize_vector_store():
    """Monthly vector store optimization"""
    # Compact database
    vector_store.compact()
    
    # Rebuild indexes
    vector_store.rebuild_indexes()
    
    # Remove orphaned chunks
    vector_store.cleanup_orphans()
    
    # Update statistics
    vector_store.update_statistics()
```

#### Backup Procedures
```bash
#!/bin/bash
# Backup script

# Set variables
BACKUP_DIR="/backup/kms"
DATE=$(date +%Y%m%d_%H%M%S)

# Create backup directory
mkdir -p $BACKUP_DIR/$DATE

# Backup database
cp -r data/vectordb $BACKUP_DIR/$DATE/

# Backup uploads
tar -czf $BACKUP_DIR/$DATE/uploads.tar.gz data/uploads/

# Backup configuration
cp .env $BACKUP_DIR/$DATE/

# Remove old backups (keep 30 days)
find $BACKUP_DIR -type d -mtime +30 -exec rm -rf {} \;

# Upload to S3 (optional)
aws s3 sync $BACKUP_DIR/$DATE s3://backup-bucket/kms/$DATE
```

### Update Procedures

#### Dependency Updates
```bash
# Check for updates
pip list --outdated

# Update specific package
pip install --upgrade package_name

# Test after update
pytest tests/

# Update requirements.txt
pip freeze > requirements.txt
```

#### System Updates
```python
def update_system():
    """System update procedure"""
    steps = [
        "1. Announce maintenance window",
        "2. Backup current system",
        "3. Pull latest code",
        "4. Update dependencies",
        "5. Run database migrations",
        "6. Run tests",
        "7. Deploy to staging",
        "8. Verify staging",
        "9. Deploy to production",
        "10. Verify production",
        "11. Monitor for issues"
    ]
    return steps
```

---

## 12. Troubleshooting

### Common Issues

#### Issue: High Memory Usage
```python
# Diagnosis
def diagnose_memory():
    import tracemalloc
    tracemalloc.start()
    
    # Run operations
    snapshot = tracemalloc.take_snapshot()
    top_stats = snapshot.statistics('lineno')
    
    for stat in top_stats[:10]:
        print(stat)

# Solution
MEMORY_FIXES = {
    "reduce_batch_size": 50,
    "enable_gc": True,
    "clear_cache_regularly": True,
    "limit_concurrent_requests": 10
}
```

#### Issue: Slow Query Response
```python
# Diagnosis
def analyze_slow_query(query):
    profile = {
        "embedding_time": 0,
        "search_time": 0,
        "rerank_time": 0,
        "llm_time": 0
    }
    
    # Profile each step
    start = time.time()
    embedding = generate_embedding(query)
    profile["embedding_time"] = time.time() - start
    
    # Continue for each step...
    return profile

# Solutions
PERFORMANCE_TUNING = {
    "enable_caching": True,
    "reduce_chunk_size": 256,
    "limit_search_results": 3,
    "use_simpler_model": "claude-3-haiku"
}
```

#### Issue: Document Processing Fails
```python
# Common causes and fixes
PROCESSING_FIXES = {
    "corrupted_file": "Re-save or recreate document",
    "unsupported_encoding": "Convert to UTF-8",
    "file_too_large": "Split into smaller files",
    "complex_layout": "Convert to simpler format",
    "password_protected": "Remove password protection"
}
```

### Debug Mode

#### Enable Debug Logging
```python
# In .env
LOG_LEVEL=DEBUG
STREAMLIT_SERVER_ENABLE_CORS=false
STREAMLIT_SERVER_ENABLE_XSRF_PROTECTION=false

# In code
import streamlit as st
st.set_option('deprecation.showPyplotGlobalUse', False)
st.set_option('server.enableCORS', False)
```

#### Debug Tools
```python
def debug_context():
    """Print debug information"""
    import pdb
    pdb.set_trace()
    
    # Or use IPython
    from IPython import embed
    embed()
    
    # Or use VS Code debugger
    import debugpy
    debugpy.listen(5678)
    debugpy.wait_for_client()
```

### Recovery Procedures

#### Database Recovery
```bash
#!/bin/bash
# Restore from backup

# Stop application
systemctl stop kms

# Backup current (corrupted) data
mv data/vectordb data/vectordb.corrupted

# Restore from backup
cp -r /backup/kms/20240101/vectordb data/

# Verify restoration
python -c "from backend.vector_store import VectorStore; vs = VectorStore(); print(vs.get_collection_stats())"

# Restart application
systemctl start kms
```

#### Emergency Rollback
```bash
#!/bin/bash
# Quick rollback procedure

# Save current version
git tag rollback-point

# Rollback to previous version
git checkout previous-release-tag

# Restore dependencies
pip install -r requirements.txt

# Restart services
docker-compose down
docker-compose up -d

# Verify
curl http://localhost:8501/health
```

---

## Appendices

### A. API Response Codes

| Code | Status | Description |
|------|--------|-------------|
| 200 | OK | Request successful |
| 201 | Created | Resource created |
| 400 | Bad Request | Invalid input |
| 401 | Unauthorized | Authentication required |
| 403 | Forbidden | Access denied |
| 404 | Not Found | Resource not found |
| 413 | Payload Too Large | File too large |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Server Error | Server error |
| 503 | Service Unavailable | Service down |

### B. File Format Specifications

| Format | MIME Type | Max Size | Parser |
|--------|-----------|----------|--------|
| PDF | application/pdf | 10 MB | PyPDF2 |
| DOCX | application/vnd.openxmlformats | 10 MB | python-docx |
| PPTX | application/vnd.openxmlformats | 10 MB | python-pptx |
| XLSX | application/vnd.openxmlformats | 10 MB | openpyxl |
| TXT | text/plain | 10 MB | Built-in |

### C. Environment Template

```env
# Production Environment Template
# Copy to .env and fill in values

# === REQUIRED ===
ANTHROPIC_API_KEY=
DEEPSEEK_API_KEY=

# === OPTIONAL ===
OPENAI_API_KEY=

# === CONFIGURATION ===
LLM_PROVIDER=anthropic
VECTOR_DB_TYPE=chromadb
VECTOR_DB_PATH=./data/vectordb
UPLOAD_PATH=./data/uploads
CACHE_PATH=./data/cache
LOG_LEVEL=INFO

# === LIMITS ===
MAX_FILE_SIZE=10485760
CHUNK_SIZE=512
CHUNK_OVERLAP=50
TOP_K_RETRIEVAL=5
SESSION_TIMEOUT=3600

# === MODELS ===
MODEL_NAME=claude-3-opus-20240229
DEEPSEEK_MODEL=deepseek-chat
EMBEDDING_MODEL=all-MiniLM-L6-v2
```

---

*End of Technical Documentation - Version 1.0*
*Last Updated: December 2024*
*© 2024 Your Organization - Confidential*