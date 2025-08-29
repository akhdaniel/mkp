"""Configuration management for the AI Helpdesk Chat application."""

import os
from pathlib import Path
from typing import Optional
from pydantic_settings import BaseSettings
from pydantic import Field
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

class Settings(BaseSettings):
    """Application settings with environment variable support."""
    
    # API Keys
    anthropic_api_key: str = Field(default="", env="ANTHROPIC_API_KEY")
    openai_api_key: Optional[str] = Field(default=None, env="OPENAI_API_KEY")
    
    # Database Configuration
    vector_db_path: str = Field(default="./data/vectordb", env="VECTOR_DB_PATH")
    vector_db_type: str = Field(default="chromadb", env="VECTOR_DB_TYPE")
    
    # File Upload Configuration
    upload_path: str = Field(default="./data/uploads", env="UPLOAD_PATH")
    max_file_size: int = Field(default=10485760, env="MAX_FILE_SIZE")  # 10MB
    allowed_extensions: list = Field(
        default=["pdf", "docx", "pptx", "txt", "xlsx"]
    )
    
    # RAG Configuration
    chunk_size: int = Field(default=512, env="CHUNK_SIZE")
    chunk_overlap: int = Field(default=50, env="CHUNK_OVERLAP")
    top_k_retrieval: int = Field(default=5, env="TOP_K_RETRIEVAL")
    
    # Model Configuration
    model_name: str = Field(
        default="claude-3-opus-20240229", 
        env="MODEL_NAME"
    )
    embedding_model: str = Field(
        default="all-MiniLM-L6-v2", 
        env="EMBEDDING_MODEL"
    )
    use_openai_embeddings: bool = Field(default=False)
    
    # Multi-LLM Support
    llm_provider: str = Field(
        default="anthropic",
        env="LLM_PROVIDER"
    )
    deepseek_api_key: Optional[str] = Field(
        default=None,
        env="DEEPSEEK_API_KEY"
    )
    deepseek_model: str = Field(
        default="deepseek-chat",
        env="DEEPSEEK_MODEL"
    )
    
    # Application Configuration
    log_level: str = Field(default="INFO", env="LOG_LEVEL")
    cache_path: str = Field(default="./data/cache", env="CACHE_PATH")
    session_timeout: int = Field(default=3600, env="SESSION_TIMEOUT")
    
    class Config:
        env_file = ".env"
        env_file_encoding = "utf-8"
    
    def __init__(self, **kwargs):
        super().__init__(**kwargs)
        # Determine if we should use OpenAI embeddings
        if self.embedding_model.startswith("text-embedding"):
            self.use_openai_embeddings = True
        
        # Create necessary directories
        self._create_directories()
    
    def _create_directories(self):
        """Create necessary directories if they don't exist."""
        directories = [
            self.vector_db_path,
            self.upload_path,
            self.cache_path,
            Path("logs")
        ]
        for directory in directories:
            Path(directory).mkdir(parents=True, exist_ok=True)
    
    def validate_api_keys(self) -> tuple[bool, list[str]]:
        """Validate that required API keys are present."""
        errors = []
        
        if not self.anthropic_api_key:
            errors.append("ANTHROPIC_API_KEY is not set")
        
        if self.use_openai_embeddings and not self.openai_api_key:
            errors.append("OPENAI_API_KEY is required for OpenAI embeddings")
        
        return len(errors) == 0, errors

# Singleton instance
settings = Settings()