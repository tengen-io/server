package resolvers

import (
    "github.com/camirmas/go_stop/models"
    "github.com/graphql-go/graphql"
)

func GetUser(p graphql.ResolveParams) (interface{}, error) {
    db := p.Context.Value("db").(models.Database)
    return db.GetUser(p.Args["username"].(string))
}

func CreateUser(p graphql.ResolveParams) (interface{}, error) {
    db := p.Context.Value("db").(models.Database)
    username := p.Args["username"].(string)
    email := p.Args["email"].(string)
    password := p.Args["password"].(string)
    passwordConfirm := p.Args["passwordConfirmation"].(string)

    user, err := db.CreateUser(username, email, password, passwordConfirm)

    if err != nil {
        return nil, err
    }

    token, err := GenerateToken(user.Id)

    return &models.AuthUser{token, user}, nil
}

func LogIn(p graphql.ResolveParams) (interface{}, error) {
    username := p.Args["username"].(string)
    password := p.Args["password"].(string)
    db := p.Context.Value("db").(models.Database)
    user, err := db.CheckPw(username, password)

    if err != nil {
        return nil, err
    }

    token, err := GenerateToken(user.Id)

    if err != nil {
        return user, err
    }

    authUser := &models.AuthUser{token, user}

    return authUser, nil
}

func CurrentUser(p graphql.ResolveParams) (interface{}, error) {
    db := p.Context.Value("db").(models.Database)

    token, ok := p.Context.Value("token").(string)

    if !ok {
        return nil, missingTokenError{}
    }

    claims, err := ValidateToken(token)

    if err != nil {
        return nil, err
    }

    return db.GetUser(claims.UserId)
}
