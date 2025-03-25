import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';

export enum NotificationType {

  INFO = 'info',
  WARNING = 'warning',
  SUCCESS = 'success',
  ERROR = 'error'
  
}

export interface Notification {
  type: NotificationType;
  message: string;
}

@Injectable({
  providedIn: 'root'
})
export class NotificationService {
  private notificationSubject = new Subject<Notification>();
  notifications$ = this.notificationSubject.asObservable();

  success(message: string): void {
    this.notify(NotificationType.SUCCESS, message);
  }

  error(message: string): void {
    this.notify(NotificationType.ERROR, message);
  }

  info(message: string): void {
    this.notify(NotificationType.INFO, message);
  }

  warning(message: string): void {
    this.notify(NotificationType.WARNING, message);
  }

  private notify(type: NotificationType, message: string): void {
    this.notificationSubject.next({ type, message });
  }
}