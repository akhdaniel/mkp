#!/usr/bin/env python
"""
Test script to verify document listing functionality
"""
import os
import sys
from pathlib import Path
from datetime import datetime
from backend.vector_store import VectorStore
from backend.document_processor import DocumentProcessor
from backend.config import settings

def test_document_listing():
    """Test that documents show up in the collection after upload"""
    print("="*60)
    print("Testing Document Listing Functionality")
    print("="*60)
    
    # Initialize components
    print("\n1. Initializing components...")
    vector_store = VectorStore()
    doc_processor = DocumentProcessor()
    
    # Check initial state
    print("\n2. Checking initial state...")
    stats = vector_store.get_collection_stats()
    print(f"   Initial documents: {stats['total_documents']}")
    print(f"   Initial chunks: {stats['total_chunks']}")
    
    # Create a test document
    print("\n3. Creating test document...")
    test_file = Path("./test_document.txt")
    test_content = """This is a test document for the knowledge management system.
    It contains multiple paragraphs to test chunking.
    
    This is the second paragraph with more content.
    The document should be properly indexed and searchable.
    
    This is the third paragraph with additional information.
    We want to ensure the document appears in the collection listing."""
    
    test_file.write_text(test_content)
    print(f"   Created {test_file}")
    
    # Process the document
    print("\n4. Processing document...")
    try:
        document = doc_processor.process_document(
            test_file,
            metadata={
                "upload_user": "test_user",
                "department": "testing",
                "source": str(test_file.absolute()),
                "upload_time": datetime.now().isoformat()
            }
        )
        print(f"   Document processed: {document['num_chunks']} chunks created")
        print(f"   Document ID: {document['id']}")
        
        # Add to vector store
        print("\n5. Adding to vector store...")
        result = vector_store.add_documents([document])
        
        if result["success"]:
            print(f"   ✅ Successfully added {result['total_chunks']} chunks")
        else:
            print(f"   ❌ Failed to add document: {result['failed_chunks']} chunks failed")
            
    except Exception as e:
        print(f"   ❌ Error processing document: {e}")
        return False
    
    # Check if document appears in collection
    print("\n6. Checking collection for document...")
    stats = vector_store.get_collection_stats()
    print(f"   Total documents: {stats['total_documents']}")
    print(f"   Total chunks: {stats['total_chunks']}")
    
    # List all documents
    print("\n7. Listing all documents in collection...")
    try:
        all_metadata = vector_store.collection.get(include=["metadatas"])
        
        if all_metadata and all_metadata.get("metadatas"):
            # Extract unique documents
            document_sources = {}
            for metadata in all_metadata["metadatas"]:
                if metadata and "source" in metadata:
                    source = metadata["source"]
                    if source not in document_sources:
                        document_sources[source] = {
                            "chunks": 0,
                            "document_id": metadata.get("document_id", "Unknown")
                        }
                    document_sources[source]["chunks"] += 1
            
            print(f"   Found {len(document_sources)} unique documents:")
            for source, info in document_sources.items():
                filename = Path(source).name
                print(f"     - {filename}: {info['chunks']} chunks (ID: {info['document_id']})")
                
            # Check if our test document is there
            test_file_found = any(
                "test_document.txt" in source 
                for source in document_sources.keys()
            )
            
            if test_file_found:
                print("\n   ✅ Test document appears in collection listing!")
            else:
                print("\n   ❌ Test document NOT found in collection listing!")
                print("   Available sources:", list(document_sources.keys()))
                
        else:
            print("   ❌ No documents found in collection")
            
    except Exception as e:
        print(f"   ❌ Error listing documents: {e}")
        return False
    
    # Clean up
    print("\n8. Cleaning up...")
    try:
        test_file.unlink()
        print(f"   Deleted {test_file}")
    except:
        pass
    
    print("\n" + "="*60)
    print("Test Complete!")
    print("="*60)
    
    return test_file_found

if __name__ == "__main__":
    success = test_document_listing()
    sys.exit(0 if success else 1)