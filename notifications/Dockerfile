# Используем официальный образ Go как базовый
FROM golang:1.22 as builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /server_app


# Копируем все остальные файлы проекта в контейнер
COPY . .

# Скачиваем все зависимости
RUN  go mod tidy

# Собираем приложение
RUN go build -o main

# Начинаем новую стадию сборки на основе минимального образа
FROM golang:1.22

# Добавляем исполняемый файл из первой стадии в корневую директорию контейнера
COPY --from=builder /server_app/main /main

# Открываем порт 1234
EXPOSE 1234

# Запускаем приложение
CMD ["/main"]