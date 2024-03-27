package metrics

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/application/observability/logger"
	"github.com/kodmain/thetiptop/api/internal/architecture/persistence"
	"github.com/kodmain/thetiptop/api/internal/architecture/serializers/prom"
	"github.com/kodmain/thetiptop/api/internal/domain/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/services"
)

var repository = persistence.NewMetricsRepository()
var service = services.NewMetricsService(repository)

// @Summary      Show request metrics.
// @Description  Retrieves data for the most frequent request.
// @Tags         Metrics
// @Accept       */*
// @Produce      text/plain
// @Success      200  {string}  entities.Metrics  "Statistics of the most frequent request"
// @Failure      404  {string}  nil           	  "No data available"
// @Router       /metrics/statistics [get]
// @Id           metrics.Statistics
func Statistics(c *fiber.Ctx) error {
	promBytes, err := getPrometheusMetrics()
	if err != nil {
		logger.Error(err)
		return sendPrometheusError(c)
	}

	return c.SendString(string(promBytes))
}

func getPrometheusMetrics() ([]byte, error) {
	allHits, maxHits, err := fetchMetricsData()
	if err != nil {
		return nil, err
	}

	hitsMeta, err := prom.NewMetricMeta("request_count", "Total number of requests for each request path.", "counter", allHits)
	if err != nil {
		return nil, err
	}

	maxMeta, err := prom.NewMetricMeta("request_max_hit", "Details of the request with the most hits.", "counter", maxHits)
	if err != nil {
		return nil, err
	}

	return prom.Marshal(hitsMeta, maxMeta)
}

func fetchMetricsData() (entities.Metrics, *entities.Metric, error) {
	allHits, err := service.GetAllRequestStats()
	if err != nil {
		return nil, nil, err
	}

	maxHits, err := service.GetMostFrequentRequest()
	if err != nil {
		return nil, nil, err
	}

	return allHits, maxHits, nil
}

func sendPrometheusError(c *fiber.Ctx) error {
	data := prom.NewSimpleMeta("request_error", "Indicates an error occurred while fetching statistics.", "counter", 1)
	promBytes, _ := prom.Marshal(data)
	return c.SendString(string(promBytes))
}

func Counter(c *fiber.Ctx) error {
	service.IncrementRequest(c.Method(), c.Path())
	return c.Next()
}
