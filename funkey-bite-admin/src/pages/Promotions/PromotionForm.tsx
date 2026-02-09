// src/pages/Promotions/PromotionForm.tsx
import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import {
  Save,
  X,
  ArrowLeft,
  Percent,
  DollarSign,
  Gift,
  Calendar,
  Hash,
  Users,
  ShoppingBag,
  Clock,
  Info,
  Sparkles,
  AlertCircle,
  Copy,
  Eye,
  TrendingUp,
  CheckCircle,
  XCircle,
  Zap
} from 'lucide-react';
import { useQuery, useMutation } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { getPromotions, createPromotion, updatePromotion, getPromotionByID } from '../../api/adminApi';
import type { Promotion, PromotionType } from '../../types';
import { usePromotionStore } from '../../stores/promotionStore';

// Validation schema
const promotionSchema = z.object({
  code: z.string()
    .min(3, 'Code must be at least 3 characters')
    .max(20, 'Code must be at most 20 characters')
    .regex(/^[A-Z0-9_-]+$/i, 'Code can only contain letters, numbers, hyphens and underscores'),
  
  title: z.string()
    .min(5, 'Title must be at least 5 characters')
    .max(100, 'Title must be at most 100 characters'),
  
  description: z.string()
    .min(10, 'Description must be at least 10 characters')
    .max(500, 'Description must be at most 500 characters')
    .optional(),
  
  promotionType: z.enum(['percentage', 'fixed', 'bogo']),
  
  discountValue: z.number()
    .min(0.01, 'Discount value must be greater than 0')
    .refine((val, ctx) => {
      const type = ctx.parent.promotionType;
      if (type === 'percentage' && val > 100) {
        return false;
      }
      return true;
    }, {
      message: 'Percentage discount cannot exceed 100%'
    }),
  
  maxDiscount: z.number()
    .min(0, 'Max discount must be positive')
    .optional()
    .nullable(),
  
  minOrderAmount: z.number()
    .min(0, 'Min order amount must be positive')
    .optional()
    .nullable(),
  
  usageLimit: z.number()
    .min(1, 'Usage limit must be at least 1')
    .optional()
    .nullable(),
  
  validFrom: z.string()
    .refine(date => !isNaN(Date.parse(date)), 'Invalid start date'),
  
  validUntil: z.string()
    .refine(date => !isNaN(Date.parse(date)), 'Invalid end date')
    .refine((date, ctx) => {
      const from = new Date(ctx.parent.validFrom);
      const until = new Date(date);
      return until > from;
    }, {
      message: 'End date must be after start date'
    }),
  
  isActive: z.boolean().default(true),
});

type PromotionFormData = z.infer<typeof promotionSchema>;

const PromotionForm: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const isEditing = !!id;
  const [isGeneratingCode, setIsGeneratingCode] = useState(false);
  const [showAdvanced, setShowAdvanced] = useState(false);
  
  const { setSelectedPromotion, addPromotion, updatePromotion: updateStorePromotion } = usePromotionStore();

  // Fetch existing promotion for editing
  const { data: existingPromotion, isLoading: loadingPromotion } = useQuery({
    queryKey: ['promotion', id],
    queryFn: async () => {
      if (!id) return null;
      try {
        const data = await getPromotionByID(parseInt(id));
        return data;
      } catch (error) {
        toast.error('Failed to load promotion');
        navigate('/promotions');
        return null;
      }
    },
    enabled: isEditing,
  });

  // Fetch all promotions to check for duplicate codes
  const { data: allPromotions = [] } = useQuery({
    queryKey: ['all-promotions'],
    queryFn: async () => {
      try {
        const data = await getPromotions({ limit: 1000 });
        return data.data || [];
      } catch (error) {
        return [];
      }
    },
  });

  const {
    control,
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    setValue,
    watch,
    reset,
    setError,
    clearErrors,
  } = useForm<PromotionFormData>({
    resolver: zodResolver(promotionSchema),
    defaultValues: {
      code: '',
      title: '',
      description: '',
      promotionType: 'percentage',
      discountValue: 10,
      maxDiscount: null,
      minOrderAmount: null,
      usageLimit: null,
      validFrom: new Date().toISOString().split('T')[0],
      validUntil: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
      isActive: true,
    },
  });

  // Set form values when editing
  useEffect(() => {
    if (existingPromotion) {
      const formData = {
        ...existingPromotion,
        validFrom: existingPromotion.validFrom.split('T')[0],
        validUntil: existingPromotion.validUntil.split('T')[0],
        maxDiscount: existingPromotion.maxDiscount || null,
        minOrderAmount: existingPromotion.minOrderAmount || null,
        usageLimit: existingPromotion.usageLimit || null,
        description: existingPromotion.description || '',
      };
      reset(formData);
    }
  }, [existingPromotion, reset]);

  // Watch form values
  const promotionType = watch('promotionType');
  const discountValue = watch('discountValue');
  const code = watch('code');
  const validFrom = watch('validFrom');
  const validUntil = watch('validUntil');

  // Generate promotion code
  const generatePromoCode = () => {
    setIsGeneratingCode(true);
    
    const prefixes = ['FUNKEY', 'GRAB', 'BITE', 'SAVE', 'TREAT', 'OFFER'];
    const suffix = Math.random().toString(36).substring(2, 6).toUpperCase();
    const prefix = prefixes[Math.floor(Math.random() * prefixes.length)];
    
    const generatedCode = `${prefix}${suffix}`;
    
    // Check if code exists
    const exists = allPromotions.some(promo => 
      promo.code.toLowerCase() === generatedCode.toLowerCase() && 
      (!isEditing || promo.id !== parseInt(id!))
    );
    
    if (exists) {
      // Retry once more
      const retrySuffix = Math.random().toString(36).substring(2, 6).toUpperCase();
      const retryCode = `${prefix}${retrySuffix}`;
      setValue('code', retryCode);
      clearErrors('code');
    } else {
      setValue('code', generatedCode);
      clearErrors('code');
    }
    
    setTimeout(() => setIsGeneratingCode(false), 500);
  };

  // Calculate end date based on duration
  const setDuration = (days: number) => {
    const fromDate = new Date(validFrom);
    const untilDate = new Date(fromDate);
    untilDate.setDate(untilDate.getDate() + days);
    setValue('validUntil', untilDate.toISOString().split('T')[0]);
  };

  // Calculate days remaining
  const getDaysRemaining = () => {
    const from = new Date(validFrom);
    const until = new Date(validUntil);
    const diffTime = until.getTime() - from.getTime();
    return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
  };

  const mutation = useMutation({
    mutationFn: async (data: PromotionFormData) => {
      const promotionData: any = {
        ...data,
        // Convert empty strings to null for optional fields
        maxDiscount: data.maxDiscount || null,
        minOrderAmount: data.minOrderAmount || null,
        usageLimit: data.usageLimit || null,
        description: data.description || null,
      };

      if (isEditing && id) {
        const result = await updatePromotion(parseInt(id), promotionData);
        updateStorePromotion(parseInt(id), result);
        return result;
      } else {
        const result = await createPromotion(promotionData);
        addPromotion(result);
        return result;
      }
    },
    onSuccess: (data) => {
      setSelectedPromotion(data);
      toast.success(isEditing ? 'Promotion updated!' : 'Promotion created!');
      navigate('/promotions');
    },
    onError: (error: any) => {
      if (error.message?.includes('already exists') || error.message?.includes('duplicate')) {
        setError('code', {
          type: 'manual',
          message: 'This promotion code already exists',
        });
      }
      toast.error(error.message || 'Failed to save promotion');
    },
  });

  const onSubmit = (data: PromotionFormData) => {
    // Check for duplicate code (except for current promotion)
    const duplicate = allPromotions.find(promo => 
      promo.code.toLowerCase() === data.code.toLowerCase() && 
      (!isEditing || promo.id !== parseInt(id!))
    );
    
    if (duplicate) {
      setError('code', {
        type: 'manual',
        message: 'This promotion code already exists',
      });
      return;
    }
    
    mutation.mutate(data);
  };

  // Preview calculations
  const getDiscountPreview = () => {
    const orderAmount = 50; // Sample order amount
    let discount = 0;
    
    switch (promotionType) {
      case 'percentage':
        discount = orderAmount * (discountValue / 100);
        if (watch('maxDiscount') && discount > watch('maxDiscount')!) {
          discount = watch('maxDiscount')!;
        }
        break;
      case 'fixed':
        discount = discountValue;
        break;
      case 'bogo':
        discount = orderAmount / 2; // Simplified BOGO calculation
        break;
    }
    
    return {
      original: orderAmount,
      discount: Math.min(discount, orderAmount),
      final: orderAmount - Math.min(discount, orderAmount),
    };
  };

  const preview = getDiscountPreview();

  if (loadingPromotion && isEditing) {
    return (
      <div className="text-center py-12">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-500 mx-auto"></div>
        <p className="mt-4 text-gray-500 dark:text-gray-400">Loading promotion...</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div className="flex items-center space-x-4">
          <button
            onClick={() => navigate('/promotions')}
            className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
          >
            <ArrowLeft className="h-5 w-5 text-gray-600 dark:text-gray-400" />
          </button>
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
              {isEditing ? 'Edit Promotion' : 'Create New Promotion'}
            </h1>
            <p className="text-gray-600 dark:text-gray-400 mt-1">
              {isEditing ? 'Update your promotion details' : 'Create discounts and special offers for customers'}
            </p>
          </div>
        </div>
        
        <div className="flex items-center space-x-3">
          <button
            onClick={() => navigate('/promotions')}
            className="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
          >
            Cancel
          </button>
          <button
            onClick={handleSubmit(onSubmit)}
            disabled={isSubmitting}
            className="flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors disabled:opacity-50"
          >
            <Save className="h-4 w-4" />
            <span>{isSubmitting ? 'Saving...' : (isEditing ? 'Update Promotion' : 'Create Promotion')}</span>
          </button>
        </div>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left Column - Basic Info */}
        <div className="lg:col-span-2 space-y-6">
          {/* Promotion Details */}
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-6">
              Promotion Details
            </h2>
            
            <div className="space-y-6">
              {/* Code & Title */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <div className="flex items-center justify-between mb-2">
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                      Promotion Code *
                    </label>
                    <button
                      type="button"
                      onClick={generatePromoCode}
                      disabled={isGeneratingCode}
                      className="text-xs text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300 flex items-center space-x-1"
                    >
                      {isGeneratingCode ? (
                        <>
                          <div className="animate-spin h-3 w-3 border-t-2 border-b-2 border-primary-500 rounded-full"></div>
                          <span>Generating...</span>
                        </>
                      ) : (
                        <>
                          <Sparkles className="h-3 w-3" />
                          <span>Generate Code</span>
                        </>
                      )}
                    </button>
                  </div>
                  <div className="relative">
                    <Hash className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                    <input
                      type="text"
                      {...register('code')}
                      className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white uppercase"
                      placeholder="e.g., FUNKEY20"
                      onChange={(e) => {
                        const value = e.target.value.toUpperCase();
                        e.target.value = value;
                        register('code').onChange(e);
                      }}
                    />
                  </div>
                  {errors.code && (
                    <p className="mt-1 text-sm text-red-600">{errors.code.message}</p>
                  )}
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Title *
                  </label>
                  <input
                    type="text"
                    {...register('title')}
                    className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                    placeholder="e.g., Summer Special - 20% Off"
                  />
                  {errors.title && (
                    <p className="mt-1 text-sm text-red-600">{errors.title.message}</p>
                  )}
                </div>
              </div>

              {/* Description */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Description
                </label>
                <textarea
                  {...register('description')}
                  rows={3}
                  className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                  placeholder="Describe this promotion for customers..."
                />
                {errors.description && (
                  <p className="mt-1 text-sm text-red-600">{errors.description.message}</p>
                )}
              </div>

              {/* Promotion Type */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-4">
                  Promotion Type *
                </label>
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                  {([
                    { type: 'percentage', icon: Percent, label: 'Percentage Off', desc: 'Discount by percentage' },
                    { type: 'fixed', icon: DollarSign, label: 'Fixed Amount', desc: 'Flat amount discount' },
                    { type: 'bogo', icon: Gift, label: 'Buy One Get One', desc: 'Buy one, get one free' },
                  ] as const).map(({ type, icon: Icon, label, desc }) => (
                    <Controller
                      key={type}
                      name="promotionType"
                      control={control}
                      render={({ field }) => (
                        <label className="cursor-pointer">
                          <input
                            type="radio"
                            className="sr-only"
                            checked={field.value === type}
                            onChange={() => field.onChange(type)}
                          />
                          <div className={`p-4 border-2 rounded-xl transition-all ${
                            field.value === type
                              ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20'
                              : 'border-gray-300 dark:border-gray-600 hover:border-primary-400 dark:hover:border-primary-500'
                          }`}>
                            <div className="flex items-center space-x-3">
                              <div className={`h-10 w-10 rounded-lg flex items-center justify-center ${
                                field.value === type
                                  ? 'bg-primary-500 text-white'
                                  : 'bg-gray-100 dark:bg-gray-700 text-gray-500'
                              }`}>
                                <Icon className="h-5 w-5" />
                              </div>
                              <div>
                                <p className="font-medium text-gray-900 dark:text-white">{label}</p>
                                <p className="text-xs text-gray-500 dark:text-gray-400">{desc}</p>
                              </div>
                            </div>
                          </div>
                        </label>
                      )}
                    />
                  ))}
                </div>
              </div>

              {/* Discount Value */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    {promotionType === 'percentage' ? 'Discount Percentage *' :
                     promotionType === 'fixed' ? 'Discount Amount ($) *' :
                     'Free Item Value ($) *'}
                  </label>
                  <div className="relative">
                    {promotionType === 'percentage' ? (
                      <Percent className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                    ) : (
                      <DollarSign className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                    )}
                    <input
                      type="number"
                      step={promotionType === 'percentage' ? 0.1 : 0.01}
                      min="0"
                      max={promotionType === 'percentage' ? 100 : undefined}
                      {...register('discountValue', { valueAsNumber: true })}
                      className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                      placeholder={promotionType === 'percentage' ? '20' : '5.00'}
                    />
                    {promotionType === 'percentage' && (
                      <span className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400">%</span>
                    )}
                  </div>
                  {errors.discountValue && (
                    <p className="mt-1 text-sm text-red-600">{errors.discountValue.message}</p>
                  )}
                </div>

                {promotionType === 'percentage' && (
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Maximum Discount ($)
                    </label>
                    <div className="relative">
                      <DollarSign className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                      <input
                        type="number"
                        step="0.01"
                        min="0"
                        {...register('maxDiscount', { valueAsNumber: true })}
                        className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                        placeholder="e.g., 10.00"
                      />
                    </div>
                    <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                      Leave empty for no maximum
                    </p>
                  </div>
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
                    placeholder="e.g., 25.00"
                  />
                </div>
                <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                  Leave empty for no minimum
                </p>
              </div>
            </div>
          </div>

          {/* Validity Period */}
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-6">
              Validity Period
            </h2>
            
            <div className="space-y-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Start Date *
                  </label>
                  <div className="relative">
                    <Calendar className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                    <input
                      type="date"
                      {...register('validFrom')}
                      min={new Date().toISOString().split('T')[0]}
                      className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                    />
                  </div>
                  {errors.validFrom && (
                    <p className="mt-1 text-sm text-red-600">{errors.validFrom.message}</p>
                  )}
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    End Date *
                  </label>
                  <div className="relative">
                    <Calendar className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                    <input
                      type="date"
                      {...register('validUntil')}
                      min={validFrom}
                      className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                    />
                  </div>
                  {errors.validUntil && (
                    <p className="mt-1 text-sm text-red-600">{errors.validUntil.message}</p>
                  )}
                </div>
              </div>

              {/* Quick Duration Buttons */}
              <div>
                <p className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
                  Quick Duration
                </p>
                <div className="flex flex-wrap gap-2">
                  {[7, 14, 30, 60, 90].map((days) => (
                    <button
                      type="button"
                      key={days}
                      onClick={() => setDuration(days)}
                      className="px-4 py-2 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 rounded-lg transition-colors"
                    >
                      {days} days
                    </button>
                  ))}
                </div>
              </div>

              {/* Days Remaining Display */}
              <div className="p-4 bg-gray-50 dark:bg-gray-700/50 rounded-lg">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-gray-700 dark:text-gray-300">
                      Promotion Duration
                    </p>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                      {getDaysRemaining()} days remaining
                    </p>
                  </div>
                  <Clock className="h-5 w-5 text-gray-400" />
                </div>
              </div>
            </div>
          </div>

          {/* Advanced Settings */}
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                Advanced Settings
              </h2>
              <button
                type="button"
                onClick={() => setShowAdvanced(!showAdvanced)}
                className="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300"
              >
                {showAdvanced ? 'Hide' : 'Show Advanced'}
              </button>
            </div>

            {showAdvanced && (
              <div className="space-y-6">
                {/* Usage Limit */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Usage Limit
                  </label>
                  <div className="relative">
                    <Users className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                    <input
                      type="number"
                      min="1"
                      {...register('usageLimit', { valueAsNumber: true })}
                      className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                      placeholder="e.g., 100"
                    />
                  </div>
                  <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                    Maximum number of times this promotion can be used
                  </p>
                </div>

                {/* Active Status */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-4">
                    Status
                  </label>
                  <div className="flex items-center space-x-6">
                    <Controller
                      name="isActive"
                      control={control}
                      render={({ field }) => (
                        <>
                          <label className="flex items-center space-x-2 cursor-pointer">
                            <div className="relative">
                              <input
                                type="radio"
                                checked={field.value === true}
                                onChange={() => field.onChange(true)}
                                className="sr-only"
                              />
                              <div className={`w-10 h-6 rounded-full transition-colors ${field.value === true ? 'bg-green-500' : 'bg-gray-300 dark:bg-gray-600'}`}>
                                <div className={`absolute top-1 w-4 h-4 rounded-full bg-white transition-transform ${field.value === true ? 'left-5' : 'left-1'}`}></div>
                              </div>
                            </div>
                            <span className="flex items-center space-x-2">
                              <CheckCircle className="h-4 w-4 text-green-500" />
                              <span className="text-sm text-gray-700 dark:text-gray-300">Active</span>
                            </span>
                          </label>
                          
                          <label className="flex items-center space-x-2 cursor-pointer">
                            <div className="relative">
                              <input
                                type="radio"
                                checked={field.value === false}
                                onChange={() => field.onChange(false)}
                                className="sr-only"
                              />
                              <div className={`w-10 h-6 rounded-full transition-colors ${field.value === false ? 'bg-red-500' : 'bg-gray-300 dark:bg-gray-600'}`}>
                                <div className={`absolute top-1 w-4 h-4 rounded-full bg-white transition-transform ${field.value === false ? 'left-5' : 'left-1'}`}></div>
                              </div>
                            </div>
                            <span className="flex items-center space-x-2">
                              <XCircle className="h-4 w-4 text-red-500" />
                              <span className="text-sm text-gray-700 dark:text-gray-300">Inactive</span>
                            </span>
                          </label>
                        </>
                      )}
                    />
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Right Column - Preview & Tips */}
        <div className="space-y-6">
          {/* Live Preview */}
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-6">
              Live Preview
            </h2>
            
            <div className="space-y-6">
              {/* Promotion Card Preview */}
              <div className="bg-gradient-to-r from-primary-500 to-primary-600 rounded-xl p-6 text-white">
                <div className="flex items-start justify-between mb-4">
                  <div>
                    <div className="flex items-center space-x-2 mb-2">
                      <Zap className="h-5 w-5" />
                      <h3 className="text-xl font-bold">{watch('title') || 'Promotion Title'}</h3>
                    </div>
                    <p className="text-primary-100 text-sm">
                      {watch('description') || 'Promotion description will appear here...'}
                    </p>
                  </div>
                  <div className="bg-white/20 backdrop-blur-sm rounded-lg px-3 py-1">
                    <span className="font-mono font-bold text-lg">
                      {watch('code') || 'CODE123'}
                    </span>
                  </div>
                </div>
                
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-primary-100 text-sm">Discount Value</p>
                    <p className="text-2xl font-bold">
                      {promotionType === 'percentage' ? `${discountValue}%` :
                       promotionType === 'fixed' ? `$${discountValue.toFixed(2)}` :
                       'BOGO'}
                    </p>
                  </div>
                  <div className="text-right">
                    <p className="text-primary-100 text-sm">Valid Until</p>
                    <p className="font-semibold">
                      {validUntil ? new Date(validUntil).toLocaleDateString('en-US', {
                        month: 'short',
                        day: 'numeric'
                      }) : 'End Date'}
                    </p>
                  </div>
                </div>
              </div>

              {/* Discount Calculator */}
              <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-4">
                <h3 className="font-medium text-gray-900 dark:text-white mb-3">
                  Discount Calculator
                </h3>
                
                <div className="space-y-3">
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-gray-600 dark:text-gray-400">Sample Order:</span>
                    <span className="font-medium">${preview.original.toFixed(2)}</span>
                  </div>
                  
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-gray-600 dark:text-gray-400">Discount Applied:</span>
                    <span className="font-medium text-green-600 dark:text-green-400">
                      -${preview.discount.toFixed(2)}
                    </span>
                  </div>
                  
                  <div className="border-t border-gray-300 dark:border-gray-600 pt-2">
                    <div className="flex items-center justify-between">
                      <span className="font-medium text-gray-900 dark:text-white">Final Amount:</span>
                      <span className="text-xl font-bold text-primary-600 dark:text-primary-400">
                        ${preview.final.toFixed(2)}
                      </span>
                    </div>
                  </div>
                </div>
              </div>

              {/* Status Badges */}
              <div className="space-y-3">
                <p className="text-sm font-medium text-gray-700 dark:text-gray-300">
                  Promotion Status
                </p>
                <div className="flex flex-wrap gap-2">
                  {watch('isActive') ? (
                    <span className="px-3 py-1 bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200 rounded-full text-sm font-medium inline-flex items-center space-x-1">
                      <CheckCircle className="h-3 w-3" />
                      <span>Active</span>
                    </span>
                  ) : (
                    <span className="px-3 py-1 bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200 rounded-full text-sm font-medium inline-flex items-center space-x-1">
                      <XCircle className="h-3 w-3" />
                      <span>Inactive</span>
                    </span>
                  )}
                  
                  {watch('minOrderAmount') && (
                    <span className="px-3 py-1 bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 rounded-full text-sm font-medium inline-flex items-center space-x-1">
                      <ShoppingBag className="h-3 w-3" />
                      <span>Min ${watch('minOrderAmount')}</span>
                    </span>
                  )}
                  
                  {watch('usageLimit') && (
                    <span className="px-3 py-1 bg-purple-100 dark:bg-purple-900 text-purple-800 dark:text-purple-200 rounded-full text-sm font-medium inline-flex items-center space-x-1">
                      <Users className="h-3 w-3" />
                      <span>Limit: {watch('usageLimit')}</span>
                    </span>
                  )}
                </div>
              </div>
            </div>
          </div>

          {/* Tips & Best Practices */}
          <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-xl p-6">
            <div className="flex items-start space-x-3">
              <Info className="h-5 w-5 text-blue-500 mt-0.5" />
              <div className="space-y-3">
                <h3 className="font-medium text-blue-800 dark:text-blue-300">
                  Promotion Tips
                </h3>
                <ul className="text-sm text-blue-700 dark:text-blue-400 space-y-2">
                  <li className="flex items-start space-x-2">
                    <TrendingUp className="h-4 w-4 text-blue-500 mt-0.5 flex-shrink-0" />
                    <span>Use clear, memorable codes (e.g., SUMMER20, FUNKEYSAVE)</span>
                  </li>
                  <li className="flex items-start space-x-2">
                    <Calendar className="h-4 w-4 text-blue-500 mt-0.5 flex-shrink-0" />
                    <span>Set reasonable expiration dates (2-4 weeks works well)</span>
                  </li>
                  <li className="flex items-start space-x-2">
                    <AlertCircle className="h-4 w-4 text-blue-500 mt-0.5 flex-shrink-0" />
                    <span>Add usage limits for high-value promotions</span>
                  </li>
                  <li className="flex items-start space-x-2">
                    <Copy className="h-4 w-4 text-blue-500 mt-0.5 flex-shrink-0" />
                    <span>Test codes before making them public</span>
                  </li>
                  <li className="flex items-start space-x-2">
                    <Eye className="h-4 w-4 text-blue-500 mt-0.5 flex-shrink-0" />
                    <span>Review performance regularly in the promotions dashboard</span>
                  </li>
                </ul>
              </div>
            </div>
          </div>

          {/* Form Summary */}
          <div className="bg-gray-50 dark:bg-gray-700/50 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <h3 className="font-medium text-gray-900 dark:text-white mb-4">
              Form Summary
            </h3>
            <div className="space-y-3">
              {Object.entries(errors).map(([field, error]) => (
                <div key={field} className="flex items-start space-x-2 text-sm">
                  <AlertCircle className="h-4 w-4 text-red-500 mt-0.5 flex-shrink-0" />
                  <div>
                    <span className="font-medium text-red-600 dark:text-red-400 capitalize">
                      {field.replace(/([A-Z])/g, ' $1').trim()}:
                    </span>
                    <span className="ml-2 text-red-600 dark:text-red-400">{error.message}</span>
                  </div>
                </div>
              ))}
              
              {Object.keys(errors).length === 0 && (
                <div className="flex items-center space-x-2 text-sm text-green-600 dark:text-green-400">
                  <CheckCircle className="h-4 w-4" />
                  <span>All required fields are valid</span>
                </div>
              )}
            </div>
          </div>
        </div>
      </form>
    </div>
  );
};

export default PromotionForm;