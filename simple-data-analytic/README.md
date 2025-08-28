# DataInsight AI - Platform Analytics Cerdas

Platform analitik data yang menggunakan AI untuk membantu bisnis dan individu menganalisis data dengan mudah.

## Features

### Current Features (MVP)
- **Multi-format Support**: Upload CSV, Excel (.xlsx, .xls), and JSON files
- **Auto Insights**: Automatic detection of patterns, outliers, and correlations
- **Natural Language Commands**: Simple commands in English or Indonesian
- **Interactive Visualizations**: Distribution plots, correlation matrices using Plotly
- **Statistical Analysis**: Summary statistics, missing value detection
- **Data Exploration**: View head/tail of data, dataset info

### Supported Commands
- `insights` / `analyze` - Generate automatic insights
- `summary` / `ringkasan` - Statistical summary with visualizations
- `correlation` / `korelasi` - Correlation matrix and heatmap
- `missing values` / `null` - Missing data analysis
- `head` / `atas` - First 10 rows
- `tail` / `bawah` - Last 10 rows  
- `info` - Dataset structure information

## Tech Stack

### Backend
- FastAPI (Python)
- Pandas for data processing
- Plotly for visualizations
- NumPy for numerical operations

### Frontend
- React.js
- Tailwind CSS for styling
- Plotly.js for interactive charts
- Chart.js for additional visualizations

## Installation & Setup

### Backend Setup

1. Navigate to backend directory:
```bash
cd backend
```

2. Create virtual environment:
```bash
python -m venv venv
```

3. Activate virtual environment:
- Windows: `venv\Scripts\activate`
- Mac/Linux: `source venv/bin/activate`

4. Install dependencies:
```bash
pip install -r requirements.txt
```

5. Run the backend server:
```bash
uvicorn main:app --reload --port 8000
```

Backend will be available at: http://localhost:8000

### Frontend Setup

1. Navigate to frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Start the development server:
```bash
npm start
```

Frontend will be available at: http://localhost:3000

## Usage

1. Open browser and navigate to http://localhost:3000
2. Upload a data file (CSV, Excel, or JSON)
3. Enter a command (e.g., "insights", "summary", "correlation")
4. Click "Run Analysis" to process
5. View results including:
   - AI-generated insights
   - Interactive visualizations
   - Data tables

## Sample Data

You can test the application with any CSV or Excel file containing numerical data. The system will automatically:
- Detect data types
- Generate appropriate visualizations
- Identify patterns and anomalies
- Calculate correlations

## API Documentation

When backend is running, visit:
- API Docs: http://localhost:8000/docs
- ReDoc: http://localhost:8000/redoc

## Future Enhancements

- OpenAI integration for natural language queries
- Real-time data streaming
- Advanced predictive analytics
- Custom dashboard builder
- Data export functionality
- User authentication & workspaces

## License

MIT