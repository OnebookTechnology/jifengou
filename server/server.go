package server

import (
	"fmt"
	captcha "github.com/OnebookTechnology/captcha-sdk"
	csst "github.com/OnebookTechnology/etcd-sdk"
	"github.com/OnebookTechnology/jifengou/mysqlservice"
	"github.com/OnebookTechnology/jifengou/server/interface"
	sms "github.com/OnebookTechnology/smssdk/sdk"
	levelLogger "github.com/cxt90730/LevelLogger-go"
	"github.com/gin-gonic/gin"
	"github.com/robfig/config"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var server *Server
var logger *levelLogger.LevelLogger

type Conf struct {
	port    int
	tcpPort int

	certPem string
	certKey string

	domain   string
	logLevel int
	logDir   string

	ueditorConf *UEditorConfig
}

type Server struct {
	ServerName string
	Conf
	Env JFGEnv

	localIP             string
	initialMembers      []string
	loginSessionTimeout int
	sessionTimeout      int

	HttpServer  *http.Server
	TcpListener *net.TCPListener
	DB          _interface.ServerDB
	Consist     _interface.Consistence
	Captcha     _interface.Captcha
	SMS         _interface.SMS

	accessLog *os.File
	errorLog  *os.File
	Logger    *levelLogger.LevelLogger

	closeChan chan bool
}

func NewService(confPath, serverName string, mode string) (*Server, error) {
	server = new(Server)
	server.ServerName = serverName
	server.closeChan = make(chan bool)

	//Configure
	err := loadByConf(confPath)
	if err != nil {
		return nil, err
	}

	//LocalIP
	server.localIP, err = getLocalIP()
	if err != nil {
		return nil, err
	}

	//Log
	if err := server.createLogs(server.logDir); err != nil {
		return nil, err
	}

	//Server
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	router := gin.New()
	router.Use(gin.LoggerWithWriter(server.accessLog), gin.RecoveryWithWriter(server.errorLog))

	//DB
	db := new(mysql.MysqlService)
	db.InitialDB(confPath, "DB")
	server.DB = db

	//Consistence
	consistService, err := csst.NewEtcdService(confPath)
	if err != nil {
		return nil, err
	}
	server.Consist = consistService

	//SMS
	smsService, err := sms.NewSMSService(confPath)
	if err != nil {
		return nil, err
	}
	server.SMS = smsService

	//Captcha
	capt, err := captcha.NewCaptchaService(confPath)
	if err != nil {
		return nil, err
	}
	server.Captcha = capt

	//Env
	//server.Env = testEnv
	//if mode == "online" {
	//	server.Env = onlineEnv
	//	server.Conf.port = 80
	//}
	server.Env = onlineEnv
	LoadRouter(router)

	//Start Bridge
	srvIPPort := ":" + strconv.Itoa(server.port)

	srv := &http.Server{
		Addr:     srvIPPort,
		Handler:  router,
		ErrorLog: nil,
	}

	server.HttpServer = srv
	fmt.Println("Start server", srvIPPort)
	logger.Info("Start server", srvIPPort)

	return server, nil
}

func loadByConf(confPath string) error {
	c, err := config.ReadDefault(confPath)
	if err != nil {
		return err
	}

	server.port, err = c.Int("server", "port")
	if err != nil {
		return err
	}

	server.logLevel, err = c.Int("server", "log_level")
	if err != nil {
		return err
	}
	server.logDir, err = c.String("server", "log_dir")
	if err != nil {
		return err
	}

	server.Conf.domain = "http://zhuzhushop.cn"

	ec := new(UEditorConfig)
	ec.ImageUrl, _ = c.String("ueditor", "ImageUrl")
	ec.ImagePath, _ = c.String("ueditor", "ImagePath")
	ec.ImageFieldName, _ = c.String("ueditor", "ImageFieldName")
	ec.ImageMaxSize, _ = c.Int("ueditor", "ImageMaxSize")
	ec.ImageActionName, _ = c.String("ueditor", "ImageActionName")
	ec.ImageAllowFiles = ImageAllowFiles
	server.ueditorConf = ec
	return nil
}

//Start all service and router
func (s *Server) Start() error {

	//Start server. The server will shutdown gracefully.
	logger.Info("##### start server ######")

	if err := s.HttpServer.ListenAndServe(); err != nil {
		logger.Error("Listen", ":"+strconv.Itoa(server.port), err)
		return err
	}
	return nil
}

//Create all logs
func (s *Server) createLogs(logDir string) error {

	logDir = strings.TrimSuffix(logDir, string(os.PathSeparator))
	os.MkdirAll(logDir, 0755)
	srvLogPath := logDir + string(os.PathSeparator) + "server.log"
	accessLogPath := logDir + string(os.PathSeparator) + "access.log"
	errorLogPath := logDir + string(os.PathSeparator) + "error.log"

	// 建立 ServerLogger
	logFile, err := os.OpenFile(srvLogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	accessLogFile, err := os.OpenFile(accessLogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	errorLogFile, err := os.OpenFile(errorLogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	logger, _ = levelLogger.NewLevelLogger(logFile, "", log.LstdFlags, s.logLevel)

	s.accessLog = accessLogFile
	s.errorLog = errorLogFile
	s.Logger = logger

	return nil
}
