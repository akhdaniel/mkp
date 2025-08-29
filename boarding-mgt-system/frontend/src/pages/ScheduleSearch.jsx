import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { format } from 'date-fns';
import { MagnifyingGlassIcon, ClockIcon, MapPinIcon } from '@heroicons/react/24/outline';
import { portsAPI } from '../api/ports';
import { bookingAPI } from '../api/booking';

export default function ScheduleSearch() {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [ports, setPorts] = useState([]);
  const [schedules, setSchedules] = useState([]);
  const [searched, setSearched] = useState(false);
  const [formData, setFormData] = useState({
    departure_port_id: '',
    arrival_port_id: '',
    departure_date: format(new Date(), 'yyyy-MM-dd'),
    passenger_count: 1,
  });

  useEffect(() => {
    loadPorts();
  }, []);

  const loadPorts = async () => {
    try {
      const data = await portsAPI.getAllPorts({ limit: 100 });
      setPorts(data.ports || []);
    } catch (error) {
      console.error('Failed to load ports:', error);
    }
  };

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSearch = async (e) => {
    e.preventDefault();
    setLoading(true);
    setSearched(true);

    try {
      const data = await bookingAPI.searchSchedules(formData);
      setSchedules(data.schedules || []);
    } catch (error) {
      console.error('Search failed:', error);
      setSchedules([]);
    } finally {
      setLoading(false);
    }
  };

  const handleSelectSchedule = (schedule) => {
    navigate('/booking/new', {
      state: {
        schedule,
        passengerCount: formData.passenger_count,
      },
    });
  };

  return (
    <div className="max-w-7xl mx-auto">
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <h3 className="text-lg font-medium leading-6 text-gray-900 mb-4">
            Search Ferry Schedules
          </h3>
          
          <form onSubmit={handleSearch} className="space-y-4">
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
              <div>
                <label htmlFor="departure_port_id" className="block text-sm font-medium text-gray-700">
                  From
                </label>
                <select
                  id="departure_port_id"
                  name="departure_port_id"
                  required
                  value={formData.departure_port_id}
                  onChange={handleChange}
                  className="mt-1 block w-full rounded-md border-gray-300 py-2 pl-3 pr-10 text-base focus:border-primary-500 focus:outline-none focus:ring-primary-500 sm:text-sm border"
                >
                  <option value="">Select departure port</option>
                  {ports.map((port) => (
                    <option key={port.id} value={port.id}>
                      {port.name} ({port.code})
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label htmlFor="arrival_port_id" className="block text-sm font-medium text-gray-700">
                  To
                </label>
                <select
                  id="arrival_port_id"
                  name="arrival_port_id"
                  required
                  value={formData.arrival_port_id}
                  onChange={handleChange}
                  className="mt-1 block w-full rounded-md border-gray-300 py-2 pl-3 pr-10 text-base focus:border-primary-500 focus:outline-none focus:ring-primary-500 sm:text-sm border"
                >
                  <option value="">Select arrival port</option>
                  {ports.map((port) => (
                    <option key={port.id} value={port.id}>
                      {port.name} ({port.code})
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label htmlFor="departure_date" className="block text-sm font-medium text-gray-700">
                  Date
                </label>
                <input
                  type="date"
                  id="departure_date"
                  name="departure_date"
                  required
                  min={format(new Date(), 'yyyy-MM-dd')}
                  value={formData.departure_date}
                  onChange={handleChange}
                  className="mt-1 block w-full rounded-md border-gray-300 px-3 py-2 shadow-sm focus:border-primary-500 focus:outline-none focus:ring-primary-500 sm:text-sm border"
                />
              </div>

              <div>
                <label htmlFor="passenger_count" className="block text-sm font-medium text-gray-700">
                  Passengers
                </label>
                <input
                  type="number"
                  id="passenger_count"
                  name="passenger_count"
                  min="1"
                  max="10"
                  required
                  value={formData.passenger_count}
                  onChange={handleChange}
                  className="mt-1 block w-full rounded-md border-gray-300 px-3 py-2 shadow-sm focus:border-primary-500 focus:outline-none focus:ring-primary-500 sm:text-sm border"
                />
              </div>
            </div>

            <div className="flex justify-end">
              <button
                type="submit"
                disabled={loading}
                className="inline-flex items-center rounded-md bg-primary-600 px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 disabled:opacity-50"
              >
                <MagnifyingGlassIcon className="-ml-0.5 mr-1.5 h-5 w-5" aria-hidden="true" />
                {loading ? 'Searching...' : 'Search'}
              </button>
            </div>
          </form>
        </div>
      </div>

      {searched && (
        <div className="mt-8">
          <h2 className="text-lg font-medium text-gray-900 mb-4">
            {schedules.length > 0 ? 'Available Schedules' : 'No schedules found'}
          </h2>

          {schedules.length > 0 && (
            <div className="bg-white shadow overflow-hidden sm:rounded-md">
              <ul className="divide-y divide-gray-200">
                {schedules.map((schedule) => (
                  <li key={schedule.id}>
                    <div className="px-4 py-4 sm:px-6 hover:bg-gray-50 cursor-pointer" onClick={() => handleSelectSchedule(schedule)}>
                      <div className="flex items-center justify-between">
                        <div className="flex-1">
                          <div className="flex items-center justify-between">
                            <div className="flex items-center">
                              <MapPinIcon className="h-5 w-5 text-gray-400 mr-2" />
                              <p className="text-sm font-medium text-gray-900">
                                Route: {schedule.route?.name || 'Direct Route'}
                              </p>
                            </div>
                            <div className="flex items-center">
                              <ClockIcon className="h-5 w-5 text-gray-400 mr-2" />
                              <p className="text-sm text-gray-500">
                                {format(new Date(`2000-01-01T${schedule.departure_time}`), 'HH:mm')} - 
                                {format(new Date(`2000-01-01T${schedule.arrival_time}`), 'HH:mm')}
                              </p>
                            </div>
                          </div>
                          <div className="mt-2 flex items-center justify-between">
                            <div className="flex items-center text-sm text-gray-500">
                              <p>
                                Vessel: {schedule.vessel?.name || 'Ferry'} | 
                                Available Seats: {schedule.available_seats}
                              </p>
                            </div>
                            <div className="flex items-center">
                              <p className="text-lg font-semibold text-primary-600">
                                ${schedule.base_price}
                              </p>
                              <span className="ml-1 text-sm text-gray-500">per person</span>
                            </div>
                          </div>
                        </div>
                        <div className="ml-4">
                          <button className="inline-flex items-center rounded-md bg-primary-100 px-3 py-2 text-sm font-semibold text-primary-700 hover:bg-primary-200">
                            Select
                          </button>
                        </div>
                      </div>
                    </div>
                  </li>
                ))}
              </ul>
            </div>
          )}
        </div>
      )}
    </div>
  );
}