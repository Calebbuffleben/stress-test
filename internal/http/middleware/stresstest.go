package middleware

import (
    "context"
    "net/http"
    "strconv"
    "time"
    "strings"

    "stress-test/internal/config"
    "stress-test/internal/stress"
    "stress-test/internal/util"
)

type StressTestMiddleware struct {
    cfg     config.Config
    limiter *stress.Limiter
}

func NewStressTestMiddleware(cfg config.Config, l *stress.Limiter) *StressTestMiddleware {
    return &StressTestMiddleware{cfg: cfg, limiter: l}
}

func (m *StressTestMiddleware) Handler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        token := strings.TrimSpace(r.Header.Get(m.cfg.RateTokenHeader))

        var res stress.Result
        var err error

        if token != "" && m.cfg.RateTokenEnabled {
            res, err = m.limiter.AllowByToken(ctx, token)
        } else if m.cfg.RateIPEnabled {
            ip := util.GetClientIP(r)
            res, err = m.limiter.AllowByIP(ctx, ip)
        } else {
            // No limiting enabled
            next.ServeHTTP(w, r)
            return
        }
        if err != nil {
            http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
            return
        }
        if !res.Allowed {
            // Respond 429 with message and Retry-After seconds
            retrySecs := int(res.RetryAfter / time.Second)
            if retrySecs <= 0 {
                retrySecs = 1
            }
            w.Header().Set("Retry-After", strconv.Itoa(retrySecs))
            w.WriteHeader(http.StatusTooManyRequests)
            _, _ = w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
            return
        }
        next.ServeHTTP(w, r)
    })
}

func (m *StressTestMiddleware) AllowByToken(ctx context.Context, token string) (stress.Result, error) {
    return m.limiter.AllowByToken(ctx, token)
}


