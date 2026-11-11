import React from 'react';

// The Funkey Grab & Bite wordmark, with a white-lettered version for dark mode.
const BrandLogo: React.FC<{ className?: string }> = ({ className = 'h-9' }) => (
  <span className={`inline-flex items-center ${className}`}>
    <img src="/brand/logo.png" alt="Funkey Grab & Bite" className="h-full w-auto dark:hidden" width={512} height={202} />
    <img src="/brand/logo-dark.png" alt="" aria-hidden="true" className="hidden h-full w-auto dark:block" width={512} height={202} />
  </span>
);

export default BrandLogo;
