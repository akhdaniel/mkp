from fastapi import FastAPI, File, UploadFile, Form, HTTPException
from fastapi.middleware.cors import CORSMiddleware
import pandas as pd
import numpy as np
import io
import json
import os
from typing import Optional, Dict, Any
import plotly.express as px
import plotly.graph_objects as go
import plotly.utils
from datetime import datetime
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

app = FastAPI(title="DataInsight AI", version="0.1.0")

# Allow CORS for frontend communication
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000"],  # Allows frontend to connect
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/")
def read_root():
    return {"message": "Welcome to DataInsight AI Backend", "version": "0.1.0"}

@app.get("/api/health")
def health_check():
    return {"status": "healthy", "timestamp": datetime.now().isoformat()}

def detect_file_type(filename: str) -> str:
    """Detect file type from extension"""
    ext = filename.lower().split('.')[-1]
    if ext in ['xlsx', 'xls']:
        return 'excel'
    elif ext == 'csv':
        return 'csv'
    elif ext == 'json':
        return 'json'
    else:
        return 'unknown'

def load_dataframe(file: UploadFile, contents: bytes) -> pd.DataFrame:
    """Load data into DataFrame based on file type"""
    file_type = detect_file_type(file.filename)
    
    if file_type == 'excel':
        return pd.read_excel(io.BytesIO(contents))
    elif file_type == 'csv':
        return pd.read_csv(io.BytesIO(contents))
    elif file_type == 'json':
        return pd.read_json(io.BytesIO(contents))
    else:
        raise ValueError(f"Unsupported file type: {file.filename}")

def generate_insights(df: pd.DataFrame) -> Dict[str, Any]:
    """Generate automatic insights from the dataframe"""
    insights = {
        "basic_info": {
            "rows": len(df),
            "columns": len(df.columns),
            "column_names": df.columns.tolist(),
            "data_types": df.dtypes.astype(str).to_dict()
        },
        "missing_values": df.isnull().sum().to_dict(),
        "numeric_summary": {}
    }
    
    # Get numeric columns summary
    numeric_cols = df.select_dtypes(include=[np.number]).columns.tolist()
    if numeric_cols:
        insights["numeric_summary"] = df[numeric_cols].describe().to_dict()
    
    # Detect outliers using IQR method
    outliers = {}
    for col in numeric_cols:
        Q1 = df[col].quantile(0.25)
        Q3 = df[col].quantile(0.75)
        IQR = Q3 - Q1
        lower = Q1 - 1.5 * IQR
        upper = Q3 + 1.5 * IQR
        outlier_count = ((df[col] < lower) | (df[col] > upper)).sum()
        if outlier_count > 0:
            outliers[col] = outlier_count
    insights["outliers"] = outliers
    
    # Detect correlations
    if len(numeric_cols) > 1:
        corr_matrix = df[numeric_cols].corr()
        high_corr = {}
        for i in range(len(corr_matrix.columns)):
            for j in range(i+1, len(corr_matrix.columns)):
                corr_val = corr_matrix.iloc[i, j]
                if abs(corr_val) > 0.7:
                    high_corr[f"{corr_matrix.columns[i]} & {corr_matrix.columns[j]}"] = round(corr_val, 3)
        insights["high_correlations"] = high_corr
    
    return insights

def create_visualization(df: pd.DataFrame, viz_type: str = "auto") -> Optional[str]:
    """Create interactive visualizations using Plotly"""
    try:
        numeric_cols = df.select_dtypes(include=[np.number]).columns.tolist()
        
        if not numeric_cols:
            return None
        
        if viz_type == "auto" or viz_type == "distribution":
            # Create distribution plots for numeric columns
            if len(numeric_cols) == 1:
                fig = px.histogram(df, x=numeric_cols[0], title=f"Distribution of {numeric_cols[0]}")
            else:
                # Create subplot for first 2 numeric columns
                fig = go.Figure()
                for col in numeric_cols[:2]:
                    fig.add_trace(go.Histogram(x=df[col], name=col, opacity=0.7))
                fig.update_layout(
                    title="Distribution Comparison",
                    xaxis_title="Value",
                    yaxis_title="Count",
                    barmode='overlay'
                )
        elif viz_type == "correlation":
            if len(numeric_cols) > 1:
                corr = df[numeric_cols].corr()
                fig = px.imshow(corr, 
                              labels=dict(x="Features", y="Features", color="Correlation"),
                              title="Correlation Matrix",
                              color_continuous_scale="RdBu",
                              aspect="auto")
            else:
                return None
        else:
            # Default scatter plot if we have 2+ numeric columns
            if len(numeric_cols) >= 2:
                fig = px.scatter(df, x=numeric_cols[0], y=numeric_cols[1], 
                               title=f"{numeric_cols[0]} vs {numeric_cols[1]}")
            else:
                fig = px.histogram(df, x=numeric_cols[0], title=f"Distribution of {numeric_cols[0]}")
        
        # Convert to JSON
        graph_json = json.dumps(fig, cls=plotly.utils.PlotlyJSONEncoder)
        return graph_json
    except Exception as e:
        print(f"Visualization error: {e}")
        return None

@app.post("/api/process")
async def process_file(
    file: UploadFile = File(...),
    prompt: str = Form(...),
):
    try:
        # Read the uploaded file
        contents = await file.read()
        df = load_dataframe(file, contents)
        
        result = {
            "table_html": "",
            "insights": {},
            "visualization": None,
            "message": ""
        }
        
        prompt_lower = prompt.lower()
        
        # Process commands
        if "head" in prompt_lower or "atas" in prompt_lower or "first" in prompt_lower:
            result["table_html"] = df.head(10).to_html(classes="table-auto w-full text-left")
            result["message"] = "Showing first 10 rows of data"
            
        elif "tail" in prompt_lower or "bawah" in prompt_lower or "last" in prompt_lower:
            result["table_html"] = df.tail(10).to_html(classes="table-auto w-full text-left")
            result["message"] = "Showing last 10 rows of data"
            
        elif "summary" in prompt_lower or "ringkasan" in prompt_lower or "describe" in prompt_lower:
            result["table_html"] = df.describe().to_html(classes="table-auto w-full text-left")
            result["message"] = "Statistical summary of numeric columns"
            result["visualization"] = create_visualization(df, "distribution")
            
        elif "info" in prompt_lower:
            buffer = io.StringIO()
            df.info(buf=buffer)
            info_str = buffer.getvalue()
            result["table_html"] = f"<pre class='bg-gray-100 p-4 rounded'>{info_str}</pre>"
            result["message"] = "Dataset information"
            
        elif "correlation" in prompt_lower or "korelasi" in prompt_lower:
            numeric_cols = df.select_dtypes(include=[np.number]).columns
            if len(numeric_cols) > 1:
                corr = df[numeric_cols].corr()
                result["table_html"] = corr.to_html(classes="table-auto w-full text-left")
                result["visualization"] = create_visualization(df, "correlation")
                result["message"] = "Correlation matrix"
            else:
                result["message"] = "Not enough numeric columns for correlation analysis"
                
        elif "missing" in prompt_lower or "null" in prompt_lower:
            missing = df.isnull().sum()
            missing_df = pd.DataFrame({'Column': missing.index, 'Missing Values': missing.values})
            result["table_html"] = missing_df[missing_df['Missing Values'] > 0].to_html(classes="table-auto w-full text-left", index=False)
            result["message"] = "Missing values summary"
            
        elif "insight" in prompt_lower or "analyze" in prompt_lower or "analisis" in prompt_lower:
            insights = generate_insights(df)
            result["insights"] = insights
            result["table_html"] = df.head(10).to_html(classes="table-auto w-full text-left")
            result["visualization"] = create_visualization(df, "auto")
            result["message"] = "Auto-generated insights"
            
        else:
            # Default: show basic info and insights
            insights = generate_insights(df)
            result["insights"] = insights
            result["table_html"] = df.head(10).to_html(classes="table-auto w-full text-left")
            result["message"] = "Data loaded successfully. Try commands: 'insights', 'summary', 'correlation', 'missing values'"
        
        return result
        
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error processing file: {str(e)}")