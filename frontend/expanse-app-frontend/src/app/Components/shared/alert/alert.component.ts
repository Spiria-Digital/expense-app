// src/app/components/shared/alert/alert.component.ts
import { Component, OnInit, OnDestroy } from '@angular/core';
import { Subscription } from 'rxjs';
import { NotificationService, Notification, NotificationType } from '../../../Core/services/notif.service';


@Component({
  selector: 'app-alert',
  standalone: false,
  template: `
    <div *ngIf="notification" class="alert" [ngClass]="alertTypeClass">
      <span>{{ notification.message }}</span>
      <button class="close-btn" (click)="closeAlert()">&times;</button>
    </div>
  `,
  styles: [`
    .alert {
      padding: 15px;
      margin-bottom: 20px;
      border: 1px solid transparent;
      border-radius: 4px;
      position: relative;
    }
    .success {
      color: #155724;
      background-color: #d4edda;
      border-color: #c3e6cb;
    }
    .error {
      color: #721c24;
      background-color: #f8d7da;
      border-color: #f5c6cb;
    }
    .info {
      color: #0c5460;
      background-color: #d1ecf1;
      border-color: #bee5eb;
    }
    .warning {
      color: #856404;
      background-color: #fff3cd;
      border-color: #ffeeba;
    }
    .close-btn {
      position: absolute;
      top: 10px;
      right: 15px;
      background: none;
      border: none;
      font-size: 20px;
      cursor: pointer;
    }
  `]
})
export class AlertComponent implements OnInit, OnDestroy {
  notification: Notification | null = null;
  alertTypeClass = '';
  private subscription: Subscription | null = null;
  private timeout: any;

  constructor(private notificationService: NotificationService) { }

  ngOnInit() {
    this.subscription = this.notificationService.notifications$.subscribe(notification => {
      this.notification = notification;
      this.setAlertClass();
      this.setAutoClose();
    });
  }

  ngOnDestroy() {
    if (this.subscription) {
      this.subscription.unsubscribe();
    }
    if (this.timeout) {
      clearTimeout(this.timeout);
    }
  }

  closeAlert() {
    this.notification = null;
    if (this.timeout) {
      clearTimeout(this.timeout);
    }
  }

  private setAlertClass() {
    if (!this.notification) return;

    switch (this.notification.type) {
      case NotificationType.SUCCESS:
        this.alertTypeClass = 'success';
        break;
      case NotificationType.ERROR:
        this.alertTypeClass = 'error';
        break;
      case NotificationType.INFO:
        this.alertTypeClass = 'info';
        break;
      case NotificationType.WARNING:
        this.alertTypeClass = 'warning';
        break;
      default:
        this.alertTypeClass = 'info';
    }
  }

  private setAutoClose() {
    if (this.timeout) {
      clearTimeout(this.timeout);
    }
    
    // Auto close after 5 seconds for non-error notifications
    if (this.notification && this.notification.type !== NotificationType.ERROR) {
      this.timeout = setTimeout(() => {
        this.notification = null;
      }, 5000);
    }
  }
}