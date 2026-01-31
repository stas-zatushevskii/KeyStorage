package logger

import (
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Log singleton value
var Log *zap.Logger
var once sync.Once

func GetLogger() {
	once.Do(func() {
		stdout := zapcore.AddSync(os.Stdout)

		level := zap.NewAtomicLevelAt(zap.InfoLevel)
		file := zapcore.AddSync(&lumberjack.Logger{
			Filename:   "../../logs/app.log",
			MaxSize:    10, // megabytes
			MaxBackups: 3,
			MaxAge:     7, // days
			Compress:   true,
		})

		productionCfg := zap.NewProductionEncoderConfig()
		productionCfg.TimeKey = "timestamp"
		productionCfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
		fileEncoder := zapcore.NewJSONEncoder(productionCfg)

		developmentCfg := zap.NewDevelopmentEncoderConfig()
		developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		developmentCfg.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")
		consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)

		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, stdout, level),
			zapcore.NewCore(fileEncoder, file, level))

		Log = zap.New(core)
	})
}

type (
	responseData struct {
		status int
		size   int
	}
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}


// WithLogging logs all requests and response from server
func WithLogging(c *chi.Mux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := &loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		c.ServeHTTP(lw, r)
		duration := time.Since(start)

		Log.Info("Handled HTTP request",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Int("status", responseData.status),
			zap.Duration("duration", duration),
			zap.Int("size", responseData.size),
		)
	})
}
