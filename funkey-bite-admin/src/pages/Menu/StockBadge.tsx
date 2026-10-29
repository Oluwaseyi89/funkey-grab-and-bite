import React from 'react';
import { Link } from 'react-router-dom';
import type { MenuItem } from '../../types';

// Below this, stock is flagged as low.
const LOW_STOCK_THRESHOLD = 5;

const stockBadge = (item: MenuItem) => {
  const stock = item.stockQuantity;
  if (stock == null) {
    return { label: 'Not tracked', className: 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300' };
  }
  if (stock <= 0) {
    return { label: 'Out of stock', className: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200' };
  }
  if (stock < LOW_STOCK_THRESHOLD) {
    return { label: `Low: ${stock} left`, className: 'bg-amber-100 text-amber-800 dark:bg-amber-900 dark:text-amber-200' };
  }
  return { label: `${stock} in stock`, className: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200' };
};

// Stock is managed in Inventory, so the badge links there.
const StockBadge: React.FC<{ item: MenuItem }> = ({ item }) => {
  const badge = stockBadge(item);
  return (
    <Link
      to="/inventory"
      title="Manage stock in Inventory"
      className={`inline-block px-2 py-1 text-xs font-medium rounded-full hover:opacity-80 ${badge.className}`}
    >
      {badge.label}
    </Link>
  );
};

export default StockBadge;
