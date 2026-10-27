import React, { useEffect, useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import { Search, Filter, RefreshCw, Package, ChefHat, Users, MapPin, Calendar } from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { getCateringRequests } from '../../api/adminApi';
import type { CateringRequest, CateringStatus } from '../../types';
import { STATUS_STYLES, eventTypeLabel, formatEventDate, naira } from './cateringLabels';

const PAGE_SIZE = 20;

const CateringList: React.FC = () => {
  const [page, setPage] = useState(1);
  const [statusFilter, setStatusFilter] = useState<CateringStatus | ''>('');
  const [typeFilter, setTypeFilter] = useState<'' | 'package' | 'custom'>('');
  const [searchQuery, setSearchQuery] = useState('');

  const { data, isLoading, isError, refetch } = useQuery({
    queryKey: ['catering-requests'],
    queryFn: async () => (await getCateringRequests({ page: 1, limit: 1000 })).data,
  });

  const requests = useMemo(() => {
    let list: CateringRequest[] = data ?? [];
    if (statusFilter) list = list.filter((r) => r.status === statusFilter);
    if (typeFilter) list = list.filter((r) => r.requestType === typeFilter);
    const query = searchQuery.trim().toLowerCase();
    if (query) {
      list = list.filter((r) =>
        [r.reference, r.contactName, r.contactPhone, r.contactEmail, r.eventName, r.venueArea, r.packageName]
          .some((value) => value?.toLowerCase().includes(query)),
      );
    }
    return list;
  }, [data, statusFilter, typeFilter, searchQuery]);

  const totalPages = Math.max(1, Math.ceil(requests.length / PAGE_SIZE));
  const pageRequests = requests.slice((page - 1) * PAGE_SIZE, page * PAGE_SIZE);
  const pendingCount = (data ?? []).filter((r) => r.status === 'pending').length;

  useEffect(() => {
    if (isError) toast.error('Failed to load catering requests');
  }, [isError]);

  return (
    <div className="space-y-6">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Catering Requests</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            {pendingCount > 0 ? `${pendingCount} waiting for a quote or reply` : 'Follow up on catering requests'}
          </p>
        </div>
        <div className="flex items-center gap-3">
          <button
            onClick={() => refetch()}
            className="flex items-center space-x-2 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
          >
            <RefreshCw className="h-4 w-4" />
            <span>Refresh</span>
          </button>
          <Link
            to="/catering/packages"
            className="flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors"
          >
            <Package className="h-4 w-4" />
            <span>Manage Packages</span>
          </Link>
        </div>
      </div>

      <div className="bg-white dark:bg-gray-800 rounded-xl p-4 border border-gray-200 dark:border-gray-700 flex flex-col lg:flex-row gap-4">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 h-5 w-5" />
          <input
            type="search"
            placeholder="Search reference, customer, phone, area…"
            value={searchQuery}
            onChange={(e) => { setSearchQuery(e.target.value); setPage(1); }}
            className="w-full pl-10 pr-4 py-2 bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg dark:text-white"
          />
        </div>
        <div className="flex items-center gap-3">
          <Filter className="h-5 w-5 text-gray-400" />
          <select
            value={statusFilter}
            onChange={(e) => { setStatusFilter(e.target.value as CateringStatus | ''); setPage(1); }}
            className="px-3 py-2 bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg dark:text-white"
          >
            <option value="">All statuses</option>
            {Object.entries(STATUS_STYLES).map(([value, { label }]) => (
              <option key={value} value={value}>{label}</option>
            ))}
          </select>
          <select
            value={typeFilter}
            onChange={(e) => { setTypeFilter(e.target.value as '' | 'package' | 'custom'); setPage(1); }}
            className="px-3 py-2 bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg dark:text-white"
          >
            <option value="">Packages & custom</option>
            <option value="package">Package requests</option>
            <option value="custom">Custom requests</option>
          </select>
        </div>
      </div>

      <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden">
        {isLoading ? (
          <div className="p-8 text-center text-gray-500">Loading catering requests…</div>
        ) : pageRequests.length === 0 ? (
          <div className="p-8 text-center">
            <ChefHat className="h-12 w-12 text-gray-400 mx-auto mb-3" />
            <p className="text-gray-600 dark:text-gray-400">
              {data?.length ? 'No requests match your filters.' : 'No catering requests yet.'}
            </p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead className="bg-gray-50 dark:bg-gray-700/50 text-left text-xs uppercase text-gray-500 dark:text-gray-400">
                <tr>
                  <th className="px-4 py-3">Request</th>
                  <th className="px-4 py-3">Event</th>
                  <th className="px-4 py-3">Guests</th>
                  <th className="px-4 py-3">Venue</th>
                  <th className="px-4 py-3">Package</th>
                  <th className="px-4 py-3">Amount</th>
                  <th className="px-4 py-3">Status</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                {pageRequests.map((request) => (
                  <tr key={request.id} className="hover:bg-gray-50 dark:hover:bg-gray-700/30">
                    <td className="px-4 py-3">
                      <Link to={`/catering/${request.id}`} className="font-mono font-medium text-primary-600 dark:text-primary-400 hover:underline">
                        {request.reference}
                      </Link>
                      <div className="text-gray-900 dark:text-white">{request.contactName}</div>
                      <div className="text-gray-500 dark:text-gray-400">{request.contactPhone}</div>
                    </td>
                    <td className="px-4 py-3">
                      <div className="font-medium text-gray-900 dark:text-white">
                        {request.eventName || eventTypeLabel(request.eventType, request.eventTypeOther)}
                      </div>
                      <div className="flex items-center text-gray-500 dark:text-gray-400">
                        <Calendar className="h-3.5 w-3.5 mr-1" />
                        {formatEventDate(request.eventDate, request.eventTime)}
                      </div>
                    </td>
                    <td className="px-4 py-3">
                      <span className="inline-flex items-center"><Users className="h-3.5 w-3.5 mr-1 text-gray-400" />{request.guestCount}</span>
                    </td>
                    <td className="px-4 py-3">
                      <span className="inline-flex items-center text-gray-700 dark:text-gray-300">
                        <MapPin className="h-3.5 w-3.5 mr-1 text-gray-400" />
                        {[request.venueArea, request.venueState].filter(Boolean).join(', ') || '—'}
                      </span>
                    </td>
                    <td className="px-4 py-3">
                      {request.requestType === 'package' || request.packageName
                        ? request.packageName ?? 'Package'
                        : <span className="text-purple-600 dark:text-purple-400">Custom</span>}
                    </td>
                    <td className="px-4 py-3">
                      {request.quotedAmount != null ? (
                        <span className="font-medium">{naira(request.quotedAmount)} <span className="text-xs text-gray-500">quoted</span></span>
                      ) : request.estimatedTotal != null ? (
                        <span>{naira(request.estimatedTotal)} <span className="text-xs text-gray-500">est.</span></span>
                      ) : '—'}
                    </td>
                    <td className="px-4 py-3">
                      <span className={`px-2.5 py-1 rounded-full text-xs font-medium ${STATUS_STYLES[request.status]?.className ?? ''}`}>
                        {STATUS_STYLES[request.status]?.label ?? request.status}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {totalPages > 1 && (
          <div className="flex items-center justify-between px-4 py-3 border-t border-gray-200 dark:border-gray-700 text-sm">
            <span className="text-gray-500">Page {page} of {totalPages}</span>
            <div className="flex gap-2">
              <button disabled={page === 1} onClick={() => setPage((p) => p - 1)} className="px-3 py-1 border rounded-lg disabled:opacity-50">Previous</button>
              <button disabled={page === totalPages} onClick={() => setPage((p) => p + 1)} className="px-3 py-1 border rounded-lg disabled:opacity-50">Next</button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default CateringList;
