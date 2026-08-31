package server

import (
	"context"
	"testing"

	"github.com/go-openapi/testify/v2/require"
	"github.com/harvester/webhook/pkg/config"
	"github.com/harvester/webhook/pkg/server/admission"
	admissionregv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type mockValidator struct {
	admission.Validator
}

func (m *mockValidator) Resource() admission.Resource {
	return admission.Resource{
		Names:      []string{"pods"},
		APIGroup:   "",
		APIVersion: "v1",
		Scope:      "Namespace",
		ObjectType: &corev1.Pod{},
		MatchConditions: []admissionregv1.MatchCondition{
			{
				Name:       "mock-validator-match-condition-1",
				Expression: "self.metadata.name == 'singleton'",
			},
		},
	}
}
func TestWebhookGeneration(t *testing.T) {
	mv := &mockValidator{}
	newWebhook := &WebhookServer{
		server: server{
			context:    context.TODO(),
			restConfig: nil,
			name:       "test-webhook",
			options: &config.Options{
				Namespace: "default",
			},
		},
		validators: []admission.Validator{mv},
		ValidationFilter: &admission.ResourceFilters{
			NamespaceSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": "default",
				},
			},
		},
	}

	validatingWebhookConfiguration := newWebhook.validatingWebhookConfiguration()
	assert := require.New(t)
	assert.NotNil(validatingWebhookConfiguration)
	assert.Len(validatingWebhookConfiguration.Webhooks, 1)
	assert.Equal("validator.default.test-webhook", validatingWebhookConfiguration.Webhooks[0].Name)
	assert.NotNil(validatingWebhookConfiguration.Webhooks[0].NamespaceSelector)
	assert.Equal("default", validatingWebhookConfiguration.Webhooks[0].NamespaceSelector.MatchLabels["name"])
	assert.Len(validatingWebhookConfiguration.Webhooks[0].MatchConditions, 1)
}
