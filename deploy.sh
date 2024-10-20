# Build the Docker image
echo "Building the Docker image..."
docker compose build
# Run the Docker Compose setup with N peers
echo "Deploying the p2p network with $1 peers..."
docker compose up --scale peer=$1 -d  # Changed from docker-compose to docker compose

echo "Deployment complete! $1 peers are now running."