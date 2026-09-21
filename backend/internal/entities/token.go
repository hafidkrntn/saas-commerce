package entities

// SessionToken holds the claims data used when generating JWT tokens.
type SessionToken struct {
	SessionId       string `json:"session_id"`
	UserId          string `json:"user_id"`
	UserGroupId     string `json:"user_group_id"`
	TenantId        string `json:"tenant_id"`
	Role            string `json:"role"`
	IsAdministrator bool   `json:"is_administrator"`
	Application     string `json:"application"`
}
