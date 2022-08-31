package grpc

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/sergalkin/go-url-shortener.git/internal/app/config"
	pb "github.com/sergalkin/go-url-shortener.git/internal/app/grpc/proto"
	"github.com/sergalkin/go-url-shortener.git/internal/app/middleware"
	"github.com/sergalkin/go-url-shortener.git/internal/app/service"
	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
	"github.com/sergalkin/go-url-shortener.git/internal/app/utils"
)

type server struct {
	pb.UnimplementedShortenerServer
	dbStorage       storage.DB
	internalService service.Internal
	shortenService  service.URLShorten
	expandService   service.URLExpand
}

// NewServer - creates new gRPC server.
func NewServer(db storage.DB, internal service.Internal, shortService service.URLShorten, expand service.URLExpand) *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterShortenerServer(
		s,
		&server{
			dbStorage: db, internalService: internal, shortenService: shortService, expandService: expand,
		},
	)
	return s
}

// Ping - checks connection with Database.
func (s server) Ping(ctx context.Context, in *pb.EmptyRequest) (*pb.PingResponse, error) {
	result := true

	if err := s.dbStorage.Ping(ctx); err != nil {
		result = false
	}

	return &pb.PingResponse{Ok: result}, nil
}

// Stats - returns amount of urls and users stored in database.
func (s server) Stats(ctx context.Context, in *pb.EmptyRequest) (*pb.StatsResponse, error) {
	urls, users, err := s.internalService.Stats()
	if err != nil {
		return &pb.StatsResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.StatsResponse{
		Urls:  int32(urls),
		Users: int32(users),
		Error: "",
	}, nil
}

// ShortenURL - receives in request long URL and returns in response short URL.
func (s server) ShortenURL(ctx context.Context, in *pb.ShortenURLRequest) (*pb.ShortenURLResponse, error) {
	response := pb.ShortenURLResponse{}

	uid, errID := getUserID(in.UserId)
	if errID != nil {
		return &pb.ShortenURLResponse{
			Error: errID.Error(),
		}, nil
	}
	response.UserId = uid

	shortURL, err := s.shortenService.ShortenURL(in.Url, uid)
	if err != nil {
		return &pb.ShortenURLResponse{Error: err.Error()}, nil
	}
	response.Result = shortURL

	return &response, nil
}

// ExpandURL - return original URL.
func (s server) ExpandURL(ctx context.Context, in *pb.ExpandURLRequest) (*pb.ExpandURLResponse, error) {
	original, err := s.expandService.ExpandURL(in.ShortUrl)
	if err != nil {
		return &pb.ExpandURLResponse{Error: err.Error()}, nil
	}

	return &pb.ExpandURLResponse{OriginalUrl: original}, nil
}

// GetUserURLs - return a list of records for specific user.
func (s server) GetUserURLs(ctx context.Context, in *pb.GetUserURLsRequest) (*pb.GetUserURLsResponse, error) {

	uid, errID := getUserID(in.UserId)
	if errID != nil {
		return &pb.GetUserURLsResponse{Error: errID.Error()}, nil
	}

	result, err := s.expandService.ExpandUserLinks(uid)
	if err != nil {
		return &pb.GetUserURLsResponse{Error: err.Error()}, nil
	}

	records := make([]*pb.GetUserURLsResponse_Record, 0, len(result))
	for _, rec := range result {
		records = append(records, &pb.GetUserURLsResponse_Record{
			ShortUrl:    fmt.Sprintf("%s/%s", config.BaseURL(), rec.ShortURL),
			OriginalUrl: rec.OriginalURL,
		})
	}

	return &pb.GetUserURLsResponse{Records: records}, nil
}

// BatchInsert - shortens a list of URLs.
func (s server) BatchInsert(ctx context.Context, in *pb.BatchInsertRequest) (*pb.BatchInsertResponse, error) {
	uid, errID := getUserID(in.UserId)
	if errID != nil {
		return &pb.BatchInsertResponse{Error: errID.Error()}, nil
	}

	if len(in.Records) == 0 {
		return &pb.BatchInsertResponse{}, nil
	}
	response := pb.BatchInsertResponse{UserId: uid}

	reqRecords := make([]storage.BatchRequest, len(in.Records))
	for i, record := range in.Records {
		reqRecords[i].CorrelationID = record.CorrelationId
		reqRecords[i].OriginalURL = record.Url
	}
	res, err := s.dbStorage.BatchInsert(reqRecords, uid)
	if err != nil {
		return &pb.BatchInsertResponse{Error: err.Error()}, nil
	}

	responseRecords := make([]*pb.BatchInsertResponse_Records, len(res))
	for i, record := range res {
		responseRecords[i] = &pb.BatchInsertResponse_Records{
			CorrelationId: record.CorrelationID,
			ShortUrl:      record.ShortURL,
		}
	}

	response.Records = responseRecords

	return &response, nil
}

// DeleteURLs - soft deletes a list of URLs.
func (s server) DeleteURLs(ctx context.Context, in *pb.DeleteURLsRequest) (*pb.DeleteURLsResponse, error) {
	if len(in.Keys) == 0 {
		return &pb.DeleteURLsResponse{}, nil
	}

	uid, errID := getUserID(in.UserId)
	if errID != nil {
		return &pb.DeleteURLsResponse{Error: errID.Error()}, nil
	}

	if err := s.dbStorage.SoftDeleteUserURLs(uid, in.Keys); err != nil {
		return &pb.DeleteURLsResponse{Error: err.Error()}, nil
	}

	return &pb.DeleteURLsResponse{}, nil
}

func getUserID(requestUserID string) (string, error) {
	var uid string

	if requestUserID != "" {
		err := utils.Decode(requestUserID, &uid)
		if err != nil {
			return "", utils.ErrGRPCWrongUserID
		}

		return uid, nil
	}

	uid = middleware.GetUUID()
	if uid == "" {
		uid = uuid.New().String()
	} else {
		err := utils.Decode(middleware.GetUUID(), &uid)
		if err != nil {
			return "", utils.ErrGRPCWrongUserID
		}
	}

	return uid, nil
}
