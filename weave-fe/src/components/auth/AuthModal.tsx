import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { motion, AnimatePresence } from 'framer-motion';
import toast from 'react-hot-toast';
import styled from 'styled-components';
import { useAuthStore } from '../../store/authStore';
import { SendEmailVerificationRequest, VerifyEmailRequest } from '../../types/auth';

interface AuthModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess?: () => void;
}

const Overlay = styled(motion.div)`
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
`;

const ModalContainer = styled(motion.div)`
  background: white;
  border-radius: 16px;
  box-shadow: 0 20px 25px rgba(0, 0, 0, 0.15);
  width: 90%;
  max-width: 400px;
  max-height: 90vh;
  overflow-y: auto;
  position: relative;
`;

const CloseButton = styled.button`
  position: absolute;
  top: 16px;
  right: 16px;
  background: none;
  border: none;
  font-size: 24px;
  color: #6b7280;
  cursor: pointer;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  transition: all 0.2s;

  &:hover {
    background: #f3f4f6;
    color: #374151;
  }
`;

const Content = styled.div`
  padding: 2rem;
`;

const Title = styled.h2`
  font-size: 1.75rem;
  font-weight: bold;
  text-align: center;
  margin-bottom: 0.5rem;
  color: #1f2937;
`;

const Subtitle = styled.p`
  text-align: center;
  color: #6b7280;
  margin-bottom: 2rem;
  line-height: 1.5;
`;

const Form = styled.form`
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
`;

const InputGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
`;

const Label = styled.label`
  font-weight: 500;
  color: #374151;
  font-size: 0.875rem;
`;

const Input = styled.input<{ hasError?: boolean }>`
  padding: 0.875rem;
  border: 2px solid ${props => props.hasError ? '#ef4444' : '#e5e7eb'};
  border-radius: 8px;
  font-size: 1rem;
  transition: all 0.2s;
  background: #f9fafb;

  &:focus {
    outline: none;
    border-color: ${props => props.hasError ? '#ef4444' : '#3b82f6'};
    background: white;
    box-shadow: 0 0 0 3px ${props => props.hasError ? 'rgba(239, 68, 68, 0.1)' : 'rgba(59, 130, 246, 0.1)'};
  }
`;

const CodeInput = styled(Input)`
  text-align: center;
  letter-spacing: 0.2em;
  font-size: 1.5rem;
  font-weight: 600;
  font-family: monospace;
`;

const ErrorText = styled.span`
  color: #ef4444;
  font-size: 0.875rem;
  margin-top: 0.25rem;
`;

const Button = styled(motion.button)<{ isLoading?: boolean; variant?: 'primary' | 'secondary' }>`
  background: ${props => {
    if (props.isLoading) return '#9ca3af';
    if (props.variant === 'secondary') return '#f3f4f6';
    return '#3b82f6';
  }};
  color: ${props => props.variant === 'secondary' ? '#374151' : 'white'};
  padding: 0.875rem 1.5rem;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 500;
  cursor: ${props => props.isLoading ? 'not-allowed' : 'pointer'};
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;

  &:hover {
    background: ${props => {
      if (props.isLoading) return '#9ca3af';
      if (props.variant === 'secondary') return '#e5e7eb';
      return '#2563eb';
    }};
  }
`;

const Divider = styled.div`
  display: flex;
  align-items: center;
  margin: 1.5rem 0;

  &::before,
  &::after {
    content: '';
    flex: 1;
    height: 1px;
    background: #e5e7eb;
  }

  span {
    margin: 0 1rem;
    color: #6b7280;
    font-size: 0.875rem;
  }
`;

const GoogleButton = styled(Button)`
  background: white;
  color: #374151;
  border: 2px solid #e5e7eb;

  &:hover {
    background: #f9fafb;
    border-color: #d1d5db;
  }

  svg {
    width: 20px;
    height: 20px;
  }
`;

const BackButton = styled.button`
  background: none;
  border: none;
  color: #6b7280;
  font-size: 0.875rem;
  cursor: pointer;
  padding: 0.5rem 0;
  transition: color 0.2s;

  &:hover {
    color: #374151;
  }
`;


type AuthStep = 'email' | 'verification';

const AuthModal: React.FC<AuthModalProps> = ({ isOpen, onClose, onSuccess }) => {
  const [step, setStep] = useState<AuthStep>('email');
  const [email, setEmail] = useState('');
  const [timeLeft, setTimeLeft] = useState(0);
  const { sendEmailVerification, verifyEmail, isLoading } = useAuthStore();

  const emailForm = useForm<SendEmailVerificationRequest>();
  const codeForm = useForm<VerifyEmailRequest>();

  // Timer effect for resend countdown
  React.useEffect(() => {
    let interval: NodeJS.Timeout;
    if (timeLeft > 0) {
      interval = setInterval(() => {
        setTimeLeft((prev) => prev - 1);
      }, 1000);
    }
    return () => clearInterval(interval);
  }, [timeLeft]);

  const handleEmailSubmit = async (data: SendEmailVerificationRequest) => {
    try {
      const response = await sendEmailVerification(data);
      setEmail(data.email);
      setStep('verification');
      setTimeLeft(response.expires_in);
      toast.success('인증번호가 발송되었습니다!');
      
      // If in development and code is provided, show it
      if (response.code) {
        toast.success(`개발용 인증번호: ${response.code}`, { duration: 10000 });
      }
    } catch (error: any) {
      toast.error(error.message || '인증번호 발송에 실패했습니다.');
    }
  };

  const handleCodeSubmit = async (data: VerifyEmailRequest) => {
    try {
      await verifyEmail(data);
      toast.success('로그인에 성공했습니다!');
      onSuccess?.();
      handleClose();
    } catch (error: any) {
      toast.error(error.message || '인증에 실패했습니다.');
    }
  };

  const handleResendCode = async () => {
    if (timeLeft > 0) return;
    
    try {
      const response = await sendEmailVerification({ email });
      setTimeLeft(response.expires_in);
      toast.success('인증번호가 재발송되었습니다!');
      
      if (response.code) {
        toast.success(`개발용 인증번호: ${response.code}`, { duration: 10000 });
      }
    } catch (error: any) {
      toast.error(error.message || '인증번호 발송에 실패했습니다.');
    }
  };

  const handleGoogleLogin = () => {
    window.location.href = `${process.env.REACT_APP_API_URL}/v1/api/auth/google/login`;
  };

  const handleClose = () => {
    setStep('email');
    setEmail('');
    setTimeLeft(0);
    emailForm.reset();
    codeForm.reset();
    onClose();
  };

  const formatTime = (seconds: number) => {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
  };

  if (!isOpen) return null;

  return (
    <AnimatePresence>
      <Overlay
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        exit={{ opacity: 0 }}
        onClick={handleClose}
      >
        <ModalContainer
          initial={{ opacity: 0, scale: 0.95, y: 20 }}
          animate={{ opacity: 1, scale: 1, y: 0 }}
          exit={{ opacity: 0, scale: 0.95, y: 20 }}
          onClick={(e) => e.stopPropagation()}
        >
          <CloseButton onClick={handleClose}>×</CloseButton>
          
          <Content>
            {step === 'email' && (
              <motion.div
                key="email"
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
              >
                <Title>Weave에 로그인</Title>
                <Subtitle>
                  이메일 주소를 입력하면 인증번호를 보내드립니다.<br />
                  아이디어가 모두의 지혜로 엮이는 곳으로 오세요.
                </Subtitle>

                <Form onSubmit={emailForm.handleSubmit(handleEmailSubmit)}>
                  <InputGroup>
                    <Label htmlFor="email">이메일 주소</Label>
                    <Input
                      id="email"
                      type="email"
                      placeholder="example@gmail.com"
                      hasError={!!emailForm.formState.errors.email}
                      {...emailForm.register('email', {
                        required: '이메일을 입력해주세요.',
                        pattern: {
                          value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                          message: '올바른 이메일 형식을 입력해주세요.',
                        },
                      })}
                    />
                    {emailForm.formState.errors.email && (
                      <ErrorText>{emailForm.formState.errors.email.message}</ErrorText>
                    )}
                  </InputGroup>

                  <Button
                    type="submit"
                    isLoading={isLoading}
                    disabled={isLoading}
                    whileHover={{ scale: 1.02 }}
                    whileTap={{ scale: 0.98 }}
                  >
                    {isLoading ? '발송 중...' : '인증번호 받기'}
                  </Button>
                </Form>

                <Divider>
                  <span>또는</span>
                </Divider>

                <GoogleButton
                  type="button"
                  onClick={handleGoogleLogin}
                  whileHover={{ scale: 1.02 }}
                  whileTap={{ scale: 0.98 }}
                >
                  <svg viewBox="0 0 24 24">
                    <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                    <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                    <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                    <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
                  </svg>
                  Google로 로그인
                </GoogleButton>
              </motion.div>
            )}

            {step === 'verification' && (
              <motion.div
                key="verification"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: 20 }}
              >
                <Title>인증번호 입력</Title>
                <Subtitle>
                  <strong>{email}</strong>로 발송된<br />
                  6자리 인증번호를 입력해주세요.
                </Subtitle>

                <Form onSubmit={codeForm.handleSubmit(handleCodeSubmit)}>
                  <InputGroup>
                    <Label htmlFor="code">인증번호</Label>
                    <CodeInput
                      id="code"
                      type="text"
                      placeholder="000000"
                      maxLength={6}
                      hasError={!!codeForm.formState.errors.code}
                      {...codeForm.register('code', {
                        required: '인증번호를 입력해주세요.',
                        pattern: {
                          value: /^\d{6}$/,
                          message: '6자리 숫자를 입력해주세요.',
                        },
                      })}
                    />
                    {codeForm.formState.errors.code && (
                      <ErrorText>{codeForm.formState.errors.code.message}</ErrorText>
                    )}
                  </InputGroup>

                  <Button
                    type="submit"
                    isLoading={isLoading}
                    disabled={isLoading}
                    whileHover={{ scale: 1.02 }}
                    whileTap={{ scale: 0.98 }}
                  >
                    {isLoading ? '인증 중...' : '로그인'}
                  </Button>
                </Form>

                <div style={{ textAlign: 'center', marginTop: '1rem' }}>
                  {timeLeft > 0 ? (
                    <p style={{ color: '#6b7280', fontSize: '0.875rem' }}>
                      {formatTime(timeLeft)} 후 재발송 가능
                    </p>
                  ) : (
                    <Button
                      type="button"
                      variant="secondary"
                      onClick={handleResendCode}
                      disabled={isLoading}
                    >
                      인증번호 재발송
                    </Button>
                  )}
                </div>

                <div style={{ textAlign: 'center', marginTop: '1rem' }}>
                  <BackButton onClick={() => setStep('email')}>
                    ← 이메일 주소 변경
                  </BackButton>
                </div>
              </motion.div>
            )}
          </Content>
        </ModalContainer>
      </Overlay>
    </AnimatePresence>
  );
};

export default AuthModal;