import React, { useEffect, useId, useRef } from 'react';
import { createPortal } from 'react-dom';
import { X } from 'lucide-react';

interface SheetProps {
  open: boolean;
  title: string;
  onClose: () => void;
  children: React.ReactNode;
  footer?: React.ReactNode;
  // On larger screens: a panel on the right (default) or a centred dialog.
  desktop?: 'right' | 'center';
}

// Open sheets, newest last: only the top one answers the keyboard.
const stack: symbol[] = [];

const focusableSelector =
  'a[href], button:not([disabled]), input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])';

// A dialog that slides up from the bottom on phones. It keeps focus inside,
// closes on Escape or a tap outside, locks the page behind it and returns
// focus to where it was.
const Sheet: React.FC<SheetProps> = ({ open, title, onClose, children, footer, desktop = 'right' }) => {
  const panel = useRef<HTMLDivElement>(null);
  const titleId = useId();
  const closeRef = useRef(onClose);
  closeRef.current = onClose;

  useEffect(() => {
    if (!open) return;
    const self = Symbol('sheet');
    stack.push(self);
    const returnTo = document.activeElement as HTMLElement | null;
    const previousOverflow = document.body.style.overflow;
    document.body.style.overflow = 'hidden';
    requestAnimationFrame(() => panel.current?.focus());

    const onKeyDown = (event: KeyboardEvent) => {
      if (stack[stack.length - 1] !== self || !panel.current) return;
      if (event.key === 'Escape') {
        event.preventDefault();
        closeRef.current();
        return;
      }
      if (event.key !== 'Tab') return;
      const items = Array.from(panel.current.querySelectorAll<HTMLElement>(focusableSelector)).filter((el) => el.offsetParent !== null);
      if (!items.length) return;
      const first = items[0];
      const last = items[items.length - 1];
      if (!panel.current.contains(document.activeElement)) {
        event.preventDefault();
        first.focus();
      } else if (event.shiftKey && document.activeElement === first) {
        event.preventDefault();
        last.focus();
      } else if (!event.shiftKey && document.activeElement === last) {
        event.preventDefault();
        first.focus();
      }
    };
    document.addEventListener('keydown', onKeyDown);

    return () => {
      document.removeEventListener('keydown', onKeyDown);
      const index = stack.indexOf(self);
      if (index >= 0) stack.splice(index, 1);
      if (!stack.length) document.body.style.overflow = previousOverflow;
      returnTo?.focus?.();
    };
  }, [open]);

  if (!open) return null;

  const position =
    desktop === 'right'
      ? 'md:inset-y-0 md:left-auto md:right-0 md:max-h-none md:w-[440px] md:rounded-none md:rounded-l-3xl md:animate-slide-left'
      : 'md:inset-auto md:left-1/2 md:top-1/2 md:w-[min(520px,92vw)] md:max-h-[85vh] md:-translate-x-1/2 md:-translate-y-1/2 md:rounded-3xl md:animate-fade-in';

  return createPortal(
    <>
      <div className="fixed inset-0 z-[80] bg-gray-950/50 animate-fade-in" aria-hidden="true" onClick={onClose} />
      <div
        ref={panel}
        role="dialog"
        aria-modal="true"
        aria-labelledby={titleId}
        tabIndex={-1}
        className={`fixed inset-x-0 bottom-0 z-[90] flex max-h-[92dvh] flex-col rounded-t-3xl bg-white shadow-lift outline-none animate-slide-up dark:bg-gray-900 ${position}`}
      >
        <div className="flex justify-center pt-2.5 md:hidden" aria-hidden="true">
          <span className="h-1.5 w-10 rounded-full bg-gray-200 dark:bg-gray-700" />
        </div>
        <header className="flex items-center justify-between gap-3 px-5 pb-2 pt-2 md:pt-5">
          <h2 id={titleId} className="text-lg font-bold text-gray-900 dark:text-white">{title}</h2>
          <button type="button" className="icon-btn -mr-2" aria-label={`Close ${title}`} onClick={onClose}>
            <X className="h-5 w-5" aria-hidden="true" />
          </button>
        </header>
        <div className="min-h-0 flex-1 overflow-y-auto overscroll-contain px-5 pb-5">{children}</div>
        {footer && <footer className="border-t border-gray-100 px-5 pb-safe pt-3 dark:border-gray-800"><div className="pb-3">{footer}</div></footer>}
      </div>
    </>,
    document.body,
  );
};

export default Sheet;
