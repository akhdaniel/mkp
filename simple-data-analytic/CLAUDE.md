# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

DataInsight AI is a bilingual (English/Indonesian) data analytics platform that helps businesses and individuals analyze data without deep technical expertise. The platform uses AI to automatically generate insights, detect patterns, and create visualizations from uploaded data.

## Development Commands

### Quick Start (Both Services)
```bash
./start.sh  # Starts both backend (port 8000) and frontend (port 3000)
```

### Backend Development
```bash
cd backend
source venv/bin/activate  # Activate virtual environment
uvicorn main:app --reload --port 8000  # Start FastAPI server with auto-reload
```

### Frontend Development
```bash
cd frontend
npm start  # Start React development server on port 3000
npm run build  # Build production bundle
```

### Testing
```bash
# Backend - No test framework currently configured
# Frontend - via Create React App:
cd frontend && npm test
```

## Architecture Overview

### Full-Stack Structure
The application follows a clear separation between frontend and backend:
- **Frontend**: React SPA serving UI and visualizations
- **Backend**: FastAPI REST API handling data processing and analysis
- **Communication**: RESTful API with CORS enabled for localhost:3000

### Key Design Patterns

1. **API Design**: All backend endpoints are prefixed with `/api/` for clear separation
2. **File Processing**: Backend handles multiple formats (CSV, Excel, JSON) with pandas
3. **Visualization Strategy**: Backend generates Plotly configurations, frontend renders them
4. **Error Handling**: HTTPException for API errors, user-friendly messages in frontend
5. **State Management**: React hooks (useState, useEffect) for frontend state

### Data Flow
1. User uploads file through React UI
2. Frontend sends file to backend `/api/analyze` endpoint
3. Backend processes with pandas, generates insights and visualizations
4. Backend returns JSON with analysis results and Plotly chart configs
5. Frontend renders results and interactive charts

## Working with Key Features

### Adding New Analysis Commands
Commands are processed in `backend/main.py` within the `process_command()` function. To add a new command:
1. Add the command logic to `process_command()`
2. Handle both English and Indonesian variants
3. Return structured response with `result`, `visualization`, and optional `error`

### Creating Visualizations
The backend generates Plotly configurations that the frontend renders. Pattern:
```python
# Backend: Generate Plotly figure
fig = px.scatter(df, x='col1', y='col2')
chart_json = json.loads(plotly.utils.PlotlyJSONEncoder().encode(fig))
return {"visualization": chart_json}
```

### Bilingual Support
The system supports English/Indonesian commands. When adding features:
- Support command aliases in both languages (e.g., "summary"/"ringkasan")
- Return user messages in the requested language when possible
- Maintain bilingual documentation in README.md and spesifikasi.md

## Current Implementation Status

### Completed Features (MVP)
- Multi-format file upload (CSV, Excel, JSON)
- Natural language command processing
- Statistical analysis (summary, correlation, missing values)
- Automatic insights generation (outliers, patterns)
- Interactive Plotly visualizations
- Bilingual command support

### Prepared but Not Integrated
- OpenAI API setup (API key configured but not used)
- Advanced pandas/numpy operations (commented in requirements.txt)

### Next Development Phase
According to spesifikasi.md, priorities include:
- OpenAI integration for natural language queries
- Real-time data streaming
- User authentication and workspaces
- Custom dashboard builder
- Data export functionality

## Important Technical Notes

### Environment Configuration
- Backend uses `python-dotenv` for environment variables
- Create `.env` file in `/backend/` for API keys
- OpenAI API key should be set as `OPENAI_API_KEY`

### Frontend Build Configuration
- Tailwind CSS configured with PostCSS
- Both npm and yarn lock files present (prefer npm for consistency)
- Production build outputs to `frontend/build/`

### CORS Configuration
Currently allows only `http://localhost:3000`. For production:
- Update CORS origins in `backend/main.py`
- Consider environment-based configuration

### Data Processing Constraints
- File size limits not explicitly set (add for production)
- All processing happens in-memory (consider streaming for large files)
- No data persistence (all analysis is session-based)

## Code Style Conventions

### Python (Backend)
- Type hints for function parameters and returns
- Docstrings for complex functions
- FastAPI route decorators with clear paths
- HTTPException for error handling

### JavaScript/React (Frontend)
- Functional components with hooks
- Tailwind utility classes for styling
- Async/await for API calls
- Component files match component names

## Debugging Tips

### Common Issues
1. **CORS errors**: Check backend is running and CORS middleware configured
2. **File upload fails**: Verify file format and check backend console for pandas errors
3. **Charts not rendering**: Ensure Plotly.js is properly imported in frontend
4. **Commands not recognized**: Check command aliases in `process_command()`

### API Testing
- FastAPI automatic docs: http://localhost:8000/docs
- ReDoc alternative: http://localhost:8000/redoc
- Test endpoints directly with Swagger UI

## Production Considerations

Before deploying to production:
1. Add comprehensive error logging
2. Implement file size limits and validation
3. Set up proper environment configuration
4. Add rate limiting for API endpoints
5. Implement user authentication
6. Configure production CORS settings
7. Set up monitoring and analytics
8. Add data persistence layer
9. Implement proper testing suite
10. Configure CI/CD pipeline