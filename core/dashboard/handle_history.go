package dashboard

import (
	"encoding/json"
	"net/http"
)

func (s *Server) HandleHistory(w http.ResponseWriter, r *http.Request) {
	type HistoryData struct {
		Clock      int32
		Directions []interface{}
	}

	data := HistoryData{
		Clock:      s.counter.Clock.GetTime(),
		Directions: make([]interface{}, 0, len(s.conf.Directions)),
	}

	for iDirection := range s.conf.Directions {
		direction := &s.conf.Directions[iDirection]

		dirData := map[string]interface{}{
			"name":           direction.Name,
			"blocked-values": direction.Store.BlockedValues(),
		}

		data.Directions = append(data.Directions, dirData)
	}

	json, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
