# Create directories
echo "Creating directories..."
mkdir -p configs internal/{domains,services,handlers,repositories} pkgs/{errs,logs,builderutil,cryptoutil,jwtutil} server/{middlewares,routes}
# Create files and add content
echo "Creating and populating files..."
# main.go
cat <<EOF > main.go
package main

import (
	"$PROJECT_NAME/configs"
	"$PROJECT_NAME/pkgs/db"
	"$PROJECT_NAME/pkgs/logs"
	"$PROJECT_NAME/server"

	"github.com/gofiber/fiber/v2"
)

func init() {
	configs.Init()
	configs.SetTimeZone()
	logs.Init()
	db.Init()
}

func main() {
	serve := server.New(fiber.New())
	serve.ListenAndServe()
}
EOF

# routes.go
cat <<EOF > server/server.go
package server

import (
	"fmt"
	"os"
	"os/signal"
	"$PROJECT_NAME/pkgs/builder"

	"github.com/gofiber/fiber/v2"
)

type server struct {
	app *fiber.App
}

func New(app *fiber.App) server {
	return server{app: app}
}

func (s server) routes() {
}

func (s server) gracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		serv := <-c
		if serv.String() == "interrupt" {
			fmt.Println("\nGracefully shutting down...")
			s.app.Shutdown()
		}
	}()
}

func (s server) ListenAndServe() {
	s.routes()
	s.gracefulShutdown()
	addr := builder.URLBuilder("addr")
	s.app.Listen(addr)
}
EOF

# configs.go
cat <<EOF > configs/configs.go
package configs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.ReadInConfig()
}

func SetTimeZone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	time.Local = ict
}
EOF

# internal/handlers/utils.go
cat <<EOF > internal/handlers/utils.go
package handlers

import (
	"$PROJECT_NAME/pkgs/errs"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func handleError(c *fiber.Ctx, err error) error {
	switch e := err.(type) {
	case errs.AppError:
		fmt.Fprintln(c, e)
		return c.SendStatus(e.Code)
	case error:
		fmt.Fprintln(c, e)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return nil
}
EOF

# pkgs/maps/data_mapping.go
cat <<EOF > pkgs/maps/data_mapping.go
package maps

import (
	"fmt"
	"reflect"
)

func Copy(src, dest any) error {
	srcValue := reflect.ValueOf(src)
	destValue := reflect.ValueOf(dest)

	// Ensure src is a struct and dest is a pointer to a struct
	if srcValue.Kind() != reflect.Struct {
		return fmt.Errorf("source must be a struct")
	}
	if destValue.Kind() != reflect.Ptr || destValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("destination must be a pointer to a struct")
	}

	destElem := destValue.Elem()

	// Map fields from src to dest, including embedded structs
	if err := mapFields(srcValue, destElem); err != nil {
		return err
	}

	return nil
}

func mapFields(srcValue, destElem reflect.Value) error {
	for i := 0; i < srcValue.NumField(); i++ {
		srcField := srcValue.Type().Field(i)
		srcFieldValue := srcValue.Field(i)

		// Handle embedded structs
		if srcField.Anonymous {
			if err := mapFields(srcFieldValue, destElem); err != nil {
				return err
			}
			continue
		}

		srcFieldName := srcField.Name

		// Find matching field in the destination
		destField := destElem.FieldByName(srcFieldName)
		if destField.IsValid() && destField.CanSet() && destField.Type() == srcFieldValue.Type() {
			destField.Set(srcFieldValue)
		}
	}
	return nil
}
EOF

# logs.go
cat <<EOF > pkgs/logs/logs.go
package logs

import (
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log     *zap.Logger
	config  zap.Config
	err     error
	Logware fiber.Handler
)

func Init() {
	if viper.GetString("server.mode") == "debug" {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.StacktraceKey = ""
	config.OutputPaths = viper.GetStringSlice("server.log")

	log, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatal(err.Error())
	}

	Logware = fiberzap.New(fiberzap.Config{
		Logger: log,
	})
}

func Info(message string, field ...zapcore.Field) {
	log.Info(message, field...)
}

func Debug(message string, field ...zapcore.Field) {
	log.Debug(message, field...)
}

func Error(message interface{}, field ...zapcore.Field) {
	switch v := message.(type) {
	case error:
		log.Error(v.Error(), field...)
	case string:
		log.Error(v, field...)
	}
}
EOF

# errs.go
cat <<EOF > pkgs/errs/errs.go
package errs

type AppError struct {
	Code    int
	Message string
}

var _ error = AppError{}

func (ae AppError) Error() string {
	return ae.Message
}
EOF

# pkgs/builder/strbuilder.go
cat <<EOF > pkgs/builder/strbuilder.go
package builder

import (
	"fmt"

	"github.com/spf13/viper"
)

func URLBuilder(connType string) string {
	switch connType {
	case "addr":
		addr := ""
		if viper.GetString("server.mode") == "debug" {
			addr = "localhost:" + viper.GetString("server.port")
		} else {
			addr = ":" + viper.GetString("server.port")
		}
		return addr

	case "postgres":
		return fmt.Sprintf(
			"host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Bangkok",
			viper.GetString("db.host"),
			viper.GetString("db.username"),
			viper.GetString("db.password"),
			viper.GetString("db.name"),
			viper.GetInt("db.port"),
		)

	case "mysql":
		return fmt.Sprintf(
			"%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True",
			viper.GetString("db.username"),
			viper.GetString("db.password"),
			viper.GetString("db.host"),
			viper.GetInt("db.port"),
			viper.GetString("db.dbname"),
		)
	default:
		return ""
	}
}
EOF


# config.yml file
RANDOM_KEY=$(openssl rand -hex 32)
cat <<EOF > config.yml
server:
  port: 8888
  mode: debug
secret:
  token: $RANDOM_KEY
db:
  user:
  password:
  host:
  port:
  name:
EOF

# .gitignore
cat <<EOF > .gitignore
**/config.y*ml
**/gen.sh
**/docker-compose.y*ml
**/compose.y*ml
EOF