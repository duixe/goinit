package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/duixe/go_rest/service/cart"
	"github.com/duixe/go_rest/service/order"
	"github.com/duixe/go_rest/service/product"
	"github.com/duixe/go_rest/service/user"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

// create an APIServer receiver that'll return new instance of the APIServer
func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error{
	router := mux.NewRouter()
	//create a subrouter which will then be passed instead of the router
	//in order to prefix the api routes with "/api/v1"
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(subrouter)

	orderStore := order.NewStore(s.db)

	cartHandler := cart.NewHandler(orderStore, productStore, userStore)
	cartHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
