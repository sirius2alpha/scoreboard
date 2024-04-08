package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"scoreboard/app/backend/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 创建一个 WebSocket 升级器
// ReadBufferSize 和 WriteBufferSize 分别设置了读取和写入 WebSocket 连接的缓冲区大小
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 记录所有打开的 WebSocket 连接
var connections []*websocket.Conn

// 判断字符串是否为 JSON 格式
func isJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

// HandleWebSocket 是一个处理 WebSocket 连接的函数
// 它接收一个 gin.Context 对象，这个对象包含了 HTTP 请求的所有信息
func HandleWebSocket(c *gin.Context) {
	// 使用 WebSocket 升级器将 HTTP GET 请求升级为 WebSocket 连接
	// c.Writer 和 c.Request 分别是 HTTP 的响应写入器和请求对象
	// 如果升级失败，打印错误信息并返回
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("failed to upgrade GET request: %v", err)
		return
	}

	connections = append(connections, ws)

	// 使用 defer 语句确保 WebSocket 连接在函数返回时关闭
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			log.Printf("error closing WebSocket connection: %v", err)
		}
		// remove the connection from the list
		for i, conn := range connections {
			if conn == ws {
				connections = append(connections[:i], connections[i+1:]...)
				break
			}
		}
	}(ws)

	// 进入一个无限循环，持续监听 WebSocket 连接上的消息
	for {
		// 使用 ws.ReadMessage() 函数读取新消息
		// 这个函数会返回一个消息类型、一个包含消息内容的字节切片bytes[]和一个错误对象
		_, message, err := ws.ReadMessage()
		// 如果读取消息时出错，处理错误
		if err != nil {
			// 如果错误是因为 WebSocket 连接已经正常关闭，就退出循环
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				break
			}
			// 如果是其他错误，打印错误信息并退出循环
			log.Printf("error reading message: %v", err)
			break
		}

		// 定义一个 map，用于存储解析后的消息
		msg := make(map[string]string)

		// 检查message是否为json格式
		if !isJSON(string(message)) {
			log.Printf("WARNING: The essage is not json: %v", message)
			break
		}

		// 使用 json.Unmarshal 函数解析消息
		// message 是一个包含 JSON 数据的字节切片
		// &msg 是一个指向 msg 的指针，json.Unmarshal 会将解析后的数据填充到这个 map 中
		err = json.Unmarshal(message, &msg)
		if err != nil {
			return
		}

		// 从解析后的消息中获取 "type" 字段的值，并将其转换为 string 类型
		// 这里假设 "type" 字段的值是一个字符串，表示消息的类型
		messageType := msg["type"]
		log.Printf("Received a new message, and the type is " + msg["type"])

		// 使用 switch 语句根据消息类型进行不同的处理
		switch messageType {
		case "NewUser":
			// 当消息类型为 "NewUser" 时，处理新用户逻辑
			// 从解析后的消息中获取 "nickname" 字段的值，并将其转换为 string 类型
			// 然后调用 services.AddnewUser 函数添加新用户
			services.AddNewUser(msg["nickname"])
			log.Printf("New user added: " + msg["nickname"])

		case "UserClick":
			// 当消息类型为 "UserClick" 时，处理用户点击逻辑
			// 从解析后的消息中获取 "nickname" 字段的值，并将其转换为 string 类型
			// 然后调用 services.HandleUserClick 函数处理用户点击
			services.HandleUserClick(msg["nickname"])
			log.Printf("User " + msg["nickname"] + " clicked")

		case "UserInactive":
			// 当消息类型为 "UserInactive" 时，处理用户不活跃逻辑
			// 从解析后的消息中获取 "nickname" 字段的值，并将其转换为 string 类型
			// 然后调用 services.HandleUserInactive 函数处理用户不活跃
			services.HandleUserInactive(msg["nickname"])
			log.Printf("User inactive: " + msg["nickname"])
		}
	}
}

func init() {
	// 使用 go 关键字启动一个新的 goroutine
	go func() {
		// 使用无限循环，使得这个 goroutine 会一直运行
		for {
			// 每隔一分钟执行一次循环体中的代码
			time.Sleep(200 * time.Millisecond) // 可根据需要调整时间间隔

			// 调用 service.UpdateClickInterval
			services.UpdateClickInterval()
			// 调用 services.CheckAllUsers() 函数检查所有用户的活跃状态
			services.CheckAllUsers()
			// 调用 services.GetRanking() 函数获取用户排名
			ranking, err := services.GetRanking()

			// 如果在获取用户排名时发生错误，打印错误信息并跳过本次循环的剩余部分
			if err != nil {
				log.Printf("error getting ranking: %v", err)
				continue
			}

			// 使用 json.Marshal 函数将用户排名转换为 JSON 格式
			message, err := json.Marshal(ranking)
			// 如果在转换为 JSON 格式时发生错误，打印错误信息并跳过本次循环的剩余部分
			if err != nil {
				log.Printf("error marshalling ranking: %v", err)
				continue
			}

			// 遍历所有打开的 WebSocket 连接
			for _, conn := range connections {
				// 使用 conn.WriteMessage 方法向每个连接发送包含用户排名的消息
				// 如果在发送消息时发生错误，打印错误信息
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					log.Printf("error writing message: %v", err)
				}
			}
		}
	}()
}
