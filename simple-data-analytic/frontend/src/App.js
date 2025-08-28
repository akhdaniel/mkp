import React, { useState, useEffect } from 'react';
import { Chart as ChartJS, CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend, PointElement, LineElement, ArcElement } from 'chart.js';
import { Bar, Line, Pie, Scatter } from 'react-chartjs-2';
import Plot from 'react-plotly.js';

// Register ChartJS components
ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  PointElement,
  LineElement,
  ArcElement,
  Title,
  Tooltip,
  Legend
);

function App() {
  const [file, setFile] = useState(null);
  const [prompt, setPrompt] = useState('insights');
  const [output, setOutput] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [insights, setInsights] = useState(null);
  const [visualization, setVisualization] = useState(null);
  const [message, setMessage] = useState('');
  const [suggestions] = useState([
    'insights', 'summary', 'correlation', 'missing values',
    'head', 'tail', 'info'
  ]);

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const handlePromptChange = (e) => {
    setPrompt(e.target.value);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!file) {
      setError('Please select a file first.');
      return;
    }

    const formData = new FormData();
    formData.append('file', file);
    formData.append('prompt', prompt);

    setIsLoading(true);
    setError('');
    setOutput('');

    try {
      const response = await fetch('http://127.0.0.1:8000/api/process', {
        method: 'POST',
        body: formData,
      });

      if (!response.ok) {
        const errData = await response.json();
        throw new Error(errData.detail || 'Something went wrong');
      }

      const data = await response.json();
      setOutput(data.table_html || '');
      setInsights(data.insights || null);
      setVisualization(data.visualization || null);
      setMessage(data.message || '');
    } catch (err) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-100 text-gray-800">
      <header className="bg-white shadow-md">
        <div className="container mx-auto px-4 py-6">
          <h1 className="text-3xl font-bold text-gray-900">DataInsight AI</h1>
          <p className="text-sm text-gray-600">Upload an Excel file and give a command to analyze.</p>
        </div>
      </header>

      <main className="container mx-auto px-4 py-8">
        <div className="bg-white p-6 rounded-lg shadow-lg">
          <form onSubmit={handleSubmit}>
            <div className="mb-4">
              <label htmlFor="file-upload" className="block text-sm font-medium text-gray-700 mb-2">
                Excel File
              </label>
              <input 
                id="file-upload" 
                type="file" 
                onChange={handleFileChange} 
                accept=".xlsx, .xls, .csv, .json"
                className="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100" 
              />
            </div>

            <div className="mb-4">
              <label htmlFor="prompt-input" className="block text-sm font-medium text-gray-700 mb-2">
                Command Prompt
              </label>
              <input 
                id="prompt-input" 
                type="text" 
                value={prompt} 
                onChange={handlePromptChange} 
                placeholder="Try: insights, summary, correlation, missing values"
                className="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              />
            </div>

            <button 
              type="submit" 
              disabled={isLoading}
              className="w-full inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:bg-gray-400"
            >
              {isLoading ? 'Processing...' : 'Run Analysis'}
            </button>
          </form>

          {error && <div className="mt-4 text-red-600 bg-red-100 p-3 rounded">Error: {error}</div>}
          
          {/* Command Suggestions */}
          <div className="mt-4">
            <p className="text-sm text-gray-600 mb-2">Quick commands:</p>
            <div className="flex flex-wrap gap-2">
              {suggestions.map(cmd => (
                <button
                  key={cmd}
                  type="button"
                  onClick={() => setPrompt(cmd)}
                  className="px-3 py-1 bg-gray-100 hover:bg-gray-200 rounded-full text-sm"
                >
                  {cmd}
                </button>
              ))}
            </div>
          </div>
        </div>

        {/* Message Display */}
        {message && (
          <div className="mt-4 bg-blue-50 border-l-4 border-blue-500 p-4">
            <p className="text-blue-700">{message}</p>
          </div>
        )}

        {/* Insights Display */}
        {insights && (
          <div className="mt-8 bg-white p-6 rounded-lg shadow-lg">
            <h2 className="text-xl font-semibold mb-4">AI Insights</h2>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
              <div className="bg-gray-50 p-4 rounded">
                <h3 className="font-semibold text-gray-700 mb-2">Dataset Overview</h3>
                <p>Rows: {insights.basic_info?.rows}</p>
                <p>Columns: {insights.basic_info?.columns}</p>
              </div>
              
              {insights.outliers && Object.keys(insights.outliers).length > 0 && (
                <div className="bg-yellow-50 p-4 rounded">
                  <h3 className="font-semibold text-gray-700 mb-2">Outliers Detected</h3>
                  {Object.entries(insights.outliers).map(([col, count]) => (
                    <p key={col}>{col}: {count} outliers</p>
                  ))}
                </div>
              )}
              
              {insights.high_correlations && Object.keys(insights.high_correlations).length > 0 && (
                <div className="bg-green-50 p-4 rounded">
                  <h3 className="font-semibold text-gray-700 mb-2">High Correlations</h3>
                  {Object.entries(insights.high_correlations).map(([pair, corr]) => (
                    <p key={pair}>{pair}: {corr}</p>
                  ))}
                </div>
              )}
            </div>
          </div>
        )}

        {/* Visualization Display */}
        {visualization && (
          <div className="mt-8 bg-white p-6 rounded-lg shadow-lg">
            <h2 className="text-xl font-semibold mb-4">Data Visualization</h2>
            <Plot
              data={JSON.parse(visualization).data}
              layout={JSON.parse(visualization).layout}
              config={{responsive: true}}
              className="w-full"
            />
          </div>
        )}

        {/* Table Output */}
        {output && (
          <div className="mt-8 bg-white p-6 rounded-lg shadow-lg">
            <h2 className="text-xl font-semibold mb-4">Data Table</h2>
            <div 
              className="prose max-w-none overflow-x-auto"
              dangerouslySetInnerHTML={{ __html: output }}
            />
          </div>
        )}
      </main>
    </div>
  );
}

export default App;