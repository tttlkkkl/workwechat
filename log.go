package workwechat

import (
	"fmt"
	"io"
	"log"
	"os"
)

// 默认的全局日志变量
var loger Loger

//默认的日志全局日志输出
func init() {
	loger = NewLogger(SetCallDepth(3))
}

// Loger 日志接口 如果要对日志进行扩展都应该遵循这个接口
type Loger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warning(msg string, err error, args ...interface{})
	Error(msg string, err error, args ...interface{})
	Fatal(msg string, err error, args ...interface{})
	Output(calldepth int, s string) error
}

// Logging 日志打印统一封装
type Logging struct {
	//Debug 调试日志
	debug *log.Logger
	//Info 重要提示
	info *log.Logger
	//Warning 错误日志
	warning *log.Logger
	//Error 严重的错误日志
	err *log.Logger
	//Fatal 致命的错误日志
	fatal *log.Logger
	// all 原始输出
	all *log.Logger
	Options
}

// LogLevel 日志级别
type LogLevel uint8

const (
	// LevelAll 全部日志
	LevelAll LogLevel = iota
	// LevelDebug 调试日志
	LevelDebug
	// LevelInfo 运行信息
	LevelInfo
	//LevelWarning 需要特别注意的信息
	LevelWarning
	// LevelErr 错误日志
	LevelErr
	// LevelFatal 致命的错误日志
	LevelFatal
	// LevelNon 不打印任何日志
	LevelNon
)

// Options 日志选项
type Options struct {
	logLevel  LogLevel
	output    io.Writer
	calldepth int
}

// Option 日志选项设置方法
type Option func(*Options)

// SetLogLevel 设置日志等级
func SetLogLevel(lev LogLevel) Option {
	return func(o *Options) {
		o.logLevel = lev
	}
}

// SetOutPut 设置日志输出路径 这个使用默认的输出到标准输入输出即可，后续用 docker 统一采集
func SetOutPut(out io.Writer) Option {
	return func(o *Options) {
		o.output = out
	}
}

// SetCallDepth 设置调用深度，以确保打印的日志输出代码所在文件的正确性
func SetCallDepth(d int) Option {
	return func(o *Options) {
		o.calldepth = d
	}
}

// SetDefaultLogger 设置默认的日志记录器
func SetDefaultLogger(opt ...Option) {
	loger = NewLogger(opt...)
}

// GetDefaultLogger 获取默认日志记录器
func GetDefaultLogger(opt ...Option) Loger {
	return loger
}

// NewLogger 获取日志打印实例
func NewLogger(opt ...Option) Loger {
	opts := new(Options)
	for _, o := range opt {
		o(opts)
	}
	if opts.logLevel == 0 {
		opts.logLevel = LevelDebug
	}
	if opts.output == nil {
		opts.output = os.Stdout
	}
	if opts.calldepth == 0 {
		opts.calldepth = 2
	}
	log := Logging{
		debug:   logFormat(opts.output, "DEBUG: "),
		info:    logFormat(opts.output, "INFO: "),
		warning: logFormat(opts.output, "WARNING: "),
		err:     logFormat(opts.output, "ERROR: "),
		fatal:   logFormat(opts.output, "FATAL: "),
		all:     log.New(opts.output, "", 0),
		Options: *opts,
	}
	return &log
}

// 设置日志格式
func logFormat(out io.Writer, prex string) *log.Logger {
	return log.New(out, prex, log.Ldate|log.Ltime|log.Lshortfile)
}

//Debug 打印调试日志
func (l *Logging) Debug(msg string, args ...interface{}) {
	if l.logLevel > LevelDebug || l.logLevel == LevelNon {
		return
	}
	l.debug.Output(l.calldepth, format(msg, nil, args...))
}

//Info 打印提示信息日志
func (l *Logging) Info(msg string, args ...interface{}) {
	if l.logLevel > LevelInfo || l.logLevel == LevelNon {
		return
	}
	l.info.Output(l.calldepth, format(msg, nil, args...))
}

//Warning 打印错误日志
func (l *Logging) Warning(msg string, err error, args ...interface{}) {
	if l.logLevel > LevelWarning || l.logLevel == LevelNon {
		return
	}
	l.warning.Output(l.calldepth, format(msg, err, args...))
}

//Error 打印严重的错误日志
func (l *Logging) Error(msg string, err error, args ...interface{}) {
	if l.logLevel > LevelErr || l.logLevel == LevelNon {
		return
	}
	l.err.Output(l.calldepth, format(msg, err, args...))
}

//Fatal 打印致命错误日志，并中断程序执行
func (l *Logging) Fatal(msg string, err error, args ...interface{}) {
	if l.logLevel > LevelFatal || l.logLevel == LevelNon {
		return
	}
	l.fatal.Output(l.calldepth, format(msg, err, args...))
	os.Exit(1)
}

// Output 日志原始输出
func (l *Logging) Output(calldepth int, s string) error {
	return l.all.Output(calldepth, s)
}

// 日志内容格式化 -- 尽量不要使用过多的反射
func format(msg string, err error, v ...interface{}) string {
	if err != nil {
		msg = msg + "|err:" + err.Error()
	}
	if v != nil {
		msg = msg + "|v:" + fmt.Sprintln(v...)
	}
	return msg
}
