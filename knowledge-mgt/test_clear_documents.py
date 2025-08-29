#!/usr/bin/env python
"""
Test script to verify the Clear All Documents functionality
"""
import os
import sys
from pathlib import Path
from backend.vector_store import VectorStore
from backend.config import settings

def test_clear_documents():
    """Test clearing documents from vector store"""
    print("="*60)
    print("Testing Clear All Documents Functionality")
    print("="*60)
    
    # Initialize vector store
    print("\n1. Initializing vector store...")
    vector_store = VectorStore()
    
    # Check initial state
    print("\n2. Checking initial state...")
    stats = vector_store.get_collection_stats()
    print(f"   Documents: {stats['total_documents']}")
    print(f"   Chunks: {stats['total_chunks']}")
    
    # Add some test documents if empty
    if stats['total_chunks'] == 0:
        print("\n3. Adding test documents...")
        test_docs = [
            {
                "id": "test_1",
                "content": "This is test document 1 about company policies.",
                "metadata": {"source": "test1.txt", "document_id": "doc1"}
            },
            {
                "id": "test_2", 
                "content": "This is test document 2 about vacation rules.",
                "metadata": {"source": "test2.txt", "document_id": "doc2"}
            },
            {
                "id": "test_3",
                "content": "This is test document 3 about expense reports.",
                "metadata": {"source": "test3.txt", "document_id": "doc3"}
            }
        ]
        
        for doc in test_docs:
            vector_store.add_documents([doc])
        
        stats = vector_store.get_collection_stats()
        print(f"   Added documents - Now have {stats['total_documents']} documents, {stats['total_chunks']} chunks")
    
    # Test clear functionality
    print("\n4. Testing clear_collection()...")
    success = vector_store.clear_collection()
    
    if success:
        print("   ✅ Clear operation returned success")
    else:
        print("   ❌ Clear operation failed")
    
    # Verify documents are cleared
    print("\n5. Verifying collection is empty...")
    stats = vector_store.get_collection_stats()
    print(f"   Documents after clear: {stats['total_documents']}")
    print(f"   Chunks after clear: {stats['total_chunks']}")
    
    if stats['total_documents'] == 0 and stats['total_chunks'] == 0:
        print("   ✅ Collection successfully cleared!")
    else:
        print("   ❌ Collection still contains data!")
    
    # Test upload directory
    print("\n6. Checking upload directory...")
    upload_dir = Path(settings.upload_path)
    if upload_dir.exists():
        files = list(upload_dir.glob("*"))
        print(f"   Upload directory has {len(files)} files")
        if files:
            print("   Files found:")
            for f in files[:5]:  # Show first 5
                print(f"     - {f.name}")
    else:
        print("   Upload directory doesn't exist")
    
    print("\n" + "="*60)
    print("Test Complete!")
    print("="*60)
    
    return stats['total_documents'] == 0 and stats['total_chunks'] == 0

if __name__ == "__main__":
    success = test_clear_documents()
    sys.exit(0 if success else 1)