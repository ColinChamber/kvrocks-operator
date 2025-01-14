package sentinel

import (
	"github.com/RocksLabs/kvrocks-operator/pkg/resources"
)

func (h *KVRocksSentinelHandler) ensureKubernetes() error {
	cm := resources.NewSentinelConfigMap(h.instance)
	err := h.k8s.CreateOrUpdateConfigMap(cm)
	if err != nil {
		return err
	}
	service := resources.NewSentinelService(h.instance)
	if err = h.k8s.CreateIfNotExistsService(service); err != nil {
		return err
	}
	dep := resources.NewSentinelDeployment(h.instance)
	if err = h.k8s.CreateIfNotExistsDeployment(dep); err != nil {
		return err
	}
	dep, err = h.k8s.GetDeployment(h.key)
	if err != nil {
		return err
	}
	if dep.Status.ReadyReplicas != *dep.Spec.Replicas {
		h.log.Info("please wait deployment ready")
		h.requeue = true
		return nil
	}
	pods, err := h.k8s.ListDeploymentPods(h.key)
	if err != nil {
		return err
	}
	for _, pod := range pods.Items {
		h.pods = append(h.pods, pod.Status.PodIP)
	}
	h.log.Info("kubernetes resources ok")
	return nil
}
