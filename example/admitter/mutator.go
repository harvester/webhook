package admitter

import (
	"github.com/sirupsen/logrus"
	admissionregv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/yaocw2020/webhook/pkg/types"
)

type poolMutator struct {
	types.DefaultMutator
}

var _ types.Mutator = &poolMutator{}

func NewPoolMutator() *poolMutator {
	return &poolMutator{}
}

func (v *poolMutator) Create(request *types.Request, newObj runtime.Object) (types.Patch, error) {
	pod := newObj.(*corev1.Pod)

	labels := pod.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}

	labels["example"] = "example"

	logrus.Info("add example label")

	return types.Patch{
		types.PatchOp{
			Op:    types.PatchOpReplace,
			Path:  "/metadata/labels",
			Value: labels,
		},
	}, nil
}

func (v *poolMutator) Resource() types.Resource {
	return types.Resource{
		Names:      []string{"pods"},
		Scope:      admissionregv1.NamespacedScope,
		APIGroup:   corev1.SchemeGroupVersion.Group,
		APIVersion: corev1.SchemeGroupVersion.Version,
		ObjectType: &corev1.Pod{},
		OperationTypes: []admissionregv1.OperationType{
			admissionregv1.Create,
		},
	}
}
