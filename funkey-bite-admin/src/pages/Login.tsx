import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Eye, EyeOff, Loader2, Lock, Mail } from 'lucide-react';
import { useAuthStore } from '../stores/authStore';
import { adminLogin } from '../api/adminApi';
import BrandLogo from '../components/ui/BrandLogo';

const loginSchema = z.object({
  email: z.string().email('Enter your email address'),
  password: z.string().min(1, 'Enter your password'),
});

type LoginFormData = z.infer<typeof loginSchema>;

const Login: React.FC = () => {
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [formError, setFormError] = useState('');
  const navigate = useNavigate();
  const { login } = useAuthStore();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>({ resolver: zodResolver(loginSchema) });

  const onSubmit = async (data: LoginFormData) => {
    setIsLoading(true);
    setFormError('');
    try {
      const response = await adminLogin(data);
      localStorage.setItem('admin_token', response.token);
      login(response.token, response.admin);
      navigate('/dashboard', { replace: true });
    } catch (error: any) {
      setFormError(
        error?.status === 429
          ? 'Too many attempts. Wait a few minutes and try again.'
          : error?.message || 'That email and password don’t match an admin account.',
      );
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex min-h-dvh items-center justify-center bg-gradient-to-b from-primary-50 to-white px-4 py-10 dark:from-gray-900 dark:to-gray-950">
      <div className="w-full max-w-sm">
        <div className="mb-6 flex flex-col items-center text-center">
          <BrandLogo className="h-14" />
          <h1 className="mt-4 text-2xl font-bold text-gray-900 dark:text-white">Admin sign in</h1>
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">Orders, menu, catering and payments</p>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className="card space-y-4 p-6" noValidate>
          <div>
            <label htmlFor="login-email" className="form-label">Email</label>
            <div className="relative">
              <Mail className="pointer-events-none absolute left-3.5 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400" aria-hidden="true" />
              <input
                id="login-email"
                {...register('email')}
                type="email"
                autoComplete="username"
                inputMode="email"
                className="input-primary pl-11"
                placeholder="you@example.com"
                aria-invalid={errors.email ? 'true' : undefined}
                aria-describedby={errors.email ? 'login-email-error' : undefined}
              />
            </div>
            {errors.email && <p id="login-email-error" className="mt-1.5 text-sm text-red-600">{errors.email.message}</p>}
          </div>

          <div>
            <label htmlFor="login-password" className="form-label">Password</label>
            <div className="relative">
              <Lock className="pointer-events-none absolute left-3.5 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400" aria-hidden="true" />
              <input
                id="login-password"
                {...register('password')}
                type={showPassword ? 'text' : 'password'}
                autoComplete="current-password"
                className="input-primary pl-11 pr-12"
                aria-invalid={errors.password ? 'true' : undefined}
                aria-describedby={errors.password ? 'login-password-error' : undefined}
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="icon-btn absolute right-0.5 top-1/2 -translate-y-1/2"
                aria-label={showPassword ? 'Hide password' : 'Show password'}
                aria-pressed={showPassword}
              >
                {showPassword ? <EyeOff className="h-5 w-5" aria-hidden="true" /> : <Eye className="h-5 w-5" aria-hidden="true" />}
              </button>
            </div>
            {errors.password && <p id="login-password-error" className="mt-1.5 text-sm text-red-600">{errors.password.message}</p>}
          </div>

          {formError && <p className="rounded-xl bg-red-50 px-3 py-2 text-sm text-red-700 dark:bg-red-500/10 dark:text-red-300" role="alert">{formError}</p>}

          <button type="submit" disabled={isLoading} className="btn-primary w-full">
            {isLoading && <Loader2 className="h-5 w-5 animate-spin" aria-hidden="true" />}
            {isLoading ? 'Signing in…' : 'Sign in'}
          </button>
        </form>
      </div>
    </div>
  );
};

export default Login;
