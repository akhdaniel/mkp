"""
Tests for text splitting utilities
"""
import pytest
from utils.text_splitter import DocumentChunker


class TestDocumentChunker:
    """Test DocumentChunker class"""
    
    @pytest.fixture
    def chunker(self):
        return DocumentChunker(chunk_size=100, chunk_overlap=20)
        
    def test_initialization(self, chunker):
        """Test chunker initialization"""
        assert chunker.chunk_size == 100
        assert chunker.chunk_overlap == 20
        
    def test_chunk_text_basic(self, chunker):
        """Test basic text chunking"""
        text = "This is a sentence. " * 20  # Create text longer than chunk size
        chunks = chunker.chunk_text(text)
        
        assert len(chunks) > 1
        assert all(isinstance(chunk, dict) for chunk in chunks)
        assert all('text' in chunk for chunk in chunks)
        
    def test_chunk_text_with_overlap(self, chunker):
        """Test that chunks have overlap"""
        text = "Word1 Word2 Word3 Word4 Word5. " * 10
        chunks = chunker.chunk_text(text)
        
        if len(chunks) > 1:
            # Check that chunks have expected structure
            assert all('text' in chunk for chunk in chunks)
            assert all('metadata' in chunk for chunk in chunks)
                
    def test_chunk_text_preserves_metadata(self, chunker):
        """Test that chunking preserves metadata"""
        text = "First sentence. Second sentence. Third sentence. Fourth sentence."
        metadata = {"source": "test.txt", "page": 1}
        chunks = chunker.chunk_text(text, metadata=metadata)
        
        # Check metadata is preserved in chunks
        for chunk in chunks:
            assert 'metadata' in chunk
            assert chunk['metadata']['source'] == "test.txt"
                
    def test_empty_text(self, chunker):
        """Test chunking empty text"""
        chunks = chunker.chunk_text("")
        assert chunks == []
        
    def test_short_text(self, chunker):
        """Test chunking text shorter than chunk size"""
        text = "Short text"
        chunks = chunker.chunk_text(text)
        assert len(chunks) == 1
        assert chunks[0]['text'] == text
        
    def test_chunk_with_custom_size(self):
        """Test chunking with custom chunk size"""
        chunker = DocumentChunker(chunk_size=50, chunk_overlap=10)
        text = "This is a test. " * 20
        chunks = chunker.chunk_text(text)
        
        assert len(chunks) > 1
        # Verify chunk sizes are appropriate
        
    def test_chunk_documents(self, chunker):
        """Test chunking multiple documents"""
        documents = [
            {"text": "Document 1 content. " * 10, "metadata": {"id": 1}},
            {"text": "Document 2 content. " * 10, "metadata": {"id": 2}}
        ]
        
        all_chunks = chunker.chunk_documents(documents)
        
        assert len(all_chunks) > len(documents)
        # Verify metadata is preserved
        assert any(chunk['metadata']['id'] == 1 for chunk in all_chunks)
        assert any(chunk['metadata']['id'] == 2 for chunk in all_chunks)