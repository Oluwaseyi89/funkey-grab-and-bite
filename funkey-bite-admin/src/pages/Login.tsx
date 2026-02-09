import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Lock, Mail, Eye, EyeOff } from 'lucide-react';
import toast from 'react-hot-toast';
import { useAuthStore } from '../stores/authStore';
import { adminLogin } from '../api/adminApi';
import type { AdminLoginCredentials } from '../types';

const loginSchema = z.object({
  email: z.string().email('Invalid email address'),
  password: z.string().min(6, 'Password must be at least 6 characters'),
});

type LoginFormData = z.infer<typeof loginSchema>;

const Login: React.FC = () => {
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();
  const { login } = useAuthStore();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
  });

  const onSubmit = async (data: LoginFormData) => {
    setIsLoading(true);
    try {
      const response = await adminLogin(data);
      
      // Store token and user
      localStorage.setItem('admin_token', response.token);
      login(response.token, response.admin);
      
      toast.success('Login successful!');
      navigate('/dashboard');
    } catch (error: any) {
      toast.error(error.message || 'Login failed. Please check your credentials.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary-50 to-white dark:from-gray-900 dark:to-gray-800 flex items-center justify-center p-4">
      <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-xl w-full max-w-md overflow-hidden">
        {/* Header */}
        <div className="bg-primary-500 p-8 text-center">
          <div className="flex justify-center mb-4">
            <div className="bg-white p-3 rounded-full">
              <span className="text-2xl font-bold text-primary-500">FG&B</span>
            </div>
          </div>
          <h1 className="text-2xl font-bold text-white">Admin Dashboard</h1>
          <p className="text-primary-100 mt-2">Sign in to manage your business</p>
        </div>

        {/* Form */}
        <div className="p-8">
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            {/* Email Field */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Email Address
              </label>
              <div className="relative">
                <Mail className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                <input
                  {...register('email')}
                  type="email"
                  className="pl-10 w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                  placeholder="admin@funkeygrabandbite.com"
                />
              </div>
              {errors.email && (
                <p className="mt-1 text-sm text-red-600">{errors.email.message}</p>
              )}
            </div>

            {/* Password Field */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Password
              </label>
              <div className="relative">
                <Lock className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                <input
                  {...register('password')}
                  type={showPassword ? 'text' : 'password'}
                  className="pl-10 pr-10 w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                  placeholder="••••••••"
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
                >
                  {showPassword ? (
                    <EyeOff className="h-5 w-5" />
                  ) : (
                    <Eye className="h-5 w-5" />
                  )}
                </button>
              </div>
              {errors.password && (
                <p className="mt-1 text-sm text-red-600">{errors.password.message}</p>
              )}
            </div>

            {/* Submit Button */}
            <button
              type="submit"
              disabled={isLoading}
              className="w-full bg-primary-500 hover:bg-primary-600 text-white font-semibold py-3 px-4 rounded-lg transition duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center"
            >
              {isLoading ? (
                <>
                  <div className="animate-spin rounded-full h-5 w-5 border-t-2 border-b-2 border-white mr-2"></div>
                  Signing in...
                </>
              ) : (
                'Sign In'
              )}
            </button>
          </form>

          {/* Demo Credentials */}
          <div className="mt-8 p-4 bg-gray-50 dark:bg-gray-700 rounded-lg">
            <p className="text-sm text-gray-600 dark:text-gray-400">
              <strong>Default Admin:</strong><br />
              Email: admin@funkey.com<br />
              Password: admin123
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Login;





















// // src/pages/Login.tsx
// import React, { useState } from 'react';
// import { useNavigate } from 'react-router-dom';
// import { useForm } from 'react-hook-form';
// import { zodResolver } from '@hookform/resolvers/zod';
// import { z } from 'zod';
// import { Lock, Mail, Eye, EyeOff } from 'lucide-react';
// import toast from 'react-hot-toast';
// import { useAuthStore } from '../stores/authStore';
// import api from '../api/axiosConfig';
// import logo from '../assets/logo.svg';

// const loginSchema = z.object({
//   email: z.string().email('Invalid email address'),
//   password: z.string().min(6, 'Password must be at least 6 characters'),
// });

// type LoginFormData = z.infer<typeof loginSchema>;

// const Login: React.FC = () => {
//   const [showPassword, setShowPassword] = useState(false);
//   const [isLoading, setIsLoading] = useState(false);
//   const navigate = useNavigate();
//   const { login } = useAuthStore();

//   const {
//     register,
//     handleSubmit,
//     formState: { errors },
//   } = useForm<LoginFormData>({
//     resolver: zodResolver(loginSchema),
//   });

//   const onSubmit = async (data: LoginFormData) => {
//     setIsLoading(true);
//     try {
//       // For now, using admin login endpoint. We'll need to add this to your Go API
//       const response = await api.post('/api/v1/admin/auth/login', data);
      
//       const { token, user } = response.data;
      
//       login(token, user);
//       toast.success('Login successful!');
//       navigate('/dashboard');
//     } catch (error: any) {
//       toast.error(error.response?.data?.message || 'Login failed');
//     } finally {
//       setIsLoading(false);
//     }
//   };

//   return (
//     <div className="min-h-screen bg-gradient-to-br from-primary-50 to-white dark:from-gray-900 dark:to-gray-800 flex items-center justify-center p-4">
//       <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-xl w-full max-w-md overflow-hidden">
//         {/* Header */}
//         <div className="bg-primary-500 p-8 text-center">
//           <div className="flex justify-center mb-4">
//             <div className="bg-white p-3 rounded-full">
//               <img src={logo} alt="Funkey Grab & Bite" className="h-12" />
//             </div>
//           </div>
//           <h1 className="text-2xl font-bold text-white">Admin Dashboard</h1>
//           <p className="text-primary-100 mt-2">Sign in to manage your business</p>
//         </div>

//         {/* Form */}
//         <div className="p-8">
//           <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
//             {/* Email Field */}
//             <div>
//               <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
//                 Email Address
//               </label>
//               <div className="relative">
//                 <Mail className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
//                 <input
//                   {...register('email')}
//                   type="email"
//                   className="pl-10 w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
//                   placeholder="admin@funkeygrabandbite.com"
//                 />
//               </div>
//               {errors.email && (
//                 <p className="mt-1 text-sm text-red-600">{errors.email.message}</p>
//               )}
//             </div>

//             {/* Password Field */}
//             <div>
//               <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
//                 Password
//               </label>
//               <div className="relative">
//                 <Lock className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
//                 <input
//                   {...register('password')}
//                   type={showPassword ? 'text' : 'password'}
//                   className="pl-10 pr-10 w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
//                   placeholder="••••••••"
//                 />
//                 <button
//                   type="button"
//                   onClick={() => setShowPassword(!showPassword)}
//                   className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
//                 >
//                   {showPassword ? (
//                     <EyeOff className="h-5 w-5" />
//                   ) : (
//                     <Eye className="h-5 w-5" />
//                   )}
//                 </button>
//               </div>
//               {errors.password && (
//                 <p className="mt-1 text-sm text-red-600">{errors.password.message}</p>
//               )}
//             </div>

//             {/* Submit Button */}
//             <button
//               type="submit"
//               disabled={isLoading}
//               className="w-full bg-primary-500 hover:bg-primary-600 text-white font-semibold py-3 px-4 rounded-lg transition duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center"
//             >
//               {isLoading ? (
//                 <>
//                   <div className="animate-spin rounded-full h-5 w-5 border-t-2 border-b-2 border-white mr-2"></div>
//                   Signing in...
//                 </>
//               ) : (
//                 'Sign In'
//               )}
//             </button>

//             {/* Forgot Password Link */}
//             <div className="text-center">
//               <a
//                 href="#"
//                 className="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300"
//                 onClick={(e) => {
//                   e.preventDefault();
//                   toast('Password reset feature coming soon!');
//                 }}
//               >
//                 Forgot your password?
//               </a>
//             </div>
//           </form>

//           {/* Demo Credentials (Remove in production) */}
//           <div className="mt-8 p-4 bg-gray-50 dark:bg-gray-700 rounded-lg">
//             <p className="text-sm text-gray-600 dark:text-gray-400">
//               <strong>Demo Credentials:</strong><br />
//               Email: admin@funkey.com<br />
//               Password: admin123
//             </p>
//           </div>
//         </div>
//       </div>
//     </div>
//   );
// };

// export default Login;