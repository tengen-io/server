package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/tengen-io/server/models"
	"github.com/tengen-io/server/pubsub"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type AuthRepository struct {
	db         *sqlx.DB
	signingKey []byte
	lifetime   time.Duration
}

func NewAuthRepository(db *sqlx.DB, signingKey []byte, lifetime time.Duration) *AuthRepository {
	return &AuthRepository{
		db:         db,
		signingKey: signingKey,
		lifetime:   lifetime,
	}
}

func (p *AuthRepository) ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return p.signingKey, nil
	})

	return token, err
}

func (p *AuthRepository) SignJWT(identity models.Identity) (string, error) {
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
func (p *AuthRepository) CheckPasswordByEmail(email, password string) (*models.Identity, error) {
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

type GameRepository struct {
	db     *sqlx.DB
	pubsub pubsub.Bus
}

func NewGameRepository(db *sqlx.DB, pubsub pubsub.Bus) *GameRepository {
	return &GameRepository{
		db, pubsub,
	}
}

// TODO(eac): Add validation
// TODO(eac): Switch to sqlx binding
func (p *GameRepository) CreateGame(identity models.Identity, gameType models.GameType, boardSize int, gameState models.GameState) (*models.Game, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, err
	}

	var rv models.Game
	ts := pq.FormatTimestamp(time.Now().UTC())

	game := tx.QueryRowx("INSERT INTO games (type, state, board_size, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, type, state, board_size", gameType, gameState, boardSize, ts, ts)
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

	payload := pubsub.Event{
		Event: "create",
		Payload: &rv,
	}

	p.pubsub.Publish(payload, "game", "game_"+rv.State.String())

	return &rv, nil
}

// TODO(eac): Validation? here or in the resolver
// TODO(eac): maybe make this return a gameUser, and have the resolver get the game?
func (p *GameRepository) CreateGameUser(gameId string, userId string, edgeType models.GameUserEdgeType) (*models.Game, error) {
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

	return &rv, nil
}

func (p *GameRepository) GetGameById(id string) (*models.Game, error) {
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

func (p *GameRepository) GetUsersForGame(id string) ([]models.GameUserEdge, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	rows, err := p.db.Query("SELECT type, user_id, name FROM game_user gu, users u WHERE game_id = $1 AND gu.user_id = u.id", idInt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rv = make([]models.GameUserEdge, 0)
	for rows.Next() {
		var i models.GameUserEdge
		err := rows.Scan(&i.Type, &i.User.Id, &i.User.Name)
		if err != nil {
			return nil, err
		}

		rv = append(rv, i)
	}

	return rv, nil
}

func (p *GameRepository) GetGamesByIds(ids []string) ([]*models.Game, error) {
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

func (p *GameRepository) GetGamesByState(states []models.GameState) ([]*models.Game, error) {
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

type IdentityRepository struct {
	db         *sqlx.DB
	BcryptCost int
}

func NewIdentityRepository(db *sqlx.DB, bcryptCost int) *IdentityRepository {
	return &IdentityRepository{
		db,
		bcryptCost,
	}
}

func (p *IdentityRepository) CreateIdentity(email string, password string, name string) (*models.Identity, error) {
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
func (p *IdentityRepository) GetIdentityById(id int32) (*models.Identity, error) {
	var identity models.Identity
	row := p.db.QueryRowx("SELECT i.id, i.email, u.id, u.name FROM identities i, users u WHERE i.id = u.identity_id AND i.id = $1", id)
	err := row.Scan(&identity.Id, &identity.Email, &identity.User.Id, &identity.User.Name)
	if err != nil {
		return nil, err
	}
	return &identity, nil
}

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db,
	}
}

// TODO(eac): switch to sqlx get
func (p *UserRepository) GetUserById(id string) (*models.User, error) {
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
