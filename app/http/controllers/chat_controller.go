package controllers

import (
	"github.com/869413421/chatgpt-web/config"
	"github.com/869413421/chatgpt-web/pkg/logger"
	"github.com/gin-gonic/gin"
	gogpt "github.com/sashabaranov/go-gpt3"
	"net/http"
	"strings"
	"sync"
)

var once sync.Once

// ChatController 首页控制器
type ChatController struct {
	BaseController
}

// NewChatController 创建控制器
func NewChatController() *ChatController {
	return &ChatController{}
}

//Index 首页
func (c *ChatController) Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main website",
	})
}

func (c *ChatController) Avatar(ctx *gin.Context) {

	once.Do(initUser)

	user := ctx.MustGet(gin.AuthUserKey).(string)
	logger.Info("user:" + user)

	u := userMap[user]
	avatar := ""
	if u.Avatar != "" {
		avatar = u.Avatar
	}
	ctx.JSON(200, gin.H{
		"avatar": avatar,
	})
}

func initUser() {
	cnf := config.LoadConfig()
	users := cnf.Users

	for _, u := range users {
		userMap[u.Name] = u
	}
}

var userMap = map[string]config.User{}

//Completion 回复
func (c *ChatController) Completion(ctx *gin.Context) {

	once.Do(initUser)

	var request gogpt.ChatCompletionRequest

	err := ctx.BindJSON(&request)
	if err != nil {
		c.ResponseJson(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	user := ctx.MustGet(gin.AuthUserKey).(string)
	logger.Info("user:" + user)
	logger.Info(request)
	if len(request.Messages) == 0 {
		c.ResponseJson(ctx, http.StatusBadRequest, "request messages required", nil)
		return
	}

	cnf := config.LoadConfig()

	client := gogpt.NewClient(cnf.ApiKey)

	botDesc := cnf.BotDesc
	if userMap[user].BotDesc != "" {
		botDesc = userMap[user].BotDesc
	}

	if request.Messages[0].Role != "system" {
		newMessage := append([]gogpt.ChatCompletionMessage{
			{Role: "system", Content: botDesc},
		}, request.Messages...)
		request.Messages = newMessage
		logger.Info(request.Messages)
	}

	if cnf.Model == gogpt.GPT3Dot5Turbo0301 || cnf.Model == gogpt.GPT3Dot5Turbo {
		request.Model = cnf.Model
		resp, err := client.CreateChatCompletion(ctx, request)
		if err != nil {
			c.ResponseJson(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		c.ResponseJson(ctx, http.StatusOK, "", gin.H{
			"reply":    resp.Choices[0].Message.Content,
			"messages": append(request.Messages, resp.Choices[0].Message),
		})
	} else {
		prompt := ""
		for _, item := range request.Messages {
			prompt += item.Content + "/n"
		}
		prompt = strings.Trim(prompt, "/n")

		logger.Info("request prompt is %s", prompt)
		req := gogpt.CompletionRequest{
			Model:            cnf.Model,
			MaxTokens:        cnf.MaxTokens,
			TopP:             cnf.TopP,
			FrequencyPenalty: cnf.FrequencyPenalty,
			PresencePenalty:  cnf.PresencePenalty,
			Prompt:           prompt,
		}

		resp, err := client.CreateCompletion(ctx, req)
		if err != nil {
			c.ResponseJson(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		c.ResponseJson(ctx, http.StatusOK, "", gin.H{
			"reply": resp.Choices[0].Text,
			"messages": append(request.Messages, gogpt.ChatCompletionMessage{
				Role:    "assistant",
				Content: resp.Choices[0].Text,
			}),
		})

		/*
			c.ResponseJson(ctx, http.StatusOK, "", gin.H{
				"reply": "123",
				"messages": append(request.Messages, gogpt.ChatCompletionMessage{
					Role:    "assistant",
					Content: "456",
				}),
			})
		*/
	}

}
