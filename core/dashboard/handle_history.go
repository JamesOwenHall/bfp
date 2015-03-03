package dashboard

import (
	"encoding/json"
	"net/http"
)

func (s *Server) HandleHistory(w http.ResponseWriter, r *http.Request) {
	data := make([]interface{}, 0, len(s.conf.Directions))

	for iDirection := range s.conf.Directions {
		direction := &s.conf.Directions[iDirection]

		dirData := map[string]interface{}{
			"name":           direction.Name,
			"blocked-values": direction.Store.BlockedValues(),
			"clock":          s.counter.Clock.GetTime(),
		}

		data = append(data, dirData)
	}

	json, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
