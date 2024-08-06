package api

import (
	"ask-me-anything/internal/store/pgstore"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (h apiHandler) readRoom(w http.ResponseWriter, r *http.Request) (room pgstore.Room, rawRoomId string, roomID uuid.UUID, ok bool) {
	rawRoomId = chi.URLParam(r, "room_id")
	roomID, err := uuid.Parse(rawRoomId)
	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return pgstore.Room{}, "", uuid.UUID{}, false
	}

	room, err = h.q.GetRoom(r.Context(), roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "room not found", http.StatusNotFound)
			return pgstore.Room{}, "", uuid.UUID{}, false
		}
		setInternalServerError(w)
		return pgstore.Room{}, "", uuid.UUID{}, false
	}

	return room, rawRoomId, roomID, true
}

func (h apiHandler) readRoomMessage(w http.ResponseWriter, r *http.Request) (message pgstore.Message, rawMessageId string, messageID uuid.UUID, rawRoomId string, ok bool) {
	_, rawRoomId, _, ok = h.readRoom(w, r)
	if !ok {
		return pgstore.Message{}, "", uuid.UUID{}, "", ok
	}
	rawMessageId = chi.URLParam(r, "message_id")
	messageID, err := uuid.Parse(rawMessageId)
	if err != nil {
		http.Error(w, "invalid message id", http.StatusBadRequest)
		return pgstore.Message{}, "", uuid.UUID{}, rawRoomId, false
	}

	message, err = h.q.GetMessage(r.Context(), messageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "message not found", http.StatusNotFound)
			return pgstore.Message{}, "", uuid.UUID{}, rawRoomId, false
		}
		setInternalServerError(w)
		return pgstore.Message{}, "", uuid.UUID{}, rawRoomId, false
	}

	return message, rawMessageId, messageID, rawRoomId, true
}

func sendJSON(w http.ResponseWriter, rawData any) {
	data, _ := json.Marshal(rawData)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func setInternalServerError(w http.ResponseWriter) {
	http.Error(w, "something went wrong", http.StatusInternalServerError)
}
