import { useEffect } from 'react';
import { useQuery } from '@tanstack/react-query';
import { getSettings } from '../api/adminApi';
import { useSettingsStore } from '../stores/settingsStore';

export const BUSINESS_SETTINGS_QUERY_KEY = ['settings'] as const;

// Business settings are only ever read from the API. The last copy is kept in
// the persisted store so the business name renders immediately on later visits
// while a fresh copy loads.
export const useBusinessSettings = () => {
  const cached = useSettingsStore((state) => state.settings);
  const setSettings = useSettingsStore((state) => state.setSettings);

  const query = useQuery({
    queryKey: BUSINESS_SETTINGS_QUERY_KEY,
    queryFn: getSettings,
    placeholderData: cached ?? undefined,
  });

  useEffect(() => {
    if (query.data && !query.isPlaceholderData) {
      setSettings(query.data);
    }
  }, [query.data, query.isPlaceholderData, setSettings]);

  return {
    settings: query.data ?? null,
    // True only while nothing (not even a cached copy) can be shown.
    isLoading: query.isLoading,
    isError: query.isError && !query.data,
    refetch: query.refetch,
  };
};

// "Funkey Grab & Bite" -> "FG&B"
export const businessInitials = (name: string) =>
  name
    .split(/\s+/)
    .filter(Boolean)
    .map((word) => (word === '&' ? '&' : word[0].toUpperCase()))
    .join('')
    .slice(0, 4);
