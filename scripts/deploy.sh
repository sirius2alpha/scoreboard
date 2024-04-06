cd ../src/frontend/
npm run build
scp -r dist root@sirius1y.top:/var/www/scoreboard/frontend/

cd ../src/backend/
go build -o main
scp main root@sirius1y.top:/var/www/scoreboard/backend/