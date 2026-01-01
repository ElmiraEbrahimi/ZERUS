package server

import (
	"net/http"

	"l2alchemy/internal/server/handlers"
)

// NewRouter constructs a new http.ServeMux and registers all routes using the
// provided handlers. This helper is useful in tests and when embedding the
// router in a custom HTTP server.
func NewRouter(counter *handlers.CounterHandler, zk *handlers.ZKHandler, oracle *handlers.OracleHandler) *http.ServeMux {
	mux := http.NewServeMux()
	RegisterRoutes(mux, counter, zk, oracle)
	return mux
}

// RegisterRoutes binds HTTP paths to handler methods. It can be called from
// main() to configure the default ServeMux.
func RegisterRoutes(mux *http.ServeMux, counter *handlers.CounterHandler, zk *handlers.ZKHandler, oracle *handlers.OracleHandler) {
	if counter != nil {
		mux.HandleFunc("/counter", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			counter.GetCounter(w, r)
		})

		mux.HandleFunc("/counter/increment", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			counter.IncrementCounter(w, r)
		})

		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			counter.Health(w, r)
		})
	}

	if zk != nil {
		// ZK key generation and status endpoints.
		// POST /circuits/keygen triggers Groth16 setup and writes pk/vk for both circuits.
		mux.HandleFunc("/circuits/keygen", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			zk.GenerateCircuitKeys(w, r)
		})

		// GET /circuits/keys returns whether pk/vk exist on disk for each circuit.
		mux.HandleFunc("/circuits/keys", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			zk.CircuitKeysStatus(w, r)
		})
	}

	if oracle != nil {
		// POST /users/default/register registers the default user on-chain.
		mux.HandleFunc("/users/default/register", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			oracle.RegisterDefaultUser(w, r)
		})

		// POST /validators/register registers all validator nodes on-chain.
		mux.HandleFunc("/validators/register", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			oracle.RegisterValidators(w, r)
		})
	}
}
