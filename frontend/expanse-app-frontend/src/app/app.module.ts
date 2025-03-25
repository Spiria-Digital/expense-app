import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';

// Interceptors
import { AuthInterceptor } from './Core/interceptors/auth.interceptor';
import { ErrorInterceptor } from './Core/interceptors/error.interceptor';

// Components will be imported here
import { LoginComponent } from './Components/auth/login/login.component';
import { RegisterComponent } from './Components/auth/register/register.component';
import { ExpenseListComponent } from './Components/expenses/expense-list/expense-list.component';
import { ExpenseFormComponent } from './Components/expenses/expense-form/expense-form.component';
import { ExpenseDetailsComponent } from './Components/expenses/expense-details/expense-details.component';
import { CreateCategoryComponent } from './Components/categories/create-category/create-category.component';
import { HeaderComponent } from './Components/layout/header/header.component';
import { SidebarComponent } from './Components/layout/sidebar/sidebar.component';
import { AlertComponent } from './Components/shared/alert/alert.component';
import { LoadingSpinnerComponent } from './Components/shared/loading-spinner/loading-spinner.component';

@NgModule({
  declarations: [
    AppComponent,
    LoginComponent,
    RegisterComponent,
    ExpenseListComponent,
    ExpenseFormComponent,
    ExpenseDetailsComponent,
    CreateCategoryComponent,
    HeaderComponent,
    SidebarComponent,
    AlertComponent,
    LoadingSpinnerComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    ReactiveFormsModule,
    FormsModule
  ],
  providers: [
    { provide: HTTP_INTERCEPTORS, useClass: AuthInterceptor, multi: true },
    { provide: HTTP_INTERCEPTORS, useClass: ErrorInterceptor, multi: true }
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }