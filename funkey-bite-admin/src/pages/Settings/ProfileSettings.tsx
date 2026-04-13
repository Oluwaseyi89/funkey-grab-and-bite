import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import {
  Save,
  User,
  Mail,
  Shield,
  Key,
  Eye,
  EyeOff,
  CheckCircle,
  AlertCircle,
  LogOut,
  Calendar,
  Clock,
  UserCog
} from 'lucide-react';
import toast from 'react-hot-toast';
import { useAuthStore } from '../../stores/authStore';
import { updateAdminPassword, adminLogout } from '../../api/adminApi';

const profileSchema = z.object({
  email: z.string().email('Valid email is required'),
});

const passwordSchema = z.object({
  currentPassword: z.string().min(6, 'Current password is required'),
  newPassword: z.string()
    .min(8, 'Password must be at least 8 characters')
    .regex(/[A-Z]/, 'Password must contain at least one uppercase letter')
    .regex(/[a-z]/, 'Password must contain at least one lowercase letter')
    .regex(/[0-9]/, 'Password must contain at least one number')
    .regex(/[^A-Za-z0-9]/, 'Password must contain at least one special character'),
  confirmPassword: z.string().min(8, 'Please confirm your password'),
}).refine((data) => data.newPassword === data.confirmPassword, {
  message: "Passwords don't match",
  path: ["confirmPassword"],
});

type ProfileFormData = z.infer<typeof profileSchema>;
type PasswordFormData = z.infer<typeof passwordSchema>;

const ProfileSettings: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'profile' | 'security'>('profile');
  const [showCurrentPassword, setShowCurrentPassword] = useState(false);
  const [showNewPassword, setShowNewPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [isChangingPassword, setIsChangingPassword] = useState(false);
  const { user, logout } = useAuthStore();

  const {
    register: registerProfile,
    handleSubmit: handleSubmitProfile,
    formState: { errors: profileErrors, isSubmitting: isProfileSubmitting },
  } = useForm<ProfileFormData>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      email: user?.email || '',
    },
  });

  const {
    register: registerPassword,
    handleSubmit: handleSubmitPassword,
    formState: { errors: passwordErrors, isSubmitting: isPasswordSubmitting },
    reset: resetPasswordForm,
  } = useForm<PasswordFormData>({
    resolver: zodResolver(passwordSchema),
  });

  const handleProfileUpdate = async (data: ProfileFormData) => {
    toast.error('Profile update endpoint not implemented');
  };

  const handlePasswordChange = async (data: PasswordFormData) => {
    setIsChangingPassword(true);
    try {
      await updateAdminPassword({
        currentPassword: data.currentPassword,
        newPassword: data.newPassword,
      });
      
      toast.success('Password updated successfully!');
      resetPasswordForm();
    } catch (error: any) {
      toast.error(error.message || 'Failed to update password');
    } finally {
      setIsChangingPassword(false);
    }
  };

  const handleLogout = async () => {
    try {
      await adminLogout();
      logout();
      window.location.href = '/login';
      toast.success('Logged out successfully');
    } catch (error) {
      logout();
      window.location.href = '/login';
    }
  };

  const formatDate = (dateString?: string) => {
    if (!dateString) return 'Never';
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  if (!user) {
    return (
      <div className="text-center py-12">
        <AlertCircle className="h-12 w-12 text-gray-400 mx-auto mb-4" />
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
          No User Data
        </h2>
        <p className="text-gray-500 dark:text-gray-400">
          Please login to view profile settings
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Profile Settings</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Manage your account information and security
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        
        <div className="lg:col-span-1 space-y-6">
          
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <div className="flex flex-col items-center text-center">
              <div className="h-24 w-24 rounded-full bg-gradient-to-br from-primary-500 to-primary-600 flex items-center justify-center mb-4">
                <User className="h-12 w-12 text-white" />
              </div>
              
              <h2 className="text-xl font-bold text-gray-900 dark:text-white">
                {user.username}
              </h2>
              
              <div className="mt-2">
                <span className={`px-3 py-1 rounded-full text-sm font-medium ${
                  user.role === 'admin' 
                    ? 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200'
                    : user.role === 'manager'
                    ? 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200'
                    : 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                }`}>
                  {user.role.charAt(0).toUpperCase() + user.role.slice(1)}
                </span>
              </div>
              
              <p className="text-gray-500 dark:text-gray-400 mt-2">{user.email}</p>
              
              <div className="mt-6 space-y-3 w-full">
                <div className="flex items-center justify-between text-sm">
                  <span className="text-gray-500 dark:text-gray-400">Status:</span>
                  <span className={`flex items-center space-x-1 ${
                    user.isActive 
                      ? 'text-green-600 dark:text-green-400' 
                      : 'text-red-600 dark:text-red-400'
                  }`}>
                    {user.isActive ? (
                      <>
                        <CheckCircle className="h-4 w-4" />
                        <span>Active</span>
                      </>
                    ) : (
                      <span>Inactive</span>
                    )}
                  </span>
                </div>
                
                <div className="flex items-center justify-between text-sm">
                  <span className="text-gray-500 dark:text-gray-400">Member since:</span>
                  <span className="text-gray-900 dark:text-white">
                    {formatDate(user.createdAt).split(',')[0]}
                  </span>
                </div>
                
                {user.lastLogin && (
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-gray-500 dark:text-gray-400">Last login:</span>
                    <span className="text-gray-900 dark:text-white">
                      {formatDate(user.lastLogin)}
                    </span>
                  </div>
                )}
              </div>
            </div>
          </div>

          
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <h3 className="font-medium text-gray-900 dark:text-white mb-4 flex items-center space-x-2">
              <UserCog className="h-5 w-5 text-primary-500" />
              <span>Account Information</span>
            </h3>
            
            <div className="space-y-3">
              <div className="flex items-center justify-between text-sm">
                <span className="text-gray-500 dark:text-gray-400">User ID:</span>
                <span className="font-mono text-gray-900 dark:text-white">{user.id}</span>
              </div>
              
              <div className="flex items-center justify-between text-sm">
                <span className="text-gray-500 dark:text-gray-400">Role:</span>
                <span className="text-gray-900 dark:text-white capitalize">{user.role}</span>
              </div>
              
              <div className="flex items-center justify-between text-sm">
                <span className="text-gray-500 dark:text-gray-400">Email:</span>
                <span className="text-gray-900 dark:text-white">{user.email}</span>
              </div>
              
              <div className="flex items-center justify-between text-sm">
                <span className="text-gray-500 dark:text-gray-400">Account created:</span>
                <span className="text-gray-900 dark:text-white">{formatDate(user.createdAt)}</span>
              </div>
            </div>
          </div>

          
          <button
            onClick={handleLogout}
            className="w-full flex items-center justify-center space-x-2 px-4 py-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 hover:bg-red-100 dark:hover:bg-red-900/30 rounded-xl transition-colors"
          >
            <LogOut className="h-5 w-5" />
            <span className="font-medium">Logout</span>
          </button>
        </div>

        
        <div className="lg:col-span-2 space-y-6">
          
          <div className="bg-white dark:bg-gray-800 rounded-xl p-2 border border-gray-200 dark:border-gray-700">
            <div className="flex space-x-2">
              <button
                onClick={() => setActiveTab('profile')}
                className={`flex-1 px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                  activeTab === 'profile'
                    ? 'bg-primary-500 text-white'
                    : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
                }`}
              >
                <div className="flex items-center justify-center space-x-2">
                  <User className="h-4 w-4" />
                  <span>Profile Information</span>
                </div>
              </button>
              
              <button
                onClick={() => setActiveTab('security')}
                className={`flex-1 px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                  activeTab === 'security'
                    ? 'bg-primary-500 text-white'
                    : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
                }`}
              >
                <div className="flex items-center justify-center space-x-2">
                  <Shield className="h-4 w-4" />
                  <span>Security & Password</span>
                </div>
              </button>
            </div>
          </div>

          
          {activeTab === 'profile' && (
            <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
              <div className="flex items-center space-x-3 mb-6">
                <User className="h-5 w-5 text-primary-500" />
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                  Profile Information
                </h2>
              </div>
              
              <form onSubmit={handleSubmitProfile(handleProfileUpdate)} className="space-y-6">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Username
                    </label>
                    <input
                      type="text"
                      value={user.username}
                      disabled
                      className="w-full px-4 py-2 bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg dark:text-white cursor-not-allowed"
                    />
                    <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                      Username cannot be changed
                    </p>
                  </div>
                  
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Email Address *
                    </label>
                    <div className="relative">
                      <Mail className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                      <input
                        type="email"
                        {...registerProfile('email')}
                        className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                        placeholder="your.email@example.com"
                      />
                    </div>
                    {profileErrors.email && (
                      <p className="mt-1 text-sm text-red-600">{profileErrors.email.message}</p>
                    )}
                  </div>
                </div>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Role
                    </label>
                    <input
                      type="text"
                      value={user.role.charAt(0).toUpperCase() + user.role.slice(1)}
                      disabled
                      className="w-full px-4 py-2 bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg dark:text-white cursor-not-allowed"
                    />
                  </div>
                  
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Account Status
                    </label>
                    <div className="flex items-center space-x-2">
                      {user.isActive ? (
                        <CheckCircle className="h-5 w-5 text-green-500" />
                      ) : (
                        <AlertCircle className="h-5 w-5 text-red-500" />
                      )}
                      <span className={user.isActive ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'}>
                        {user.isActive ? 'Active' : 'Inactive'}
                      </span>
                    </div>
                  </div>
                </div>
                
                <div className="pt-6 border-t border-gray-200 dark:border-gray-700">
                  <button
                    type="submit"
                    disabled={isProfileSubmitting}
                    className="flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors disabled:opacity-50"
                  >
                    <Save className="h-4 w-4" />
                    <span>{isProfileSubmitting ? 'Saving...' : 'Update Profile'}</span>
                  </button>
                  <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">
                    Note: Profile updates are limited to email only. Contact super-admin for other changes.
                  </p>
                </div>
              </form>
            </div>
          )}

          
          {activeTab === 'security' && (
            <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
              <div className="flex items-center space-x-3 mb-6">
                <Shield className="h-5 w-5 text-primary-500" />
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                  Security & Password
                </h2>
              </div>
              
              <form onSubmit={handleSubmitPassword(handlePasswordChange)} className="space-y-6">
                
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Current Password *
                  </label>
                  <div className="relative">
                    <Key className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                    <input
                      type={showCurrentPassword ? 'text' : 'password'}
                      {...registerPassword('currentPassword')}
                      className="w-full pl-10 pr-10 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                      placeholder="••••••••"
                    />
                    <button
                      type="button"
                      onClick={() => setShowCurrentPassword(!showCurrentPassword)}
                      className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
                    >
                      {showCurrentPassword ? (
                        <EyeOff className="h-5 w-5" />
                      ) : (
                        <Eye className="h-5 w-5" />
                      )}
                    </button>
                  </div>
                  {passwordErrors.currentPassword && (
                    <p className="mt-1 text-sm text-red-600">{passwordErrors.currentPassword.message}</p>
                  )}
                </div>
                
                
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    New Password *
                  </label>
                  <div className="relative">
                    <Key className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                    <input
                      type={showNewPassword ? 'text' : 'password'}
                      {...registerPassword('newPassword')}
                      className="w-full pl-10 pr-10 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                      placeholder="••••••••"
                    />
                    <button
                      type="button"
                      onClick={() => setShowNewPassword(!showNewPassword)}
                      className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
                    >
                      {showNewPassword ? (
                        <EyeOff className="h-5 w-5" />
                      ) : (
                        <Eye className="h-5 w-5" />
                      )}
                    </button>
                  </div>
                  {passwordErrors.newPassword && (
                    <p className="mt-1 text-sm text-red-600">{passwordErrors.newPassword.message}</p>
                  )}
                </div>
                
                
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Confirm New Password *
                  </label>
                  <div className="relative">
                    <Key className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                    <input
                      type={showConfirmPassword ? 'text' : 'password'}
                      {...registerPassword('confirmPassword')}
                      className="w-full pl-10 pr-10 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                      placeholder="••••••••"
                    />
                    <button
                      type="button"
                      onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                      className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
                    >
                      {showConfirmPassword ? (
                        <EyeOff className="h-5 w-5" />
                      ) : (
                        <Eye className="h-5 w-5" />
                      )}
                    </button>
                  </div>
                  {passwordErrors.confirmPassword && (
                    <p className="mt-1 text-sm text-red-600">{passwordErrors.confirmPassword.message}</p>
                  )}
                </div>
                
                
                <div className="bg-gray-50 dark:bg-gray-700/50 rounded-lg p-4">
                  <h4 className="text-sm font-medium text-gray-900 dark:text-white mb-2">
                    Password Requirements:
                  </h4>
                  <ul className="text-sm text-gray-600 dark:text-gray-400 space-y-1">
                    <li>• At least 8 characters long</li>
                    <li>• At least one uppercase letter (A-Z)</li>
                    <li>• At least one lowercase letter (a-z)</li>
                    <li>• At least one number (0-9)</li>
                    <li>• At least one special character (!@#$%^&* etc.)</li>
                  </ul>
                </div>
                
                
                <div className="flex items-center space-x-3 pt-6 border-t border-gray-200 dark:border-gray-700">
                  <button
                    type="submit"
                    disabled={isPasswordSubmitting || isChangingPassword}
                    className="flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors disabled:opacity-50"
                  >
                    <Save className="h-4 w-4" />
                    <span>{(isPasswordSubmitting || isChangingPassword) ? 'Updating...' : 'Change Password'}</span>
                  </button>
                  
                  <button
                    type="button"
                    onClick={() => resetPasswordForm()}
                    className="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                  >
                    Clear
                  </button>
                </div>
              </form>
            </div>
          )}

          
          <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-xl p-6">
            <div className="flex items-start space-x-3">
              <Shield className="h-5 w-5 text-blue-500 mt-0.5 flex-shrink-0" />
              <div className="space-y-3">
                <h3 className="font-medium text-blue-800 dark:text-blue-300">
                  Security Tips
                </h3>
                <ul className="text-sm text-blue-700 dark:text-blue-400 space-y-2">
                  <li>• Use a strong, unique password that you don't use elsewhere</li>
                  <li>• Change your password regularly (every 90 days)</li>
                  <li>• Never share your password with anyone</li>
                  <li>• Log out when using shared computers</li>
                  <li>• Contact support immediately if you suspect unauthorized access</li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProfileSettings;