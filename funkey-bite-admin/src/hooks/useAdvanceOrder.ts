import { useMutation, useQueryClient } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { updateOrderStatus } from '../api/adminApi';
import { ORDER_STATUS } from '../components/ui/StatusPill';

// Moves an order to a new status and refreshes every list that shows it.
export const useAdvanceOrder = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, status }: { id: number; status: string; orderNumber?: string }) => updateOrderStatus(id, { status }),
    onSuccess: (_data, { status, orderNumber }) => {
      toast.success(`${orderNumber ? `#${orderNumber}` : 'Order'}: ${ORDER_STATUS[status]?.label ?? status}`);
      ['orders', 'order', 'recent-orders', 'dashboard-overview', 'calendar-orders'].forEach((key) =>
        queryClient.invalidateQueries({ queryKey: [key] }),
      );
    },
    onError: (error: Error) => {
      toast.error(error.message || 'The order wasn’t updated. Please try again.');
    },
  });
};
