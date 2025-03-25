export interface Expense {
    id?: number;
    title: string;
    amount?: number;
    description?: string;
    date?: string;
    categoryId?: number;
    merchant?: string;
    ownerID?: number;
  }
  
  export interface CreateExpenseRequest {
    title: string;
    amount?: number;
    description?: string;
    date?: string;
    categoryId?: number;
    merchant?: string;
  }
  