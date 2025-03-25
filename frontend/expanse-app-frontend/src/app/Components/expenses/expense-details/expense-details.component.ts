import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { ExpenseService } from '../../../Core/services/expense.service';
import { CategoryService } from '../../../Core/services/category.service';
import { NotificationService } from '../../../Core/services/notif.service';
import { Expense } from '../../../Models/expense.model';
import { Category } from '../../../Models/category.model';

@Component({
  selector: 'app-expense-details',
  standalone: false,
  template: `
    <div class="expense-details-container">
      <div class="header">
        <h2>Expense Details</h2>
        <div class="header-actions">
          <button class="btn btn-secondary" routerLink="/expenses">Back to List</button>
          <button class="btn btn-primary" *ngIf="expense" (click)="editExpense()">Edit</button>
          <button class="btn btn-danger" *ngIf="expense" (click)="confirmDelete()">Delete</button>
        </div>
      </div>
      
      <app-loading-spinner [show]="isLoading" message="Loading expense details..."></app-loading-spinner>
      
      <div class="expense-card" *ngIf="!isLoading && expense">
        <h3>{{ expense.title }}</h3>
        
        <div class="expense-info">
          <div class="info-item">
            <span class="label">Amount:</span>
            <span class="value">{{ expense.amount | currency }}</span>
          </div>
          
          <div class="info-item">
            <span class="label">Date:</span>
            <span class="value">{{ expense.date | date }}</span>
          </div>
          
          <div class="info-item">
            <span class="label">Category:</span>
            <span class="value">{{ getCategoryName(expense.categoryId) }}</span>
          </div>
          
          <div class="info-item">
            <span class="label">Merchant:</span>
            <span class="value">{{ expense.merchant || '-' }}</span>
          </div>
          
          <div class="info-item" *ngIf="expense.description">
            <span class="label">Description:</span>
            <p class="description">{{ expense.description }}</p>
          </div>
        </div>
      </div>
    </div>
    
    <div class="modal" *ngIf="showDeleteModal">
      <div class="modal-content">
        <h3>Confirm Delete</h3>
        <p>Are you sure you want to delete the expense "{{ expense?.title }}"?</p>
        <div class="modal-actions">
          <button class="btn btn-secondary" (click)="cancelDelete()">Cancel</button>
          <button class="btn btn-danger" (click)="deleteExpense()">Delete</button>
        </div>
      </div>
    </div>
  `,
  styles: [`
    .expense-details-container {
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
    
    .header-actions {
      display: flex;
      gap: 10px;
    }
    
    h2 {
      margin: 0;
      color: #333;
    }
    
    .expense-card {
      padding: 20px;
      border: 1px solid #e0e0e0;
      border-radius: 8px;
    }
    
    h3 {
      margin-top: 0;
      margin-bottom: 20px;
      color: #333;
      font-size: 1.5rem;
    }
    
    .expense-info {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
      gap: 20px;
    }
    
    .info-item {
      margin-bottom: 15px;
    }
    
    .label {
      font-weight: 500;
      color: #666;
      display: block;
      margin-bottom: 5px;
    }
    
    .value {
      color: #333;
      font-size: 1.1rem;
    }
    
    .description {
      margin-top: 5px;
      white-space: pre-line;
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
export class ExpenseDetailsComponent implements OnInit {
  expense: Expense | null = null;
  categories: Category[] = [];
  isLoading = true;
  showDeleteModal = false;

  constructor(
    private expenseService: ExpenseService,
    private categoryService: CategoryService,
    private notificationService: NotificationService,
    private route: ActivatedRoute,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.loadCategories();
    this.route.paramMap.subscribe(params => {
      const id = params.get('id');
      if (id) {
        this.loadExpense(+id);
      } else {
        this.router.navigate(['/expenses']);
      }
    });
  }

  loadExpense(id: number): void {
    this.isLoading = true;
    this.expenseService.getExpense(id).subscribe({
      next: (expense) => {
        this.expense = expense;
        this.isLoading = false;
      },
      error: () => {
        this.notificationService.error('Failed to load expense');
        this.router.navigate(['/expenses']);
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

  editExpense(): void {
    if (this.expense?.id) {
      this.router.navigate(['/expenses', this.expense.id, 'edit']);
    }
  }

  confirmDelete(): void {
    this.showDeleteModal = true;
  }

  cancelDelete(): void {
    this.showDeleteModal = false;
  }

  deleteExpense(): void {
    if (!this.expense?.id) {
      return;
    }

    const expenseId = this.expense.id;
    this.expenseService.deleteExpense(expenseId).subscribe({
      next: () => {
        this.notificationService.success('Expense deleted successfully');
        this.router.navigate(['/expenses']);
      },
      error: () => {
        this.showDeleteModal = false;
      },
      complete: () => {
        this.showDeleteModal = false;
      }
    });
  }
}