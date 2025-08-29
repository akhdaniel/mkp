"""Embedding generation for documents and queries."""

from typing import List, Optional, Union
from abc import ABC, abstractmethod
import numpy as np
import hashlib

try:
    from sentence_transformers import SentenceTransformer
    SENTENCE_TRANSFORMERS_AVAILABLE = True
except ImportError:
    SENTENCE_TRANSFORMERS_AVAILABLE = False

from backend.config import settings


class EmbeddingModel(ABC):
    """Abstract base class for embedding models."""
    
    @abstractmethod
    def encode(self, texts: Union[str, List[str]]) -> np.ndarray:
        """Generate embeddings for text(s)."""
        pass
    
    @abstractmethod
    def get_embedding_dimension(self) -> int:
        """Get the dimension of the embeddings."""
        pass


class SimpleHashEmbedding(EmbeddingModel):
    """Simple hash-based embedding for testing without heavy dependencies."""
    
    def __init__(self, dimension: int = 384):
        """Initialize the simple hash embedding model.
        
        Args:
            dimension: Dimension of embeddings
        """
        self.dimension = dimension
    
    def encode(self, texts: Union[str, List[str]]) -> np.ndarray:
        """Generate simple hash-based embeddings.
        
        Args:
            texts: Single text or list of texts
        
        Returns:
            Numpy array of embeddings
        """
        if isinstance(texts, str):
            texts = [texts]
        
        embeddings = []
        for text in texts:
            # Create a deterministic hash-based embedding
            hash_obj = hashlib.sha256(text.encode())
            hash_bytes = hash_obj.digest()
            
            # Convert to float array
            embedding = []
            for i in range(0, len(hash_bytes), 4):
                if len(embedding) >= self.dimension:
                    break
                chunk = hash_bytes[i:i+4]
                value = int.from_bytes(chunk, byteorder='big') / (2**32)
                embedding.append(value * 2 - 1)  # Scale to [-1, 1]
            
            # Pad or truncate to dimension
            if len(embedding) < self.dimension:
                embedding.extend([0.0] * (self.dimension - len(embedding)))
            else:
                embedding = embedding[:self.dimension]
            
            embeddings.append(embedding)
        
        return np.array(embeddings)
    
    def get_embedding_dimension(self) -> int:
        """Get the dimension of the embeddings."""
        return self.dimension


class SentenceTransformerEmbedding(EmbeddingModel):
    """Sentence Transformer embedding model."""
    
    def __init__(self, model_name: str = "all-MiniLM-L6-v2"):
        """Initialize the sentence transformer model.
        
        Args:
            model_name: Name of the sentence transformer model
        """
        if not SENTENCE_TRANSFORMERS_AVAILABLE:
            raise ImportError(
                "sentence-transformers is not installed. "
                "Install it with: pip install sentence-transformers"
            )
        self.model = SentenceTransformer(model_name)
        self.dimension = self.model.get_sentence_embedding_dimension()
    
    def encode(self, texts: Union[str, List[str]]) -> np.ndarray:
        """Generate embeddings for text(s).
        
        Args:
            texts: Single text or list of texts
        
        Returns:
            Numpy array of embeddings
        """
        if isinstance(texts, str):
            texts = [texts]
        
        embeddings = self.model.encode(
            texts,
            convert_to_numpy=True,
            show_progress_bar=False
        )
        
        return embeddings
    
    def get_embedding_dimension(self) -> int:
        """Get the dimension of the embeddings."""
        return self.dimension


class OpenAIEmbedding(EmbeddingModel):
    """OpenAI embedding model."""
    
    def __init__(self, model_name: str = "text-embedding-3-small"):
        """Initialize the OpenAI embedding model.
        
        Args:
            model_name: Name of the OpenAI embedding model
        """
        try:
            import openai
            self.client = openai.OpenAI(api_key=settings.openai_api_key)
        except ImportError:
            raise ImportError("OpenAI package not installed")
        
        if not settings.openai_api_key:
            raise ValueError("OpenAI API key not configured")
        
        self.model_name = model_name
        # Dimensions for different OpenAI models
        self.dimensions = {
            "text-embedding-3-small": 1536,
            "text-embedding-3-large": 3072,
            "text-embedding-ada-002": 1536
        }
        self.dimension = self.dimensions.get(model_name, 1536)
    
    def encode(self, texts: Union[str, List[str]]) -> np.ndarray:
        """Generate embeddings for text(s).
        
        Args:
            texts: Single text or list of texts
        
        Returns:
            Numpy array of embeddings
        """
        if isinstance(texts, str):
            texts = [texts]
        
        # OpenAI API call
        response = self.client.embeddings.create(
            model=self.model_name,
            input=texts
        )
        
        # Extract embeddings
        embeddings = [item.embedding for item in response.data]
        
        return np.array(embeddings)
    
    def get_embedding_dimension(self) -> int:
        """Get the dimension of the embeddings."""
        return self.dimension


class EmbeddingFactory:
    """Factory for creating embedding models."""
    
    @staticmethod
    def create_embedding_model(
        model_name: Optional[str] = None
    ) -> EmbeddingModel:
        """Create an embedding model based on configuration.
        
        Args:
            model_name: Optional model name override
        
        Returns:
            Embedding model instance
        """
        model_name = model_name or settings.embedding_model
        
        if settings.use_openai_embeddings or model_name.startswith("text-embedding"):
            try:
                return OpenAIEmbedding(model_name)
            except Exception as e:
                print(f"Warning: Could not initialize OpenAI embeddings: {e}")
                print("Falling back to simple embeddings")
                return SimpleHashEmbedding()
        elif SENTENCE_TRANSFORMERS_AVAILABLE:
            try:
                return SentenceTransformerEmbedding(model_name)
            except Exception as e:
                print(f"Warning: Could not initialize SentenceTransformers: {e}")
                print("Falling back to simple embeddings")
                return SimpleHashEmbedding()
        else:
            print("Note: Using simple hash-based embeddings. For better results, install sentence-transformers")
            return SimpleHashEmbedding()


class EmbeddingManager:
    """Manages embedding generation and caching."""
    
    def __init__(self, model: Optional[EmbeddingModel] = None):
        """Initialize the embedding manager.
        
        Args:
            model: Optional embedding model instance
        """
        self.model = model or EmbeddingFactory.create_embedding_model()
        self._cache = {}
    
    def embed_texts(
        self,
        texts: Union[str, List[str]],
        use_cache: bool = True
    ) -> np.ndarray:
        """Generate embeddings for texts with optional caching.
        
        Args:
            texts: Text or list of texts to embed
            use_cache: Whether to use caching
        
        Returns:
            Numpy array of embeddings
        """
        if isinstance(texts, str):
            texts = [texts]
        
        if use_cache:
            # Check cache and generate only for missing texts
            embeddings = []
            texts_to_embed = []
            indices_to_embed = []
            
            for i, text in enumerate(texts):
                cache_key = self._get_cache_key(text)
                if cache_key in self._cache:
                    embeddings.append(self._cache[cache_key])
                else:
                    texts_to_embed.append(text)
                    indices_to_embed.append(i)
                    embeddings.append(None)
            
            # Generate embeddings for uncached texts
            if texts_to_embed:
                new_embeddings = self.model.encode(texts_to_embed)
                for idx, embedding in zip(indices_to_embed, new_embeddings):
                    embeddings[idx] = embedding
                    # Cache the embedding
                    cache_key = self._get_cache_key(texts[idx])
                    self._cache[cache_key] = embedding
            
            return np.array(embeddings)
        else:
            return self.model.encode(texts)
    
    def embed_query(self, query: str) -> np.ndarray:
        """Generate embedding for a query.
        
        Args:
            query: Query text
        
        Returns:
            Query embedding
        """
        # Queries typically shouldn't be cached
        return self.model.encode(query)[0]
    
    def get_dimension(self) -> int:
        """Get the embedding dimension."""
        return self.model.get_embedding_dimension()
    
    def clear_cache(self):
        """Clear the embedding cache."""
        self._cache.clear()
    
    def _get_cache_key(self, text: str) -> str:
        """Generate a cache key for text.
        
        Args:
            text: Text to generate key for
        
        Returns:
            Cache key
        """
        import hashlib
        return hashlib.md5(text.encode()).hexdigest()