import { useCallback, useRef } from 'react';

const newKey = () =>
  typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function'
    ? `admin-${crypto.randomUUID()}`
    : `admin-${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 12)}`;

// One Idempotency-Key per attempt at an action. The same key is sent if the
// admin tries again after a failure (e.g. a timeout where the server did save
// it), so the action is applied once; call reset() after it succeeds so the
// next action gets a new key.
export const useIdempotencyKey = () => {
  const key = useRef(newKey());
  const current = useCallback(() => key.current, []);
  const reset = useCallback(() => {
    key.current = newKey();
  }, []);
  return { current, reset };
};
