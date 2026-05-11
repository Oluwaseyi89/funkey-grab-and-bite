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

const WEBSOCKET_RECONNECT_BASE_DELAY_MS = 1000;
const WEBSOCKET_RECONNECT_MAX_DELAY_MS = 30000;

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
  const reconnectAttemptRef = useRef(0);
  const isConnectingRef = useRef(false);
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
        isConnectingRef.current = false;
        reconnectAttemptRef.current = 0;
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

    const scheduleReconnect = () => {
      if (!shouldReconnectRef.current) return;

      reconnectAttemptRef.current += 1;
      const exponentialDelay = Math.min(
        WEBSOCKET_RECONNECT_BASE_DELAY_MS * (2 ** (reconnectAttemptRef.current - 1)),
        WEBSOCKET_RECONNECT_MAX_DELAY_MS,
      );
      const jitter = Math.floor(Math.random() * 500);
      const delay = exponentialDelay + jitter;

      if (reconnectTimerRef.current) {
        clearTimeout(reconnectTimerRef.current);
      }
      reconnectTimerRef.current = setTimeout(connect, delay);
    };

    const connect = () => {
      if (isConnectingRef.current) {
        return;
      }

      const existingWs = wsRef.current;
      if (existingWs && (existingWs.readyState === WebSocket.OPEN || existingWs.readyState === WebSocket.CONNECTING)) {
        return;
      }

      const token = localStorage.getItem('admin_token');
      if (!token) {
        reconnectAttemptRef.current = 0;
        setIsConnected(false);
        return;
      }

      isConnectingRef.current = true;
      const wsUrl = buildRealtimeUrl();
      const ws = new WebSocket(wsUrl.toString(), [`jwt.${token}`]);
      wsRef.current = ws;

      ws.onopen = () => {
        isConnectingRef.current = false;
        reconnectAttemptRef.current = 0;
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
        isConnectingRef.current = false;
        setIsConnected(false);
        wsRef.current = null;

        if (!shouldReconnectRef.current) return;
        scheduleReconnect();
      };

      ws.onerror = () => {
        ws.close();
      };
    };

    connect();

    return () => {
      shouldReconnectRef.current = false;
      isConnectingRef.current = false;
      reconnectAttemptRef.current = 0;
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
