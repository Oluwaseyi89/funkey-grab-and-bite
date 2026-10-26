import React, { useState } from 'react';
import { X } from 'lucide-react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { restockItem, updateStock } from '../../api/adminApi';
import type { InventoryItem } from '../../types';

export type StockAdjustMode = 'add' | 'set' | 'subtract';

const MODES: { value: StockAdjustMode; label: string; help: string }[] = [
  { value: 'add', label: 'Restock', help: 'Add newly received stock to the current amount.' },
  { value: 'set', label: 'Set stock', help: 'Replace the current amount after a count. Use 0 to mark it sold out.' },
  { value: 'subtract', label: 'Remove', help: 'Take out spoiled, wasted or miscounted stock.' },
];

const DEFAULT_REASON: Record<StockAdjustMode, string> = {
  add: 'Restocked',
  set: 'Stock count',
  subtract: 'Stock removed',
};

interface StockAdjustDialogProps {
  item: InventoryItem;
  initialMode?: StockAdjustMode;
  onClose: () => void;
}

// Changes an item's stock through the API, then refreshes every view that
// shows stock (inventory, low stock, alerts and the menu).
const StockAdjustDialog: React.FC<StockAdjustDialogProps> = ({ item, initialMode = 'add', onClose }) => {
  const queryClient = useQueryClient();
  const [mode, setMode] = useState<StockAdjustMode>(initialMode);
  const [quantity, setQuantity] = useState('');
  const [reason, setReason] = useState('');

  const amount = Number(quantity);
  const isWholeNumber = quantity.trim() !== '' && Number.isInteger(amount);
  const quantityError = !isWholeNumber
    ? 'Enter a whole number'
    : mode === 'set'
      ? amount < 0 ? 'Stock cannot be negative' : ''
      : amount <= 0
        ? 'Enter a number greater than 0'
        : mode === 'subtract' && amount > item.currentStock
          ? `Only ${item.currentStock} in stock`
          : '';

  const resultingStock = !isWholeNumber
    ? item.currentStock
    : mode === 'add' ? item.currentStock + amount : mode === 'set' ? amount : item.currentStock - amount;

  const mutation = useMutation({
    mutationFn: () => {
      const payload = { menuItemId: item.menuItemId, quantity: amount, reason: reason.trim() || DEFAULT_REASON[mode] };
      return mode === 'add' ? restockItem(payload) : updateStock({ ...payload, operation: mode });
    },
    onSuccess: () => {
      toast.success(`${item.name} now has ${resultingStock} ${item.unit} in stock`);
      for (const key of ['inventory', 'low-stock', 'inventory-alerts', 'menu-items']) {
        queryClient.invalidateQueries({ queryKey: [key] });
      }
      onClose();
    },
    onError: (error: Error & { details?: string }) => {
      toast.error(error.details || error.message || 'Failed to update stock');
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!quantityError) mutation.mutate();
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50" role="presentation">
      <form
        onSubmit={handleSubmit}
        role="dialog"
        aria-modal="true"
        aria-labelledby="stock-dialog-title"
        className="bg-white dark:bg-gray-800 rounded-xl w-full max-w-md p-6 space-y-5"
      >
        <div className="flex items-start justify-between">
          <div>
            <h2 id="stock-dialog-title" className="text-lg font-semibold text-gray-900 dark:text-white">
              Update stock: {item.name}
            </h2>
            <p className="text-sm text-gray-500 dark:text-gray-400">
              Currently {item.currentStock} {item.unit} in stock
            </p>
          </div>
          <button type="button" onClick={onClose} aria-label="Close" className="p-1 text-gray-400 hover:text-gray-600">
            <X className="h-5 w-5" />
          </button>
        </div>

        <div className="grid grid-cols-3 gap-2" role="radiogroup" aria-label="Change type">
          {MODES.map((option) => (
            <button
              key={option.value}
              type="button"
              role="radio"
              aria-checked={mode === option.value}
              onClick={() => setMode(option.value)}
              className={`px-3 py-2 text-sm rounded-lg border transition-colors ${
                mode === option.value
                  ? 'border-primary-500 bg-primary-50 text-primary-700 dark:bg-primary-900/30 dark:text-primary-300'
                  : 'border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700'
              }`}
            >
              {option.label}
            </button>
          ))}
        </div>
        <p className="text-sm text-gray-500 dark:text-gray-400 -mt-2">
          {MODES.find((option) => option.value === mode)?.help}
        </p>

        <div>
          <label htmlFor="stock-quantity" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
            {mode === 'set' ? 'New stock level' : 'Quantity'} ({item.unit})
          </label>
          <input
            id="stock-quantity"
            type="number"
            inputMode="numeric"
            min={mode === 'set' ? 0 : 1}
            step={1}
            value={quantity}
            onChange={(e) => setQuantity(e.target.value)}
            autoFocus
            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
          />
          {quantity !== '' && quantityError && <p className="mt-1 text-sm text-red-600">{quantityError}</p>}
          {!quantityError && (
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              Stock after this change: <strong>{resultingStock}</strong> {item.unit}
            </p>
          )}
        </div>

        <div>
          <label htmlFor="stock-reason" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
            Reason <span className="text-gray-400">(optional)</span>
          </label>
          <input
            id="stock-reason"
            type="text"
            maxLength={200}
            value={reason}
            onChange={(e) => setReason(e.target.value)}
            placeholder={DEFAULT_REASON[mode]}
            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
          />
        </div>

        <div className="flex justify-end space-x-3">
          <button
            type="button"
            onClick={onClose}
            className="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={!!quantityError || mutation.isPending}
            className="px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors disabled:opacity-50"
          >
            {mutation.isPending ? 'Saving…' : 'Save'}
          </button>
        </div>
      </form>
    </div>
  );
};

export default StockAdjustDialog;
