package service

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/smtp"
	"net/url"
	"time"
	"weibo/models"
	"weibo/pkg/logging"
	"weibo/pkg/redis"
	"weibo/utils"
)

// Login
// @Summary 用户登录
// @Tags 用户模块
// @param email formData string false "邮箱"
// @param password formData string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /login [post]
func Login(c *gin.Context) {
	emailF := c.PostForm("email")
	password := c.PostForm("password")

	logging.Info(fmt.Sprintf("func Login() 参数: email=%s", emailF))

	account, err := models.FindAccountByEmail(emailF)
	if err != nil {
		logging.Error(fmt.Sprintf("func Login() 错误: err=%v", err))
		utils.Error(c, "服务器错误", http.StatusInternalServerError)
		return
	}

	if account == nil {
		utils.Fail(c, "登录失败")
		return
	}

	flag := utils.ValidPassword(password, account.Salt, account.Password)
	if !flag {
		utils.Fail(c, "登录失败")
		return
	}

	returnData := make(map[string]interface{})

	token, err := utils.GenerateToken(emailF, password)
	if err != nil {
		logging.Error(fmt.Sprintf("func Login() 错误: err=%v", err))
		utils.Error(c, "服务器错误", http.StatusInternalServerError)
		return
	}

	returnData["token"] = token

	utils.Ok(c, returnData)
}

// Code
// @Summary 请求发送注册邮件
// @Tags 用户模块
// @param email formData string false "邮箱"
// @param password formData string false "密码"
// @param token formData string false "token"
// @Success 200 {string} json{"code","message"}
// @Router /code [post]
func Code(c *gin.Context) {
	emailF := c.PostForm("email")
	password := c.PostForm("password")
	tokenF := c.PostForm("token")

	// 向 Google reCAPTCHA 服务器发送请求验证 token
	resp, err := http.PostForm("https://recaptcha.net/recaptcha/api/siteverify", url.Values{
		"secret":   {viper.GetString("reCAPTCHA.Secret")},
		"response": {tokenF},
	})

	if err != nil {
		logging.Error(fmt.Sprintf("func Code() 错误: err=%v", err))
		utils.Error(c, "服务器错误", http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logging.Error(fmt.Sprintf("func Code() 错误: err=%v", err))
			utils.Error(c, "服务器错误", http.StatusInternalServerError)
			return
		}
	}(resp.Body)

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logging.Error(fmt.Sprintf("func Code() 错误: err=%v", err))
		utils.Error(c, "服务器错误", http.StatusInternalServerError)
		return
	}

	// 解析 JSON 响应
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		logging.Error(fmt.Sprintf("func Code() 错误: err=%v", err))
		utils.Error(c, "服务器错误", http.StatusInternalServerError)
		return
	}

	// 获取验证结果
	success, ok := result["success"].(bool)
	if !ok {
		logging.Error(fmt.Sprintf("func Code() 错误: err=%v", ok))
		utils.Error(c, "服务器错误", http.StatusInternalServerError)
		return
	}

	// 根据验证结果返回相应的信息
	if !success {
		utils.Fail(c, "人机身份检测失败，请重新验证")
		return
	}

	logging.Info(fmt.Sprintf("func Code() 参数: email=%s", emailF))

	if !utils.IsValidEmail(emailF) {
		utils.Fail(c, "邮箱格式错误")
		return
	}

	if !utils.IsValidPassword(password) {
		utils.Fail(c, "密码长度是6到30，包含且只能数字与字母")
		return
	}

	isExist, err := models.ExistAccountByEmail(emailF)
	if err != nil {
		logging.Error(fmt.Sprintf("func Code() 错误: err=%v", err))
		utils.Error(c, "服务器错误", http.StatusInternalServerError)
		return
	}

	if isExist {
		utils.Fail(c, "邮箱号已注册")
		return
	}

	em := email.NewEmail()
	ssl := true
	em.From = viper.GetString("email.Username") // 邮件的发件人地址
	em.To = []string{emailF}                    // 邮件的收件人地址
	em.Subject = "用户注册"                         // 邮件的主题
	auth := smtp.PlainAuth("", viper.GetString("email.Username"), viper.GetString("email.Password"), viper.GetString("email.Host"))

	token, _ := utils.GenerateToken(emailF, password)

	// 设置邮件内容
	content := "邮件有效期：3分钟\n点击链接完成验证：" + "http://" + viper.GetString("server.Ip") + ":" + viper.GetString("server.HttpPort") + "/register?token=" + token
	em.Text = []byte(content) // 邮件的文本内容

	if ssl {
		err = em.SendWithTLS(viper.GetString("email.Host")+":465", auth, &tls.Config{ServerName: viper.GetString("email.Host")})
	} else {
		err = em.Send(viper.GetString("email.Host")+":"+viper.GetString("email.Port"), auth)
	}

	//Host：邮件服务器的主机名、Username：用于认证的电子邮件用户名、Password：用于认证的电子邮件密码
	if err != nil {
		logging.Error(fmt.Sprintf("func Code() 错误: err=%v", err))
		utils.Error(c, "服务器错误", http.StatusInternalServerError)
		return
	}

	err = redis.SetEmailToken(token, 3*time.Minute)
	if err != nil {
		logging.Error(fmt.Sprintf("func Code() 错误: err=%v", err))
		utils.Error(c, "服务器错误", http.StatusInternalServerError)
		return
	}

	logging.Info(fmt.Sprintf("func Code() %s 邮件发送成功", emailF))
	utils.Ok(c, nil)
}

// Register
// @Summary 注册用户
// @Tags 用户模块
// @param token query string false "token"
// @Success 200 {string} json{"code","message"}
// @Router /register [get]
func Register(c *gin.Context) {
	account := models.Account{}
	token := c.Query("token")

	claims, err := utils.ParseToken(token)
	if err != nil {
		logging.Error(fmt.Sprintf("func Register() 错误: err=%v", err))
		utils.Error(c, "服务器错误", http.StatusInternalServerError)
		return
	}

	exists, err := redis.EmailTokenExists(token)
	if err != nil {
		logging.Error(fmt.Sprintf("func Register() 错误: err=%v", err))
		utils.Error(c, "服务器错误", http.StatusInternalServerError)
		return
	}

	if !exists {
		utils.Fail(c, "邮件有效期已过")
		return
	}

	isExist, err := models.ExistAccountByEmail(claims.Email)
	if err != nil {
		logging.Error(fmt.Sprintf("func Register() 错误: err=%v", err))
		utils.Error(c, "服务器错误", http.StatusInternalServerError)
		return
	}

	if isExist {
		utils.Fail(c, "邮箱号已注册")
		return
	}

	account.Email = claims.Email
	salt := fmt.Sprintf("%06d", rand.Int31())
	account.Password = utils.MakePassword(claims.Password, salt)
	account.Salt = salt

	models.CreateUser(account)

	logging.Info(fmt.Sprintf("func Register() %s 邮箱号注册成功", claims.Email))
	utils.Ok(c, "您的账号已注册成功，可以关闭此页面了")
}
