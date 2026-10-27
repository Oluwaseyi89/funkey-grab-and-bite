import type { BusinessSettings, CateringContract, CateringInvoice, CateringReceipt } from '../../types';
import { PAYMENT_METHOD_LABELS, formatEventDate, naira } from './cateringLabels';

const escape = (value: string | number | undefined | null) =>
  String(value ?? '')
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');

const day = (iso: string) => new Date(iso).toLocaleDateString('en-NG', { day: 'numeric', month: 'short', year: 'numeric' });

// Opens a printable invoice or receipt in a new window and prints it.
export const printCateringDocument = (
  contract: CateringContract,
  document: { invoice: CateringInvoice } | { receipt: CateringReceipt },
  business: BusinessSettings | null,
) => {
  const { request, quote } = contract;
  if (!quote) return;

  const isReceipt = 'receipt' in document;
  const title = isReceipt ? `Receipt ${document.receipt.receiptNumber}` : `Invoice ${document.invoice.invoiceNumber} (v${document.invoice.version})`;
  const issued = isReceipt ? document.receipt.issuedAt : document.invoice.issuedAt;
  const paid = isReceipt ? document.receipt.totalPaid : document.invoice.amountPaid;
  const balance = isReceipt ? 0 : document.invoice.balance;
  // An invoice shows the payments made by the time it was issued.
  const payments = contract.payments.filter((p) => isReceipt || new Date(p.createdAt) <= new Date(issued));
  const address = business
    ? [business.address, business.city, business.state, business.country].filter(Boolean).join(', ')
    : '';

  const lines = quote.lineItems.length > 0 ? quote.lineItems : [{ description: 'Catering service', amount: quote.amount }];

  const html = `<!doctype html><html><head><meta charset="utf-8"><title>${escape(title)}</title>
<style>
  body { font-family: Arial, sans-serif; color: #1f2937; max-width: 720px; margin: 24px auto; padding: 0 16px; }
  table { width: 100%; border-collapse: collapse; }
  td, th { padding: 6px 8px; }
  .lines td, .lines th { border-bottom: 1px solid #e5e7eb; }
  .right { text-align: right; }
  .muted { color: #6b7280; }
</style></head><body>
  <h2 style="margin-bottom:4px">${escape(title)}</h2>
  <p class="muted" style="margin-top:0">Issued ${escape(day(issued))} · Catering request ${escape(request.reference)}</p>
  <table style="margin:16px 0"><tr>
    <td style="vertical-align:top;width:50%"><strong>From</strong><br>${escape(business?.businessName)}<br>${escape(address)}<br>${escape(business?.phoneNumber)}<br>${escape(business?.email)}</td>
    <td style="vertical-align:top;width:50%"><strong>For</strong><br>${escape(request.contactName)}${request.organization ? ` (${escape(request.organization)})` : ''}<br>${escape(request.contactPhone)}<br>${escape(request.contactEmail)}<br>
      Event: ${escape(formatEventDate(request.eventDate, request.eventTime))}, ${escape(request.guestCount)} guests<br>${escape(request.venueAddress)}</td>
  </tr></table>
  <table class="lines"><thead><tr style="background:#f3f4f6"><th style="text-align:left">Description</th><th class="right">Amount</th></tr></thead><tbody>
    ${lines.map((l) => `<tr><td>${escape(l.description)}</td><td class="right">${escape(naira(l.amount))}</td></tr>`).join('')}
  </tbody></table>
  <table style="margin-top:12px">
    <tr><td class="right"><strong>Total</strong></td><td class="right" style="width:160px"><strong>${escape(naira(quote.amount))}</strong></td></tr>
    <tr><td class="right">Deposit to confirm (${escape(quote.depositPercent)}%)</td><td class="right">${escape(naira(quote.depositAmount))}</td></tr>
    <tr><td class="right">Paid</td><td class="right">${escape(naira(paid))}</td></tr>
    <tr><td class="right"><strong>Balance</strong></td><td class="right"><strong>${escape(naira(balance))}</strong></td></tr>
  </table>
  ${payments.length > 0 ? `<h3>Payments</h3><table class="lines">${payments
    .map((p) => `<tr><td>${escape(day(p.paidAt))}</td><td>${escape(PAYMENT_METHOD_LABELS[p.method] ?? p.method)}</td><td>${escape(p.reference)}</td><td class="right">${escape(naira(p.amount))}</td></tr>`)
    .join('')}</table>` : ''}
  ${isReceipt ? '<p style="margin-top:24px">Paid in full. Thank you.</p>' : `<h3>Terms</h3><p class="muted" style="white-space:pre-line;font-size:13px">${escape(quote.terms)}</p>`}
</body></html>`;

  const win = window.open('', '_blank', 'width=800,height=900');
  if (!win) return;
  win.document.open();
  win.document.write(html);
  win.document.close();
  win.focus();
  win.print();
};
