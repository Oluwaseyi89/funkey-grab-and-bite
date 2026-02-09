export interface InventoryItem {
    id: number;
    menuItemId: number;
    name: string;
    currentStock: number;
    minimumStock: number;
    reorderPoint: number;
    unit: string;
    isActive: boolean;
    lastRestocked: string;
    createdAt: string;
    updatedAt: string;
  }
  
  export interface InventoryAlert {
    id: number;
    inventoryItemId: number;
    alertType: string;
    message: string;
    isResolved: boolean;
    createdAt: string;
    resolvedAt?: string;
    readAt?: string;
  }
  
  export interface InventoryUpdate {
    menuItemId: number;
    quantity: number;
    operation: 'add' | 'subtract' | 'set';
    reason: string;
    notes?: string;
  }
  
  export interface InventoryHistory {
    id: number;
    inventoryItemId: number;
    previousStock: number;
    newStock: number;
    change: number;
    operation: string;
    reason: string;
    notes: string;
    createdAt: string;
  }