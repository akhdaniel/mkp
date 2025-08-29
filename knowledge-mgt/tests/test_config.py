"""
Tests for configuration management
"""
import os
import pytest
from unittest.mock import patch
from backend.config import Settings


class TestSettings:
    """Test Settings configuration"""
    
    def test_default_settings(self):
        """Test that Settings loads with default values"""
        settings = Settings()
        assert settings.vector_db_type == "chromadb"
        assert settings.chunk_size == 512
        assert settings.chunk_overlap == 50
        assert settings.top_k_retrieval == 5
        assert settings.max_file_size == 10485760
        
    def test_settings_from_env(self):
        """Test Settings loads from environment variables"""
        with patch.dict(os.environ, {
            'VECTOR_DB_TYPE': 'qdrant',
            'CHUNK_SIZE': '1024',
            'TOP_K_RETRIEVAL': '10'
        }):
            settings = Settings()
            assert settings.vector_db_type == "qdrant"
            assert settings.chunk_size == 1024
            assert settings.top_k_retrieval == 10
            
    def test_paths_creation(self):
        """Test that required paths are set"""
        settings = Settings()
        assert settings.upload_path is not None
        assert settings.vector_db_path is not None
        assert settings.cache_path is not None
        
    def test_model_configuration(self):
        """Test model configuration settings"""
        settings = Settings()
        assert settings.model_name == "claude-3-opus-20240229"
        assert settings.embedding_model is not None