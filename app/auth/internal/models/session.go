package models

type Session struct {
	ID           string
	UserID       string
	RefreshToken string
}

func (s *Session) Verify(refreshToken string) bool {
	return s.RefreshToken == refreshToken
}
