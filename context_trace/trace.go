package context_trace

import (
	"bytes"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/context"
)

func genTrace() (string, error) {
	var buffer bytes.Buffer
	var t int64 = time.Now().UnixNano() / 1e6
	var r int64 = rand.Int63n(100000)
	if _, err := buffer.WriteString(strconv.FormatInt(t, 10)); err != nil {
		return "", err
	}
	if _, err := buffer.WriteString(strconv.FormatInt(r, 10)); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

// func GenCtx(ctx context.Context) context.Context {
// backgroundCtx := context.Background()
// md, ok := metadata.FromIncomingContext(ctx)
// if ctx.Value(trace) != nil {
// tempCtx := context.WithValue(backgroundCtx, trace, ctx.Value(trace).(string))
// return metadata.NewOutgoingContext(tempCtx, metadata.Pairs(trace, ctx.Value(trace).(string)))
// }
// if ok {
// if mds, ok := md[trace]; ok {
// if len(mds) != 0 {
// tempCtx := metadata.NewOutgoingContext(backgroundCtx, metadata.Pairs(trace, mds[0]))
// return context.WithValue(tempCtx, trace, mds[0])
// }
// }
// }
// logId := genLogId()
// tempCtx := metadata.NewOutgoingContext(backgroundCtx, metadata.Pairs(trace, logId))
// return context.WithValue(tempCtx, trace, logId)
// }

func GetCtxTrace(ctx context.Context) (context.Context, string, error) {
	if ctx.Value(trace) != nil {
		t, ok := ctx.Value(trace).(string)
		if ok {
			return ctx, t, nil
		}
	}

	t, err := genTrace()
	if err != nil {
		return ctx, "", err
	}
	return context.WithValue(ctx, trace, t), t, nil
}

func GetHeaderTrace(req *http.Request) (*http.Request, string, error) {
	if req == nil {
		return nil, "", nil
	}

	if t := req.Header.Get(trace); t != "" {
		return req, t, nil
	}

	t, err := genTrace()
	if err != nil {
		return req, "", err
	}

	return req.WithContext(context.WithValue(req.Context(), trace, t)), t, nil
}
