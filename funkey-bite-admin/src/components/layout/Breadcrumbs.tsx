import React from 'react';
import { Link, useLocation } from 'react-router-dom';

// Labels for path segments; anything else is title-cased.
const SEGMENT_LABELS: Record<string, string> = {
  dashboard: 'Dashboard',
  orders: 'Orders',
  calendar: 'Calendar',
  menu: 'Menu',
  categories: 'Categories',
  new: 'New Item',
  edit: 'Edit',
  customers: 'Customers',
  catering: 'Catering',
  packages: 'Packages',
  team: 'Team',
  inventory: 'Inventory',
  alerts: 'Alerts',
  promotions: 'Promotions',
  reports: 'Reports',
  settings: 'Settings',
  profile: 'Profile',
  admins: 'Admin Users',
};

// Sections whose records have their own page at /<section>/:id. Others (e.g.
// promotions) only have an edit page, so their ID crumb is not a link.
const SECTIONS_WITH_DETAIL_PAGES = new Set(['orders', 'menu', 'customers', 'catering']);

const labelFor = (segment: string) => {
  if (/^\d+$/.test(segment)) return `#${segment}`;
  return SEGMENT_LABELS[segment] ?? segment.charAt(0).toUpperCase() + segment.slice(1).replace(/-/g, ' ');
};

// Shows where the admin is, e.g. "Menu / #9 / Edit", from the current path.
const Breadcrumbs: React.FC = () => {
  const { pathname } = useLocation();
  const segments = pathname.split('/').filter(Boolean);

  const crumbs = segments.map((segment, index) => {
    const isId = /^\d+$/.test(segment);
    const linkable = index === 0 || (isId && index === 1 && SECTIONS_WITH_DETAIL_PAGES.has(segments[0]));
    return {
      label: labelFor(segment),
      href: '/' + segments.slice(0, index + 1).join('/'),
      linkable,
    };
  });

  if (crumbs.length === 0) {
    crumbs.push({ label: 'Dashboard', href: '/dashboard', linkable: true });
  }

  return (
    <nav aria-label="Breadcrumb" className="text-sm text-gray-500 dark:text-gray-400">
      <ol className="flex flex-wrap items-center gap-1">
        {crumbs.map((crumb, index) => {
          const isLast = index === crumbs.length - 1;
          return (
            <li key={crumb.href} className="flex items-center gap-1">
              {index > 0 && <span aria-hidden="true">/</span>}
              {isLast ? (
                <span aria-current="page" className="text-gray-700 dark:text-gray-200">
                  {crumb.label}
                </span>
              ) : crumb.linkable ? (
                <Link to={crumb.href} className="hover:text-primary-500 transition-colors">
                  {crumb.label}
                </Link>
              ) : (
                <span>{crumb.label}</span>
              )}
            </li>
          );
        })}
      </ol>
    </nav>
  );
};

export default Breadcrumbs;
