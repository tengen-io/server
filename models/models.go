package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"io"
	"strconv"
	"time"
)

type NodeFields struct {
	Id        string `json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Identity struct {
	NodeFields
	Email string `json:"email"`
	User
}

func (Identity) IsNode() {}

type User struct {
	NodeFields
	Name string `json:"name"`
}

func (User) IsNode() {}

type GameType int8

const (
	Standard GameType = iota
)

func (g GameType) String() string {
	switch g {
	case Standard:
		return "STANDARD"
	default:
		return "UNKNOWN"
	}
}

func GameTypeForString(str string) (GameType, error) {
	switch str {
	case "STANDARD":
		return Standard, nil
	default:
		return 0, fmt.Errorf("unknown gametype %s", str)
	}
}

func (g *GameType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return errors.New("cannot unmarshal non-string as GameType")
	}

	var err error
	*g, err = GameTypeForString(str)
	return err
}

func (g GameType) MarshalGQL(w io.Writer) {
	fmt.Fprintf(w, strconv.Quote(g.String()))
}

func (g *GameType) Scan(value interface{}) error {
	val, ok := value.([]byte)
	if !ok {
		return errors.New("cannot scan non-string as gametype")
	}

	var err error
	*g, err = GameTypeForString(string(val))
	if err != nil {
		return err
	}
	return nil
}

func (g GameType) Value() (driver.Value, error) {
	return g.String(), nil
}

type GameState int8

const (
	Invitation GameState = iota
	InProgress
	Finished
)

func GameStateForString(str string) (GameState, error) {
	switch str {
	case "INVITATION":
		return Invitation, nil
	case "IN_PROGRESS":
		return InProgress, nil
	case "FINISHED":
		return Finished, nil
	default:
		return 0, fmt.Errorf("unknown gamestate %s", str)
	}
}

func (g GameState) String() string {
	switch g {
	case Invitation:
		return "INVITATION"
	case InProgress:
		return "IN_PROGRESS"
	case Finished:
		return "FINISHED"
	default:
		return "UNKNOWN"
	}
}

func (g *GameState) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return errors.New("cannot unmarshal non-string as GameState")
	}

	var err error
	*g, err = GameStateForString(str)
	return err
}

func (g GameState) MarshalGQL(w io.Writer) {
	fmt.Fprintf(w, strconv.Quote(g.String()))
}

func (g *GameState) Scan(value interface{}) error {
	val, ok := value.([]byte)
	if !ok {
		return errors.New("cannot scan non-[]byte as gamestate")
	}

	var err error
	*g, err = GameStateForString(string(val))
	if err != nil {
		return nil
	}
	return nil
}

func (g GameState) Value() (driver.Value, error) {
	return g.String(), nil
}

type Game struct {
	NodeFields
	BoardSize int       `json:"boardSize",db:"board_size"`
	Type      GameType  `json:"type"`
	State     GameState `json:"state"`
}

func (Game) IsNode() {}

func MarshalTimestamp(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.FormatInt(t.UTC().Unix(), 10))
	})
}

func UnmarshalTimestamp(v interface{}) (time.Time, error) {
	if conv, ok := v.(int64); ok {
		return time.Unix(conv, 0), nil
	}

	return time.Time{}, errors.New("could not convert timestamp to int64")
}
