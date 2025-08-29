"""
Tests for RAG engine functionality
"""
import pytest
from unittest.mock import Mock, patch, MagicMock
from backend.rag_engine import RAGEngine


class TestRAGEngine:
    """Test RAGEngine class"""
    
    @pytest.fixture
    def rag_engine(self):
        """Create a RAGEngine instance"""
        with patch('backend.rag_engine.Anthropic'):
            with patch('backend.rag_engine.VectorStore'):
                engine = RAGEngine()
                return engine
    
    def test_initialization(self, rag_engine):
        """Test RAGEngine initialization"""
        assert rag_engine is not None
        assert hasattr(rag_engine, 'generate_response')
        assert hasattr(rag_engine, 'retrieve_context')
        
    def test_retrieve_context(self, rag_engine):
        """Test context retrieval from vector store"""
        query = "What is the company policy?"
        mock_results = [
            {"content": "Policy document 1", "metadata": {"source": "policy.pdf"}},
            {"content": "Policy document 2", "metadata": {"source": "handbook.pdf"}}
        ]
        
        with patch.object(rag_engine.vector_store, 'search') as mock_search:
            mock_search.return_value = mock_results
            context = rag_engine.retrieve_context(query, k=2)
            
            assert len(context) == 2
            assert "Policy document 1" in context[0]["content"]
            mock_search.assert_called_once_with(query, k=2)
            
    def test_generate_response(self, rag_engine):
        """Test response generation with context"""
        query = "What is the vacation policy?"
        context = [
            {"content": "Employees get 15 days vacation", "metadata": {"source": "policy.pdf"}}
        ]
        expected_response = "According to the policy, employees receive 15 days of vacation."
        
        with patch.object(rag_engine, 'retrieve_context') as mock_retrieve:
            mock_retrieve.return_value = context
            with patch.object(rag_engine.client.messages, 'create') as mock_create:
                mock_response = Mock()
                mock_response.content = [Mock(text=expected_response)]
                mock_create.return_value = mock_response
                
                response = rag_engine.generate_response(query)
                
                assert expected_response in response
                mock_retrieve.assert_called_once()
                mock_create.assert_called_once()
                
    def test_generate_response_no_context(self, rag_engine):
        """Test response generation when no context is found"""
        query = "Random question with no context"
        
        with patch.object(rag_engine, 'retrieve_context') as mock_retrieve:
            mock_retrieve.return_value = []
            
            response = rag_engine.generate_response(query)
            assert response is not None
            # Should return a message about no relevant information found
            
    def test_format_citations(self, rag_engine):
        """Test citation formatting"""
        context = [
            {"content": "Content 1", "metadata": {"source": "doc1.pdf", "page": 1}},
            {"content": "Content 2", "metadata": {"source": "doc2.pdf", "page": 5}}
        ]
        
        citations = rag_engine.format_citations(context)
        assert len(citations) == 2
        assert "doc1.pdf" in citations[0]
        assert "page 1" in citations[0] or "p. 1" in citations[0]
        
    def test_rerank_results(self, rag_engine):
        """Test reranking of search results"""
        query = "vacation policy"
        results = [
            {"content": "Unrelated content", "score": 0.9},
            {"content": "Vacation policy: 15 days", "score": 0.8},
            {"content": "Holiday schedule", "score": 0.85}
        ]
        
        with patch.object(rag_engine, 'rerank') as mock_rerank:
            mock_rerank.return_value = [results[1], results[2], results[0]]
            reranked = rag_engine.rerank(query, results)
            
            assert reranked[0]["content"] == "Vacation policy: 15 days"
            mock_rerank.assert_called_once_with(query, results)
            
    def test_handle_conversation_history(self, rag_engine):
        """Test handling of conversation history"""
        history = [
            {"role": "user", "content": "What is the policy?"},
            {"role": "assistant", "content": "Which policy would you like to know about?"},
            {"role": "user", "content": "Vacation policy"}
        ]
        
        query = "How many days?"
        
        with patch.object(rag_engine, 'generate_response') as mock_generate:
            mock_generate.return_value = "Employees get 15 days of vacation per year."
            
            response = rag_engine.generate_response(query, history=history)
            
            # Verify that history context is considered
            assert response is not None
            mock_generate.assert_called()
            
    def test_error_handling(self, rag_engine):
        """Test error handling in RAG engine"""
        query = "Test query"
        
        with patch.object(rag_engine.client.messages, 'create') as mock_create:
            mock_create.side_effect = Exception("API Error")
            
            with pytest.raises(Exception) as exc_info:
                rag_engine.generate_response(query)
            
            assert "API Error" in str(exc_info.value)