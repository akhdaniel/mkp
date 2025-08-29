"""Main Streamlit application for AI Helpdesk Chat."""

import streamlit as st
from pathlib import Path
import time
from datetime import datetime
from typing import Optional

from backend.config import settings
from backend.document_processor import DocumentProcessor
from backend.vector_store import VectorStore
from backend.rag_engine import RAGEngine


# Page configuration
st.set_page_config(
    page_title="AI Helpdesk Chat",
    page_icon="🤖",
    layout="wide",
    initial_sidebar_state="expanded"
)

# Initialize session state
if "messages" not in st.session_state:
    st.session_state.messages = []
if "vector_store" not in st.session_state:
    st.session_state.vector_store = VectorStore()
if "rag_engine" not in st.session_state:
    st.session_state.rag_engine = RAGEngine(st.session_state.vector_store)
if "doc_processor" not in st.session_state:
    st.session_state.doc_processor = DocumentProcessor()
if "current_collection" not in st.session_state:
    st.session_state.current_collection = "documents"


def validate_api_keys():
    """Validate that required API keys are configured."""
    valid, errors = settings.validate_api_keys()
    if not valid:
        st.error("⚠️ Configuration Error")
        for error in errors:
            st.error(f"• {error}")
        st.info("Please set the required API keys in your .env file")
        st.stop()


def display_chat_interface():
    """Display the main chat interface."""
    st.title("🤖 AI Helpdesk Chat")
    st.markdown("Ask questions about your uploaded documents")
    
    # Display chat messages
    for message in st.session_state.messages:
        with st.chat_message(message["role"]):
            st.markdown(message["content"])
            
            # Display sources if available
            if message["role"] == "assistant" and "sources" in message:
                if message["sources"]:
                    with st.expander("📚 Sources", expanded=False):
                        for source in message["sources"]:
                            st.markdown(f"• **{source['title']}** (Score: {source.get('relevance_score', 0):.2f})")
    
    # Chat input
    if prompt := st.chat_input("Ask a question about your documents..."):
        # Add user message to chat
        st.session_state.messages.append({"role": "user", "content": prompt})
        
        # Display user message
        with st.chat_message("user"):
            st.markdown(prompt)
        
        # Generate response
        with st.chat_message("assistant"):
            with st.spinner("Searching documents and generating response..."):
                # Get conversation history
                conversation_history = [
                    {"role": msg["role"], "content": msg["content"]}
                    for msg in st.session_state.messages[-10:]  # Last 10 messages
                ]
                
                # Generate response using RAG
                response = st.session_state.rag_engine.generate_response(
                    query=prompt,
                    conversation_history=conversation_history,
                    use_hybrid_search=True
                )
                
                # Display response
                st.markdown(response["answer"])
                
                # Display sources
                if response["sources"]:
                    with st.expander("📚 Sources", expanded=False):
                        for source in response["sources"]:
                            st.markdown(
                                f"• **{source['title']}** "
                                f"(Score: {source.get('relevance_score', 0):.2f})"
                            )
                
                # Display confidence
                confidence = response["confidence"]
                confidence_color = "green" if confidence > 0.7 else "orange" if confidence > 0.4 else "red"
                st.markdown(
                    f"Confidence: :{confidence_color}[{confidence:.0%}] | "
                    f"Sources: {response['num_sources']} | "
                    f"Time: {response['processing_time']:.2f}s"
                )
                
                # Add assistant message to chat
                st.session_state.messages.append({
                    "role": "assistant",
                    "content": response["answer"],
                    "sources": response["sources"],
                    "confidence": confidence
                })


def display_sidebar():
    """Display the sidebar with document management and settings."""
    with st.sidebar:
        st.title("📁 Document Management")
        
        # File upload section
        st.subheader("Upload Documents")
        uploaded_files = st.file_uploader(
            "Choose files",
            type=["pdf", "docx", "pptx", "txt", "xlsx"],
            accept_multiple_files=True,
            help="Upload documents to add to the knowledge base"
        )
        
        if uploaded_files:
            if st.button("📤 Process Documents", type="primary"):
                with st.spinner("Processing documents..."):
                    process_uploaded_files(uploaded_files)
        
        st.divider()
        
        # Collection stats
        st.subheader("📊 Collection Statistics")
        stats = st.session_state.vector_store.get_collection_stats()
        
        col1, col2 = st.columns(2)
        with col1:
            st.metric("Documents", stats["total_documents"])
        with col2:
            st.metric("Chunks", stats["total_chunks"])
        
        # List documents in collection
        if stats["total_documents"] > 0:
            st.subheader("📄 Documents in Collection")
            try:
                # Get all unique documents from the collection
                all_metadata = st.session_state.vector_store.collection.get(include=["metadatas"])
                
                if all_metadata and all_metadata.get("metadatas"):
                    # Extract unique document sources
                    document_sources = {}
                    for metadata in all_metadata["metadatas"]:
                        if metadata and "source" in metadata:
                            source = metadata["source"]
                            if source not in document_sources:
                                document_sources[source] = {
                                    "chunks": 0,
                                    "upload_time": metadata.get("upload_time", "Unknown"),
                                    "document_id": metadata.get("document_id", "Unknown")
                                }
                            document_sources[source]["chunks"] += 1
                    
                    # Display documents
                    for source, info in document_sources.items():
                        # Extract filename from path
                        filename = Path(source).name
                        with st.expander(f"📄 {filename}", expanded=False):
                            st.text(f"Chunks: {info['chunks']}")
                            st.text(f"Document ID: {info['document_id']}")
                            if st.button(f"🗑️ Remove", key=f"remove_{filename}"):
                                # Remove this document
                                if remove_document_from_collection(source):
                                    st.success(f"Removed {filename}")
                                    st.rerun()
                                else:
                                    st.error(f"Failed to remove {filename}")
            except Exception as e:
                st.error(f"Error listing documents: {str(e)}")
        
        # Clear collection button with confirmation
        if st.button("🗑️ Clear All Documents", type="secondary", use_container_width=True):
            st.session_state.show_clear_confirm = True
        
        # Show confirmation dialog
        if st.session_state.get('show_clear_confirm', False):
            st.warning("⚠️ This will permanently delete all documents from the vector database!")
            col1, col2 = st.columns(2)
            with col1:
                if st.button("✅ Yes, Clear All", type="primary", use_container_width=True):
                    try:
                        # Clear the vector store collection
                        success = st.session_state.vector_store.clear_collection()
                        
                        if success:
                            # Clear uploaded files from disk
                            upload_dir = Path(settings.upload_path)
                            if upload_dir.exists():
                                for file in upload_dir.glob("*"):
                                    try:
                                        file.unlink()
                                        print(f"Deleted file: {file}")
                                    except Exception as e:
                                        print(f"Could not delete {file}: {e}")
                            
                            # Clear session state
                            st.session_state.show_clear_confirm = False
                            if 'uploaded_files' in st.session_state:
                                st.session_state.uploaded_files = []
                            if 'uploaded_documents' in st.session_state:
                                st.session_state.uploaded_documents = []
                            if 'messages' in st.session_state:
                                st.session_state.messages = []
                            
                            st.success("✅ All documents and chat history cleared successfully!")
                            st.balloons()
                            time.sleep(1)
                            st.rerun()
                        else:
                            st.error("Failed to clear documents. Please try again.")
                            st.session_state.show_clear_confirm = False
                    except Exception as e:
                        st.error(f"Error clearing documents: {str(e)}")
                        st.session_state.show_clear_confirm = False
            with col2:
                if st.button("❌ Cancel", use_container_width=True):
                    st.session_state.show_clear_confirm = False
                    st.rerun()
        
        st.divider()
        
        # Settings section
        st.subheader("⚙️ Settings")
        
        # Search settings
        with st.expander("Search Settings", expanded=False):
            st.slider(
                "Top K Results",
                min_value=1,
                max_value=10,
                value=settings.top_k_retrieval,
                key="top_k_setting",
                help="Number of document chunks to retrieve"
            )
            
            st.checkbox(
                "Use Hybrid Search",
                value=True,
                key="hybrid_search",
                help="Combine semantic and keyword search"
            )
        
        # Model settings
        with st.expander("Model Settings", expanded=False):
            st.text_input(
                "Claude Model",
                value=settings.model_name,
                disabled=True,
                help="Claude model being used"
            )
            
            st.text_input(
                "Embedding Model",
                value=settings.embedding_model,
                disabled=True,
                help="Embedding model for document vectorization"
            )
        
        st.divider()
        
        # Chat controls
        st.subheader("💬 Chat Controls")
        
        if st.button("🧹 Clear Chat History"):
            st.session_state.messages = []
            st.success("Chat history cleared!")
            st.rerun()
        
        if st.session_state.messages:
            if st.button("📝 Get Conversation Summary"):
                with st.spinner("Generating summary..."):
                    summary = st.session_state.rag_engine.get_conversation_summary(
                        st.session_state.messages
                    )
                    st.text_area("Summary", summary, height=150)
        
        # Export chat
        if st.session_state.messages:
            if st.button("💾 Export Chat"):
                export_chat_history()


def remove_document_from_collection(source_path):
    """Remove a specific document from the collection by source path."""
    try:
        # Get all documents with this source
        all_data = st.session_state.vector_store.collection.get(
            where={"source": source_path},
            include=["ids"]
        )
        
        if all_data and all_data.get("ids"):
            # Delete all chunks for this document
            st.session_state.vector_store.collection.delete(ids=all_data["ids"])
            
            # Try to delete the physical file
            try:
                file_path = Path(source_path)
                if file_path.exists():
                    file_path.unlink()
            except Exception as e:
                print(f"Could not delete physical file: {e}")
            
            return True
        return False
    except Exception as e:
        print(f"Error removing document: {e}")
        return False


def process_uploaded_files(uploaded_files):
    """Process uploaded files and add to vector store."""
    progress_bar = st.progress(0)
    status_text = st.empty()
    
    successful = 0
    failed = 0
    
    for idx, uploaded_file in enumerate(uploaded_files):
        try:
            # Update progress
            progress = (idx + 1) / len(uploaded_files)
            progress_bar.progress(progress)
            status_text.text(f"Processing {uploaded_file.name}...")
            
            # Validate file
            is_valid, error_msg = st.session_state.doc_processor.validate_file(
                uploaded_file.name,
                uploaded_file.size
            )
            
            if not is_valid:
                st.error(f"❌ {uploaded_file.name}: {error_msg}")
                failed += 1
                continue
            
            # Save file
            file_path = st.session_state.doc_processor.save_uploaded_file(
                uploaded_file.read(),
                uploaded_file.name
            )
            
            # Process document
            document = st.session_state.doc_processor.process_document(
                file_path,
                metadata={
                    "upload_user": "admin",
                    "department": "general",
                    "source": str(file_path),  # Ensure source is set
                    "upload_time": datetime.now().isoformat()
                }
            )
            
            # Add to vector store
            result = st.session_state.vector_store.add_documents([document])
            
            if result["success"]:
                successful += 1
                st.success(f"✅ {uploaded_file.name} processed successfully ({document['num_chunks']} chunks)")
                # Track uploaded files in session state
                if 'uploaded_documents' not in st.session_state:
                    st.session_state.uploaded_documents = []
                st.session_state.uploaded_documents.append({
                    "filename": uploaded_file.name,
                    "path": str(file_path),
                    "chunks": document['num_chunks'],
                    "document_id": document['id']
                })
            else:
                failed += 1
                st.error(f"❌ {uploaded_file.name} failed to index: {result['failed_chunks']} chunks failed")
        
        except Exception as e:
            failed += 1
            st.error(f"❌ Error processing {uploaded_file.name}: {str(e)}")
    
    # Clear progress
    progress_bar.empty()
    status_text.empty()
    
    # Summary
    st.success(f"Processing complete! Success: {successful}, Failed: {failed}")
    
    # Refresh stats
    time.sleep(1)
    st.rerun()


def export_chat_history():
    """Export chat history as markdown."""
    export_text = "# AI Helpdesk Chat History\n\n"
    export_text += f"Exported: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n\n"
    
    for message in st.session_state.messages:
        role = "User" if message["role"] == "user" else "Assistant"
        export_text += f"## {role}\n{message['content']}\n\n"
        
        if message["role"] == "assistant" and "sources" in message:
            if message["sources"]:
                export_text += "**Sources:**\n"
                for source in message["sources"]:
                    export_text += f"- {source['title']}\n"
                export_text += "\n"
    
    st.download_button(
        label="📥 Download Chat History",
        data=export_text,
        file_name=f"chat_history_{datetime.now().strftime('%Y%m%d_%H%M%S')}.md",
        mime="text/markdown"
    )


def main():
    """Main application entry point."""
    # Validate configuration
    validate_api_keys()
    
    # Display sidebar
    display_sidebar()
    
    # Display main chat interface
    display_chat_interface()


if __name__ == "__main__":
    main()