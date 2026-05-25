package auth

import (
	"time"

	"github.com/besart951/code-links/apps/auth-service/backend/internal/domain"
)

type LoginProtection struct {
	EmailFailureLimit int
	IPFailureLimit    int
	Window            time.Duration
	Lockout           time.Duration
}

func DefaultLoginProtection() LoginProtection {
	return LoginProtection{
		EmailFailureLimit: 5,
		IPFailureLimit:    10,
		Window:            15 * time.Minute,
		Lockout:           15 * time.Minute,
	}
}

func (p LoginProtection) normalized() LoginProtection {
	defaults := DefaultLoginProtection()
	if p.EmailFailureLimit <= 0 {
		p.EmailFailureLimit = defaults.EmailFailureLimit
	}
	if p.IPFailureLimit <= 0 {
		p.IPFailureLimit = defaults.IPFailureLimit
	}
	if p.Window <= 0 {
		p.Window = defaults.Window
	}
	if p.Lockout <= 0 {
		p.Lockout = defaults.Lockout
	}
	return p
}

func (p LoginProtection) rateLimited(counts domain.LoginFailureCounts) bool {
	p = p.normalized()
	return counts.Email >= p.EmailFailureLimit || counts.IP >= p.IPFailureLimit
}
