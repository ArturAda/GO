## 0) Полезные команды
```bash
docker compose up -d --build   # собрать и запустить в фоне
docker compose ps              # статус
docker compose logs -f         # логи сервиса
docker compose exec balance_server sh   # войти внутрь контейнера
docker compose down            # остановить (том сохраняется)
docker compose down -v         # остановить и удалить тома (данные пропадут)
```

---

## 1) Режим `file` — персистентность в volume

### Пример запуска (чистый старт)
```bash
docker compose down -v               # обнуляем данные и останавливаем всё (можно и без флага -v, чтобы данные сохранились)
docker compose up -d --build         # запускаем в фоне (в .yml STORAGE_TYPE=file)
curl -i http://localhost:8080/live   # проверка что жив
```

### Примеры
```bash
# 1) live
curl -i http://localhost:8080/live

# 2) deposit alice 1000 (= 10.00)
curl -i -X POST http://localhost:8080/api/v1/deposit \
  -H 'Content-Type: application/json' \
  -d '{"user_id":"alice","amount":1000}'

# 3) deposit bob 500 (= 5.00)
curl -i -X POST http://localhost:8080/api/v1/deposit \
  -H 'Content-Type: application/json' \
  -d '{"user_id":"bob","amount":500}'

# 4) transfer alice -> bob 300 (= 3.00)
curl -i -X POST http://localhost:8080/api/v1/transfer \
  -H 'Content-Type: application/json' \
  -d '{"from":"alice","to":"bob","amount":300}'

# 5) get balance alice (ожидаем 700 = 7.00)
curl -i http://localhost:8080/api/v1/balance/alice

# 6) get balance bob (ожидаем 800 = 8.00)
curl -i http://localhost:8080/api/v1/balance/bob

# 7) посмотреть файл в контейнере (лежит в /files/balances.json)
docker compose exec balance_server sh -lc 'ls -l /files && echo && cat /files/balances.json || true'
```

### Проверка персистентности
```bash
docker compose down
docker compose up -d

# балансы должны сохраниться:
curl -s http://localhost:8080/api/v1/balance/alice
curl -s http://localhost:8080/api/v1/balance/bob
```

---

## 2) Режим `memory` — данные НЕ сохраняются

###
```bash
docker compose down
docker compose build

docker compose run --rm --service-ports \
  -e HTTP_ADDRESS=":8080" \
  -e STORAGE_TYPE="memory" \
  balance_server
```
- Останавливать: `Ctrl+C`, или из другого терминала `docker stop <ID>`.

В **другом терминале** выполняем проверки (порт слушается на `localhost:8080`).

### Примеры
```bash
# 1) live
curl -i http://localhost:8080/live

# 2) баланс несуществующей alice -> 404
curl -i http://localhost:8080/api/v1/balance/alice

# 3) deposit alice 200 (=2.00)
curl -i -X POST http://localhost:8080/api/v1/deposit \
  -H 'Content-Type: application/json' \
  -d '{"user_id":"alice","amount":200}'

# 4) попытка перевести больше, чем есть у alice -> 409
curl -i -X POST http://localhost:8080/api/v1/transfer \
  -H 'Content-Type: application/json' \
  -d '{"from":"alice","to":"bob","amount":9999}'

# 5) перевод самому себе -> 400
curl -i -X POST http://localhost:8080/api/v1/transfer \
  -H 'Content-Type: application/json' \
  -d '{"from":"alice","to":"alice","amount":10}'

# 6) некорректная сумма (0) -> 400
curl -i -X POST http://localhost:8080/api/v1/deposit \
  -H 'Content-Type: application/json' \
  -d '{"user_id":"charlie","amount":0}'
```