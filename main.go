package main

import (
	"CSI2132/controllers"
	"CSI2132/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Printf("Starting up server...\n")

	us, err := models.NewUserService() // UserService is created and a connection to the database should be opened
	if err != nil {
		panic(err)
	}
	defer us.Close() // Ensures that database connection closes no matter what

	// Check if connection is successful
	if err = us.Ping(); err != nil {
		panic(err)
	} else {
		fmt.Printf("Connection to db was successful\n")
	}

	// Instantiating controller types
	landingController := controllers.NewLanding()
	dashboardController := controllers.NewDashboard()
	usersController := controllers.NewUsers(us)

	// router is created
	router := mux.NewRouter()

	// Serving all static CSS/JS files needed for styling our application
	bootstrapCSS := http.FileServer(http.Dir("./static/bootstrap/css/"))
	router.PathPrefix("/bootstrap/css/").Handler(http.StripPrefix("/bootstrap/css/", bootstrapCSS))

	customCSS := http.FileServer(http.Dir("./static/css/"))
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", customCSS))

	jqry := http.FileServer(http.Dir("./static/jquery/"))
	router.PathPrefix("/jquery/").Handler(http.StripPrefix("/jquery/", jqry))

	bootstrapJS := http.FileServer(http.Dir("./static/bootstrap/js/"))
	router.PathPrefix("/bootstrap/js/").Handler(http.StripPrefix("/bootstrap/js/", bootstrapJS))

	images := http.FileServer(http.Dir("./static/images/"))
	router.PathPrefix("/images/").Handler(http.StripPrefix("/images/", images))
	us.GetProperty()
	// Routes handled by landingController
	router.Handle("/", landingController.WelcomeView).Methods("GET")

	// Routes for handling dashboards
	router.Handle("/home", dashboardController.HomeView).Methods("GET")
	router.Handle("/guest", dashboardController.GuestView).Methods("GET")
	router.Handle("/host", dashboardController.HostView).Methods("GET")

	// routes handled by users controller
	router.HandleFunc("/signup", usersController.New).Methods("GET")
	router.HandleFunc("/signup", usersController.Create).Methods("POST")

	router.Handle("/signin", usersController.SignInView).Methods("GET")
	router.HandleFunc("/signout", usersController.SignOut).Methods("GET")
	router.HandleFunc("/signin", usersController.SignIn).Methods("POST")

	router.Handle("/property", usersController.PropertyView).Methods("GET")
	router.HandleFunc("/property", usersController.CreateProperty).Methods("POST")

	router.Handle("/search", usersController.Search).Methods("GET")

	fmt.Printf("Server has started...\n")

	// this starts the server
	log.Fatal(http.ListenAndServe(":8000", router))
}
