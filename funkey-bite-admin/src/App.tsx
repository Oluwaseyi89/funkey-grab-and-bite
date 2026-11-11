import { BrowserRouter as Router } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Toaster } from 'react-hot-toast';
import { SocketProvider } from './contexts/SocketContext';
import { ThemeProvider } from './contexts/ThemeContext';
import { AuthProvider } from './contexts/AuthContext';
import AppRoutes from './routes/AppRoutes';
import ErrorBoundary from './components/common/ErrorBoundary';
import { ConfirmProvider } from './components/ui/ConfirmProvider';
import './styles/globals.css';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
      staleTime: 5 * 60 * 1000, // 5 minutes
    },
  },
});

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Router>
        <ThemeProvider>
          <AuthProvider>
            <SocketProvider>
              <ConfirmProvider>
                <ErrorBoundary fullScreen>
                  <AppRoutes />
                </ErrorBoundary>
              </ConfirmProvider>
              {/* Top-centre on phones, where it's seen; top-right on larger screens. */}
              <Toaster
                position="top-center"
                containerStyle={{ top: 64 }}
                toastOptions={{
                  className: '!rounded-2xl !shadow-lift dark:!bg-gray-800 dark:!text-white',
                  duration: 4000,
                }}
              />
            </SocketProvider>
          </AuthProvider>
        </ThemeProvider>
      </Router>
    </QueryClientProvider>
  );
}

export default App;