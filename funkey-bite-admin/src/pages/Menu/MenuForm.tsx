import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { 
  Save, 
  X, 
  Upload, 
  ChefHat, 
  Tag, 
  DollarSign, 
  Clock,
  Info,
  Image as ImageIcon,
  Plus,
  Trash2
} from 'lucide-react';
import { useQuery, useMutation } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { getCategories, createMenuItem, updateMenuItem, getMenuItemByID } from '../../api/adminApi';
import type { MenuItem, MenuCategory, NutritionalInfo } from '../../types';

const menuItemSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  description: z.string().min(10, 'Description must be at least 10 characters'),
  categoryId: z.number().min(1, 'Please select a category'),
  price: z.number().min(0.01, 'Price must be greater than 0'),
  preparationTime: z.number().min(1, 'Preparation time must be at least 1 minute'),
  isAvailable: z.boolean().default(true),
  isPreOrder: z.boolean().default(false),
  tags: z.array(z.string()).default([]),
  nutritionalInfo: z.object({
    calories: z.number().min(0).optional(),
    protein: z.number().min(0).optional(),
    carbs: z.number().min(0).optional(),
    fat: z.number().min(0).optional(),
  }).optional(),
});

type MenuItemFormData = z.infer<typeof menuItemSchema>;
type MenuItemFormInput = z.input<typeof menuItemSchema>;
type MenuItemFormOutput = z.output<typeof menuItemSchema>;

const MenuForm: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const isEditing = !!id;
  const [imagePreview, setImagePreview] = useState<string>('');
  const [tagInput, setTagInput] = useState('');
  const [showNutrition, setShowNutrition] = useState(false);

  const { data: categories = [] } = useQuery({
    queryKey: ['categories'],
    queryFn: async () => {
      try {
        const data = await getCategories();
        return Array.isArray(data) ? data : [];
      } catch (error) {
        toast.error('Failed to load categories');
        return [];
      }
    },
  });

  const { data: existingItem, isLoading: loadingItem } = useQuery({
    queryKey: ['menu-item', id],
    queryFn: async () => {
      if (!id) return null;
      try {
        const data = await getMenuItemByID(parseInt(id));
        return data;
      } catch (error) {
        toast.error('Failed to load menu item');
        navigate('/menu');
        return null;
      }
    },
    enabled: isEditing,
  });

  const {
    control,
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    setValue,
    watch,
    reset,
  } = useForm<MenuItemFormInput, unknown, MenuItemFormOutput>({
    resolver: zodResolver(menuItemSchema),
    defaultValues: {
      name: '',
      description: '',
      categoryId: 0,
      price: 0,
      preparationTime: 15,
      isAvailable: true,
      isPreOrder: false,
      tags: [],
    },
  });

  useEffect(() => {
    if (existingItem) {
      reset({
        name: existingItem.name,
        description: existingItem.description,
        categoryId: existingItem.categoryId,
        price: existingItem.price,
        preparationTime: existingItem.preparationTime,
        isAvailable: existingItem.isAvailable,
        isPreOrder: existingItem.isPreOrder,
        tags: existingItem.tags || [],
        nutritionalInfo: existingItem.nutritionalInfo,
      });
      if (existingItem.imageUrl) {
        setImagePreview(existingItem.imageUrl);
      }
      if (existingItem.nutritionalInfo) {
        setShowNutrition(true);
      }
    }
  }, [existingItem, reset]);

  const tags = watch('tags') ?? [];
  const nutritionalInfo = watch('nutritionalInfo');

  const addTag = () => {
    if (tagInput.trim() && !tags.includes(tagInput.trim())) {
      setValue('tags', [...tags, tagInput.trim()]);
      setTagInput('');
    }
  };

  const removeTag = (tagToRemove: string) => {
    setValue('tags', tags.filter(tag => tag !== tagToRemove));
  };

  const handleImageUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onloadend = () => {
        setImagePreview(reader.result as string);
      };
      reader.readAsDataURL(file);
    }
  };

  const mutation = useMutation({
    mutationFn: async (data: MenuItemFormOutput) => {
      const menuItemData: any = {
        ...data,
        imageUrl: imagePreview || 'https://via.placeholder.com/400x300?text=Menu+Item',
      };

      if (isEditing && id) {
        return await updateMenuItem(parseInt(id), menuItemData);
      } else {
        return await createMenuItem(menuItemData);
      }
    },
    onSuccess: () => {
      toast.success(isEditing ? 'Menu item updated!' : 'Menu item created!');
      navigate('/menu');
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to save menu item');
    },
  });

  const onSubmit = (data: MenuItemFormOutput) => {
    mutation.mutate(data);
  };

  if (loadingItem && isEditing) {
    return (
      <div className="text-center py-12">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-500 mx-auto"></div>
        <p className="mt-4 text-gray-500 dark:text-gray-400">Loading menu item...</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            {isEditing ? 'Edit Menu Item' : 'Add New Menu Item'}
          </h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            {isEditing ? 'Update your menu item details' : 'Create a new menu item for your restaurant'}
          </p>
        </div>
        
        <div className="flex items-center space-x-3">
          <button
            onClick={() => navigate('/menu')}
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
            <span>{isSubmitting ? 'Saving...' : 'Save Menu Item'}</span>
          </button>
        </div>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        
        <div className="lg:col-span-2 space-y-6">
          
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-6">
              Basic Information
            </h2>
            
            <div className="space-y-6">
              
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Item Name *
                </label>
                <input
                  type="text"
                  {...register('name')}
                  className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                  placeholder="e.g., Classic Shawarma Wrap"
                />
                {errors.name && (
                  <p className="mt-1 text-sm text-red-600">{errors.name.message}</p>
                )}
              </div>

              
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Description *
                </label>
                <textarea
                  {...register('description')}
                  rows={4}
                  className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                  placeholder="Describe your menu item in detail..."
                />
                {errors.description && (
                  <p className="mt-1 text-sm text-red-600">{errors.description.message}</p>
                )}
              </div>

              
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Category *
                  </label>
                  <select
                    {...register('categoryId', { valueAsNumber: true })}
                    className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                  >
                    <option value={0}>Select a category</option>
                    {categories.map((category) => (
                      <option key={category.id} value={category.id}>
                        {category.name}
                      </option>
                    ))}
                  </select>
                  {errors.categoryId && (
                    <p className="mt-1 text-sm text-red-600">{errors.categoryId.message}</p>
                  )}
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Price ($) *
                  </label>
                  <div className="relative">
                    <DollarSign className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                    <input
                      type="number"
                      step="0.01"
                      min="0"
                      {...register('price', { valueAsNumber: true })}
                      className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                      placeholder="0.00"
                    />
                  </div>
                  {errors.price && (
                    <p className="mt-1 text-sm text-red-600">{errors.price.message}</p>
                  )}
                </div>
              </div>

              
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Preparation Time (minutes) *
                  </label>
                  <div className="relative">
                    <Clock className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                    <input
                      type="number"
                      min="1"
                      {...register('preparationTime', { valueAsNumber: true })}
                      className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                      placeholder="15"
                    />
                  </div>
                  {errors.preparationTime && (
                    <p className="mt-1 text-sm text-red-600">{errors.preparationTime.message}</p>
                  )}
                </div>

                <div className="space-y-4">
                  <div className="flex items-center space-x-3">
                    <Controller
                      name="isAvailable"
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
                            <div className={`w-10 h-6 rounded-full transition-colors ${field.value ? 'bg-primary-500' : 'bg-gray-300 dark:bg-gray-600'}`}>
                              <div className={`absolute top-1 w-4 h-4 rounded-full bg-white transition-transform ${field.value ? 'left-5' : 'left-1'}`}></div>
                            </div>
                          </div>
                          <span className="text-sm text-gray-700 dark:text-gray-300">Available for order</span>
                        </label>
                      )}
                    />
                  </div>

                  <div className="flex items-center space-x-3">
                    <Controller
                      name="isPreOrder"
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
                          <span className="text-sm text-gray-700 dark:text-gray-300">Pre-order only</span>
                        </label>
                      )}
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>

          
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-6">
              Tags
            </h2>
            
            <div className="space-y-4">
              <div className="flex items-center space-x-2">
                <Tag className="h-5 w-5 text-gray-400" />
                <input
                  type="text"
                  value={tagInput}
                  onChange={(e) => setTagInput(e.target.value)}
                  onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), addTag())}
                  className="flex-1 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                  placeholder="Add tags (e.g., spicy, vegetarian, popular)"
                />
                <button
                  type="button"
                  onClick={addTag}
                  className="px-4 py-2 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 rounded-lg transition-colors"
                >
                  <Plus className="h-4 w-4" />
                </button>
              </div>

              {tags.length > 0 && (
                <div className="flex flex-wrap gap-2">
                  {tags.map((tag) => (
                    <span
                      key={tag}
                      className="inline-flex items-center space-x-2 px-3 py-1 bg-primary-100 dark:bg-primary-900 text-primary-800 dark:text-primary-200 rounded-full text-sm"
                    >
                      <span>{tag}</span>
                      <button
                        type="button"
                        onClick={() => removeTag(tag)}
                        className="hover:text-red-500"
                      >
                        <X className="h-3 w-3" />
                      </button>
                    </span>
                  ))}
                </div>
              )}
            </div>
          </div>

          
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                Nutritional Information
              </h2>
              <button
                type="button"
                onClick={() => setShowNutrition(!showNutrition)}
                className="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300"
              >
                {showNutrition ? 'Hide' : 'Add Nutritional Info'}
              </button>
            </div>

            {showNutrition && (
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Calories
                  </label>
                  <input
                    type="number"
                    min="0"
                    {...register('nutritionalInfo.calories', { valueAsNumber: true })}
                    className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                    placeholder="0"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Protein (g)
                  </label>
                  <input
                    type="number"
                    min="0"
                    {...register('nutritionalInfo.protein', { valueAsNumber: true })}
                    className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                    placeholder="0"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Carbs (g)
                  </label>
                  <input
                    type="number"
                    min="0"
                    {...register('nutritionalInfo.carbs', { valueAsNumber: true })}
                    className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                    placeholder="0"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Fat (g)
                  </label>
                  <input
                    type="number"
                    min="0"
                    {...register('nutritionalInfo.fat', { valueAsNumber: true })}
                    className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                    placeholder="0"
                  />
                </div>
              </div>
            )}
          </div>
        </div>

        
        <div className="space-y-6">
          
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-6">
              Item Image
            </h2>
            
            <div className="space-y-4">
              
              <div className="relative h-64 bg-gray-100 dark:bg-gray-700 rounded-lg overflow-hidden">
                {imagePreview ? (
                  <img
                    src={imagePreview}
                    alt="Preview"
                    className="w-full h-full object-cover"
                  />
                ) : (
                  <div className="w-full h-full flex flex-col items-center justify-center">
                    <ImageIcon className="h-12 w-12 text-gray-400 mb-4" />
                    <p className="text-gray-500 dark:text-gray-400">No image selected</p>
                  </div>
                )}
                
                
                {imagePreview && (
                  <button
                    type="button"
                    onClick={() => setImagePreview('')}
                    className="absolute top-3 right-3 p-2 bg-red-500 hover:bg-red-600 text-white rounded-full"
                  >
                    <Trash2 className="h-4 w-4" />
                  </button>
                )}
              </div>

              
              <label className="block">
                <input
                  type="file"
                  accept="image/*"
                  onChange={handleImageUpload}
                  className="hidden"
                />
                <div className="flex items-center justify-center space-x-2 px-4 py-3 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg hover:border-primary-500 dark:hover:border-primary-500 cursor-pointer transition-colors">
                  <Upload className="h-5 w-5 text-gray-400" />
                  <span className="text-gray-600 dark:text-gray-400">
                    {imagePreview ? 'Change Image' : 'Upload Image'}
                  </span>
                </div>
              </label>
              
              <p className="text-sm text-gray-500 dark:text-gray-400">
                Recommended: 400x300px JPG or PNG. Max 2MB.
              </p>
            </div>
          </div>

          
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-6">
              Live Preview
            </h2>
            
            <div className="space-y-4">
              <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-4">
                <div className="flex items-start justify-between mb-2">
                  <h3 className="font-semibold text-gray-900 dark:text-white truncate">
                    {watch('name') || 'Item Name'}
                  </h3>
                  <span className="font-bold text-primary-600 dark:text-primary-400">
                    ${watch('price') || '0.00'}
                  </span>
                </div>
                
                <p className="text-sm text-gray-500 dark:text-gray-400 mb-3 line-clamp-2">
                  {watch('description') || 'Item description will appear here...'}
                </p>
                
                
                <div className="flex flex-wrap gap-2">
                  {watch('isPreOrder') && (
                    <span className="px-2 py-1 text-xs bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 rounded">
                      Pre-order
                    </span>
                  )}
                  <span className={`px-2 py-1 text-xs rounded ${
                    watch('isAvailable')
                      ? 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200'
                      : 'bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200'
                  }`}>
                    {watch('isAvailable') ? 'Available' : 'Unavailable'}
                  </span>
                  <span className="px-2 py-1 text-xs bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-400 rounded">
                    {watch('preparationTime') || 15} min
                  </span>
                </div>
              </div>
              
              <p className="text-sm text-gray-500 dark:text-gray-400">
                This is how your menu item will appear to customers.
              </p>
            </div>
          </div>

          
          <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-xl p-6">
            <div className="flex items-start space-x-3">
              <Info className="h-5 w-5 text-blue-500 mt-0.5" />
              <div className="space-y-2">
                <h3 className="font-medium text-blue-800 dark:text-blue-300">
                  Tips for Better Menu Items
                </h3>
                <ul className="text-sm text-blue-700 dark:text-blue-400 space-y-1">
                  <li>• Use high-quality, appetizing images</li>
                  <li>• Write clear, descriptive names</li>
                  <li>• Include accurate preparation times</li>
                  <li>• Add relevant tags for filtering</li>
                  <li>• Set realistic prices</li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </form>
    </div>
  );
};

export default MenuForm;