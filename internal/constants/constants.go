package constants

const (
	// User Roles
	RoleAdmin  = "admin"
	RoleUser   = "user"
	RoleGuest  = "guest"

	// Status
	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusDeleted  = "deleted"
)

// Error Messages
const (
	ErrUnauthorized     = "Unauthorized"
	ErrForbidden        = "Forbidden"
	ErrNotFound         = "Resource not found"
	ErrValidation       = "Validation error"
	ErrInternalServer   = "Internal server error"
	ErrInvalidEmail     = "Invalid email format"
	ErrPasswordTooShort = "Password must be at least 6 characters"
	ErrDuplicateEmail   = "Email already registered"
)

// Success Messages
const (
	SuccessRegister = "Register successfully"
	SuccessLogin    = "Login successfully"
	SuccessGetUser  = "Get user successfully"
	SuccessCreate   = "Create successfully"
	SuccessUpdate   = "Update successfully"
	SuccessDelete   = "Delete successfully"
	SuccessGetAll   = "Get all data successfully"
)
