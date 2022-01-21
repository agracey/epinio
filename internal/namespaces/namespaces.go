// Package namespaces encapsulates all the functionality around Epinio-controlled namespaces
// TODO: Consider moving this + the applications + the services packages under
// "models".
package namespaces

import (
	"context"
	"fmt"

	"github.com/epinio/epinio/helpers/kubernetes"
	"github.com/epinio/epinio/internal/duration"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Namespace represents an epinio-controlled namespace in the system
type Namespace struct {
	Name string
}

func List(ctx context.Context, kubeClient *kubernetes.Cluster) ([]Namespace, error) {
	listOptions := metav1.ListOptions{
		LabelSelector: kubernetes.EpinioNamespaceLabelKey + "=" + kubernetes.EpinioNamespaceLabelValue,
	}

	namespaceList, err := kubeClient.Kubectl.CoreV1().Namespaces().List(ctx, listOptions)
	if err != nil {
		return []Namespace{}, err
	}

	result := []Namespace{}
	for _, namespace := range namespaceList.Items {
		result = append(result, Namespace{Name: namespace.ObjectMeta.Name})
	}

	return result, nil
}

// Exists checks if the named epinio-controlled namespace exists or
// not, and returns an appropriate boolean flag
func Exists(ctx context.Context, kubeClient *kubernetes.Cluster, lookupNamespace string) (bool, error) {
	namespaces, err := List(ctx, kubeClient)
	if err != nil {
		return false, err
	}
	fmt.Printf("Lookup Namespace: %+v", lookupNamespace)
	for _, namespace := range namespaces {
		fmt.Printf("Namespace: %+v", namespace)
		if namespace.Name == lookupNamespace {
			return true, nil
		}
	}

	return false, nil
}

// Create generates a new epinio-controlled namespace, i.e. a kube
// namespace plus a service account.
func Create(ctx context.Context, kubeClient *kubernetes.Cluster, namespace string) error {
	if _, err := kubeClient.Kubectl.CoreV1().Namespaces().Create(
		ctx,
		&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
				Labels: map[string]string{
					"kubed-sync":                       "registry-creds", // Instruct kubed to copy image pull secrets over.
					kubernetes.EpinioNamespaceLabelKey: kubernetes.EpinioNamespaceLabelValue,
				},
				Annotations: map[string]string{
					"linkerd.io/inject": "enabled",
				},
			},
		},
		metav1.CreateOptions{},
	); err != nil {
		if apierrors.IsAlreadyExists(err) {
			return errors.Errorf("Namespace '%s' name cannot be used. Please try another name", namespace)
		}
		return err
	}

	if err := createServiceAccount(ctx, kubeClient, namespace); err != nil {
		return errors.Wrap(err, "failed to create a service account for apps")
	}

	return nil
}

// Delete destroys an epinio-controlled namespace, i.e. the associated
// kube namespace and service account.
func Delete(ctx context.Context, kubeClient *kubernetes.Cluster, namespace string) error {
	err := kubeClient.Kubectl.CoreV1().Namespaces().Delete(ctx, namespace, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return kubeClient.WaitForNamespaceMissing(ctx, nil, namespace, duration.ToNamespaceDeletion())
}

// createServiceAccount is a helper to `Create` which creates the
// service account applications pushed to the namespace need for
// permission handling.
func createServiceAccount(ctx context.Context, kubeClient *kubernetes.Cluster, targetNamespace string) error {
	automountServiceAccountToken := true
	_, err := kubeClient.Kubectl.CoreV1().ServiceAccounts(targetNamespace).Create(
		ctx,
		&corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name: targetNamespace,
			},
			ImagePullSecrets: []corev1.LocalObjectReference{
				{Name: "registry-creds"},
			},
			AutomountServiceAccountToken: &automountServiceAccountToken,
		}, metav1.CreateOptions{})

	return err
}
