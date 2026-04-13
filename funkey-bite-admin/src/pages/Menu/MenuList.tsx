import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { 
  Search, 
  Filter, 
  Plus, 
  Edit, 
  Trash2, 
  Eye, 
  ChefHat,
  Tag,
  DollarSign,
  Clock,
  CheckCircle,
  XCircle,
  ChevronDown,
  Grid,
  List
} from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { getMenuItems, deleteMenuItem } from '../../api/adminApi';
import type { MenuItem } from '../../types';

const MenuList: React.FC = () => {
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(20);
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');
  const [searchQuery, setSearchQuery] = useState('');
  const [categoryFilter, setCategoryFilter] = useState<string>('');

  const { data: menuItems, isLoading, refetch } = useQuery({
    queryKey: ['menu-items', page, limit, categoryFilter, searchQuery],
    queryFn: async () => {
      try {
        const params: any = { page, limit };
        if (categoryFilter) params.categoryId = categoryFilter;
        if (searchQuery) params.query = searchQuery;
        
        const data = await getMenuItems(params);
        return Array.isArray(data) ? data : [];
      } catch (error) {
        toast.error('Failed to load menu items');
        return [];
      }
    },
  });

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this menu item?')) return;
    
    try {
      await deleteMenuItem(id);
      toast.success('Menu item deleted successfully');
      refetch();
    } catch (error) {
      toast.error('Failed to delete menu item');
    }
  };

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(price);
  };

  return (
    <div className="space-y-6">
      
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Menu Management</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Manage your restaurant menu items and categories
          </p>
        </div>
        
        <div className="flex items-center space-x-3">
          <Link
            to="/menu/new"
            className="flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors"
          >
            <Plus className="h-4 w-4" />
            <span>Add Menu Item</span>
          </Link>
        </div>
      </div>

      
      <div className="bg-white dark:bg-gray-800 rounded-xl p-4 border border-gray-200 dark:border-gray-700">
        <div className="flex flex-col lg:flex-row gap-4">
          
          <div className="flex-1">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
              <input
                type="search"
                placeholder="Search menu items by name or description..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full pl-10 pr-4 py-2 bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:text-white"
              />
            </div>
          </div>

          
          <div className="flex items-center space-x-4">
            <button
              onClick={() => setViewMode(viewMode === 'grid' ? 'list' : 'grid')}
              className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
            >
              {viewMode === 'grid' ? (
                <List className="h-5 w-5 text-gray-400" />
              ) : (
                <Grid className="h-5 w-5 text-gray-400" />
              )}
            </button>
            
            <Filter className="h-5 w-5 text-gray-400" />
            <select
              value={categoryFilter}
              onChange={(e) => setCategoryFilter(e.target.value)}
              className="px-4 py-2 bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:text-white"
            >
              <option value="">All Categories</option>
              <option value="1">Chips & Chicken</option>
              <option value="2">Noodles</option>
              <option value="3">Shawarma</option>
              <option value="4">Drinks</option>
              <option value="5">Soup & Food Bowls</option>
            </select>
          </div>
        </div>
      </div>

      
      {isLoading ? (
        <div className="text-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-500 mx-auto"></div>
          <p className="mt-4 text-gray-500 dark:text-gray-400">Loading menu items...</p>
        </div>
      ) : menuItems?.length === 0 ? (
        <div className="text-center py-12">
          <ChefHat className="h-16 w-16 text-gray-400 mx-auto mb-4" />
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
            No Menu Items Found
          </h2>
          <p className="text-gray-500 dark:text-gray-400 mb-6">
            Get started by adding your first menu item.
          </p>
          <Link
            to="/menu/new"
            className="inline-flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors"
          >
            <Plus className="h-4 w-4" />
            <span>Add Menu Item</span>
          </Link>
        </div>
      ) : viewMode === 'grid' ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {menuItems?.map((item) => (
            <div
              key={item.id}
              className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden hover:shadow-lg transition-shadow"
            >
              
              <div className="h-48 bg-gray-200 dark:bg-gray-700 relative">
                {item.imageUrl ? (
                  <img
                    src={item.imageUrl}
                    alt={item.name}
                    className="w-full h-full object-cover"
                  />
                ) : (
                  <div className="w-full h-full flex items-center justify-center">
                    <ChefHat className="h-12 w-12 text-gray-400" />
                  </div>
                )}
                
                
                <div className="absolute top-3 right-3">
                  <span className={`px-3 py-1 rounded-full text-xs font-medium ${
                    item.isAvailable 
                      ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                      : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                  }`}>
                    {item.isAvailable ? 'Available' : 'Unavailable'}
                  </span>
                </div>
                
                
                {item.isPreOrder && (
                  <div className="absolute top-3 left-3">
                    <span className="px-3 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200">
                      Pre-order
                    </span>
                  </div>
                )}
              </div>

              
              <div className="p-4">
                <div className="flex items-start justify-between mb-2">
                  <h3 className="font-semibold text-gray-900 dark:text-white truncate">
                    {item.name}
                  </h3>
                  <span className="font-bold text-primary-600 dark:text-primary-400">
                    {formatPrice(item.price)}
                  </span>
                </div>
                
                <p className="text-sm text-gray-500 dark:text-gray-400 mb-3 line-clamp-2">
                  {item.description}
                </p>
                
                
                {item.tags && item.tags.length > 0 && (
                  <div className="flex flex-wrap gap-1 mb-3">
                    {item.tags.slice(0, 3).map((tag, index) => (
                      <span
                        key={index}
                        className="px-2 py-1 text-xs bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-400 rounded"
                      >
                        {tag}
                      </span>
                    ))}
                    {item.tags.length > 3 && (
                      <span className="px-2 py-1 text-xs text-gray-400">
                        +{item.tags.length - 3}
                      </span>
                    )}
                  </div>
                )}
                
                
                <div className="flex items-center text-sm text-gray-500 dark:text-gray-400 mb-4">
                  <Clock className="h-4 w-4 mr-1" />
                  <span>{item.preparationTime} min</span>
                </div>
                
                
                <div className="flex items-center justify-between pt-4 border-t border-gray-200 dark:border-gray-700">
                  <div className="flex items-center space-x-2">
                    <Link
                      to={`/menu/${item.id}/edit`}
                      className="p-1 text-gray-400 hover:text-primary-500 transition-colors"
                      title="Edit"
                    >
                      <Edit className="h-4 w-4" />
                    </Link>
                    <button
                      onClick={() => handleDelete(item.id)}
                      className="p-1 text-gray-400 hover:text-red-500 transition-colors"
                      title="Delete"
                    >
                      <Trash2 className="h-4 w-4" />
                    </button>
                  </div>
                  <Link
                    to={`/menu/${item.id}`}
                    className="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300 font-medium"
                  >
                    View Details →
                  </Link>
                </div>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 dark:bg-gray-700/50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Item
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Category
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Price
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Status
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Prep Time
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                {menuItems?.map((item) => (
                  <tr key={item.id} className="hover:bg-gray-50 dark:hover:bg-gray-700/30">
                    <td className="px-6 py-4">
                      <div className="flex items-center space-x-3">
                        <div className="h-10 w-10 rounded-lg bg-gray-200 dark:bg-gray-700 flex items-center justify-center">
                          {item.imageUrl ? (
                            <img
                              src={item.imageUrl}
                              alt={item.name}
                              className="h-full w-full object-cover rounded-lg"
                            />
                          ) : (
                            <ChefHat className="h-5 w-5 text-gray-400" />
                          )}
                        </div>
                        <div>
                          <div className="font-medium text-gray-900 dark:text-white">
                            {item.name}
                          </div>
                          <div className="text-sm text-gray-500 dark:text-gray-400 truncate max-w-xs">
                            {item.description}
                          </div>
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <span className="px-2 py-1 text-xs rounded-full bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300">
                        Category {item.categoryId}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      <div className="font-semibold text-gray-900 dark:text-white">
                        {formatPrice(item.price)}
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex items-center space-x-2">
                        {item.isAvailable ? (
                          <CheckCircle className="h-4 w-4 text-green-500" />
                        ) : (
                          <XCircle className="h-4 w-4 text-red-500" />
                        )}
                        <span className={`text-sm ${item.isAvailable ? 'text-green-600' : 'text-red-600'}`}>
                          {item.isAvailable ? 'Available' : 'Unavailable'}
                        </span>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex items-center text-gray-600 dark:text-gray-400">
                        <Clock className="h-4 w-4 mr-1" />
                        <span>{item.preparationTime} min</span>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex items-center space-x-2">
                        <Link
                          to={`/menu/${item.id}`}
                          className="p-1 text-gray-400 hover:text-primary-500 transition-colors"
                          title="View"
                        >
                          <Eye className="h-4 w-4" />
                        </Link>
                        <Link
                          to={`/menu/${item.id}/edit`}
                          className="p-1 text-gray-400 hover:text-primary-500 transition-colors"
                          title="Edit"
                        >
                          <Edit className="h-4 w-4" />
                        </Link>
                        <button
                          onClick={() => handleDelete(item.id)}
                          className="p-1 text-gray-400 hover:text-red-500 transition-colors"
                          title="Delete"
                        >
                          <Trash2 className="h-4 w-4" />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      
      {menuItems && menuItems.length > 0 && (
        <div className="flex items-center justify-between">
          <div className="text-sm text-gray-500 dark:text-gray-400">
            Showing {menuItems.length} items
          </div>
          <div className="flex items-center space-x-2">
            <button
              onClick={() => setLimit(20)}
              className={`px-3 py-1 rounded-lg ${limit === 20 ? 'bg-primary-500 text-white' : 'hover:bg-gray-100 dark:hover:bg-gray-700'}`}
            >
              20
            </button>
            <button
              onClick={() => setLimit(50)}
              className={`px-3 py-1 rounded-lg ${limit === 50 ? 'bg-primary-500 text-white' : 'hover:bg-gray-100 dark:hover:bg-gray-700'}`}
            >
              50
            </button>
            <button
              onClick={() => setLimit(100)}
              className={`px-3 py-1 rounded-lg ${limit === 100 ? 'bg-primary-500 text-white' : 'hover:bg-gray-100 dark:hover:bg-gray-700'}`}
            >
              100
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default MenuList;