import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable, tap } from 'rxjs';
import { Router } from '@angular/router';
import { environment } from '../../../environments/environment.uat';
import { 
  UserLoginRequest, 
  UserLoginResponse, 
  UserRegistrationRequest, 
  UserRegistrationResponse
} from '../../Models/auth.model';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private tokenExpirationTimer: any;
  private userSubject = new BehaviorSubject<{ token: string } | null>(null);
  user$ = this.userSubject.asObservable();

  constructor(private http: HttpClient, private router: Router) {
    this.autoLogin();
  }

  get token(): string | null {
    return this.userSubject.value?.token || null;
  }

  get isAuthenticated(): boolean {
    return !!this.token;
  }

  register(userData: UserRegistrationRequest): Observable<UserRegistrationResponse> {
    return this.http.post<UserRegistrationResponse>(
      `${environment.apiUrl}/auth/register`, 
      userData
    );
  }

  login(userData: UserLoginRequest): Observable<UserLoginResponse> {
    return this.http.post<UserLoginResponse>(
      `${environment.apiUrl}/auth/login`, 
      userData
    ).pipe(
      tap(response => {
        this.handleAuthentication(response.token);
      })
    );
  }

  autoLogin() {
    const userData = this.getUserFromLocalStorage();
    if (!userData) {
      return;
    }
    this.userSubject.next(userData);
  }

  logout() {
    this.userSubject.next(null);
    localStorage.removeItem('userData');
    if (this.tokenExpirationTimer) {
      clearTimeout(this.tokenExpirationTimer);
    }
    this.tokenExpirationTimer = null;
    this.router.navigate(['/login']);
  }

  private handleAuthentication(token: string) {
    const userData = { token };
    this.userSubject.next(userData);
    localStorage.setItem('userData', JSON.stringify(userData));
  }

  private getUserFromLocalStorage(): { token: string } | null {
    const userData = localStorage.getItem('userData');
    if (!userData) {
      return null;
    }
    return JSON.parse(userData);
  }
}