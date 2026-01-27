# KeyStorage

KeyStorage — консольное TUI-приложение для работы с пользовательскими данными и файлами.
Проект состоит из клиентского TUI и backend-сервиса с инфраструктурой в Docker.

---

## Requirements

- Go 1.21+
- Docker
- Docker Compose

---

## Getting started

### 1. Build TUI client

```bash
cd cmd/client
go build -o keystorage