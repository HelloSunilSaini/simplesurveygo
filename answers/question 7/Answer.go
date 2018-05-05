//it is in main.go



//---------------------------------------------------------------------------------------------

package main

import (
	"time"
	"fmt"
	"net/http"
	"simplesurveygo/servicehandlers"
	"simplesurveygo/dao"
	"reflect"
	"gopkg.in/mgo.v2/bson"
)

func main() {

	// Serves the html pages
	http.Handle("/", http.FileServer(http.Dir("./static")))

	pingHandler := servicehandlers.PingHandler{}
	authHandler := servicehandlers.UserValidationHandler{}
	sessionHandler := servicehandlers.SessionHandler{}
	surveyHandler := servicehandlers.SurveyHandler{}
	userSurveyHandler := servicehandlers.UserSurveyHandler{}
	signupHandler := servicehandlers.SignupHandler{}

	// Serves the API content
	http.Handle("/api/v1/ping/", pingHandler)

	http.Handle("/api/v1/signup/",  timerAndAuth(signupHandler))
	http.Handle("/api/v1/authenticate/",  timerAndAuth(authHandler))
	http.Handle("/api/v1/validate/", timerAndAuth(sessionHandler))

	http.Handle("/api/v1/survey/{surveyname}", timerAndAuth(surveyHandler))
	http.Handle("/api/v1/survey/", surveyHandler)

	http.Handle("/api/v1/usersurvey/", timerAndAuth(userSurveyHandler))

	// Start Server
	http.ListenAndServe(":3000", nil)
}


func timerAndAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		tokens, _ := r.Header["Token"]
		if r.Method == "OPTIONS"{
			formatted := bson.M{
				"responseData": "{}",
				"message": "UnAuthorized OPTION call",
				"status":       true}
			data, _ := servicehandlers.MarshalResponse(formatted)
			w.Header().Set("Content-Length", fmt.Sprint(len(data)))
			w.WriteHeader(401)
			fmt.Fprint(w, string(data))

		}else if (len(tokens) > 0){
			token := tokens[0]
			user := dao.GetSessionDetails(token)
			if (reflect.DeepEqual(user,(dao.UserCredentials{}))){
				formatted := bson.M{
					"responseData": "{}",
					"message": "UnAuthorized",
					"status":       true}
				data, _ := servicehandlers.MarshalResponse(formatted)
				w.Header().Set("Content-Length", fmt.Sprint(len(data)))
				w.WriteHeader(401)
				fmt.Fprint(w, string(data))
			}
			
		}else{
			h.ServeHTTP(w, r)
		}
		duration := time.Now().Sub(startTime)
		fmt.Println(duration)
	})
}