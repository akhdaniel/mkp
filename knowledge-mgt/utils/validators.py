"""Input validation utilities."""

import re
from pathlib import Path
from typing import Optional, Tuple


def validate_file_extension(filename: str, allowed_extensions: list) -> bool:
    """Validate file extension.
    
    Args:
        filename: Name of the file
        allowed_extensions: List of allowed extensions
    
    Returns:
        True if extension is valid
    """
    extension = Path(filename).suffix.lower().lstrip('.')
    return extension in allowed_extensions


def validate_file_size(file_size: int, max_size: int) -> bool:
    """Validate file size.
    
    Args:
        file_size: Size of file in bytes
        max_size: Maximum allowed size in bytes
    
    Returns:
        True if size is valid
    """
    return 0 < file_size <= max_size


def sanitize_text(text: str) -> str:
    """Sanitize text input.
    
    Args:
        text: Text to sanitize
    
    Returns:
        Sanitized text
    """
    # Remove control characters
    text = re.sub(r'[\x00-\x1f\x7f-\x9f]', '', text)
    
    # Normalize whitespace
    text = ' '.join(text.split())
    
    return text.strip()


def validate_api_key(api_key: str, prefix: Optional[str] = None) -> bool:
    """Validate API key format.
    
    Args:
        api_key: API key to validate
        prefix: Expected prefix for the key
    
    Returns:
        True if key appears valid
    """
    if not api_key or not isinstance(api_key, str):
        return False
    
    # Check minimum length
    if len(api_key) < 20:
        return False
    
    # Check prefix if provided
    if prefix and not api_key.startswith(prefix):
        return False
    
    # Check for valid characters (alphanumeric, dash, underscore)
    if not re.match(r'^[a-zA-Z0-9_-]+$', api_key):
        return False
    
    return True


def validate_query(query: str, min_length: int = 3, max_length: int = 1000) -> Tuple[bool, Optional[str]]:
    """Validate user query.
    
    Args:
        query: User query text
        min_length: Minimum query length
        max_length: Maximum query length
    
    Returns:
        Tuple of (is_valid, error_message)
    """
    if not query or not query.strip():
        return False, "Query cannot be empty"
    
    query = query.strip()
    
    if len(query) < min_length:
        return False, f"Query must be at least {min_length} characters"
    
    if len(query) > max_length:
        return False, f"Query must be no more than {max_length} characters"
    
    return True, None


def validate_metadata(metadata: dict) -> dict:
    """Validate and clean metadata dictionary.
    
    Args:
        metadata: Metadata dictionary
    
    Returns:
        Cleaned metadata
    """
    cleaned = {}
    
    for key, value in metadata.items():
        # Ensure key is string
        key = str(key)
        
        # Skip None values
        if value is None:
            continue
        
        # Convert values to appropriate types
        if isinstance(value, (str, int, float, bool)):
            cleaned[key] = value
        elif isinstance(value, (list, tuple)):
            cleaned[key] = list(value)
        elif isinstance(value, dict):
            cleaned[key] = validate_metadata(value)
        else:
            cleaned[key] = str(value)
    
    return cleaned