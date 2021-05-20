package util

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	runtimeutil "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var logger = zap.New()

func Do(ctx context.Context, fn func(msg []byte) ([][]byte, error)) error {
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		in, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Error(err, "failed to marshal message")
			w.WriteHeader(500)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		msgs, err := fn(in)
		if err != nil {
			logger.Error(err, "failed execute")
			w.WriteHeader(500)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		for _, out := range msgs {
			resp, err := http.Post("http://localhost:3569/messages", "application/octet-stream", bytes.NewBuffer(out))
			if err != nil {
				logger.Error(err, "failed to post message")
				w.WriteHeader(500)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
			if resp.StatusCode != 200 {
				err := fmt.Errorf("failed to post message %s", resp.Status)
				logger.Error(err, "failed to post message")
				w.WriteHeader(500)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
			logger.V(6).Info("do", "in", string(in), "out", string(out))
		}
		w.WriteHeader(200)
	})

	go func() {
		defer runtimeutil.HandleCrash()
		if err := http.ListenAndServe(":8080", nil); err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()

	return nil
}
