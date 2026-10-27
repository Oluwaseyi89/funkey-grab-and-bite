import React, { useState } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import {
  ArrowLeft,
  Calendar,
  Clock,
  Users,
  MapPin,
  Phone,
  Mail,
  Printer,
  ChefHat,
  Building2,
  Utensils,
  AlertCircle,
  FileText,
  Wallet,
  History,
} from 'lucide-react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import {
  changeCateringStatus,
  getCateringContract,
  recordCateringPayment,
  sendCateringQuote,
  updateCateringRequest,
} from '../../api/adminApi';
import { useBusinessSettings } from '../../hooks/useBusinessSettings';
import type { CateringContract, CateringRequest, CateringStatus } from '../../types';
import {
  AccountSummary,
  DocumentsList,
  PaymentForm,
  PaymentsLedger,
  QuoteComposer,
  QuoteHistory,
  StatusHistory,
} from './CateringContractPanels';
import {
  DIETARY_LABELS,
  MANUAL_STATUSES,
  SERVICE_STYLE_LABELS,
  STATUS_STYLES,
  eventTypeLabel,
  formatEventDate,
  naira,
} from './cateringLabels';

const Section: React.FC<{ title: string; icon: React.ReactNode; children: React.ReactNode }> = ({ title, icon, children }) => (
  <section className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
    <h2 className="flex items-center gap-2 text-lg font-semibold text-gray-900 dark:text-white mb-4">
      {icon}
      {title}
    </h2>
    {children}
  </section>
);

const Row: React.FC<{ label: string; children: React.ReactNode }> = ({ label, children }) => (
  <div className="flex justify-between gap-4 py-1.5 text-sm">
    <span className="text-gray-500 dark:text-gray-400">{label}</span>
    <span className="text-right text-gray-900 dark:text-white">{children}</span>
  </div>
);

const STATUS_ACTION_LABELS: Partial<Record<CateringStatus, string>> = {
  pending: 'Reopen',
  declined: 'Decline',
  completed: 'Mark completed',
  cancelled: 'Cancel booking',
};

interface NotesFormProps {
  request: CateringRequest;
  saving: boolean;
  onSave: (notes: string) => void;
}

// Keyed on the request's last update, so it starts from the saved notes.
const NotesForm: React.FC<NotesFormProps> = ({ request, saving, onSave }) => {
  const [notes, setNotes] = useState(request.adminNotes ?? '');

  return (
    <>
      <textarea
        value={notes}
        onChange={(e) => setNotes(e.target.value)}
        rows={4}
        maxLength={2000}
        placeholder="Only visible to admins"
        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white"
      />
      <div className="flex justify-end mt-2 print:hidden">
        <button
          disabled={notes === (request.adminNotes ?? '') || saving}
          onClick={() => onSave(notes)}
          className="px-3 py-2 text-sm bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 rounded-lg disabled:opacity-50"
        >
          Save notes
        </button>
      </div>
    </>
  );
};

const CateringDetails: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const requestId = Number(id);
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const { settings } = useBusinessSettings();

  const { data: contract, isLoading, isError } = useQuery({
    queryKey: ['catering-request', requestId],
    queryFn: () => getCateringContract(requestId),
    enabled: Number.isInteger(requestId) && requestId > 0,
  });
  const request = contract?.request;

  const onContract = (message: string) => (updated: CateringContract) => {
    queryClient.setQueryData(['catering-request', requestId], updated);
    queryClient.invalidateQueries({ queryKey: ['catering-requests'] });
    toast.success(message);
  };
  const onError = (error: Error) => toast.error(error.message || 'Something went wrong');

  const notesMutation = useMutation({
    mutationFn: (adminNotes: string) => updateCateringRequest(requestId, { adminNotes }),
    onSuccess: (updated: CateringRequest) => {
      queryClient.setQueryData<CateringContract>(['catering-request', requestId], (current) =>
        current ? { ...current, request: updated } : current,
      );
      toast.success('Notes saved');
    },
    onError,
  });
  const quoteMutation = useMutation({
    mutationFn: (input: Parameters<typeof sendCateringQuote>[1]) => sendCateringQuote(requestId, input),
    onSuccess: onContract('Quote sent to the customer'),
    onError,
  });
  const paymentMutation = useMutation({
    mutationFn: (input: Parameters<typeof recordCateringPayment>[1]) => recordCateringPayment(requestId, input),
    onSuccess: onContract('Payment recorded'),
    onError,
  });
  const statusMutation = useMutation({
    mutationFn: ({ status, note }: { status: CateringStatus; note?: string }) => changeCateringStatus(requestId, status, note),
    onSuccess: onContract('Status updated'),
    onError,
  });
  const saving = quoteMutation.isPending || paymentMutation.isPending || statusMutation.isPending;

  if (isLoading) {
    return <div className="p-8 text-center text-gray-500">Loading catering request…</div>;
  }

  if (isError || !contract || !request) {
    return (
      <div className="text-center py-12">
        <AlertCircle className="h-12 w-12 text-gray-400 mx-auto mb-3" />
        <p className="text-gray-700 dark:text-gray-300 mb-4">This catering request could not be found.</p>
        <button onClick={() => navigate('/catering')} className="px-4 py-2 bg-primary-500 text-white rounded-lg">
          Back to Catering
        </button>
      </div>
    );
  }

  const formKey = `${request.id}-${request.updatedAt ?? ''}`;
  const status = STATUS_STYLES[request.status];
  const venue = [request.venueAddress, request.venueArea, request.venueState].filter(Boolean).join(', ');
  const mapsUrl = venue ? `https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(venue)}` : '';

  const canQuote = request.status === 'pending' || request.status === 'quoted';
  const paymentDue = (request.status === 'accepted' || request.status === 'confirmed') && contract.account;

  const changeStatus = (next: CateringStatus) => {
    let note: string | undefined;
    if (next === 'declined' || next === 'cancelled') {
      const reason = window.prompt(`Reason for ${next === 'declined' ? 'declining' : 'cancelling'} (shown in the history)`);
      if (reason === null) return;
      note = reason.trim() || undefined;
    }
    statusMutation.mutate({ status: next, note });
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div className="flex items-center gap-3">
          <button onClick={() => navigate('/catering')} aria-label="Back" className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 print:hidden">
            <ArrowLeft className="h-5 w-5" />
          </button>
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
              <span className="font-mono">{request.reference}</span>
              <span className={`px-3 py-1 rounded-full text-xs font-medium ${status?.className ?? ''}`}>{status?.label ?? request.status}</span>
            </h1>
            <p className="text-gray-600 dark:text-gray-400">
              {request.requestType === 'package' ? `Package: ${request.packageName ?? '—'}` : 'Custom catering request'}
              {' · '}Received {new Date(request.createdAt).toLocaleString('en-NG', { dateStyle: 'medium', timeStyle: 'short' })}
            </p>
          </div>
        </div>
        <div className="flex flex-wrap gap-2 print:hidden">
          <a href={`tel:${request.contactPhone}`} className="flex items-center gap-2 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700">
            <Phone className="h-4 w-4" /> Call
          </a>
          {request.contactEmail && (
            <a href={`mailto:${request.contactEmail}?subject=${encodeURIComponent(`Your catering request ${request.reference}`)}`} className="flex items-center gap-2 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700">
              <Mail className="h-4 w-4" /> Email
            </a>
          )}
          <button onClick={() => window.print()} className="flex items-center gap-2 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700">
            <Printer className="h-4 w-4" /> Print
          </button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-6">
          <Section title="Event" icon={<Calendar className="h-5 w-5 text-primary-500" />}>
            <Row label="Occasion">{request.eventName || '—'}</Row>
            <Row label="Event type">{eventTypeLabel(request.eventType, request.eventTypeOther)}</Row>
            <Row label="Date & start time">{formatEventDate(request.eventDate, request.eventTime)}</Row>
            <Row label="Duration">{request.durationHours != null ? `${request.durationHours} hours` : '—'}</Row>
            <Row label="Guests"><span className="inline-flex items-center gap-1"><Users className="h-4 w-4 text-gray-400" />{request.guestCount}</span></Row>
          </Section>

          <Section title="Venue" icon={<MapPin className="h-5 w-5 text-primary-500" />}>
            <Row label="Venue">{request.venueName || '—'}</Row>
            <Row label="Address">{venue || '—'}</Row>
            {request.venueNotes && <Row label="Access notes">{request.venueNotes}</Row>}
            {mapsUrl && (
              <a href={mapsUrl} target="_blank" rel="noopener noreferrer" className="inline-block mt-2 text-sm text-primary-600 dark:text-primary-400 hover:underline print:hidden">
                Open in Google Maps
              </a>
            )}
          </Section>

          <Section title="Food & service" icon={<Utensils className="h-5 w-5 text-primary-500" />}>
            <Row label="Service style">{request.serviceStyle ? SERVICE_STYLE_LABELS[request.serviceStyle] : '—'}</Row>
            {request.packageId && (
              <Row label="Package">
                <Link to="/catering/packages" className="text-primary-600 dark:text-primary-400 hover:underline">{request.packageName}</Link>
              </Row>
            )}
            <div className="py-1.5 text-sm">
              <span className="text-gray-500 dark:text-gray-400">
                {request.requestType === 'package' ? 'Dishes chosen' : 'Dish preferences'}
              </span>
              {request.selectedItems.length > 0 ? (
                <ul className="mt-1 flex flex-wrap gap-2">
                  {request.selectedItems.map((item) => (
                    <li key={item.menuItemId} className="px-2.5 py-1 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200">
                      {item.name}
                    </li>
                  ))}
                </ul>
              ) : (
                <p className="mt-1 text-gray-900 dark:text-white">None selected</p>
              )}
            </div>
            {request.menuNotes && <Row label="Menu notes">{request.menuNotes}</Row>}
            <Row label="Dietary needs">
              {request.dietaryRequirements.length > 0
                ? request.dietaryRequirements.map((d) => DIETARY_LABELS[d] ?? d).join(', ')
                : 'None'}
            </Row>
            {request.dietaryNotes && <Row label="Dietary notes">{request.dietaryNotes}</Row>}
            {request.specialRequests && <Row label="Special requests">{request.specialRequests}</Row>}
          </Section>
        </div>

        <div className="space-y-6">
          <Section title="Customer" icon={<ChefHat className="h-5 w-5 text-primary-500" />}>
            <Row label="Name">{request.contactName}</Row>
            {request.organization && (
              <Row label="Organisation"><span className="inline-flex items-center gap-1"><Building2 className="h-4 w-4 text-gray-400" />{request.organization}</span></Row>
            )}
            <Row label="Phone"><a href={`tel:${request.contactPhone}`} className="text-primary-600 dark:text-primary-400">{request.contactPhone}</a></Row>
            {request.contactEmail && (
              <Row label="Email"><a href={`mailto:${request.contactEmail}`} className="text-primary-600 dark:text-primary-400 break-all">{request.contactEmail}</a></Row>
            )}
          </Section>

          <Section title="Price" icon={<span className="text-primary-500 font-bold">₦</span>}>
            <Row label="Package estimate">{naira(request.estimatedTotal)}</Row>
            <Row label="Customer budget">{naira(request.budget)}</Row>
            {contract.quote && (
              <Row label={contract.quote.status === 'accepted' ? 'Accepted quote' : `Open quote v${contract.quote.version}`}>
                {naira(contract.quote.amount)}
                {contract.quote.isExpired && <span className="ml-1 text-xs text-red-600">expired</span>}
              </Row>
            )}
            {canQuote && (
              <div className="mt-4">
                <QuoteComposer
                  key={`${formKey}-${contract.quote?.id ?? 0}-${settings ? 'terms' : ''}`}
                  contract={contract}
                  defaultTerms={settings?.cateringTerms ?? ''}
                  saving={saving}
                  onSend={quoteMutation.mutate}
                />
              </div>
            )}
          </Section>

          {contract.account && (
            <Section title="Account" icon={<Wallet className="h-5 w-5 text-primary-500" />}>
              <AccountSummary contract={contract} />
              {paymentDue && contract.account.balance > 0 && (
                <div className="mt-4 pt-4 border-t border-gray-100 dark:border-gray-700">
                  <PaymentForm key={`${contract.account.balance}`} balance={contract.account.balance} saving={saving} onRecord={paymentMutation.mutate} />
                </div>
              )}
              <div className="mt-4 pt-4 border-t border-gray-100 dark:border-gray-700">
                <h3 className="text-sm font-semibold text-gray-900 dark:text-white mb-2">Payments</h3>
                <PaymentsLedger contract={contract} />
              </div>
            </Section>
          )}

          <Section title="Invoices & receipt" icon={<FileText className="h-5 w-5 text-primary-500" />}>
            <DocumentsList contract={contract} business={settings} />
          </Section>

          <Section title="Status" icon={<Clock className="h-5 w-5 text-primary-500" />}>
            <div className="flex flex-wrap gap-2 print:hidden">
              {MANUAL_STATUSES[request.status].length === 0 && (
                <p className="text-sm text-gray-500">This booking is completed.</p>
              )}
              {MANUAL_STATUSES[request.status].map((next) => (
                <button
                  key={next}
                  disabled={saving}
                  onClick={() => changeStatus(next)}
                  className={`px-3 py-2 text-sm rounded-lg border transition-colors disabled:opacity-50 ${
                    next === 'declined' || next === 'cancelled'
                      ? 'border-red-300 text-red-700 hover:bg-red-50 dark:text-red-300 dark:hover:bg-red-900/20'
                      : 'border-gray-300 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700'
                  }`}
                >
                  {STATUS_ACTION_LABELS[next] ?? next}
                </button>
              ))}
            </div>
            <p className="mt-2 text-xs text-gray-500 print:hidden">
              Quoted, accepted, confirmed and paid are set by quotes, the customer's acceptance and payments.
            </p>
          </Section>

          <Section title="History" icon={<History className="h-5 w-5 text-primary-500" />}>
            <StatusHistory contract={contract} />
            <div className="mt-4 pt-4 border-t border-gray-100 dark:border-gray-700">
              <h3 className="text-sm font-semibold text-gray-900 dark:text-white mb-2">Quotes</h3>
              <QuoteHistory contract={contract} />
            </div>
          </Section>

          <Section title="Internal notes" icon={<AlertCircle className="h-5 w-5 text-primary-500" />}>
            <NotesForm key={formKey} request={request} saving={notesMutation.isPending} onSave={notesMutation.mutate} />
          </Section>
        </div>
      </div>
    </div>
  );
};

export default CateringDetails;
