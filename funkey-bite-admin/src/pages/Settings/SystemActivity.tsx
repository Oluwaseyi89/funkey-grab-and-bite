import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { ChevronDown, ChevronRight, RefreshCw } from 'lucide-react';
import toast from 'react-hot-toast';
import { getAuditLog, getBackgroundJobs, retryBackgroundJob } from '../../api/adminApi';
import SettingsTabs from '../../components/settings/SettingsTabs';
import type { AuditEntry } from '../../types';

const formatDate = (value: string) =>
  new Date(value).toLocaleString('en-NG', { dateStyle: 'medium', timeStyle: 'short' });

// What each kind of background job does, in plain words.
const JOB_KINDS: Record<string, string> = {
  'order.confirmation': 'Order confirmation email',
  'order.invoice': 'Pay-on-delivery invoice email',
  'order.payment_check': 'Check an online transfer was paid',
  'order.paid_email': 'Payment confirmed email',
  'order.paid_sms': 'Payment confirmed text',
  'email.send': 'Catering quote, invoice or receipt email',
};

const ACTOR_LABELS: Record<string, string> = {
  admin: 'Admin',
  customer: 'Customer',
  system: 'System',
  paystack: 'Paystack',
};

// Friendly names for the detailed audit actions; others show as recorded.
const ACTION_LABELS: Record<string, string> = {
  'order.created': 'Order placed',
  'order.payment_received': 'Online payment received',
  'order.payment_needs_review': 'Payment sent for review',
  'order.payment_expired': 'Unpaid order cancelled',
  'order.payment_review_accept': 'Payment accepted',
  'order.payment_review_refund': 'Payment marked for refund',
  'order.payment_review_refunded': 'Refund recorded',
};

const entityLink = (entry: AuditEntry) => {
  if (!entry.entityId) return null;
  const paths: Record<string, string> = { order: '/orders', orders: '/orders', catering: '/catering', menu: '/menu', customers: '/customers' };
  const base = paths[entry.entityType];
  return base ? `${base}/${entry.entityId}` : null;
};

const FailedJobs: React.FC = () => {
  const queryClient = useQueryClient();
  const { data, isLoading } = useQuery({
    queryKey: ['background-jobs', 'dead'],
    queryFn: () => getBackgroundJobs({ status: 'dead', limit: 50 }),
  });
  const retry = useMutation({
    mutationFn: retryBackgroundJob,
    onSuccess: () => {
      toast.success('Queued to try again');
      queryClient.invalidateQueries({ queryKey: ['background-jobs'] });
    },
    onError: (error: Error) => toast.error(error.message || 'Failed to retry'),
  });

  const jobs = data?.jobs ?? [];
  return (
    <section className="space-y-3">
      <div>
        <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Failed background tasks</h2>
        <p className="text-sm text-gray-600 dark:text-gray-400">
          Emails and payment checks are retried automatically for about an hour. These ran out of attempts.
        </p>
      </div>
      <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden">
        {isLoading ? (
          <div className="p-6 text-center text-gray-500">Loading…</div>
        ) : jobs.length === 0 ? (
          <div className="p-6 text-center text-gray-500 dark:text-gray-400">No failed tasks.</div>
        ) : (
          <ul className="divide-y divide-gray-200 dark:divide-gray-700">
            {jobs.map((job) => (
              <li key={job.id} className="p-4 flex flex-col md:flex-row md:items-center justify-between gap-3">
                <div className="min-w-0">
                  <div className="font-medium text-gray-900 dark:text-white">{JOB_KINDS[job.kind] ?? job.kind}</div>
                  <div className="text-xs text-gray-500 dark:text-gray-400">
                    {job.attempts} of {job.maxAttempts} attempts · last tried {formatDate(job.updatedAt)}
                  </div>
                  {job.lastError && (
                    <div className="mt-1 text-xs text-red-700 dark:text-red-300 break-words">{job.lastError}</div>
                  )}
                </div>
                <button
                  onClick={() => retry.mutate(job.id)}
                  disabled={retry.isPending}
                  className="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg border border-gray-300 dark:border-gray-600 text-sm hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50 shrink-0"
                >
                  <RefreshCw className="h-4 w-4" /> Retry
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>
    </section>
  );
};

const AuditLog: React.FC = () => {
  const [actorType, setActorType] = useState('');
  const [entityType, setEntityType] = useState('');
  const [page, setPage] = useState(1);
  const [expanded, setExpanded] = useState<number | null>(null);

  const { data, isLoading } = useQuery({
    queryKey: ['audit-log', actorType, entityType, page],
    queryFn: () => getAuditLog({ actorType: actorType || undefined, entityType: entityType || undefined, page, limit: 25 }),
  });
  const entries = data?.entries ?? [];
  const totalPages = data?.pagination.totalPages ?? 0;

  const select = (setter: (v: string) => void) => (e: React.ChangeEvent<HTMLSelectElement>) => {
    setter(e.target.value);
    setPage(1);
  };

  return (
    <section className="space-y-3">
      <div>
        <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Activity log</h2>
        <p className="text-sm text-gray-600 dark:text-gray-400">Every change, who made it and when.</p>
      </div>
      <div className="flex flex-wrap gap-3">
        <select aria-label="Who" value={actorType} onChange={select(setActorType)} className="px-3 py-2 bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg dark:text-white text-sm">
          <option value="">Anyone</option>
          {Object.entries(ACTOR_LABELS).map(([value, label]) => <option key={value} value={value}>{label}</option>)}
        </select>
        <select aria-label="What" value={entityType} onChange={select(setEntityType)} className="px-3 py-2 bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg dark:text-white text-sm">
          <option value="">Everything</option>
          <option value="order">Orders (payments)</option>
          <option value="orders">Orders (admin changes)</option>
          <option value="catering">Catering</option>
          <option value="menu">Menu</option>
          <option value="inventory">Inventory</option>
          <option value="settings">Settings</option>
        </select>
      </div>
      <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden">
        {isLoading ? (
          <div className="p-6 text-center text-gray-500">Loading…</div>
        ) : entries.length === 0 ? (
          <div className="p-6 text-center text-gray-500 dark:text-gray-400">No activity recorded yet.</div>
        ) : (
          <ul className="divide-y divide-gray-200 dark:divide-gray-700 text-sm">
            {entries.map((entry) => {
              const link = entityLink(entry);
              const open = expanded === entry.id;
              const hasDetails = entry.before || entry.after || entry.metadata;
              return (
                <li key={entry.id} className="p-4">
                  <div className="flex flex-col md:flex-row md:items-center justify-between gap-2">
                    <div className="min-w-0">
                      <div className="font-medium text-gray-900 dark:text-white break-words">
                        {ACTION_LABELS[entry.action] ?? entry.action}
                      </div>
                      <div className="text-xs text-gray-500 dark:text-gray-400">
                        {ACTOR_LABELS[entry.actorType] ?? entry.actorType}
                        {entry.actorId && ` #${entry.actorId}`}
                        {entry.ipAddress && ` · ${entry.ipAddress}`}
                        {' · '}
                        {link ? (
                          <Link to={link} className="text-primary-600 dark:text-primary-400 hover:underline">
                            {entry.entityType} #{entry.entityId}
                          </Link>
                        ) : (
                          <>{entry.entityType}{entry.entityId && ` #${entry.entityId}`}</>
                        )}
                      </div>
                    </div>
                    <div className="flex items-center gap-3 shrink-0">
                      <span className="text-xs text-gray-500 dark:text-gray-400">{formatDate(entry.occurredAt)}</span>
                      {hasDetails && (
                        <button
                          onClick={() => setExpanded(open ? null : entry.id)}
                          aria-expanded={open}
                          className="inline-flex items-center gap-1 text-xs text-primary-600 dark:text-primary-400"
                        >
                          {open ? <ChevronDown className="h-3 w-3" /> : <ChevronRight className="h-3 w-3" />} Details
                        </button>
                      )}
                    </div>
                  </div>
                  {open && (
                    <pre className="mt-2 p-3 bg-gray-50 dark:bg-gray-900 rounded-lg text-xs overflow-x-auto text-gray-700 dark:text-gray-300">
                      {JSON.stringify({ before: entry.before, after: entry.after, details: entry.metadata }, null, 2)}
                    </pre>
                  )}
                </li>
              );
            })}
          </ul>
        )}
        {totalPages > 1 && (
          <div className="flex items-center justify-between px-4 py-3 border-t border-gray-200 dark:border-gray-700 text-sm">
            <button disabled={page <= 1} onClick={() => setPage((p) => p - 1)} className="px-3 py-1 rounded-lg disabled:opacity-50 hover:bg-gray-100 dark:hover:bg-gray-700">Previous</button>
            <span className="text-gray-500 dark:text-gray-400">Page {page} of {totalPages}</span>
            <button disabled={page >= totalPages} onClick={() => setPage((p) => p + 1)} className="px-3 py-1 rounded-lg disabled:opacity-50 hover:bg-gray-100 dark:hover:bg-gray-700">Next</button>
          </div>
        )}
      </div>
    </section>
  );
};

const SystemActivity: React.FC = () => (
  <div className="space-y-6">
    <SettingsTabs />
    <div>
      <h1 className="text-2xl font-bold text-gray-900 dark:text-white">System activity</h1>
      <p className="text-gray-600 dark:text-gray-400 mt-1">Background tasks that failed, and a record of every change.</p>
    </div>
    <FailedJobs />
    <AuditLog />
  </div>
);

export default SystemActivity;
