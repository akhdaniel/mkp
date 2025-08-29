"""
DeepSeek RAG Engine implementation
Uses DeepSeek's OpenAI-compatible API for response generation
"""
import os
from typing import List, Dict, Any, Optional
from openai import OpenAI
from backend.vector_store import VectorStore
from backend.config import Settings


class DeepSeekEngine:
    """RAG Engine using DeepSeek model"""
    
    def __init__(self, api_key: Optional[str] = None):
        """Initialize DeepSeek client with OpenAI-compatible interface"""
        self.settings = Settings()
        
        # DeepSeek uses OpenAI-compatible API
        self.api_key = api_key or os.getenv("DEEPSEEK_API_KEY")
        if not self.api_key:
            raise ValueError("DeepSeek API key not found. Set DEEPSEEK_API_KEY in .env file")
            
        self.client = OpenAI(
            api_key=self.api_key,
            base_url="https://api.deepseek.com/v1"  # DeepSeek endpoint
        )
        
        # Initialize vector store
        self.vector_store = VectorStore()
        
        # Model selection - DeepSeek models
        self.model_name = os.getenv("DEEPSEEK_MODEL", "deepseek-chat")  # or "deepseek-coder" for code-focused tasks
        
    def retrieve_context(self, query: str, k: int = None) -> List[Dict[str, Any]]:
        """Retrieve relevant context from vector store"""
        k = k or self.settings.top_k_retrieval
        return self.vector_store.search(query, k=k)
    
    def format_context(self, context_docs: List[Dict[str, Any]]) -> str:
        """Format context documents for prompt"""
        if not context_docs:
            return "No relevant context found."
        
        formatted_context = []
        for i, doc in enumerate(context_docs, 1):
            content = doc.get('content', '')
            metadata = doc.get('metadata', {})
            source = metadata.get('source', 'Unknown')
            page = metadata.get('page', 'N/A')
            
            formatted_context.append(
                f"[Document {i}]\n"
                f"Source: {source} (Page {page})\n"
                f"Content: {content}\n"
            )
        
        return "\n---\n".join(formatted_context)
    
    def generate_response(
        self, 
        query: str, 
        context: Optional[List[Dict[str, Any]]] = None,
        history: Optional[List[Dict[str, str]]] = None,
        temperature: float = 0.7,
        max_tokens: int = 2000
    ) -> str:
        """Generate response using DeepSeek model"""
        
        # Retrieve context if not provided
        if context is None:
            context = self.retrieve_context(query)
        
        # Format context for prompt
        formatted_context = self.format_context(context)
        
        # Build system prompt
        system_prompt = """You are a helpful AI assistant for an internal knowledge management system.
Your role is to answer questions based on the provided document context.
Always cite your sources when providing information from the documents.
If the context doesn't contain relevant information, say so clearly."""
        
        # Build messages
        messages = [
            {"role": "system", "content": system_prompt}
        ]
        
        # Add conversation history if provided
        if history:
            for msg in history[-10:]:  # Limit to last 10 messages for context window
                messages.append(msg)
        
        # Add current query with context
        user_message = f"""Context from documents:
{formatted_context}

User Question: {query}

Please provide a comprehensive answer based on the context above. Include citations to specific documents when referencing information."""
        
        messages.append({"role": "user", "content": user_message})
        
        try:
            # Call DeepSeek API
            response = self.client.chat.completions.create(
                model=self.model_name,
                messages=messages,
                temperature=temperature,
                max_tokens=max_tokens,
                stream=False
            )
            
            return response.choices[0].message.content
            
        except Exception as e:
            return f"Error generating response: {str(e)}"
    
    def format_citations(self, context: List[Dict[str, Any]]) -> List[str]:
        """Format citations from context documents"""
        citations = []
        for doc in context:
            metadata = doc.get('metadata', {})
            source = metadata.get('source', 'Unknown')
            page = metadata.get('page', '')
            
            if page:
                citation = f"{source}, p. {page}"
            else:
                citation = source
            
            if citation not in citations:
                citations.append(citation)
        
        return citations
    
    def rerank(self, query: str, documents: List[Dict[str, Any]]) -> List[Dict[str, Any]]:
        """
        Rerank documents for relevance using DeepSeek
        This is a simple implementation - can be enhanced with dedicated reranking
        """
        # For now, return documents as-is (already ranked by vector similarity)
        # Could implement LLM-based reranking here if needed
        return documents
    
    def generate_response_with_citations(
        self, 
        query: str,
        history: Optional[List[Dict[str, str]]] = None
    ) -> Dict[str, Any]:
        """Generate response with separate citations"""
        
        # Retrieve and potentially rerank context
        context = self.retrieve_context(query)
        context = self.rerank(query, context)
        
        # Generate response
        response_text = self.generate_response(
            query=query,
            context=context,
            history=history
        )
        
        # Format citations
        citations = self.format_citations(context)
        
        return {
            "response": response_text,
            "citations": citations,
            "context_used": len(context)
        }