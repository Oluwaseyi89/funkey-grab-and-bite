import React from 'react';
import { Outlet } from 'react-router-dom';
import Sidebar from './Sidebar';
import Topbar from './Topbar';
import { useSocket } from '../../contexts/SocketContext';
import { Bell, Wifi, WifiOff } from 'lucide-react';

const Layout: React.FC = () => {
  const { isConnected } = useSocket();

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      {/* Topbar */}
      <Topbar />

      <div className="flex pt-16">
        {/* Sidebar */}
        <Sidebar />

        {/* Main Content */}
        <main className="flex-1 p-4 md:p-6 lg:p-8 overflow-x-hidden">
          {/* Status Bar */}
          <div className="mb-6 flex flex-wrap items-center justify-between gap-4">
            {/* Connection Status */}
            <div className="flex items-center space-x-3">
              <div className="flex items-center space-x-2">
                {isConnected ? (
                  <Wifi className="h-4 w-4 text-green-500" />
                ) : (
                  <WifiOff className="h-4 w-4 text-red-500" />
                )}
                <span className="text-sm text-gray-600 dark:text-gray-400">
                  {isConnected ? 'Connected' : 'Disconnected'}
                </span>
              </div>
              
              {/* Notifications Badge */}
              <button className="relative p-2 text-gray-600 dark:text-gray-400 hover:text-primary-500 transition-colors">
                <Bell className="h-5 w-5" />
                <span className="absolute -top-1 -right-1 h-5 w-5 bg-red-500 text-white text-xs rounded-full flex items-center justify-center">
                  3
                </span>
              </button>
            </div>

            {/* Breadcrumb placeholder */}
            <div className="text-sm text-gray-500 dark:text-gray-400">
              Dashboard / Overview
            </div>
          </div>

          {/* Page Content */}
          <div className="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 p-4 md:p-6">
            <Outlet />
          </div>

          {/* Footer */}
          <footer className="mt-8 text-center text-sm text-gray-500 dark:text-gray-400">
            <p>Funkey Grab & Bite Admin Dashboard © {new Date().getFullYear()}</p>
            <p className="mt-1">Version 1.0.0</p>
          </footer>
        </main>
      </div>
    </div>
  );
};

export default Layout;























// // src/components/layout/Layout.tsx
// import React from 'react';
// import { Outlet } from 'react-router-dom';
// import Sidebar from './Sidebar';
// import Topbar from './Topbar';
// import { useSocket } from '../../contexts/SocketContext';
// import { Bell } from 'lucide-react';

// const Layout: React.FC = () => {
//   const { socket, isConnected } = useSocket();

//   return (
//     <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
//       {/* Topbar */}
//       <Topbar />

//       <div className="flex">
//         {/* Sidebar */}
//         <Sidebar />

//         {/* Main Content */}
//         <main className="flex-1 p-6 lg:p-8">
//           {/* Connection Status */}
//           {socket && (
//             <div className="mb-6 flex items-center justify-between">
//               <div className="flex items-center space-x-2">
//                 <div
//                   className={`h-3 w-3 rounded-full ${
//                     isConnected ? 'bg-green-500 animate-pulse' : 'bg-red-500'
//                   }`}
//                 />
//                 <span className="text-sm text-gray-600 dark:text-gray-400">
//                   {isConnected ? 'Connected to server' : 'Disconnected'}
//                 </span>
//               </div>
              
//               {/* Notifications */}
//               <button className="p-2 text-gray-600 dark:text-gray-400 hover:text-primary-500">
//                 <Bell className="h-5 w-5" />
//               </button>
//             </div>
//           )}

//           {/* Page Content */}
//           <Outlet />
//         </main>
//       </div>
//     </div>
//   );
// };

// export default Layout;