import React from 'react';
import { Link } from 'react-router-dom';
import { ArrowRightIcon, TicketIcon, CalendarIcon, ShieldCheckIcon } from '@heroicons/react/24/outline';

export default function Home() {
  return (
    <div className="relative">
      {/* Hero section */}
      <div className="relative bg-white pb-16 pt-8 sm:pb-24">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="text-center">
            <h1 className="text-4xl font-bold tracking-tight text-gray-900 sm:text-5xl md:text-6xl">
              <span className="block">Ferry Booking Made</span>
              <span className="block text-primary-600">Simple and Secure</span>
            </h1>
            <p className="mx-auto mt-3 max-w-md text-base text-gray-500 sm:text-lg md:mt-5 md:max-w-3xl md:text-xl">
              Book your ferry tickets online with FerryFlow. Search schedules, compare prices, and secure your seats in minutes.
            </p>
            <div className="mx-auto mt-5 max-w-md sm:flex sm:justify-center md:mt-8">
              <div className="rounded-md shadow">
                <Link
                  to="/search"
                  className="flex w-full items-center justify-center rounded-md bg-primary-600 px-8 py-3 text-base font-medium text-white hover:bg-primary-700 md:px-10 md:py-4 md:text-lg"
                >
                  Search Schedules
                  <ArrowRightIcon className="ml-2 h-5 w-5" />
                </Link>
              </div>
              <div className="mt-3 rounded-md shadow sm:ml-3 sm:mt-0">
                <Link
                  to="/register"
                  className="flex w-full items-center justify-center rounded-md bg-white px-8 py-3 text-base font-medium text-primary-600 hover:bg-gray-50 md:px-10 md:py-4 md:text-lg border border-primary-600"
                >
                  Get Started
                </Link>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Features section */}
      <div className="bg-gray-50 py-16">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="text-center">
            <h2 className="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
              Why Choose FerryFlow?
            </h2>
            <p className="mt-4 text-lg text-gray-600">
              Everything you need for a smooth ferry journey
            </p>
          </div>

          <div className="mt-12 grid gap-8 sm:grid-cols-2 lg:grid-cols-3">
            <div className="bg-white rounded-lg p-6 shadow">
              <div className="flex items-center justify-center h-12 w-12 rounded-md bg-primary-500 text-white">
                <CalendarIcon className="h-6 w-6" />
              </div>
              <h3 className="mt-4 text-lg font-medium text-gray-900">Easy Scheduling</h3>
              <p className="mt-2 text-base text-gray-500">
                Browse all available ferry schedules and routes in one place. Filter by date, time, and destination.
              </p>
            </div>

            <div className="bg-white rounded-lg p-6 shadow">
              <div className="flex items-center justify-center h-12 w-12 rounded-md bg-primary-500 text-white">
                <TicketIcon className="h-6 w-6" />
              </div>
              <h3 className="mt-4 text-lg font-medium text-gray-900">Digital Tickets</h3>
              <p className="mt-2 text-base text-gray-500">
                Receive instant QR code tickets on your phone. No need to print - just scan and board.
              </p>
            </div>

            <div className="bg-white rounded-lg p-6 shadow">
              <div className="flex items-center justify-center h-12 w-12 rounded-md bg-primary-500 text-white">
                <ShieldCheckIcon className="h-6 w-6" />
              </div>
              <h3 className="mt-4 text-lg font-medium text-gray-900">Secure Payment</h3>
              <p className="mt-2 text-base text-gray-500">
                Your payment information is encrypted and secure. Multiple payment options available.
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Stats section */}
      <div className="bg-white py-16">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="text-center">
            <h2 className="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
              Trusted by Thousands
            </h2>
            <div className="mt-10 grid grid-cols-2 gap-8 md:grid-cols-4">
              <div>
                <p className="text-4xl font-bold text-primary-600">10+</p>
                <p className="mt-2 text-base text-gray-500">Ferry Operators</p>
              </div>
              <div>
                <p className="text-4xl font-bold text-primary-600">50+</p>
                <p className="mt-2 text-base text-gray-500">Routes Available</p>
              </div>
              <div>
                <p className="text-4xl font-bold text-primary-600">1000+</p>
                <p className="mt-2 text-base text-gray-500">Daily Bookings</p>
              </div>
              <div>
                <p className="text-4xl font-bold text-primary-600">99%</p>
                <p className="mt-2 text-base text-gray-500">Customer Satisfaction</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* CTA section */}
      <div className="bg-primary-700">
        <div className="mx-auto max-w-2xl px-4 py-16 text-center sm:px-6 sm:py-20 lg:px-8">
          <h2 className="text-3xl font-bold tracking-tight text-white sm:text-4xl">
            <span className="block">Ready to book your next journey?</span>
          </h2>
          <p className="mt-4 text-lg leading-6 text-primary-200">
            Join thousands of travelers who trust FerryFlow for their ferry bookings.
          </p>
          <Link
            to="/register"
            className="mt-8 inline-flex w-full items-center justify-center rounded-md border border-transparent bg-white px-5 py-3 text-base font-medium text-primary-600 hover:bg-primary-50 sm:w-auto"
          >
            Sign up for free
          </Link>
        </div>
      </div>
    </div>
  );
}