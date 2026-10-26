import React from 'react';
import type { Invoice } from '../../types';

const naira = (amount: number) =>
  `₦${amount.toLocaleString('en-NG', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;

const collectedLabel: Record<string, string> = {
  cash: 'in cash',
  pos: 'by POS',
  transfer: 'by bank transfer',
};

// The invoice for a pay-on-delivery order, laid out for printing on white paper.
const InvoiceDocument: React.FC<{ invoice: Invoice }> = ({ invoice }) => {
  const isPaid = invoice.paymentStatus === 'paid';
  const handover = invoice.orderType === 'delivery' ? 'delivery' : 'pickup';

  return (
    <article className="bg-white text-gray-900 p-8 text-sm">
      <div className="flex items-start justify-between pb-6 border-b border-gray-300">
        <div>
          <h1 className="text-2xl font-bold">Invoice {invoice.invoiceNumber}</h1>
          <p className="text-gray-600">
            Order {invoice.orderNumber} · {new Date(invoice.issuedAt).toLocaleString('en-NG', { dateStyle: 'medium', timeStyle: 'short' })}
          </p>
        </div>
        <p className="font-semibold">{isPaid ? 'PAID' : `PAYMENT DUE ON ${handover.toUpperCase()}`}</p>
      </div>

      <div className="grid grid-cols-2 gap-8 py-6 border-b border-gray-300">
        <div>
          <p className="text-xs uppercase text-gray-500 mb-1">From</p>
          <p className="font-medium">{invoice.business.name}</p>
          {invoice.business.address && <p>{invoice.business.address}</p>}
          <p>{invoice.business.phone}</p>
          {invoice.business.email && <p>{invoice.business.email}</p>}
        </div>
        <div>
          <p className="text-xs uppercase text-gray-500 mb-1">Bill to</p>
          <p className="font-medium">{invoice.customer.name}</p>
          <p>{invoice.customer.phone}</p>
          {invoice.customer.email && <p>{invoice.customer.email}</p>}
          {invoice.customer.address && <p className="mt-1">Deliver to: {invoice.customer.address}</p>}
        </div>
      </div>

      <table className="w-full mt-4">
        <thead>
          <tr className="text-left border-b border-gray-300">
            <th className="py-2">Item</th>
            <th className="py-2 text-right">Qty</th>
            <th className="py-2 text-right">Price</th>
            <th className="py-2 text-right">Amount</th>
          </tr>
        </thead>
        <tbody>
          {invoice.items.map((line, index) => (
            <tr key={index} className="border-b border-gray-200">
              <td className="py-2">{line.name}</td>
              <td className="py-2 text-right">{line.quantity}</td>
              <td className="py-2 text-right">{naira(line.unitPrice)}</td>
              <td className="py-2 text-right">{naira(line.total)}</td>
            </tr>
          ))}
        </tbody>
      </table>

      <dl className="mt-4 ml-auto w-64 space-y-1">
        <div className="flex justify-between"><dt>Subtotal</dt><dd>{naira(invoice.subtotal)}</dd></div>
        {invoice.deliveryFee > 0 && <div className="flex justify-between"><dt>Delivery fee</dt><dd>{naira(invoice.deliveryFee)}</dd></div>}
        {invoice.tax > 0 && <div className="flex justify-between"><dt>Tax</dt><dd>{naira(invoice.tax)}</dd></div>}
        {invoice.discount > 0 && <div className="flex justify-between"><dt>Discount</dt><dd>-{naira(invoice.discount)}</dd></div>}
        <div className="flex justify-between pt-2 border-t border-gray-300 font-bold text-base">
          <dt>{isPaid ? 'Total paid' : 'Total due'}</dt>
          <dd>{naira(invoice.total)}</dd>
        </div>
      </dl>

      <p className="mt-8">
        {isPaid
          ? `Payment received ${collectedLabel[invoice.paymentCollectedVia ?? ''] ?? ''}${invoice.paidAt ? ` on ${new Date(invoice.paidAt).toLocaleString('en-NG')}` : ''}. Thank you!`
          : `Please pay ${naira(invoice.total)} on ${handover}. Cash, POS and bank transfer are accepted.`}
      </p>
    </article>
  );
};

export default InvoiceDocument;
