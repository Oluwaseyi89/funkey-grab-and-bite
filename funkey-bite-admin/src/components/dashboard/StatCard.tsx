import React from 'react';
import { ArrowDownRight, ArrowUpRight, Minus } from 'lucide-react';
import type { LucideIcon } from 'lucide-react';

interface StatCardProps {
  label: string;
  value: string;
  icon: LucideIcon;
  // % change against the comparison, or null when there's nothing to compare.
  change?: number | null;
  comparison?: string;
  loading?: boolean;
}

// A headline number with its real change against a stated comparison.
const StatCard: React.FC<StatCardProps> = ({ label, value, icon: Icon, change, comparison, loading }) => {
  const up = (change ?? 0) > 0.5;
  const down = (change ?? 0) < -0.5;
  const Trend = up ? ArrowUpRight : down ? ArrowDownRight : Minus;

  return (
    <div className="card p-4">
      <div className="flex items-center justify-between">
        <p className="text-sm font-medium text-gray-500 dark:text-gray-400">{label}</p>
        <span className="flex h-9 w-9 items-center justify-center rounded-xl bg-primary-50 text-primary-600 dark:bg-primary-500/10 dark:text-primary-300">
          <Icon className="h-[18px] w-[18px]" aria-hidden="true" />
        </span>
      </div>
      {loading ? (
        <div className="mt-2 space-y-2">
          <div className="skeleton h-7 w-28" />
          <div className="skeleton h-4 w-36" />
        </div>
      ) : (
        <>
          <p className="mt-1 font-display text-2xl font-bold tabular-nums text-gray-900 dark:text-white">{value}</p>
          {comparison && (
            <p className="mt-1 flex items-center gap-1 text-xs text-gray-500 dark:text-gray-400">
              {change != null ? (
                <span className={`inline-flex items-center font-semibold ${up ? 'text-leaf-600' : down ? 'text-red-600' : 'text-gray-500'}`}>
                  <Trend className="h-3.5 w-3.5" aria-hidden="true" />
                  {Math.abs(change).toFixed(0)}%
                </span>
              ) : (
                <span className="font-semibold text-gray-400">No data</span>
              )}
              <span>{comparison}</span>
            </p>
          )}
        </>
      )}
    </div>
  );
};

export default StatCard;
