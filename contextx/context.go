package contextx

import "context"

type key string

const (
	RequestIDKey       key = "request_id"
	CorrelationIDKey   key = "correlation_id"
	SubjectUUIDKey     key = "subject_uuid"
	SubjectFullNameKey key = "subject_full_name"
	SubjectEmailKey    key = "subject_email"
	RoleCodeKey        key = "role_code"
	PermissionCodesKey key = "permission_codes"
	AuthTypeKey        key = "auth_type"
	OverviewKey        key = "overview"
)

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, RequestIDKey, id)
}

func WithSubjectUUID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, SubjectUUIDKey, id)
}

func GetSubjectUUID(ctx context.Context) string {
	return getString(ctx, SubjectUUIDKey)
}

func GetSubjectUUIDNil(ctx context.Context) *string {
	return getStringNil(ctx, SubjectUUIDKey)
}

func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, CorrelationIDKey, id)
}

func WithSubjectFullName(ctx context.Context, fullName string) context.Context {
	return context.WithValue(ctx, SubjectFullNameKey, fullName)
}

func WithSubjectEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, SubjectEmailKey, email)
}

func GetSubjectFullName(ctx context.Context) string {
	return getString(ctx, SubjectFullNameKey)
}

func GetSubjectEmail(ctx context.Context) string {
	return getString(ctx, SubjectEmailKey)
}

func WithRoleCode(ctx context.Context, code string) context.Context {
	return context.WithValue(ctx, RoleCodeKey, code)
}

func WithPermissionCodes(ctx context.Context, codes map[string]bool) context.Context {
	return context.WithValue(ctx, PermissionCodesKey, codes)
}

func WithAuthType(ctx context.Context, authType string) context.Context {
	return context.WithValue(ctx, AuthTypeKey, authType)
}

func GetRequestID(ctx context.Context) string {
	return getString(ctx, RequestIDKey)
}

func GetCorrelationID(ctx context.Context) string {
	return getString(ctx, CorrelationIDKey)
}

func GetRoleCode(ctx context.Context) string {
	return getString(ctx, RoleCodeKey)
}

func GetPermissionCodes(ctx context.Context) map[string]bool {
	if v, ok := ctx.Value(PermissionCodesKey).(map[string]bool); ok {
		return v
	}
	return nil
}

func GetAuthType(ctx context.Context) string {
	return getString(ctx, AuthTypeKey)
}

func getString(ctx context.Context, k key) string {
	if v, ok := ctx.Value(k).(string); ok {
		return v
	}

	return ""
}

func getStringNil(ctx context.Context, k key) *string {
	if v, ok := ctx.Value(k).(string); ok {
		return &v
	}

	return nil
}

func getBool(ctx context.Context, k key) bool {
	if v, ok := ctx.Value(k).(bool); ok {
		return v
	}

	return false
}

func WithOverview(ctx context.Context, value bool) context.Context {
	return context.WithValue(ctx, OverviewKey, value)
}

func GetOverview(ctx context.Context) bool {
	return getBool(ctx, OverviewKey)
}
