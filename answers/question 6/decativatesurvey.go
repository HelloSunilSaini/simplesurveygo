// -- code in srevicehandlers servey.go  -> Post call and dao survey.go


// ------------------------- dao survey.go -----------------------------------------

package dao

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

type Question struct {
	QuestionString string   `json:"questionString" bson:"questionString"`
	Options        []string `json:"options" bson:"options"`
}

type Answer struct {
	Question Question `json:"question" bson:"question"`
	Answer   string   `json:"answer" bson:"answer"`
}

type Survey struct {
	SurveyName  string     `json:"surveyName" bson:"surveyName"`
	Heading     string     `json:"heading" bson:"heading"`
	Description string     `json:"description" bson:"description"`
	Questions   []Question `json:"questions" bson:"questions"`
	Status      bool       `json:"status" bson:"status"`
}

type SurveyResponse struct {
	UserName string   `json:"userName" bson:"userName"`
	Survey   Survey   `json:"survey" bson:"survey"`
	Answers  []Answer `json:"answers" bson:"answers"`
}

func GetActiveSurveys() interface{} {
	session := MgoSession.Clone()
	defer session.Close()

	var response []interface{}
	clctn := session.DB("simplesurveys").C("survey")
	query := clctn.Find(bson.M{"status": true})
	err := query.All(&response)

	if err != nil {
		return nil
	} else {
		return response
	}
}

func GetSurveysForUser(userName string) interface{} {
	session := MgoSession.Clone()
	defer session.Close()

	var response []interface{}
	clctn := session.DB("simplesurveys").C("survey_response")
	query := clctn.Find(bson.M{"userName": userName})
	err := query.All(&response)

	if err != nil {
		return nil
	} else {
		return response
	}
}

func GetSurveyByName(surveyName string) interface{} {
	fmt.Println("GetSurveyByName:" + surveyName)
	session := MgoSession.Clone()
	defer session.Close()

	var response interface{}
	clctn := session.DB("simplesurveys").C("survey")
	query := clctn.Find(bson.M{"surveyname": surveyName})
	err := query.One(&response)

	if err != nil {
		return nil
	} else {
		return response
	}
}

func InsertUserResponse(userResponse SurveyResponse) {
	session := MgoSession.Clone()
	defer session.Close()

	clctn := session.DB("simplesurveys").C("survey_response")
	clctn.Insert(userResponse)
}

type DeactivateSurveyStruct struct{
	SurveyId string
	Expired string
}

func DeactivateSurvey(t DeactivateSurveyStruct){
	session := MgoSession.Clone()
	defer session.Close()
	clctn := session.DB("simplesurveys").C("survey")
	clctn.Update(bson.M{"_id": bson.ObjectIdHex(t.SurveyId)}, bson.M{"$set": bson.M{"status":t.Expired}})
}

// ------------------ srevicehandlers survey.go ---------------------------------------


package servicehandlers

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"simplesurveygo/dao"
)

type SurveyHandler struct {
}

func (p SurveyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := methodRouter(p, w, r)
	response.(SrvcRes).RenderResponse(w)

}

func (p SurveyHandler) Get(r *http.Request) SrvcRes {
	params := r.URL.Query()
	ids, present := params["name"]

	fmt.Println(params)

	if present {
		id := ids[0]
		return Response200OK(dao.GetSurveyByName(id))
	} else {
		return Response200OK(dao.GetActiveSurveys())
	}

}

func (p SurveyHandler) Put(r *http.Request) SrvcRes {
	return ResponseNotImplemented()
}

func (p SurveyHandler) Post(r *http.Request) SrvcRes {
	body, err := ioutil.ReadAll(r.Body)

	var userResponse []dao.DeactivateSurveyStruct
	err = json.Unmarshal(body, &userResponse)

	for _,v := range userResponse{
		go dao.DeactivateSurvey(v)
	}
	if err == nil {
		return Simple200OK("Status updated successfully")
	} else {
		return InternalServerError("Something went wrong")
	}
}
