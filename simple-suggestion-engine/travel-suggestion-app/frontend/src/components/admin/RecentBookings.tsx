import '../../styles/admin/RecentBookings.css';

interface Booking {
  packageId: number;
  purchaseDate: string;
  status: string;
  userId: string;
  packageDetails: any;
}

interface RecentBookingsProps {
  bookings: Booking[];
}

const RecentBookings = ({ bookings }: RecentBookingsProps) => {
  return (
    <div className="recent-bookings">
      <h2>Recent Bookings</h2>
      <div className="bookings-list">
        {bookings.map((booking, index) => (
          <div key={index} className="booking-item">
            <div className="booking-date">
              {new Date(booking.purchaseDate).toLocaleDateString()}
            </div>
            <div className="booking-details">
              <h4>{booking.packageDetails?.title}</h4>
              <p>
                {booking.packageDetails?.destination}, {booking.packageDetails?.country}
              </p>
              <p className="user-id">User: {booking.userId}</p>
            </div>
            <div className="booking-price">
              ${booking.packageDetails?.price}
            </div>
            <span className={`booking-status ${booking.status}`}>
              {booking.status}
            </span>
          </div>
        ))}
      </div>
    </div>
  );
};

export default RecentBookings;