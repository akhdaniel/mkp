# User Acceptance Testing (UAT) Test Cases

## Document Information
- **System**: Knowledge Management System (KMS)
- **Version**: 1.0
- **Date**: December 2024
- **Test Environment**: UAT Environment
- **Test Data**: Sample organizational documents

---

## Test Case Summary

| Test Suite | Test Cases | Priority | Estimated Time |
|------------|------------|----------|----------------|
| Document Management | 8 | Critical | 45 minutes |
| Search & Retrieval | 6 | Critical | 30 minutes |
| Chat Interface | 5 | High | 25 minutes |
| Multi-Format Support | 5 | High | 30 minutes |
| Performance | 4 | Medium | 20 minutes |
| User Experience | 4 | Medium | 20 minutes |
| **Total** | **32** | - | **170 minutes** |

---

## Test Suite 1: Document Management

### TC-DM-001: Single Document Upload
**Priority**: Critical  
**Preconditions**: User has access to system, test PDF document available

**Test Steps**:
1. Navigate to the KMS application
2. Click "Browse files" in the sidebar
3. Select a PDF document (< 10MB)
4. Click "Process Documents"
5. Wait for processing to complete

**Expected Results**:
- File upload progress bar displays
- Success message shows with chunk count
- Document appears in "Documents in Collection" list
- Collection statistics update correctly

**Pass Criteria**: Document successfully uploaded and listed

---

### TC-DM-002: Multiple Document Upload
**Priority**: Critical  
**Preconditions**: Multiple test documents available (PDF, DOCX, PPTX)

**Test Steps**:
1. Click "Browse files" in sidebar
2. Select 3-5 documents using Ctrl/Cmd+Click
3. Click "Process Documents"
4. Monitor processing progress

**Expected Results**:
- All files process sequentially
- Individual success/failure messages display
- All successful documents appear in collection
- Statistics show correct totals

**Pass Criteria**: All valid documents uploaded successfully

---

### TC-DM-003: Invalid File Upload
**Priority**: High  
**Preconditions**: Invalid file types (.exe, .zip) and oversized file (>10MB)

**Test Steps**:
1. Attempt to upload .exe file
2. Attempt to upload file > 10MB
3. Attempt to upload corrupted PDF

**Expected Results**:
- System rejects invalid file types
- Error message for oversized files
- Clear error messages displayed
- No partial uploads in collection

**Pass Criteria**: All invalid files rejected with appropriate errors

---

### TC-DM-004: Individual Document Removal
**Priority**: Critical  
**Preconditions**: At least 3 documents in collection

**Test Steps**:
1. Expand "Documents in Collection"
2. Select a document to remove
3. Click the 🗑️ Remove button
4. Verify removal

**Expected Results**:
- Document removed from list immediately
- Collection statistics update
- Related chunks removed from vector store
- Other documents remain unaffected

**Pass Criteria**: Selected document removed successfully

---

### TC-DM-005: Clear All Documents
**Priority**: Critical  
**Preconditions**: Multiple documents in collection

**Test Steps**:
1. Click "Clear All Documents" button
2. Read warning message
3. Click "Yes, Clear All"
4. Verify collection is empty

**Expected Results**:
- Confirmation dialog appears
- All documents removed after confirmation
- Statistics show 0 documents, 0 chunks
- Success message with celebration animation

**Pass Criteria**: All documents cleared successfully

---

### TC-DM-006: Clear All Documents - Cancel
**Priority**: Medium  
**Preconditions**: Documents in collection

**Test Steps**:
1. Click "Clear All Documents"
2. Click "Cancel" in confirmation dialog
3. Verify documents remain

**Expected Results**:
- Confirmation dialog closes
- No documents removed
- Collection remains unchanged

**Pass Criteria**: Cancel operation works correctly

---

### TC-DM-007: Document Update/Replace
**Priority**: High  
**Preconditions**: Outdated document in collection, updated version available

**Test Steps**:
1. Remove old document version
2. Upload new document version
3. Test with relevant query
4. Verify updated information returned

**Expected Results**:
- Old document successfully removed
- New document successfully uploaded
- Queries return updated information
- Citations reference new version

**Pass Criteria**: Document update process works correctly

---

### TC-DM-008: Duplicate Document Handling
**Priority**: Medium  
**Preconditions**: Document already in collection

**Test Steps**:
1. Upload a document
2. Upload the same document again
3. Check collection listing

**Expected Results**:
- System processes duplicate
- Both instances may appear (or single updated)
- No system errors
- Clear indication of duplicates

**Pass Criteria**: Duplicate handling works without errors

---

## Test Suite 2: Search & Retrieval

### TC-SR-001: Basic Question Answering
**Priority**: Critical  
**Preconditions**: HR policy document uploaded

**Test Steps**:
1. Type "What is the vacation policy?"
2. Press Enter or click Send
3. Review response
4. Check source citations

**Expected Results**:
- Response generated within 10 seconds
- Answer is relevant and accurate
- Sources cited with relevance scores
- Confidence indicator displayed

**Pass Criteria**: Accurate answer with proper citations

---

### TC-SR-002: No Relevant Documents
**Priority**: High  
**Preconditions**: Limited documents uploaded

**Test Steps**:
1. Ask about topic not in documents
2. Example: "What is the company stock price?"
3. Review response

**Expected Results**:
- System indicates no relevant information found
- No hallucinated information
- Suggestion to upload relevant documents
- Professional response tone

**Pass Criteria**: Appropriate handling of no-match queries

---

### TC-SR-003: Multi-Document Query
**Priority**: High  
**Preconditions**: Multiple related documents uploaded

**Test Steps**:
1. Ask question spanning multiple documents
2. Example: "Compare vacation and sick leave policies"
3. Review comprehensive response

**Expected Results**:
- Information synthesized from multiple sources
- All relevant documents cited
- Coherent combined response
- Clear source attribution

**Pass Criteria**: Successfully combines information from multiple documents

---

### TC-SR-004: Follow-up Questions
**Priority**: Critical  
**Preconditions**: Initial question answered

**Test Steps**:
1. Ask initial question about vacation policy
2. Ask follow-up: "How do I apply for it?"
3. Ask another: "What's the approval process?"

**Expected Results**:
- Context maintained across questions
- Follow-ups understood correctly
- Relevant answers to each question
- Conversation flows naturally

**Pass Criteria**: Context preserved in conversation

---

### TC-SR-005: Complex Query Handling
**Priority**: Medium  
**Preconditions**: Technical documents uploaded

**Test Steps**:
1. Ask complex multi-part question
2. Example: "What are the expense limits for travel, and how long does reimbursement take?"
3. Review structured response

**Expected Results**:
- Both parts of question answered
- Response well-structured
- All relevant information included
- Clear formatting

**Pass Criteria**: Complex queries handled appropriately

---

### TC-SR-006: Search Settings Adjustment
**Priority**: Medium  
**Preconditions**: Documents uploaded, settings accessible

**Test Steps**:
1. Expand "Search Settings"
2. Adjust "Top K Results" from 5 to 10
3. Ask a broad question
4. Verify more sources retrieved

**Expected Results**:
- Settings adjustment saved
- More document chunks retrieved
- Response potentially more comprehensive
- Performance acceptable

**Pass Criteria**: Search settings affect retrieval as expected

---

## Test Suite 3: Chat Interface

### TC-CI-001: Message Display
**Priority**: High  
**Preconditions**: Clean session

**Test Steps**:
1. Type and send a message
2. Receive AI response
3. Send another message
4. Scroll through history

**Expected Results**:
- User messages right-aligned (blue)
- AI responses left-aligned (gray)
- Timestamps visible
- Smooth scrolling

**Pass Criteria**: Messages display correctly with proper formatting

---

### TC-CI-002: Clear Chat History
**Priority**: Medium  
**Preconditions**: Active conversation with multiple messages

**Test Steps**:
1. Click "Clear Chat History" button
2. Confirm action if prompted
3. Verify chat area is empty

**Expected Results**:
- All messages removed
- Success notification
- Chat ready for new conversation
- Document collection unaffected

**Pass Criteria**: Chat history cleared successfully

---

### TC-CI-003: Export Chat History
**Priority**: Medium  
**Preconditions**: Conversation with 5+ messages

**Test Steps**:
1. Click "Export Chat" button
2. Save file when prompted
3. Open exported file
4. Verify content completeness

**Expected Results**:
- File downloads successfully
- Markdown format preserved
- All messages included
- Sources and timestamps preserved

**Pass Criteria**: Complete chat export with formatting

---

### TC-CI-004: Conversation Summary
**Priority**: Low  
**Preconditions**: Long conversation (10+ messages)

**Test Steps**:
1. Click "Get Conversation Summary"
2. Wait for generation
3. Review summary accuracy

**Expected Results**:
- Summary generated within 15 seconds
- Key points captured
- Concise and accurate
- Displayed in text area

**Pass Criteria**: Accurate conversation summary generated

---

### TC-CI-005: Input Validation
**Priority**: Medium  
**Preconditions**: System ready

**Test Steps**:
1. Try sending empty message
2. Type very long message (1000+ chars)
3. Use special characters
4. Paste formatted text

**Expected Results**:
- Empty messages not sent
- Long messages handled appropriately
- Special characters preserved
- Formatting stripped appropriately

**Pass Criteria**: Input handling works correctly

---

## Test Suite 4: Multi-Format Support

### TC-MF-001: PDF Document Processing
**Priority**: Critical  
**Preconditions**: Sample PDF with text, images, tables

**Test Steps**:
1. Upload complex PDF document
2. Process document
3. Query about PDF content
4. Verify text extraction quality

**Expected Results**:
- PDF uploads successfully
- Text extracted accurately
- Tables preserved in some form
- Queries answered correctly

**Pass Criteria**: PDF content searchable and retrievable

---

### TC-MF-002: Word Document Processing
**Priority**: Critical  
**Preconditions**: .docx file with formatting

**Test Steps**:
1. Upload Word document
2. Process and index
3. Query document content
4. Verify formatting preserved

**Expected Results**:
- DOCX uploads successfully
- Content extracted completely
- Headers/sections recognized
- Accurate responses to queries

**Pass Criteria**: Word documents fully supported

---

### TC-MF-003: PowerPoint Processing
**Priority**: High  
**Preconditions**: .pptx presentation file

**Test Steps**:
1. Upload PowerPoint file
2. Process presentation
3. Query slide content
4. Verify slide text extracted

**Expected Results**:
- PPTX uploads successfully
- Slide content extracted
- Speaker notes included
- Content searchable

**Pass Criteria**: PowerPoint content accessible

---

### TC-MF-004: Excel Processing
**Priority**: Medium  
**Preconditions**: .xlsx file with data

**Test Steps**:
1. Upload Excel file
2. Process spreadsheet
3. Query about data
4. Verify data extraction

**Expected Results**:
- XLSX uploads successfully
- Data extracted from sheets
- Column headers preserved
- Data queryable

**Pass Criteria**: Excel data accessible via queries

---

### TC-MF-005: Plain Text Processing
**Priority**: High  
**Preconditions**: .txt file

**Test Steps**:
1. Upload text file
2. Process document
3. Query content
4. Verify complete extraction

**Expected Results**:
- TXT uploads successfully
- Full content preserved
- Line breaks maintained
- Content searchable

**Pass Criteria**: Text files fully supported

---

## Test Suite 5: Performance Testing

### TC-PF-001: Response Time - Simple Query
**Priority**: High  
**Preconditions**: 10+ documents uploaded

**Test Steps**:
1. Ask simple factual question
2. Measure response time
3. Repeat with different questions

**Expected Results**:
- Response within 5 seconds
- Consistent performance
- No timeout errors
- Smooth UI during processing

**Pass Criteria**: Response time < 5 seconds

---

### TC-PF-002: Response Time - Complex Query
**Priority**: Medium  
**Preconditions**: 50+ documents uploaded

**Test Steps**:
1. Ask complex analytical question
2. Measure response time
3. Monitor system responsiveness

**Expected Results**:
- Response within 15 seconds
- Progress indicator shown
- UI remains responsive
- No crashes or freezes

**Pass Criteria**: Response time < 15 seconds

---

### TC-PF-003: Large Collection Performance
**Priority**: Medium  
**Preconditions**: 100+ documents uploaded

**Test Steps**:
1. Navigate through document list
2. Perform searches
3. Upload additional documents
4. Test system responsiveness

**Expected Results**:
- Document list loads quickly
- Scrolling remains smooth
- Searches complete reasonably
- No significant degradation

**Pass Criteria**: Acceptable performance with large collection

---

### TC-PF-004: Concurrent Operations
**Priority**: Low  
**Preconditions**: Documents uploaded

**Test Steps**:
1. Start document upload
2. While uploading, ask a question
3. Try to remove a document
4. Monitor system behavior

**Expected Results**:
- Operations queue appropriately
- No data corruption
- Clear status indicators
- Graceful handling

**Pass Criteria**: System handles concurrent operations

---

## Test Suite 6: User Experience

### TC-UX-001: First-Time User Flow
**Priority**: High  
**Preconditions**: Fresh system, new user

**Test Steps**:
1. Access system for first time
2. Follow UI prompts
3. Upload first document
4. Ask first question

**Expected Results**:
- Clear onboarding flow
- Intuitive interface
- Helpful prompts/tooltips
- Successful first interaction

**Pass Criteria**: New user can use system without training

---

### TC-UX-002: Error Recovery
**Priority**: High  
**Preconditions**: System access

**Test Steps**:
1. Trigger various errors
2. Upload invalid file
3. Ask while no documents
4. Disconnect network briefly

**Expected Results**:
- Clear error messages
- Recovery suggestions provided
- No data loss
- Graceful degradation

**Pass Criteria**: Users can recover from errors easily

---

### TC-UX-003: Help and Documentation
**Priority**: Medium  
**Preconditions**: System access

**Test Steps**:
1. Look for help options
2. Hover over UI elements
3. Check for tooltips
4. Access any documentation

**Expected Results**:
- Help text available
- Tooltips on key features
- Documentation accessible
- Examples provided

**Pass Criteria**: Adequate help available in UI

---

### TC-UX-004: Responsive Design
**Priority**: Low  
**Preconditions**: Access on different screen sizes

**Test Steps**:
1. Access on desktop (1920x1080)
2. Access on laptop (1366x768)
3. Resize browser window
4. Test all functions

**Expected Results**:
- Layout adapts to screen size
- All features accessible
- No horizontal scrolling
- Readable text at all sizes

**Pass Criteria**: UI works on common screen sizes

---

## Test Execution Summary

### Test Environment Requirements
- Browser: Chrome 90+, Firefox 88+, Safari 14+, Edge 90+
- Network: Stable internet connection
- Test Data: Set of organizational documents
- Users: 2-3 test users with different roles

### Test Data Requirements
- HR Policy document (PDF, 2MB)
- Employee Handbook (DOCX, 1MB)
- Training Presentation (PPTX, 5MB)
- Budget Spreadsheet (XLSX, 500KB)
- Meeting Notes (TXT, 50KB)
- Oversized file (>10MB) for negative testing
- Corrupted PDF for error testing

### Severity Definitions
- **Critical**: System unusable, data loss, security issue
- **High**: Major feature broken, significant workaround needed
- **Medium**: Feature impaired, workaround available
- **Low**: Minor issue, cosmetic, enhancement

### Pass/Fail Criteria
- **Pass**: 100% of Critical tests pass, 90% of High priority tests pass
- **Conditional Pass**: 80% of all tests pass, plan for fixes
- **Fail**: Any Critical test fails, <80% pass rate

---

## Defect Tracking Template

**Defect ID**: UAT-XXX  
**Date Found**: [Date]  
**Tester**: [Name]  
**Test Case**: [TC-XX-XXX]  
**Severity**: [Critical/High/Medium/Low]  
**Status**: [Open/In Progress/Fixed/Closed]  

**Description**:  
[Clear description of the issue]

**Steps to Reproduce**:  
1. [Step 1]
2. [Step 2]
3. [Step 3]

**Expected Result**:  
[What should happen]

**Actual Result**:  
[What actually happened]

**Screenshots/Evidence**:  
[Attach if applicable]

**Environment**:  
- Browser: [Version]
- OS: [Version]
- Test Data: [Used]

---

*End of UAT Test Cases Document*