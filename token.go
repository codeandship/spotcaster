package spotcaster

// {"access_token":"abc","token_type":"Bearer","expires_in":3600,"scope":"streaming ugc-image-upload user-read-email"}

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	ExpiresAt   int64  `json:"expires_at,omitempty"`
}
