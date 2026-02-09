import React from 'react';
import type { LucideIcon } from 'lucide-react';

interface StatCardProps {
  title: string;
  value: string | number;
  icon: LucideIcon;
  change?: number;
  color: 'primary' | 'blue' | 'green' | 'yellow' | 'red' | 'purple';
  loading?: boolean;
}

const StatCard: React.FC<StatCardProps> = ({ 
  title, 
  value, 
  icon: Icon, 
  change, 
  color,
  loading = false 
}) => {
  const colorClasses = {
    primary: 'bg-primary-500',
    blue: 'bg-blue-500',
    green: 'bg-green-500',
    yellow: 'bg-yellow-500',
    red: 'bg-red-500',
    purple: 'bg-purple-500',
  };

  const textColorClasses = {
    primary: 'text-primary-500',
    blue: 'text-blue-500',
    green: 'text-green-500',
    yellow: 'text-yellow-500',
    red: 'text-red-500',
    purple: 'text-purple-500',
  };

  const bgColorClasses = {
    primary: 'bg-primary-50 dark:bg-primary-900/20',
    blue: 'bg-blue-50 dark:bg-blue-900/20',
    green: 'bg-green-50 dark:bg-green-900/20',
    yellow: 'bg-yellow-50 dark:bg-yellow-900/20',
    red: 'bg-red-50 dark:bg-red-900/20',
    purple: 'bg-purple-50 dark:bg-purple-900/20',
  };

  if (loading) {
    return (
      <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700 animate-pulse">
        <div className="flex items-center justify-between">
          <div className="space-y-2">
            <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-24"></div>
            <div className="h-8 bg-gray-200 dark:bg-gray-700 rounded w-32"></div>
          </div>
          <div className="h-12 w-12 rounded-lg bg-gray-200 dark:bg-gray-700"></div>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700 hover:shadow-md transition-shadow">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm text-gray-600 dark:text-gray-400 mb-1">{title}</p>
          <p className="text-2xl font-bold text-gray-900 dark:text-white">{value}</p>
          
          {change !== undefined && (
            <div className="flex items-center mt-2">
              <span className={`text-sm font-medium ${change >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                {change >= 0 ? '+' : ''}{change}%
              </span>
              <span className="text-sm text-gray-500 dark:text-gray-400 ml-2">from yesterday</span>
            </div>
          )}
        </div>
        
        <div className={`${bgColorClasses[color]} p-3 rounded-lg`}>
          <Icon className={`h-6 w-6 ${textColorClasses[color]}`} />
        </div>
      </div>
    </div>
  );
};

export default StatCard;