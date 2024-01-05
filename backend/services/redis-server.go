package services

// 导入需要的包，包括 Redis 客户端库、日志和格式化输出
import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
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
}

// 获取用户点击排行榜的函数，打印出点击次数最多的前 10 个用户
// 定义一个结构体来存储用户ID和点击次数
type UserScore struct {
	ID    string
	Score float64
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

	// 返回切片，切片中的元素顺序与它们在Redis中的顺序相同
	return result, nil
}

// 处理不活跃用户的函数，当用户不活跃时，将其从 sorted set 中删除
func HandleUserInactive(userID string) (err error) {
	_, err = client.ZScore("users", userID).Result()
	if err != nil {
		log.Printf("the user " + userID + " don't exist")
		// 如果用户不存在，就直接返回
		return err
	}
	// 从 sorted set 中删除用户
	client.ZRem("users", userID)
	return err
}
