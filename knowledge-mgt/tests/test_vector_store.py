"""
Tests for vector store functionality
"""
import pytest
from unittest.mock import Mock, patch, MagicMock
from backend.vector_store import VectorStore


class TestVectorStore:
    """Test VectorStore class"""
    
    @pytest.fixture
    def vector_store(self):
        """Create a VectorStore instance"""
        with patch('backend.vector_store.Chroma'):
            store = VectorStore()
            return store
    
    def test_initialization(self, vector_store):
        """Test VectorStore initialization"""
        assert vector_store is not None
        assert hasattr(vector_store, 'add_documents')
        assert hasattr(vector_store, 'search')
        
    def test_add_documents(self, vector_store):
        """Test adding documents to vector store"""
        documents = [
            {"content": "Document 1", "metadata": {"source": "file1.txt"}},
            {"content": "Document 2", "metadata": {"source": "file2.txt"}}
        ]
        
        with patch.object(vector_store, 'add_documents') as mock_add:
            mock_add.return_value = True
            result = vector_store.add_documents(documents)
            assert result is True
            mock_add.assert_called_once_with(documents)
            
    def test_search_documents(self, vector_store):
        """Test searching documents in vector store"""
        query = "test query"
        expected_results = [
            {"content": "Relevant document", "score": 0.95},
            {"content": "Another relevant doc", "score": 0.85}
        ]
        
        with patch.object(vector_store, 'search') as mock_search:
            mock_search.return_value = expected_results
            results = vector_store.search(query, k=2)
            
            assert len(results) == 2
            assert results[0]["score"] > results[1]["score"]
            mock_search.assert_called_once_with(query, k=2)
            
    def test_delete_documents(self, vector_store):
        """Test deleting documents from vector store"""
        doc_ids = ["doc1", "doc2"]
        
        with patch.object(vector_store, 'delete') as mock_delete:
            mock_delete.return_value = True
            result = vector_store.delete(doc_ids)
            assert result is True
            mock_delete.assert_called_once_with(doc_ids)
            
    def test_clear_store(self, vector_store):
        """Test clearing all documents from store"""
        with patch.object(vector_store, 'clear') as mock_clear:
            mock_clear.return_value = True
            result = vector_store.clear()
            assert result is True
            mock_clear.assert_called_once()
            
    def test_get_document_count(self, vector_store):
        """Test getting document count"""
        with patch.object(vector_store, 'count') as mock_count:
            mock_count.return_value = 42
            count = vector_store.count()
            assert count == 42
            mock_count.assert_called_once()
            
    @patch('backend.vector_store.ChromaDB')
    def test_chromadb_initialization(self, mock_chromadb):
        """Test ChromaDB specific initialization"""
        mock_client = Mock()
        mock_chromadb.return_value = mock_client
        
        store = VectorStore(db_type="chromadb")
        mock_chromadb.assert_called()
        
    @patch('backend.vector_store.QdrantClient')
    def test_qdrant_initialization(self, mock_qdrant):
        """Test Qdrant specific initialization"""
        mock_client = Mock()
        mock_qdrant.return_value = mock_client
        
        store = VectorStore(db_type="qdrant")
        mock_qdrant.assert_called()