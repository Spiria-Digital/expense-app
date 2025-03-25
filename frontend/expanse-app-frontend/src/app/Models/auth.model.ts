export interface UserLoginRequest {
    email: string;
    password: string;
  }
  
  export interface UserLoginResponse {
    token: string;
  }
  
  export interface UserRegistrationRequest {
    email: string;
    firstName: string;
    lastName: string;
    password: string;
  }
  
  export interface UserRegistrationResponse {
    message: string;
  }
  
  export interface OutgoingUser {
    id: number;
    firstName: string;
    lastName: string;
    isAdmin: boolean;
  }