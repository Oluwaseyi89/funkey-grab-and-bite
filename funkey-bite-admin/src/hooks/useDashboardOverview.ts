import { useQuery } from '@tanstack/react-query';
import { getDashboardOverview } from '../api/adminApi';
import { useAuthStore } from '../stores/authStore';
import type { DashboardOverview } from '../types';

// Today's figures and what needs doing. Shared by the dashboard and the
// navigation badges; live events refresh it, and it re-checks every minute.
export const useDashboardOverview = () => {
  const signedIn = useAuthStore((state) => state.isAuthenticated);
  return useQuery<DashboardOverview>({
    queryKey: ['dashboard-overview'],
    queryFn: getDashboardOverview,
    enabled: signedIn,
    refetchInterval: 60_000,
    staleTime: 30_000,
  });
};

export type BadgeKey = 'newOrders' | 'paymentsToReview' | 'lowStock' | 'cateringToQuote';

export const badgeCounts = (overview?: DashboardOverview): Record<BadgeKey, number> => ({
  newOrders: overview?.queue.pending ?? 0,
  paymentsToReview: overview?.paymentsToReview ?? 0,
  lowStock: overview?.lowStock ?? 0,
  cateringToQuote: overview?.cateringToQuote ?? 0,
});
