package constants

const (
	// Storage
	STORAGE_ENDPOINT_UPLOAD                    = "/storage/upload"
	STORAGE_ENDPOINT_UPLOAD_STAGE              = "/storage/upload/stage"
	STORAGE_ENDPOINT_UPLOAD_COMMIT             = "/storage/upload/commit"
	STORAGE_ENDPOINT_UPLOAD_WITH_SECRET        = "/storage/secret/upload"
	STORAGE_ENDPOINT_UPLOAD_STAGE_WITH_SECRET  = "/storage/secret/upload/stage"
	STORAGE_ENDPOINT_UPLOAD_COMMIT_WITH_SECRET = "/storage/secret/upload/commit"
	STORAGE_ENDPOINT_DOWNLOAD                  = "/storage/download/:file_id"
	STORAGE_ENDPOINT_REQUEST_DOWNLOAD          = "/storage/download/request/:file_id"
	STORAGE_ENDPOINT_CREATE_SECRET             = "/storage/secret"
	STORAGE_ENDPOINT_RETRIEVE_SECRET           = "/storage/secret/retrieve"
	STORAGE_ENDPOINT_RESET_PIN_CODE            = "/storage/secret/pin"

	// Share
	SHARE_ENDPOINT_DOWNLOAD = "/download/:file_id"
)
