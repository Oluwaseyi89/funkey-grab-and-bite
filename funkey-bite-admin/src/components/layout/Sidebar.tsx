import React, { useState } from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import {
  Home,
  Package,
  Users,
  ChefHat,
  ShoppingCart,
  BarChart3,
  Settings,
  Tag,
  BellRing,
  LogOut,
  Menu as MenuIcon,
  X,
  FolderOpen,
  PieChart,
  Package2,
} from 'lucide-react';
import { useAuthStore } from '../../stores/authStore';
import { useTheme } from '../../contexts/ThemeContext';

const Sidebar: React.FC = () => {
  const [isOpen, setIsOpen] = useState(false);
  const { logout } = useAuthStore();
  const { theme } = useTheme();
  const location = useLocation();

  const navigation = [
    { name: 'Dashboard', href: '/dashboard', icon: Home },
    { name: 'Orders', href: '/orders', icon: Package },
    { name: 'Menu', href: '/menu', icon: ChefHat },
    { name: 'Customers', href: '/customers', icon: Users },
    { name: 'Catering', href: '/catering', icon: ShoppingCart },
    { name: 'Inventory', href: '/inventory', icon: Package2 },
    { name: 'Promotions', href: '/promotions', icon: Tag },
    { name: 'Reports', href: '/reports', icon: BarChart3 },
    { name: 'Settings', href: '/settings', icon: Settings },
  ];

  return (
    <>
      {/* Mobile Toggle Button */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="lg:hidden fixed top-4 left-4 z-50 p-2 rounded-lg bg-primary-500 text-white shadow-lg"
      >
        {isOpen ? <X className="h-6 w-6" /> : <MenuIcon className="h-6 w-6" />}
      </button>

      {/* Overlay for mobile */}
      {isOpen && (
        <div
          className="lg:hidden fixed inset-0 bg-black bg-opacity-50 z-40"
          onClick={() => setIsOpen(false)}
        />
      )}

      {/* Sidebar */}
      <aside
        className={`${
          isOpen ? 'translate-x-0' : '-translate-x-full'
        } lg:translate-x-0 fixed lg:static inset-y-0 left-0 z-40 w-64 bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 transform transition-transform duration-200 ease-in-out flex flex-col h-full`}
      >
        {/* Logo */}
        <div className="p-6 border-b border-gray-200 dark:border-gray-700">
          <div className="flex items-center space-x-3">
            <div className="h-10 w-10 rounded-lg bg-primary-500 flex items-center justify-center">
              <span className="text-white font-bold text-lg">FG&B</span>
            </div>
            <div>
              <h1 className="text-xl font-bold text-gray-900 dark:text-white">
                Funkey Grab & Bite
              </h1>
              <p className="text-sm text-gray-500 dark:text-gray-400">Admin Panel</p>
            </div>
          </div>
        </div>

        {/* Navigation */}
        <nav className="flex-1 p-4 space-y-1 overflow-y-auto">
          {navigation.map((item) => {
            const Icon = item.icon;
            const isActive = location.pathname === item.href || 
                           location.pathname.startsWith(item.href + '/');
            
            return (
              <NavLink
                key={item.name}
                to={item.href}
                onClick={() => setIsOpen(false)}
                className={`
                  flex items-center space-x-3 px-4 py-3 rounded-lg transition-colors
                  ${isActive
                    ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-600 dark:text-primary-400'
                    : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
                  }
                `}
              >
                <Icon className="h-5 w-5 flex-shrink-0" />
                <span className="font-medium">{item.name}</span>
                {isActive && (
                  <div className="ml-auto h-2 w-2 rounded-full bg-primary-500" />
                )}
              </NavLink>
            );
          })}
        </nav>

        {/* User Info & Logout */}
        <div className="p-4 border-t border-gray-200 dark:border-gray-700">
          <div className="flex items-center space-x-3 mb-4">
            <div className="h-10 w-10 rounded-full bg-primary-100 dark:bg-primary-900 flex items-center justify-center">
              <span className="text-primary-600 dark:text-primary-300 font-semibold">
                A
              </span>
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-gray-900 dark:text-white truncate">
                Admin User
              </p>
              <p className="text-xs text-gray-500 dark:text-gray-400 truncate">
                admin@funkey.com
              </p>
            </div>
          </div>
          
          <button
            onClick={() => {
              logout();
              localStorage.removeItem('admin_token');
              window.location.href = '/login';
            }}
            className="flex items-center space-x-3 w-full px-4 py-3 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors"
          >
            <LogOut className="h-5 w-5" />
            <span className="font-medium">Logout</span>
          </button>
        </div>
      </aside>
    </>
  );
};

export default Sidebar;


























// // src/components/layout/Sidebar.tsx
// import React, { useState } from 'react';
// import { NavLink } from 'react-router-dom';
// import {
//   Home,
//   Package,
//   Users,
//   ChefHat,
//   ShoppingCart,
//   BarChart3,
//   Settings,
//   Tag,
//   BellRing,
//   LogOut,
//   Menu,
//   X,
// } from 'lucide-react';
// import { useAuthStore } from '../../stores/authStore';
// import logo from '../../assets/logo.svg';

// const Sidebar: React.FC = () => {
//   const [isOpen, setIsOpen] = useState(false);
//   const { logout } = useAuthStore();

//   const navigation = [
//     { name: 'Dashboard', href: '/dashboard', icon: Home },
//     { name: 'Orders', href: '/orders', icon: Package },
//     { name: 'Menu', href: '/menu', icon: ChefHat },
//     { name: 'Customers', href: '/customers', icon: Users },
//     { name: 'Catering', href: '/catering', icon: ShoppingCart },
//     { name: 'Inventory', href: '/inventory', icon: BarChart3 },
//     { name: 'Promotions', href: '/promotions', icon: Tag },
//     { name: 'Reports', href: '/reports', icon: BarChart3 },
//     { name: 'Settings', href: '/settings', icon: Settings },
//   ];

//   return (
//     <>
//       {/* Mobile Toggle Button */}
//       <button
//         onClick={() => setIsOpen(!isOpen)}
//         className="lg:hidden fixed top-4 left-4 z-50 p-2 rounded-lg bg-primary-500 text-white"
//       >
//         {isOpen ? <X className="h-6 w-6" /> : <Menu className="h-6 w-6" />}
//       </button>

//       {/* Overlay for mobile */}
//       {isOpen && (
//         <div
//           className="lg:hidden fixed inset-0 bg-black bg-opacity-50 z-40"
//           onClick={() => setIsOpen(false)}
//         />
//       )}

//       {/* Sidebar */}
//       <aside
//         className={`${
//           isOpen ? 'translate-x-0' : '-translate-x-full'
//         } lg:translate-x-0 fixed lg:static inset-y-0 left-0 z-40 w-64 bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 transform transition-transform duration-200 ease-in-out`}
//       >
//         {/* Logo */}
//         <div className="p-6 border-b border-gray-200 dark:border-gray-700">
//           <div className="flex items-center space-x-3">
//             <img src={logo} alt="Logo" className="h-10" />
//             <div>
//               <h1 className="text-xl font-bold text-gray-900 dark:text-white">
//                 Funkey Grab & Bite
//               </h1>
//               <p className="text-sm text-gray-500 dark:text-gray-400">Admin Panel</p>
//             </div>
//           </div>
//         </div>

//         {/* Navigation */}
//         <nav className="p-4 space-y-1">
//           {navigation.map((item) => {
//             const Icon = item.icon;
//             return (
//               <NavLink
//                 key={item.name}
//                 to={item.href}
//                 onClick={() => setIsOpen(false)}
//                 className={({ isActive }) =>
//                   `flex items-center space-x-3 px-4 py-3 rounded-lg transition-colors ${
//                     isActive
//                       ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-600 dark:text-primary-400'
//                       : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
//                   }`
//                 }
//               >
//                 <Icon className="h-5 w-5" />
//                 <span className="font-medium">{item.name}</span>
//               </NavLink>
//             );
//           })}
//         </nav>

//         {/* Logout Button */}
//         <div className="absolute bottom-0 left-0 right-0 p-4 border-t border-gray-200 dark:border-gray-700">
//           <button
//             onClick={logout}
//             className="flex items-center space-x-3 w-full px-4 py-3 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors"
//           >
//             <LogOut className="h-5 w-5" />
//             <span className="font-medium">Logout</span>
//           </button>
//         </div>
//       </aside>
//     </>
//   );
// };

// export default Sidebar;