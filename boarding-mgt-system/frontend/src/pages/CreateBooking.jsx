import React, { useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { bookingAPI } from '../api/booking';
import { CheckCircleIcon } from '@heroicons/react/24/outline';

export default function CreateBooking() {
  const location = useLocation();
  const navigate = useNavigate();
  const { user } = useAuth();
  const { schedule, passengerCount } = location.state || {};
  
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const [bookingData, setBookingData] = useState(null);
  
  const [passengers, setPassengers] = useState(
    Array.from({ length: passengerCount || 1 }, () => ({
      name: '',
      type: 'adult',
      seat_number: '',
    }))
  );
  
  const [paymentMethod, setPaymentMethod] = useState('credit_card');
  const [specialRequirements, setSpecialRequirements] = useState('');

  if (!schedule) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500">No schedule selected. Please search for schedules first.</p>
        <button
          onClick={() => navigate('/search')}
          className="mt-4 text-primary-600 hover:text-primary-500"
        >
          Go to Search
        </button>
      </div>
    );
  }

  const handlePassengerChange = (index, field, value) => {
    const updated = [...passengers];
    updated[index][field] = value;
    setPassengers(updated);
  };

  const calculateTotal = () => {
    return passengers.reduce((total, passenger) => {
      let price = schedule.base_price;
      if (passenger.type === 'child') price *= 0.5;
      if (passenger.type === 'infant') price = 0;
      if (passenger.type === 'senior') price *= 0.8;
      return total + price;
    }, 0);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const bookingRequest = {
        schedule_id: schedule.id,
        passengers: passengers,
        payment_method: paymentMethod,
        special_requirements: specialRequirements,
      };

      const response = await bookingAPI.createBooking(bookingRequest);
      setBookingData(response);
      setSuccess(true);
    } catch (error) {
      setError(error.response?.data?.error || 'Failed to create booking');
    } finally {
      setLoading(false);
    }
  };

  if (success && bookingData) {
    return (
      <div className="max-w-3xl mx-auto">
        <div className="bg-white shadow sm:rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <div className="text-center">
              <CheckCircleIcon className="mx-auto h-12 w-12 text-green-500" />
              <h3 className="mt-2 text-lg font-medium text-gray-900">Booking Confirmed!</h3>
              <p className="mt-1 text-sm text-gray-500">
                Your booking reference is: <span className="font-bold">{bookingData.booking_reference}</span>
              </p>
              <div className="mt-6">
                <button
                  onClick={() => navigate(`/bookings/${bookingData.id}`)}
                  className="inline-flex items-center rounded-md bg-primary-600 px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-primary-500"
                >
                  View Booking Details
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-3xl mx-auto">
      <h2 className="text-2xl font-bold text-gray-900 mb-6">Complete Your Booking</h2>
      
      <div className="bg-white shadow sm:rounded-lg mb-6">
        <div className="px-4 py-5 sm:p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Schedule Details</h3>
          <dl className="grid grid-cols-2 gap-4">
            <div>
              <dt className="text-sm font-medium text-gray-500">Route</dt>
              <dd className="text-sm text-gray-900">{schedule.route?.name || 'Direct Route'}</dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-gray-500">Date</dt>
              <dd className="text-sm text-gray-900">{schedule.departure_date}</dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-gray-500">Departure Time</dt>
              <dd className="text-sm text-gray-900">{schedule.departure_time}</dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-gray-500">Arrival Time</dt>
              <dd className="text-sm text-gray-900">{schedule.arrival_time}</dd>
            </div>
          </dl>
        </div>
      </div>

      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="bg-white shadow sm:rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <h3 className="text-lg font-medium text-gray-900 mb-4">Passenger Information</h3>
            
            {passengers.map((passenger, index) => (
              <div key={index} className="border-b border-gray-200 pb-4 mb-4 last:border-0">
                <h4 className="text-sm font-medium text-gray-700 mb-3">Passenger {index + 1}</h4>
                <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      Full Name
                    </label>
                    <input
                      type="text"
                      required
                      value={passenger.name}
                      onChange={(e) => handlePassengerChange(index, 'name', e.target.value)}
                      className="mt-1 block w-full rounded-md border-gray-300 px-3 py-2 shadow-sm focus:border-primary-500 focus:outline-none focus:ring-primary-500 sm:text-sm border"
                    />
                  </div>
                  
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      Type
                    </label>
                    <select
                      value={passenger.type}
                      onChange={(e) => handlePassengerChange(index, 'type', e.target.value)}
                      className="mt-1 block w-full rounded-md border-gray-300 py-2 pl-3 pr-10 text-base focus:border-primary-500 focus:outline-none focus:ring-primary-500 sm:text-sm border"
                    >
                      <option value="adult">Adult</option>
                      <option value="child">Child (50% off)</option>
                      <option value="infant">Infant (Free)</option>
                      <option value="senior">Senior (20% off)</option>
                    </select>
                  </div>
                  
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      Seat Number (Optional)
                    </label>
                    <input
                      type="text"
                      value={passenger.seat_number}
                      onChange={(e) => handlePassengerChange(index, 'seat_number', e.target.value)}
                      className="mt-1 block w-full rounded-md border-gray-300 px-3 py-2 shadow-sm focus:border-primary-500 focus:outline-none focus:ring-primary-500 sm:text-sm border"
                    />
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-white shadow sm:rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <h3 className="text-lg font-medium text-gray-900 mb-4">Payment Information</h3>
            
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700">
                  Payment Method
                </label>
                <select
                  value={paymentMethod}
                  onChange={(e) => setPaymentMethod(e.target.value)}
                  className="mt-1 block w-full rounded-md border-gray-300 py-2 pl-3 pr-10 text-base focus:border-primary-500 focus:outline-none focus:ring-primary-500 sm:text-sm border"
                >
                  <option value="credit_card">Credit Card</option>
                  <option value="debit_card">Debit Card</option>
                  <option value="cash">Cash (Pay at Terminal)</option>
                </select>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700">
                  Special Requirements (Optional)
                </label>
                <textarea
                  rows={3}
                  value={specialRequirements}
                  onChange={(e) => setSpecialRequirements(e.target.value)}
                  className="mt-1 block w-full rounded-md border-gray-300 px-3 py-2 shadow-sm focus:border-primary-500 focus:outline-none focus:ring-primary-500 sm:text-sm border"
                  placeholder="Any special requirements or requests..."
                />
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white shadow sm:rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <div className="flex justify-between items-center">
              <div>
                <h3 className="text-lg font-medium text-gray-900">Total Amount</h3>
                <p className="text-2xl font-bold text-primary-600">${calculateTotal().toFixed(2)}</p>
              </div>
              
              <button
                type="submit"
                disabled={loading}
                className="inline-flex items-center rounded-md bg-primary-600 px-6 py-3 text-base font-semibold text-white shadow-sm hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 disabled:opacity-50"
              >
                {loading ? 'Processing...' : 'Confirm Booking'}
              </button>
            </div>
            
            {error && (
              <div className="mt-4 rounded-md bg-red-50 p-4">
                <p className="text-sm text-red-800">{error}</p>
              </div>
            )}
          </div>
        </div>
      </form>
    </div>
  );
}