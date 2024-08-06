package constants

const (
	STORAGE_ENDPOINT_UPLOAD               = "/storage/upload"
	STORAGE_ENDPOINT_UPLOAD_STAGE         = "/storage/upload/stage"
	STORAGE_ENDPOINT_UPLOAD_COMMIT        = "/storage/upload/commit"
	STORAGE_ENDPOINT_UPLOAD_WITH_SECRET   = "/storage/secret/upload"
	STORAGE_ENDPOINT_DOWNLOAD             = "/storage/download/:file_id"
	STORAGE_ENDPOINT_DOWNLOAD_WITH_SECRET = "/storage/secret/download/:file_id"
	STORAGE_ENDPOINT_CREATE_SECRET        = "/storage/secret"
	STORAGE_ENDPOINT_RETRIEVE_SECRET      = "storage/secret/retrieve"
	STORAGE_ENDPOINT_RESET_PIN_CODE       = "/storage/secret/pin"
)
