import React, { createContext, useContext, useEffect, useMemo, useRef, useState } from 'react';
import { useAuthStore } from '../stores/authStore';

type Handler = (payload: any) => void;

type RealtimeSocket = {
  on: (event: string, handler: Handler) => void;
  off: (event: string, handler: Handler) => void;
};

interface SocketContextType {
  socket: RealtimeSocket | null;
  isConnected: boolean;
  isRealtimeEnabled: boolean;
}

const SocketContext = createContext<SocketContextType>({
  socket: null,
  isConnected: false,
  isRealtimeEnabled: false,
});

export const useSocket = () => useContext(SocketContext);

const BASE_DELAY_MS = 1000;
const MAX_DELAY_MS = 30000;

const buildRealtimeUrl = () => {
  const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';
  const url = new URL(apiUrl);
  url.protocol = url.protocol === 'https:' ? 'wss:' : 'ws:';
  url.pathname = '/api/v1/admin/realtime/ws';
  return url.toString();
};

// The live connection for admin events (new orders, payments, catering).
// It opens as soon as someone signs in, reconnects with backoff when it
// drops, and closes on sign-out.
export const SocketProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const realtimeMode = (import.meta.env.VITE_REALTIME_MODE || 'websocket').toLowerCase();
  const isRealtimeEnabled = realtimeMode === 'websocket' || realtimeMode === 'socket-io';
  const token = useAuthStore((state) => state.token);
  const [isConnected, setIsConnected] = useState(false);
  const listeners = useRef<Map<string, Set<Handler>>>(new Map());

  const socket = useMemo<RealtimeSocket>(
    () => ({
      on: (event, handler) => {
        const set = listeners.current.get(event) ?? new Set();
        set.add(handler);
        listeners.current.set(event, set);
      },
      off: (event, handler) => {
        listeners.current.get(event)?.delete(handler);
      },
    }),
    [],
  );

  useEffect(() => {
    const authToken = token || localStorage.getItem('admin_token');
    if (!isRealtimeEnabled || !authToken) {
      setIsConnected(false);
      return;
    }

    // Everything below belongs to this run of the effect only. A later run
    // (new token, or React mounting twice in development) starts its own, and
    // events from an old connection are ignored.
    let active = true;
    let ws: WebSocket | null = null;
    let attempts = 0;
    let timer: ReturnType<typeof setTimeout> | undefined;

    const connect = () => {
      if (!active) return;
      ws = new WebSocket(buildRealtimeUrl(), [`jwt.${authToken}`]);
      const current = ws;

      current.onopen = () => {
        if (!active) return;
        attempts = 0;
        setIsConnected(true);
      };
      current.onmessage = (event) => {
        if (!active) return;
        try {
          const message = JSON.parse(event.data);
          if (typeof message?.event !== 'string') return;
          listeners.current.get(message.event)?.forEach((handler) => handler(message.data));
        } catch {
          // Not an event we understand.
        }
      };
      current.onclose = () => {
        if (!active) return;
        setIsConnected(false);
        attempts += 1;
        const delay = Math.min(BASE_DELAY_MS * 2 ** (attempts - 1), MAX_DELAY_MS) + Math.floor(Math.random() * 500);
        timer = setTimeout(connect, delay);
      };
      current.onerror = () => current.close();
    };

    // Reconnect straight away when the phone wakes up or comes back online.
    const reconnectNow = () => {
      if (!active || document.visibilityState !== 'visible') return;
      if (ws && (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING)) return;
      clearTimeout(timer);
      attempts = 0;
      connect();
    };

    connect();
    document.addEventListener('visibilitychange', reconnectNow);
    window.addEventListener('online', reconnectNow);

    return () => {
      active = false;
      clearTimeout(timer);
      document.removeEventListener('visibilitychange', reconnectNow);
      window.removeEventListener('online', reconnectNow);
      ws?.close();
      setIsConnected(false);
    };
  }, [isRealtimeEnabled, token]);

  return (
    <SocketContext.Provider value={{ socket: isRealtimeEnabled ? socket : null, isConnected, isRealtimeEnabled }}>
      {children}
    </SocketContext.Provider>
  );
};
