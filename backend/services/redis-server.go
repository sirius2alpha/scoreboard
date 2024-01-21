package services

// 导入需要的包，包括 Redis 客户端库、日志和格式化输出
import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// 创建一个 Redis 客户端，连接到本地的 Redis 服务器
var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379", // Redis 服务器的地址
	Password: "",               // Redis 服务器的密码，这里没有设置密码
	DB:       0,                // 使用的 Redis 数据库编号，这里使用默认的数据库 0
})

// 初始化函数，会在包被导入时自动执行
func init() {
	// 如果没有创建过users的sorted set，就创建一个
	if client.Exists("users").Val() == 0 {
		// 创建一个 sorted set，用于存储用户的点击次数
		// 这里先添加了一个初始用户 "user:1"，点击次数为 0
		client.ZAdd("users", redis.Z{
			Score:  10,        // 用户的点击次数
			Member: "default", // 用户的 ID
		})
	} else {
		// 如果已经存在，就清空它
		client.ZRemRangeByRank("users", 0, -1)
	}

	// 初始化点击时间和点击间隔时间
	if client.Exists("clickTime").Val() == 0 {
		client.HSet("clickTime", "default", "")
	} else {
		// 如果已经存在，就清空它
		client.Del("clickTime")
	}
	if client.Exists("clickInterval").Val() == 0 {
		client.HSet("clickInterval", "default", "0")
	} else {
		// 如果已经存在，就清空它
		client.Del("clickInterval")
	}
}

// 处理新用户的函数，当有新用户时，将其添加到 sorted set 中
func AddNewUser(userID string) {
	// 将新用户添加到 sorted set 中，初始点击次数为 0
	client.ZAdd("users", redis.Z{
		Score:  0,      // 新用户的点击次数
		Member: userID, // 新用户的 ID
	})
}

// 处理用户点击的函数，当用户点击时，增加其在 sorted set 中的分数（即点击次数）
func HandleUserClick(userID string) {
	// 增加用户的点击次数
	client.ZIncrBy("users", 1, userID)

	// 记录点击时间
	clickTime := time.Now().Format("2006-01-02 15:04:05") // 获取当前时间，并格式化为字符串
	client.HSet("clickTime", userID, clickTime)           // 将点击时间存储到 Redis 中
}

func UpdateClickInterval() {
	// 获取所有用户的 ID
	userIDs, _ := client.ZRange("users", 0, -1).Result()
	for _, userID := range userIDs {
		// 获取当前时间
		now := time.Now()
		// 获取上一次点击的时间
		lastClickTimeStr, err := client.HGet("clickTime", userID).Result()
		if err != nil {
			// 处理错误
			log.Printf("error getting last click time for user %s: %v", userID, err)
			continue
		}

		// 设定时区为中国标准时间
		loc, _ := time.LoadLocation("Asia/Shanghai")
		// 将上一次点击时间的字符串转换为 time.Time 类型
		lastClickTime, err := time.ParseInLocation("2006-01-02 15:04:05", lastClickTimeStr, loc)
		if err != nil {
			// 处理错误
			log.Printf("error parsing last click time for user %s: %v", userID, err)
			continue
		}

		// 打印调试信息，之前发现时区不一致导致计算时间间隔出现了一些问题
		// log.Printf("User: %s, Now: %s, Last Click: %s", userID, now, lastClickTime)

		// 计算时间间隔
		clickInterval := now.Sub(lastClickTime).Seconds()

		// 将时间间隔存储到 Redis 中
		client.HSet("clickInterval", userID, strconv.FormatInt(int64(clickInterval), 10))
	}
}

// 获取用户点击间隔时间的函数
func GetClickInterval(userID string) (int64, error) {
	// 从 Redis 中获取用户的点击间隔时间
	clickIntervalStr, err := client.HGet("clickInterval", userID).Result()
	if err != nil {
		return 0, fmt.Errorf("error getting click interval: %v", err)
	}
	// 将点击间隔时间的字符串转换为 int64 类型
	clickInterval, _ := strconv.ParseInt(clickIntervalStr, 10, 64)
	// 返回点击间隔时间
	return clickInterval, nil
}

// 获取用户点击时间的函数
func GetClickTime(userID string) (string, error) {
	// 从 Redis 中获取用户的点击时间
	clickTime, err := client.HGet("clickTime", userID).Result()
	if err != nil {
		return "", fmt.Errorf("error getting click time: %v", err)
	}
	// 返回点击时间
	return clickTime, nil
}

// 获取用户点击排行榜的函数，打印出点击次数最多的前 10 个用户
// 定义一个结构体来存储用户ID和点击次数，还有用户的上一次点击时间和点击间隔时间
type UserScore struct {
	ID            string
	Score         float64
	ClickTime     string
	ClickInterval int64
}

func GetRanking() ([]UserScore, error) {
	// 从Redis的sorted set中获取点击次数最多的前10个用户
	ranking, err := client.ZRevRangeWithScores("users", 0, 9).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting ranking: %v", err)
	}

	// 创建一个切片来存储用户ID和点击次数
	result := make([]UserScore, len(ranking))
	for i, user := range ranking {
		// 将用户ID和点击次数存入切片
		result[i] = UserScore{ID: user.Member.(string), Score: user.Score}
	}

	// 在result中添加上用户的上一次点击时间和点击间隔时间
	for i, user := range result {
		// 获取用户的上一次点击时间
		lastClickTimeStr, _ := GetClickTime(user.ID)
		clickInterval, _ := GetClickInterval(user.ID)

		// 把他们添加到result中
		result[i].ClickTime = lastClickTimeStr
		result[i].ClickInterval = clickInterval
		//log.Print(result[i])
	}

	// 返回切片，切片中的元素顺序与它们在Redis中的顺序相同
	return result, nil
}

// 检查所有用户的活跃状态，如果用户不活跃，就将其从 sorted set 中删除
func CheckAllUsers() {
	// 获取所有用户的 ID
	userIDs, _ := client.ZRange("users", 0, -1).Result()
	for _, userID := range userIDs {
		// 对每个用户调用 HandleUserInactive 函数
		HandleUserInactive(userID)
	}
}

// 处理不活跃用户的函数，当用户上一次点击间隔时间超过20s的时候就删除该用户
func HandleUserInactive(userID string) {
	// 获取用户的点击间隔时间
	clickIntervalStr, _ := client.HGet("clickInterval", userID).Result()

	// 将点击间隔时间的字符串转换为 float64 类型
	clickInterval, _ := strconv.ParseFloat(clickIntervalStr, 64)

	// 如果点击间隔时间超过20秒，就从 sorted set 中删除用户
	if clickInterval > 20 {
		client.ZRem("users", userID)
	}
	return
}
