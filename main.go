package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/op/go-logging"
	"x-ui/config"
	"x-ui/database"
	"x-ui/logger"
	"x-ui/web"
	"x-ui/web/global"
	"x-ui/web/service"
)

func runWebServer() {
	log.Printf("%v %v", config.GetName(), config.GetVersion())

	switchLogger := func(level logging.Level) {
		logger.SetLevel(level)
		log.Println("Log level set to:", level)
	}

	eLogLevel := os.Getenv("XUI_LOG_LEVEL")
	if eLogLevel == "" {
		eLogLevel = "info"
	}

	switch eLogLevel {
	case "debug":
		switchLogger(logging.DEBUG)
	case "info":
		switchLogger(logging.INFO)
	case "warn":
		switchLogger(logging.WARNING)
	case "error":
		switchLogger(logging.ERROR)
	default:
		switchLogger(logging.INFO)
	}

	err := database.InitDB(config.GetDBPath())
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}

	var server *web.Server
	server = web.NewServer()
	global.SetWebServer(server)
	err = server.Start()
	if err != nil {
		log.Fatal("Error starting web server:", err)
	}
}

func resetSetting() {
	err := database.InitDB(config.GetDBPath())
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}

	settingService := service.SettingService{}
	err = settingService.ResetSettings()
	if err != nil {
		fmt.Println("Error resetting settings:", err)
		return
	}
	fmt.Println("Settings successfully reset to defaults.")
}

func showSetting(show bool) {
	if show {
		err := database.InitDB(config.GetDBPath())
		if err != nil {
			fmt.Println("Error initializing database:", err)
			return
		}
		settingService := service.SettingService{}
		port, err := settingService.GetPort()
		if err != nil {
			fmt.Println("Error retrieving port:", err)
			return
		}
		userService := service.UserService{}
		userModel, err := userService.GetFirstUser()
		if err != nil {
			fmt.Println("Error retrieving user:", err)
			return
		}
		fmt.Printf("username: %v\n", userModel.Username)
		fmt.Printf("password: %v\n", userModel.Password)
		fmt.Printf("port: %v\n", port)
	}
}

func main() {
	// Define CLI flags for administrative operations
	var showSettingFlag bool
	var resetSettingFlag bool

	flag.BoolVar(&showSettingFlag, "show", false, "Show current settings (username, password, port)")
	flag.BoolVar(&resetSettingFlag, "reset", false, "Reset all settings to default values")
	flag.Parse()

	switch {
	case resetSettingFlag:
		resetSetting()
	case showSettingFlag:
		showSetting(true)
	default:
		runWebServer()
	}
}
