package service

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	sharedDto "github.com/Cxons/unischedulebackend/internal/shared/dto"
	status "github.com/Cxons/unischedulebackend/pkg/statuscodes"
	storage_go "github.com/supabase-community/storage-go"
)







type SupabaseResponse = sharedDto.ResponseDto


type SupabaseService interface{
	InitClient()
	CreateSignedUrl(fileName string) (SupabaseResponse, string, error) 
}



type supabaseService struct {
	supabaseUrl string
	supabaseSecretKey string
	storageClient *storage_go.Client
	logger *slog.Logger
}


func NewSupabaseService(supabaseUrl string, supabaseSecretKey string,logger *slog.Logger)*supabaseService{
	return &supabaseService{
		supabaseUrl: supabaseUrl,
		supabaseSecretKey: supabaseSecretKey,
		logger: logger,
	}
}


func (sps *supabaseService) InitClient(){
	storageClient := storage_go.NewClient(sps.supabaseUrl,sps.supabaseSecretKey,nil)
	sps.storageClient = storageClient
	if sps.storageClient == nil {
		sps.logger.Error("failed to initialize Supabase storage client")
	}
}


func (sps *supabaseService) CreateSignedUrl(fileName string) (SupabaseResponse, string, error) {
    // Get current time
    currentTime := time.Now()

    // Convert timestamp to Unix string
    ts := strconv.FormatInt(currentTime.Unix(), 10)

    // Construct the full storage path
    path := "public/university_logos/" + ts + "-" + fileName

	sps.logger.Info("the path","path:",path)

    // Create the signed upload URL
    result, err := sps.storageClient.CreateSignedUploadUrl("Unischedule", path)
   if err != nil {
    if serr, ok := err.(*storage_go.StorageError); ok {
        sps.logger.Error("supabase storage error",
            "status", serr.Status,
            "message", serr.Message,
        )
    } else {
        sps.logger.Error("error creating signed upload url",
            "type", fmt.Sprintf("%T", err),
            "err", err,
        )
    }
    return SupabaseResponse{}, status.InternalServerError.Message, err
}

    // Return a standard response
    return SupabaseResponse{
        Message:           "The signed URL is generated successfully",
        Data:              sps.supabaseUrl + result.Url,
        StatusCode:        status.OK.Code,
        StatusCodeMessage: status.OK.Message,
    }, status.OK.Message, nil
}
