import React, { useState } from 'react';
import { Plus, Printer, Trash2 } from 'lucide-react';
import type {
  BusinessSettings,
  CateringContract,
  CateringPaymentInput,
  CateringQuoteInput,
  CateringQuoteLine,
} from '../../types';
import { PAYMENT_METHOD_LABELS, STATUS_STYLES, naira } from './cateringLabels';
import { printCateringDocument } from './printCateringDocument';

const inputClass =
  'w-full min-w-0 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white';
const labelClass = 'block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1';
const primaryButton = 'px-3 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg disabled:opacity-50';

const dateTime = (iso: string) => new Date(iso).toLocaleString('en-NG', { dateStyle: 'medium', timeStyle: 'short' });
const dateOnly = (iso: string) => new Date(iso).toLocaleDateString('en-NG', { dateStyle: 'medium' });

const Row: React.FC<{ label: string; children: React.ReactNode; strong?: boolean }> = ({ label, children, strong }) => (
  <div className={`flex justify-between gap-4 py-1.5 text-sm ${strong ? 'font-semibold' : ''}`}>
    <span className="text-gray-500 dark:text-gray-400">{label}</span>
    <span className="text-right text-gray-900 dark:text-white">{children}</span>
  </div>
);

interface QuoteComposerProps {
  contract: CateringContract;
  defaultTerms: string;
  saving: boolean;
  onSend: (input: CateringQuoteInput) => void;
}

// The admin prices a request. Line items are optional but must add up to
// the amount; terms start from the default in settings.
export const QuoteComposer: React.FC<QuoteComposerProps> = ({ contract, defaultTerms, saving, onSend }) => {
  const { request, quote } = contract;
  const suggested = quote?.amount ?? request.estimatedTotal;
  const [amount, setAmount] = useState(suggested != null ? String(suggested) : '');
  const [lines, setLines] = useState<{ description: string; amount: string }[]>(
    quote?.lineItems.map((l) => ({ description: l.description, amount: String(l.amount) })) ?? [],
  );
  const [deposit, setDeposit] = useState(String(quote?.depositPercent ?? 70));
  const [validDays, setValidDays] = useState('7');
  const [message, setMessage] = useState('');
  const [terms, setTerms] = useState<string | undefined>(quote?.terms);

  const amountValue = Number(amount);
  const lineTotal = lines.reduce((sum, l) => sum + (Number(l.amount) || 0), 0);
  const depositValue = Number(deposit);
  const daysValue = Number(validDays);
  const termsValue = terms ?? defaultTerms;

  const errors: string[] = [];
  if (!(amountValue > 0)) errors.push('Enter the quoted amount');
  if (lines.length > 0 && Math.abs(lineTotal - amountValue) > 0.009) {
    errors.push(`Lines add up to ${naira(lineTotal)}, not ${naira(amountValue || 0)}`);
  }
  if (lines.some((l) => !l.description.trim())) errors.push('Every line needs a description');
  if (!Number.isInteger(depositValue) || depositValue < 1 || depositValue > 100) errors.push('Deposit must be 1–100%');
  if (!Number.isInteger(daysValue) || daysValue < 1 || daysValue > 60) errors.push('Validity must be 1–60 days');
  if (!termsValue.trim()) errors.push('Enter the terms');

  const updateLine = (index: number, field: 'description' | 'amount', value: string) =>
    setLines((current) => current.map((l, i) => (i === index ? { ...l, [field]: value } : l)));

  const submit = () => {
    const lineItems: CateringQuoteLine[] = lines.map((l) => ({ description: l.description.trim(), amount: Number(l.amount) }));
    onSend({
      amount: amountValue,
      lineItems,
      depositPercent: depositValue,
      validDays: daysValue,
      message: message.trim() || undefined,
      // Omitted terms use the default from settings.
      terms: terms !== undefined ? terms : undefined,
    });
  };

  return (
    <div className="space-y-3 print:hidden">
      <div>
        <label htmlFor="quote-amount" className={labelClass}>Total (₦)</label>
        <input id="quote-amount" type="number" min={0} step={500} value={amount} onChange={(e) => setAmount(e.target.value)} className={inputClass} />
      </div>

      <div>
        <div className="flex items-center justify-between mb-1">
          <span className={labelClass}>Breakdown (optional)</span>
          <button
            type="button"
            onClick={() => setLines((current) => [...current, { description: '', amount: '' }])}
            className="flex items-center gap-1 text-sm text-primary-600 dark:text-primary-400"
          >
            <Plus className="h-4 w-4" /> Add line
          </button>
        </div>
        {lines.map((line, index) => (
          <div key={index} className="flex gap-2 mb-2">
            <input
              aria-label={`Line ${index + 1} description`}
              value={line.description}
              maxLength={200}
              placeholder="e.g. Food for 80 guests"
              onChange={(e) => updateLine(index, 'description', e.target.value)}
              className={`${inputClass} flex-[2]`}
            />
            <input
              aria-label={`Line ${index + 1} amount`}
              type="number"
              min={0}
              value={line.amount}
              onChange={(e) => updateLine(index, 'amount', e.target.value)}
              className={`${inputClass} flex-1`}
            />
            <button
              type="button"
              aria-label="Remove line"
              onClick={() => setLines((current) => current.filter((_, i) => i !== index))}
              className="p-2 text-gray-500 hover:text-red-600"
            >
              <Trash2 className="h-4 w-4" />
            </button>
          </div>
        ))}
        {lines.length > 0 && <p className="text-xs text-gray-500">Lines total {naira(lineTotal)}</p>}
      </div>

      <div className="grid grid-cols-2 gap-3">
        <div>
          <label htmlFor="quote-deposit" className={labelClass}>Deposit %</label>
          <input id="quote-deposit" type="number" min={1} max={100} value={deposit} onChange={(e) => setDeposit(e.target.value)} className={inputClass} />
          {amountValue > 0 && depositValue > 0 && (
            <p className="mt-1 text-xs text-gray-500">{naira(Math.round(amountValue * depositValue) / 100)} to confirm</p>
          )}
        </div>
        <div>
          <label htmlFor="quote-valid" className={labelClass}>Valid for (days)</label>
          <input id="quote-valid" type="number" min={1} max={60} value={validDays} onChange={(e) => setValidDays(e.target.value)} className={inputClass} />
        </div>
      </div>

      <div>
        <label htmlFor="quote-message" className={labelClass}>Message to customer (optional)</label>
        <textarea id="quote-message" rows={2} maxLength={1000} value={message} onChange={(e) => setMessage(e.target.value)} className={inputClass} />
      </div>

      <div>
        <label htmlFor="quote-terms" className={labelClass}>Terms</label>
        <textarea id="quote-terms" rows={5} maxLength={5000} value={termsValue} onChange={(e) => setTerms(e.target.value)} className={`${inputClass} text-xs`} />
        {terms !== undefined && terms !== defaultTerms && (
          <button type="button" onClick={() => setTerms(undefined)} className="mt-1 text-xs text-primary-600 dark:text-primary-400">
            Use the default terms
          </button>
        )}
      </div>

      {errors.length > 0 && amount !== '' && <p className="text-sm text-red-600">{errors[0]}</p>}
      <button disabled={errors.length > 0 || saving} onClick={submit} className={`${primaryButton} w-full`}>
        {quote ? `Send revised quote (v${quote.version + 1})` : 'Send quote'}
      </button>
      <p className="text-xs text-gray-500">
        The customer gets the quote by email with a link to accept it, and can also accept it on the website with their reference and phone number.
      </p>
    </div>
  );
};

export const AccountSummary: React.FC<{ contract: CateringContract }> = ({ contract }) => {
  const { account, quote, request } = contract;
  if (!account || !quote) return null;
  return (
    <div>
      <Row label="Invoice">{request.invoiceNumber}</Row>
      <Row label="Total">{naira(account.total)}</Row>
      <Row label={`Deposit (${quote.depositPercent}%)`}>{naira(account.depositDue)}</Row>
      <Row label="Paid">{naira(account.amountPaid)}</Row>
      <Row label="Balance" strong>{naira(account.balance)}</Row>
      {account.depositRemaining > 0 && (
        <p className="mt-2 text-sm text-amber-700 dark:text-amber-300">{naira(account.depositRemaining)} more confirms the booking.</p>
      )}
      {request.contractedAt && <p className="mt-2 text-xs text-gray-500">Booking confirmed {dateTime(request.contractedAt)}</p>}
      {request.paidInFullAt && <p className="text-xs text-gray-500">Paid in full {dateTime(request.paidInFullAt)}</p>}
      {contract.pendingCharge && (
        <p className="mt-2 text-xs text-gray-500">
          Customer has an online transfer of {naira(contract.pendingCharge.amount)} pending ({contract.pendingCharge.reference}).
        </p>
      )}
    </div>
  );
};

interface PaymentFormProps {
  balance: number;
  saving: boolean;
  onRecord: (input: CateringPaymentInput) => void;
}

// Records money the business received directly. Keyed on the balance, so it
// resets after each payment.
export const PaymentForm: React.FC<PaymentFormProps> = ({ balance, saving, onRecord }) => {
  const [amount, setAmount] = useState(String(balance));
  const [method, setMethod] = useState<'cash' | 'bank_transfer'>('cash');
  const [reference, setReference] = useState('');
  const [paidOn, setPaidOn] = useState('');
  const [note, setNote] = useState('');

  const value = Number(amount);
  let error = '';
  if (!(value > 0)) error = 'Enter the amount received';
  else if (value > balance + 0.009) error = `That is more than the ${naira(balance)} owed`;
  else if (method === 'bank_transfer' && !reference.trim()) error = 'Enter the transfer reference';

  return (
    <div className="space-y-3 print:hidden">
      <div className="grid grid-cols-2 gap-3">
        <div>
          <label htmlFor="payment-amount" className={labelClass}>Amount (₦)</label>
          <input id="payment-amount" type="number" min={0} value={amount} onChange={(e) => setAmount(e.target.value)} className={inputClass} />
        </div>
        <div>
          <label htmlFor="payment-method" className={labelClass}>Method</label>
          <select id="payment-method" value={method} onChange={(e) => setMethod(e.target.value as 'cash' | 'bank_transfer')} className={inputClass}>
            <option value="cash">Cash</option>
            <option value="bank_transfer">Bank transfer</option>
          </select>
        </div>
      </div>
      <div className="grid grid-cols-2 gap-3">
        <div>
          <label htmlFor="payment-reference" className={labelClass}>Reference{method === 'cash' ? ' (optional)' : ''}</label>
          <input id="payment-reference" maxLength={100} value={reference} onChange={(e) => setReference(e.target.value)} className={inputClass} />
        </div>
        <div>
          <label htmlFor="payment-date" className={labelClass}>Date received</label>
          <input id="payment-date" type="date" value={paidOn} max={new Date().toISOString().slice(0, 10)} onChange={(e) => setPaidOn(e.target.value)} className={inputClass} />
        </div>
      </div>
      <div>
        <label htmlFor="payment-note" className={labelClass}>Note (optional)</label>
        <input id="payment-note" maxLength={500} value={note} onChange={(e) => setNote(e.target.value)} className={inputClass} />
      </div>
      {error && amount !== '' && <p className="text-sm text-red-600">{error}</p>}
      <button
        disabled={!!error || saving}
        onClick={() =>
          onRecord({
            amount: value,
            method,
            reference: reference.trim() || undefined,
            note: note.trim() || undefined,
            paidOn: paidOn || undefined,
          })
        }
        className={`${primaryButton} w-full`}
      >
        Record payment
      </button>
      <p className="text-xs text-gray-500">The customer is emailed an updated invoice, or the receipt once fully paid.</p>
    </div>
  );
};

export const PaymentsLedger: React.FC<{ contract: CateringContract }> = ({ contract }) =>
  contract.payments.length === 0 ? (
    <p className="text-sm text-gray-500">No payments yet.</p>
  ) : (
    <ul className="divide-y divide-gray-100 dark:divide-gray-700">
      {contract.payments.map((p) => (
        <li key={p.id} className="py-2 text-sm flex justify-between gap-4">
          <div>
            <p className="text-gray-900 dark:text-white">{PAYMENT_METHOD_LABELS[p.method] ?? p.method}</p>
            <p className="text-xs text-gray-500">
              {dateOnly(p.paidAt)}
              {p.reference ? ` · ${p.reference}` : ''}
              {p.note ? ` · ${p.note}` : ''}
            </p>
          </div>
          <span className="font-medium text-gray-900 dark:text-white">{naira(p.amount)}</span>
        </li>
      ))}
    </ul>
  );

export const DocumentsList: React.FC<{ contract: CateringContract; business: BusinessSettings | null }> = ({ contract, business }) => {
  const invoices = contract.invoices ?? [];
  if (invoices.length === 0 && !contract.receipt) {
    return <p className="text-sm text-gray-500">Issued when the customer accepts a quote.</p>;
  }
  const button = 'flex items-center gap-1 text-sm text-primary-600 dark:text-primary-400 hover:underline print:hidden';
  return (
    <ul className="divide-y divide-gray-100 dark:divide-gray-700">
      {contract.receipt && (
        <li className="py-2 flex justify-between items-center gap-4 text-sm">
          <div>
            <p className="font-medium text-gray-900 dark:text-white">Receipt {contract.receipt.receiptNumber}</p>
            <p className="text-xs text-gray-500">{dateTime(contract.receipt.issuedAt)} · {naira(contract.receipt.totalPaid)} paid</p>
          </div>
          <button className={button} onClick={() => printCateringDocument(contract, { receipt: contract.receipt! }, business)}>
            <Printer className="h-4 w-4" /> Print
          </button>
        </li>
      )}
      {invoices.map((invoice) => (
        <li key={invoice.id} className="py-2 flex justify-between items-center gap-4 text-sm">
          <div>
            <p className="font-medium text-gray-900 dark:text-white">
              Invoice {invoice.invoiceNumber} v{invoice.version}
            </p>
            <p className="text-xs text-gray-500">
              {dateTime(invoice.issuedAt)} · {invoice.reason === 'accepted' ? 'on acceptance' : 'after payment'} · balance {naira(invoice.balance)}
            </p>
          </div>
          <button className={button} onClick={() => printCateringDocument(contract, { invoice }, business)}>
            <Printer className="h-4 w-4" /> Print
          </button>
        </li>
      ))}
    </ul>
  );
};

const QUOTE_STATUS_LABELS = { sent: 'Open', accepted: 'Accepted', superseded: 'Replaced', withdrawn: 'Withdrawn' } as const;

export const QuoteHistory: React.FC<{ contract: CateringContract }> = ({ contract }) => {
  const quotes = contract.quotes ?? [];
  if (quotes.length === 0) return <p className="text-sm text-gray-500">No quote sent yet.</p>;
  return (
    <ul className="divide-y divide-gray-100 dark:divide-gray-700">
      {quotes.map((q) => (
        <li key={q.id} className="py-2 text-sm">
          <div className="flex justify-between gap-4">
            <span className="font-medium text-gray-900 dark:text-white">
              v{q.version} · {naira(q.amount)} · {q.depositPercent}% deposit
            </span>
            <span className="text-xs text-gray-500">{q.isExpired ? 'Expired' : QUOTE_STATUS_LABELS[q.status]}</span>
          </div>
          <p className="text-xs text-gray-500">
            Sent {dateTime(q.createdAt)} · valid until {q.validUntil}
            {q.acceptedAt ? ` · accepted ${dateTime(q.acceptedAt)} via ${q.acceptedVia === 'email_link' ? 'email link' : 'status page'}` : ''}
          </p>
          {q.lineItems.length > 0 && (
            <ul className="mt-1 text-xs text-gray-600 dark:text-gray-300">
              {q.lineItems.map((l, i) => (
                <li key={i}>{l.description}: {naira(l.amount)}</li>
              ))}
            </ul>
          )}
        </li>
      ))}
    </ul>
  );
};

const ACTOR_LABELS = { admin: 'Admin', customer: 'Customer', system: 'System' } as const;

export const StatusHistory: React.FC<{ contract: CateringContract }> = ({ contract }) => {
  const history = contract.history ?? [];
  if (history.length === 0) return <p className="text-sm text-gray-500">Received {dateTime(contract.request.createdAt)}.</p>;
  return (
    <ol className="space-y-2">
      {history.map((h) => (
        <li key={h.id} className="text-sm">
          <p className="text-gray-900 dark:text-white">
            {STATUS_STYLES[h.toStatus]?.label ?? h.toStatus}
            <span className="text-xs text-gray-500"> · {ACTOR_LABELS[h.actorType]} · {dateTime(h.createdAt)}</span>
          </p>
          {h.note && <p className="text-xs text-gray-500">{h.note}</p>}
        </li>
      ))}
    </ol>
  );
};
