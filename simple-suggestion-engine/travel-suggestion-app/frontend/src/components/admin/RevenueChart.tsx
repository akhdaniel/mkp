import '../../styles/admin/RevenueChart.css';

interface RevenueData {
  month: string;
  revenue: number;
  bookings: number;
}

interface RevenueChartProps {
  data: RevenueData[];
}

const RevenueChart = ({ data }: RevenueChartProps) => {
  const maxRevenue = Math.max(...data.map(d => d.revenue));

  return (
    <div className="revenue-chart">
      <div className="chart-bars">
        {data.map((item) => (
          <div key={item.month} className="bar-container">
            <div className="bar-wrapper">
              <div 
                className="revenue-bar"
                style={{ height: `${(item.revenue / maxRevenue) * 200}px` }}
              >
                <span className="bar-value">${(item.revenue / 1000).toFixed(1)}k</span>
              </div>
            </div>
            <div className="bar-label">
              <span className="month">{item.month}</span>
              <span className="bookings">{item.bookings} bookings</span>
            </div>
          </div>
        ))}
      </div>
      <div className="chart-legend">
        <span className="legend-item">
          <span className="legend-color revenue"></span>
          Revenue
        </span>
      </div>
    </div>
  );
};

export default RevenueChart;