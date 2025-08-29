import React, { useState, useRef } from 'react';
import './App.css';
import axios from 'axios';

function App() {
  const [selectedImage, setSelectedImage] = useState(null);
  const [previewUrl, setPreviewUrl] = useState('');
  const [extractedData, setExtractedData] = useState(null);
  const [documents, setDocuments] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [activeTab, setActiveTab] = useState('upload');
  const fileInputRef = useRef(null);

  // Handle image selection
  const handleImageChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      setSelectedImage(file);
      setPreviewUrl(URL.createObjectURL(file));
      setExtractedData(null);
      setError('');
    }
  };

  // Handle image upload
  const handleUpload = async () => {
    if (!selectedImage) {
      setError('Please select an image first');
      return;
    }

    setLoading(true);
    setError('');

    const formData = new FormData();
    formData.append('image', selectedImage);

    try {
      const response = await axios.post('/api/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      });

      setExtractedData(response.data.data);
      // Refresh documents list
      fetchDocuments();
    } catch (err) {
      setError('Failed to upload image: ' + (err.response?.data?.error || err.message));
    } finally {
      setLoading(false);
    }
  };

  // Fetch all documents
  const fetchDocuments = async () => {
    try {
      const response = await axios.get('/api/documents');
      setDocuments(response.data);
    } catch (err) {
      setError('Failed to fetch documents: ' + err.message);
    }
  };

  // Load documents on component mount
  React.useEffect(() => {
    fetchDocuments();
  }, []);

  // Trigger file input click
  const triggerFileInput = () => {
    fileInputRef.current.click();
  };

  // Format date for display
  const formatDate = (dateString) => {
    if (!dateString) return 'N/A';
    return new Date(dateString).toLocaleDateString();
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>ID Document Capture</h1>
        <p>Upload and extract information from ID documents</p>
      </header>

      <main className="App-main">
        <div className="tabs">
          <button 
            className={activeTab === 'upload' ? 'active' : ''} 
            onClick={() => setActiveTab('upload')}
          >
            Upload Document
          </button>
          <button 
            className={activeTab === 'history' ? 'active' : ''} 
            onClick={() => setActiveTab('history')}
          >
            Document History
          </button>
        </div>

        {error && <div className="error-message">{error}</div>}

        {activeTab === 'upload' && (
          <div className="upload-section">
            <div className="upload-area" onClick={triggerFileInput}>
              <input
                type="file"
                ref={fileInputRef}
                onChange={handleImageChange}
                accept="image/*"
                style={{ display: 'none' }}
              />
              {previewUrl ? (
                <img src={previewUrl} alt="Preview" className="image-preview" />
              ) : (
                <div className="upload-placeholder">
                  <p>Click to select an ID document image</p>
                  <p className="upload-hint">Supported formats: JPG, PNG, etc.</p>
                </div>
              )}
            </div>

            {selectedImage && (
              <div className="upload-actions">
                <button onClick={handleUpload} disabled={loading}>
                  {loading ? 'Processing...' : 'Extract Information'}
                </button>
              </div>
            )}

            {extractedData && (
              <div className="extraction-results">
                <h2>Extracted Information</h2>
                <div className="data-grid">
                  <div className="data-item">
                    <label>Name:</label>
                    <span>{extractedData.name || 'N/A'}</span>
                  </div>
                  <div className="data-item">
                    <label>Document Number:</label>
                    <span>{extractedData.document_number || 'N/A'}</span>
                  </div>
                  <div className="data-item">
                    <label>Birth Date:</label>
                    <span>{formatDate(extractedData.birth_date)}</span>
                  </div>
                  <div className="data-item">
                    <label>Issue Date:</label>
                    <span>{formatDate(extractedData.issue_date)}</span>
                  </div>
                  <div className="data-item">
                    <label>Expiry Date:</label>
                    <span>{formatDate(extractedData.expiry_date)}</span>
                  </div>
                </div>
              </div>
            )}
          </div>
        )}

        {activeTab === 'history' && (
          <div className="history-section">
            <h2>Document History</h2>
            {documents.length === 0 ? (
              <p className="no-documents">No documents uploaded yet</p>
            ) : (
              <div className="documents-grid">
                {documents.map((doc) => (
                  <div key={doc.id} className="document-card">
                    <div className="card-header">
                      <h3>{doc.name || 'Unnamed Document'}</h3>
                      <span className="document-id">ID: {doc.document_id}</span>
                    </div>
                    {doc.image_path && (
                      <img 
                        src={doc.image_path} 
                        alt="Document" 
                        className="document-thumbnail"
                      />
                    )}
                    <div className="card-details">
                      <p><strong>Number:</strong> {doc.document_number || 'N/A'}</p>
                      <p><strong>Birth Date:</strong> {formatDate(doc.birth_date)}</p>
                      <p><strong>Issue Date:</strong> {formatDate(doc.issue_date)}</p>
                      <p><strong>Expiry Date:</strong> {formatDate(doc.expiry_date)}</p>
                      <p className="upload-date">Uploaded: {new Date(doc.created_at).toLocaleString()}</p>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}
      </main>
    </div>
  );
}

export default App;