import React, { createContext, useCallback, useContext, useRef, useState } from 'react';
import Sheet from './Sheet';

export interface ConfirmOptions {
  title: string;
  message?: React.ReactNode;
  confirmLabel?: string;
  cancelLabel?: string;
  tone?: 'danger' | 'primary';
  // Asks for a short reason (e.g. why a booking is declined).
  reason?: { label: string; placeholder?: string; required?: boolean };
}

export interface ConfirmResult {
  confirmed: boolean;
  reason: string;
}

type Confirm = (options: ConfirmOptions) => Promise<ConfirmResult>;

const ConfirmContext = createContext<Confirm | null>(null);

// const { confirmed, reason } = await confirm({...}) anywhere in the app, in
// place of the browser's confirm() and prompt() boxes.
export const useConfirm = (): Confirm => {
  const confirm = useContext(ConfirmContext);
  if (!confirm) throw new Error('useConfirm must be used inside ConfirmProvider');
  return confirm;
};

export const ConfirmProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [options, setOptions] = useState<ConfirmOptions | null>(null);
  const [reason, setReason] = useState('');
  const resolver = useRef<((result: ConfirmResult) => void) | null>(null);

  const confirm = useCallback<Confirm>((next) => {
    setReason('');
    setOptions(next);
    return new Promise<ConfirmResult>((resolve) => {
      resolver.current = resolve;
    });
  }, []);

  const finish = (confirmed: boolean) => {
    resolver.current?.({ confirmed, reason: reason.trim() });
    resolver.current = null;
    setOptions(null);
  };

  const needsReason = Boolean(options?.reason?.required) && reason.trim().length === 0;

  return (
    <ConfirmContext.Provider value={confirm}>
      {children}
      <Sheet
        open={!!options}
        title={options?.title ?? ''}
        desktop="center"
        onClose={() => finish(false)}
        footer={
          <div className="grid grid-cols-2 gap-3">
            <button type="button" className="btn-secondary" onClick={() => finish(false)}>
              {options?.cancelLabel ?? 'Cancel'}
            </button>
            <button
              type="button"
              className={options?.tone === 'danger' ? 'btn-danger' : 'btn-primary'}
              disabled={needsReason}
              onClick={() => finish(true)}
            >
              {options?.confirmLabel ?? 'Confirm'}
            </button>
          </div>
        }
      >
        {options?.message && <div className="text-sm text-gray-600 dark:text-gray-300">{options.message}</div>}
        {options?.reason && (
          <div className="mt-4">
            <label htmlFor="confirm-reason" className="form-label">{options.reason.label}</label>
            <textarea
              id="confirm-reason"
              className="input-primary resize-none"
              rows={3}
              maxLength={500}
              placeholder={options.reason.placeholder}
              value={reason}
              onChange={(e) => setReason(e.target.value)}
            />
          </div>
        )}
      </Sheet>
    </ConfirmContext.Provider>
  );
};
