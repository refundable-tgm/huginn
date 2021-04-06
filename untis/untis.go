package untis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const ClientName = "Refundable"
const URL = "https://neilo.webuntis.com/WebUntis/jsonrpc.do?school=tgm"

var activeClients map[string]Client

type Client struct {
	Username string
	Password string
	SessionID string
	PersonType int
	PersonID int
	Closed bool
	Authenticated bool
}

type Lesson struct {
	Start time.Time
	End time.Time
	ClassIDs []int
	Classes []string
	TeacherIDs []int
	Teachers []string
	RoomIDs []int
	Rooms []string
}

func CreateClient(username, password string) *Client {
	client := Client{
		Username:   username,
		Password:   password,
		SessionID:  "",
		PersonType: -1,
		PersonID:   -1,
		Closed:     false,
		Authenticated: false,
	}
	activeClients[username] = client
	return &client
}

func (client *Client) Authenticate() error {
	id := getID()
	body, _ := json.Marshal(map[string]interface{}{
		"id": id,
		"method": "authenticate",
		"params": map[string]string {
			"user": client.Username,
			"password": client.Password,
			"client": ClientName,
		},
		"jsonrpc": "2.0",
	})
	resp, err := http.Post(URL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	r := struct{
		JSONRPC string `json:"jsonrpc"`
		ID string `json:"id"`
		Result map[string]string `json:"result"`
	}{}
	err = json.Unmarshal(respBody, &r)
	if err != nil {
		return err
	}
	personType, _ := strconv.Atoi(r.Result["personType"])
	personID, _ := strconv.Atoi(r.Result["personId"])
	if r.ID == strconv.Itoa(id) {
		client.SessionID = r.Result["sessionID"]
		client.PersonType = personType
		client.PersonID = personID
		client.Authenticated = true
		return nil
	} else {
		return fmt.Errorf("IDs not matching")
	}
}

func (client Client) GetTimetable(start, end time.Time) ([]Lesson, error) {
	if !client.Authenticated {
		return nil, fmt.Errorf("not authenticated")
	}
	id := getID()
	smonth := strconv.Itoa(int(start.Month()))
	if len(smonth) == 1 {
		smonth = "0" + smonth
	}
	sday := strconv.Itoa(start.Day())
	if len(sday) == 1 {
		sday = "0" + sday
	}
	emonth := strconv.Itoa(int(end.Month()))
	if len(emonth) == 1 {
		emonth = "0" + emonth
	}
	eday := strconv.Itoa(end.Day())
	if len(eday) == 1 {
		eday = "0" + eday
	}
	body, _ := json.Marshal(map[string]interface{}{
		"id": id,
		"method": "getTimetable",
		"params": map[string]interface{}{
			"id": client.PersonID,
			"type": client.PersonType,
			"startDate": strconv.Itoa(start.Year()) + smonth + sday,
			"endDate": strconv.Itoa(end.Year()) + emonth + eday,
		},
		"jsonrpc": "2.0",
	})
	req, err :=  http.NewRequest("POST", URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("JSESSIONID", client.SessionID)
	repClient := &http.Client{}
	resp, err := repClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r := struct{
		JSONRPC string `json:"jsonrpc"`
		ID string `json:"id"`
		Result []struct{
			ID int `json:"id"`
			Date int `json:"date"`
			StartTime int `json:"startTime"`
			EndTime int `json:"endTime"`
			Kl []struct{
				ID int `json:"id"`
			} `json:"kl"`
			Te []struct{
				ID int `json:"id"`
			} `json:"te"`
			Su []struct{
				ID int `json:"id"`
			} `json:"su"`
			Ro []struct{
				ID int `json:"id"`
			} `json:"ro"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(respBody, &r)
	if err != nil {
		return nil, err
	}
	rid, _ := strconv.Atoi(r.ID)
	if rid == id {
		lessons := make([]Lesson, 0)
		for _, l := range r.Result {
			date := strconv.Itoa(l.Date)
			year, _ := strconv.Atoi(date[0:4])
			month, _ := strconv.Atoi(date[4:6])
			day, _ := strconv.Atoi(date[6:8])
			startTime := strconv.Itoa(l.StartTime)
			startHour, _ := strconv.Atoi(startTime[0:2])
			startMinute, _ := strconv.Atoi(startTime[2:4])
			endTime := strconv.Itoa(l.EndTime)
			endHour, _ := strconv.Atoi(endTime[0:2])
			endMinute, _ := strconv.Atoi(endTime[2:4])
			classIDArr := make([]int, 0)
			for _, kls := range l.Kl {
				classIDArr = append(classIDArr, kls.ID)
			}
			classArr, err := client.ResolveClass(classIDArr)
			if err != nil {
				return nil, err
			}
			teachIDArr := make([]int, 0)
			for _, tes := range l.Te {
				teachIDArr = append(teachIDArr, tes.ID)
			}
			teachArr, err := client.ResolveTeacher(teachIDArr)
			if err != nil {
				return nil, err
			}
			roomIDArr := make([]int, 0)
			for _, ros := range l.Ro {
				roomIDArr = append(roomIDArr, ros.ID)
			}
			roomArr, err := client.ResolveRoom(roomIDArr)
			if err != nil {
				return nil, err
			}
			lessons = append(lessons, Lesson{
				Start:      time.Date(year, time.Month(month), day, startHour, startMinute, 0, 0, time.UTC),
				End:        time.Date(year, time.Month(month), day, endHour, endMinute, 0, 0, time.UTC),
				ClassIDs:   classIDArr,
				Classes:    classArr,
				TeacherIDs: teachIDArr,
				Teachers:   teachArr,
				RoomIDs:    roomIDArr,
				Rooms:      roomArr,
			})
		}
		return lessons, nil
	}
	return nil, fmt.Errorf("ids not matching")
}

func (client Client) ResolveTeacher(ids []int) ([]string, error) {
	if !client.Authenticated {
		return nil, fmt.Errorf("not authenticated")
	}
	id := getID()
	body, _ := json.Marshal(map[string]interface{}{
		"id": id,
		"method": "getTeachers",
		"params": map[string]string{},
		"jsonrpc": "2.0",
	})
	req, err :=  http.NewRequest("POST", URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("JSESSIONID", client.SessionID)
	repClient := &http.Client{}
	resp, err := repClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r := struct{
		JSONRPC string `json:"jsonrpc"`
		ID string `json:"id"`
		Result []struct{
			ID int `json:"id"`
			Name string `json:"name"`
			Forename string `json:"foreName"`
			Longname string `json:"longName"`
			ForeColor string `json:"foreColor"`
			BackColor string `json:"backColor"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(respBody, &r)
	if err != nil {
		return nil, err
	}
	teacher := make([]string, 0)
	for _, id := range ids {
		for _, res := range r.Result {
			if id == res.ID {
				teacher = append(teacher, res.Forename + res.Name)
			}
		}
	}
	return teacher, nil
}

func (client Client) ResolveRoom (ids []int) ([]string, error) {
	if !client.Authenticated {
		return nil, fmt.Errorf("not authenticated")
	}
	id := getID()
	body, _ := json.Marshal(map[string]interface{}{
		"id": id,
		"method": "getRooms",
		"params": map[string]string{},
		"jsonrpc": "2.0",
	})
	req, err :=  http.NewRequest("POST", URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("JSESSIONID", client.SessionID)
	repClient := &http.Client{}
	resp, err := repClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r := struct{
		JSONRPC string `json:"jsonrpc"`
		ID string `json:"id"`
		Result []struct{
			ID int `json:"id"`
			Name string `json:"name"`
			Longname string `json:"longName"`
			ForeColor string `json:"foreColor"`
			BackColor string `json:"backColor"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(respBody, &r)
	if err != nil {
		return nil, err
	}
	rid, _ := strconv.Atoi(r.ID)
	if id == rid {
		rooms := make([]string, 0)
		for _, id := range ids {
			for _, res := range r.Result {
				if id == res.ID {
					rooms = append(rooms, res.Name)
				}
			}
		}
		return rooms, nil
	} else {
		return nil, fmt.Errorf("ids not matching")
	}
}

func (client Client) ResolveClass(ids []int) ([]string, error) {
	if !client.Authenticated {
		return nil, fmt.Errorf("not authenticated")
	}
	id := getID()
	body, _ := json.Marshal(map[string]interface{}{
		"id": id,
		"method": "getKlassen",
		"params": map[string]string{},
		"jsonrpc": "2.0",
	})
	req, err :=  http.NewRequest("POST", URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("JSESSIONID", client.SessionID)
	repClient := &http.Client{}
	resp, err := repClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r := struct{
		JSONRPC string `json:"jsonrpc"`
		ID string `json:"id"`
		Result []struct{
			ID int `json:"id"`
			Name string `json:"name"`
			Longname string `json:"longName"`
			ForeColor string `json:"foreColor"`
			BackColor string `json:"backColor"`
			Teacher1 int `json:"teacher1"`
			Teacher2 int `json:"teacher2"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(respBody, &r)
	if err != nil {
		return nil, err
	}
	rid, _ := strconv.Atoi(r.ID)
	if rid == id {
		classes := make([]string, 0)
		for _, id := range ids {
			for _, res := range r.Result {
				if id == res.ID {
					classes = append(classes, res.Name)
				}
			}
		}
		return classes, nil
	} else {
		return nil, fmt.Errorf("ids not matching")
	}
}

func (client *Client) Close() error {
	if !client.Authenticated {
		return fmt.Errorf("not authenticated")
	}
	delete(activeClients, client.Username)
	client.Closed = true
	id := getID()
	body, _ := json.Marshal(map[string]interface{}{
		"id": id,
		"method": "logout",
		"params": map[string]string {},
		"jsonrpc": "2.0",
	})
	req, err :=  http.NewRequest("POST", URL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("JSESSIONID", client.SessionID)
	repClient := &http.Client{}
	_, err = repClient.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func getID() int {
	return rand.Intn(math.MaxInt64)
}

func GetLessonNrByStart(start time.Time) int {
	switch start.Hour() {
	case 8:
		if start.Minute() == 0 {
			return 1
		} else if start.Minute() == 50{
			return 2
		}
	case 9:
		return 3
	case 10:
		return 4
	case 11:
		return 5
	case 12:
		return 6
	case 13:
		return 7
	case 14:
		return 8
	case 15:
		return 9
	case 16:
		return 10
	case 17:
		if start.Minute() == 0 {
			return 11
		} else if start.Minute() == 45 {
			return 12
		}
	case 18:
		return 13
	case 19:
		return 14
	case 20:
		return 15
	}
	return -1
}

func GetLessonNrByEnd(end time.Time) int {
	switch end.Hour() {
	case 8:
		return 1
	case 9:
		return 2
	case 10:
		return 3
	case 11:
		return 4
	case 12:
		return 5
	case 13:
		return 6
	case 14:
		return 7
	case 15:
		return 8
	case 16:
		if end.Minute() == 0 {
			return 9
		} else if end.Minute() == 50 {
			return 10
		}
	case 17:
		return 11
	case 18:
		return 12
	case 19:
		return 13
	case 20:
		return 14
	case 21:
		return 15
	}
	return -1
}


