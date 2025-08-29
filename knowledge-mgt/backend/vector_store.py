"""Vector database operations for document storage and retrieval."""

from typing import List, Dict, Any, Optional, Tuple
from pathlib import Path
import json
import numpy as np
from datetime import datetime

import chromadb
from chromadb.config import Settings as ChromaSettings
from backend.config import settings
from backend.embeddings import EmbeddingManager


class VectorStore:
    """Manages vector database operations using ChromaDB."""
    
    def __init__(self, collection_name: str = "documents"):
        """Initialize the vector store.
        
        Args:
            collection_name: Name of the collection to use
        """
        self.collection_name = collection_name
        self.embedding_manager = EmbeddingManager()
        
        # Initialize ChromaDB client
        self.client = chromadb.PersistentClient(
            path=settings.vector_db_path,
            settings=ChromaSettings(
                anonymized_telemetry=False,
                allow_reset=True
            )
        )
        
        # Get or create collection
        self.collection = self._get_or_create_collection()
    
    def _get_or_create_collection(self):
        """Get existing collection or create new one."""
        try:
            collection = self.client.get_collection(name=self.collection_name)
        except:
            collection = self.client.create_collection(
                name=self.collection_name,
                metadata={"created_at": datetime.now().isoformat()}
            )
        return collection
    
    def add_documents(self, documents: List[Dict[str, Any]]) -> Dict[str, Any]:
        """Add documents to the vector store.
        
        Args:
            documents: List of processed documents with chunks
        
        Returns:
            Status dictionary with results
        """
        total_chunks = 0
        failed_chunks = 0
        
        for document in documents:
            doc_id = document["id"]
            chunks = document["chunks"]
            
            # Prepare data for insertion
            texts = []
            metadatas = []
            ids = []
            
            for chunk in chunks:
                chunk_id = f"{doc_id}_{chunk['metadata']['chunk_index']}"
                texts.append(chunk["content"])
                
                # Combine chunk metadata with document metadata
                chunk_metadata = {
                    **chunk["metadata"],
                    "document_title": document["title"],
                    "document_id": doc_id,
                    "source": document["metadata"].get("source", document["metadata"].get("filename", "unknown"))
                }
                
                # Ensure source is present
                if "source" not in chunk_metadata and "filename" in document["metadata"]:
                    chunk_metadata["source"] = document["metadata"]["filename"]
                    
                metadatas.append(chunk_metadata)
                ids.append(chunk_id)
            
            if texts:
                try:
                    # Generate embeddings
                    embeddings = self.embedding_manager.embed_texts(texts)
                    
                    # Add to ChromaDB
                    self.collection.add(
                        documents=texts,
                        embeddings=embeddings.tolist(),
                        metadatas=metadatas,
                        ids=ids
                    )
                    
                    total_chunks += len(texts)
                except Exception as e:
                    print(f"Error adding chunks for document {doc_id}: {e}")
                    failed_chunks += len(texts)
        
        return {
            "total_documents": len(documents),
            "total_chunks": total_chunks,
            "failed_chunks": failed_chunks,
            "success": failed_chunks == 0
        }
    
    def search(
        self,
        query: str,
        top_k: int = 5,
        filter_metadata: Optional[Dict[str, Any]] = None
    ) -> List[Dict[str, Any]]:
        """Search for relevant documents using semantic similarity.
        
        Args:
            query: Search query
            top_k: Number of results to return
            filter_metadata: Optional metadata filters
        
        Returns:
            List of search results with scores
        """
        # Generate query embedding
        query_embedding = self.embedding_manager.embed_query(query)
        
        # Prepare where clause for filtering
        where = filter_metadata if filter_metadata else None
        
        # Search in ChromaDB
        results = self.collection.query(
            query_embeddings=[query_embedding.tolist()],
            n_results=min(top_k, self.collection.count()),
            where=where,
            include=["documents", "metadatas", "distances"]
        )
        
        # Format results
        formatted_results = []
        if results and results["ids"] and results["ids"][0]:
            for idx in range(len(results["ids"][0])):
                formatted_results.append({
                    "chunk_id": results["ids"][0][idx],
                    "content": results["documents"][0][idx],
                    "metadata": results["metadatas"][0][idx],
                    "score": 1 - results["distances"][0][idx],  # Convert distance to similarity
                    "distance": results["distances"][0][idx]
                })
        
        return formatted_results
    
    def hybrid_search(
        self,
        query: str,
        top_k: int = 5,
        filter_metadata: Optional[Dict[str, Any]] = None,
        keyword_weight: float = 0.3
    ) -> List[Dict[str, Any]]:
        """Perform hybrid search combining semantic and keyword search.
        
        Args:
            query: Search query
            top_k: Number of results to return
            filter_metadata: Optional metadata filters
            keyword_weight: Weight for keyword search (0-1)
        
        Returns:
            List of search results with combined scores
        """
        # Semantic search
        semantic_results = self.search(query, top_k * 2, filter_metadata)
        
        # Simple keyword search (BM25-like scoring)
        keyword_results = self._keyword_search(query, top_k * 2, filter_metadata)
        
        # Combine and rerank results
        combined_results = self._combine_search_results(
            semantic_results,
            keyword_results,
            keyword_weight
        )
        
        # Return top k results
        return combined_results[:top_k]
    
    def _keyword_search(
        self,
        query: str,
        top_k: int,
        filter_metadata: Optional[Dict[str, Any]] = None
    ) -> List[Dict[str, Any]]:
        """Perform keyword-based search.
        
        Args:
            query: Search query
            top_k: Number of results to return
            filter_metadata: Optional metadata filters
        
        Returns:
            List of search results with keyword scores
        """
        # Get all documents (simplified for MVP)
        all_docs = self.collection.get(
            where=filter_metadata,
            include=["documents", "metadatas"]
        )
        
        if not all_docs["ids"]:
            return []
        
        # Calculate keyword scores
        query_terms = query.lower().split()
        results = []
        
        for idx, doc in enumerate(all_docs["documents"]):
            doc_lower = doc.lower()
            score = 0
            
            # Simple TF scoring
            for term in query_terms:
                score += doc_lower.count(term)
            
            if score > 0:
                results.append({
                    "chunk_id": all_docs["ids"][idx],
                    "content": doc,
                    "metadata": all_docs["metadatas"][idx],
                    "score": score,
                    "distance": 1 - (score / len(query_terms))  # Normalize
                })
        
        # Sort by score
        results.sort(key=lambda x: x["score"], reverse=True)
        
        return results[:top_k]
    
    def _combine_search_results(
        self,
        semantic_results: List[Dict[str, Any]],
        keyword_results: List[Dict[str, Any]],
        keyword_weight: float
    ) -> List[Dict[str, Any]]:
        """Combine semantic and keyword search results.
        
        Args:
            semantic_results: Results from semantic search
            keyword_results: Results from keyword search
            keyword_weight: Weight for keyword scores
        
        Returns:
            Combined and reranked results
        """
        # Create result map
        result_map = {}
        
        # Add semantic results
        for result in semantic_results:
            chunk_id = result["chunk_id"]
            result_map[chunk_id] = result.copy()
            result_map[chunk_id]["semantic_score"] = result["score"]
            result_map[chunk_id]["keyword_score"] = 0
        
        # Add keyword results
        for result in keyword_results:
            chunk_id = result["chunk_id"]
            if chunk_id in result_map:
                result_map[chunk_id]["keyword_score"] = result["score"]
            else:
                result_map[chunk_id] = result.copy()
                result_map[chunk_id]["semantic_score"] = 0
                result_map[chunk_id]["keyword_score"] = result["score"]
        
        # Calculate combined scores
        semantic_weight = 1 - keyword_weight
        for chunk_id, result in result_map.items():
            # Normalize scores
            max_semantic = max(
                (r.get("semantic_score", 0) for r in result_map.values()),
                default=1
            )
            max_keyword = max(
                (r.get("keyword_score", 0) for r in result_map.values()),
                default=1
            )
            
            norm_semantic = result["semantic_score"] / max_semantic if max_semantic > 0 else 0
            norm_keyword = result["keyword_score"] / max_keyword if max_keyword > 0 else 0
            
            # Combined score
            result["combined_score"] = (
                semantic_weight * norm_semantic +
                keyword_weight * norm_keyword
            )
            result["score"] = result["combined_score"]
        
        # Sort by combined score
        results = list(result_map.values())
        results.sort(key=lambda x: x["combined_score"], reverse=True)
        
        return results
    
    def get_document_by_id(self, document_id: str) -> Optional[Dict[str, Any]]:
        """Retrieve all chunks for a specific document.
        
        Args:
            document_id: Document ID
        
        Returns:
            Document data or None if not found
        """
        results = self.collection.get(
            where={"document_id": document_id},
            include=["documents", "metadatas"]
        )
        
        if not results["ids"]:
            return None
        
        # Reconstruct document
        chunks = []
        for idx in range(len(results["ids"])):
            chunks.append({
                "chunk_id": results["ids"][idx],
                "content": results["documents"][idx],
                "metadata": results["metadatas"][idx]
            })
        
        # Sort by chunk index
        chunks.sort(key=lambda x: x["metadata"].get("chunk_index", 0))
        
        return {
            "document_id": document_id,
            "title": chunks[0]["metadata"].get("document_title", ""),
            "chunks": chunks,
            "num_chunks": len(chunks)
        }
    
    def delete_document(self, document_id: str) -> bool:
        """Delete a document from the vector store.
        
        Args:
            document_id: Document ID to delete
        
        Returns:
            Success status
        """
        try:
            # Get all chunk IDs for the document
            results = self.collection.get(
                where={"document_id": document_id}
            )
            
            if results["ids"]:
                # Delete all chunks
                self.collection.delete(ids=results["ids"])
            
            return True
        except Exception as e:
            print(f"Error deleting document {document_id}: {e}")
            return False
    
    def get_collection_stats(self) -> Dict[str, Any]:
        """Get statistics about the vector store collection.
        
        Returns:
            Collection statistics
        """
        try:
            # Ensure we have a valid collection
            if not self.collection:
                self.collection = self._get_or_create_collection()
            
            # Get all metadata from the collection
            all_metadata = self.collection.get(include=["metadatas"])
            
            # Check if collection is empty
            if not all_metadata or not all_metadata.get("metadatas"):
                return {
                    "total_chunks": 0,
                    "total_documents": 0,
                    "collection_name": self.collection_name,
                    "embedding_dimension": self.embedding_manager.get_dimension()
                }
            
            # Count unique documents
            document_ids = set()
            document_sources = set()
            for metadata in all_metadata["metadatas"]:
                if metadata:
                    if "document_id" in metadata:
                        document_ids.add(metadata["document_id"])
                    if "source" in metadata:
                        document_sources.add(metadata["source"])
            
            # Use document_ids if available, otherwise use sources
            unique_docs = len(document_ids) if document_ids else len(document_sources)
            
            return {
                "total_chunks": self.collection.count(),
                "total_documents": unique_docs,
                "collection_name": self.collection_name,
                "embedding_dimension": self.embedding_manager.get_dimension()
            }
        except Exception as e:
            print(f"Error getting collection stats: {e}")
            return {
                "total_chunks": 0,
                "total_documents": 0,
                "collection_name": self.collection_name,
                "embedding_dimension": self.embedding_manager.get_dimension()
            }
    
    def clear_collection(self):
        """Clear all documents from the collection."""
        try:
            # Delete the existing collection
            self.client.delete_collection(name=self.collection_name)
            print(f"Deleted collection: {self.collection_name}")
            
            # Recreate an empty collection
            self.collection = self._get_or_create_collection()
            print(f"Recreated empty collection: {self.collection_name}")
            
            # Reset document count
            self.document_count = 0
            
            # Clear any cached data
            if hasattr(self, '_cache'):
                self._cache = {}
                
            return True
        except Exception as e:
            print(f"Error clearing collection: {e}")
            # Try alternative approach - delete all documents
            try:
                if self.collection:
                    # Get all document IDs and delete them
                    results = self.collection.get()
                    if results and 'ids' in results and results['ids']:
                        self.collection.delete(ids=results['ids'])
                        print(f"Deleted {len(results['ids'])} documents from collection")
                        self.document_count = 0
                        return True
            except Exception as e2:
                print(f"Alternative deletion also failed: {e2}")
            return False