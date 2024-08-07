package xtype

import "mime/multipart"

type File struct {
	multipart.FileHeader
}

func NewFile(fh *multipart.FileHeader) *File {
	return &File{
		FileHeader: *fh,
	}
}
