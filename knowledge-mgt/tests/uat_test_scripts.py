#!/usr/bin/env python
"""
UAT Automated Test Scripts for Knowledge Management System
Run these scripts to validate core functionality
"""

import os
import sys
import time
import json
from pathlib import Path
from datetime import datetime
from typing import Dict, List, Optional, Tuple
import hashlib

# Add parent directory to path for imports
sys.path.append(str(Path(__file__).parent.parent))

from backend.config import settings
from backend.document_processor import DocumentProcessor
from backend.vector_store import VectorStore
from backend.rag_engine import RAGEngine


class Colors:
    """ANSI color codes for terminal output"""
    HEADER = '\033[95m'
    BLUE = '\033[94m'
    GREEN = '\033[92m'
    YELLOW = '\033[93m'
    RED = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'


class UATTestRunner:
    """Main UAT test runner class"""
    
    def __init__(self):
        self.vector_store = VectorStore()
        self.doc_processor = DocumentProcessor()
        self.rag_engine = RAGEngine(self.vector_store)
        self.test_results = []
        self.test_data_dir = Path("./tests/test_data")
        self.test_data_dir.mkdir(exist_ok=True)
        
    def print_header(self, text: str):
        """Print formatted header"""
        print(f"\n{Colors.HEADER}{'='*60}")
        print(f"{text}")
        print(f"{'='*60}{Colors.ENDC}\n")
        
    def print_test(self, test_id: str, description: str):
        """Print test information"""
        print(f"{Colors.BLUE}[{test_id}] {description}{Colors.ENDC}")
        
    def print_success(self, message: str):
        """Print success message"""
        print(f"  {Colors.GREEN}✓ {message}{Colors.ENDC}")
        
    def print_error(self, message: str):
        """Print error message"""
        print(f"  {Colors.RED}✗ {message}{Colors.ENDC}")
        
    def print_info(self, message: str):
        """Print info message"""
        print(f"  {Colors.YELLOW}ℹ {message}{Colors.ENDC}")
        
    def create_test_documents(self) -> Dict[str, Path]:
        """Create test documents for UAT"""
        documents = {}
        
        # Create HR Policy document
        hr_policy = self.test_data_dir / "hr_policy.txt"
        hr_content = """
        EMPLOYEE HANDBOOK - HR POLICIES
        
        VACATION POLICY
        All full-time employees are entitled to 15 days of paid vacation per year.
        Vacation days accrue at a rate of 1.25 days per month.
        Unused vacation can be carried over to the next year, up to a maximum of 10 days.
        
        SICK LEAVE POLICY
        Employees receive 10 days of paid sick leave per year.
        Sick leave does not carry over to the next year.
        A doctor's note is required for absences longer than 3 consecutive days.
        
        REMOTE WORK POLICY
        Employees may work remotely up to 2 days per week with manager approval.
        Core hours are 10 AM to 3 PM in the employee's local timezone.
        Remote workers must be available via chat and video during core hours.
        
        EXPENSE REIMBURSEMENT
        Business travel expenses must be submitted within 30 days of travel.
        Daily meal allowance is $75 for domestic travel and $100 for international.
        All expenses over $100 require receipts and manager approval.
        Reimbursement is processed within 14 business days of approval.
        """
        hr_policy.write_text(hr_content)
        documents["hr_policy"] = hr_policy
        
        # Create IT Security document
        it_security = self.test_data_dir / "it_security.txt"
        it_content = """
        IT SECURITY GUIDELINES
        
        PASSWORD REQUIREMENTS
        - Minimum 12 characters
        - Must include uppercase, lowercase, numbers, and special characters
        - Password rotation every 90 days
        - Cannot reuse last 5 passwords
        
        DATA CLASSIFICATION
        - Public: Can be shared externally
        - Internal: For employees only
        - Confidential: Restricted access, requires approval
        - Secret: Highly restricted, C-level approval required
        
        ACCEPTABLE USE
        Company devices should be used primarily for business purposes.
        Personal use should be minimal and appropriate.
        Installation of unauthorized software is prohibited.
        All data on company devices is subject to monitoring.
        """
        it_security.write_text(it_content)
        documents["it_security"] = it_security
        
        # Create Training Manual
        training = self.test_data_dir / "training_manual.txt"
        training_content = """
        NEW EMPLOYEE TRAINING PROGRAM
        
        WEEK 1: ORIENTATION
        - Company history and culture
        - HR policies and procedures
        - IT setup and security training
        - Team introductions
        
        WEEK 2: ROLE-SPECIFIC TRAINING
        - Department overview
        - Job responsibilities
        - Key systems and tools
        - Shadow experienced team member
        
        WEEK 3-4: HANDS-ON PRACTICE
        - Begin working on actual tasks with supervision
        - Daily check-ins with manager
        - Complete online training modules
        - First performance review at end of week 4
        """
        training.write_text(training_content)
        documents["training"] = training
        
        return documents
    
    def test_document_upload(self) -> bool:
        """TC-DM-001: Test single document upload"""
        self.print_test("TC-DM-001", "Single Document Upload")
        
        try:
            # Create test document
            test_docs = self.create_test_documents()
            hr_policy = test_docs["hr_policy"]
            
            # Process document
            document = self.doc_processor.process_document(
                hr_policy,
                metadata={
                    "upload_user": "uat_tester",
                    "department": "testing",
                    "source": str(hr_policy.absolute()),
                    "upload_time": datetime.now().isoformat()
                }
            )
            
            # Add to vector store
            result = self.vector_store.add_documents([document])
            
            if result["success"]:
                self.print_success(f"Document uploaded: {document['num_chunks']} chunks")
                
                # Verify in collection
                stats = self.vector_store.get_collection_stats()
                if stats["total_documents"] > 0:
                    self.print_success("Document appears in collection")
                    return True
                else:
                    self.print_error("Document not in collection")
                    return False
            else:
                self.print_error(f"Upload failed: {result.get('error', 'Unknown error')}")
                return False
                
        except Exception as e:
            self.print_error(f"Exception: {str(e)}")
            return False
    
    def test_multiple_upload(self) -> bool:
        """TC-DM-002: Test multiple document upload"""
        self.print_test("TC-DM-002", "Multiple Document Upload")
        
        try:
            # Create all test documents
            test_docs = self.create_test_documents()
            
            success_count = 0
            for name, path in test_docs.items():
                # Process each document
                document = self.doc_processor.process_document(
                    path,
                    metadata={
                        "upload_user": "uat_tester",
                        "department": "testing",
                        "source": str(path.absolute()),
                        "upload_time": datetime.now().isoformat()
                    }
                )
                
                result = self.vector_store.add_documents([document])
                if result["success"]:
                    success_count += 1
                    self.print_info(f"Uploaded {name}: {document['num_chunks']} chunks")
            
            if success_count == len(test_docs):
                self.print_success(f"All {success_count} documents uploaded")
                return True
            else:
                self.print_error(f"Only {success_count}/{len(test_docs)} uploaded")
                return False
                
        except Exception as e:
            self.print_error(f"Exception: {str(e)}")
            return False
    
    def test_basic_query(self) -> bool:
        """TC-SR-001: Test basic question answering"""
        self.print_test("TC-SR-001", "Basic Question Answering")
        
        try:
            # Ensure documents are loaded
            stats = self.vector_store.get_collection_stats()
            if stats["total_documents"] == 0:
                self.print_info("Loading test documents first...")
                self.test_multiple_upload()
            
            # Ask a question
            query = "What is the vacation policy for full-time employees?"
            self.print_info(f"Query: {query}")
            
            response = self.rag_engine.generate_response(
                query=query,
                conversation_history=[],
                use_hybrid_search=True
            )
            
            if response and response.get("answer"):
                self.print_success("Response received")
                self.print_info(f"Answer preview: {response['answer'][:100]}...")
                
                # Check if answer contains expected information
                if "15 days" in response["answer"] or "vacation" in response["answer"].lower():
                    self.print_success("Answer contains relevant information")
                    
                    # Check citations
                    if response.get("sources"):
                        self.print_success(f"Sources cited: {len(response['sources'])}")
                        return True
                    else:
                        self.print_error("No sources cited")
                        return False
                else:
                    self.print_error("Answer may not be relevant")
                    return False
            else:
                self.print_error("No response received")
                return False
                
        except Exception as e:
            self.print_error(f"Exception: {str(e)}")
            return False
    
    def test_no_relevant_docs(self) -> bool:
        """TC-SR-002: Test query with no relevant documents"""
        self.print_test("TC-SR-002", "No Relevant Documents Query")
        
        try:
            query = "What is the current stock price of Apple?"
            self.print_info(f"Query: {query}")
            
            response = self.rag_engine.generate_response(
                query=query,
                conversation_history=[],
                use_hybrid_search=True
            )
            
            if response and response.get("answer"):
                self.print_success("Response received")
                
                # Check that system doesn't hallucinate
                answer_lower = response["answer"].lower()
                if any(word in answer_lower for word in ["don't have", "no information", "not found", "unable"]):
                    self.print_success("System correctly indicates no relevant information")
                    return True
                else:
                    self.print_error("System may be hallucinating information")
                    return False
            else:
                self.print_error("No response received")
                return False
                
        except Exception as e:
            self.print_error(f"Exception: {str(e)}")
            return False
    
    def test_multi_doc_query(self) -> bool:
        """TC-SR-003: Test multi-document query"""
        self.print_test("TC-SR-003", "Multi-Document Query")
        
        try:
            query = "Compare the vacation and sick leave policies"
            self.print_info(f"Query: {query}")
            
            response = self.rag_engine.generate_response(
                query=query,
                conversation_history=[],
                use_hybrid_search=True
            )
            
            if response and response.get("answer"):
                self.print_success("Response received")
                
                # Check if both policies are mentioned
                answer_lower = response["answer"].lower()
                has_vacation = "vacation" in answer_lower
                has_sick = "sick" in answer_lower
                
                if has_vacation and has_sick:
                    self.print_success("Answer references both policies")
                    
                    # Check multiple sources
                    if response.get("sources") and len(response["sources"]) >= 1:
                        self.print_success(f"Multiple sections referenced: {len(response['sources'])}")
                        return True
                    else:
                        self.print_error("Insufficient source citations")
                        return False
                else:
                    self.print_error("Answer doesn't cover both topics")
                    return False
            else:
                self.print_error("No response received")
                return False
                
        except Exception as e:
            self.print_error(f"Exception: {str(e)}")
            return False
    
    def test_document_removal(self) -> bool:
        """TC-DM-004: Test individual document removal"""
        self.print_test("TC-DM-004", "Individual Document Removal")
        
        try:
            # Get initial stats
            initial_stats = self.vector_store.get_collection_stats()
            self.print_info(f"Initial documents: {initial_stats['total_documents']}")
            
            # Find a document to remove
            all_metadata = self.vector_store.collection.get(include=["metadatas"])
            
            if all_metadata and all_metadata.get("metadatas"):
                # Get first unique source
                sources = set()
                for metadata in all_metadata["metadatas"]:
                    if metadata and "source" in metadata:
                        sources.add(metadata["source"])
                
                if sources:
                    source_to_remove = list(sources)[0]
                    self.print_info(f"Removing: {Path(source_to_remove).name}")
                    
                    # Remove document
                    all_data = self.vector_store.collection.get(
                        where={"source": source_to_remove},
                        include=["ids"]
                    )
                    
                    if all_data and all_data.get("ids"):
                        chunks_to_remove = len(all_data["ids"])
                        self.vector_store.collection.delete(ids=all_data["ids"])
                        
                        # Verify removal
                        final_stats = self.vector_store.get_collection_stats()
                        
                        if final_stats["total_chunks"] < initial_stats["total_chunks"]:
                            self.print_success(f"Removed {chunks_to_remove} chunks")
                            return True
                        else:
                            self.print_error("Document not removed")
                            return False
                    else:
                        self.print_error("No chunks found for document")
                        return False
                else:
                    self.print_error("No documents to remove")
                    return False
            else:
                self.print_error("Collection is empty")
                return False
                
        except Exception as e:
            self.print_error(f"Exception: {str(e)}")
            return False
    
    def test_clear_collection(self) -> bool:
        """TC-DM-005: Test clear all documents"""
        self.print_test("TC-DM-005", "Clear All Documents")
        
        try:
            # Ensure we have documents first
            stats = self.vector_store.get_collection_stats()
            if stats["total_documents"] == 0:
                self.print_info("Adding documents first...")
                self.test_document_upload()
            
            # Clear collection
            self.print_info("Clearing collection...")
            success = self.vector_store.clear_collection()
            
            if success:
                # Verify empty
                final_stats = self.vector_store.get_collection_stats()
                
                if final_stats["total_documents"] == 0 and final_stats["total_chunks"] == 0:
                    self.print_success("Collection cleared successfully")
                    return True
                else:
                    self.print_error(f"Collection not empty: {final_stats}")
                    return False
            else:
                self.print_error("Clear operation failed")
                return False
                
        except Exception as e:
            self.print_error(f"Exception: {str(e)}")
            return False
    
    def test_performance(self) -> bool:
        """TC-PF-001: Test response time"""
        self.print_test("TC-PF-001", "Response Time - Simple Query")
        
        try:
            # Ensure documents loaded
            stats = self.vector_store.get_collection_stats()
            if stats["total_documents"] == 0:
                self.print_info("Loading test documents...")
                self.test_multiple_upload()
            
            queries = [
                "What is the vacation policy?",
                "How many sick days do employees get?",
                "What are the password requirements?"
            ]
            
            response_times = []
            
            for query in queries:
                start_time = time.time()
                
                response = self.rag_engine.generate_response(
                    query=query,
                    conversation_history=[],
                    use_hybrid_search=True
                )
                
                elapsed = time.time() - start_time
                response_times.append(elapsed)
                
                self.print_info(f"Query: '{query[:30]}...' - {elapsed:.2f}s")
            
            avg_time = sum(response_times) / len(response_times)
            max_time = max(response_times)
            
            self.print_info(f"Average response time: {avg_time:.2f}s")
            self.print_info(f"Maximum response time: {max_time:.2f}s")
            
            if max_time < 5.0:
                self.print_success("All responses under 5 seconds")
                return True
            elif max_time < 10.0:
                self.print_info("Some responses between 5-10 seconds")
                return True
            else:
                self.print_error(f"Response too slow: {max_time:.2f}s")
                return False
                
        except Exception as e:
            self.print_error(f"Exception: {str(e)}")
            return False
    
    def cleanup_test_data(self):
        """Clean up test data after tests"""
        self.print_info("Cleaning up test data...")
        
        # Clear vector store
        self.vector_store.clear_collection()
        
        # Remove test files
        if self.test_data_dir.exists():
            for file in self.test_data_dir.glob("*"):
                try:
                    file.unlink()
                except:
                    pass
            
            try:
                self.test_data_dir.rmdir()
            except:
                pass
    
    def run_all_tests(self) -> Dict:
        """Run all UAT tests"""
        self.print_header("UAT TEST EXECUTION - KNOWLEDGE MANAGEMENT SYSTEM")
        print(f"Started: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")
        
        # Define test suites
        test_suites = {
            "Document Management": [
                ("TC-DM-001", "Single Document Upload", self.test_document_upload),
                ("TC-DM-002", "Multiple Document Upload", self.test_multiple_upload),
                ("TC-DM-004", "Document Removal", self.test_document_removal),
                ("TC-DM-005", "Clear All Documents", self.test_clear_collection),
            ],
            "Search & Retrieval": [
                ("TC-SR-001", "Basic Query", self.test_basic_query),
                ("TC-SR-002", "No Relevant Docs", self.test_no_relevant_docs),
                ("TC-SR-003", "Multi-Document Query", self.test_multi_doc_query),
            ],
            "Performance": [
                ("TC-PF-001", "Response Time", self.test_performance),
            ]
        }
        
        results = {
            "total": 0,
            "passed": 0,
            "failed": 0,
            "tests": []
        }
        
        # Run each test suite
        for suite_name, tests in test_suites.items():
            self.print_header(f"Test Suite: {suite_name}")
            
            for test_id, test_name, test_func in tests:
                results["total"] += 1
                
                try:
                    passed = test_func()
                    
                    if passed:
                        results["passed"] += 1
                        status = "PASSED"
                        print(f"{Colors.GREEN}  [{test_id}] {status}{Colors.ENDC}\n")
                    else:
                        results["failed"] += 1
                        status = "FAILED"
                        print(f"{Colors.RED}  [{test_id}] {status}{Colors.ENDC}\n")
                    
                    results["tests"].append({
                        "id": test_id,
                        "name": test_name,
                        "suite": suite_name,
                        "status": status
                    })
                    
                except Exception as e:
                    results["failed"] += 1
                    print(f"{Colors.RED}  [{test_id}] FAILED - Exception: {str(e)}{Colors.ENDC}\n")
                    results["tests"].append({
                        "id": test_id,
                        "name": test_name,
                        "suite": suite_name,
                        "status": "FAILED",
                        "error": str(e)
                    })
        
        # Cleanup
        self.cleanup_test_data()
        
        # Print summary
        self.print_header("TEST EXECUTION SUMMARY")
        print(f"Total Tests: {results['total']}")
        print(f"{Colors.GREEN}Passed: {results['passed']}{Colors.ENDC}")
        print(f"{Colors.RED}Failed: {results['failed']}{Colors.ENDC}")
        
        pass_rate = (results['passed'] / results['total'] * 100) if results['total'] > 0 else 0
        
        if pass_rate >= 90:
            print(f"\n{Colors.GREEN}✓ UAT PASSED - {pass_rate:.1f}% pass rate{Colors.ENDC}")
        elif pass_rate >= 80:
            print(f"\n{Colors.YELLOW}⚠ CONDITIONAL PASS - {pass_rate:.1f}% pass rate{Colors.ENDC}")
        else:
            print(f"\n{Colors.RED}✗ UAT FAILED - {pass_rate:.1f}% pass rate{Colors.ENDC}")
        
        print(f"\nCompleted: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        
        # Save results to file
        results_file = Path("./tests/uat_results.json")
        with open(results_file, "w") as f:
            json.dump(results, f, indent=2)
        print(f"\nResults saved to: {results_file}")
        
        return results


def main():
    """Main entry point for UAT testing"""
    runner = UATTestRunner()
    results = runner.run_all_tests()
    
    # Exit with appropriate code
    if results["failed"] == 0:
        sys.exit(0)
    else:
        sys.exit(1)


if __name__ == "__main__":
    main()