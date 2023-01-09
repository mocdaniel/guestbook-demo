package main

import (
	"context"
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/gomodule/redigo/redis"
	_ "github.com/lib/pq"
	"github.com/mocdaniel/guestbook-app/app/internal/data"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

//go:embed migrations/*.sql
var migrations embed.FS

type databaseCfg struct {
	host         string
	port         int
	user         string
	password     string
	database     string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type redisCfg struct {
	host     string
	port     int
	password string
}

type config struct {
	port  int
	db    databaseCfg
	redis redisCfg
}

type application struct {
	config         *config
	logger         *log.Logger
	models         data.Models
	sessionManager *scs.SessionManager
	frontend       fs.FS
}

var sessionManager *scs.SessionManager

func main() {

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	cfg, err := parseConfig()
	if err != nil {
		logger.Fatal(err)
	}

	logger.Printf("Configuration parsed...")

	logger.Printf("Connecting to database...")
	db, err := openDB(cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Printf("Database connection pool established.")

	logger.Printf("Applying pending database schema upgrades...")

	migrationDBDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Fatal(err)
	}

	migrationSourceDriver, err := iofs.New(migrations, "migrations")
	if err != nil {
		logger.Fatal(err)
	}

	migrator, err := migrate.NewWithInstance("embedded", migrationSourceDriver, cfg.db.database, migrationDBDriver)
	if err != nil {
		logger.Fatal(err)
	}

	err = migrator.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			logger.Printf("No pending database schema upgrades found.")
		} else {
			logger.Fatal(err)
		}
	} else {
		logger.Printf("Successfully applied database schema upgrades.")
	}

	logger.Printf("Connecting to redis cache...")

	var dial func() (redis.Conn, error)
	if cfg.redis.password == "" {
		dial = func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%v:%v", cfg.redis.host, cfg.redis.port))
		}
	} else {
		dial = func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%v:%v", cfg.redis.host, cfg.redis.port), redis.DialPassword(cfg.redis.password))
		}
	}

	pool := &redis.Pool{
		MaxIdle: 10,
		Dial:    dial,
	}

	conn := pool.Get()
	err = conn.Err()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Printf("Redis connection seems fine...")

	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Store = redisstore.New(pool)

	fe, err := getFrontendAssets()
	if err != nil {
		logger.Fatal(err)
	}

	app := &application{
		config:         cfg,
		logger:         logger,
		models:         data.NewModels(db),
		sessionManager: sessionManager,
		frontend:       fe,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      sessionManager.LoadAndSave(app.routes()),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Starting server listening on %s", srv.Addr)

	err = srv.ListenAndServe()
	logger.Fatal(err)
}

func openDB(cfg *config, logger *log.Logger) (*sql.DB, error) {
	dsn := "postgres://" + cfg.db.user + ":" + cfg.db.password + "@" + cfg.db.host + ":" + fmt.Sprint(cfg.db.port) + "/" + cfg.db.database + "?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func parseConfig() (*config, error) {
	// Viper set defaults
	viper.SetDefault("port", 8080)
	viper.SetDefault("db-port", 5432)
	viper.SetDefault("db-max-open-conns", 25)
	viper.SetDefault("db-max-idle-conns", 25)
	viper.SetDefault("db-host", "localhost")
	viper.SetDefault("db-user", "guestbook")
	viper.SetDefault("db-password", "password")
	viper.SetDefault("db-name", "guestbook")
	viper.SetDefault("db-max-idle-time", "15m")
	viper.SetDefault("redis-host", "localhost")
	viper.SetDefault("redis-port", 6379)
	viper.SetDefault("redis-password", "")

	// Define command line options
	flag.Int("port", 8080, "Webserver port")
	flag.Int("db-port", 5432, "PostgreSQL port")
	flag.Int("db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.Int("db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.Int("redis-port", 6379, "Redis port")
	flag.String("db-host", "localhost", "PostgreSQL server address")
	flag.String("db-user", "guestbook", "PostgreSQL user")
	flag.String("db-password", "password", "PostgreSQL password")
	flag.String("db-name", "guestbook", "PostgreSQL database")
	flag.String("db-max-idle-time", "15m", "PostgreSQL max idle time")
	flag.String("redis-host", "localhost", "Redis server address")
	flag.String("redis-password", "", "Redis password")

	// Parse command line options into Viper
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.BoolP("help", "h", false, "Prints this overview")
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return nil, err
	}

	if viper.GetBool("help") {
		pflag.Usage()
		os.Exit(0)
	}

	// Parse ENV options into viper with higher priority than command line
	viper.AutomaticEnv()
	viper.SetEnvPrefix("guestbook")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Create config struct

	cfg := config{
		port: viper.GetInt("port"),
		db: databaseCfg{
			host:         viper.GetString("db-host"),
			port:         viper.GetInt("db-port"),
			user:         viper.GetString("db-user"),
			password:     viper.GetString("db-password"),
			database:     viper.GetString("db-name"),
			maxOpenConns: viper.GetInt("db-max-open-conns"),
			maxIdleConns: viper.GetInt("db-max-idle-conns"),
			maxIdleTime:  viper.GetString("db-max-idle-time"),
		},
		redis: redisCfg{
			host:     viper.GetString("redis-host"),
			port:     viper.GetInt("redis-port"),
			password: viper.GetString("redis-password"),
		},
	}

	return &cfg, nil
}
