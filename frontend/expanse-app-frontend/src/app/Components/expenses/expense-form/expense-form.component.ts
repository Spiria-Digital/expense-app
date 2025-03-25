import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { ExpenseService } from '../../../Core/services/expense.service';
import { CategoryService } from '../../../Core/services/category.service';
import { NotificationService } from '../../../Core/services/notif.service';
import { Category } from '../../../Models/category.model';
import { Expense, CreateExpenseRequest } from '../../../Models/expense.model';

@Component({
  selector: 'app-expense-form',
  standalone: false,
  template: `
    <div class="expense-form-container">
      <div class="header">
        <h2>{{ isEditMode ? 'Edit Expense' : 'Add New Expense' }}</h2>
        <button class="btn btn-secondary" routerLink="/expenses">Back to List</button>
      </div>
      
      <app-loading-spinner [show]="isLoading" [message]="loadingMessage"></app-loading-spinner>
      
      <form *ngIf="!isLoading" [formGroup]="expenseForm" (ngSubmit)="onSubmit()">
        <div class="form-group">
          <label for="title">Title *</label>
          <input 
            type="text" 
            id="title" 
            formControlName="title" 
            class="form-control"
            [ngClass]="{'invalid': isFieldInvalid('title')}" 
          />
          <div class="error-message" *ngIf="isFieldInvalid('title')">
            <span *ngIf="expenseForm.get('title')?.errors?.['required']">Title is required</span>
          </div>
        </div>
        
        <div class="form-group">
          <label for="amount">Amount</label>
          <input 
            type="number" 
            id="amount" 
            formControlName="amount" 
            class="form-control"
            step="0.01"
          />
        </div>
        
        <div class="form-group">
          <label for="date">Date</label>
          <input 
            type="date" 
            id="date" 
            formControlName="date" 
            class="form-control"
          />
        </div>
        
        <div class="form-group">
          <label for="categoryId">Category</label>
          <select 
            id="categoryId" 
            formControlName="categoryId" 
            class="form-control"
          >
            <option [ngValue]="null">Uncategorized</option>
            <option *ngFor="let category of categories" [ngValue]="category.id">
              {{ category.name }}
            </option>
          </select>
        </div>
        
        <div class="form-group">
          <label for="merchant">Merchant</label>
          <input 
            type="text" 
            id="merchant" 
            formControlName="merchant" 
            class="form-control"
          />
        </div>
        
        <div class="form-group">
          <label for="description">Description</label>
          <textarea 
            id="description" 
            formControlName="description" 
            class="form-control"
            rows="4"
          ></textarea>
        </div>
        
        <div class="form-actions">
          <button 
            type="button" 
            class="btn btn-secondary" 
            routerLink="/expenses"
          >
            Cancel
          </button>
          <button 
            type="submit" 
            class="btn btn-primary" 
            [disabled]="expenseForm.invalid || isSaving"
          >
            {{ isSaving ? (isEditMode ? 'Updating...' : 'Creating...') : (isEditMode ? 'Update' : 'Create') }}
          </button>
        </div>
      </form>
    </div>
  `,
  styles: [`
    .expense-form-container {
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
    
    .form-group {
      margin-bottom: 20px;
    }
    
    label {
      display: block;
      margin-bottom: 8px;
      font-weight: 500;
    }
    
    .form-control {
      width: 100%;
      padding: 10px;
      border: 1px solid #ddd;
      border-radius: 4px;
      font-size: 16px;
      background-color: white;
    }
    
    textarea.form-control {
      resize: vertical;
    }
    
    .invalid {
      border-color: #dc3545;
    }
    
    .error-message {
      color: #dc3545;
      font-size: 14px;
      margin-top: 5px;
    }
    
    .form-actions {
      display: flex;
      justify-content: flex-end;
      gap: 10px;
      margin-top: 30px;
    }
    
    .btn {
      padding: 10px 20px;
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
    
    .btn-primary:disabled {
      background-color: #c5cae9;
      cursor: not-allowed;
    }
    
    .btn-secondary {
      background-color: #f5f5f5;
      color: #333;
    }
    
    .btn-secondary:hover {
      background-color: #e0e0e0;
    }
  `]
})
export class ExpenseFormComponent implements OnInit {
  expenseForm!: FormGroup;
  categories: Category[] = [];
  isLoading = true;
  isSaving = false;
  isEditMode = false;
  expenseId: number | null = null;
  loadingMessage = 'Loading...';

  constructor(
    private fb: FormBuilder,
    private expenseService: ExpenseService,
    private categoryService: CategoryService,
    private notificationService: NotificationService,
    private route: ActivatedRoute,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.initForm();
    this.loadCategories();
    
    // Check if we're in edit mode
    this.route.paramMap.subscribe(params => {
      const id = params.get('id');
      if (id) {
        this.expenseId = +id;
        this.isEditMode = true;
        this.loadExpense(this.expenseId);
      } else {
        this.isLoading = false;
      }
    });
  }

  initForm(): void {
    this.expenseForm = this.fb.group({
      title: ['', Validators.required],
      amount: [null],
      date: [this.formatDate(new Date())],
      categoryId: [null],
      merchant: [''],
      description: ['']
    });
  }

  loadCategories(): void {
    this.categoryService.getCategories().subscribe({
      next: (categories) => {
        this.categories = categories;
      },
      error: () => {
        this.notificationService.error('Failed to load categories');
      }
    });
  }

  loadExpense(id: number): void {
    this.loadingMessage = 'Loading expense...';
    this.isLoading = true;
    
    this.expenseService.getExpense(id).subscribe({
      next: (expense) => {
        this.updateForm(expense);
        this.isLoading = false;
      },
      error: () => {
        this.notificationService.error('Failed to load expense');
        this.router.navigate(['/expenses']);
      }
    });
  }

  updateForm(expense: Expense): void {
    this.expenseForm.patchValue({
      title: expense.title,
      amount: expense.amount,
      date: expense.date ? this.formatDate(new Date(expense.date)) : null,
      categoryId: expense.categoryId,
      merchant: expense.merchant,
      description: expense.description
    });
  }

  formatDate(date: Date): string {
    const year = date.getFullYear();
    const month = (date.getMonth() + 1).toString().padStart(2, '0');
    const day = date.getDate().toString().padStart(2, '0');
    return `${year}-${month}-${day}`;
  }

  isFieldInvalid(field: string): boolean {
    const formControl = this.expenseForm.get(field);
    return !!formControl && formControl.invalid && (formControl.dirty || formControl.touched);
  }

  onSubmit(): void {
    if (this.expenseForm.invalid) {
      return;
    }

    this.isSaving = true;
    const expenseData: CreateExpenseRequest = this.expenseForm.value;

    if (this.isEditMode && this.expenseId) {
      this.updateExpense(this.expenseId, expenseData);
    } else {
      this.createExpense(expenseData);
    }
  }

  createExpense(expenseData: CreateExpenseRequest): void {
    this.expenseService.createExpense(expenseData).subscribe({
      next: () => {
        this.notificationService.success('Expense created successfully');
        this.router.navigate(['/expenses']);
      },
      error: () => {
        this.isSaving = false;
      },
      complete: () => {
        this.isSaving = false;
      }
    });
  }

  updateExpense(id: number, expenseData: CreateExpenseRequest): void {
    this.expenseService.updateExpense(id, expenseData).subscribe({
      next: () => {
        this.notificationService.success('Expense updated successfully');
        this.router.navigate(['/expenses']);
      },
      error: () => {
        this.isSaving = false;
      },
      complete: () => {
        this.isSaving = false;
      }
    });
  }
}