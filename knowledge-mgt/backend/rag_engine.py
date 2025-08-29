"""RAG (Retrieval-Augmented Generation) engine for query processing."""

from typing import List, Dict, Any, Optional, Tuple
from datetime import datetime
import json

from anthropic import Anthropic
from backend.config import settings
from backend.vector_store import VectorStore


class RAGEngine:
    """Orchestrates the RAG pipeline for generating responses."""
    
    def __init__(self, vector_store: Optional[VectorStore] = None):
        """Initialize the RAG engine.
        
        Args:
            vector_store: Optional vector store instance
        """
        self.vector_store = vector_store or VectorStore()
        self.anthropic_client = Anthropic(api_key=settings.anthropic_api_key)
        self.model_name = settings.model_name
        self.top_k = settings.top_k_retrieval
    
    def generate_response(
        self,
        query: str,
        conversation_history: Optional[List[Dict[str, str]]] = None,
        filter_metadata: Optional[Dict[str, Any]] = None,
        use_hybrid_search: bool = True
    ) -> Dict[str, Any]:
        """Generate a response to a user query.
        
        Args:
            query: User query
            conversation_history: Previous conversation messages
            filter_metadata: Optional metadata filters for search
            use_hybrid_search: Whether to use hybrid search
        
        Returns:
            Response dictionary with answer and metadata
        """
        start_time = datetime.now()
        
        # Retrieve relevant documents
        if use_hybrid_search:
            search_results = self.vector_store.hybrid_search(
                query=query,
                top_k=self.top_k,
                filter_metadata=filter_metadata
            )
        else:
            search_results = self.vector_store.search(
                query=query,
                top_k=self.top_k,
                filter_metadata=filter_metadata
            )
        
        # Rerank results if needed
        reranked_results = self._rerank_results(query, search_results)
        
        # Build context from search results
        context, sources = self._build_context(reranked_results)
        
        # Generate response using Claude
        response_text = self._generate_claude_response(
            query=query,
            context=context,
            conversation_history=conversation_history
        )
        
        # Calculate confidence score
        confidence_score = self._calculate_confidence(reranked_results)
        
        # Prepare response
        response = {
            "answer": response_text,
            "sources": sources,
            "confidence": confidence_score,
            "num_sources": len(sources),
            "search_results": reranked_results[:3],  # Top 3 for reference
            "processing_time": (datetime.now() - start_time).total_seconds(),
            "timestamp": datetime.now().isoformat()
        }
        
        return response
    
    def _rerank_results(
        self,
        query: str,
        search_results: List[Dict[str, Any]]
    ) -> List[Dict[str, Any]]:
        """Rerank search results for better relevance.
        
        Args:
            query: User query
            search_results: Initial search results
        
        Returns:
            Reranked results
        """
        # Simple reranking based on query term overlap
        query_terms = set(query.lower().split())
        
        for result in search_results:
            content_terms = set(result["content"].lower().split())
            overlap = len(query_terms & content_terms)
            
            # Boost score based on term overlap
            result["rerank_score"] = result["score"] * (1 + overlap * 0.1)
        
        # Sort by reranked score
        search_results.sort(key=lambda x: x["rerank_score"], reverse=True)
        
        return search_results
    
    def _build_context(
        self,
        search_results: List[Dict[str, Any]]
    ) -> Tuple[str, List[Dict[str, str]]]:
        """Build context string from search results.
        
        Args:
            search_results: Search results to build context from
        
        Returns:
            Tuple of (context string, source citations)
        """
        if not search_results:
            return "", []
        
        context_parts = []
        sources = []
        seen_docs = set()
        
        for idx, result in enumerate(search_results):
            doc_id = result["metadata"].get("document_id", "unknown")
            doc_title = result["metadata"].get("document_title", "Unknown Document")
            chunk_index = result["metadata"].get("chunk_index", 0)
            
            # Add to context
            context_parts.append(
                f"[Source {idx + 1} - {doc_title} (Section {chunk_index + 1})]:\n"
                f"{result['content']}\n"
            )
            
            # Track unique sources
            if doc_id not in seen_docs:
                sources.append({
                    "document_id": doc_id,
                    "title": doc_title,
                    "filename": result["metadata"].get("filename", ""),
                    "relevance_score": result.get("rerank_score", result["score"])
                })
                seen_docs.add(doc_id)
        
        context = "\n---\n".join(context_parts)
        return context, sources
    
    def _generate_claude_response(
        self,
        query: str,
        context: str,
        conversation_history: Optional[List[Dict[str, str]]] = None
    ) -> str:
        """Generate response using Claude API.
        
        Args:
            query: User query
            context: Retrieved context
            conversation_history: Previous conversation
        
        Returns:
            Generated response text
        """
        # Build system prompt
        system_prompt = """You are an AI assistant helping employees with questions about internal documents and procedures. 
        Use the provided context to answer questions accurately. Always cite the source when providing information.
        If the context doesn't contain enough information to answer the question, say so clearly.
        Be concise but thorough in your responses."""
        
        # Build messages
        messages = []
        
        # Add conversation history if provided
        if conversation_history:
            for msg in conversation_history[-5:]:  # Last 5 messages for context
                messages.append({
                    "role": msg["role"],
                    "content": msg["content"]
                })
        
        # Build the user message with context
        if context:
            user_message = f"""Context from relevant documents:
{context}

Question: {query}

Please provide a helpful answer based on the context above. Cite specific sources when possible."""
        else:
            user_message = f"""Question: {query}

I couldn't find any relevant documents to answer this question. Please provide a general response or suggest what documents might be helpful."""
        
        messages.append({
            "role": "user",
            "content": user_message
        })
        
        try:
            # Call Claude API
            response = self.anthropic_client.messages.create(
                model=self.model_name,
                max_tokens=1000,
                temperature=0.3,
                system=system_prompt,
                messages=messages
            )
            
            return response.content[0].text
        
        except Exception as e:
            return f"Error generating response: {str(e)}"
    
    def _calculate_confidence(
        self,
        search_results: List[Dict[str, Any]]
    ) -> float:
        """Calculate confidence score based on search results.
        
        Args:
            search_results: Search results
        
        Returns:
            Confidence score between 0 and 1
        """
        if not search_results:
            return 0.0
        
        # Calculate based on top result scores
        top_scores = [r.get("rerank_score", r["score"]) for r in search_results[:3]]
        
        if not top_scores:
            return 0.0
        
        # Average of top scores (assuming scores are 0-1)
        avg_score = sum(top_scores) / len(top_scores)
        
        # Boost confidence if multiple high-scoring results
        if len(top_scores) >= 3 and all(s > 0.7 for s in top_scores):
            avg_score = min(1.0, avg_score * 1.2)
        
        return round(avg_score, 2)
    
    def process_feedback(
        self,
        query: str,
        response: Dict[str, Any],
        feedback: str,
        rating: Optional[int] = None
    ) -> bool:
        """Process user feedback on a response.
        
        Args:
            query: Original query
            response: Generated response
            feedback: User feedback text
            rating: Optional rating (1-5)
        
        Returns:
            Success status
        """
        # In a production system, this would store feedback for analysis
        # For MVP, we'll just log it
        feedback_data = {
            "timestamp": datetime.now().isoformat(),
            "query": query,
            "response_summary": response["answer"][:200],
            "feedback": feedback,
            "rating": rating,
            "confidence": response.get("confidence", 0)
        }
        
        # Could write to a feedback log file or database
        print(f"Feedback received: {json.dumps(feedback_data, indent=2)}")
        
        return True
    
    def get_conversation_summary(
        self,
        conversation_history: List[Dict[str, str]]
    ) -> str:
        """Generate a summary of the conversation.
        
        Args:
            conversation_history: List of conversation messages
        
        Returns:
            Summary text
        """
        if not conversation_history:
            return "No conversation history available."
        
        # Build conversation text
        conversation_text = ""
        for msg in conversation_history:
            role = "User" if msg["role"] == "user" else "Assistant"
            conversation_text += f"{role}: {msg['content']}\n\n"
        
        # Generate summary using Claude
        try:
            response = self.anthropic_client.messages.create(
                model=self.model_name,
                max_tokens=500,
                temperature=0.3,
                system="Summarize the following conversation concisely, highlighting key topics discussed and any important information provided.",
                messages=[{
                    "role": "user",
                    "content": f"Please summarize this conversation:\n\n{conversation_text}"
                }]
            )
            
            return response.content[0].text
        
        except Exception as e:
            return f"Error generating summary: {str(e)}"