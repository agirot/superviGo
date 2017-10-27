# superviGo
Personal terminal dashboard supervisor (Sentry,redis,rabbitMQ,HTTP)

##Features
- Connector for RabbitMQ, Sentry 7, Redis, endpoint http
- Highlight data depending of your configuration

##Overview

![superviGo overview](https://raw.githubusercontent.com/agirot/superviGo/master/assets/screenshot.png)

## config.json

```json
{
  "rabbitmq": {
    "host": "localhost",
    "web_port": "15672",
    "user": "user",
    "password": "password",
    "queue": {
      "worker.low": {
        "alert_warning_message": 10000,
        "alert_critical_message": 50000
      },
      "worker.mail": {
        "alert_warning_message": 300,
        "alert_critical_message": 500
      },
      "online.task": {
        "alert_warning_message": 1000,
        "alert_critical_message": 5000
      }
    }
  },
  "sentry": {
    "host": "localhost",
    "organization_slug": "XXXXXX",
    "project_slug": "XXXXXX",
    "api_key": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
    "max_item": 6
  },
  "redis": {
    "host": "localhost:6379",
    "info": [
      "used_memory_human",
      "instantaneous_input_kbps",
      "instantaneous_output_kbps"
    ]
  },
  "http-status": {
    "order API": {
      "url": "https://localhost/order",
      "http_code_expected": 200
    },
    "status API": {
      "url": "https://localhost/api/v1/status",
      "basic_auth_user": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
      "basic_auth_password": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
      "http_code_expected": 200
    }
  }
}
```