"""
Tests for document processing functionality
"""
import os
import pytest
from pathlib import Path
from unittest.mock import Mock, patch, MagicMock
from backend.document_processor import DocumentProcessor


class TestDocumentProcessor:
    """Test DocumentProcessor class"""
    
    @pytest.fixture
    def processor(self):
        """Create a DocumentProcessor instance"""
        return DocumentProcessor()
    
    def test_initialization(self, processor):
        """Test DocumentProcessor initialization"""
        assert processor is not None
        assert hasattr(processor, 'process_document')
        
    def test_supported_formats(self, processor):
        """Test that processor supports expected file formats"""
        supported = {'.pdf', '.docx', '.pptx', '.txt', '.xlsx'}
        assert hasattr(processor, 'supported_formats')
        # Check if method or property exists to get supported formats
        
    @patch('backend.document_processor.os.path.exists')
    def test_process_text_file(self, mock_exists, processor):
        """Test processing a text file"""
        mock_exists.return_value = True
        test_file = "test.txt"
        test_content = "This is test content for document processing."
        
        with patch('builtins.open', mock_open(read_data=test_content)):
            result = processor.process_document(test_file)
            assert result is not None
            # Verify chunks were created
            
    def test_chunk_text(self, processor):
        """Test text chunking functionality"""
        text = "This is a long text. " * 100  # Create long text
        chunks = processor.chunk_text(text, chunk_size=50, overlap=10)
        
        assert len(chunks) > 1
        assert all(isinstance(chunk, str) for chunk in chunks)
        
    def test_process_nonexistent_file(self, processor):
        """Test handling of non-existent file"""
        with pytest.raises(FileNotFoundError):
            processor.process_document("nonexistent.txt")
            
    def test_process_unsupported_format(self, processor):
        """Test handling of unsupported file format"""
        with pytest.raises(ValueError):
            processor.process_document("file.xyz")
            
    @patch('backend.document_processor.PyPDFLoader')
    def test_process_pdf(self, mock_loader, processor):
        """Test PDF processing"""
        mock_loader_instance = Mock()
        mock_loader_instance.load.return_value = [
            Mock(page_content="Page 1 content"),
            Mock(page_content="Page 2 content")
        ]
        mock_loader.return_value = mock_loader_instance
        
        result = processor.process_document("test.pdf")
        assert result is not None
        mock_loader.assert_called_once()


def mock_open(read_data=""):
    """Helper to create mock file object"""
    import builtins
    from unittest.mock import mock_open as _mock_open
    return _mock_open(read_data=read_data)