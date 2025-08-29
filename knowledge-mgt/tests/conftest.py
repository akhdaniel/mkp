"""
Pytest configuration and fixtures for the test suite
"""
import os
import sys
import pytest
from pathlib import Path
from unittest.mock import Mock, patch

# Add project root to Python path
project_root = Path(__file__).parent.parent
sys.path.insert(0, str(project_root))

# Set test environment variables
os.environ['TESTING'] = 'true'
os.environ['ANTHROPIC_API_KEY'] = 'test-api-key'
os.environ['VECTOR_DB_PATH'] = './test_data/vectordb'
os.environ['UPLOAD_PATH'] = './test_data/uploads'
os.environ['CACHE_PATH'] = './test_data/cache'


@pytest.fixture(autouse=True)
def setup_test_environment():
    """Setup test environment before each test"""
    # Create test directories
    test_dirs = [
        './test_data/uploads',
        './test_data/vectordb', 
        './test_data/cache'
    ]
    for dir_path in test_dirs:
        Path(dir_path).mkdir(parents=True, exist_ok=True)
    
    yield
    
    # Cleanup after tests (optional)
    # Can add cleanup code here if needed


@pytest.fixture
def mock_anthropic_client():
    """Mock Anthropic client for testing"""
    with patch('anthropic.Anthropic') as mock:
        client = Mock()
        mock.return_value = client
        yield client


@pytest.fixture
def mock_openai_client():
    """Mock OpenAI client for testing"""
    with patch('openai.OpenAI') as mock:
        client = Mock()
        mock.return_value = client
        yield client


@pytest.fixture
def mock_chromadb_client():
    """Mock ChromaDB client for testing"""
    with patch('chromadb.Client') as mock:
        client = Mock()
        mock.return_value = client
        yield client


@pytest.fixture
def sample_documents():
    """Sample documents for testing"""
    return [
        {
            "content": "This is the first test document about company policies.",
            "metadata": {
                "source": "policy.pdf",
                "page": 1,
                "chunk_id": "doc1_chunk1"
            }
        },
        {
            "content": "This is the second test document about vacation rules.",
            "metadata": {
                "source": "handbook.pdf", 
                "page": 5,
                "chunk_id": "doc2_chunk1"
            }
        },
        {
            "content": "This is the third test document about expense reports.",
            "metadata": {
                "source": "finance.docx",
                "page": 2,
                "chunk_id": "doc3_chunk1"
            }
        }
    ]


@pytest.fixture
def sample_query():
    """Sample query for testing"""
    return "What is the company vacation policy?"


@pytest.fixture
def sample_context():
    """Sample context for RAG testing"""
    return [
        {
            "content": "Employees are entitled to 15 days of paid vacation per year.",
            "metadata": {"source": "policy.pdf", "page": 10}
        },
        {
            "content": "Vacation days must be approved by your manager in advance.",
            "metadata": {"source": "handbook.pdf", "page": 25}
        }
    ]


@pytest.fixture
def temp_test_file(tmp_path):
    """Create a temporary test file"""
    test_file = tmp_path / "test_document.txt"
    test_file.write_text("This is test content for document processing.")
    return str(test_file)


@pytest.fixture
def mock_streamlit():
    """Mock Streamlit for testing the UI"""
    with patch('streamlit.st') as mock_st:
        # Setup common streamlit attributes
        mock_st.session_state = {}
        mock_st.sidebar = Mock()
        mock_st.columns = Mock(return_value=[Mock(), Mock()])
        mock_st.file_uploader = Mock(return_value=None)
        mock_st.text_input = Mock(return_value="")
        mock_st.button = Mock(return_value=False)
        mock_st.success = Mock()
        mock_st.error = Mock()
        mock_st.info = Mock()
        mock_st.spinner = Mock()
        
        yield mock_st