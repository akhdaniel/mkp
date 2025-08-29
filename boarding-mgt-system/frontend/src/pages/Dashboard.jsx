import React from 'react';
import { UsersIcon, TicketIcon, CurrencyDollarIcon, CalendarIcon } from '@heroicons/react/24/outline';

export default function Dashboard() {
  const stats = [
    { name: 'Total Bookings', value: '1,234', icon: TicketIcon },
    { name: 'Total Revenue', value: '$45,678', icon: CurrencyDollarIcon },
    { name: 'Active Schedules', value: '42', icon: CalendarIcon },
    { name: 'Registered Users', value: '892', icon: UsersIcon },
  ];

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900 mb-6">Admin Dashboard</h2>
      
      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => (
          <div key={stat.name} className="bg-white overflow-hidden shadow rounded-lg">
            <div className="p-5">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <stat.icon className="h-6 w-6 text-gray-400" aria-hidden="true" />
                </div>
                <div className="ml-5 w-0 flex-1">
                  <dl>
                    <dt className="text-sm font-medium text-gray-500 truncate">{stat.name}</dt>
                    <dd className="text-lg font-semibold text-gray-900">{stat.value}</dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="mt-8 bg-white shadow sm:rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Recent Activity</h3>
          <p className="text-sm text-gray-500">
            Dashboard features coming soon. This will include booking management, schedule management, and reporting tools.
          </p>
        </div>
      </div>
    </div>
  );
}