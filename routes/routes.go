package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/QianMason/drone-cloud-api/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Services struct {
	TrackingIP string
	RegistryIP string
}

func setupCors(w *http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Header.Get("Origin"))
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CRSF-Token, Authorization")
}

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	trackingip := os.Getenv("TRACKINGIP")
	registryip := os.Getenv("REGISTRYIP")

	s := &Services{TrackingIP: trackingip, RegistryIP: registryip}
	fmt.Println("trackingip:", s.TrackingIP, "registryip:", s.RegistryIP)
	r.HandleFunc("/tracking", s.trackingPostHandler).Methods("POST")
	r.HandleFunc("/tracking/{id}", s.trackingGetHandler).Methods("GET")
	r.HandleFunc("/register", s.registerPostHandler).Methods("POST")
	return r
}

//theoretically, this should be able to access a cache of memory somewhere
func (s *Services) trackingGetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("tracking get handler called")
	url := "http://" + s.TrackingIP + r.URL.Path
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error getting from tracking component")
		fmt.Println("error:", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(body))
}

func (s *Services) trackingPostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("tracking post handler called")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("err tph bodyread:", err)
		return
	}
	//repackage request body into a new request
	url := "http://" + s.TrackingIP + r.URL.Path
	if err != nil {
		fmt.Println("json marshaling error")
		fmt.Println(err)
		return
	}
	//send second post request to tracking component
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("post request issue")
		fmt.Println(err)
		return
	}
	//read out response body from second post
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	//unmarshalling request body into structs

	fmt.Println("exiting create post handler")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	w.Write([]byte(body))

}

func (s *Services) registerPostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("register post handler called")
	/*
		/register
		FORWARDING REQUEST BY REPACKAGING BODY
	*/
	//read out request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("err tph bodyread:", err)
		return
	}
	//unmarshal body into a struct
	rd := &models.RegisterDrone{}

	err = json.Unmarshal(body, rd)
	if err != nil {
		fmt.Println("in here error unmarshalling:", err)
		return
	}
	//repackage request body into a new request for registration
	url := "http://" + s.RegistryIP + r.URL.Path
	if err != nil {
		fmt.Println("json marshaling error")
		fmt.Println(err)
		return
	}
	//send second post request to tracking component
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("post request issue")
		fmt.Println(err)
		return
	}
	//read out response body from second post
	secondBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	//body should be the droneID
	dID := string(secondBody)
	fmt.Println("dID:", dID)
	/*
		/create
		CREATING OBJECT TO MARSHAL AND SEND
	*/
	trackingDevice := models.TrackingDevice{
		DroneID: dID,
		Lat:     rd.Lat,
		Lng:     rd.Lng,
	}
	//repackage request body into a new request for registration
	trackingCreateURL := "http://" + s.TrackingIP + "/create"
	jsn, err := json.Marshal(trackingDevice)

	//send second post request to tracking component
	trackingResp, err := http.Post(trackingCreateURL, "application/json", bytes.NewBuffer(jsn)
	if err != nil {
		fmt.Println("tracking resp issue")
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	fmt.Println("exiting register post handler")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	w.Write(secondBody)
}

// func (s *Services) redirectRequestBody(host string, endpoint string)

// func (db *DroneDB) droneHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("in index handler")
// 	setupCors(&w, r)
// 	if r.Method == "OPTIONS" {
// 		return
// 	}
// 	clientOptions := options.Client().
// 		ApplyURI("mongodb+srv://thunderpurtz:" + db.password + "@cluster0.14i4y.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	client, err := mongo.Connect(ctx, clientOptions)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer func() {
// 		if err = client.Disconnect(ctx); err != nil {
// 			panic(err)
// 		}
// 	}()
// 	collection := client.Database("DronePlatform").Collection("trackingData")
// 	cur, err := collection.Find(ctx, bson.D{})
// 	if err != nil {
// 		fmt.Println(cur)
// 		fmt.Println("error with cur")
// 		fmt.Println(err)
// 		return
// 	}
// 	drones := []models.Drone{}

// 	for cur.Next(ctx) {
// 		d := models.Drone{}
// 		err = cur.Decode(&d)
// 		fmt.Println("d lat:", d.Coordinates.Lat, ":d lng:", d.Coordinates.Lng)
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 		drones = append(drones, d)
// 	}
// 	cur.Close(ctx)
// 	if len(drones) == 0 {
// 		w.WriteHeader(500)
// 		w.Write([]byte("No data found."))
// 		return
// 	}
// 	jsn, err := json.Marshal(drones)
// 	fmt.Println("jsn:", jsn)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(200)
// 	w.Write(jsn)
// }

// func getUUID() string {
// 	uid := strings.Replace(uuid.New().String(), "-", "", -1)
// 	fmt.Println("New UUID:", uid)
// 	return uid
// }

// func userGetHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("userget" + strconv.Itoa(globalcount))
// 	globalcount += 1
// 	session, _ := sessions.Store.Get(r, "session")
// 	untypedUserId := session.Values["user_id"]
// 	currentUserId, ok := untypedUserId.(int64)
// 	fmt.Println(currentUserId)
// 	if !ok {
// 		utils.InternalServerError(w)
// 		return
// 	}
// 	vars := mux.Vars(r) //hashmap of variable names and content passed for that variable
// 	username := vars["username"]
// 	fmt.Println("username", username)

// 	currentPageUserString := strings.TrimLeft(r.URL.Path, "/")
// 	currentPageUser, err := models.GetUserByUsername(currentPageUserString)
// 	if err != nil {
// 		utils.InternalServerError(w)
// 		return
// 	}
// 	currentPageUserID, err := currentPageUser.GetId()
// 	if err != nil {
// 		utils.InternalServerError(w)
// 		return
// 	}
// 	updates, err := models.GetUpdates(currentPageUserID)
// 	if err != nil {
// 		utils.InternalServerError(w)
// 		return
// 	}

// 	utils.ExecuteTemplate(w, "index.html", struct {
// 		Title       string
// 		Updates     []*models.Update
// 		DisplayForm bool
// 	}{
// 		Title:       username,
// 		Updates:     updates,
// 		DisplayForm: currentPageUserID == currentUserId,
// 	})

// }

// func indexHandler(w http.ResponseWriter, r *http.Request) {
// 	updates, err := models.GetAllUpdates()
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte("Internal server error"))
// 		return
// 	}
// 	utils.ExecuteTemplate(w, "index.html", struct {
// 		Title       string
// 		Updates     []*models.Update
// 		DisplayForm bool
// 	}{
// 		Title:       "All updates",
// 		Updates:     updates,
// 		DisplayForm: true,
// 	})
// 	fmt.Println("get")
// }

// func postHandlerHelper(w http.ResponseWriter, r *http.Request) error {
// 	session, _ := sessions.Store.Get(r, "session")
// 	untypedUserID := session.Values["user_id"]
// 	userID, ok := untypedUserID.(int64)
// 	if !ok {
// 		return utils.InternalServer
// 	}
// 	currentPageUserString := strings.TrimLeft(r.URL.Path, "/")
// 	currentPageUser, err := models.GetUserByUsername(currentPageUserString)
// 	if err != nil {
// 		return utils.InternalServer
// 	}
// 	currentPageUserID, err := currentPageUser.GetId()
// 	if err != nil {
// 		return utils.InternalServer
// 	}
// 	if currentPageUserID != userID {
// 		return utils.BadPostError
// 	}
// 	r.ParseForm()
// 	body := r.PostForm.Get("adddrone")
// 	fmt.Println(body)
// 	err = models.PostUpdates(userID, body)
// 	if err != nil {
// 		return utils.InternalServer
// 	}
// 	return nil
// }

// func postHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("post handler called")
// 	err := postHandlerHelper(w, r)
// 	if err == utils.InternalServer {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte("Internal server error"))
// 	}
// 	http.Redirect(w, r, "/", 302)
// }

// func UserPostHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("user post handler called")
// 	fmt.Println(r.URL.Path)
// 	err := postHandlerHelper(w, r)
// 	if err == utils.BadPostError {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte("Cannot write to another user's page"))
// 	}
// 	http.Redirect(w, r, r.URL.Path, 302)
// }

// func loginGetHandler(w http.ResponseWriter, r *http.Request) {
// 	utils.ExecuteTemplate(w, "login.html", nil)
// }

// func loginPostHandler(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()
// 	username := r.PostForm.Get("username")
// 	password := r.PostForm.Get("password")

// 	user, err := models.AuthenticateUser(username, password)
// 	if err != nil {
// 		switch err {
// 		case models.InvalidLogin:
// 			utils.ExecuteTemplate(w, "login.html", "User or Pass Incorrect")
// 		default:
// 			w.WriteHeader(http.StatusInternalServerError)
// 			w.Write([]byte("Internal server error"))
// 		}
// 		return
// 	}
// 	userId, err := user.GetId()
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte("Internal server error"))
// 		return
// 	}
// 	sessions.GetSession(w, r, "session", userId)
// 	http.Redirect(w, r, "/", 302)
// }

// func logoutGetHandler(w http.ResponseWriter, r *http.Request) {
// 	sessions.EndSession(w, r)
// 	http.Redirect(w, r, "/login", 302)
// }

// func registerGetHandler(w http.ResponseWriter, r *http.Request) {
// 	utils.ExecuteTemplate(w, "register.html", nil)
// }

// func registerPostHandler(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()
// 	username := r.PostForm.Get("username")
// 	password := r.PostForm.Get("password")
// 	err := models.RegisterUser(username, password)
// 	if err == models.UserNameTaken {
// 		utils.ExecuteTemplate(w, "register.html", "username taken")
// 		return
// 	}
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte("Internal server error"))
// 		return
// 	}
// 	http.Redirect(w, r, "/login", 302)
// }
