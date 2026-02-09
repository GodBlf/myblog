package database

import (
	"fmt"
	"log"
	"myblog/util"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	ormlog "gorm.io/gorm/logger"
)

type StringContextKey string

var (
	blog_mysql      *gorm.DB
	blog_mysql_once sync.Once
	dblog           ormlog.Interface

	blog_redis      *redis.Client
	blog_redis_once sync.Once
)

func init() { //db log
	dblog = ormlog.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		ormlog.Config{
			SlowThreshold: 100 * time.Millisecond, // 慢 SQL 阈值
			LogLevel:      ormlog.Info,            // Log level, Silent 表示不输出日志
			Colorful:      true,                   // 禁用彩色打印 (注：代码中为 true，通常开启彩色)
		},
	)
}

func createMysqlDB(dbname, host, user, pass string, port int) *gorm.DB {
	// data source name 是 tester:123456@tcp(localhost:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname) // mb4兼容emoji表情符号
	var err error
	db, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{Logger: dblog, PrepareStmt: true}, //gorm config
	) // 启用 PrepareStmt, SQL预编译, 提高查询效率
	if err != nil {
		zap.L().Panic(("connect to mysql db failed"), zap.Error(err), zap.String("dsn", dsn)) // panic() os.Exit(2)
	}
	// 设置数据库连接池参数, 提高并发性能
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(100) // 设置数据库连接池最大连接数
	sqlDB.SetMaxIdleConns(20)  // 连接池最大允许的空闲连接数, 如果没有sql任务需要执行的连接数大于20, 超过的连接会被连接池关闭。
	zap.L().Info("connect to mysql db", zap.String("dbname", dbname))
	return db
}

func GetBlogDBConnection() *gorm.DB { //链接myblog数据库
	//懒汉初始化使用sync.Once保证只执行一次初始化数据库连接的操作, 避免重复创建连接池导致资源浪费和性能问题
	blog_mysql_once.Do(func() {
		dbName := "myblog"
		viper := util.CreateConfig("mysql")
		host := viper.GetString(dbName + ".host")
		port := viper.GetInt(dbName + ".port")
		user := viper.GetString(dbName + ".user")
		pass := viper.GetString(dbName + ".password")
		blog_mysql = createMysqlDB(dbName, host, user, pass, port)
	})
	return blog_mysql
}

func InitBlogDBConnection() *gorm.DB { //初始化链接myblog数据库的变量
	//懒汉初始化,使用sync.Once保证只执行一次初始化数据库连接的操作, 避免重复创建连接池导致资源浪费和性能问题
	blog_mysql_once.Do(func() {
		dbName := "myblog"
		viper := util.CreateConfig("mysql")
		host := viper.GetString(dbName + ".host")
		port := viper.GetInt(dbName + ".port")
		user := viper.GetString(dbName + ".user")
		pass := viper.GetString(dbName + ".password")
		blog_mysql = createMysqlDB(dbName, host, user, pass, port)
	})
	return blog_mysql
}

// Redis
func InitRedisClient() *redis.Client {
	blog_redis_once.Do(func() {
		config := util.CreateConfig("redis")
		client := redis.NewClient(&redis.Options{
			Addr:     config.GetString("addr"),
			Password: config.GetString("password"),
			DB:       config.GetInt("db"),
		})
		result, err := client.Ping().Result()
		if err != nil {
			zap.L().Panic("failed to connect to redis", zap.Error(err))
			panic(err)
		}
		zap.L().Info("connected to redis", zap.String("result", result))
		blog_redis = client
	})
	return blog_redis
}
