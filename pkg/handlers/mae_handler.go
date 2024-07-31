package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gitlayzer/MAE/pkg/controller"

	admissionv1 "k8s.io/api/admission/v1"
)

func MaeHandler(c *gin.Context) {
	// AdmissionReview 对象
	var admissionReview admissionv1.AdmissionReview

	// 读取请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 确保请求体在读取后正确关闭
	defer c.Request.Body.Close()

	// 解析请求体为 AdmissionReview 对象
	if err = json.Unmarshal(body, &admissionReview); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用 mutate 函数, 得到 AdmissionResponse 对象
	admissionResponse := controller.Mutate(&admissionReview)

	// 将 AdmissionResponse 赋值给 AdmissionReview.Response
	admissionReview.Response = &admissionv1.AdmissionResponse{
		UID:       admissionReview.Request.UID,
		Allowed:   admissionResponse.Allowed,
		Patch:     admissionResponse.Patch,
		PatchType: admissionResponse.PatchType,
	}

	// 序列化 AdmissionReview 对象为 JSON 格式
	respBytes, err := json.Marshal(admissionReview)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 写入响应
	c.Data(http.StatusOK, "application/json", respBytes)
}
