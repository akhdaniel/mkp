import React, { useState, useEffect } from 'react';
import Plot from 'react-plotly.js';

function Dashboard() {
  const [selectedPeriod, setSelectedPeriod] = useState('7d');
  const [kpiData, setKpiData] = useState(null);
  const [isLoading, setIsLoading] = useState(false);

  // Sample KPI data
  const kpis = {
    totalRevenue: { value: 'Rp 125,000,000', change: '+12.5%', trend: 'up' },
    totalOrders: { value: '1,248', change: '+8.3%', trend: 'up' },
    avgOrderValue: { value: 'Rp 100,160', change: '-2.1%', trend: 'down' },
    customerSatisfaction: { value: '4.7/5.0', change: '+0.2', trend: 'up' },
  };

  // Sample revenue data for line chart
  const revenueData = {
    '7d': {
      dates: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
      values: [12000000, 14500000, 13200000, 15800000, 16200000, 21000000, 32280000]
    },
    '30d': {
      dates: Array.from({length: 30}, (_, i) => `Day ${i+1}`),
      values: Array.from({length: 30}, () => Math.floor(Math.random() * 20000000) + 10000000)
    },
    '90d': {
      dates: Array.from({length: 90}, (_, i) => `Day ${i+1}`),
      values: Array.from({length: 90}, () => Math.floor(Math.random() * 25000000) + 8000000)
    }
  };

  // Sample product performance data
  const productData = [
    { name: 'Product A', sales: 450, revenue: 45000000, growth: '+15%' },
    { name: 'Product B', sales: 320, revenue: 32000000, growth: '+8%' },
    { name: 'Product C', sales: 280, revenue: 21000000, growth: '-3%' },
    { name: 'Product D', sales: 198, revenue: 15840000, growth: '+22%' },
    { name: 'Product E', sales: 150, revenue: 11250000, growth: '+5%' },
  ];

  // Sample customer segmentation data
  const segmentData = {
    labels: ['New Customers', 'Returning', 'VIP', 'At Risk'],
    values: [35, 40, 15, 10],
    colors: ['#3b82f6', '#10b981', '#f59e0b', '#ef4444']
  };

  // Sample geographic data
  const geoData = [
    { region: 'Jakarta', value: 42000000, percentage: 33.6 },
    { region: 'Surabaya', value: 28000000, percentage: 22.4 },
    { region: 'Bandung', value: 18000000, percentage: 14.4 },
    { region: 'Medan', value: 15000000, percentage: 12.0 },
    { region: 'Semarang', value: 12000000, percentage: 9.6 },
    { region: 'Others', value: 10000000, percentage: 8.0 },
  ];

  // Sample correlation matrix
  const correlationMatrix = {
    z: [
      [1.0, 0.85, 0.73, -0.45, 0.62],
      [0.85, 1.0, 0.69, -0.38, 0.55],
      [0.73, 0.69, 1.0, -0.52, 0.48],
      [-0.45, -0.38, -0.52, 1.0, -0.33],
      [0.62, 0.55, 0.48, -0.33, 1.0]
    ],
    x: ['Price', 'Quality', 'Reviews', 'Returns', 'Sales'],
    y: ['Price', 'Quality', 'Reviews', 'Returns', 'Sales']
  };

  // Fetch live data (simulated)
  useEffect(() => {
    const fetchDashboardData = async () => {
      setIsLoading(true);
      // Simulate API call
      setTimeout(() => {
        setKpiData(kpis);
        setIsLoading(false);
      }, 500);
    };
    
    fetchDashboardData();
    
    // Set up auto-refresh every 30 seconds
    const interval = setInterval(fetchDashboardData, 30000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-4">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">Analytics Dashboard</h1>
              <p className="text-sm text-gray-600 mt-1">Real-time business metrics and insights</p>
            </div>
            <div className="flex items-center space-x-4">
              <select
                value={selectedPeriod}
                onChange={(e) => setSelectedPeriod(e.target.value)}
                className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="7d">Last 7 days</option>
                <option value="30d">Last 30 days</option>
                <option value="90d">Last 90 days</option>
              </select>
              <button className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
                Export Report
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* KPI Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between mb-2">
              <h3 className="text-sm font-medium text-gray-600">Total Revenue</h3>
              <span className={`text-xs px-2 py-1 rounded-full ${
                kpis.totalRevenue.trend === 'up' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
              }`}>
                {kpis.totalRevenue.change}
              </span>
            </div>
            <p className="text-2xl font-bold text-gray-900">{kpis.totalRevenue.value}</p>
            <p className="text-xs text-gray-500 mt-2">vs previous period</p>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between mb-2">
              <h3 className="text-sm font-medium text-gray-600">Total Orders</h3>
              <span className={`text-xs px-2 py-1 rounded-full ${
                kpis.totalOrders.trend === 'up' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
              }`}>
                {kpis.totalOrders.change}
              </span>
            </div>
            <p className="text-2xl font-bold text-gray-900">{kpis.totalOrders.value}</p>
            <p className="text-xs text-gray-500 mt-2">vs previous period</p>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between mb-2">
              <h3 className="text-sm font-medium text-gray-600">Avg Order Value</h3>
              <span className={`text-xs px-2 py-1 rounded-full ${
                kpis.avgOrderValue.trend === 'up' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
              }`}>
                {kpis.avgOrderValue.change}
              </span>
            </div>
            <p className="text-2xl font-bold text-gray-900">{kpis.avgOrderValue.value}</p>
            <p className="text-xs text-gray-500 mt-2">vs previous period</p>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between mb-2">
              <h3 className="text-sm font-medium text-gray-600">Customer Satisfaction</h3>
              <span className={`text-xs px-2 py-1 rounded-full ${
                kpis.customerSatisfaction.trend === 'up' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
              }`}>
                {kpis.customerSatisfaction.change}
              </span>
            </div>
            <p className="text-2xl font-bold text-gray-900">{kpis.customerSatisfaction.value}</p>
            <p className="text-xs text-gray-500 mt-2">vs previous period</p>
          </div>
        </div>

        {/* Charts Row 1 */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
          {/* Revenue Trend */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold mb-4">Revenue Trend</h3>
            <Plot
              data={[
                {
                  x: revenueData[selectedPeriod].dates,
                  y: revenueData[selectedPeriod].values,
                  type: 'scatter',
                  mode: 'lines+markers',
                  marker: { color: '#3b82f6' },
                  line: { shape: 'spline', smoothing: 1.3 },
                  fill: 'tozeroy',
                  fillcolor: 'rgba(59, 130, 246, 0.1)'
                }
              ]}
              layout={{
                autosize: true,
                margin: { l: 40, r: 20, t: 20, b: 40 },
                xaxis: { showgrid: false },
                yaxis: { 
                  tickformat: ',.',
                  showgrid: true,
                  gridcolor: '#e5e7eb'
                },
                hovermode: 'x unified',
                showlegend: false
              }}
              config={{ responsive: true, displayModeBar: false }}
              className="w-full"
              style={{ width: '100%', height: '300px' }}
            />
          </div>

          {/* Customer Segmentation */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold mb-4">Customer Segmentation</h3>
            <Plot
              data={[
                {
                  labels: segmentData.labels,
                  values: segmentData.values,
                  type: 'pie',
                  hole: 0.4,
                  marker: {
                    colors: segmentData.colors
                  },
                  textinfo: 'label+percent',
                  textposition: 'outside'
                }
              ]}
              layout={{
                autosize: true,
                margin: { l: 20, r: 20, t: 20, b: 20 },
                showlegend: false
              }}
              config={{ responsive: true, displayModeBar: false }}
              className="w-full"
              style={{ width: '100%', height: '300px' }}
            />
          </div>
        </div>

        {/* Product Performance Table */}
        <div className="bg-white rounded-lg shadow mb-6">
          <div className="px-6 py-4 border-b">
            <h3 className="text-lg font-semibold">Top Products Performance</h3>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Product Name
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Units Sold
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Revenue
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Growth
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {productData.map((product, index) => (
                  <tr key={index} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                      {product.name}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {product.sales}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      Rp {product.revenue.toLocaleString('id-ID')}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      <span className={`px-2 py-1 text-xs rounded-full ${
                        product.growth.startsWith('+') ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                      }`}>
                        {product.growth}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Charts Row 2 */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Geographic Distribution */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold mb-4">Revenue by Region</h3>
            <Plot
              data={[
                {
                  x: geoData.map(d => d.region),
                  y: geoData.map(d => d.value),
                  type: 'bar',
                  marker: {
                    color: ['#3b82f6', '#10b981', '#f59e0b', '#8b5cf6', '#ec4899', '#6b7280']
                  },
                  text: geoData.map(d => `${d.percentage}%`),
                  textposition: 'outside'
                }
              ]}
              layout={{
                autosize: true,
                margin: { l: 40, r: 20, t: 20, b: 80 },
                xaxis: { 
                  showgrid: false,
                  tickangle: -45
                },
                yaxis: { 
                  tickformat: ',.',
                  showgrid: true,
                  gridcolor: '#e5e7eb'
                },
                showlegend: false
              }}
              config={{ responsive: true, displayModeBar: false }}
              className="w-full"
              style={{ width: '100%', height: '300px' }}
            />
          </div>

          {/* Correlation Matrix */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold mb-4">Metrics Correlation</h3>
            <Plot
              data={[
                {
                  z: correlationMatrix.z,
                  x: correlationMatrix.x,
                  y: correlationMatrix.y,
                  type: 'heatmap',
                  colorscale: 'RdBu',
                  reversescale: true,
                  showscale: true
                }
              ]}
              layout={{
                autosize: true,
                margin: { l: 60, r: 40, t: 20, b: 60 },
                xaxis: { side: 'bottom' },
                yaxis: { autorange: 'reversed' }
              }}
              config={{ responsive: true, displayModeBar: false }}
              className="w-full"
              style={{ width: '100%', height: '300px' }}
            />
          </div>
        </div>

        {/* AI Insights Section */}
        <div className="mt-6 bg-gradient-to-r from-blue-50 to-indigo-50 rounded-lg shadow p-6">
          <div className="flex items-start space-x-3">
            <div className="flex-shrink-0">
              <svg className="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </div>
            <div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">AI-Powered Insights</h3>
              <ul className="space-y-2 text-sm text-gray-700">
                <li className="flex items-start">
                  <span className="text-green-500 mr-2">•</span>
                  <span>Revenue increased by 12.5% compared to the previous period, with Saturday showing the highest sales.</span>
                </li>
                <li className="flex items-start">
                  <span className="text-yellow-500 mr-2">•</span>
                  <span>Product C is underperforming with -3% growth. Consider promotional strategies or inventory adjustments.</span>
                </li>
                <li className="flex items-start">
                  <span className="text-blue-500 mr-2">•</span>
                  <span>Strong correlation (0.85) detected between product quality ratings and sales volume.</span>
                </li>
                <li className="flex items-start">
                  <span className="text-purple-500 mr-2">•</span>
                  <span>Jakarta region contributes 33.6% of total revenue - opportunity to expand in other regions.</span>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}

export default Dashboard;