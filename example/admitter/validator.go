package admitter

import (
	"github.com/sirupsen/logrus"
	admissionregv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/yaocw2020/webhook/pkg/types"
)

type poolValidator struct {
	types.DefaultValidator
}

var _ types.Validator = &poolValidator{}

func NewPoolValidator() *poolValidator {
	return &poolValidator{}
}

func (v *poolValidator) Delete(request *types.Request, oldObj runtime.Object) error {
	pod := oldObj.(*corev1.Pod)
	logrus.Infof("delete pod %s/%s", pod.Namespace, pod.Name)
	return nil
}

func (v *poolValidator) Resource() types.Resource {
	return types.Resource{
		Names:      []string{"pods"},
		Scope:      admissionregv1.NamespacedScope,
		APIGroup:   corev1.SchemeGroupVersion.Group,
		APIVersion: corev1.SchemeGroupVersion.Version,
		ObjectType: &corev1.Pod{},
		OperationTypes: []admissionregv1.OperationType{
			admissionregv1.Delete,
		},
	}
}

func (v *poolValidator) ToAdmitter() types.Admitter {
	return types.Validator2Admitter(v)
}
