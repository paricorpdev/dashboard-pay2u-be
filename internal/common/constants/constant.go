package constants

// Application constants
const (
	Response         = "response"
	AuthUserLocalKey = "user"
)

const (
	ResponseNameSuccess = "success"
	ResponseNameError   = "error"
)

// JWT constants
const (
	TokenTypeBearer       = "Bearer"
	PrefixJwtAccessToken  = "access_token"
	PrefixJwtRefreshToken = "refresh_token"
)

// Standard Error messages (Single Source of Truth)
const (
	ErrInvalidRequest          = "invalid request body"
	ErrInvalidRequestSignature = "invalid request signature"
	ErrDataNotFound            = "data not found"
	ErrUserNotFound            = "user not found"
	ErrUserAlreadyExists       = "user already exists"
	ErrEmailAlreadyExists      = "email already exists"
	ErrUserInactive            = "user is inactive"
	ErrRoleNotFound            = "role not found"
	ErrUnauthorized            = "Unauthorized user"
	ErrUsernameOrPassword      = "invalid username or password"
	ErrFailedGenerate          = "failed to generate token"
	ErrInternalServerError     = "internal server error"
	ErrForbiddenAccess         = "forbidden: missing required access"
)

// Status constants
const (
	StatusActive   = "ACTIVE"
	StatusInactive = "INACTIVE"
)

// Role codes
const (
	RoleSuperAdmin = "ROLE_SUPER_ADMIN"
	RoleAdmin      = "ROLE_ADMIN"
	RoleUser       = "ROLE_USER"
)

// Audit Log actions & entities
const (
	AuditActionLogin      = "LOGIN"
	AuditActionCreate     = "CREATE"
	AuditActionUpdate     = "UPDATE"
	AuditActionDelete     = "DELETE"
	AuditActionActivate   = "ACTIVATE"
	AuditActionDeactivate = "DEACTIVATE"
	AuditActionResetPass  = "RESET_PASSWORD"

	AuditEntityUser = "users"
	AuditEntityRole = "roles"
)
