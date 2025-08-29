"""
Tests for embeddings functionality
"""
import pytest
import numpy as np
from unittest.mock import Mock, patch, MagicMock
from backend.embeddings import (
    EmbeddingModel,
    OpenAIEmbedding,
    SentenceTransformerEmbedding,
    EmbeddingFactory,
    EmbeddingManager
)


class TestEmbeddingGenerator:
    """Test base EmbeddingGenerator class"""
    
    def test_abstract_base_class(self):
        """Test that EmbeddingGenerator is abstract"""
        with pytest.raises(TypeError):
            EmbeddingGenerator()
            

class TestOpenAIEmbeddings:
    """Test OpenAI embeddings"""
    
    @pytest.fixture
    def embeddings(self):
        with patch('backend.embeddings.OpenAI'):
            return OpenAIEmbeddings(api_key="test-key")
            
    def test_initialization(self, embeddings):
        """Test OpenAI embeddings initialization"""
        assert embeddings is not None
        assert hasattr(embeddings, 'generate')
        assert hasattr(embeddings, 'generate_batch')
        
    @patch('backend.embeddings.OpenAI')
    def test_generate_single_embedding(self, mock_openai):
        """Test generating a single embedding"""
        mock_client = Mock()
        mock_response = Mock()
        mock_response.data = [Mock(embedding=[0.1, 0.2, 0.3, 0.4, 0.5])]
        mock_client.embeddings.create.return_value = mock_response
        mock_openai.return_value = mock_client
        
        embeddings = OpenAIEmbedding(api_key="test-key")
        result = embeddings.generate("test text")
        
        assert len(result) == 5
        assert result == [0.1, 0.2, 0.3, 0.4, 0.5]
        mock_client.embeddings.create.assert_called_once()
        
    @patch('backend.embeddings.OpenAI')
    def test_generate_batch_embeddings(self, mock_openai):
        """Test generating batch embeddings"""
        mock_client = Mock()
        mock_response = Mock()
        mock_response.data = [
            Mock(embedding=[0.1, 0.2, 0.3]),
            Mock(embedding=[0.4, 0.5, 0.6])
        ]
        mock_client.embeddings.create.return_value = mock_response
        mock_openai.return_value = mock_client
        
        embeddings = OpenAIEmbedding(api_key="test-key")
        texts = ["text 1", "text 2"]
        results = embeddings.generate_batch(texts)
        
        assert len(results) == 2
        assert len(results[0]) == 3
        assert len(results[1]) == 3
        
    @patch('backend.embeddings.OpenAI')
    def test_embedding_dimension(self, mock_openai):
        """Test embedding dimension property"""
        mock_client = Mock()
        mock_openai.return_value = mock_client
        
        embeddings = OpenAIEmbeddings(api_key="test-key", model="text-embedding-3-small")
        # Default dimension for text-embedding-3-small
        assert embeddings.dimension == 1536
        

class TestSentenceTransformerEmbeddings:
    """Test SentenceTransformer embeddings"""
    
    @pytest.fixture
    def embeddings(self):
        with patch('backend.embeddings.SentenceTransformer'):
            return SentenceTransformerEmbeddings(model_name="all-MiniLM-L6-v2")
            
    def test_initialization(self, embeddings):
        """Test SentenceTransformer embeddings initialization"""
        assert embeddings is not None
        assert hasattr(embeddings, 'generate')
        assert hasattr(embeddings, 'generate_batch')
        
    @patch('backend.embeddings.SentenceTransformer')
    def test_generate_single_embedding(self, mock_st):
        """Test generating a single embedding with SentenceTransformer"""
        mock_model = Mock()
        mock_model.encode.return_value = np.array([0.1, 0.2, 0.3, 0.4, 0.5])
        mock_st.return_value = mock_model
        
        embeddings = SentenceTransformerEmbedding(model_name="all-MiniLM-L6-v2")
        result = embeddings.generate("test text")
        
        assert len(result) == 5
        np.testing.assert_array_almost_equal(result, [0.1, 0.2, 0.3, 0.4, 0.5])
        mock_model.encode.assert_called_once()
        
    @patch('backend.embeddings.SentenceTransformer')
    def test_generate_batch_embeddings(self, mock_st):
        """Test generating batch embeddings with SentenceTransformer"""
        mock_model = Mock()
        mock_model.encode.return_value = np.array([
            [0.1, 0.2, 0.3],
            [0.4, 0.5, 0.6]
        ])
        mock_st.return_value = mock_model
        
        embeddings = SentenceTransformerEmbedding(model_name="all-MiniLM-L6-v2")
        texts = ["text 1", "text 2"]
        results = embeddings.generate_batch(texts)
        
        assert len(results) == 2
        assert len(results[0]) == 3
        assert len(results[1]) == 3
        
    @patch('backend.embeddings.SentenceTransformer')
    def test_embedding_dimension(self, mock_st):
        """Test embedding dimension property"""
        mock_model = Mock()
        mock_model.get_sentence_embedding_dimension.return_value = 384
        mock_st.return_value = mock_model
        
        embeddings = SentenceTransformerEmbedding(model_name="all-MiniLM-L6-v2")
        assert embeddings.dimension == 384
        

class TestGetEmbeddingModel:
    """Test embedding model factory function"""
    
    @patch('backend.embeddings.OpenAI')
    def test_get_openai_model(self, mock_openai):
        """Test getting OpenAI embedding model"""
        model = get_embedding_model(
            model_type="openai",
            api_key="test-key",
            model_name="text-embedding-3-small"
        )
        assert isinstance(model, OpenAIEmbeddings)
        
    @patch('backend.embeddings.SentenceTransformer')
    def test_get_sentence_transformer_model(self, mock_st):
        """Test getting SentenceTransformer model"""
        model = get_embedding_model(
            model_type="sentence-transformer",
            model_name="all-MiniLM-L6-v2"
        )
        assert isinstance(model, SentenceTransformerEmbeddings)
        
    def test_get_invalid_model_type(self):
        """Test getting invalid model type"""
        with pytest.raises(ValueError):
            get_embedding_model(model_type="invalid")
            
    @patch.dict('os.environ', {'OPENAI_API_KEY': 'env-test-key'})
    @patch('backend.embeddings.OpenAI')
    def test_get_model_with_env_api_key(self, mock_openai):
        """Test getting model with API key from environment"""
        model = get_embedding_model(model_type="openai")
        assert isinstance(model, OpenAIEmbeddings)
        
    @patch('backend.embeddings.SentenceTransformer')
    def test_default_model_names(self, mock_st):
        """Test default model names are used"""
        model = get_embedding_model(model_type="sentence-transformer")
        mock_st.assert_called_with("all-MiniLM-L6-v2")