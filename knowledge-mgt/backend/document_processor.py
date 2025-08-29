"""Document processing and ingestion pipeline."""

import os
import hashlib
from pathlib import Path
from typing import List, Dict, Any, Optional
from datetime import datetime
import json

from utils.file_handlers import FileHandlerFactory
from utils.text_splitter import DocumentChunker
from backend.config import settings


class DocumentProcessor:
    """Handles document ingestion, processing, and chunking."""
    
    def __init__(self):
        """Initialize the document processor."""
        self.chunker = DocumentChunker(
            chunk_size=settings.chunk_size,
            chunk_overlap=settings.chunk_overlap
        )
        self.upload_path = Path(settings.upload_path)
        self.upload_path.mkdir(parents=True, exist_ok=True)
    
    def process_document(
        self,
        file_path: Path,
        metadata: Optional[Dict[str, Any]] = None
    ) -> Dict[str, Any]:
        """Process a single document.
        
        Args:
            file_path: Path to the document file
            metadata: Optional metadata to attach to the document
        
        Returns:
            Processed document with chunks and metadata
        """
        # Get appropriate file handler
        handler = FileHandlerFactory.get_handler(file_path)
        
        # Extract text and metadata
        text = handler.extract_text(file_path)
        file_metadata = handler.get_metadata(file_path)
        
        # Generate document ID
        doc_id = self._generate_doc_id(file_path, text)
        
        # Combine metadata
        combined_metadata = {
            "document_id": doc_id,
            "filename": file_path.name,
            "file_type": file_path.suffix.lower(),
            "file_size": file_path.stat().st_size,
            "upload_date": datetime.now().isoformat(),
            **file_metadata,
            **(metadata or {})
        }
        
        # Chunk the document
        chunks = self.chunker.chunk_with_headers(text, combined_metadata)
        
        # Prepare document object
        document = {
            "id": doc_id,
            "title": combined_metadata.get("title", file_path.stem),
            "content": text,
            "chunks": chunks,
            "metadata": combined_metadata,
            "num_chunks": len(chunks),
            "total_tokens": sum(
                self.chunker._token_length(chunk["content"]) 
                for chunk in chunks
            )
        }
        
        return document
    
    def process_batch(
        self,
        file_paths: List[Path],
        metadata: Optional[Dict[str, Any]] = None
    ) -> List[Dict[str, Any]]:
        """Process multiple documents in batch.
        
        Args:
            file_paths: List of paths to document files
            metadata: Optional metadata to attach to all documents
        
        Returns:
            List of processed documents
        """
        documents = []
        errors = []
        
        for file_path in file_paths:
            try:
                document = self.process_document(file_path, metadata)
                documents.append(document)
            except Exception as e:
                errors.append({
                    "file": str(file_path),
                    "error": str(e)
                })
        
        if errors:
            # Log errors but continue with successful documents
            for error in errors:
                print(f"Error processing {error['file']}: {error['error']}")
        
        return documents
    
    def save_uploaded_file(
        self,
        file_content: bytes,
        filename: str
    ) -> Path:
        """Save an uploaded file to the upload directory.
        
        Args:
            file_content: File content as bytes
            filename: Original filename
        
        Returns:
            Path to the saved file
        """
        # Sanitize filename
        safe_filename = self._sanitize_filename(filename)
        
        # Generate unique filename if exists
        file_path = self.upload_path / safe_filename
        if file_path.exists():
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
            name_parts = safe_filename.rsplit('.', 1)
            if len(name_parts) == 2:
                safe_filename = f"{name_parts[0]}_{timestamp}.{name_parts[1]}"
            else:
                safe_filename = f"{safe_filename}_{timestamp}"
            file_path = self.upload_path / safe_filename
        
        # Save file
        with open(file_path, 'wb') as f:
            f.write(file_content)
        
        return file_path
    
    def validate_file(
        self,
        filename: str,
        file_size: int
    ) -> tuple[bool, Optional[str]]:
        """Validate file before processing.
        
        Args:
            filename: Name of the file
            file_size: Size of the file in bytes
        
        Returns:
            Tuple of (is_valid, error_message)
        """
        # Check file extension
        extension = Path(filename).suffix.lower().lstrip('.')
        if extension not in settings.allowed_extensions:
            return False, f"File type '{extension}' is not supported"
        
        # Check file size
        if file_size > settings.max_file_size:
            max_size_mb = settings.max_file_size / (1024 * 1024)
            return False, f"File size exceeds maximum of {max_size_mb:.1f} MB"
        
        return True, None
    
    def _generate_doc_id(self, file_path: Path, content: str) -> str:
        """Generate a unique document ID.
        
        Args:
            file_path: Path to the document
            content: Document content
        
        Returns:
            Unique document ID
        """
        # Create hash from filename and content
        hash_input = f"{file_path.name}:{content[:1000]}"
        return hashlib.md5(hash_input.encode()).hexdigest()
    
    def _sanitize_filename(self, filename: str) -> str:
        """Sanitize filename for safe storage.
        
        Args:
            filename: Original filename
        
        Returns:
            Sanitized filename
        """
        # Remove potentially dangerous characters
        import re
        # Keep only alphanumeric, dots, hyphens, and underscores
        safe_name = re.sub(r'[^a-zA-Z0-9._-]', '_', filename)
        # Remove multiple consecutive underscores
        safe_name = re.sub(r'_+', '_', safe_name)
        return safe_name
    
    def get_document_stats(self, document: Dict[str, Any]) -> Dict[str, Any]:
        """Get statistics about a processed document.
        
        Args:
            document: Processed document
        
        Returns:
            Document statistics
        """
        return {
            "document_id": document["id"],
            "title": document["title"],
            "num_chunks": document["num_chunks"],
            "total_tokens": document["total_tokens"],
            "avg_chunk_tokens": (
                document["total_tokens"] / document["num_chunks"] 
                if document["num_chunks"] > 0 else 0
            ),
            "file_type": document["metadata"]["file_type"],
            "file_size": document["metadata"]["file_size"],
            "upload_date": document["metadata"]["upload_date"]
        }