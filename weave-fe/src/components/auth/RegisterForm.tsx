import React from 'react';
import { useForm } from 'react-hook-form';
import { Link, useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';
import toast from 'react-hot-toast';
import styled from 'styled-components';
import { RegisterRequest } from '@/types/auth';
import { useAuthStore } from '@/store/authStore';

const FormContainer = styled(motion.div)`
  max-width: 400px;
  margin: 0 auto;
  padding: 2rem;
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
`;

const Title = styled.h1`
  font-size: 2rem;
  font-weight: bold;
  text-align: center;
  margin-bottom: 0.5rem;
  color: #1f2937;
`;

const Subtitle = styled.p`
  text-align: center;
  color: #6b7280;
  margin-bottom: 2rem;
`;

const Form = styled.form`
  display: flex;
  flex-direction: column;
  gap: 1rem;
`;

const InputGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
`;

const Label = styled.label`
  font-weight: 500;
  color: #374151;
`;

const Input = styled.input<{ hasError?: boolean }>`
  padding: 0.75rem;
  border: 2px solid ${props => props.hasError ? '#ef4444' : '#d1d5db'};
  border-radius: 8px;
  font-size: 1rem;
  transition: border-color 0.2s;

  &:focus {
    outline: none;
    border-color: ${props => props.hasError ? '#ef4444' : '#3b82f6'};
  }
`;

const ErrorText = styled.span`
  color: #ef4444;
  font-size: 0.875rem;
`;

const Button = styled(motion.button)<{ isLoading?: boolean }>`
  background: ${props => props.isLoading ? '#9ca3af' : '#10b981'};
  color: white;
  padding: 0.75rem;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 500;
  cursor: ${props => props.isLoading ? 'not-allowed' : 'pointer'};
  transition: background-color 0.2s;

  &:hover {
    background: ${props => props.isLoading ? '#9ca3af' : '#059669'};
  }
`;

const LinkText = styled.p`
  text-align: center;
  margin-top: 1.5rem;
  color: #6b7280;

  a {
    color: #3b82f6;
    text-decoration: none;
    font-weight: 500;

    &:hover {
      text-decoration: underline;
    }
  }
`;

const PasswordHint = styled.div`
  font-size: 0.875rem;
  color: #6b7280;
  margin-top: 0.25rem;
`;

const RegisterForm: React.FC = () => {
  const navigate = useNavigate();
  const { register: registerUser, isLoading } = useAuthStore();
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<RegisterRequest & { confirmPassword: string }>();

  const watchPassword = watch('password');

  const onSubmit = async (data: RegisterRequest & { confirmPassword: string }) => {
    if (data.password !== data.confirmPassword) {
      toast.error('비밀번호가 일치하지 않습니다.');
      return;
    }

    try {
      await registerUser({
        username: data.username,
        email: data.email,
        password: data.password,
      });
      toast.success('회원가입이 완료되었습니다! 로그인해주세요.');
      navigate('/login');
    } catch (error: any) {
      toast.error(error.message || '회원가입에 실패했습니다.');
    }
  };

  return (
    <FormContainer
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5 }}
    >
      <Title>Weave 회원가입</Title>
      <Subtitle>당신의 창의적인 여정을 시작하세요</Subtitle>

      <Form onSubmit={handleSubmit(onSubmit)}>
        <InputGroup>
          <Label htmlFor="username">사용자명</Label>
          <Input
            id="username"
            type="text"
            hasError={!!errors.username}
            {...register('username', {
              required: '사용자명을 입력해주세요.',
              minLength: {
                value: 3,
                message: '사용자명은 최소 3자 이상이어야 합니다.',
              },
              maxLength: {
                value: 50,
                message: '사용자명은 최대 50자까지 가능합니다.',
              },
              pattern: {
                value: /^[a-zA-Z0-9_]+$/,
                message: '사용자명은 영문, 숫자, 언더스코어만 사용 가능합니다.',
              },
            })}
          />
          {errors.username && <ErrorText>{errors.username.message}</ErrorText>}
        </InputGroup>

        <InputGroup>
          <Label htmlFor="email">이메일</Label>
          <Input
            id="email"
            type="email"
            hasError={!!errors.email}
            {...register('email', {
              required: '이메일을 입력해주세요.',
              pattern: {
                value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                message: '올바른 이메일 형식을 입력해주세요.',
              },
            })}
          />
          {errors.email && <ErrorText>{errors.email.message}</ErrorText>}
        </InputGroup>

        <InputGroup>
          <Label htmlFor="password">비밀번호</Label>
          <Input
            id="password"
            type="password"
            hasError={!!errors.password}
            {...register('password', {
              required: '비밀번호를 입력해주세요.',
              minLength: {
                value: 8,
                message: '비밀번호는 최소 8자 이상이어야 합니다.',
              },
              pattern: {
                value: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]/,
                message: '비밀번호는 대소문자, 숫자, 특수문자를 포함해야 합니다.',
              },
            })}
          />
          {errors.password && <ErrorText>{errors.password.message}</ErrorText>}
          <PasswordHint>
            8자 이상, 대소문자, 숫자, 특수문자 포함
          </PasswordHint>
        </InputGroup>

        <InputGroup>
          <Label htmlFor="confirmPassword">비밀번호 확인</Label>
          <Input
            id="confirmPassword"
            type="password"
            hasError={!!errors.confirmPassword}
            {...register('confirmPassword', {
              required: '비밀번호 확인을 입력해주세요.',
              validate: (value) =>
                value === watchPassword || '비밀번호가 일치하지 않습니다.',
            })}
          />
          {errors.confirmPassword && (
            <ErrorText>{errors.confirmPassword.message}</ErrorText>
          )}
        </InputGroup>

        <Button
          type="submit"
          isLoading={isLoading}
          disabled={isLoading}
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
        >
          {isLoading ? '회원가입 중...' : '회원가입'}
        </Button>
      </Form>

      <LinkText>
        이미 계정이 있나요?{' '}
        <Link to="/login">로그인</Link>
      </LinkText>
    </FormContainer>
  );
};

export default RegisterForm;