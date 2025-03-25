import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { AuthGuard } from './Core/guards/auth.guard';
import { LoginComponent } from './Components/auth/login/login.component';
import { RegisterComponent } from './Components/auth/register/register.component';
import { ExpenseListComponent } from './Components/expenses/expense-list/expense-list.component';
import { ExpenseFormComponent } from './Components/expenses/expense-form/expense-form.component';
import { ExpenseDetailsComponent } from './Components/expenses/expense-details/expense-details.component';
import { CreateCategoryComponent } from './Components/categories/create-category/create-category.component';

const routes: Routes = [
  { path: '', redirectTo: '/expenses', pathMatch: 'full' },
  { path: 'login', component: LoginComponent },
  { path: 'register', component: RegisterComponent },
  { 
    path: 'expenses', 
    component: ExpenseListComponent, 
    canActivate: [AuthGuard] 
  },
  { 
    path: 'expenses/new', 
    component: ExpenseFormComponent, 
    canActivate: [AuthGuard] 
  },
  { 
    path: 'expenses/:id', 
    component: ExpenseDetailsComponent, 
    canActivate: [AuthGuard] 
  },
  { 
    path: 'expenses/:id/edit', 
    component: ExpenseFormComponent, 
    canActivate: [AuthGuard] 
  },
  { 
    path: 'categories/new', 
    component: CreateCategoryComponent, 
    canActivate: [AuthGuard] 
  },
  { path: '**', redirectTo: '/expenses' }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }