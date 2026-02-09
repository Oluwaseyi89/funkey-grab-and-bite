// src/pages/Settings/GeneralSettings.tsx
import React, { useState, useEffect } from 'react';
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
  DollarSign,
  Percent,
  CheckCircle,
  XCircle,
  Calendar,
  Building,
  Info,
  AlertCircle,
  RefreshCw
} from 'lucide-react';
import { useQuery, useMutation } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { getSettings, updateSettings } from '../../api/adminApi';
import { useSettingsStore } from '../../stores/settingsStore';
import type { BusinessSettings, OpeningHours } from '../../types';

// Validation schema
const settingsSchema = z.object({
  businessName: z.string().min(2, 'Business name must be at least 2 characters'),
  phoneNumber: z.string().min(5, 'Phone number is required'),
  email: z.string().email('Valid email is required'),
  address: z.string().min(5, 'Address is required'),
  deliveryFee: z.number().min(0, 'Delivery fee cannot be negative'),
  minOrderAmount: z.number().min(0, 'Minimum order amount cannot be negative'),
  taxRate: z.number().min(0).max(100, 'Tax rate must be between 0-100'),
  isDeliveryOpen: z.boolean(),
  isPickupOpen: z.boolean(),
});

type SettingsFormData = z.infer<typeof settingsSchema>;

// Days of week for opening hours
const DAYS_OF_WEEK = [
  { id: 'Monday', label: 'Monday' },
  { id: 'Tuesday', label: 'Tuesday' },
  { id: 'Wednesday', label: 'Wednesday' },
  { id: 'Thursday', label: 'Thursday' },
  { id: 'Friday', label: 'Friday' },
  { id: 'Saturday', label: 'Saturday' },
  { id: 'Sunday', label: 'Sunday' },
];

// Time options for dropdown
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

const GeneralSettings: React.FC = () => {
  const [isSaving, setIsSaving] = useState(false);
  const [openingHours, setOpeningHours] = useState<OpeningHours[]>([]);
  const { setSettings, updateSettings: updateStoreSettings } = useSettingsStore();

  // Fetch settings
  const { data: settings, isLoading, refetch } = useQuery({
    queryKey: ['settings'],
    queryFn: async () => {
      try {
        const data = await getSettings();
        setSettings(data);
        setOpeningHours(data.openingHours || []);
        return data;
      } catch (error) {
        toast.error('Failed to load settings');
        return null;
      }
    },
  });

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
      deliveryFee: 0,
      minOrderAmount: 0,
      taxRate: 0,
      isDeliveryOpen: true,
      isPickupOpen: true,
    },
  });

  // Reset form when settings load
  useEffect(() => {
    if (settings) {
      reset({
        businessName: settings.businessName,
        phoneNumber: settings.phoneNumber,
        email: settings.email,
        address: settings.address,
        deliveryFee: settings.deliveryFee,
        minOrderAmount: settings.minOrderAmount,
        taxRate: settings.taxRate,
        isDeliveryOpen: settings.isDeliveryOpen,
        isPickupOpen: settings.isPickupOpen,
      });
      setOpeningHours(settings.openingHours || []);
    }
  }, [settings, reset]);

  // Watch form values
  const isDeliveryOpen = watch('isDeliveryOpen');
  const isPickupOpen = watch('isPickupOpen');
  const deliveryFee = watch('deliveryFee');
  const minOrderAmount = watch('minOrderAmount');
  const taxRate = watch('taxRate');

  // Update opening hours
  const updateDayHours = (day: string, field: keyof OpeningHours, value: any) => {
    setOpeningHours(prev => prev.map(hour => 
      hour.day === day ? { ...hour, [field]: value } : hour
    ));
  };

  // Toggle day open/closed
  const toggleDayOpen = (day: string) => {
    setOpeningHours(prev => prev.map(hour =>
      hour.day === day ? { ...hour, isOpen: !hour.isOpen } : hour
    ));
  };

  // Reset to default hours
  const resetToDefaultHours = () => {
    const defaultHours: OpeningHours[] = DAYS_OF_WEEK.map((day, index) => ({
      day: day.id,
      openTime: index < 5 ? '08:00' : index === 5 ? '09:00' : '10:00',
      closeTime: index < 5 ? '22:00' : '23:00',
      isOpen: true,
    }));
    setOpeningHours(defaultHours);
  };

  // Mutation for saving settings
  const mutation = useMutation({
    mutationFn: async (data: SettingsFormData) => {
      // Combine form data with opening hours
      const settingsData = {
        ...data,
        openingHours: openingHours,
      };
      
      return await updateSettings(settingsData);
    },
    onSuccess: (data) => {
      setSettings(data);
      updateStoreSettings(data);
      toast.success('Settings updated successfully!');
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to update settings');
    },
  });

  const onSubmit = (data: SettingsFormData) => {
    mutation.mutate(data);
  };

  // Format time for display
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

  return (
    <div className="space-y-6">
      {/* Header */}
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
        {/* Business Information */}
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center space-x-3 mb-6">
            <Building className="h-5 w-5 text-primary-500" />
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
              Business Information
            </h2>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* Business Name */}
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
                  placeholder="Funkey Grab & Bite"
                />
              </div>
              {errors.businessName && (
                <p className="mt-1 text-sm text-red-600">{errors.businessName.message}</p>
              )}
            </div>

            {/* Phone Number */}
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
                  placeholder="+1 (234) 567-890"
                />
              </div>
              {errors.phoneNumber && (
                <p className="mt-1 text-sm text-red-600">{errors.phoneNumber.message}</p>
              )}
            </div>

            {/* Email */}
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
                  placeholder="contact@funkeygrabandbite.com"
                />
              </div>
              {errors.email && (
                <p className="mt-1 text-sm text-red-600">{errors.email.message}</p>
              )}
            </div>

            {/* Address */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Business Address *
              </label>
              <div className="relative">
                <MapPin className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                <input
                  type="text"
                  {...register('address')}
                  className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                  placeholder="123 Fast Food Street, Food City"
                />
              </div>
              {errors.address && (
                <p className="mt-1 text-sm text-red-600">{errors.address.message}</p>
              )}
            </div>
          </div>
        </div>

        {/* Opening Hours */}
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
              onClick={resetToDefaultHours}
              className="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300"
            >
              Reset to Default
            </button>
          </div>
          
          <div className="space-y-4">
            {DAYS_OF_WEEK.map((day) => {
              const dayHours = openingHours.find(h => h.day === day.id) || {
                day: day.id,
                openTime: '08:00',
                closeTime: '22:00',
                isOpen: true,
              };
              
              return (
                <div key={day.id} className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 p-4 border border-gray-200 dark:border-gray-700 rounded-lg">
                  <div className="flex items-center space-x-4">
                    <Controller
                      name=""
                      render={() => (
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
                      )}
                    />
                  </div>
                  
                  <div className="flex items-center space-x-3">
                    {/* Open Time */}
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
                    
                    {/* Close Time */}
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
                    
                    {/* Status Badge */}
                    <span className={`px-2 py-1 text-xs rounded-full ${
                      dayHours.isOpen 
                        ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                        : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                    }`}>
                      {dayHours.isOpen ? 'Open' : 'Closed'}
                    </span>
                  </div>
                </div>
              );
            })}
          </div>
        </div>

        {/* Order Settings */}
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center space-x-3 mb-6">
            <Package className="h-5 w-5 text-primary-500" />
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
              Order Settings
            </h2>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* Delivery Fee */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Delivery Fee ($)
              </label>
              <div className="relative">
                <DollarSign className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                <input
                  type="number"
                  step="0.01"
                  min="0"
                  {...register('deliveryFee', { valueAsNumber: true })}
                  className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                  placeholder="2.99"
                />
              </div>
              {errors.deliveryFee && (
                <p className="mt-1 text-sm text-red-600">{errors.deliveryFee.message}</p>
              )}
            </div>

            {/* Minimum Order Amount */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Minimum Order Amount ($)
              </label>
              <div className="relative">
                <DollarSign className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                <input
                  type="number"
                  step="0.01"
                  min="0"
                  {...register('minOrderAmount', { valueAsNumber: true })}
                  className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                  placeholder="10.00"
                />
              </div>
              {errors.minOrderAmount && (
                <p className="mt-1 text-sm text-red-600">{errors.minOrderAmount.message}</p>
              )}
            </div>

            {/* Tax Rate */}
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
                  placeholder="8.5"
                />
              </div>
              {errors.taxRate && (
                <p className="mt-1 text-sm text-red-600">{errors.taxRate.message}</p>
              )}
            </div>
          </div>
        </div>

        {/* Service Availability */}
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center space-x-3 mb-6">
            <Truck className="h-5 w-5 text-primary-500" />
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
              Service Availability
            </h2>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* Delivery Service */}
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
                  Delivery fee: ${deliveryFee.toFixed(2)}
                </div>
              )}
            </div>

            {/* Pickup Service */}
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
                  Minimum order: ${minOrderAmount.toFixed(2)}
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Order Summary Preview */}
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
                  <span>$50.00</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span>Delivery Fee:</span>
                  <span>${deliveryFee.toFixed(2)}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span>Tax ({taxRate}%):</span>
                  <span>${(50 * (taxRate / 100)).toFixed(2)}</span>
                </div>
                <div className="border-t border-gray-300 dark:border-gray-600 pt-2">
                  <div className="flex justify-between font-medium">
                    <span>Total:</span>
                    <span>${(50 + deliveryFee + (50 * (taxRate / 100))).toFixed(2)}</span>
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
                
                <div className="text-sm text-gray-600 dark:text-gray-400">
                  Minimum order for pickup: ${minOrderAmount.toFixed(2)}
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Action Buttons */}
        <div className="flex items-center justify-end space-x-3">
          <button
            type="button"
            onClick={() => reset()}
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