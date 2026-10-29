import React from 'react';
import { AlertTriangle, RefreshCw } from 'lucide-react';

interface ErrorBoundaryProps {
  children: React.ReactNode;
  /** Clears the error when it changes, e.g. the current path, so navigating away recovers. */
  resetKey?: string;
  /** Fill the screen, for errors outside the admin layout. */
  fullScreen?: boolean;
}

interface ErrorBoundaryState {
  error: Error | null;
}

// Shows a message instead of a blank screen when a page fails to render.
class ErrorBoundary extends React.Component<ErrorBoundaryProps, ErrorBoundaryState> {
  state: ErrorBoundaryState = { error: null };

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { error };
  }

  componentDidCatch(error: Error, info: React.ErrorInfo) {
    console.error('Page failed to render', error, info.componentStack);
  }

  componentDidUpdate(prevProps: ErrorBoundaryProps) {
    if (this.state.error && prevProps.resetKey !== this.props.resetKey) {
      this.setState({ error: null });
    }
  }

  render() {
    const { error } = this.state;
    if (!error) {
      return this.props.children;
    }

    return (
      <div
        role="alert"
        className={`flex items-center justify-center p-6 ${
          this.props.fullScreen ? 'min-h-screen bg-gray-50 dark:bg-gray-900' : 'min-h-[50vh]'
        }`}
      >
        <div className="max-w-md text-center space-y-4">
          <AlertTriangle className="h-12 w-12 text-amber-500 mx-auto" />
          <h1 className="text-xl font-semibold text-gray-900 dark:text-white">
            Something went wrong on this page
          </h1>
          <p className="text-gray-600 dark:text-gray-400">
            Try again, or reload if the problem continues. Other pages are not affected.
          </p>
          <details className="text-left text-sm text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-800 rounded-lg p-3">
            <summary className="cursor-pointer">Error details</summary>
            <p className="mt-2 font-mono break-words">{error.message}</p>
          </details>
          <div className="flex items-center justify-center gap-3">
            <button
              onClick={() => this.setState({ error: null })}
              className="flex items-center gap-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors"
            >
              <RefreshCw className="h-4 w-4" />
              Try again
            </button>
            <button
              onClick={() => window.location.reload()}
              className="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
            >
              Reload page
            </button>
            <a
              href="/dashboard"
              className="px-4 py-2 text-primary-600 dark:text-primary-400 hover:underline"
            >
              Dashboard
            </a>
          </div>
        </div>
      </div>
    );
  }
}

export default ErrorBoundary;
