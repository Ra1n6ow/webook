package web

import (
	"net/http"
	"time"

	"github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ra1n6ow/webook/internal/domain"
	"github.com/ra1n6ow/webook/internal/service"
)

const (
	bizLogin          = "login"
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	nicknameRegexPattern = `^.{4,16}$`
	birthdayRegexPattern = `^(19|20)\d{2}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])$`
	introRegexPattern    = `^.{5,1024}$`
)

type UserHandler struct {
	emailRexExp          *regexp2.Regexp
	passwordRexExp       *regexp2.Regexp
	nicknameRegexPattern *regexp2.Regexp
	introRegexPattern    *regexp2.Regexp
	birthdayRegexPattern *regexp2.Regexp
	userSvc              service.UserServicer
	codeSvc              service.CodeServicer
}

func NewUserHandler(userSvc service.UserServicer, codeSvc service.CodeServicer) *UserHandler {
	return &UserHandler{
		emailRexExp:          regexp2.MustCompile(emailRegexPattern, regexp2.None),
		passwordRexExp:       regexp2.MustCompile(passwordRegexPattern, regexp2.None),
		nicknameRegexPattern: regexp2.MustCompile(nicknameRegexPattern, regexp2.None),
		introRegexPattern:    regexp2.MustCompile(introRegexPattern, regexp2.None),
		birthdayRegexPattern: regexp2.MustCompile(birthdayRegexPattern, regexp2.None),
		userSvc:              userSvc,
		codeSvc:              codeSvc,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", h.SignUp)
	ug.POST("/login", h.LoginJWT)
	ug.POST("/edit", h.Edit)
	ug.GET("/profile", h.Profile)

	// 手机验证码登录相关功能
	ug.POST("/login_sms/code/send", h.SendSMSLoginCode)
	ug.POST("/login_sms", h.LoginSMS)
}

func (h *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok, err := h.codeSvc.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码不对，请重新输入",
		})
		return
	}
	u, err := h.userSvc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	h.setJWTToken(ctx, u.Id)
	ctx.JSON(http.StatusOK, Result{
		Msg: "登录成功",
	})
}

func (h *UserHandler) SendSMSLoginCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 你这边可以校验 Req
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "请输入手机号码",
		})
		return
	}
	err := h.codeSvc.Send(ctx, bizLogin, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "短信发送太频繁，请稍后再试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 补日志的
	}
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	type signUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req signUpReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	isEmail, err := h.emailRexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "非法邮箱格式")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入密码不对")
		return
	}

	isPassword, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码必须包含字母、数字、特殊字符，并且不少于八位")
		return
	}

	ud := domain.User{
		Email:    req.Email,
		Password: req.Password,
	}

	err = h.userSvc.SignUp(ctx, ud)
	switch err {
	case nil:
		ctx.String(http.StatusOK, "注册成功")
		return
	case service.ErrDuplicateEmail:
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	default:
		ctx.String(http.StatusOK, "系统错误")
		return
	}
}

func (h *UserHandler) LoginJWT(ctx *gin.Context) {
	type loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req loginReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	u, err := h.userSvc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		h.setJWTToken(ctx, u.Id)

		ctx.String(http.StatusOK, "登录成功")
		return
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或密码错误")
		return
	default:
		ctx.String(http.StatusOK, "系统错误")
		return
	}
}

func (h *UserHandler) Login(ctx *gin.Context) {
	type loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req loginReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	u, err := h.userSvc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		session := sessions.Default(ctx)
		session.Set("userId", u.Id)
		session.Options(sessions.Options{
			MaxAge: 60 * 60 * 24 * 30,
			// 生产环境需要开启 HttpOnly 和 Secure
			// 防止 XSS 攻击，只能通过 http/https 访问
			// HttpOnly: true,
			// 只允许 https 访问
			// Secure:   true,
		})
		session.Save()
		ctx.String(http.StatusOK, "登录成功")
		return
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或密码错误")
		return
	default:
		ctx.String(http.StatusOK, "系统错误")
		return
	}
}

func (h *UserHandler) Edit(ctx *gin.Context) {
	type editReq struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}

	var req editReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	userId, ok := ctx.Get("userId")
	if !ok {
		ctx.String(http.StatusOK, "用户未登录")
		return
	}
	userIdInt64, ok := userId.(int64)
	if !ok {
		ctx.String(http.StatusOK, "用户ID类型错误")
		return
	}

	/* 从 session 中获取 userId
	session := sessions.Default(ctx)
	userId, ok := session.Get("userId").(int64)
	if !ok {
		ctx.String(http.StatusOK, "用户未登录")
		return
	}
	*/

	_, err := h.userSvc.Profile(ctx, userIdInt64)
	if err != nil {
		ctx.String(http.StatusOK, "session 错误, 未找到用户")
		return
	}

	isNickname, err := h.nicknameRegexPattern.MatchString(req.Nickname)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isNickname {
		ctx.String(http.StatusOK, "昵称格式错误")
		return
	}

	isAboutMe, err := h.introRegexPattern.MatchString(req.AboutMe)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isAboutMe {
		ctx.String(http.StatusOK, "简介格式错误")
		return
	}

	isBirthday, err := h.birthdayRegexPattern.MatchString(req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if !isBirthday {
		ctx.String(http.StatusOK, "生日格式错误")
		return
	}

	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "日期格式错误")
		return
	}

	ud := domain.User{
		Id:       userIdInt64,
		Nickname: req.Nickname,
		Birthday: birthday,
		AboutMe:  req.AboutMe,
	}

	err = h.userSvc.Edit(ctx, ud)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "编辑成功")
}

func (h *UserHandler) Profile(ctx *gin.Context) {
	userId, ok := ctx.Get("userId")
	if !ok {
		ctx.String(http.StatusOK, "用户未登录")
		return
	}
	userIdInt64, ok := userId.(int64)
	if !ok {
		ctx.String(http.StatusOK, "用户ID类型错误")
		return
	}
	/*  从 session 中获取 userId
	session := sessions.Default(ctx)
	userId, ok := session.Get("userId").(int64)
	if !ok {
		ctx.String(http.StatusOK, "用户未登录 from session")
		return
	}
	*/

	u, err := h.userSvc.Profile(ctx, userIdInt64)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, u)
}

func (h *UserHandler) setJWTToken(ctx *gin.Context, userId int64) {
	uc := UserClaims{
		Uid:       userId,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			// 2 小时过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 120)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenStr, err := token.SignedString([]byte("1234567890"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	ctx.Header("x-jwt-token", tokenStr)
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}
