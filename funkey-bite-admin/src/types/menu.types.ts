export interface MenuCategory {
    id: number;
    name: string;
    description: string;
    displayOrder: number;
    isActive: boolean;
    createdAt: string;
  }
  
  export interface NutritionalInfo {
    calories: number;
    protein: number;
    carbs: number;
    fat: number;
  }
  
  export interface MenuItem {
    id: number;
    categoryId: number;
    name: string;
    description: string;
    price: number;
    imageUrl: string;
    isAvailable: boolean;
    isPreOrder: boolean;
    preparationTime: number;
    tags: string[];
    nutritionalInfo?: NutritionalInfo;
    createdAt: string;
    updatedAt?: string;
  }

