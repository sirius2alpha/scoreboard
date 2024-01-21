cd ./frontend/
npm run build
scp -r dist root@sirius1y.top:~/scoreboard/frontend

cd ../backend/
go build -o main
ssh root@sirius1y.top "mkdir -p ~/scoreboard/backend/"
scp main root@sirius1y.top:~/scoreboard/backend/