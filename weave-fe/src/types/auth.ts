export interface User {
  id: string;
  username: string;
  email: string;
  profile_image?: string;
  bio?: string;
  is_verified: boolean;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface UserProfile extends User {
  followers_count: number;
  following_count: number;
  weaves_count: number;
  contributions_count: number;
}

export interface SendEmailVerificationRequest {
  email: string;
}

export interface VerifyEmailRequest {
  code: string;
}

export interface EmailVerificationResponse {
  message: string;
  expires_in: number;
  code?: string; // Only for development
}

export interface LoginResponse {
  user: User;
  token: string;
}

export interface UpdateProfileRequest {
  profile_image?: string;
  bio?: string;
}

export interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}

export interface AuthActions {
  sendEmailVerification: (data: SendEmailVerificationRequest) => Promise<EmailVerificationResponse>;
  verifyEmail: (data: VerifyEmailRequest) => Promise<void>;
  logout: () => void;
  updateProfile: (data: UpdateProfileRequest) => Promise<void>;
  setLoading: (loading: boolean) => void;
}