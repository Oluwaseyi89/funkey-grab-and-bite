import React, { createContext, useContext, useEffect, useMemo, useRef, useState } from 'react';

type RealtimeSocket = {
  on: (event: string, handler: (payload: any) => void) => void;
  off: (event: string, handler: (payload: any) => void) => void;
  disconnect: () => void;
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

const WEBSOCKET_RECONNECT_DELAY_MS = 3000;

const buildRealtimeUrl = () => {
  const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';
  const baseUrl = new URL(apiUrl);
  baseUrl.protocol = baseUrl.protocol === 'https:' ? 'wss:' : 'ws:';
  baseUrl.pathname = '/api/v1/admin/realtime/ws';
  return baseUrl;
};

export const SocketProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const realtimeMode = (import.meta.env.VITE_REALTIME_MODE || 'websocket').toLowerCase();
  const isRealtimeEnabled = realtimeMode === 'websocket' || realtimeMode === 'socket-io';
  const [isConnected, setIsConnected] = useState(false);
  const [socket, setSocket] = useState<RealtimeSocket | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const shouldReconnectRef = useRef(true);
  const listenersRef = useRef<Map<string, Set<(payload: any) => void>>>(new Map());

  const socketApi = useMemo<RealtimeSocket>(
    () => ({
      on: (event, handler) => {
        const current = listenersRef.current.get(event) || new Set();
        current.add(handler);
        listenersRef.current.set(event, current);
      },
      off: (event, handler) => {
        const current = listenersRef.current.get(event);
        if (!current) return;
        current.delete(handler);
        if (current.size === 0) {
          listenersRef.current.delete(event);
        }
      },
      disconnect: () => {
        shouldReconnectRef.current = false;
        if (reconnectTimerRef.current) {
          clearTimeout(reconnectTimerRef.current);
          reconnectTimerRef.current = null;
        }
        wsRef.current?.close();
        wsRef.current = null;
        setIsConnected(false);
      },
    }),
    [],
  );

  useEffect(() => {
    if (!isRealtimeEnabled) {
      setSocket(null);
      setIsConnected(false);
      console.info('Realtime is disabled. Set VITE_REALTIME_MODE=websocket to enable live admin events.');
      return;
    }

    shouldReconnectRef.current = true;
    setSocket(socketApi);

    const connect = () => {
      const token = localStorage.getItem('admin_token');
      if (!token) {
        setIsConnected(false);
        return;
      }

      const wsUrl = buildRealtimeUrl();
      wsUrl.searchParams.set('token', token);

      const ws = new WebSocket(wsUrl.toString());
      wsRef.current = ws;

      ws.onopen = () => {
        setIsConnected(true);
      };

      ws.onmessage = (event) => {
        try {
          const payload = JSON.parse(event.data);
          const eventName = payload?.event;
          if (typeof eventName !== 'string') return;

          const handlers = listenersRef.current.get(eventName);
          if (!handlers || handlers.size === 0) return;

          handlers.forEach((handler) => {
            handler(payload.data);
          });
        } catch (error) {
          console.error('Failed to parse realtime message', error);
        }
      };

      ws.onclose = () => {
        setIsConnected(false);
        wsRef.current = null;

        if (!shouldReconnectRef.current) return;
        reconnectTimerRef.current = setTimeout(connect, WEBSOCKET_RECONNECT_DELAY_MS);
      };

      ws.onerror = () => {
        ws.close();
      };
    };

    connect();

    return () => {
      shouldReconnectRef.current = false;
      if (reconnectTimerRef.current) {
        clearTimeout(reconnectTimerRef.current);
        reconnectTimerRef.current = null;
      }
      wsRef.current?.close();
      wsRef.current = null;
      setIsConnected(false);
    };
  }, [isRealtimeEnabled, socketApi]);

  return (
    <SocketContext.Provider value={{ socket, isConnected, isRealtimeEnabled }}>
      {children}
    </SocketContext.Provider>
  );
};
