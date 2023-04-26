package pod_test

import (
	"context"
	"crypto/tls"
	"os"
	"testing"

	"github.com/infraboard/mcube/client/negotiator"
	"github.com/infraboard/mcube/client/rest"
	"github.com/infraboard/mcube/logger/zap"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	ctx = context.Background()
)

// admission.Handler Webhook处理接口
// admission.Response
// admission.Request
func TestPodMutateHook(t *testing.T) {
	c := rest.NewRESTClient()
	c.SetBaseURL("https://localhost:9443")

	payload, err := os.ReadFile("./test/request.json")
	if err != nil {
		t.Fatal(err)
	}
	c.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	resp := make(map[string]admission.Response)
	err = c.Post("/mutate--v1-pod").
		Body(payload).
		Do(ctx).
		ContentType(negotiator.MIME_JSON).
		Into(&resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

func init() {
	// 设置日志模式
	zap.DevelopmentSetup()
}
