#!/bin/bash

# Blog Aggregator Deployment Script
echo "🚀 Starting Blog Aggregator Deployment..."

# Check if .env file exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp env.example .env
    echo "⚠️  Please edit .env file with your configuration before running again!"
    echo "   - Change JWT_SECRET to a strong random string"
    echo "   - Update database credentials if needed"
    exit 1
fi

# Build and start services
echo "🔨 Building and starting services..."
docker-compose down
docker-compose build --no-cache
docker-compose up -d

# Wait for services to be ready
echo "⏳ Waiting for services to start..."
sleep 10

# Check if services are running
echo "🔍 Checking service status..."
docker-compose ps

# Add sample feeds
echo "📰 Adding sample RSS feeds..."
sleep 5

# Function to add a feed
add_feed() {
    local title="$1"
    local url="$2"
    
    curl -X POST http://localhost:8080/feeds \
        -H "Content-Type: application/json" \
        -d "{\"title\": \"$title\", \"url\": \"$url\"}" \
        -s > /dev/null
}

# Add popular RSS feeds
echo "Adding sample feeds..."
add_feed "TechCrunch" "https://techcrunch.com/feed/"
add_feed "Hacker News" "https://hnrss.org/frontpage"
add_feed "Dev.to" "https://dev.to/feed"
add_feed "Go Blog" "https://blog.golang.org/feed.atom"
add_feed "Reddit Programming" "https://www.reddit.com/r/programming.rss"

echo "✅ Deployment complete!"
echo ""
echo "🌐 Your Blog Aggregator is running at: http://localhost:8080"
echo "📚 API Documentation: http://localhost:8080/swagger/index.html"
echo "💊 Health Check: http://localhost:8080/healthz"
echo ""
echo "📊 To view logs: docker-compose logs -f"
echo "🛑 To stop: docker-compose down"
echo ""
echo "🔑 Sample API calls:"
echo "   Register: curl -X POST http://localhost:8080/users/register -H 'Content-Type: application/json' -d '{\"username\":\"testuser\",\"email\":\"test@example.com\",\"password\":\"password123\"}'"
echo "   Login: curl -X POST http://localhost:8080/login -H 'Content-Type: application/json' -d '{\"username\":\"testuser\",\"password\":\"password123\"}'"
echo "   List Feeds: curl http://localhost:8080/feeds"
echo "   List Posts: curl http://localhost:8080/posts"
