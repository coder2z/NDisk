package rpc

import (
	NFilePb "github.com/myxy99/ndisk/pkg/pb/nfile"
)

type Server struct{}

func (s Server) FileUpload(server NFilePb.NFileService_FileUploadServer) error {
	panic("implement me")
}

func (s Server) FileDownload(info *NFilePb.FileInfo, server NFilePb.NFileService_FileDownloadServer) error {
	panic("implement me")
}
