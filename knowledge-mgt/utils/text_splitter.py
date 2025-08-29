"""Text splitting utilities for document chunking."""

from typing import List, Optional, Dict, Any
from langchain.text_splitter import RecursiveCharacterTextSplitter
import tiktoken


class DocumentChunker:
    """Handles document chunking with various strategies."""
    
    def __init__(
        self,
        chunk_size: int = 512,
        chunk_overlap: int = 50,
        model_name: str = "gpt-3.5-turbo"
    ):
        """Initialize the document chunker.
        
        Args:
            chunk_size: Maximum size of each chunk in tokens
            chunk_overlap: Number of overlapping tokens between chunks
            model_name: Model name for tokenizer
        """
        self.chunk_size = chunk_size
        self.chunk_overlap = chunk_overlap
        self.model_name = model_name
        
        # Initialize tokenizer
        try:
            self.encoding = tiktoken.encoding_for_model(model_name)
        except:
            self.encoding = tiktoken.get_encoding("cl100k_base")
        
        # Initialize text splitter
        self.text_splitter = RecursiveCharacterTextSplitter(
            chunk_size=chunk_size,
            chunk_overlap=chunk_overlap,
            length_function=self._token_length,
            separators=["\n\n", "\n", ". ", " ", ""]
        )
    
    def _token_length(self, text: str) -> int:
        """Calculate the number of tokens in text."""
        return len(self.encoding.encode(text))
    
    def chunk_text(
        self, 
        text: str, 
        metadata: Optional[Dict[str, Any]] = None
    ) -> List[Dict[str, Any]]:
        """Split text into chunks with metadata.
        
        Args:
            text: Text to split
            metadata: Optional metadata to attach to each chunk
        
        Returns:
            List of chunks with text and metadata
        """
        if not text or not text.strip():
            return []
        
        # Split text into chunks
        chunks = self.text_splitter.split_text(text)
        
        # Prepare chunk documents with metadata
        chunk_docs = []
        for i, chunk in enumerate(chunks):
            chunk_metadata = {
                "chunk_index": i,
                "total_chunks": len(chunks),
                **(metadata or {})
            }
            
            chunk_docs.append({
                "content": chunk,
                "metadata": chunk_metadata
            })
        
        return chunk_docs
    
    def chunk_with_headers(
        self, 
        text: str, 
        metadata: Optional[Dict[str, Any]] = None
    ) -> List[Dict[str, Any]]:
        """Split text while preserving section headers.
        
        This method attempts to keep section headers with their content
        for better context preservation.
        
        Args:
            text: Text to split
            metadata: Optional metadata to attach to each chunk
        
        Returns:
            List of chunks with text and metadata
        """
        # Split by common header patterns
        sections = self._split_by_headers(text)
        
        all_chunks = []
        for section_idx, section in enumerate(sections):
            section_chunks = self.chunk_text(
                section['content'],
                {
                    **(metadata or {}),
                    "section_title": section.get('title', ''),
                    "section_index": section_idx
                }
            )
            all_chunks.extend(section_chunks)
        
        return all_chunks
    
    def _split_by_headers(self, text: str) -> List[Dict[str, str]]:
        """Split text by header patterns.
        
        Args:
            text: Text to split
        
        Returns:
            List of sections with title and content
        """
        import re
        
        # Common header patterns
        header_patterns = [
            r'^#{1,6}\s+(.+)$',  # Markdown headers
            r'^(.+)\n[=\-]{3,}$',  # Underlined headers
            r'^\d+\.\s+(.+)$',  # Numbered sections
            r'^[A-Z][A-Z\s]+:$'  # ALL CAPS headers
        ]
        
        sections = []
        current_section = {"title": "Introduction", "content": ""}
        
        lines = text.split('\n')
        for line in lines:
            is_header = False
            header_text = ""
            
            for pattern in header_patterns:
                match = re.match(pattern, line, re.MULTILINE)
                if match:
                    is_header = True
                    header_text = match.group(1) if match.lastindex else line
                    break
            
            if is_header and current_section['content'].strip():
                # Save current section and start new one
                sections.append(current_section)
                current_section = {"title": header_text, "content": ""}
            else:
                current_section['content'] += line + '\n'
        
        # Add the last section
        if current_section['content'].strip():
            sections.append(current_section)
        
        # If no sections were found, return the whole text as one section
        if not sections:
            sections = [{"title": "Document", "content": text}]
        
        return sections