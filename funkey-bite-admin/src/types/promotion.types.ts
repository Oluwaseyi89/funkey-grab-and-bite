export type PromotionType = 'percentage' | 'fixed' | 'bogo';
export type PromotionStatus = 'active' | 'inactive' | 'expired';

export interface Promotion {
  id: number;
  code: string;
  title: string;
  description: string;
  promotionType: PromotionType;
  discountValue: number;
  maxDiscount?: number;
  minOrderAmount?: number;
  validFrom: string;
  validUntil: string;
  usageLimit?: number;
  usedCount: number;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface PromotionValidation {
  isValid: boolean;
  message?: string;
  discount: number;
  promotionId: number;
}

export interface PromotionUsage {
  id: number;
  promotionId: number;
  orderId: number;
  customerId?: number;
  discountApplied: number;
  createdAt: string;
}