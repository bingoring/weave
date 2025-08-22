package services

import (
	"fmt"
	"log"

	"weave-module/config"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	cfg    *config.Config
	dialer *gomail.Dialer
}

func NewEmailService(cfg *config.Config) *EmailService {
	// Create SMTP dialer
	var dialer *gomail.Dialer
	if cfg.External.Email.SMTPHost != "" {
		d := gomail.NewDialer(
			cfg.External.Email.SMTPHost,
			587, // Default SMTP port
			cfg.External.Email.SMTPUsername,
			cfg.External.Email.SMTPPassword,
		)
		dialer = d
	}

	return &EmailService{
		cfg:    cfg,
		dialer: dialer,
	}
}

// SendEmail sends an email using SMTP
func (s *EmailService) SendEmail(to, subject, body string, data map[string]string) error {
	if s.dialer == nil {
		log.Printf("Email service not configured, skipping email to %s", to)
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.External.Email.SMTPUsername)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Add custom headers if needed
	if data != nil {
		for key, value := range data {
			if key == "reply_to" {
				m.SetHeader("Reply-To", value)
			}
		}
	}

	if err := s.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email to %s: %w", to, err)
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}

// SendWelcomeEmail sends a welcome email to new users
func (s *EmailService) SendWelcomeEmail(to, username string) error {
	subject := "Welcome to Weave!"
	body := s.generateWelcomeEmailHTML(username)
	
	return s.SendEmail(to, subject, body, nil)
}

// SendPasswordResetEmail sends a password reset email
func (s *EmailService) SendPasswordResetEmail(to, resetToken string) error {
	subject := "Reset Your Weave Password"
	body := s.generatePasswordResetEmailHTML(resetToken)
	
	return s.SendEmail(to, subject, body, nil)
}

// SendVerificationEmail sends an email verification email
func (s *EmailService) SendVerificationEmail(to, verificationToken string) error {
	subject := "Verify Your Weave Account"
	body := s.generateVerificationEmailHTML(verificationToken)
	
	return s.SendEmail(to, subject, body, nil)
}

// SendNotificationEmail sends a notification email
func (s *EmailService) SendNotificationEmail(to, title, message string, data map[string]string) error {
	subject := fmt.Sprintf("Weave Notification: %s", title)
	body := s.generateNotificationEmailHTML(title, message, data)
	
	return s.SendEmail(to, subject, body, data)
}

// HTML email templates
func (s *EmailService) generateWelcomeEmailHTML(username string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Welcome to Weave</title>
    <style>
        body { font-family: Arial, sans-serif; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #3b82f6 0%%, #1d4ed8 100%%); color: white; padding: 30px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: white; padding: 30px; border: 1px solid #e5e7eb; border-top: none; border-radius: 0 0 8px 8px; }
        .button { display: inline-block; background: #3b82f6; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin: 20px 0; }
        .footer { text-align: center; margin-top: 30px; color: #6b7280; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to Weave! 🎉</h1>
            <p>The Living Document Network</p>
        </div>
        <div class="content">
            <h2>Hello %s!</h2>
            <p>Welcome to Weave, where your ideas become everyone's wisdom. We're excited to have you join our community of creators and collaborators.</p>
            
            <h3>What's Next?</h3>
            <ul>
                <li>📝 Create your first Weave</li>
                <li>🌟 Explore trending content</li>
                <li>👥 Follow interesting creators</li>
                <li>💡 Contribute to existing Weaves</li>
            </ul>
            
            <a href="https://weave.app/dashboard" class="button">Start Creating</a>
            
            <p>If you have any questions, feel free to reach out to our support team.</p>
            
            <p>Happy weaving!</p>
            <p>The Weave Team</p>
        </div>
        <div class="footer">
            <p>© 2024 Weave. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`, username)
}

func (s *EmailService) generatePasswordResetEmailHTML(resetToken string) string {
	resetURL := fmt.Sprintf("https://weave.app/reset-password?token=%s", resetToken)
	
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Reset Your Password</title>
    <style>
        body { font-family: Arial, sans-serif; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #ef4444; color: white; padding: 30px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: white; padding: 30px; border: 1px solid #e5e7eb; border-top: none; border-radius: 0 0 8px 8px; }
        .button { display: inline-block; background: #ef4444; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin: 20px 0; }
        .warning { background: #fef2f2; border: 1px solid #fecaca; padding: 15px; border-radius: 6px; margin: 15px 0; }
        .footer { text-align: center; margin-top: 30px; color: #6b7280; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔐 Password Reset Request</h1>
        </div>
        <div class="content">
            <h2>Reset Your Password</h2>
            <p>We received a request to reset your Weave account password. If you didn't make this request, you can safely ignore this email.</p>
            
            <div class="warning">
                <strong>⚠️ Security Notice:</strong> This reset link will expire in 1 hour for your security.
            </div>
            
            <p>Click the button below to reset your password:</p>
            
            <a href="%s" class="button">Reset Password</a>
            
            <p>Or copy and paste this link into your browser:</p>
            <p style="word-break: break-all; color: #6b7280;">%s</p>
            
            <p>If you're having trouble, contact our support team.</p>
            
            <p>Best regards,<br>The Weave Team</p>
        </div>
        <div class="footer">
            <p>© 2024 Weave. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`, resetURL, resetURL)
}

func (s *EmailService) generateVerificationEmailHTML(verificationToken string) string {
	verifyURL := fmt.Sprintf("https://weave.app/verify-email?token=%s", verificationToken)
	
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Verify Your Email</title>
    <style>
        body { font-family: Arial, sans-serif; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #10b981; color: white; padding: 30px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: white; padding: 30px; border: 1px solid #e5e7eb; border-top: none; border-radius: 0 0 8px 8px; }
        .button { display: inline-block; background: #10b981; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin: 20px 0; }
        .footer { text-align: center; margin-top: 30px; color: #6b7280; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>✅ Verify Your Email</h1>
        </div>
        <div class="content">
            <h2>Almost There!</h2>
            <p>Please verify your email address to complete your Weave account setup and start creating amazing content.</p>
            
            <a href="%s" class="button">Verify Email Address</a>
            
            <p>Or copy and paste this link into your browser:</p>
            <p style="word-break: break-all; color: #6b7280;">%s</p>
            
            <p>Once verified, you'll have full access to all Weave features.</p>
            
            <p>Welcome to the community!<br>The Weave Team</p>
        </div>
        <div class="footer">
            <p>© 2024 Weave. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`, verifyURL, verifyURL)
}

func (s *EmailService) generateNotificationEmailHTML(title, message string, data map[string]string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>%s</title>
    <style>
        body { font-family: Arial, sans-serif; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #6366f1; color: white; padding: 30px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: white; padding: 30px; border: 1px solid #e5e7eb; border-top: none; border-radius: 0 0 8px 8px; }
        .notification { background: #f8fafc; border-left: 4px solid #6366f1; padding: 20px; margin: 20px 0; }
        .footer { text-align: center; margin-top: 30px; color: #6b7280; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔔 %s</h1>
        </div>
        <div class="content">
            <div class="notification">
                <p>%s</p>
            </div>
            
            <p>Check out what's happening in your Weave community!</p>
            
            <p>Best regards,<br>The Weave Team</p>
        </div>
        <div class="footer">
            <p>© 2024 Weave. All rights reserved.</p>
            <p><a href="https://weave.app/settings">Manage email preferences</a></p>
        </div>
    </div>
</body>
</html>
`, title, title, message)
}