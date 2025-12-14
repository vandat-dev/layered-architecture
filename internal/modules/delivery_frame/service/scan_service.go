package service

import (
	"app/global"
	"app/internal/modules/delivery_frame/dto"
	"app/internal/modules/delivery_frame/model"
	"app/internal/modules/delivery_frame/repo"
	"app/internal/third_party/redis"
	"app/internal/third_party/s3"
	"app/pkg/response"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IScanService interface {
	CreateScan() *response.ServiceResult
	UploadImage(ctx context.Context, deviceID, scanID string, file *multipart.FileHeader) *response.ServiceResult
	GetImages(ctx context.Context, scanID string) *response.ServiceResult
	DeleteScanFolder(ctx context.Context, scanID string) *response.ServiceResult
	ProcessKafkaFrame(ctx context.Context, deviceID, scanID string, frameData []byte) error
}

type scanService struct {
	repo       repo.IScanRepository
	s3Provider *s3.S3Provider
	redis      *redis.RedisProvider
	localCache sync.Map
}

func NewScanService(repo repo.IScanRepository, s3Provider *s3.S3Provider, redis *redis.RedisProvider) IScanService {
	return &scanService{
		repo:       repo,
		s3Provider: s3Provider,
		redis:      redis,
	}
}

// 1. HELPERS & CACHE (Using Redis String - Optimized)
func (s *scanService) getScanFromCache(ctx context.Context, scanID string) (string, bool) {
	key := fmt.Sprintf("scan:%s", scanID)
	data, err := s.redis.Get(ctx, key)
	if err != nil || data == "" {
		return "", false
	}
	return data, true
}

func (s *scanService) saveScanToCache(ctx context.Context, scanID string, scan *model.Scan) {
	key := fmt.Sprintf("scan:%s", scanID)
	if data, err := json.Marshal(scan); err == nil {
		_ = s.redis.Set(ctx, key, data, 5*time.Minute)
	}
}

// Helper: Validity check flow (RAM -> Redis -> DB)
// scanID: string (used for cache lookup)
// UUID is parsed only when DB access is required
func (s *scanService) isScanValidLazy(ctx context.Context, scanID string) bool {
	// 1. Check RAM cache
	if _, ok := s.localCache.Load(scanID); ok {
		return true
	}

	// 2. Check Redis cache
	if _, ok := s.getScanFromCache(ctx, scanID); ok {
		s.localCache.Store(scanID, true)
		return true
	}

	// 3. Check DB (UUID parsing only here)
	scanUUID, err := uuid.Parse(scanID)
	if err != nil {
		return false // Invalid UUID string -> definitely not in DB
	}

	scan, err := s.repo.GetScanByID(scanUUID)
	if err != nil || scan == nil {
		return false
	}

	// 4. Valid -> save to Redis & RAM cache
	s.saveScanToCache(ctx, scanID, scan)
	s.localCache.Store(scanID, true)
	return true
}

// Helper: Update DB path (Redis String lock, UUID for DB update)
func (s *scanService) tryUpdateDBPath(ctx context.Context, scanID, deviceID string) {
	key := fmt.Sprintf("scan:db_saved:%s", scanID)

	// Redis lock using string key
	isFirstTime, err := s.redis.SetNX(ctx, key, "1", 24*time.Hour)

	if err == nil && isFirstTime {
		// Parse UUID only when update is required
		if scanUUID, err := uuid.Parse(scanID); err == nil {
			folderPath := fmt.Sprintf("%s/%s", deviceID, scanID)
			_ = s.repo.UpdateScanImagePath(scanUUID, folderPath)
		}
	}
}

// ProcessKafkaFrame 2. KAFKA METHODS (New Logic - String-based processing)
func (s *scanService) ProcessKafkaFrame(ctx context.Context, deviceID string, scanID string, frameData []byte) error {
	// 1. Validate scan (string-based, no UUID parsing yet)
	if !s.isScanValidLazy(ctx, scanID) {
		global.Logger.Warn("ProcessKafkaFrame: Invalid scanID", zap.String("scanID", scanID), zap.String("deviceID", deviceID))
		return fmt.Errorf("scan_id invalid or not found: %s", scanID)
	}

	// 2. Get sequence number from Redis
	seqKey := fmt.Sprintf("scan:seq:%s", scanID)
	seqID, err := s.redis.Incr(ctx, seqKey)
	if err != nil {
		global.Logger.Error("ProcessKafkaFrame: Failed to increment sequence", zap.Error(err), zap.String("scanID", scanID))
		seqID = time.Now().UnixNano()
	}

	// 3. Generate filename & async upload
	fileName := fmt.Sprintf("%06d.webp", seqID)
	objectKey := fmt.Sprintf("%s/%s/%s", deviceID, scanID, fileName)

	go func(data []byte, sID, dID, oKey string) {
		bgCtx := context.Background()

		// a. Upload to S3
		if _, err := s.s3Provider.UploadBytes(bgCtx, data, oKey); err != nil {
			global.Logger.Error("ProcessKafkaFrame: Async upload failed", zap.Error(err), zap.String("objectKey", oKey))
			return
		}

		// b. Update DB path if needed
		s.tryUpdateDBPath(bgCtx, sID, dID)

	}(frameData, scanID, deviceID, objectKey)

	return nil
}

// CreateScan 3. API METHODS
func (s *scanService) CreateScan() *response.ServiceResult {
	scan := &model.Scan{ID: uuid.New()}
	if err := s.repo.CreateScan(scan); err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}
	return response.NewServiceResult(dto.ScanResponseDto{ID: scan.ID})
}

func (s *scanService) UploadImage(ctx context.Context, deviceID, scanID string, file *multipart.FileHeader) *response.ServiceResult {
	// API layer should validate UUID format early
	scanUUID, err := uuid.Parse(scanID)
	if err != nil {
		global.Logger.Warn("UploadImage: Invalid scanID format", zap.String("scanID", scanID))
		return response.NewServiceErrorWithCode(400, response.ErrCodeInvalidParams)
	}

	if _, ok := s.getScanFromCache(ctx, scanID); !ok {
		scan, err := s.repo.GetScanByID(scanUUID)
		if err != nil {
			global.Logger.Warn("UploadImage: Scan not found in DB", zap.String("scanID", scanID))
			return response.NewServiceErrorWithCode(404, response.ErrCodeDataNotFound)
		}
		s.saveScanToCache(ctx, scanID, scan)
	}

	objectName := fmt.Sprintf("%s/%s/%s", deviceID, scanID, file.Filename)
	url, err := s.s3Provider.UploadFile(ctx, file, objectName)
	if err != nil {
		global.Logger.Error("UploadImage: Failed to upload file", zap.Error(err), zap.String("objectName", objectName))
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	folderPath := fmt.Sprintf("%s/%s", deviceID, scanID)
	if err := s.repo.UpdateScanImagePath(scanUUID, folderPath); err != nil {
		global.Logger.Error("UploadImage: Failed to update DB image path", zap.Error(err), zap.String("scanID", scanID))
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	global.Logger.Info("UploadImage: Success", zap.String("scanID", scanID), zap.String("url", url))
	return response.NewServiceResult(map[string]string{"url": url})
}

func (s *scanService) GetImages(ctx context.Context, scanID string) *response.ServiceResult {
	scanUUID, err := uuid.Parse(scanID)
	if err != nil {
		global.Logger.Warn("GetImages: Invalid scanID format", zap.String("scanID", scanID))
		return response.NewServiceErrorWithCode(400, response.ErrCodeInvalidParams)
	}

	scan, err := s.repo.GetScanByID(scanUUID)
	if err != nil {
		global.Logger.Warn("GetImages: Scan not found", zap.String("scanID", scanID))
		return response.NewServiceErrorWithCode(404, response.ErrCodeDataNotFound)
	}

	if scan.ImagePath == "" {
		return response.NewServiceResult([]string{})
	}

	prefix := scan.ImagePath
	if prefix[len(prefix)-1] != '/' {
		prefix += "/"
	}

	images, err := s.s3Provider.ListObjects(ctx, prefix)
	if err != nil {
		global.Logger.Error("GetImages: Failed to list objects", zap.Error(err), zap.String("prefix", prefix))
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	return response.NewServiceResult(images)
}

func (s *scanService) DeleteScanFolder(ctx context.Context, scanID string) *response.ServiceResult {
	scanUUID, err := uuid.Parse(scanID)
	if err != nil {
		global.Logger.Warn("DeleteScanFolder: Invalid scanID format", zap.String("scanID", scanID))
		return response.NewServiceErrorWithCode(400, response.ErrCodeInvalidParams)
	}

	scan, err := s.repo.GetScanByID(scanUUID)
	if err != nil {
		global.Logger.Warn("DeleteScanFolder: Scan not found", zap.String("scanID", scanID))
		return response.NewServiceErrorWithCode(404, response.ErrCodeDataNotFound)
	}

	if scan.ImagePath == "" {
		return response.NewServiceResult(nil)
	}

	if err := s.s3Provider.RemoveFolder(ctx, scan.ImagePath); err != nil {
		global.Logger.Error("DeleteScanFolder: Failed to remove folder", zap.Error(err), zap.String("path", scan.ImagePath))
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	global.Logger.Info("DeleteScanFolder: Success", zap.String("scanID", scanID))
	return response.NewServiceResult(nil)
}
