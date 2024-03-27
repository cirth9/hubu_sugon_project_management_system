chcp 65001
cd project-user
docker build -t project-user:latest .
cd ..
docker-compose up -d