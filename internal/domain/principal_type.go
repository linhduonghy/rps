package domain

type PrincipalType string

const (
	PrincipalTypeUser  PrincipalType = "USER"
	PrincipalTypeRole  PrincipalType = "ROLE"
	PrincipalTypeGroup PrincipalType = "GROUP"
	PrincipalTypeAny   PrincipalType = "*"
)
