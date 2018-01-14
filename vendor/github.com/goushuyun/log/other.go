package log

import (
	"time"

	"google.golang.org/grpc/metadata"

	"golang.org/x/net/context"
)

// ======== 兼容 qiniu/log   ===============
func SetOutputLevel(l interface{}) {
	n, _ := l.(int)
	v = Level(n)
}

const (
	Lmodule = iota
	Llevel
	LstdFlags
)

func SetFlags(flag int) {

}

// ======== 兼容 wothing/log ===============

// TraceIn and TraceOut use in function in and out,reduce code line
// Example:
//	func test() {
//		user := User{Name: "zhangsan", Age: 21, School: "xayddx"}
//		service := "verification.GetVerifiCode"
//		defer log.TraceOut(log.TraceIn("12345", service, "user:%v", user))
//		....
//	}

// TraceIn 方法入口打印日志
func TraceIn(tag string, method string, format string, m ...interface{}) (string, string, time.Time) {
	startTime := time.Now()
	std.Tprintf(Linfo, tag, "calling "+method+", "+format, m...)
	return tag, method, startTime
}

// TraceCtx 方法入口打印日志
func TraceCtx(ctx context.Context, method string, format string, m ...interface{}) (string, string, time.Time) {
	tag := "-"
	if md, ok := metadata.FromContext(ctx); ok {
		if md["tid"] != nil && len(md["tid"]) > 0 {
			tag = md["tid"][0]
		}
	}
	startTime := time.Now()
	std.Tprintf(Linfo, tag, "calling "+method+", "+format, m...)
	return tag, method, startTime
}

// TraceOut 方法退出记录下消耗时间
func TraceOut(tag string, method string, startTime time.Time) {
	std.Tprintf(Linfo, tag, "finished "+method+", took %v", time.Since(startTime))
}

func Println(m ...interface{}) { std.Tprintf(Lprint, "", "", m...) }

func getTracerIDFromCtx(ctx context.Context) string {
	nid := "00000000-0000-0000-0000-000000000000"

	if ctx == nil {
		return nid
	}

	if md, ok := metadata.FromContext(ctx); ok {
		if md["tid"] != nil && len(md["tid"]) > 0 {
			return md["tid"][0]
		}
	}
	return nid
}

func CtxDebugf(ctx context.Context, format string, m ...interface{}) {
	std.Tprintf(Ldebug, getTracerIDFromCtx(ctx), format, m...)
}

func CtxDebug(ctx context.Context, m ...interface{}) {
	std.Tprintf(Ldebug, getTracerIDFromCtx(ctx), "", m...)
}

func CtxInfof(ctx context.Context, format string, m ...interface{}) {
	std.Tprintf(Linfo, getTracerIDFromCtx(ctx), format, m...)
}

func CtxInfo(ctx context.Context, m ...interface{}) {
	std.Tprintf(Linfo, getTracerIDFromCtx(ctx), "", m...)
}

func CtxWarnf(ctx context.Context, format string, m ...interface{}) {
	std.Tprintf(Lwarn, getTracerIDFromCtx(ctx), format, m...)
}

func CtxWarn(ctx context.Context, m ...interface{}) {
	std.Tprintf(Lwarn, getTracerIDFromCtx(ctx), "", m...)
}

func CtxErrorf(ctx context.Context, format string, m ...interface{}) {
	std.Tprintf(Lerror, getTracerIDFromCtx(ctx), format, m...)
}

func CtxError(ctx context.Context, m ...interface{}) {
	std.Tprintf(Lerror, getTracerIDFromCtx(ctx), "", m...)
}

func CtxFatal(ctx context.Context, m ...interface{}) {
	std.Tprintf(Lfatal, getTracerIDFromCtx(ctx), "", m...)
}

func CtxFatalf(ctx context.Context, format string, m ...interface{}) {
	std.Tprintf(Lfatal, getTracerIDFromCtx(ctx), format, m...)
}

func CtxFatalln(ctx context.Context, m ...interface{}) {
	std.Tprintf(Lfatal, getTracerIDFromCtx(ctx), "", m...)
}

func CtxPanic(ctx context.Context, m ...interface{}) {
	std.Tprintf(Lpanic, getTracerIDFromCtx(ctx), "", m...)
}

func CtxPanicf(ctx context.Context, format string, m ...interface{}) {
	std.Tprintf(Lpanic, getTracerIDFromCtx(ctx), format, m...)
}

func CtxPanicln(ctx context.Context, m ...interface{}) {
	std.Tprintf(Lpanic, getTracerIDFromCtx(ctx), "", m...)
}

func CtxStack(ctx context.Context, m ...interface{}) {
	std.Tprintf(Lstack, getTracerIDFromCtx(ctx), "", m...)
}

func CtxPrint(ctx context.Context, m ...interface{}) {
	std.Tprintf(Lprint, getTracerIDFromCtx(ctx), "", m...)
}

func CtxPrintf(ctx context.Context, format string, m ...interface{}) {
	std.Tprintf(Lprint, getTracerIDFromCtx(ctx), format, m...)
}

func CtxPrintln(ctx context.Context, m ...interface{}) {
	std.Tprintf(Lprint, getTracerIDFromCtx(ctx), "", m...)
}
