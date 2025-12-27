package common

import (
	"distributed-disk-register-with-grpc/internal/storage"
	"errors"
	"strconv"
	"strings"
)

type Command interface {
	Execute() (string, error)
}

type SetCommand struct {
	ID   int
	Text string
}

func (sc *SetCommand) Execute() (string, error) {
	err := storage.WriteMessage(sc.ID, sc.Text)
	if err != nil {
		return "", err
	}
	return "OK", nil
}

type GetCommand struct {
	ID int
}

func (gc *GetCommand) Execute() (string, error) {
	msg, err := storage.ReadMessage(gc.ID)
	if err != nil {
		return "NOT_FOUND", err
	}
	return msg, nil
}

func ParseCommand(msg string) (Command, error) {
	msg = strings.TrimSpace(msg)
	cmdParts := strings.Fields(msg)
	if len(cmdParts) < 2 {
		return nil, errors.New("Empty Command")
	}

	id, err := strconv.Atoi(cmdParts[1])
	if err != nil {
		return nil, errors.New("Invalid Id")
	}

	switch strings.ToUpper(cmdParts[0]) {
	case "SET":
		if len(cmdParts) < 3 {
			return nil, errors.New("Invalid Command")
		}
		cmd := SetCommand{
			ID:   id,
			Text: strings.Join(cmdParts[2:], " "),
		}
		return &cmd, nil
	case "GET":
		cmd := GetCommand{
			ID: id,
		}
		return &cmd, nil
	default:
		return nil, errors.New("Unknown Command")
	}

}
