package models

import "time"

type TokenUsage struct {
	TokenID       int32     `json:"token_id"`
	UserID        int32     `json:"user_id"`
	LastUsedAt    time.Time `json:"last_used_at"`
	UsageCount    int32     `json:"usage_count"`
	LastUsedIP    string    `json:"last_used_ip"`
	LastUserAgent string    `json:"last_user_agent"`
}
