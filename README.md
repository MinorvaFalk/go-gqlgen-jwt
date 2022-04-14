# go-gqlgen-jwt
Go using gqlgen and JWT Authentication\
**!!! Reminder, this might not the best practice for JWT Implementation !!!**

### External Library Used :
- Gqlgen [github.com/99designs/gqlgen]
- JWT for Authentication [github.com/golang-jwt/jwt]
- *Gin http Framework [github.com/gin-gonic/gin]

> *optional
> you can prefer another library other than gin

### Project Structure :
```
├───graph
│   ├───generated
│   └───model
├───middleware
└───utils
```

### Authentication process :
1. Register `Gin Context` for gqlgen.\
Refer [https://gqlgen.com/recipes/gin/] for more information.\
2. Create `schema directive` inside .graphql file.\
Refer [https://gqlgen.com/reference/directives/] or open `/graph/schema.graphqls`\
example code:
```graphql
directive @mustAuth(auth: Boolean!) on FIELD_DEFINITION
```
4. Run `go generate ./...` to generate directive.
5. Implement directive inside `graphql controller`.
```go
func graphqlHandler(middleware *middleware.Middlewares) gin.HandlerFunc {
	c := generated.Config{ Resolvers: &graph.Resolver{} }
  
  // Implement your directive
	c.Directives.MustAuth = func(ctx context.Context, obj interface{}, next graphql.Resolver, auth bool) (res interface{}, err error) {
		// Do something...

		return next(ctx)
	}

	h := handler.NewDefaultServer(generated.NewExecutableSchema(c))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
```
7. Add `JWT Validation` inside your directive function.\
**First** you need to get Gin Context.\
example: 
```go
c.Directives.MustAuth = func(ctx context.Context, obj interface{}, next graphql.Resolver, auth bool) (res interface{}, err error) {
  gc, err := GinContextFromContext(ctx)
  if err != nil {
    return nil, err
  }

  // Do something...
 }
```
**Second** check JWTValidation using gin context.\
example:
```go
c.Directives.MustAuth = func(ctx context.Context, obj interface{}, next graphql.Resolver, auth bool) (res interface{}, err error) {
  // Pretend we already get the context.
  gc := contextFromGin

  checkAuth := JwtMiddleware(gc)
  if checkAuth != nil {
    return nil, checkAuth
  }

  return next(ctx)
}

// JwtMiddlewareFunction
func JwtMiddleware(c *gin.Context) error {
  header := c.GetHeader("Authorization")
  if !strings.Contains(header, "Bearer") {
    return errors.New("invalid token supplied")
  }

  tokenString := header[len("Bearer "):]	

  // Handle your validation here... 
}
```
### Notes:
This might not be the best implementation, but since you need to implement auth per query this might be the way.
Feel free to create any `issue` if there is any problem
