package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type GameState int8

const (
	GameStateNegotiation GameState = iota
	GameStateInProgress
	GameStateFinished
)

func GameStateForString(str string) (GameState, error) {
	switch str {
	case "NEGOTIATION":
		return GameStateNegotiation, nil
	case "IN_PROGRESS":
		return GameStateInProgress, nil
	case "FINISHED":
		return GameStateFinished, nil
	default:
		return 0, fmt.Errorf("unknown gamestate %s", str)
	}
}

func (g GameState) String() string {
	switch g {
	case GameStateNegotiation:
		return "NEGOTIATION"
	case GameStateInProgress:
		return "IN_PROGRESS"
	case GameStateFinished:
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
