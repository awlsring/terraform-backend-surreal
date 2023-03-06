package server

import (
	"fmt"
	"regexp"

	"github.com/awlsring/terraform-backend-surreal/pkg/config"
	"github.com/awlsring/terraform-backend-surreal/pkg/state"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var SYNTAX_VERSION = 4
var NAME_REGEX = "^[a-zA-Z0-9_]+$"

func Authenticator(users map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, pass, ok := c.Request.BasicAuth()
		if ok {
			if pass == users[user] {
				c.Next()
				return
			}
		}
		c.Writer.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		c.AbortWithStatus(401)
	}
}

type UriParams struct {
	Project string `uri:"project" binding:"required"`
	Stack   string `uri:"stack" binding:"required"`
}

func GetParams(c *gin.Context) (UriParams, error) {
	var params UriParams
	if err := c.ShouldBindUri(&params); err != nil {
		return params, err
	}

	re := regexp.MustCompile("^[a-zA-Z0-9_]+$")
	if !re.MatchString(params.Project) {
		return params, fmt.Errorf("project name must only contain letter, number or '_'. Was: %s", params.Project)
	}

	if !re.MatchString(params.Stack) {
		return params, fmt.Errorf("stack name must only contain letter, number or '_'. Was: %s", params.Project)
	}

	return params, nil
}

func Start(cfg *config.Config, dao state.StateDao) {
	level, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetLevel(level)
	log.Debug("Debug logging enabled")

	log.Debug("Starting Gin Server in mode: ", cfg.Gin)
	gin.SetMode(cfg.Gin)
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "healthy"})
	})

	stacks := router.Group("")
	stacks.Use(Authenticator(cfg.Users))

	stacks.GET("/:project/:stack", func(c *gin.Context) {
		params, err := GetParams(c)
		if err != nil {
			c.JSON(400, gin.H{"message": err})
			return
		}
		key := fmt.Sprintf("%s%s", params.Project, params.Stack)
		log.Infof("Recieved GetState request for project - stack: %s-%s", params.Project, params.Stack)

		entity, err := dao.Read(key)
		if err != nil {
			entity = state.Entity{
				Id: state.STACK_ID(key),
				State: state.TfState{
					Version: SYNTAX_VERSION,
				},
				Locked: false,
			}
			err := dao.Create(entity)
			if err != nil {
				log.Error(err)
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(200, entity.State)
	})

	stacks.POST("/:project/:stack", func(c *gin.Context) {
		log.Debug("Recieved PostState request")
		params, err := GetParams(c)
		if err != nil {
			c.JSON(400, gin.H{"message": err})
			return
		}

		key := fmt.Sprintf("%s%s", params.Project, params.Stack)
		log.Infof("Recieved PostState request for project - stack: %s-%s", params.Project, params.Stack)

		var state state.TfState
		if err := c.BindJSON(&state); err != nil {
			log.Error(err)
			c.Status(400)
			return
		}

		entity, err := dao.Read(key)
		if err != nil {
			log.Error(err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		entity.State = state
		err = dao.Update(entity)
		if err != nil {
			log.Error(err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.Status(200)
	})

	stacks.DELETE("/:project/:stack", func(c *gin.Context) {
		params, err := GetParams(c)
		if err != nil {
			c.JSON(400, gin.H{"message": err})
			return
		}
		key := fmt.Sprintf("%s%s", params.Project, params.Stack)
		log.Infof("Recieved DeleteState request for project - stack: %s-%s", params.Project, params.Stack)

		err = dao.Delete(key)
		if err != nil {
			log.Error(err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.Status(200)
	})

	stacks.Handle("LOCK", "/:project/:stack", func(c *gin.Context) {
		params, err := GetParams(c)
		if err != nil {
			c.JSON(400, gin.H{"message": err})
			return
		}
		key := fmt.Sprintf("%s%s", params.Project, params.Stack)
		log.Infof("Recieved LockState request for project - stack: %s-%s", params.Project, params.Stack)

		entity, err := dao.Read(key)
		if err != nil {
			log.Error(err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		switch entity.Locked {
		case true:
			c.Status(423)
			return
		case false:
			entity.Locked = true
			dao.Update(entity)
			c.Status(200)
			return
		}
	})

	stacks.Handle("UNLOCK", "/:project/:stack", func(c *gin.Context) {
		params, err := GetParams(c)
		if err != nil {
			c.JSON(400, gin.H{"message": err})
			return
		}
		key := fmt.Sprintf("%s%s", params.Project, params.Stack)
		log.Infof("Recieved UnlockState request for project - stack: %s-%s", params.Project, params.Stack)

		entity, err := dao.Read(key)
		if err != nil {
			log.Error(err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		entity.Locked = false
		err = dao.Update(entity)
		if err != nil {
			log.Error(err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Status(200)
	})

	log.Infof("Starting server on port %d", cfg.Port)
	router.Run(fmt.Sprintf(":%d", cfg.Port))
}
