package errorhandler

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
)

type ErrorHanlder struct {
	mux       *runtime.ServeMux
	marshaler runtime.Marshaler
	w         http.ResponseWriter
	r         *http.Request
	err       error
}

func New(mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) *ErrorHanlder {
	return &ErrorHanlder{
		mux:       mux,
		marshaler: marshaler,
		w:         w,
		r:         r,
		err:       err,
	}
}

func CallHanlder(
	ctx context.Context,
	mux *runtime.ServeMux,
	marshaler runtime.Marshaler,
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	h := New(mux, marshaler, w, r, err)
	h.Handle(ctx)
}

func (e *ErrorHanlder) Handle(ctx context.Context) {
	s, ok := status.FromError(e.err)
	if !ok {
		http.Error(e.w, e.err.Error(), http.StatusInternalServerError)
		return
	}
	e.w.WriteHeader(runtime.HTTPStatusFromCode(s.Code()))

	_ = e.marshaler.NewEncoder(e.w).Encode(map[string]interface{}{
		"error": s.Message(),
		"code":  s.Code().String(),
	})
}
