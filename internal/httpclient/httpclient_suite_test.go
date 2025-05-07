package httpclient_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/pershin-daniil/life-log/internal/httpclient"
)

func TestHttpclient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Httpclient Suite")
}

var _ = Describe("HTTPClient", func() {
	var (
		server *httptest.Server
		client *httpclient.Client
		called bool
	)

	BeforeEach(func() {
		client = httpclient.New(slog.Default())
	})

	AfterEach(func() {
		if server != nil {
			server.Close()
		}
	})

	It("sends JSON and parses response", func() {
		called = false
		server = httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				Expect(r.Method).To(Equal(http.MethodPost))
				Expect(r.Header.Get("Content-Type")).To(Equal("application/json"))

				var body map[string]any
				err := json.NewDecoder(r.Body).Decode(&body)
				Expect(err).ToNot(HaveOccurred())
				Expect(body["message"]).To(Equal("hi"))

				w.WriteHeader(http.StatusOK)
				_, err = w.Write([]byte(`{"ok": true}`))
				Expect(err).ToNot(HaveOccurred())
			}),
		)

		type response struct {
			OK bool `json:"ok"`
		}
		var out response
		err := client.Do(
			context.Background(),
			http.MethodPost,
			server.URL,
			map[string]any{"message": "hi"},
			&out,
		)

		Expect(err).ToNot(HaveOccurred())
		Expect(out.OK).To(BeTrue())
		Expect(called).To(BeTrue())
	})
})
