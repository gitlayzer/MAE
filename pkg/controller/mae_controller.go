package controller

import (
	"encoding/json"
	"fmt"

	"github.com/gitlayzer/MAE/tools"

	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Mutate 实现 admission webhook 功能
func Mutate(a *admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	// 获取 AdmissionReview 请求对象
	req := a.Request

	// 定义一个 patches 数组，用于存储需要修改的镜像
	var patches []map[string]string

	// 分别兼容 Pod, Deployment, statefulSet, DaemonSet, replicaSet, job, cronJob 等资源创建请求
	switch req.Kind.Kind {
	case "Deployment":
		var deploy appsv1.Deployment
		if err := json.Unmarshal(req.Object.Raw, &deploy); err != nil {
			return &admissionv1.AdmissionResponse{
				Allowed: false,
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
		}
		patches = patchContainers(deploy.Spec.Template.Spec.Containers)
	case "StatefulSet":
		var statefulSet appsv1.StatefulSet
		if err := json.Unmarshal(req.Object.Raw, &statefulSet); err != nil {
			return &admissionv1.AdmissionResponse{
				Allowed: false,
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
		}
		patches = patchContainers(statefulSet.Spec.Template.Spec.Containers)
	case "DaemonSet":
		var daemonSet appsv1.DaemonSet
		if err := json.Unmarshal(req.Object.Raw, &daemonSet); err != nil {
			return &admissionv1.AdmissionResponse{
				Allowed: false,
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
		}
		patches = patchContainers(daemonSet.Spec.Template.Spec.Containers)
	case "ReplicaSet":
		var replicaSet appsv1.ReplicaSet
		if err := json.Unmarshal(req.Object.Raw, &replicaSet); err != nil {
			return &admissionv1.AdmissionResponse{
				Allowed: false,
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
		}
		patches = patchContainers(replicaSet.Spec.Template.Spec.Containers)
	default:
		return &admissionv1.AdmissionResponse{Allowed: true}
	}

	if len(patches) == 0 {
		return &admissionv1.AdmissionResponse{Allowed: true}
	}

	patchBytes, err := json.Marshal(patches)
	if err != nil {
		return &admissionv1.AdmissionResponse{
			Allowed: false,
			Result: &metav1.Status{
				Message: "Failed to marshal patches: " + err.Error(),
			},
		}
	}

	return &admissionv1.AdmissionResponse{
		Allowed:   true,
		Patch:     patchBytes,
		PatchType: func() *admissionv1.PatchType { pt := admissionv1.PatchTypeJSONPatch; return &pt }(),
	}
}

// patchContainers 遍历 containers
func patchContainers(containers []corev1.Container) []map[string]string {
	var patches []map[string]string

	for i, container := range containers {
		if tools.RegexMatched(container.Image) {
			newImage := "docker.cloudimages.asia" + "/" + container.Image
			patch := map[string]string{
				"op":    "replace",
				"path":  fmt.Sprintf("/spec/template/spec/containers/%d/image", i),
				"value": newImage,
			}
			patches = append(patches, patch)
		}
	}

	return patches
}
