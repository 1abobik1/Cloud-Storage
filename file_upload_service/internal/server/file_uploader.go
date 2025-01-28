package server

import (
	"context"
	"mime/multipart"

	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/domain/models"
	pb "github.com/1abobik1/proto-Cloud-Storage/gen/go/file_uploader"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FileUploaderServiceI interface {
	UploadFile(ctx context.Context, file multipart.File, fileInfo *models.FileModel) (int, error)
	GetFileInfo(ctx context.Context, FileId int) (*models.FileModel, error)
	DownloadFile(ctx context.Context, FileId int, fileModel models.FileModel) error
}

type FileUploaderServerAPI struct {
	pb.UnimplementedFileUploaderServer
	service FileUploaderServiceI
}

func RegisterFileUploaderServ(gRPC *grpc.Server, service FileUploaderServiceI) {
	pb.RegisterFileUploaderServer(gRPC, &FileUploaderServerAPI{service: service})
}

func (s *FileUploaderServerAPI) UploadFile(ctx context.Context, req *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
	if len(req.GetFileContent()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "FileContent is required")
	}

	if len(req.GetFileName()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "FileName is required")
	}

	if req.GetFileType() == pb.FileType_UNKNOWN {
		return nil, status.Error(codes.InvalidArgument, "Unrecognized file type")
	}


	return &pb.UploadFileResponse{
		FileId:  "12",
		Message: "created",
	}, nil
}

func (s *FileUploaderServerAPI) GetFileInfo(ctx context.Context, req *pb.GetFileInfoRequest) (*pb.GetFileInfoResponse, error) {
	if len(req.GetFileId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "FileId is required")
	}

	return &pb.GetFileInfoResponse{
		FileName:   "example.txt",
		FileType:   pb.FileType_TEXT,
		FileSize:   1024,
		UploadTime: "2025-01-27T12:00:00Z",
	}, nil
}

func (s *FileUploaderServerAPI) DownloadFile(req *pb.DownloadFileRequest, stream pb.FileUploader_DownloadFileServer) error {
	if len(req.GetFileId()) == 0 {
		return status.Error(codes.InvalidArgument, "FileId is required")
	}

	for i := 0; i < 3; i++ {
		if err := stream.Send(&pb.DownloadFileResponse{
			FileChunk: []byte("chunk"),
		}); err != nil {
			return err
		}
	}
	return nil
}
