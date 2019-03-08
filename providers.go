package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/tengen-io/server/models"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type AuthProvider struct {
	db         *sqlx.DB
	signingKey []byte
	lifetime   time.Duration
}

func NewAuthProvider(db *sqlx.DB, signingKey []byte, lifetime time.Duration) *AuthProvider {
	return &AuthProvider{
		db:         db,
		signingKey: signingKey,
		lifetime:   lifetime,
	}
}

func (p *AuthProvider) ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return p.signingKey, nil
	})

	return token, err
}

func (p *AuthProvider) SignJWT(identity models.Identity) (string, error) {
	// TODO(eac): reintroduce custom claims for ID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Id:        identity.Id,
		NotBefore: time.Now().Unix(),
		ExpiresAt: time.Now().Add(p.lifetime * time.Second).Unix(),
		Issuer:    "tengen.io",
	})

	ss, err := token.SignedString(p.signingKey)
	if err != nil {
		return "", err
	}

	return ss, nil
}

// TODO(eac): Figure out how to use dbx structs for nested structures
func (p *AuthProvider) CheckPasswordByEmail(email, password string) (*models.Identity, error) {
	var passwordHash string
	err := p.db.Get(&passwordHash, "SELECT password_hash FROM identities WHERE email = $1", email)
	if err != nil {
		return nil, err
	}

	passwordBytes := []byte(password)
	hashBytes := []byte(passwordHash)

	err = bcrypt.CompareHashAndPassword(hashBytes, passwordBytes)
	if err != nil {
		return nil, err
	}

	var rv models.Identity
	row := p.db.QueryRowx("SELECT i.id, i.email, u.id, u.name FROM identities i, users u WHERE i.id = u.identity_id AND email = $1", email)
	err = row.Scan(&rv.Id, &rv.Email, &rv.User.Id, &rv.User.Name)

	if err != nil {
		return nil, err
	}

	return &rv, nil
}

type GameProvider struct {
	db *sqlx.DB
}

func NewGameProvider(db *sqlx.DB) *GameProvider {
	return &GameProvider{
		db,
	}
}

// TODO(eac): Add validation
// TODO(eac): Switch to sqlx binding
func (p *GameProvider) CreateInvitation(identity models.Identity, input models.CreateGameInvitationInput) (*models.Game, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, err
	}

	var rv models.Game
	ts := pq.FormatTimestamp(time.Now().UTC())

	game := tx.QueryRowx("INSERT INTO games (type, state, board_size, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, type, state, board_size", input.Type, models.Invitation, input.BoardSize, ts, ts)
	err = game.Scan(&rv.Id, &rv.Type, &rv.State, &rv.BoardSize)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = tx.Exec("INSERT INTO game_user (game_id, user_id, type, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)", rv.Id, identity.User.Id, models.GameUserEdgeTypeOwner, ts, ts)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &rv, nil
}

// TODO(eac): Validation? here or in the resolver
func (p *GameProvider) CreateGameUser(gameId string, userId string, edgeType models.GameUserEdgeType) (*models.JoinGamePayload, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, err
	}

	var rv models.Game
	ts := pq.FormatTimestamp(time.Now().UTC())
	_, err = tx.Exec("INSERT INTO game_user (game_id, user_id, type, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)", gameId, userId, edgeType, ts, ts)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	game := tx.QueryRow("SELECT id, type, state FROM games WHERE id = $1", gameId)
	err = game.Scan(&rv.Id, &rv.Type, &rv.State)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &models.JoinGamePayload{
		Game: &rv,
	}, nil
}

func (p *GameProvider) GetGameById(id string) (*models.Game, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var game models.Game
	err = p.db.Get(&game, "SELECT * FROM games WHERE id = $1", idInt)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (p *GameProvider) GetUsersForGame(id string) ([]*models.GameUserEdge, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var rv = make([]*models.GameUserEdge, 1)
	err = p.db.Select(&rv, "SELECT * FROM game_user WHERE game_id = $1", idInt)
	if err != nil {
		return nil, err
	}

	return rv, nil
}

func (p *GameProvider) GetGamesByIds(ids []string) ([]*models.Game, error) {
	idInts := make([]int, len(ids))
	for i, id := range ids {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		idInts[i] = idInt
	}

	query, fragArgs, err := sqlx.In("SELECT * FROM games WHERE id IN (?)", idInts)
	if err != nil {
		return nil, err
	}

	args := make([]interface{}, 0)
	args = append(args, fragArgs...)

	query = p.db.Rebind(query)
	rows, err := p.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rv := make([]*models.Game, 0)
	for rows.Next() {
		var i models.Game
		err := rows.StructScan(&i)
		if err != nil {
			return nil, err
		}
		rv = append(rv, &i)
	}

	return rv, nil
}

func (p *GameProvider) GetGamesByState(states []models.GameState) ([]*models.Game, error) {
	query, fragArgs, err := sqlx.In("SELECT * FROM games WHERE state IN (?)", states)
	if err != nil {
		return nil, err
	}

	args := make([]interface{}, 0)
	args = append(args, fragArgs...)

	query = p.db.Rebind(query)
	rows, err := p.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rv := make([]*models.Game, 0)
	for rows.Next() {
		var i models.Game
		err := rows.StructScan(&i)
		if err != nil {
			return nil, err
		}
		rv = append(rv, &i)
	}

	return rv, nil
}

type IdentityProvider struct {
	db         *sqlx.DB
	BcryptCost int
}

func NewIdentityProvider(db *sqlx.DB, bcryptCost int) *IdentityProvider {
	return &IdentityProvider{
		db,
		bcryptCost,
	}
}

func (p *IdentityProvider) CreateIdentity(email string, password string, name string) (*models.Identity, error) {
	// TODO(eac): re-add validation
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), p.BcryptCost)
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, err
	}

	var rv models.Identity
	ts := pq.FormatTimestamp(time.Now().UTC())

	// TODO(eac): do a precondition check for duplicate users to save autoincrement IDs
	identity := tx.QueryRowx("INSERT INTO identities (email, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id, email", email, passwordHash, ts, ts)
	err = identity.Scan(&rv.Id, &rv.Email)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	user := tx.QueryRowx("INSERT INTO users (identity_id, name, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id, name", rv.Id, name, ts, ts)
	err = user.Scan(&rv.User.Id, &rv.User.Name)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &rv, nil
}

// TODO(eac): switch to sqlx.get
func (p *IdentityProvider) GetIdentityById(id int32) (*models.Identity, error) {
	var identity models.Identity
	row := p.db.QueryRowx("SELECT i.id, i.email, u.id, u.name FROM identities i, users u WHERE i.id = u.identity_id AND i.id = $1", id)
	err := row.Scan(&identity.Id, &identity.Email, &identity.User.Id, &identity.User.Name)
	if err != nil {
		return nil, err
	}
	return &identity, nil
}

type UserProvider struct {
	db *sqlx.DB
}

func NewUserProvider(db *sqlx.DB) *UserProvider {
	return &UserProvider{
		db,
	}
}

// TODO(eac): switch to sqlx get
func (p *UserProvider) GetUserById(id string) (*models.User, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var rvId int
	var user models.User
	row := p.db.QueryRow("SELECT id, name FROM users WHERE id = $1", idInt)
	err = row.Scan(&rvId, &user.Name)
	user.Id = strconv.Itoa(rvId)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
