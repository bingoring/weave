import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import {
  AuthState,
  AuthActions,
  SendEmailVerificationRequest,
  VerifyEmailRequest,
  EmailVerificationResponse,
  UpdateProfileRequest,
} from '../types/auth';
import AuthService from '../services/authService';

interface AuthStore extends AuthState, AuthActions {}

export const useAuthStore = create<AuthStore>()(
  persist(
    (set) => ({
      // Initial state
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,

      // Actions
      sendEmailVerification: async (data: SendEmailVerificationRequest): Promise<EmailVerificationResponse> => {
        set({ isLoading: true });
        try {
          const response = await AuthService.sendEmailVerification(data);
          set({ isLoading: false });
          return response;
        } catch (error) {
          set({ isLoading: false });
          throw error;
        }
      },

      verifyEmail: async (data: VerifyEmailRequest) => {
        set({ isLoading: true });
        try {
          const response = await AuthService.verifyEmail(data);
          
          // Save to localStorage
          AuthService.setToken(response.token);
          AuthService.setUserData(response.user);
          
          set({
            user: response.user,
            token: response.token,
            isAuthenticated: true,
            isLoading: false,
          });
        } catch (error) {
          set({ isLoading: false });
          throw error;
        }
      },

      logout: () => {
        AuthService.logout();
        set({
          user: null,
          token: null,
          isAuthenticated: false,
          isLoading: false,
        });
      },

      updateProfile: async (data: UpdateProfileRequest) => {
        set({ isLoading: true });
        try {
          const updatedUser = await AuthService.updateProfile(data);
          
          // Update localStorage
          AuthService.setUserData(updatedUser);
          
          set({
            user: updatedUser,
            isLoading: false,
          });
        } catch (error) {
          set({ isLoading: false });
          throw error;
        }
      },

      setLoading: (loading: boolean) => {
        set({ isLoading: loading });
      },
    }),
    {
      name: 'auth-store',
      partialize: (state) => ({
        user: state.user,
        token: state.token,
        isAuthenticated: state.isAuthenticated,
      }),
      onRehydrateStorage: () => (state) => {
        // Validate stored token on app load
        if (state?.token && !AuthService.isAuthenticated()) {
          state.logout();
        }
      },
    }
  )
);

export default useAuthStore;