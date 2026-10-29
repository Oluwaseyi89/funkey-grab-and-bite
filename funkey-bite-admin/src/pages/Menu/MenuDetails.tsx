import React from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { ArrowLeft, ChefHat, Clock, Edit, Tag, Trash2 } from 'lucide-react';
import toast from 'react-hot-toast';
import { deleteMenuItem, getCategories, getMenuItemByID } from '../../api/adminApi';
import StockBadge from './StockBadge';

const formatPrice = (price: number) =>
  new Intl.NumberFormat('en-NG', { style: 'currency', currency: 'NGN' }).format(price);

const formatDate = (value: string) =>
  new Date(value).toLocaleString('en-NG', { dateStyle: 'medium', timeStyle: 'short' });

const Detail: React.FC<{ label: string; children: React.ReactNode }> = ({ label, children }) => (
  <div>
    <dt className="text-sm text-gray-500 dark:text-gray-400">{label}</dt>
    <dd className="mt-1 font-medium text-gray-900 dark:text-white">{children}</dd>
  </div>
);

const MenuDetails: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const itemId = Number(id);

  const { data: item, isLoading, isError } = useQuery({
    queryKey: ['menu-item', id],
    queryFn: () => getMenuItemByID(itemId),
    enabled: Number.isInteger(itemId) && itemId > 0,
    retry: false,
  });

  const { data: categories = [] } = useQuery({
    queryKey: ['categories'],
    queryFn: async () => {
      const data = await getCategories();
      return Array.isArray(data) ? data : [];
    },
  });

  const handleDelete = async () => {
    if (!item || !confirm('Are you sure you want to delete this menu item?')) return;

    try {
      await deleteMenuItem(item.id);
      toast.success('Menu item deleted successfully');
      navigate('/menu');
    } catch {
      toast.error('Failed to delete menu item');
    }
  };

  if (isLoading) {
    return (
      <div className="space-y-4">
        <div className="h-6 bg-gray-200 dark:bg-gray-700 rounded w-48 animate-pulse"></div>
        <div className="h-64 bg-gray-200 dark:bg-gray-700 rounded-xl animate-pulse"></div>
      </div>
    );
  }

  if (isError || !item) {
    return (
      <div className="text-center py-12">
        <ChefHat className="h-16 w-16 text-gray-400 mx-auto mb-4" />
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">Menu Item Not Found</h2>
        <p className="text-gray-500 dark:text-gray-400 mb-6">
          This menu item doesn't exist or could not be loaded.
        </p>
        <Link
          to="/menu"
          className="inline-flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors"
        >
          <ArrowLeft className="h-4 w-4" />
          <span>Back to Menu</span>
        </Link>
      </div>
    );
  }

  const category = categories.find((c) => c.id === item.categoryId);
  const nutrition = item.nutritionalInfo;
  const hasNutrition = nutrition && Object.values(nutrition).some((value) => Number(value) > 0);

  return (
    <div className="space-y-6">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div className="flex items-center space-x-4">
          <Link
            to="/menu"
            className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
            title="Back to Menu"
          >
            <ArrowLeft className="h-5 w-5" />
          </Link>
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white">{item.name}</h1>
            <p className="text-gray-600 dark:text-gray-400">{category?.name ?? 'Uncategorised'}</p>
          </div>
        </div>

        <div className="flex items-center space-x-3">
          <button
            onClick={handleDelete}
            className="flex items-center space-x-2 px-4 py-2 border border-red-300 text-red-600 dark:border-red-800 dark:text-red-400 rounded-lg hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors"
          >
            <Trash2 className="h-4 w-4" />
            <span>Delete</span>
          </button>
          <Link
            to={`/menu/${item.id}/edit`}
            className="flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors"
          >
            <Edit className="h-4 w-4" />
            <span>Edit Item</span>
          </Link>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="h-72 lg:h-full min-h-72 bg-gray-200 dark:bg-gray-700 rounded-xl overflow-hidden">
          {item.imageUrl ? (
            <img src={item.imageUrl} alt={item.name} className="w-full h-full object-cover" />
          ) : (
            <div className="w-full h-full flex items-center justify-center">
              <ChefHat className="h-16 w-16 text-gray-400" />
            </div>
          )}
        </div>

        <div className="space-y-6">
          <div className="flex flex-wrap items-center gap-2">
            <span className="text-2xl font-bold text-primary-600 dark:text-primary-400 mr-2">
              {formatPrice(item.price)}
            </span>
            <span
              className={`px-3 py-1 rounded-full text-xs font-medium ${
                item.isAvailable
                  ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                  : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
              }`}
            >
              {item.isAvailable ? 'Available' : 'Unavailable'}
            </span>
            {item.isPreOrder && (
              <span className="px-3 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200">
                Pre-order
              </span>
            )}
          </div>

          <p className="text-gray-700 dark:text-gray-300 whitespace-pre-line">
            {item.description || 'No description.'}
          </p>

          <dl className="grid grid-cols-2 gap-4">
            <Detail label="Stock">
              <StockBadge item={item} />
            </Detail>
            <Detail label="Preparation time">
              <span className="inline-flex items-center">
                <Clock className="h-4 w-4 mr-1 text-gray-400" />
                {item.preparationTime} min
              </span>
            </Detail>
            <Detail label="Added">{formatDate(item.createdAt)}</Detail>
            {/* Items never edited come back with the zero date (year 1). */}
            {item.updatedAt && !item.updatedAt.startsWith('0001') && (
              <Detail label="Last updated">{formatDate(item.updatedAt)}</Detail>
            )}
          </dl>

          {item.tags && item.tags.length > 0 && (
            <div>
              <h3 className="text-sm text-gray-500 dark:text-gray-400 mb-2 flex items-center">
                <Tag className="h-4 w-4 mr-1" />
                Tags
              </h3>
              <div className="flex flex-wrap gap-2">
                {item.tags.map((tag) => (
                  <span
                    key={tag}
                    className="px-2 py-1 text-xs bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-400 rounded"
                  >
                    {tag}
                  </span>
                ))}
              </div>
            </div>
          )}

          {hasNutrition && nutrition && (
            <div>
              <h3 className="text-sm text-gray-500 dark:text-gray-400 mb-2">Nutrition per serving</h3>
              <dl className="grid grid-cols-4 gap-2 text-center">
                {([
                  ['Calories', nutrition.calories, ''],
                  ['Protein', nutrition.protein, 'g'],
                  ['Carbs', nutrition.carbs, 'g'],
                  ['Fat', nutrition.fat, 'g'],
                ] as const).map(([label, value, unit]) => (
                  <div key={label} className="bg-gray-50 dark:bg-gray-700 rounded-lg p-3">
                    <dd className="font-semibold text-gray-900 dark:text-white">{value ?? 0}{unit}</dd>
                    <dt className="text-xs text-gray-500 dark:text-gray-400">{label}</dt>
                  </div>
                ))}
              </dl>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default MenuDetails;
