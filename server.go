package main

import (
	"context"
	"fmt"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/MinorvaFalk/go-gqlgen-jwt/graph"
	"github.com/MinorvaFalk/go-gqlgen-jwt/graph/generated"
	"github.com/MinorvaFalk/go-gqlgen-jwt/middleware"
	"github.com/MinorvaFalk/go-gqlgen-jwt/utils"
	"github.com/gin-gonic/gin"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	jwtAuth := utils.NewJwtAuth()
	middleware := middleware.NewMiddlewares(jwtAuth)

	// Setting up Gin
	r := gin.Default()

	r.Use(middleware.GinContextToContextMiddleware())

	r.POST("/query", graphqlHandler(middleware))
	r.GET("/", playgroundHandler())
	r.Run(":" + port)
}

// Defining the Graphql handler
func graphqlHandler(middleware *middleware.Middlewares) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	c := generated.Config{
		Resolvers: &graph.Resolver{
			Middleware: middleware,
		},
	}
	c.Directives.MustAuth = func(ctx context.Context, obj interface{}, next graphql.Resolver, auth bool) (res interface{}, err error) {
		gc, err := GinContextFromContext(ctx)
		if err != nil {
			return nil, err
		}

		checkAuth := middleware.JwtMiddlewareV2(gc)
		if checkAuth != nil {
			return nil, checkAuth
		}

		return next(ctx)
	}

	h := handler.NewDefaultServer(generated.NewExecutableSchema(c))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value("GinContextKey")
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}
