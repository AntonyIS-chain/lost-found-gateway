## * Overview of the API Gateway
### The API Gateway will handle:

* Authentication & Authorization (JWT, OAuth2)
* Rate Limiting (Prevent abuse)
* Logging & Monitoring
* Load Balancing (Forward requests to services)
* Request Transformation (Modify headers, query params)


## 🛠️ Tech Stack
* Framework: Gin (lightweight, fast)
* Authentication: JWT / OAuth2
* Rate Limiting: golang.org/x/time/rate
* Message Broker: Kafka / RabbitMQ (optional for async processing)
* Configuration Management: Viper
* Logging: Zap / Logrus
* Monitoring: Prometheus & Grafana
* Circuit Breaker: Resilience (e.g., Hystrix, gobreaker)
