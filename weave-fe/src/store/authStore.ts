import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import {
  User,
  AuthState,
  AuthActions,
  LoginRequest,
  RegisterRequest,
  UpdateProfileRequest,
} from '@/types/auth';
import AuthService from '@/services/authService';

interface AuthStore extends AuthState, AuthActions {}

export const useAuthStore = create<AuthStore>()(
  persist(
    (set, get) => ({
      // Initial state
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,

      // Actions
      login: async (credentials: LoginRequest) => {
        set({ isLoading: true });
        try {
          const response = await AuthService.login(credentials);
          
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

      register: async (data: RegisterRequest) => {
        set({ isLoading: true });
        try {
          const user = await AuthService.register(data);
          
          // After registration, user needs to login
          set({
            user: null,
            token: null,
            isAuthenticated: false,
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