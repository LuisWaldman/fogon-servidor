set FOGON_PUERTO=8080    
set FOGON_DB=mongodb://localhost:27017/
docker build -t fogon-servidor .
docker run --name fogonserver-container -p 8080:8080 fogon-servidor