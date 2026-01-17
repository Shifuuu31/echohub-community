#!/bin/zsh

# Start Docker rootless daemon
echo "🚀 Starting Docker daemon (rootless) in background..."

# Set environment variables if not already set
export PATH=$HOME/bin:$PATH
export DOCKER_HOST=unix://$XDG_RUNTIME_DIR/docker.sock

# Start the rootless Docker daemon in background
nohup dockerd-rootless.sh > ~/docker-rootless.log 2>&1 &

echo "✅ Docker daemon started in background (log: ~/docker-rootless.log)"
echo "📋 Check status with: docker ps"
echo "📋 View logs with: tail -f ~/docker-rootless.log"
