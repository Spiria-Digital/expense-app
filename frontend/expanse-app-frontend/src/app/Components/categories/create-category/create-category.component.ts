import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { CategoryService } from '../../../Core/services/category.service';
import { NotificationService } from '../../../Core/services/notif.service';
import { Category } from '../../../Models/category.model';

@Component({
  selector: 'app-create-category',
  standalone: false,
  template: `
    <div class="category-form-container">
      <div class="header">
        <h2>Create New Category</h2>
        <button class="btn btn-secondary" routerLink="/expenses">Back to Expenses</button>
      </div>
      
      <div class="row">
        <div class="col">
          <form [formGroup]="categoryForm" (ngSubmit)="onSubmit()">
            <div class="form-group">
              <label for="name">Category Name *</label>
              <input 
                type="text" 
                id="name" 
                formControlName="name" 
                class="form-control"
                [ngClass]="{'invalid': isFieldInvalid('name')}" 
              />
              <div class="error-message" *ngIf="isFieldInvalid('name')">
                <span *ngIf="categoryForm.get('name')?.errors?.['required']">Category name is required</span>
              </div>
            </div>
            
            <div class="form-actions">
              <button 
                type="submit" 
                class="btn btn-primary" 
                [disabled]="categoryForm.invalid || isSubmitting"
              >
                {{ isSubmitting ? 'Creating...' : 'Create Category' }}
              </button>
            </div>
          </form>
        </div>
        
        <div class="col">
          <div class="categories-list">
            <h3>Existing Categories</h3>
            <app-loading-spinner [show]="isLoading" message="Loading categories..."></app-loading-spinner>
            
            <div *ngIf="!isLoading && categories.length === 0" class="empty-state">
              <p>No categories found. Create your first category!</p>
            </div>
            
            <ul *ngIf="!isLoading && categories.length > 0">
              <li *ngFor="let category of categories">{{ category.name }}</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  `,
  styles: [`
    .category-form-container {
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
    
    h2, h3 {
      margin: 0;
      color: #333;
    }
    
    .row {
      display: flex;
      gap: 30px;
    }
    
    .col {
      flex: 1;
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
    
    .categories-list {
      border: 1px solid #e0e0e0;
      border-radius: 8px;
      padding: 20px;
    }
    
    .categories-list h3 {
      margin-bottom: 20px;
    }
    
    .empty-state {
      color: #666;
      padding: 20px 0;
    }
    
    ul {
      list-style-type: none;
      padding: 0;
      margin: 0;
    }
    
    li {
      padding: 10px 0;
      border-bottom: 1px solid #eee;
    }
    
    li:last-child {
      border-bottom: none;
    }
  `]
})
export class CreateCategoryComponent implements OnInit {
  categoryForm!: FormGroup;
  categories: Category[] = [];
  isLoading = true;
  isSubmitting = false;

  constructor(
    private fb: FormBuilder,
    private categoryService: CategoryService,
    private notificationService: NotificationService,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.initForm();
    this.loadCategories();
  }

  initForm(): void {
    this.categoryForm = this.fb.group({
      name: ['', Validators.required]
    });
  }

  loadCategories(): void {
    this.isLoading = true;
    this.categoryService.getCategories().subscribe({
      next: (categories) => {
        this.categories = categories;
        this.isLoading = false;
      },
      error: () => {
        this.isLoading = false;
      }
    });
  }

  isFieldInvalid(field: string): boolean {
    const formControl = this.categoryForm.get(field);
    return !!formControl && formControl.invalid && (formControl.dirty || formControl.touched);
  }

  onSubmit(): void {
    if (this.categoryForm.invalid) {
      return;
    }

    this.isSubmitting = true;
    const category: Category = this.categoryForm.value;

    this.categoryService.createCategory(category).subscribe({
      next: (newCategory) => {
        this.categories.push(newCategory);
        this.notificationService.success('Category created successfully');
        this.categoryForm.reset();
      },
      error: () => {
        this.isSubmitting = false;
      },
      complete: () => {
        this.isSubmitting = false;
      }
    });
  }
}