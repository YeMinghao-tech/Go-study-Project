package router

import (
	"mall/adaptor"
	"mall/api/admin"
	"mall/api/customer"
	"mall/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 路由接口：定义路由注册、链路追踪过滤、访问日志过滤规范
type IRouter interface {
	Register(engine *gin.Engine)            // 注册所有路由
	SpanFilter(r *gin.Context) bool         // 链路追踪过滤（是否开启Span）
	AccessRecordFilter(r *gin.Context) bool // 访问日志过滤（是否记录日志）
}

// 路由管理器：实现IRouter接口，聚合所有路由相关依赖
type Router struct {
	FullPPROF bool           // 是否开启PPROF性能分析
	rootPath  string         // 全局路由根路径（/api/mall）
	conf      *config.Config // 全局配置
	checkFunc func() error   // 服务健康检查函数
	admin     *admin.Ctrl    // 管理员模块控制器
	customer  *customer.Ctrl // 客户模块控制器
}

func NewRouter(conf *config.Config, adaptor adaptor.IAdaptor, checkFunc func() error) *Router {
	return &Router{
		FullPPROF: conf.Server.EnablePprof,   // 从配置读取是否开启PPROF
		rootPath:  "/api/mall",               // 全局路由根路径
		conf:      conf,                      // 注入全局配置
		checkFunc: checkFunc,                 // 注入健康检查函数
		admin:     admin.NewCtrl(adaptor),    // 创建管理员模块控制器（注入适配器）
		customer:  customer.NewCtrl(adaptor), // 创建客户模块控制器（注入适配器）
	}
}
func (r *Router) checkServer() func(*gin.Context) {
	return func(ctx *gin.Context) {
		err := r.checkFunc() // 执行外部传入的健康检查逻辑
		if err != nil {
			// 检查失败：返回500错误
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		// 检查成功：返回200空响应
		ctx.JSON(http.StatusOK, gin.H{})
	}
}

func (r *Router) Register(app *gin.Engine) {
	// 1. 全局中间件：注册认证中间件（传入SpanFilter控制链路追踪）
	app.Use(AuthMiddleware(r.SpanFilter))
	// 2. 开启PPROF性能分析（从配置控制）
	if r.conf.Server.EnablePprof {
		SetupPprof(app, "/debug/pprof")
	}
	// 3. 健康检查接口：所有请求方法（GET/POST等）都支持
	app.Any("/ping", r.checkServer())

	// 4. 全局路由分组：所有业务接口都挂载在 /api/mall 下
	root := app.Group(r.rootPath)
	// 5. 注册具体业务路由
	r.route(root)
}

func (r *Router) SpanFilter(ctx *gin.Context) bool {
	return true
}

func (r *Router) AccessRecordFilter(ctx *gin.Context) bool {
	return true
}

func (r *Router) route(root *gin.RouterGroup) {

	//root.GET("/hello", r.admin.HelloWorld)
	// 管理员模块路由分组：/api/mall/admin
	adminRoot := root.Group("/admin")
	// 具体接口：/api/mall/admin/user/info（GET方法）
	adminRoot.GET("/user/info", r.admin.GetUserInfo)
	// 可扩展：添加更多admin接口，比如 adminRoot.POST("/user/add", r.admin.AddUser)
}
