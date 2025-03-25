import { Component } from '@angular/core';
import { AuthService } from '../../../Core/services/auth.service';

@Component({
  selector: 'app-header',
  standalone: false,
  template: `
    <header class="header">
      <div class="logo">
        <h1>Expense Manager</h1>
      </div>
      <div class="actions">
        <button class="logout-btn" (click)="logout()">Logout</button>
      </div>
    </header>
  `,
  styles: [`
    .header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 0 20px;
      height: 60px;
      background-color: #3f51b5;
      color: white;
      box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
    }
    .logo h1 {
      margin: 0;
      font-size: 1.5rem;
      font-weight: 500;
    }
    .logout-btn {
      background-color: transparent;
      border: 1px solid white;
      color: white;
      padding: 8px 16px;
      border-radius: 4px;
      cursor: pointer;
      font-size: 14px;
      transition: all 0.2s;
    }
    .logout-btn:hover {
      background-color: white;
      color: #3f51b5;
    }
  `]
})
export class HeaderComponent {
  constructor(private authService: AuthService) {}

  logout(): void {
    this.authService.logout();
  }
}