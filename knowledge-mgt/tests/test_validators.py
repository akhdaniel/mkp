"""
Tests for input validators
"""
import pytest
from pathlib import Path
from utils.validators import (
    validate_file_extension,
    validate_file_size,
    sanitize_text,
    validate_api_key,
    validate_query,
    validate_metadata
)


class TestValidators:
    """Test input validation functions"""
    
    def test_validate_file_extension_valid(self):
        """Test file extension validation with valid types"""
        allowed = ['.pdf', '.docx', '.pptx', '.txt', '.xlsx']
        assert validate_file_extension("document.pdf", allowed) is True
        assert validate_file_extension("report.docx", allowed) is True
        assert validate_file_extension("presentation.pptx", allowed) is True
        assert validate_file_extension("notes.txt", allowed) is True
        assert validate_file_extension("data.xlsx", allowed) is True
        
    def test_validate_file_extension_invalid(self):
        """Test file extension validation with invalid types"""
        allowed = ['.pdf', '.docx', '.pptx', '.txt', '.xlsx']
        assert validate_file_extension("script.py", allowed) is False
        assert validate_file_extension("image.jpg", allowed) is False
        assert validate_file_extension("video.mp4", allowed) is False
        assert validate_file_extension("archive.zip", allowed) is False
        
    def test_validate_file_extension_case_insensitive(self):
        """Test that file extension validation is case insensitive"""
        allowed = ['.pdf', '.docx', '.xlsx']
        assert validate_file_extension("Document.PDF", allowed) is True
        assert validate_file_extension("Report.DOCX", allowed) is True
        assert validate_file_extension("Data.XLSX", allowed) is True
        
    def test_validate_file_size_valid(self):
        """Test file size validation with valid size"""
        assert validate_file_size(1024, 10485760) is True  # 1KB
        assert validate_file_size(1048576, 10485760) is True  # 1MB
        assert validate_file_size(10485760, 10485760) is True  # 10MB (max)
        
    def test_validate_file_size_invalid(self):
        """Test file size validation with invalid size"""
        assert validate_file_size(10485761, 10485760) is False  # Just over 10MB
        assert validate_file_size(20971520, 10485760) is False  # 20MB
        assert validate_file_size(-1, 10485760) is False  # Negative size
        
    def test_validate_file_size_custom_max(self):
        """Test file size validation with custom max size"""
        assert validate_file_size(5242880, 5242880) is True  # 5MB with 5MB max
        assert validate_file_size(5242881, 5242880) is False  # Over 5MB max
        
    def test_sanitize_text(self):
        """Test text sanitization"""
        # Test removing special characters
        assert sanitize_text("Hello\x00World") == "HelloWorld"
        assert sanitize_text("Test\nLine\rBreak") == "Test\nLine\nBreak"
        
        # Test stripping whitespace
        text = "  Some text  "
        assert sanitize_text(text).strip() == "Some text"
        
        # Test handling empty/None
        assert sanitize_text("") == ""
        assert sanitize_text("   ") == "   "
        
    def test_validate_api_key_valid(self):
        """Test API key validation with valid keys"""
        assert validate_api_key("sk-proj-abcd1234efgh5678ijkl") == (True, None)
        assert validate_api_key("anthropic-api-key-12345") == (True, None)
        assert validate_api_key("a" * 20) == (True, None)  # Minimum length key
        
    def test_validate_api_key_invalid(self):
        """Test API key validation with invalid keys"""
        valid, msg = validate_api_key("")
        assert valid is False
        assert msg is not None
        
        valid, msg = validate_api_key("short")
        assert valid is False
        assert msg is not None
        
    def test_validate_api_key_with_prefix(self):
        """Test API key validation with specific prefix"""
        assert validate_api_key("sk-proj-12345", "sk-proj-") == (True, None)
        valid, msg = validate_api_key("wrong-prefix-12345", "sk-proj-")
        assert valid is False
        assert "prefix" in msg.lower()
        
    def test_validate_query_valid(self):
        """Test query validation with valid queries"""
        assert validate_query("What is the vacation policy?")[0] is True
        assert validate_query("How do I submit expenses?")[0] is True
        assert validate_query("ABC")[0] is True  # Minimum length
        
    def test_validate_query_invalid(self):
        """Test query validation with invalid queries"""
        valid, msg = validate_query("")
        assert valid is False
        assert msg is not None
        
        valid, msg = validate_query("ab")  # Too short (less than 3)
        assert valid is False
        assert "at least" in msg.lower()
        
        valid, msg = validate_query("a" * 1001)  # Too long
        assert valid is False
        assert "exceed" in msg.lower()
        
    def test_validate_metadata(self):
        """Test metadata validation and sanitization"""
        metadata = {
            "source": "test.pdf",
            "page": 1,
            "extra_field": "value"
        }
        
        validated = validate_metadata(metadata)
        assert "source" in validated
        assert validated["source"] == "test.pdf"
        assert "page" in validated
        assert validated["page"] == 1
        
    def test_validate_metadata_sanitization(self):
        """Test that metadata values are sanitized"""
        metadata = {
            "source": "test\x00file.pdf",
            "description": "  Some description  "
        }
        
        validated = validate_metadata(metadata)
        assert "\x00" not in validated["source"]
        assert validated["description"].strip() == "Some description"