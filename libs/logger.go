package libs

import (
	"os"
	"log"
	"runtime"
	"fmt"
	"path"
)

var (
	Log Logging
	log_file *os.File
	err_log_file *os.File
)

type LogLevel int

type Logging struct  {
	Level LogLevel
	NormalLog  *log.Logger
	ErrorLog *log.Logger
	FatalLog *log.Logger
}


const (
	VERBOSE LogLevel = iota
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
)

const (
	SERVER_TAG = " [ SERV ] "
	VERBOSE_TAG =" [ VERB ] "
	DEBUG_TAG = " [ DEBU ] "
	INFO_TAG = " [ INFO ] %s:%d: "
	WARNING_TAG = " [ WARN ] %s:%d: "
	ERROR_TAG = " [ ERRO ] %s:%d: "
	FATAL_TAG = " [ FATA ] "
)

func (logging *Logging) Verbose(v ...interface{})  {
	if(logging.Level <= VERBOSE){
		logging.NormalLog.SetPrefix(VERBOSE_TAG)
		logging.NormalLog.Println(v...)
	}
}

func (logging *Logging) Verbosef(format string, v ...interface{}) {
	if(logging.Level <= VERBOSE){
		logging.NormalLog.SetPrefix(VERBOSE_TAG)
		logging.NormalLog.Printf(format,v...)
	}
}


func (logging *Logging) Debug(v ...interface{})  {
	if(logging.Level <= DEBUG){
		logging.NormalLog.SetPrefix(DEBUG_TAG)
		logging.NormalLog.Println(v...)
	}
}

func (logging *Logging) Debugf(format string, v ...interface{}) {
	if(logging.Level <= DEBUG){
		logging.NormalLog.SetPrefix(DEBUG_TAG)
		logging.NormalLog.Printf(format,v...)
	}
}

func (logging *Logging) Info(v ...interface{})  {
	if(logging.Level <= INFO){
		_, file, line, _ := runtime.Caller(1)
		logging.NormalLog.SetPrefix(fmt.Sprintf(INFO_TAG,path.Base(file),line))
		logging.NormalLog.Println(v...)
	}
}

func (logging *Logging) Infof(format string, v ...interface{})  {
	if(logging.Level <= INFO){
		_, file, line, _ := runtime.Caller(1)
		logging.NormalLog.SetPrefix(fmt.Sprintf(INFO_TAG,path.Base(file),line))
		logging.NormalLog.Printf(format,v...)
	}
}

func (logging *Logging) Warning(v ...interface{})  {
	if(logging.Level <= WARNING){
		_, file, line, _ := runtime.Caller(1)
		logging.NormalLog.SetPrefix(fmt.Sprintf(WARNING_TAG,path.Base(file),line))
		logging.NormalLog.Println(v...)
	}
}

func (logging *Logging) Warningf(format string, v ...interface{})  {
	if(logging.Level <= WARNING){
		_, file, line, _ := runtime.Caller(1)
		logging.NormalLog.SetPrefix(fmt.Sprintf(WARNING_TAG,path.Base(file),line))
		logging.NormalLog.Printf(format,v...)
	}
}

func (logging *Logging) Error(v ...interface{})  {
	if(logging.Level <= ERROR){
		_, file, line, _ := runtime.Caller(1)
		logging.ErrorLog.SetPrefix(fmt.Sprintf(ERROR_TAG,path.Base(file),line))
		logging.ErrorLog.Println(v...)
	}
}

func (logging *Logging) Errorf(format string, v ...interface{})  {
	if(logging.Level <= ERROR){
		_, file, line, _ := runtime.Caller(1)
		logging.ErrorLog.SetPrefix(fmt.Sprintf(ERROR_TAG,path.Base(file),line))
		logging.ErrorLog.Printf(format,v...)
	}
}

func (logging *Logging) Fatal(v ...interface{})  {
	if(logging.Level <= FATAL){
		logging.FatalLog.SetPrefix(FATAL_TAG)
		logging.FatalLog.Panic(v...)
	}
}

func (logging *Logging) Fatalf(format string, v ...interface{})  {
	if(logging.Level <= FATAL){
		logging.FatalLog.SetPrefix(FATAL_TAG)
		logging.FatalLog.Panicf(format,v...)
	}
}

func LoadLoggerModule()  {
	if(Config.LogToFile){
		output, err := os.OpenFile(Config.LogFilePath, os.O_WRONLY  | os.O_SYNC | os.O_CREATE | os.O_APPEND, 0666)
		if err != nil{
			log.Fatalln(err)
			os.Exit(1)
		}else{
			log_file = output
		}

		output_err, err := os.OpenFile(Config.ErrLogFilePath, os.O_WRONLY  | os.O_SYNC | os.O_CREATE | os.O_APPEND, 0666)
		if err != nil{
			log.Fatalln(err)
			os.Exit(1)
		}else{
			err_log_file = output_err
		}

	}else{
		log_file = os.Stdout
		err_log_file = os.Stderr
	}
	normalLog := log.New(log_file,"",0)
	errorLog := log.New(err_log_file,"",0)
	fatalLog := log.New(err_log_file,"",log.Lshortfile|log.LstdFlags)

	var level LogLevel
	switch Config.LogLevel {
	case "verbose":
		level = VERBOSE
	case "debug":
		level =  DEBUG
	case "info":
		level =  INFO
	case "warning":
		level =  WARNING
	case "error":
		level =  ERROR
	case "fatal":
		level =  FATAL
	default:
		level = DEBUG
	}

	Log = Logging{level,normalLog,errorLog,fatalLog}

	PrintModuleLoaded("Logger")

}

func ReleaseLoggerModule()  {
	if Config.LogToFile {
		if log_file != nil{
			log_file.Close()
		}

		if err_log_file != nil {
			err_log_file.Close()
		}
	}
	PrintModuleRelease("Logger")
}