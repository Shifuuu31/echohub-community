#!/bin/zsh

# Build and run EchoHub Community Forum Docker container

echo "🔨 Building Docker image..."
docker build -t echohub-community-app .

if [ $? -eq 0 ]; then
    echo "✅ Docker image built successfully!"
    echo "🚀 Starting container..."
    docker run -p 8080:8080 echohub-community-app
else
    echo "❌ Docker build failed!"
    exit 1
fi
