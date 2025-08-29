"""
Multi-LLM RAG Engine supporting Anthropic Claude and DeepSeek models
"""
import os
from typing import List, Dict, Any, Optional, Tuple
from datetime import datetime
import json

from anthropic import Anthropic
from openai import OpenAI
from backend.config import settings
from backend.vector_store import VectorStore


class MultiLLMRAGEngine:
    """RAG Engine with support for multiple LLM providers"""
    
    def __init__(self, vector_store: Optional[VectorStore] = None, provider: str = None):
        """Initialize the Multi-LLM RAG engine.
        
        Args:
            vector_store: Optional vector store instance
            provider: LLM provider ('anthropic' or 'deepseek')
        """
        self.vector_store = vector_store or VectorStore()
        self.top_k = settings.top_k_retrieval
        
        # Determine provider from environment or parameter
        self.provider = provider or os.getenv("LLM_PROVIDER", "anthropic").lower()
        
        # Initialize the appropriate client
        if self.provider == "deepseek":
            self._init_deepseek()
        else:  # Default to Anthropic
            self._init_anthropic()
    
    def _init_anthropic(self):
        """Initialize Anthropic Claude client"""
        self.provider = "anthropic"
        api_key = settings.anthropic_api_key or os.getenv("ANTHROPIC_API_KEY")
        if not api_key:
            raise ValueError("ANTHROPIC_API_KEY not found in environment")
        
        self.client = Anthropic(api_key=api_key)
        self.model_name = settings.model_name
        
    def _init_deepseek(self):
        """Initialize DeepSeek client using OpenAI-compatible interface"""
        self.provider = "deepseek"
        api_key = os.getenv("DEEPSEEK_API_KEY")
        if not api_key:
            raise ValueError("DEEPSEEK_API_KEY not found in environment")
        
        self.client = OpenAI(
            api_key=api_key,
            base_url="https://api.deepseek.com/v1"
        )
        self.model_name = os.getenv("DEEPSEEK_MODEL", "deepseek-chat")
    
    def retrieve_context(
        self,
        query: str,
        filter_metadata: Optional[Dict[str, Any]] = None,
        use_hybrid_search: bool = True
    ) -> List[Dict[str, Any]]:
        """Retrieve relevant documents from vector store"""
        
        if use_hybrid_search and hasattr(self.vector_store, 'hybrid_search'):
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
        
        return search_results
    
    def _format_context_for_prompt(self, documents: List[Dict[str, Any]]) -> str:
        """Format retrieved documents for inclusion in prompt"""
        if not documents:
            return "No relevant documents found."
        
        context_parts = []
        for i, doc in enumerate(documents, 1):
            content = doc.get('content', '')
            metadata = doc.get('metadata', {})
            source = metadata.get('source', 'Unknown')
            page = metadata.get('page', 'N/A')
            
            context_parts.append(
                f"[Document {i}]\n"
                f"Source: {source} (Page {page})\n"
                f"Content: {content}"
            )
        
        return "\n\n---\n\n".join(context_parts)
    
    def _generate_anthropic_response(
        self,
        query: str,
        context: str,
        conversation_history: Optional[List[Dict[str, str]]] = None
    ) -> str:
        """Generate response using Anthropic Claude"""
        
        # Build message history
        messages = []
        
        # Add conversation history if provided
        if conversation_history:
            for msg in conversation_history[-10:]:  # Limit context
                role = msg.get('role', 'user')
                content = msg.get('content', '')
                # Anthropic uses 'user' and 'assistant' roles
                if role in ['user', 'assistant']:
                    messages.append({"role": role, "content": content})
        
        # Add current query with context
        user_message = f"""Based on the following context from our knowledge base, please answer the question.
        
Context:
{context}

Question: {query}

Please provide a comprehensive answer based on the context. If the context doesn't contain relevant information, please state that clearly."""
        
        messages.append({"role": "user", "content": user_message})
        
        # Generate response
        response = self.client.messages.create(
            model=self.model_name,
            messages=messages,
            max_tokens=2000,
            temperature=0.7
        )
        
        return response.content[0].text
    
    def _generate_deepseek_response(
        self,
        query: str,
        context: str,
        conversation_history: Optional[List[Dict[str, str]]] = None
    ) -> str:
        """Generate response using DeepSeek"""
        
        # Build messages for OpenAI-compatible format
        messages = [
            {
                "role": "system",
                "content": "You are a helpful AI assistant for an internal knowledge management system. Answer questions based on the provided context from documents."
            }
        ]
        
        # Add conversation history if provided
        if conversation_history:
            for msg in conversation_history[-10:]:  # Limit context
                role = msg.get('role', 'user')
                content = msg.get('content', '')
                if role in ['user', 'assistant']:
                    messages.append({"role": role, "content": content})
        
        # Add current query with context
        user_message = f"""Context from knowledge base:
{context}

Question: {query}

Please provide a comprehensive answer based on the context above. Include specific references to the documents when answering."""
        
        messages.append({"role": "user", "content": user_message})
        
        # Generate response
        response = self.client.chat.completions.create(
            model=self.model_name,
            messages=messages,
            temperature=0.7,
            max_tokens=2000
        )
        
        return response.choices[0].message.content
    
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
        documents = self.retrieve_context(
            query=query,
            filter_metadata=filter_metadata,
            use_hybrid_search=use_hybrid_search
        )
        
        # Format context
        context = self._format_context_for_prompt(documents)
        
        # Generate response based on provider
        try:
            if self.provider == "deepseek":
                response_text = self._generate_deepseek_response(
                    query=query,
                    context=context,
                    conversation_history=conversation_history
                )
            else:
                response_text = self._generate_anthropic_response(
                    query=query,
                    context=context,
                    conversation_history=conversation_history
                )
        except Exception as e:
            response_text = f"Error generating response: {str(e)}"
        
        # Format citations
        citations = []
        for doc in documents:
            metadata = doc.get('metadata', {})
            source = metadata.get('source', 'Unknown')
            page = metadata.get('page', '')
            
            if page:
                citation = f"{source}, p. {page}"
            else:
                citation = source
            
            if citation not in citations:
                citations.append(citation)
        
        # Calculate response time
        response_time = (datetime.now() - start_time).total_seconds()
        
        return {
            "response": response_text,
            "citations": citations,
            "documents_used": len(documents),
            "provider": self.provider,
            "model": self.model_name,
            "response_time": response_time
        }
    
    def switch_provider(self, provider: str):
        """Switch between LLM providers
        
        Args:
            provider: 'anthropic' or 'deepseek'
        """
        if provider.lower() == "deepseek":
            self._init_deepseek()
        elif provider.lower() == "anthropic":
            self._init_anthropic()
        else:
            raise ValueError(f"Unsupported provider: {provider}. Use 'anthropic' or 'deepseek'")
    
    def get_provider_info(self) -> Dict[str, str]:
        """Get information about the current provider"""
        return {
            "provider": self.provider,
            "model": self.model_name,
            "status": "active"
        }