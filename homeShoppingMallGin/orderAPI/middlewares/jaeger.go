package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"homeShoppingMallGin/orderAPI/settings"
	"log"
)

func JaegerTrace() func(c *gin.Context) {
	return func(c *gin.Context) {
		// jaeger
		cfg := jaegercfg.Configuration{
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans:           true,
				LocalAgentHostPort: fmt.Sprintf("%s:%d", settings.Conf.JaegerConfig.Host, settings.Conf.JaegerConfig.Port),
			},
			ServiceName: settings.Conf.JaegerConfig.Name,
		}
		tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
		if err != nil {
			log.Printf("Could not initialize jaeger tracer: %s", err.Error())
			return
		}
		opentracing.SetGlobalTracer(tracer)
		defer closer.Close()

		startSpan := opentracing.StartSpan(c.Request.URL.Path)
		defer startSpan.Finish()

		c.Set("tracer", tracer)
		c.Set("parentSpan", startSpan)
		c.Next()
	}
}
