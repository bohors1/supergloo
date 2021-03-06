package testutils

import (
	"log"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/gomega"
	"github.com/solo-io/go-utils/errors"
	"github.com/solo-io/go-utils/kubeutils"
	"github.com/solo-io/supergloo/cli/pkg/helpers/clients"
	kubev1 "k8s.io/api/core/v1"
	apiexts "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func WaitForNamespaceTeardown(ns string) {
	EventuallyWithOffset(1, func() (bool, error) {
		namespaces, err := MustKubeClient().CoreV1().Namespaces().List(v1.ListOptions{})
		if err != nil {
			// namespace is gone
			return false, err
		}
		for _, n := range namespaces.Items {
			if n.Name == ns {
				return false, nil
			}
		}
		return true, nil
	}, time.Second*180).Should(BeTrue())
}

func TeardownSuperGloo(kube kubernetes.Interface) {
	kube.CoreV1().Namespaces().Delete("supergloo-system", nil)
	clusterroles, err := kube.RbacV1beta1().ClusterRoles().List(metav1.ListOptions{})
	if err == nil {
		for _, cr := range clusterroles.Items {
			if strings.Contains(cr.Name, "supergloo") {
				kube.RbacV1beta1().ClusterRoles().Delete(cr.Name, nil)
			}
		}
	}
	clusterrolebindings, err := kube.RbacV1beta1().ClusterRoleBindings().List(metav1.ListOptions{})
	if err == nil {
		for _, cr := range clusterrolebindings.Items {
			if strings.Contains(cr.Name, "supergloo") {
				kube.RbacV1beta1().ClusterRoleBindings().Delete(cr.Name, nil)
			}
		}
	}
	webhooks, err := kube.AdmissionregistrationV1beta1().MutatingWebhookConfigurations().List(metav1.ListOptions{})
	if err == nil {
		for _, wh := range webhooks.Items {
			if strings.Contains(wh.Name, "supergloo") {
				kube.AdmissionregistrationV1beta1().MutatingWebhookConfigurations().Delete(wh.Name, nil)
			}
		}
	}

	cfg, err := kubeutils.GetConfig("", "")
	Expect(err).NotTo(HaveOccurred())

	exts, err := apiexts.NewForConfig(cfg)
	Expect(err).NotTo(HaveOccurred())

	crds, err := exts.ApiextensionsV1beta1().CustomResourceDefinitions().List(metav1.ListOptions{})
	if err == nil {
		for _, cr := range crds.Items {
			if strings.Contains(cr.Name, "supergloo") {
				exts.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(cr.Name, nil)
			}
		}
	}
}

// remove supergloo controller pod(s)
func DeleteSuperglooPods(kube kubernetes.Interface, superglooNamespace string) {
	podsToDelete := []string{
		"supergloo",
		// TODO(EItanya): add this back in once it's part of the helm chart
		// "mesh-discovery",
	}
	for _, pod := range podsToDelete {
		// wait until pod is gone
		Eventually(func() error {
			dep, err := kube.ExtensionsV1beta1().Deployments(superglooNamespace).Get(pod, metav1.GetOptions{})
			if err != nil {
				return err
			}
			dep.Spec.Replicas = proto.Int(0)
			_, err = kube.ExtensionsV1beta1().Deployments(superglooNamespace).Update(dep)
			if err != nil {
				return err
			}
			pods, err := kube.CoreV1().Pods(superglooNamespace).List(metav1.ListOptions{})
			if err != nil {
				return err
			}
			for _, p := range pods.Items {
				if strings.HasPrefix(p.Name, pod) {
					return errors.Errorf("%s pods still exist", pod)
				}
			}
			return nil
		}, time.Second*60).ShouldNot(HaveOccurred())
	}

}

func WaitUntilPodsRunning(timeout time.Duration, namespace string, podPrefixes ...string) error {
	pods := clients.MustKubeClient().CoreV1().Pods(namespace)
	getPodStatus := func(prefix string) (*kubev1.PodPhase, error) {
		list, err := pods.List(metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		for _, pod := range list.Items {
			if strings.HasPrefix(pod.Name, prefix) {
				return &pod.Status.Phase, nil
			}
		}
		return nil, errors.Errorf("pod with prefix %v not found", prefix)
	}
	failed := time.After(timeout)
	notYetRunning := make(map[string]kubev1.PodPhase)
	for {
		select {
		case <-failed:
			return errors.Errorf("timed out waiting for pods to come online: %v", notYetRunning)
		case <-time.After(time.Second / 2):
			notYetRunning = make(map[string]kubev1.PodPhase)
			for _, prefix := range podPrefixes {
				stat, err := getPodStatus(prefix)
				if err != nil {
					log.Printf("failed to get pod status: %v", err)
					notYetRunning[prefix] = kubev1.PodUnknown
					continue
				}
				if *stat != kubev1.PodRunning {
					notYetRunning[prefix] = *stat
				}
			}
			if len(notYetRunning) == 0 {
				return nil
			}
		}

	}
}
