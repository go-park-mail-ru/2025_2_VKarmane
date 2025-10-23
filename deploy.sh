#!/bin/bash

# Это скрипт для деплоя сервиса на сервер
# Для запуска нужно установить Go и задать переменные окружения
# Запуск: ./deploy.sh
# Переменные окружения:
#   SERVER_USER - пользователь на сервере
#   SERVER_IP - IP сервера
#   SERVER_PATH - путь на сервере

# Устанавливаем конфиг
SERVER_IP="217.16.23.67"
SERVER_USER="${SERVER_USER:-ubuntu}"
SERVER_PATH="/home/ubuntu"
BINARY_NAME="vkarmane-api"
SSH_KEY="bin/2025-2-VKarmane-ZquJHDDn.pem"

# Проверяем, что JWT_SECRET задан явно и не является дефолтным
DEFAULT_SECRETS=("your-production-secret-key-change-this" "your-secret-key" "your-production-secret-key" "default-secret" "secret" "password" "123456" "test")

if [ -z "$JWT_SECRET" ]; then
    echo "❌ JWT_SECRET must be explicitly set!"
    echo "   Please set JWT_SECRET environment variable:"
    echo "   export JWT_SECRET='your-secure-secret-key'"
    echo "   or run: JWT_SECRET='your-secure-secret-key' ./deploy.sh"
    echo ""
    echo "   For security reasons, JWT_SECRET is required."
    exit 1
fi

# Проверяем, что JWT_SECRET не является одним из дефолтных значений
for default_secret in "${DEFAULT_SECRETS[@]}"; do
    if [ "$JWT_SECRET" = "$default_secret" ]; then
        echo "❌ JWT_SECRET cannot be a default/insecure value: '$JWT_SECRET'"
        echo "   Please set a secure, unique JWT_SECRET:"
        echo "   export JWT_SECRET='your-secure-unique-secret-key'"
        echo "   or run: JWT_SECRET='your-secure-unique-secret-key' ./deploy.sh"
        echo ""
        echo "   For security reasons, default JWT secrets are not allowed."
        exit 1
    fi
done

# Проверяем минимальную длину JWT_SECRET
if [ ${#JWT_SECRET} -lt 16 ]; then
    echo "❌ JWT_SECRET must be at least 16 characters long!"
    echo "   Current length: ${#JWT_SECRET} characters"
    echo "   Please set a longer, secure JWT_SECRET:"
    echo "   export JWT_SECRET='your-secure-secret-key-at-least-16-chars'"
    exit 1
fi

# Остальные переменные окружения
ENV="${ENV:-development}"
HOST="${HOST:-0.0.0.0}"
PORT="${PORT:-8080}"
LOG_LEVEL="${LOG_LEVEL:-info}"

# Начинаем запуск сервиса
echo "🚀 Deploying with environment variables..."
echo "🔐 JWT_SECRET is set (security check passed)"

if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

echo "📦 Building binary for Linux..."
GOOS=linux GOARCH=amd64 go build -o $BINARY_NAME cmd/api/main.go

if [ ! -f "$BINARY_NAME" ]; then
    echo "❌ Error: Failed to build binary"
    exit 1
fi

echo "✅ Binary built successfully"

echo "🔍 Checking server connection..."
if ! ssh -i $SSH_KEY -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o ConnectTimeout=10 $SERVER_USER@$SERVER_IP "echo 'Connection successful'" 2>/dev/null; then
    echo "❌ Cannot connect to server $SERVER_IP"
    exit 1
fi

echo "🛑 Stopping old process..."
ssh -i $SSH_KEY -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null $SERVER_USER@$SERVER_IP "pkill -f $BINARY_NAME || true"

echo "📤 Uploading binary..."
scp -i $SSH_KEY -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null $BINARY_NAME $SERVER_USER@$SERVER_IP:$SERVER_PATH/

ssh -i $SSH_KEY -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null $SERVER_USER@$SERVER_IP "mkdir -p $SERVER_PATH/logs"

ssh -i $SSH_KEY -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null $SERVER_USER@$SERVER_IP "chmod +x $SERVER_PATH/$BINARY_NAME"

echo "📝 Creating startup script on server..."
ssh -i $SSH_KEY -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null $SERVER_USER@$SERVER_IP "cat > $SERVER_PATH/start.sh << 'EOF'
#!/bin/bash
export JWT_SECRET='$JWT_SECRET'
export ENV='$ENV'
export HOST='$HOST'
export PORT='$PORT'
export LOG_LEVEL='$LOG_LEVEL'
nohup ./$BINARY_NAME > logs/app.log 2>&1 &
echo \$! > logs/pid.txt
EOF"

ssh -i $SSH_KEY -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null $SERVER_USER@$SERVER_IP "chmod +x $SERVER_PATH/start.sh"

echo "🚀 Starting application..."
ssh -i $SSH_KEY -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null $SERVER_USER@$SERVER_IP "cd $SERVER_PATH && ./start.sh"

# Ждем и проверяем работу сервиса
sleep 3

if ssh -i $SSH_KEY -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null $SERVER_USER@$SERVER_IP "pgrep -f $BINARY_NAME > /dev/null"; then
    echo "✅ Service started successfully!"
    echo "📍 API available at: http://$SERVER_IP:$PORT"
    echo "🌐 CORS configured for: http://$SERVER_IP:8000"
    echo ""
    echo "🔧 Environment variables:"
    echo "  JWT_SECRET: $JWT_SECRET"
    echo "  ENV: $ENV"
    echo "  HOST: $HOST"
    echo "  PORT: $PORT"
    echo "  LOG_LEVEL: $LOG_LEVEL"
    echo ""
    echo "🔍 Check logs: ssh -i $SSH_KEY ubuntu@$SERVER_IP 'tail -f $SERVER_PATH/logs/app.log'"
    echo "🛑 Stop service: ssh -i $SSH_KEY ubuntu@$SERVER_IP 'pkill -f $BINARY_NAME'"
    echo "🔄 Restart: ssh -i $SSH_KEY ubuntu@$SERVER_IP 'cd $SERVER_PATH && ./start.sh'"
else
    echo "❌ Failed to start service. Check logs:"
    ssh -i $SSH_KEY -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null $SERVER_USER@$SERVER_IP "tail -20 $SERVER_PATH/logs/app.log"
    exit 1
fi

# Удаляем локальный бинарник
rm -f $BINARY_NAME

echo "🎉 Deployment completed!"
