import { apiRequest } from './api';
import {
  User,
  UserProfile,
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  UpdateProfileRequest,
} from '@/types/auth';

export class AuthService {
  // Authentication endpoints
  static async login(credentials: LoginRequest): Promise<LoginResponse> {
    return apiRequest<LoginResponse>('POST', '/auth/login', credentials);
  }

  static async register(data: RegisterRequest): Promise<User> {
    return apiRequest<User>('POST', '/auth/register', data);
  }

  static logout(): void {
    localStorage.removeItem('auth_token');
    localStorage.removeItem('user_data');
  }

  // User profile endpoints
  static async getProfile(): Promise<UserProfile> {
    return apiRequest<UserProfile>('GET', '/users/profile');
  }

  static async getUserById(userId: string): Promise<User> {
    return apiRequest<User>('GET', `/users/${userId}`);
  }

  static async updateProfile(data: UpdateProfileRequest): Promise<User> {
    return apiRequest<User>('PUT', '/users/profile', data);
  }

  // Follow system
  static async followUser(userId: string): Promise<void> {
    return apiRequest<void>('POST', `/users/${userId}/follow`);
  }

  static async unfollowUser(userId: string): Promise<void> {
    return apiRequest<void>('DELETE', `/users/${userId}/follow`);
  }

  static async getFollowers(userId: string, page = 1, limit = 20) {
    return apiRequest('GET', `/users/${userId}/followers`, null, {
      params: { page, limit }
    });
  }

  static async getFollowing(userId: string, page = 1, limit = 20) {
    return apiRequest('GET', `/users/${userId}/following`, null, {
      params: { page, limit }
    });
  }

  // Search users
  static async searchUsers(query: string, page = 1, limit = 20) {
    return apiRequest('GET', '/users/search', null, {
      params: { q: query, page, limit }
    });
  }

  // Token management
  static setToken(token: string): void {
    localStorage.setItem('auth_token', token);
  }

  static getToken(): string | null {
    return localStorage.getItem('auth_token');
  }

  static isAuthenticated(): boolean {
    const token = this.getToken();
    return !!token;
  }

  // User data management
  static setUserData(user: User): void {
    localStorage.setItem('user_data', JSON.stringify(user));
  }

  static getUserData(): User | null {
    const userData = localStorage.getItem('user_data');
    return userData ? JSON.parse(userData) : null;
  }

  static clearUserData(): void {
    localStorage.removeItem('user_data');
  }
}

export default AuthService;