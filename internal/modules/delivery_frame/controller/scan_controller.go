package controller

import (
	"app/internal/modules/delivery_frame/service"
	"app/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScanController struct {
	service service.IScanService
}

func NewScanController(service service.IScanService) *ScanController {
	return &ScanController{service: service}
}

// CreateScan godoc
// @Summary Create a new scan
// @Description Create a new scan record and return the scan details
// @Tags delivery-frame
// @Accept json
// @Produce json
// @Success 200 {object} dto.ScanResponseDto
// @Failure 500 {object} response.Response
// @Router /delivery-frame/scan [post]
func (c *ScanController) CreateScan(ctx *gin.Context) {
	result := c.service.CreateScan()
	response.HandleServiceResult(ctx, result)
}

// UploadImage godoc
// @Summary Upload an image for a scan
// @Description Upload an image to MinIO and update the scan record
// @Tags delivery-frame
// @Accept multipart/form-data
// @Produce json
// @Param device_id formData string true "Device ID"
// @Param scan_id formData string true "Scan ID"
// @Param image formData file true "Image File"
// @Success 200 {object} map[string]string
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /delivery-frame/upload [post]
func (c *ScanController) UploadImage(ctx *gin.Context) {
	deviceID := ctx.PostForm("device_id")
	scanID := ctx.PostForm("scan_id")
	file, err := ctx.FormFile("image")
	if err != nil {
		response.DataDetailResponse(ctx, http.StatusBadRequest, response.ErrCodeInvalidParams, nil)
		return
	}

	if deviceID == "" || scanID == "" {
		response.DataDetailResponse(ctx, http.StatusBadRequest, response.ErrCodeInvalidParams, nil)
		return
	}

	result := c.service.UploadImage(ctx.Request.Context(), deviceID, scanID, file)
	response.HandleServiceResult(ctx, result)
}

// GetImages godoc
// @Summary Get images for a scan
// @Description Retrieve a list of image URLs for a specific scan
// @Tags delivery-frame
// @Accept json
// @Produce json
// @Param scan_id query string true "Scan ID"
// @Success 200 {object} map[string][]string
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /delivery-frame/images [get]
func (c *ScanController) GetImages(ctx *gin.Context) {
	scanID := ctx.Query("scan_id")
	if scanID == "" {
		response.DataDetailResponse(ctx, http.StatusBadRequest, response.ErrCodeInvalidParams, nil)
		return
	}

	result := c.service.GetImages(ctx.Request.Context(), scanID)
	response.HandleServiceResult(ctx, result)
}

// DeleteFolder godoc
// @Summary Delete scan folder
// @Description Delete the folder associated with a scan from MinIO
// @Tags delivery-frame
// @Accept json
// @Produce json
// @Param scan_id query string true "Scan ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /delivery-frame/folder [delete]
func (c *ScanController) DeleteFolder(ctx *gin.Context) {
	scanID := ctx.Query("scan_id")
	if scanID == "" {
		response.DataDetailResponse(ctx, http.StatusBadRequest, response.ErrCodeInvalidParams, nil)
		return
	}

	result := c.service.DeleteScanFolder(ctx.Request.Context(), scanID)
	response.HandleServiceResult(ctx, result)
}
