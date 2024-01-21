# 进入后端目录并启动后端服务
cd backend
brew services start redis
go run main.go &
echo "backend started"

# 进入前端目录并启动前端服务
cd ../frontend
npm install
npm run dev &
echo "frontend started"

# 等待所有后台命令完成
wait