// src/pages/Orders/OrderCalendar.tsx
import React, { useState, useEffect } from 'react';
import { useQuery } from '@tanstack/react-query';
import { 
  Calendar as CalendarIcon, 
  ChevronLeft, 
  ChevronRight, 
  Filter, 
  Package, 
  ShoppingCart,
  Users,
  DollarSign,
  MoreVertical,
  Eye
} from 'lucide-react';
import { Link } from 'react-router-dom';
import { getOrders, getCateringRequests } from '../../api/adminApi';
import type { Order, CateringRequest, OrderStatus, CateringStatus } from '../../types';

const OrderCalendar: React.FC = () => {
  const [currentDate, setCurrentDate] = useState(new Date());
  const [view, setView] = useState<'month' | 'week' | 'day'>('month');
  const [selectedType, setSelectedType] = useState<'all' | 'orders' | 'catering'>('all');
  const [selectedEvent, setSelectedEvent] = useState<{ type: 'order' | 'catering', data: any } | null>(null);

  // Fetch orders
  const { data: ordersData, isLoading: ordersLoading } = useQuery({
    queryKey: ['calendar-orders'],
    queryFn: async () => {
      const data = await getOrders({ page: 1, limit: 1000 });
      return data.data || [];
    },
  });

  // Fetch catering requests
  const { data: cateringData, isLoading: cateringLoading } = useQuery({
    queryKey: ['calendar-catering'],
    queryFn: async () => {
      const data = await getCateringRequests({ page: 1, limit: 1000 });
      return data.data || [];
    },
  });

  const isLoading = ordersLoading || cateringLoading;
  const orders = ordersData || [];
  const cateringRequests = cateringData || [];

  // Navigation functions
  const goToPrevious = () => {
    const newDate = new Date(currentDate);
    if (view === 'month') {
      newDate.setMonth(newDate.getMonth() - 1);
    } else if (view === 'week') {
      newDate.setDate(newDate.getDate() - 7);
    } else {
      newDate.setDate(newDate.getDate() - 1);
    }
    setCurrentDate(newDate);
  };

  const goToNext = () => {
    const newDate = new Date(currentDate);
    if (view === 'month') {
      newDate.setMonth(newDate.getMonth() + 1);
    } else if (view === 'week') {
      newDate.setDate(newDate.getDate() + 7);
    } else {
      newDate.setDate(newDate.getDate() + 1);
    }
    setCurrentDate(newDate);
  };

  const goToToday = () => {
    setCurrentDate(new Date());
  };

  // Get events for a specific date
  const getEventsForDate = (date: Date) => {
    const dateStr = date.toISOString().split('T')[0];
    
    const events: any[] = [];
    
    // Get orders for the date
    if (selectedType === 'all' || selectedType === 'orders') {
      const dayOrders = orders.filter(order => {
        const orderDate = new Date(order.createdAt).toISOString().split('T')[0];
        return orderDate === dateStr;
      });
      
      dayOrders.forEach(order => {
        events.push({
          type: 'order' as const,
          data: order,
          time: new Date(order.createdAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
          title: `Order #${order.orderNumber}`,
          subtitle: order.customerName,
          amount: order.totalAmount,
          status: order.status,
        });
      });
    }
    
    // Get catering events for the date
    if (selectedType === 'all' || selectedType === 'catering') {
      const dayCatering = cateringRequests.filter(catering => {
        const eventDate = new Date(catering.eventDate).toISOString().split('T')[0];
        return eventDate === dateStr;
      });
      
      dayCatering.forEach(catering => {
        events.push({
          type: 'catering' as const,
          data: catering,
          time: catering.eventTime || 'All Day',
          title: catering.eventName || 'Catering Event',
          subtitle: catering.contactName,
          guests: catering.guestCount,
          status: catering.status,
        });
      });
    }
    
    return events.sort((a, b) => {
      if (a.time === 'All Day') return -1;
      if (b.time === 'All Day') return 1;
      return a.time.localeCompare(b.time);
    });
  };

  // Get days in month
  const getDaysInMonth = (year: number, month: number) => {
    return new Date(year, month + 1, 0).getDate();
  };

  // Get first day of month
  const getFirstDayOfMonth = (year: number, month: number) => {
    return new Date(year, month, 1).getDay();
  };

  // Generate calendar days
  const generateCalendarDays = () => {
    const year = currentDate.getFullYear();
    const month = currentDate.getMonth();
    const daysInMonth = getDaysInMonth(year, month);
    const firstDay = getFirstDayOfMonth(year, month);
    
    const days = [];
    
    // Previous month days
    const prevMonthDays = getDaysInMonth(year, month - 1);
    for (let i = firstDay - 1; i >= 0; i--) {
      const day = new Date(year, month - 1, prevMonthDays - i);
      days.push({
        date: day,
        isCurrentMonth: false,
        events: getEventsForDate(day),
      });
    }
    
    // Current month days
    for (let i = 1; i <= daysInMonth; i++) {
      const day = new Date(year, month, i);
      days.push({
        date: day,
        isCurrentMonth: true,
        isToday: day.toDateString() === new Date().toDateString(),
        events: getEventsForDate(day),
      });
    }
    
    // Next month days (to fill 42 slots total)
    const totalCells = 42; // 6 weeks
    const nextMonthDays = totalCells - days.length;
    for (let i = 1; i <= nextMonthDays; i++) {
      const day = new Date(year, month + 1, i);
      days.push({
        date: day,
        isCurrentMonth: false,
        events: getEventsForDate(day),
      });
    }
    
    return days;
  };

  const getStatusColor = (status: OrderStatus | CateringStatus) => {
    switch (status) {
      case 'completed':
      case 'confirmed':
        return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200';
      case 'cancelled':
      case 'declined':
        return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200';
      case 'pending':
        return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200';
      case 'preparing':
        return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200';
      case 'ready':
        return 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200';
      default:
        return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200';
    }
  };

  const getEventIcon = (type: 'order' | 'catering') => {
    return type === 'order' ? Package : Users;
  };

  const calendarDays = generateCalendarDays();
  const weekDays = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
  const monthYear = currentDate.toLocaleDateString('en-US', { month: 'long', year: 'numeric' });

  // Get today's events summary
  const todayEvents = getEventsForDate(new Date());
  const todaySummary = {
    orders: todayEvents.filter(e => e.type === 'order').length,
    catering: todayEvents.filter(e => e.type === 'catering').length,
    revenue: todayEvents
      .filter(e => e.type === 'order')
      .reduce((sum, event) => sum + (event.amount || 0), 0),
  };

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div className="h-8 bg-gray-200 dark:bg-gray-700 rounded w-48 animate-pulse"></div>
          <div className="flex items-center space-x-2">
            {[...Array(3)].map((_, i) => (
              <div key={i} className="h-10 w-24 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
            ))}
          </div>
        </div>
        <div className="h-96 bg-gray-200 dark:bg-gray-700 rounded-xl animate-pulse"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Order Calendar</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Visualize orders and catering events by date
          </p>
        </div>
        
        <div className="flex items-center space-x-3">
          <button
            onClick={goToToday}
            className="px-4 py-2 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 rounded-lg transition-colors"
          >
            Today
          </button>
          <div className="flex items-center border border-gray-300 dark:border-gray-600 rounded-lg overflow-hidden">
            <button
              onClick={goToPrevious}
              className="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
            >
              <ChevronLeft className="h-4 w-4" />
            </button>
            <div className="px-4 py-2 text-center min-w-48">
              <span className="font-semibold">{monthYear}</span>
            </div>
            <button
              onClick={goToNext}
              className="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
            >
              <ChevronRight className="h-4 w-4" />
            </button>
          </div>
        </div>
      </div>

      {/* Today's Summary */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">Today's Orders</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white mt-1">
                {todaySummary.orders}
              </p>
            </div>
            <div className="h-12 w-12 rounded-lg bg-blue-100 dark:bg-blue-900 flex items-center justify-center">
              <Package className="h-6 w-6 text-blue-600 dark:text-blue-400" />
            </div>
          </div>
        </div>
        
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">Today's Catering</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white mt-1">
                {todaySummary.catering}
              </p>
            </div>
            <div className="h-12 w-12 rounded-lg bg-green-100 dark:bg-green-900 flex items-center justify-center">
              <Users className="h-6 w-6 text-green-600 dark:text-green-400" />
            </div>
          </div>
        </div>
        
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">Today's Revenue</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white mt-1">
                ${todaySummary.revenue.toFixed(2)}
              </p>
            </div>
            <div className="h-12 w-12 rounded-lg bg-primary-100 dark:bg-primary-900 flex items-center justify-center">
              <DollarSign className="h-6 w-6 text-primary-600 dark:text-primary-400" />
            </div>
          </div>
        </div>
      </div>

      {/* Filters */}
      <div className="bg-white dark:bg-gray-800 rounded-xl p-4 border border-gray-200 dark:border-gray-700">
        <div className="flex flex-col lg:flex-row items-center justify-between gap-4">
          <div className="flex items-center space-x-4">
            <Filter className="h-5 w-5 text-gray-400" />
            <div className="flex items-center space-x-2">
              {(['all', 'orders', 'catering'] as const).map((type) => (
                <button
                  key={type}
                  onClick={() => setSelectedType(type)}
                  className={`px-4 py-2 rounded-lg capitalize ${
                    selectedType === type
                      ? 'bg-primary-500 text-white'
                      : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
                  }`}
                >
                  {type}
                </button>
              ))}
            </div>
          </div>
          
          <div className="flex items-center space-x-2">
            {(['month', 'week', 'day'] as const).map((viewType) => (
              <button
                key={viewType}
                onClick={() => setView(viewType)}
                className={`px-4 py-2 rounded-lg capitalize ${
                  view === viewType
                    ? 'bg-primary-500 text-white'
                    : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
                }`}
              >
                {viewType}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Calendar */}
      <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden">
        {/* Weekday Headers */}
        <div className="grid grid-cols-7 border-b border-gray-200 dark:border-gray-700">
          {weekDays.map((day) => (
            <div key={day} className="p-4 text-center font-medium text-gray-500 dark:text-gray-400">
              {day}
            </div>
          ))}
        </div>

        {/* Calendar Days */}
        <div className="grid grid-cols-7">
          {calendarDays.map((day, index) => (
            <div
              key={index}
              className={`min-h-32 p-2 border border-gray-200 dark:border-gray-700 ${
                day.isCurrentMonth ? 'bg-white dark:bg-gray-800' : 'bg-gray-50 dark:bg-gray-900/50'
              } ${day.isToday ? 'ring-2 ring-primary-500' : ''}`}
            >
              <div className="flex justify-between items-center mb-2">
                <span className={`font-medium ${
                  day.isCurrentMonth 
                    ? 'text-gray-900 dark:text-white' 
                    : 'text-gray-400 dark:text-gray-500'
                } ${day.isToday ? 'text-primary-600 dark:text-primary-400' : ''}`}>
                  {day.date.getDate()}
                </span>
                {day.events.length > 0 && (
                  <span className="h-5 w-5 rounded-full bg-primary-500 text-white text-xs flex items-center justify-center">
                    {day.events.length}
                  </span>
                )}
              </div>
              
              <div className="space-y-1 max-h-24 overflow-y-auto">
                {day.events.slice(0, 3).map((event, eventIndex) => {
                  const Icon = getEventIcon(event.type);
                  return (
                    <div
                      key={eventIndex}
                      onClick={() => setSelectedEvent(event)}
                      className="p-2 rounded text-xs cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors truncate"
                      title={`${event.time} - ${event.title}`}
                    >
                      <div className="flex items-center space-x-1">
                        <Icon className="h-3 w-3 flex-shrink-0" />
                        <span className="font-medium truncate">{event.title}</span>
                      </div>
                      <div className="flex items-center justify-between mt-1">
                        <span className="text-gray-500 truncate">{event.subtitle}</span>
                        <span className={`px-1 rounded text-xs ${getStatusColor(event.status)}`}>
                          {event.status}
                        </span>
                      </div>
                    </div>
                  );
                })}
                {day.events.length > 3 && (
                  <div className="text-xs text-gray-500 text-center">
                    +{day.events.length - 3} more
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Legend */}
      <div className="bg-white dark:bg-gray-800 rounded-xl p-4 border border-gray-200 dark:border-gray-700">
        <div className="flex flex-wrap items-center gap-4">
          <div className="flex items-center space-x-2">
            <div className="h-4 w-4 rounded bg-blue-100 dark:bg-blue-900 border border-blue-300 dark:border-blue-700"></div>
            <span className="text-sm text-gray-600 dark:text-gray-400">Orders</span>
          </div>
          <div className="flex items-center space-x-2">
            <div className="h-4 w-4 rounded bg-green-100 dark:bg-green-900 border border-green-300 dark:border-green-700"></div>
            <span className="text-sm text-gray-600 dark:text-gray-400">Catering</span>
          </div>
          <div className="flex items-center space-x-2">
            <div className="h-4 w-4 rounded-full bg-primary-500"></div>
            <span className="text-sm text-gray-600 dark:text-gray-400">Today</span>
          </div>
        </div>
      </div>

      {/* Event Details Modal */}
      {selectedEvent && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white dark:bg-gray-800 rounded-xl max-w-md w-full max-h-[90vh] overflow-y-auto">
            <div className="p-6">
              <div className="flex items-center justify-between mb-6">
                <h2 className="text-xl font-bold text-gray-900 dark:text-white">
                  {selectedEvent.type === 'order' ? 'Order Details' : 'Catering Details'}
                </h2>
                <button
                  onClick={() => setSelectedEvent(null)}
                  className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
                >
                  ✕
                </button>
              </div>
              
              {selectedEvent.type === 'order' ? (
                <div className="space-y-4">
                  <div>
                    <p className="text-sm text-gray-500 dark:text-gray-400">Order Number</p>
                    <p className="font-semibold">#{selectedEvent.data.orderNumber}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-500 dark:text-gray-400">Customer</p>
                    <p className="font-semibold">{selectedEvent.data.customerName}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-500 dark:text-gray-400">Total Amount</p>
                    <p className="font-semibold text-primary-600">
                      ${selectedEvent.data.totalAmount.toFixed(2)}
                    </p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-500 dark:text-gray-400">Status</p>
                    <span className={`px-3 py-1 rounded-full text-sm ${getStatusColor(selectedEvent.data.status)}`}>
                      {selectedEvent.data.status}
                    </span>
                  </div>
                  <div className="pt-4 border-t border-gray-200 dark:border-gray-700">
                    <Link
                      to={`/orders/${selectedEvent.data.id}`}
                      className="flex items-center justify-center space-x-2 w-full px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors"
                    >
                      <Eye className="h-4 w-4" />
                      <span>View Full Details</span>
                    </Link>
                  </div>
                </div>
              ) : (
                <div className="space-y-4">
                  <div>
                    <p className="text-sm text-gray-500 dark:text-gray-400">Event</p>
                    <p className="font-semibold">{selectedEvent.data.eventName || 'Catering Event'}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-500 dark:text-gray-400">Contact Person</p>
                    <p className="font-semibold">{selectedEvent.data.contactName}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-500 dark:text-gray-400">Date & Time</p>
                    <p className="font-semibold">
                      {new Date(selectedEvent.data.eventDate).toLocaleDateString()} •{' '}
                      {selectedEvent.data.eventTime || 'All Day'}
                    </p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-500 dark:text-gray-400">Guests</p>
                    <p className="font-semibold">{selectedEvent.data.guestCount} people</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-500 dark:text-gray-400">Status</p>
                    <span className={`px-3 py-1 rounded-full text-sm ${getStatusColor(selectedEvent.data.status)}`}>
                      {selectedEvent.data.status}
                    </span>
                  </div>
                  <div className="pt-4 border-t border-gray-200 dark:border-gray-700">
                    <Link
                      to={`/catering/${selectedEvent.data.id}`}
                      className="flex items-center justify-center space-x-2 w-full px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors"
                    >
                      <Eye className="h-4 w-4" />
                      <span>View Catering Details</span>
                    </Link>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default OrderCalendar;