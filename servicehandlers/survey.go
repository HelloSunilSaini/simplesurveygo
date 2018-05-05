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
