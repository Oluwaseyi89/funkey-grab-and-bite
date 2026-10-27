import React, { useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import { ArrowLeft, Plus, Pencil, Users, Clock, X, Utensils } from 'lucide-react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import {
  createCateringPackage,
  getCateringPackages,
  getMenuItems,
  updateCateringPackage,
} from '../../api/adminApi';
import type { CateringPackage, CateringPackageInput, ServiceStyle } from '../../types';
import { SERVICE_STYLE_LABELS, naira } from './cateringLabels';

const emptyInput = (): CateringPackageInput => ({
  name: '',
  description: '',
  pricePerPerson: 0,
  minGuests: 20,
  maxGuests: undefined,
  serviceStyle: 'buffet',
  minNoticeDays: 3,
  inclusions: [],
  dishChoices: 0,
  menuItemIds: [],
  isActive: true,
  displayOrder: 0,
});

const toInput = (pkg: CateringPackage): CateringPackageInput => ({
  name: pkg.name,
  description: pkg.description,
  pricePerPerson: pkg.pricePerPerson,
  minGuests: pkg.minGuests,
  maxGuests: pkg.maxGuests,
  serviceStyle: pkg.serviceStyle,
  minNoticeDays: pkg.minNoticeDays,
  inclusions: pkg.inclusions,
  dishChoices: pkg.dishChoices,
  menuItemIds: pkg.menuItems.map((item) => item.menuItemId),
  isActive: pkg.isActive,
  displayOrder: pkg.displayOrder,
});

const inputClass = 'w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white';
const labelClass = 'block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1';

interface PackageFormProps {
  pkg: CateringPackage | null;
  onClose: () => void;
}

const PackageForm: React.FC<PackageFormProps> = ({ pkg, onClose }) => {
  const queryClient = useQueryClient();
  const [form, setForm] = useState<CateringPackageInput>(pkg ? toInput(pkg) : emptyInput());
  const [inclusionsText, setInclusionsText] = useState((pkg?.inclusions ?? []).join('\n'));
  const [dishSearch, setDishSearch] = useState('');

  const { data: menuItems = [] } = useQuery({
    queryKey: ['menu-items', 'catering-picker'],
    queryFn: () => getMenuItems({ page: 1, limit: 500 }),
  });

  const set = <K extends keyof CateringPackageInput>(key: K, value: CateringPackageInput[K]) =>
    setForm((prev) => ({ ...prev, [key]: value }));

  const toggleDish = (id: number) =>
    set('menuItemIds', form.menuItemIds.includes(id)
      ? form.menuItemIds.filter((x) => x !== id)
      : [...form.menuItemIds, id]);

  const error =
    form.name.trim().length < 2 ? 'Enter a package name'
      : form.pricePerPerson < 0 ? 'Price cannot be negative'
        : form.minGuests < 1 ? 'Minimum guests must be at least 1'
          : form.maxGuests != null && form.maxGuests < form.minGuests ? 'Maximum guests cannot be less than minimum'
            : form.dishChoices > form.menuItemIds.length ? `Link at least ${form.dishChoices} ${form.dishChoices === 1 ? 'dish' : 'dishes'} for customers to choose from`
              : '';

  const mutation = useMutation({
    mutationFn: () => {
      const payload: CateringPackageInput = {
        ...form,
        inclusions: inclusionsText.split('\n').map((line) => line.trim()).filter(Boolean),
      };
      return pkg ? updateCateringPackage(pkg.id, payload) : createCateringPackage(payload);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['catering-packages'] });
      toast.success(pkg ? 'Package updated' : 'Package created');
      onClose();
    },
    onError: (err: Error) => toast.error(err.message || 'Failed to save package'),
  });

  const filteredMenu = useMemo(() => {
    const query = dishSearch.trim().toLowerCase();
    return query ? menuItems.filter((item) => item.name.toLowerCase().includes(query)) : menuItems;
  }, [menuItems, dishSearch]);

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50" role="presentation">
      <form
        role="dialog"
        aria-modal="true"
        aria-labelledby="package-form-title"
        onSubmit={(e) => { e.preventDefault(); if (!error) mutation.mutate(); }}
        className="bg-white dark:bg-gray-800 rounded-xl w-full max-w-2xl max-h-[90vh] overflow-y-auto p-6 space-y-4"
      >
        <div className="flex items-center justify-between">
          <h2 id="package-form-title" className="text-lg font-semibold text-gray-900 dark:text-white">
            {pkg ? `Edit ${pkg.name}` : 'New catering package'}
          </h2>
          <button type="button" onClick={onClose} aria-label="Close" className="p-1 text-gray-400 hover:text-gray-600"><X className="h-5 w-5" /></button>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div className="sm:col-span-2">
            <label htmlFor="pkg-name" className={labelClass}>Package name *</label>
            <input id="pkg-name" className={inputClass} value={form.name} maxLength={120} onChange={(e) => set('name', e.target.value)} placeholder="e.g. Owambe Party Trays" />
          </div>
          <div className="sm:col-span-2">
            <label htmlFor="pkg-description" className={labelClass}>Description</label>
            <textarea id="pkg-description" className={inputClass} rows={2} maxLength={1000} value={form.description} onChange={(e) => set('description', e.target.value)} placeholder="Who it suits and what makes it special" />
          </div>
          <div>
            <label htmlFor="pkg-price" className={labelClass}>Price per guest (₦) *</label>
            <input id="pkg-price" type="number" min={0} step={100} className={inputClass} value={form.pricePerPerson} onChange={(e) => set('pricePerPerson', Number(e.target.value))} />
          </div>
          <div>
            <label htmlFor="pkg-style" className={labelClass}>Service style *</label>
            <select id="pkg-style" className={inputClass} value={form.serviceStyle} onChange={(e) => set('serviceStyle', e.target.value as ServiceStyle)}>
              {Object.entries(SERVICE_STYLE_LABELS).map(([value, label]) => <option key={value} value={value}>{label}</option>)}
            </select>
          </div>
          <div>
            <label htmlFor="pkg-min" className={labelClass}>Minimum guests *</label>
            <input id="pkg-min" type="number" min={1} max={1000} className={inputClass} value={form.minGuests} onChange={(e) => set('minGuests', Number(e.target.value))} />
          </div>
          <div>
            <label htmlFor="pkg-max" className={labelClass}>Maximum guests</label>
            <input id="pkg-max" type="number" min={1} max={1000} className={inputClass} value={form.maxGuests ?? ''} placeholder="No limit" onChange={(e) => set('maxGuests', e.target.value === '' ? undefined : Number(e.target.value))} />
          </div>
          <div>
            <label htmlFor="pkg-notice" className={labelClass}>Minimum notice (days) *</label>
            <input id="pkg-notice" type="number" min={0} max={180} className={inputClass} value={form.minNoticeDays} onChange={(e) => set('minNoticeDays', Number(e.target.value))} />
          </div>
          <div>
            <label htmlFor="pkg-order" className={labelClass}>Display order</label>
            <input id="pkg-order" type="number" min={0} className={inputClass} value={form.displayOrder} onChange={(e) => set('displayOrder', Number(e.target.value))} />
          </div>
          <div className="sm:col-span-2">
            <label htmlFor="pkg-inclusions" className={labelClass}>What's included (one per line)</label>
            <textarea id="pkg-inclusions" className={inputClass} rows={3} value={inclusionsText} onChange={(e) => setInclusionsText(e.target.value)} placeholder={'Bottled water\n2 servers\nChafing dishes'} />
          </div>
        </div>

        <fieldset className="space-y-2">
          <legend className={labelClass}>Dishes customers choose from</legend>
          <div className="flex items-center gap-3">
            <label htmlFor="pkg-choices" className="text-sm text-gray-600 dark:text-gray-400">Customers pick</label>
            <input id="pkg-choices" type="number" min={0} max={20} className="w-20 px-2 py-1 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white" value={form.dishChoices} onChange={(e) => set('dishChoices', Number(e.target.value))} />
            <span className="text-sm text-gray-600 dark:text-gray-400">of the {form.menuItemIds.length} selected dishes (0 = fixed menu)</span>
          </div>
          <input type="search" className={inputClass} placeholder="Search menu items…" value={dishSearch} onChange={(e) => setDishSearch(e.target.value)} />
          <div className="max-h-48 overflow-y-auto border border-gray-200 dark:border-gray-700 rounded-lg divide-y divide-gray-100 dark:divide-gray-700">
            {filteredMenu.length === 0 && <p className="p-3 text-sm text-gray-500">No menu items found.</p>}
            {filteredMenu.map((item) => (
              <label key={item.id} className="flex items-center gap-3 px-3 py-2 text-sm cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700/50">
                <input type="checkbox" checked={form.menuItemIds.includes(item.id)} onChange={() => toggleDish(item.id)} />
                <span className="flex-1 text-gray-900 dark:text-white">{item.name}</span>
                {!item.isAvailable && <span className="text-xs text-red-600">Unavailable</span>}
              </label>
            ))}
          </div>
        </fieldset>

        <label className="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
          <input type="checkbox" checked={form.isActive} onChange={(e) => set('isActive', e.target.checked)} />
          Show this package to customers
        </label>

        {error && <p className="text-sm text-red-600">{error}</p>}

        <div className="flex justify-end gap-3">
          <button type="button" onClick={onClose} className="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700">Cancel</button>
          <button type="submit" disabled={!!error || mutation.isPending} className="px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg disabled:opacity-50">
            {mutation.isPending ? 'Saving…' : 'Save package'}
          </button>
        </div>
      </form>
    </div>
  );
};

const CateringPackages: React.FC = () => {
  const [editing, setEditing] = useState<CateringPackage | null | 'new'>(null);

  const { data: packages = [], isLoading } = useQuery({
    queryKey: ['catering-packages'],
    queryFn: getCateringPackages,
  });

  return (
    <div className="space-y-6">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div className="flex items-center gap-3">
          <Link to="/catering" aria-label="Back to catering requests" className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700">
            <ArrowLeft className="h-5 w-5" />
          </Link>
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Catering Packages</h1>
            <p className="text-gray-600 dark:text-gray-400">What customers can choose on the catering page. Custom requests are always available.</p>
          </div>
        </div>
        <button onClick={() => setEditing('new')} className="flex items-center gap-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg">
          <Plus className="h-4 w-4" /> New Package
        </button>
      </div>

      {isLoading ? (
        <div className="p-8 text-center text-gray-500">Loading packages…</div>
      ) : packages.length === 0 ? (
        <div className="bg-white dark:bg-gray-800 rounded-xl p-8 border border-gray-200 dark:border-gray-700 text-center">
          <Utensils className="h-12 w-12 text-gray-400 mx-auto mb-3" />
          <p className="text-gray-700 dark:text-gray-300 mb-1">No packages yet.</p>
          <p className="text-sm text-gray-500">Until you add one, customers can still send custom catering requests.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
          {packages.map((pkg) => (
            <article key={pkg.id} className={`bg-white dark:bg-gray-800 rounded-xl p-5 border ${pkg.isActive ? 'border-gray-200 dark:border-gray-700' : 'border-dashed border-gray-300 dark:border-gray-600 opacity-70'}`}>
              <div className="flex items-start justify-between gap-2">
                <div>
                  <h2 className="font-semibold text-gray-900 dark:text-white">{pkg.name}</h2>
                  <p className="text-sm text-gray-500 dark:text-gray-400">{SERVICE_STYLE_LABELS[pkg.serviceStyle]}</p>
                </div>
                <span className={`px-2 py-1 rounded-full text-xs font-medium ${pkg.isActive ? 'bg-green-100 text-green-800 dark:bg-green-900/40 dark:text-green-200' : 'bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300'}`}>
                  {pkg.isActive ? 'Shown' : 'Hidden'}
                </span>
              </div>
              <p className="mt-2 text-xl font-bold text-primary-600 dark:text-primary-400">{naira(pkg.pricePerPerson)}<span className="text-sm font-normal text-gray-500"> / guest</span></p>
              {pkg.description && <p className="mt-2 text-sm text-gray-600 dark:text-gray-300 line-clamp-2">{pkg.description}</p>}
              <div className="mt-3 flex flex-wrap gap-3 text-sm text-gray-600 dark:text-gray-400">
                <span className="inline-flex items-center gap-1"><Users className="h-4 w-4" />{pkg.minGuests}{pkg.maxGuests ? `–${pkg.maxGuests}` : '+'} guests</span>
                <span className="inline-flex items-center gap-1"><Clock className="h-4 w-4" />{pkg.minNoticeDays} days' notice</span>
              </div>
              <p className="mt-2 text-sm text-gray-600 dark:text-gray-400">
                {pkg.dishChoices > 0
                  ? `Customer picks ${pkg.dishChoices} of ${pkg.menuItems.length} dishes`
                  : pkg.menuItems.length > 0 ? `Fixed menu: ${pkg.menuItems.map((i) => i.name).join(', ')}` : 'No dishes linked'}
              </p>
              <button onClick={() => setEditing(pkg)} className="mt-4 flex items-center gap-2 px-3 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700">
                <Pencil className="h-4 w-4" /> Edit
              </button>
            </article>
          ))}
        </div>
      )}

      {editing && <PackageForm key={editing === 'new' ? 'new' : editing.id} pkg={editing === 'new' ? null : editing} onClose={() => setEditing(null)} />}
    </div>
  );
};

export default CateringPackages;
