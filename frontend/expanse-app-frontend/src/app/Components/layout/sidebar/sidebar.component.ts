import { Component } from '@angular/core';

@Component({
  selector: 'app-sidebar',
  template: `
    <aside class="sidebar">
      <nav>
        <ul>
          <li>
            <a routerLink="/expenses" routerLinkActive="active">
              <i class="icon">📊</i>
              <span>Expenses</span>
            </a>
          </li>
          <li>
            <a routerLink="/expenses/new" routerLinkActive="active">
              <i class="icon">➕</i>
              <span>Add Expense</span>
            </a>
          </li>
          <li>
            <a routerLink="/categories/new" routerLinkActive="active">
              <i class="icon">🏷️</i>
              <span>Add Category</span>
            </a>
          </li>
        </ul>
      </nav>
    </aside>
  `,
  styles: [`
    .sidebar {
      width: 250px;
      height: 100%;
      background-color: #f5f5f5;
      border-right: 1px solid #e0e0e0;
    }
    
    ul {
      list-style: none;
      padding: 0;
      margin: 0;
    }
    
    li {
      margin: 0;
    }
    
    a {
      display: flex;
      align-items: center;
      padding: 15px 20px;
      color: #333;
      text-decoration: none;
      transition: background-color 0.2s;
    }
    
    a:hover {
      background-color: #e0e0e0;
    }
    
    .active {
      background-color: #e0e0e0;
      font-weight: 500;
    }
    
    .icon {
      margin-right: 10px;
      font-size: 18px;
    }
  `]
})
export class SidebarComponent {}