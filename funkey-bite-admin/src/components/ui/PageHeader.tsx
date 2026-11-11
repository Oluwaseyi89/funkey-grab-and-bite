import React from 'react';

interface PageHeaderProps {
  title: string;
  description?: React.ReactNode;
  // Buttons on the right (below the title on phones).
  actions?: React.ReactNode;
}

// The title row at the top of each page.
const PageHeader: React.FC<PageHeaderProps> = ({ title, description, actions }) => (
  <div className="mb-5 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
    <div className="min-w-0">
      <h1 className="text-2xl font-bold text-gray-900 dark:text-white md:text-3xl">{title}</h1>
      {description && <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">{description}</p>}
    </div>
    {actions && <div className="flex flex-wrap items-center gap-2">{actions}</div>}
  </div>
);

export default PageHeader;
