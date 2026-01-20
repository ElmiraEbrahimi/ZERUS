package server

import (
	"net/http"

	"l2alchemy/internal/server/handlers"
)

// NewRouter constructs a new http.ServeMux and registers all routes using the
// provided handlers. This helper is useful in tests and when embedding the
// router in a custom HTTP server.
func NewRouter(counter *handlers.CounterHandler, zk *handlers.ZKHandler, oracle *handlers.OracleHandler, messaging *handlers.MessengerHandler) *http.ServeMux {
	mux := http.NewServeMux()
	RegisterRoutes(mux, counter, zk, oracle, messaging)
	return mux
}

// RegisterRoutes binds HTTP paths to handler methods. It can be called from
// main() to configure the default ServeMux.
func RegisterRoutes(mux *http.ServeMux, counter *handlers.CounterHandler, zk *handlers.ZKHandler, oracle *handlers.OracleHandler, messaging *handlers.MessengerHandler) {
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

		// GET /users/default/balance returns the default user's balance.
		mux.HandleFunc("/users/default/balance", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			oracle.GetDefaultUserBalance(w, r)
		})

		// POST /users/default/burn triggers a burn for the default user.
		mux.HandleFunc("/users/default/burn", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			oracle.BurnDefaultUser(w, r)
		})

		// POST /users/default/withdraw triggers a withdraw for the default user.
		mux.HandleFunc("/users/default/withdraw", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			oracle.WithdrawDefaultUser(w, r)
		})

		// POST /validators/register registers all validator nodes on-chain.
		mux.HandleFunc("/validators/register", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			oracle.RegisterValidators(w, r)
		})

		// POST /validators/replace triggers ReplaceAccountTx for a validator node.
		mux.HandleFunc("/validators/replace", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			oracle.ReplaceValidatorAccount(w, r)
		})

		// POST /validators/exit triggers ExitTx for a validator node.
		mux.HandleFunc("/validators/exit", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			oracle.ExitValidatorAccount(w, r)
		})

		// POST /validators/withdraw triggers WithdrawAccountTx for a validator node.
		mux.HandleFunc("/validators/withdraw", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			oracle.WithdrawValidatorAccount(w, r)
		})
	}

	if messaging != nil {
		mux.HandleFunc("/messaging/l2/send", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			messaging.SendL2ToL1(w, r)
		})

		mux.HandleFunc("/messaging/l2/last-from-l1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			messaging.GetL2LastFromL1(w, r)
		})

		mux.HandleFunc("/messaging/l1/send", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			messaging.SendL1ToL2(w, r)
		})

		mux.HandleFunc("/messaging/l1/send-direct", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			messaging.SendL1ToL2Direct(w, r)
		})

		mux.HandleFunc("/messaging/l1/receive", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			messaging.ReceiveFromL2(w, r)
		})

		mux.HandleFunc("/messaging/l1/last-from-l2", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			messaging.GetL1LastFromL2(w, r)
		})
	}
}
