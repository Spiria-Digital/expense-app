import { Component, OnInit } from '@angular/core';
import { AuthService } from './Core/services/auth.service';

@Component({
  selector: 'app-root',
  standalone: false,
  template: `
    <div class="app-container">
      <app-header *ngIf="isAuthenticated"></app-header>
      <div class="content-wrapper">
        <app-sidebar *ngIf="isAuthenticated" class="sidebar"></app-sidebar>
        <main class="main-content">
          <app-alert></app-alert>
          <router-outlet></router-outlet>
        </main>
      </div>
    </div>
  `,
  styles: [`
    .app-container {
      display: flex;
      flex-direction: column;
      height: 100vh;
    }
    .content-wrapper {
      display: flex;
      flex: 1;
      overflow: hidden;
    }
    .sidebar {
      width: 250px;
      overflow-y: auto;
    }
    .main-content {
      flex: 1;
      padding: 20px;
      overflow-y: auto;
    }
  `]
})
export class AppComponent implements OnInit {
  isAuthenticated = false;

  constructor(private authService: AuthService) {}

  ngOnInit() {
    this.authService.user$.subscribe(user => {
      this.isAuthenticated = !!user;
    });
  }
}