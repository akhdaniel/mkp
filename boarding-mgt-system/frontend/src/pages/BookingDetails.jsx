import React from 'react';
import { useParams } from 'react-router-dom';

export default function BookingDetails() {
  const { id } = useParams();
  
  return (
    <div className="max-w-3xl mx-auto">
      <h2 className="text-2xl font-bold text-gray-900 mb-6">Booking Details</h2>
      <div className="bg-white shadow sm:rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <p className="text-gray-500">Booking details for ID: {id}</p>
          <p className="text-sm text-gray-400 mt-2">Full booking details will be displayed here</p>
        </div>
      </div>
    </div>
  );
}