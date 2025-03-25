import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Expense } from '../../../Models/expense.model';
import { Category } from '../../../Models/category.model';
import { ExpenseService } from '../../../Core/services/expense.service';
import { CategoryService } from '../../../Core/services/category.service';
import { NotificationService } from '../../../Core/services/notif.service';

@Component({
  selector: 'app-expense-list',
  standalone: false,
  template: `
    <div class="expense-list-container">
      <div class="header">
        <h2>Expenses</h2>
        <button class="btn btn-primary" routerLink="/expenses/new">Add Expense</button>
      </div>
      
      <app-loading-spinner [show]="isLoading" message="Loading expenses..."></app-loading-spinner>
      
      <div *ngIf="!isLoading && expenses.length === 0" class="empty-state">
        <p>No expenses found. Create your first expense!</p>
        <button class="btn btn-primary" routerLink="/expenses/new">Add Expense</button>
      </div>
      
      <div class="expense-list" *ngIf="!isLoading && expenses.length > 0">
        <table>
          <thead>
            <tr>
              <th>Title</th>
              <th>Amount</th>
              <th>Date</th>
              <th>Category</th>
              <th>Merchant</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr *ngFor="let expense of expenses">
              <td>{{ expense.title }}</td>
              <td>{{ expense.amount | currency }}</td>
              <td>{{ expense.date | date }}</td>
              <td>{{ getCategoryName(expense.categoryId) }}</td>
              <td>{{ expense.merchant || '-' }}</td>
              <td class="actions">
                <button class="btn-icon" (click)="viewExpense(expense.id!)">
                  <i class="icon">👁️</i>
                </button>
                <button class="btn-icon" (click)="editExpense(expense.id!)">
                  <i class="icon">✏️</i>
                </button>
                <button class="btn-icon delete" (click)="confirmDelete(expense)">
                  <i class="icon">🗑️</i>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    
    <div class="modal" *ngIf="showDeleteModal">
      <div class="modal-content">
        <h3>Confirm Delete</h3>
        <p>Are you sure you want to delete the expense "{{ expenseToDelete?.title }}"?</p>
        <div class="modal-actions">
          <button class="btn btn-secondary" (click)="cancelDelete()">Cancel</button>
          <button class="btn btn-danger" (click)="deleteExpense()">Delete</button>
        </div>
      </div>
    </div>
  `,
  styles: [`
    .expense-list-container {
      background-color: white;
      border-radius: 8px;
      box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
      padding: 20px;
    }
    
    .header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 20px;
    }
    
    h2 {
      margin: 0;
      color: #333;
    }
    
    .empty-state {
      text-align: center;
      padding: 40px 0;
      color: #666;
    }
    
    .empty-state p {
      margin-bottom: 20px;
    }
    
    table {
      width: 100%;
      border-collapse: collapse;
    }
    
    th, td {
      padding: 12px 15px;
      text-align: left;
      border-bottom: 1px solid #e0e0e0;
    }
    
    th {
      background-color: #f5f5f5;
      font-weight: 500;
    }
    
    .actions {
      white-space: nowrap;
    }
    
    .btn {
      padding: 8px 16px;
      border: none;
      border-radius: 4px;
      font-size: 14px;
      cursor: pointer;
      transition: background-color 0.2s;
    }
    
    .btn-primary {
      background-color: #3f51b5;
      color: white;
    }
    
    .btn-primary:hover {
      background-color: #303f9f;
    }
    
    .btn-secondary {
      background-color: #f5f5f5;
      color: #333;
    }
    
    .btn-secondary:hover {
      background-color: #e0e0e0;
    }
    
    .btn-danger {
      background-color: #f44336;
      color: white;
    }
    
    .btn-danger:hover {
      background-color: #d32f2f;
    }
    
    .btn-icon {
      background: none;
      border: none;
      cursor: pointer;
      padding: 5px;
      margin: 0 3px;
      border-radius: 4px;
      transition: background-color 0.2s;
    }
    
    .btn-icon:hover {
      background-color: #f5f5f5;
    }
    
    .btn-icon.delete:hover {
      color: #f44336;
    }
    
    .modal {
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background-color: rgba(0, 0, 0, 0.5);
      display: flex;
      justify-content: center;
      align-items: center;
      z-index: 1000;
    }
    
    .modal-content {
      background-color: white;
      border-radius: 8px;
      padding: 20px;
      width: 90%;
      max-width: 500px;
      box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
    }
    
    .modal-content h3 {
      margin-top: 0;
    }
    
    .modal-actions {
      display: flex;
      justify-content: flex-end;
      margin-top: 20px;
      gap: 10px;
    }
  `]
})
export class ExpenseListComponent implements OnInit {
  expenses: Expense[] = [];
  categories: Category[] = [];
  isLoading = true;
  showDeleteModal = false;
  expenseToDelete: Expense | null = null;

  constructor(
    private expenseService: ExpenseService,
    private categoryService: CategoryService,
    private notificationService: NotificationService,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.loadExpenses();
    this.loadCategories();
  }

  loadExpenses(): void {
    this.isLoading = true;
    this.expenseService.getExpenses().subscribe({
      next: (expenses) => {
        this.expenses = expenses;
        this.isLoading = false;
      },
      error: () => {
        this.isLoading = false;
      }
    });
  }

  loadCategories(): void {
    this.categoryService.getCategories().subscribe({
      next: (categories) => {
        this.categories = categories;
      }
    });
  }

  getCategoryName(categoryId?: number): string {
    if (!categoryId) return 'Uncategorized';
    const category = this.categories.find(c => c.id === categoryId);
    return category ? category.name : 'Uncategorized';
  }

  viewExpense(id: number): void {
    this.router.navigate(['/expenses', id]);
  }

  editExpense(id: number): void {
    this.router.navigate(['/expenses', id, 'edit']);
  }

  confirmDelete(expense: Expense): void {
    this.expenseToDelete = expense;
    this.showDeleteModal = true;
  }

  cancelDelete(): void {
    this.showDeleteModal = false;
    this.expenseToDelete = null;
  }

  deleteExpense(): void {
    if (!this.expenseToDelete || !this.expenseToDelete.id) {
      return;
    }

    const expenseId = this.expenseToDelete.id;
    this.expenseService.deleteExpense(expenseId).subscribe({
      next: () => {
        this.expenses = this.expenses.filter(e => e.id !== expenseId);
        this.notificationService.success('Expense deleted successfully');
        this.showDeleteModal = false;
        this.expenseToDelete = null;
      },
      error: () => {
        this.showDeleteModal = false;
        this.expenseToDelete = null;
      }
    });
  }
}