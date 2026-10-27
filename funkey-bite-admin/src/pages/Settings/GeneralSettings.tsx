import React, { useState, useEffect, useMemo } from 'react';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import {
  Save,
  Store,
  Phone,
  Mail,
  MapPin,
  Truck,
  Package,
  Clock,
  Banknote,
  Percent,
  CheckCircle,
  XCircle,
  Calendar,
  Building,
  Info,
  AlertCircle,
  RefreshCw,
  FileText
} from 'lucide-react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { updateSettings } from '../../api/adminApi';
import { useSettingsStore } from '../../stores/settingsStore';
import { useBusinessSettings, BUSINESS_SETTINGS_QUERY_KEY } from '../../hooks/useBusinessSettings';
import type { BusinessSettings, OpeningHours } from '../../types';

const settingsSchema = z.object({
  businessName: z.string().min(2, 'Business name must be at least 2 characters'),
  phoneNumber: z
    .string()
    .trim()
    .regex(/^\+?[0-9][0-9 ()-]{5,18}$/, 'Enter a valid phone number, e.g. +234 801 234 5678'),
  email: z.string().trim().email('Valid email is required'),
  address: z.string().trim().min(3, 'Street address is required').max(500, 'Street address is too long'),
  city: z.string().trim().min(1, 'City is required').max(100, 'City is too long'),
  state: z.string().trim().min(1, 'State is required').max(100, 'State is too long'),
  postalCode: z.string().trim().max(20, 'Postal code is too long'),
  country: z.string().trim().min(1, 'Country is required').max(100, 'Country is too long'),
  deliveryFee: z.number().min(0, 'Delivery fee cannot be negative'),
  minOrderAmount: z.number().min(0, 'Minimum order amount cannot be negative'),
  taxRate: z.number().min(0).max(100, 'Tax rate must be between 0-100'),
  isDeliveryOpen: z.boolean(),
  isPickupOpen: z.boolean(),
  payOnDeliveryEnabled: z.boolean(),
  cateringTerms: z
    .string()
    .trim()
    .min(1, 'Catering terms are required')
    .max(5000, 'Catering terms must be 5,000 characters or fewer'),
});

type SettingsFormData = z.infer<typeof settingsSchema>;

const DAYS_OF_WEEK = [
  { id: 'Monday', label: 'Monday' },
  { id: 'Tuesday', label: 'Tuesday' },
  { id: 'Wednesday', label: 'Wednesday' },
  { id: 'Thursday', label: 'Thursday' },
  { id: 'Friday', label: 'Friday' },
  { id: 'Saturday', label: 'Saturday' },
  { id: 'Sunday', label: 'Sunday' },
];

// Every day must be sent to the API exactly once, so fill in any day the
// stored settings are missing (as closed) and keep the Monday-first order.
const completeWeek = (hours: OpeningHours[] = []): OpeningHours[] =>
  DAYS_OF_WEEK.map((day) => {
    const existing = hours.find((h) => h.day.toLowerCase() === day.id.toLowerCase());
    return existing
      ? { ...existing, day: day.id }
      : { day: day.id, openTime: '08:00', closeTime: '22:00', isOpen: false };
  });

const CLOCK_PATTERN = /^([01]\d|2[0-3]):[0-5]\d$/;

// Mirrors the API rules: open days need valid, different times. A close time
// at or before the open time means the business closes after midnight.
const validateOpeningHours = (hours: OpeningHours[]): Record<string, string> => {
  const errors: Record<string, string> = {};
  for (const h of hours) {
    if (!h.isOpen) continue;
    if (!CLOCK_PATTERN.test(h.openTime) || !CLOCK_PATTERN.test(h.closeTime)) {
      errors[h.day] = 'Choose both an opening and a closing time';
    } else if (h.openTime === h.closeTime) {
      errors[h.day] = 'Opening and closing times cannot be the same';
    }
  }
  return errors;
};

const closesNextDay = (h: OpeningHours) =>
  h.isOpen && CLOCK_PATTERN.test(h.openTime) && CLOCK_PATTERN.test(h.closeTime) && h.closeTime <= h.openTime;

const TIME_OPTIONS = Array.from({ length: 48 }, (_, i) => {
  const hour = Math.floor(i / 2);
  const minute = i % 2 === 0 ? '00' : '30';
  const period = hour < 12 ? 'AM' : 'PM';
  const displayHour = hour === 0 ? 12 : hour > 12 ? hour - 12 : hour;
  return {
    value: `${hour.toString().padStart(2, '0')}:${minute}`,
    label: `${displayHour}:${minute} ${period}`,
  };
});

const toFormValues = (settings: BusinessSettings): SettingsFormData => ({
  businessName: settings.businessName,
  phoneNumber: settings.phoneNumber,
  email: settings.email,
  address: settings.address,
  city: settings.city ?? '',
  state: settings.state ?? '',
  postalCode: settings.postalCode ?? '',
  country: settings.country ?? '',
  deliveryFee: settings.deliveryFee,
  minOrderAmount: settings.minOrderAmount,
  taxRate: settings.taxRate,
  isDeliveryOpen: settings.isDeliveryOpen,
  isPickupOpen: settings.isPickupOpen,
  payOnDeliveryEnabled: settings.payOnDeliveryEnabled,
  cateringTerms: settings.cateringTerms ?? '',
});

const GeneralSettings: React.FC = () => {
  const [isSaving, setIsSaving] = useState(false);
  const [openingHours, setOpeningHours] = useState<OpeningHours[]>([]);
  const setSettings = useSettingsStore((state) => state.setSettings);
  const queryClient = useQueryClient();

  const { settings, isLoading, isError, refetch } = useBusinessSettings();

  const {
    control,
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
    watch,
  } = useForm<SettingsFormData>({
    resolver: zodResolver(settingsSchema),
    defaultValues: {
      businessName: '',
      phoneNumber: '',
      email: '',
      address: '',
      city: '',
      state: '',
      postalCode: '',
      country: '',
      deliveryFee: 0,
      minOrderAmount: 0,
      taxRate: 0,
      isDeliveryOpen: true,
      isPickupOpen: true,
      payOnDeliveryEnabled: true,
      cateringTerms: '',
    },
  });

  useEffect(() => {
    if (settings) {
      reset(toFormValues(settings));
      setOpeningHours(completeWeek(settings.openingHours));
    }
  }, [settings, reset]);

  const hoursErrors = useMemo(() => validateOpeningHours(openingHours), [openingHours]);

  const isDeliveryOpen = watch('isDeliveryOpen');
  const isPickupOpen = watch('isPickupOpen');
  const payOnDeliveryEnabled = watch('payOnDeliveryEnabled');
  const deliveryFee = watch('deliveryFee');
  const minOrderAmount = watch('minOrderAmount');
  const taxRate = watch('taxRate');

  const updateDayHours = (day: string, field: keyof OpeningHours, value: any) => {
    setOpeningHours(prev => prev.map(hour => 
      hour.day === day ? { ...hour, [field]: value } : hour
    ));
  };

  const toggleDayOpen = (day: string) => {
    setOpeningHours(prev => prev.map(hour =>
      hour.day === day ? { ...hour, isOpen: !hour.isOpen } : hour
    ));
  };

  const resetHoursToSaved = () => {
    setOpeningHours(completeWeek(settings?.openingHours));
  };

  const resetAllToSaved = () => {
    if (!settings) return;
    reset(toFormValues(settings));
    resetHoursToSaved();
  };

  const mutation = useMutation({
    mutationFn: async (data: SettingsFormData) => {
      const settingsData = {
        ...data,
        openingHours: openingHours,
      };
      
      return await updateSettings(settingsData);
    },
    onSuccess: (data) => {
      // Shares the cache with the sidebar and layout, so they update too.
      queryClient.setQueryData(BUSINESS_SETTINGS_QUERY_KEY, data);
      setSettings(data);
      toast.success('Settings updated successfully!');
    },
    onError: (error: any) => {
      // The API puts the specific validation failure in details.
      toast.error(error.details || error.message || 'Failed to update settings');
    },
  });

  const onSubmit = (data: SettingsFormData) => {
    if (Object.keys(hoursErrors).length > 0) {
      toast.error('Fix the highlighted opening hours before saving');
      return;
    }
    mutation.mutate(data);
  };

  const formatTime = (time: string) => {
    const [hour, minute] = time.split(':');
    const hourNum = parseInt(hour);
    const period = hourNum >= 12 ? 'PM' : 'AM';
    const displayHour = hourNum === 0 ? 12 : hourNum > 12 ? hourNum - 12 : hourNum;
    return `${displayHour}:${minute} ${period}`;
  };

  if (isLoading) {
    return (
      <div className="text-center py-12">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-500 mx-auto"></div>
        <p className="mt-4 text-gray-500 dark:text-gray-400">Loading settings...</p>
      </div>
    );
  }

  if (isError || !settings) {
    return (
      <div className="text-center py-12">
        <AlertCircle className="h-12 w-12 text-red-500 mx-auto" />
        <p className="mt-4 text-gray-700 dark:text-gray-300">Business settings could not be loaded.</p>
        <button
          type="button"
          onClick={() => refetch()}
          className="mt-4 inline-flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors"
        >
          <RefreshCw className="h-4 w-4" />
          <span>Try again</span>
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">General Settings</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Configure your business details and operations
          </p>
        </div>
        
        <div className="flex items-center space-x-3">
          <button
            onClick={() => refetch()}
            className="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors flex items-center space-x-2"
          >
            <RefreshCw className="h-4 w-4" />
            <span>Refresh</span>
          </button>
          <button
            onClick={handleSubmit(onSubmit)}
            disabled={isSubmitting || isSaving}
            className="flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors disabled:opacity-50"
          >
            <Save className="h-4 w-4" />
            <span>{(isSubmitting || isSaving) ? 'Saving...' : 'Save Settings'}</span>
          </button>
        </div>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center space-x-3 mb-6">
            <Building className="h-5 w-5 text-primary-500" />
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
              Business Information
            </h2>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Business Name *
              </label>
              <div className="relative">
                <Store className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                <input
                  type="text"
                  {...register('businessName')}
                  className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                  placeholder="Your business name"
                />
              </div>
              {errors.businessName && (
                <p className="mt-1 text-sm text-red-600">{errors.businessName.message}</p>
              )}
            </div>

            
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Phone Number *
              </label>
              <div className="relative">
                <Phone className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                <input
                  type="tel"
                  {...register('phoneNumber')}
                  className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                  placeholder="+234 801 234 5678"
                />
              </div>
              {errors.phoneNumber && (
                <p className="mt-1 text-sm text-red-600">{errors.phoneNumber.message}</p>
              )}
            </div>

            
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Email Address *
              </label>
              <div className="relative">
                <Mail className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                <input
                  type="email"
                  {...register('email')}
                  className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                  placeholder="contact@example.com"
                />
              </div>
              {errors.email && (
                <p className="mt-1 text-sm text-red-600">{errors.email.message}</p>
              )}
            </div>

            
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Street Address *
              </label>
              <div className="relative">
                <MapPin className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                <input
                  type="text"
                  {...register('address')}
                  className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                  placeholder="Street and number"
                />
              </div>
              {errors.address && (
                <p className="mt-1 text-sm text-red-600">{errors.address.message}</p>
              )}
            </div>

            {([
              { name: 'city', label: 'City *', placeholder: 'City' },
              { name: 'state', label: 'State *', placeholder: 'State' },
              { name: 'postalCode', label: 'Postal Code', placeholder: 'Optional' },
              { name: 'country', label: 'Country *', placeholder: 'Country' },
            ] as const).map((field) => (
              <div key={field.name}>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  {field.label}
                </label>
                <input
                  type="text"
                  {...register(field.name)}
                  className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                  placeholder={field.placeholder}
                />
                {errors[field.name] && (
                  <p className="mt-1 text-sm text-red-600">{errors[field.name]?.message}</p>
                )}
              </div>
            ))}
          </div>
        </div>

        
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between mb-6">
            <div className="flex items-center space-x-3">
              <Clock className="h-5 w-5 text-primary-500" />
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                Opening Hours
              </h2>
            </div>
            <button
              type="button"
              onClick={resetHoursToSaved}
              className="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300"
            >
              Reset to Saved
            </button>
          </div>
          
          <div className="space-y-4">
            <p className="text-sm text-gray-500 dark:text-gray-400">
              Times are in the business's local time ({settings.timezone}). Pick a closing time earlier than the
              opening time for hours that run past midnight.
            </p>
            {DAYS_OF_WEEK.map((day) => {
              const dayHours = openingHours.find(h => h.day === day.id) || {
                day: day.id,
                openTime: '08:00',
                closeTime: '22:00',
                isOpen: false,
              };
              const dayError = hoursErrors[day.id];
              
              return (
                <div key={day.id}>
                  <div className={`flex flex-col sm:flex-row sm:items-center justify-between gap-4 p-4 border rounded-lg ${dayError ? 'border-red-400 dark:border-red-500' : 'border-gray-200 dark:border-gray-700'}`}>
                    <div className="flex items-center space-x-4">
                      <label className="flex items-center space-x-2 cursor-pointer">
                        <div className="relative">
                          <input
                            type="checkbox"
                            checked={dayHours.isOpen}
                            onChange={() => toggleDayOpen(day.id)}
                            className="sr-only"
                          />
                          <div className={`w-10 h-6 rounded-full transition-colors ${dayHours.isOpen ? 'bg-green-500' : 'bg-gray-300 dark:bg-gray-600'}`}>
                            <div className={`absolute top-1 w-4 h-4 rounded-full bg-white transition-transform ${dayHours.isOpen ? 'left-5' : 'left-1'}`}></div>
                          </div>
                        </div>
                        <span className={`font-medium ${dayHours.isOpen ? 'text-gray-900 dark:text-white' : 'text-gray-400'}`}>
                          {day.label}
                        </span>
                      </label>
                    </div>
                  
                    <div className="flex items-center space-x-3">
                    
                      <div className="flex items-center space-x-2">
                        <span className="text-sm text-gray-500 dark:text-gray-400">Open:</span>
                        <select
                          value={dayHours.openTime}
                          onChange={(e) => updateDayHours(day.id, 'openTime', e.target.value)}
                          disabled={!dayHours.isOpen}
                          className="px-3 py-1 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white disabled:opacity-50"
                        >
                          {TIME_OPTIONS.map(time => (
                            <option key={time.value} value={time.value}>
                              {time.label}
                            </option>
                          ))}
                        </select>
                      </div>
                    
                    
                      <div className="flex items-center space-x-2">
                        <span className="text-sm text-gray-500 dark:text-gray-400">Close:</span>
                        <select
                          value={dayHours.closeTime}
                          onChange={(e) => updateDayHours(day.id, 'closeTime', e.target.value)}
                          disabled={!dayHours.isOpen}
                          className="px-3 py-1 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white disabled:opacity-50"
                        >
                          {TIME_OPTIONS.map(time => (
                            <option key={time.value} value={time.value}>
                              {time.label}
                            </option>
                          ))}
                        </select>
                      </div>
                    
                    
                      <span className={`px-2 py-1 text-xs rounded-full ${
                        dayHours.isOpen 
                          ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                          : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                      }`}>
                        {dayHours.isOpen ? (closesNextDay(dayHours) ? 'Closes next day' : 'Open') : 'Closed'}
                      </span>
                    </div>
                  </div>
                  {dayError && (
                    <p className="mt-1 text-sm text-red-600">{dayError}</p>
                  )}
                </div>
              );
            })}
          </div>
        </div>

        
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center space-x-3 mb-6">
            <Package className="h-5 w-5 text-primary-500" />
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
              Order Settings
            </h2>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Delivery Fee (₦)
              </label>
              <div className="relative">
                <Banknote className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                <input
                  type="number"
                  step="0.01"
                  min="0"
                  {...register('deliveryFee', { valueAsNumber: true })}
                  className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                  placeholder="0.00"
                />
              </div>
              {errors.deliveryFee && (
                <p className="mt-1 text-sm text-red-600">{errors.deliveryFee.message}</p>
              )}
            </div>

            
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Minimum Order Amount (₦)
              </label>
              <div className="relative">
                <Banknote className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                <input
                  type="number"
                  step="0.01"
                  min="0"
                  {...register('minOrderAmount', { valueAsNumber: true })}
                  className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                  placeholder="0.00"
                />
              </div>
              {errors.minOrderAmount && (
                <p className="mt-1 text-sm text-red-600">{errors.minOrderAmount.message}</p>
              )}
            </div>

            
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Tax Rate (%)
              </label>
              <div className="relative">
                <Percent className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                <input
                  type="number"
                  step="0.1"
                  min="0"
                  max="100"
                  {...register('taxRate', { valueAsNumber: true })}
                  className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                  placeholder="0"
                />
              </div>
              {errors.taxRate && (
                <p className="mt-1 text-sm text-red-600">{errors.taxRate.message}</p>
              )}
            </div>
          </div>
        </div>

        
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center space-x-3 mb-6">
            <Truck className="h-5 w-5 text-primary-500" />
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
              Service Availability
            </h2>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            
            <div className="p-4 border border-gray-200 dark:border-gray-700 rounded-lg">
              <div className="flex items-center justify-between mb-3">
                <div className="flex items-center space-x-3">
                  <Truck className="h-5 w-5 text-blue-500" />
                  <div>
                    <h3 className="font-medium text-gray-900 dark:text-white">Delivery Service</h3>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                      Accept delivery orders
                    </p>
                  </div>
                </div>
                <Controller
                  name="isDeliveryOpen"
                  control={control}
                  render={({ field }) => (
                    <label className="flex items-center space-x-2 cursor-pointer">
                      <div className="relative">
                        <input
                          type="checkbox"
                          checked={field.value}
                          onChange={field.onChange}
                          className="sr-only"
                        />
                        <div className={`w-10 h-6 rounded-full transition-colors ${field.value ? 'bg-blue-500' : 'bg-gray-300 dark:bg-gray-600'}`}>
                          <div className={`absolute top-1 w-4 h-4 rounded-full bg-white transition-transform ${field.value ? 'left-5' : 'left-1'}`}></div>
                        </div>
                      </div>
                      <span className={`text-sm ${field.value ? 'text-blue-600 dark:text-blue-400' : 'text-gray-400'}`}>
                        {field.value ? 'Available' : 'Unavailable'}
                      </span>
                    </label>
                  )}
                />
              </div>
              {isDeliveryOpen && (
                <div className="text-sm text-gray-600 dark:text-gray-400">
                  Delivery fee: ₦{deliveryFee.toFixed(2)}
                </div>
              )}
            </div>

            
            <div className="p-4 border border-gray-200 dark:border-gray-700 rounded-lg">
              <div className="flex items-center justify-between mb-3">
                <div className="flex items-center space-x-3">
                  <Package className="h-5 w-5 text-green-500" />
                  <div>
                    <h3 className="font-medium text-gray-900 dark:text-white">Pickup Service</h3>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                      Accept pickup orders
                    </p>
                  </div>
                </div>
                <Controller
                  name="isPickupOpen"
                  control={control}
                  render={({ field }) => (
                    <label className="flex items-center space-x-2 cursor-pointer">
                      <div className="relative">
                        <input
                          type="checkbox"
                          checked={field.value}
                          onChange={field.onChange}
                          className="sr-only"
                        />
                        <div className={`w-10 h-6 rounded-full transition-colors ${field.value ? 'bg-green-500' : 'bg-gray-300 dark:bg-gray-600'}`}>
                          <div className={`absolute top-1 w-4 h-4 rounded-full bg-white transition-transform ${field.value ? 'left-5' : 'left-1'}`}></div>
                        </div>
                      </div>
                      <span className={`text-sm ${field.value ? 'text-green-600 dark:text-green-400' : 'text-gray-400'}`}>
                        {field.value ? 'Available' : 'Unavailable'}
                      </span>
                    </label>
                  )}
                />
              </div>
              {isPickupOpen && (
                <div className="text-sm text-gray-600 dark:text-gray-400">
                  Minimum order: ₦{minOrderAmount.toFixed(2)}
                </div>
              )}
            </div>

            
            <div className="p-4 border border-gray-200 dark:border-gray-700 rounded-lg">
              <div className="flex items-center justify-between mb-3">
                <div className="flex items-center space-x-3">
                  <Banknote className="h-5 w-5 text-amber-500" />
                  <div>
                    <h3 className="font-medium text-gray-900 dark:text-white">Pay on Delivery</h3>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                      Let customers pay on delivery or at pickup
                    </p>
                  </div>
                </div>
                <Controller
                  name="payOnDeliveryEnabled"
                  control={control}
                  render={({ field }) => (
                    <label className="flex items-center space-x-2 cursor-pointer">
                      <div className="relative">
                        <input
                          type="checkbox"
                          checked={field.value}
                          onChange={field.onChange}
                          className="sr-only"
                        />
                        <div className={`w-10 h-6 rounded-full transition-colors ${field.value ? 'bg-amber-500' : 'bg-gray-300 dark:bg-gray-600'}`}>
                          <div className={`absolute top-1 w-4 h-4 rounded-full bg-white transition-transform ${field.value ? 'left-5' : 'left-1'}`}></div>
                        </div>
                      </div>
                      <span className={`text-sm ${field.value ? 'text-amber-600 dark:text-amber-400' : 'text-gray-400'}`}>
                        {field.value ? 'Allowed' : 'Off'}
                      </span>
                    </label>
                  )}
                />
              </div>
              <div className="text-sm text-gray-600 dark:text-gray-400">
                {payOnDeliveryEnabled
                  ? 'Customers can choose to pay on delivery or at pickup, or pay online. Each order gets an invoice.'
                  : 'Customers must pay online to place an order.'}
              </div>
            </div>
          </div>
        </div>

        
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center space-x-3 mb-2">
            <FileText className="h-5 w-5 text-primary-500" />
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Catering Terms</h2>
          </div>
          <p className="text-sm text-gray-500 dark:text-gray-400 mb-4">
            The default terms on every catering quote. You can change them for a single quote when you send it; quotes already sent keep their terms.
          </p>
          <label htmlFor="cateringTerms" className="sr-only">Catering terms</label>
          <textarea
            id="cateringTerms"
            rows={10}
            {...register('cateringTerms')}
            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white text-sm"
          />
          {errors.cateringTerms && (
            <p className="mt-1 text-sm text-red-600">{errors.cateringTerms.message}</p>
          )}
        </div>

        <div className="bg-gray-50 dark:bg-gray-700/50 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center space-x-3 mb-6">
            <Info className="h-5 w-5 text-primary-500" />
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
              Order Calculation Preview
            </h2>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-4">
              <div className="text-sm text-gray-600 dark:text-gray-400">
                Based on your current settings:
              </div>
              
              <div className="space-y-2">
                <div className="flex justify-between text-sm">
                  <span>Subtotal (Sample):</span>
                  <span>₦50.00</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span>Delivery Fee:</span>
                  <span>₦{deliveryFee.toFixed(2)}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span>Tax ({taxRate}%):</span>
                  <span>₦{(50 * (taxRate / 100)).toFixed(2)}</span>
                </div>
                <div className="border-t border-gray-300 dark:border-gray-600 pt-2">
                  <div className="flex justify-between font-medium">
                    <span>Total:</span>
                    <span>₦{(50 + deliveryFee + (50 * (taxRate / 100))).toFixed(2)}</span>
                  </div>
                </div>
              </div>
            </div>
            
            <div className="space-y-4">
              <div className="text-sm text-gray-600 dark:text-gray-400">
                Business Status:
              </div>
              
              <div className="space-y-3">
                <div className="flex items-center space-x-2">
                  {isDeliveryOpen ? (
                    <CheckCircle className="h-4 w-4 text-green-500" />
                  ) : (
                    <XCircle className="h-4 w-4 text-red-500" />
                  )}
                  <span className="text-sm">Delivery: {isDeliveryOpen ? 'Available' : 'Unavailable'}</span>
                </div>
                
                <div className="flex items-center space-x-2">
                  {isPickupOpen ? (
                    <CheckCircle className="h-4 w-4 text-green-500" />
                  ) : (
                    <XCircle className="h-4 w-4 text-red-500" />
                  )}
                  <span className="text-sm">Pickup: {isPickupOpen ? 'Available' : 'Unavailable'}</span>
                </div>

                <div className="flex items-center space-x-2">
                  {payOnDeliveryEnabled ? (
                    <CheckCircle className="h-4 w-4 text-green-500" />
                  ) : (
                    <XCircle className="h-4 w-4 text-red-500" />
                  )}
                  <span className="text-sm">Pay on delivery: {payOnDeliveryEnabled ? 'Allowed' : 'Off (online payment only)'}</span>
                </div>
                
                <div className="text-sm text-gray-600 dark:text-gray-400">
                  Minimum order for pickup: ₦{minOrderAmount.toFixed(2)}
                </div>
              </div>
            </div>
          </div>
        </div>

        
        <div className="flex items-center justify-end space-x-3">
          <button
            type="button"
            onClick={resetAllToSaved}
            className="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
          >
            Reset Changes
          </button>
          <button
            type="submit"
            disabled={isSubmitting || isSaving}
            className="flex items-center space-x-2 px-6 py-3 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors disabled:opacity-50"
          >
            <Save className="h-4 w-4" />
            <span className="font-medium">Save All Settings</span>
          </button>
        </div>
      </form>
    </div>
  );
};

export default GeneralSettings;