export interface Notification {
    id: number;
    userId: number;
    type: string;
    title: string;
    message: string;
    isRead: boolean;
    referenceId?: number;
    referenceType?: string;
    createdAt: string;
    readAt?: string;
  }