package constants

const (
	FIELD_STORAGE_ID           = "id"
	FIELD_STORAGE_UUID         = "_id"
	FIELD_STORAGE_DOWNLOAD_URL = "download_url"
	FIELD_STORAGE_TYPE         = "type"
	FIELD_STORAGE_TOKEN        = "token"
	FIELD_STORAGE_LIFE_TIME    = "life_time"
	FIELD_STORAGE_EXT          = "ext"
	FIELD_STORAGE_CREATED_BY   = "created_by"
	FIELD_STORAGE_CREATED_AT   = "created_at"
)

var (
	STORAGE_TYPE_ALLOWED       = []string{"image", "video", "audio", "document", "other"}
	STORAGE_MEDIA_TYPE_ALLOWED = []string{"image", "video", "audio"}
	STORAGE_MEDIA_EXT_ALLOWED  = []string{"jpg", "jpeg", "png", "gif", "mp4", "mov", "avi", "mp3", "wav"}
)

/*
"pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx", "txt", "zip", "rar", "7z", "tar", "gz", "bz2", "xz", "iso", "apk", "exe", "msi", "deb", "rpm", "dmg", "pkg", "app", "jar", "war", "ear", "html", "css", "js", "json", "xml", "yaml", "yml", "toml", "ini", "conf", "cfg", "log", "md", "markdown", "csv", "tsv", "sql", "db", "dbf", "sqlite", "sqlite3", "mdb", "accdb", "psd", "ai", "sketch", "svg", "eps", "cdr", "ico", "webp", "tiff", "tif", "bmp", "webm", "flv", "ogg", "ogv", "ogm", "ogx", "3gp", "3g2", "wmv", "asf", "flac", "ape", "wma", "aac", "m4a", "m4r", "m4p", "m4b"
*/
