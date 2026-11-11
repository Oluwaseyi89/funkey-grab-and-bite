// The sound and browser notification for live alerts.

let audio: AudioContext | null = null;

// Browsers only allow sound after the person has tapped the page once, so the
// audio context is created on the first tap and reused.
export const unlockAudio = () => {
  if (audio || typeof window === 'undefined') return;
  const Ctx = window.AudioContext || (window as unknown as { webkitAudioContext?: typeof AudioContext }).webkitAudioContext;
  if (Ctx) audio = new Ctx();
};

// A short, bright two-note chime, made in the browser (no sound file).
// Urgent alerts play it twice.
export const playChime = (urgent = false) => {
  if (!audio) return;
  if (audio.state === 'suspended') void audio.resume();
  const notes = urgent ? [880, 660, 880, 660] : [660, 990];
  notes.forEach((frequency, index) => {
    const start = audio!.currentTime + index * 0.16;
    const osc = audio!.createOscillator();
    const gain = audio!.createGain();
    osc.type = 'sine';
    osc.frequency.value = frequency;
    gain.gain.setValueAtTime(0.0001, start);
    gain.gain.exponentialRampToValueAtTime(0.25, start + 0.02);
    gain.gain.exponentialRampToValueAtTime(0.0001, start + 0.3);
    osc.connect(gain).connect(audio!.destination);
    osc.start(start);
    osc.stop(start + 0.32);
  });
};

export const browserNotificationsSupported = () => typeof window !== 'undefined' && 'Notification' in window;

export const browserPermission = (): NotificationPermission | 'unsupported' =>
  browserNotificationsSupported() ? Notification.permission : 'unsupported';

// Asks for permission. Must be called from a tap.
export const requestBrowserPermission = async (): Promise<NotificationPermission | 'unsupported'> => {
  if (!browserNotificationsSupported()) return 'unsupported';
  return Notification.requestPermission();
};

// Shows a system notification while the admin isn't looking at the tab.
// Tapping it brings the tab forward and opens the link.
export const showBrowserNotification = (title: string, body: string, onClick?: () => void) => {
  if (!browserNotificationsSupported() || Notification.permission !== 'granted') return;
  if (document.visibilityState === 'visible') return;
  try {
    const n = new Notification(title, { body, icon: '/brand/icon-192.png', tag: `${title}-${body}` });
    n.onclick = () => {
      window.focus();
      onClick?.();
      n.close();
    };
  } catch {
    // Some phones only allow notifications from a service worker.
  }
};
