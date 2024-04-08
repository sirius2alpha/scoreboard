package services

// 导入需要的包，包括 Redis 客户端库、日志和格式化输出
import (
	"fmt"
	"log"
	"scoreboard/app/backend/core"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

// 创建一个 Redis 客户端，连接到本地的 Redis 服务器
var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

type UserScore struct {
	ID            string
	Score         float64 // 存储在users中
	ClickTime     string  // 存储在clickTime中
	ClickInterval int64   // 存储在clickIntervala中
}

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	// 清空数据库中users
	if client.Exists("users").Val() != 0 {
		client.ZRemRangeByRank("users", 0, -1)
	}

	// 清空数据库中的clickTime
	if client.Exists("clickTime").Val() != 0 {
		client.Del("clickTime")
	}

	// 清空数据库中的clickInterval
	if client.Exists("clickInterval").Val() != 0 {
		client.Del("clickInterval")
	}
}

// 处理新用户的函数，当有新用户时，将其添加到users中
func AddNewUser(userID string) {
	client.ZAdd("users", redis.Z{
		Score:  0,      // 新用户的点击次数
		Member: userID, // 新用户的 ID
	})
}

// 处理用户点击的函数，当用户点击时，增加其在users中的点击次数
func HandleUserClick(userID string) {
	// 增加用户的点击次数
	client.ZIncrBy("users", 1, userID)

	// 记录点击时间
	clickTime := time.Now().Format("2006-01-02 15:04:05") // 获取当前时间，并格式化为字符串
	client.HSet("clickTime", userID, clickTime)           // 将点击时间存储到 Redis 中
}

// 更新用户点击间隔时间的函数
func UpdateClickInterval() {
	// 获取所有用户的 ID
	userIDs, _ := client.ZRange("users", 0, -1).Result()
	for _, userID := range userIDs {
		now := time.Now()
		lastClickTimeStr, err := client.HGet("clickTime", userID).Result() // 获取用户的上一次点击时间
		if err != nil {
			log.Printf("[Error] Getting last click time for user[%s] failed in UpdateClickInterval(): %v", userID, err)
			continue
		}

		// 设定时区为中国标准时间
		loc, _ := time.LoadLocation("Asia/Shanghai")
		// 将上一次点击时间的字符串转换为 time.Time 类型
		lastClickTime, err := time.ParseInLocation("2006-01-02 15:04:05", lastClickTimeStr, loc)
		if err != nil {
			log.Printf("[Error] Parsing last click time for user [%s] failed in UpdateClickInterval(): %v", userID, err)
			continue
		}

		// 计算时间间隔
		clickInterval := now.Sub(lastClickTime).Seconds()

		// 将时间间隔存储到 clickInterval 中
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

// 从Redis的users中获取点击次数最多的前 ranking_number 个用户
func GetRanking() ([]UserScore, error) {
	// 读取配置文件
	viper.AddConfigPath("../../conf/")
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	ranking_number := viper.GetInt64("ranking_number")

	ranking, err := client.ZRevRangeWithScores("users", 0, ranking_number-1).Result()
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
	}

	// 返回切片，切片中的元素顺序与它们在Redis中的顺序相同
	return result, nil
}

// 检查所有用户的活跃状态，如果用户不活跃，就将其从 users 中删除
func CheckAllUsers() {
	// 获取所有用户的 ID
	userIDs, _ := client.ZRange("users", 0, -1).Result()
	for _, userID := range userIDs {
		// 对每个用户调用 HandleUserInactive 函数
		HandleUserInactive(userID)
	}
}

// 处理不活跃用户的函数，当用户上一次点击间隔时间超过 max_keep_seconds 的时候就删除该用户
func HandleUserInactive(userID string) {

	max_keep_seconds := core.AppConfig.GetFloat64("redis.max_keep_seconds")

	// 获取用户的点击间隔时间
	clickIntervalStr, _ := client.HGet("clickInterval", userID).Result()

	// 将点击间隔时间的字符串转换为 float64 类型
	clickInterval, _ := strconv.ParseFloat(clickIntervalStr, 64)

	// 如果点击间隔时间超过 max_keep_seconds 秒，就从users中删除用户
	if clickInterval > max_keep_seconds {
		client.ZRem("users", userID)
	}
}
