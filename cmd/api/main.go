package main

import (
	"aws-Api-Go/models"
	"log"
	"net/http"
	"strconv"

	"aws-Api-Go/middleware"
	"os"

	"aws-Api-Go/mocks"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var users = mocks.Users

func getLogLevel(debug bool) zapcore.Level {
	if debug {
		return zapcore.DebugLevel
	}
	return zapcore.InfoLevel
}

func getTimeEncoder(development bool) zapcore.TimeEncoder {
	if development {
		return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		}
	}
	return zapcore.EpochMillisTimeEncoder
}

func initLogger() *zap.Logger {
	development := os.Getenv("DEVELOPPEMENT") == "1"
	debug := os.Getenv("DEBUG") == "1"

	var cfg zap.Config
	if development {
		cfg = zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(getLogLevel(debug))
		cfg.Encoding = "console"
		cfg.EncoderConfig = zap.NewDevelopmentEncoderConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg = zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(getLogLevel(debug))
		cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	}

	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = getTimeEncoder(development)

	opts := []zap.Option{
		zap.AddStacktrace(zap.ErrorLevel),
	}
	if !development {
		opts = append(opts, zap.WithCaller(false))
	}

	logger, err := cfg.Build(opts...)

	if err != nil {
		panic(err)
	}
	return logger
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	logger := initLogger()
	defer logger.Sync()

	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	r.Use(middleware.LogMiddleware(logger))

	r.GET("/users", func(c *gin.Context) {

		if len(users) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No users found"})
			return
		} else {
			logger.Debug("Fetched users",
				zap.Int("user_count", len(users)),
			)
		}
		c.JSON(http.StatusOK, users)
	})

	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		for _, user := range users {
			if user.ID == userID {
				c.JSON(http.StatusOK, user)
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	})

	r.POST("/users", func(c *gin.Context) {
		var newUser models.User
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newUser.ID = len(users) + 1
		users = append(users, newUser)
		c.JSON(http.StatusCreated, newUser)
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
