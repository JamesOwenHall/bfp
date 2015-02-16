package dashboard

import (
	"encoding/json"
	"net/http"
)

func (s *Server) HandleHistory(w http.ResponseWriter, r *http.Request) {
	data := make([]interface{}, 0, 2*len(s.conf.Directions))

	for iDirection := range s.conf.Directions {
		direction := &s.conf.Directions[iDirection]
		history := &direction.History

		dirData := map[string]interface{}{
			"name":           direction.Name,
			"short-history":  history.Short.Read(),
			"long-history":   history.Long.Read(),
			"blocked-values": direction.Store.BlockedValues(),
		}

		data = append(data, dirData)
	}

	json, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
