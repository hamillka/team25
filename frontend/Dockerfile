# Базовый образ с Node.js (используем LTS версию)
FROM node:18-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем package.json и package-lock.json (или yarn.lock, если используете Yarn)
COPY package*.json ./

# Устанавливаем зависимости
RUN npm install

# Копируем весь исходный код
COPY . .

# Собираем приложение для продакшена
RUN npm run build

# Второй этап: образ для запуска
FROM node:18-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Устанавливаем serve для раздачи статических файлов
RUN npm install -g serve

# Копируем собранные файлы из предыдущего этапа
COPY --from=builder /app/build ./build

# Указываем порт, который будет открыт
EXPOSE 3000

# Запускаем приложение
CMD ["serve", "-s", "build", "-l", "3000"]