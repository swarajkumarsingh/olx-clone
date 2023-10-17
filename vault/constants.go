package vault

const (
	KeyOwnershipXAPIKey   = "ownership_api_key"
	KeyAdminXAPIKey       = "admin_api_key"
	KeyBaseURL            = "base_url"
	KeyFrontendBaseURL    = "frontend_base_url"
	KeyHypervergeAPPKey   = "hyperverge_app_key"
	KeyHypervergeAPPIDKey = "hyperverge_app_id_key"
	KeyVaultEncryptionKey = "vault_encryption_key"
	KeyCORSEnabled        = "cors_enabled"
	KeyCORSOriginDomains  = "cors_allowed_origin"
)

var VaultKeys = []string{
	KeyOwnershipXAPIKey,
	KeyBaseURL,
	KeyFrontendBaseURL,
	KeyHypervergeAPPKey,
	KeyHypervergeAPPIDKey,
	KeyCORSEnabled,
	KeyCORSOriginDomains,
}
