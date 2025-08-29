"""
Tests for file handlers
"""
import pytest
from pathlib import Path
from unittest.mock import Mock, patch, mock_open
from utils.file_handlers import (
    PDFHandler, 
    DocxHandler, 
    PptxHandler, 
    TxtHandler, 
    XlsxHandler,
    FileHandlerFactory
)


class TestFileHandlers:
    """Test file handler utilities"""
    
    def test_get_file_handler_pdf(self):
        """Test getting PDF handler"""
        factory = FileHandlerFactory()
        handler = factory.get_handler(".pdf")
        assert isinstance(handler, PDFHandler)
        
    def test_get_file_handler_docx(self):
        """Test getting DOCX handler"""
        factory = FileHandlerFactory()
        handler = factory.get_handler(".docx")
        assert isinstance(handler, DocxHandler)
        
    def test_get_file_handler_pptx(self):
        """Test getting PPTX handler"""
        factory = FileHandlerFactory()
        handler = factory.get_handler(".pptx")
        assert isinstance(handler, PptxHandler)
        
    def test_get_file_handler_txt(self):
        """Test getting TXT handler"""
        factory = FileHandlerFactory()
        handler = factory.get_handler(".txt")
        assert isinstance(handler, TxtHandler)
        
    def test_get_file_handler_xlsx(self):
        """Test getting XLSX handler"""
        factory = FileHandlerFactory()
        handler = factory.get_handler(".xlsx")
        assert isinstance(handler, XlsxHandler)
        
    def test_get_file_handler_unsupported(self):
        """Test unsupported file type"""
        factory = FileHandlerFactory()
        with pytest.raises(ValueError):
            factory.get_handler(".xyz")
            

class TestPDFHandler:
    """Test PDF file handler"""
    
    @pytest.fixture
    def handler(self):
        return PDFHandler()
        
    @patch('utils.file_handlers.PyPDF2.PdfReader')
    def test_extract_text(self, mock_reader, handler):
        """Test text extraction from PDF"""
        mock_pdf = Mock()
        mock_page = Mock()
        mock_page.extract_text.return_value = "Page content"
        mock_pdf.pages = [mock_page]
        mock_reader.return_value = mock_pdf
        
        with patch('builtins.open', mock_open()):
            text = handler.extract_text("test.pdf")
            
        assert "Page content" in text
        mock_reader.assert_called_once()
        

class TestDocxHandler:
    """Test DOCX file handler"""
    
    @pytest.fixture
    def handler(self):
        return DocxHandler()
        
    @patch('utils.file_handlers.Document')
    def test_extract_text(self, mock_document, handler):
        """Test text extraction from DOCX"""
        mock_doc = Mock()
        mock_paragraph = Mock()
        mock_paragraph.text = "Paragraph content"
        mock_doc.paragraphs = [mock_paragraph]
        mock_document.return_value = mock_doc
        
        text = handler.extract_text("test.docx")
        
        assert "Paragraph content" in text
        mock_document.assert_called_once_with("test.docx")
        

class TestTxtHandler:
    """Test TXT file handler"""
    
    @pytest.fixture
    def handler(self):
        return TxtHandler()
        
    def test_extract_text(self, handler):
        """Test text extraction from TXT file"""
        test_content = "This is test content"
        
        with patch('builtins.open', mock_open(read_data=test_content)):
            text = handler.extract_text("test.txt")
            
        assert text == test_content
        

class TestPptxHandler:
    """Test PPTX file handler"""
    
    @pytest.fixture
    def handler(self):
        return PptxHandler()
        
    @patch('utils.file_handlers.Presentation')
    def test_extract_text(self, mock_presentation, handler):
        """Test text extraction from PPTX"""
        mock_prs = Mock()
        mock_slide = Mock()
        mock_shape = Mock()
        mock_shape.has_text_frame = True
        mock_shape.text = "Slide content"
        mock_slide.shapes = [mock_shape]
        mock_prs.slides = [mock_slide]
        mock_presentation.return_value = mock_prs
        
        text = handler.extract_text("test.pptx")
        
        assert "Slide content" in text
        mock_presentation.assert_called_once_with("test.pptx")
        

class TestXlsxHandler:
    """Test XLSX file handler"""
    
    @pytest.fixture
    def handler(self):
        return XlsxHandler()
        
    @patch('utils.file_handlers.pd.read_excel')
    def test_extract_text(self, mock_read_excel, handler):
        """Test text extraction from XLSX"""
        import pandas as pd
        mock_df = pd.DataFrame({
            'Column1': ['Value1', 'Value2'],
            'Column2': ['Value3', 'Value4']
        })
        mock_read_excel.return_value = mock_df
        
        text = handler.extract_text("test.xlsx")
        
        assert "Column1" in text
        assert "Value1" in text
        mock_read_excel.assert_called_once_with("test.xlsx", sheet_name=None)