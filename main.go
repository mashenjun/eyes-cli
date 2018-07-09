package main

import (
	"fmt"
	"os"

	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"qiniupkg.com/x/log.v7"
	"strconv"
	"strings"
	"time"
	"gopkg.in/cheggaaa/pb.v1"

)

const (
	LoginUrl         = "/api/v1/login"
	CreateFacelibUrl = "/api/v1/internal/new_facelib"
	AddFaceImgUrl    = "/api/v1/internal/add_face"
	EndAddingUrl     = "/api/v1/internal/end_adding"
)

const (
	TrackingPerfix = "tracking_"
	ErrorLogPerfix = "error_"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "batch, b",
			Usage: "upload files in batch mode otherwise upload single file",
		},
		cli.StringFlag{
			Name:  "metadata, m",
			Usage: "define path to the metadata json file, used under batch mode",
		},
		cli.StringFlag{
			Name:  "file, f",
			Usage: "name of the target file, used under single mode",
		},
		cli.StringFlag{
			Name:  "group, g",
			Usage: "name of the target Argues group, used under single mode",
		},
		cli.StringFlag{
			Name:  "username, u",
			Usage: "username",
		},
		cli.StringFlag{
			Name:  "password, p",
			Usage: "password",
		},
		cli.StringFlag{
			Name:  "endpoint, e",
			Usage: "endpoint url",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "add",
			Usage: "Add single face image to eyes-argus",
			Action: func(c *cli.Context) error {
				userName := c.String("username")
				if userName == "" {
					fmt.Println("username is not defined")
					return nil
				}
				password := c.String("password")
				if password == "" {
					fmt.Println("password is not defined")
					return nil
				}
				endpoint := c.String("endpoint")
				if endpoint == "" {
					fmt.Println("endpoint is not defined")
					return nil
				}
				targetFileName := c.String("file")
				if targetFileName == "" {
					fmt.Println("target file (-f/--file) is not defined under single mode")
					return nil
				}
				targeGroup := c.String("group")
				if targeGroup == "" {
					fmt.Println("target group (-g/--group) is not defined under single mode")
					return nil
				}

				// TODO processing single upload
				return nil
			},
		},
		{
			Name:  "upload",
			Usage: "Add face images to eyes-argus",
			Action: func(c *cli.Context) error {
				// processing
				userName := c.String("username")
				if userName == "" {
					fmt.Println("username is not defined")
					return nil
				}
				password := c.String("password")
				if password == "" {
					fmt.Println("password is not defined")
					return nil
				}
				endpoint := c.String("endpoint")
				if endpoint == "" {
					fmt.Println("endpoint is not defined")
					return nil
				}

				metaFileName := c.String("metadata")
				if metaFileName == "" {
					fmt.Println("metadata.json file (-m/--metadata) is not defined under batch mode")
					return nil
				}
				// processing batch upload
				count, total, errLogFileName, err := batchUpload(userName, password, metaFileName, endpoint)
				if err!= nil {
					if total == 0 {
						fmt.Println("")
					}
					fmt.Printf("ABORT| %v/%v uplaoded, err:%v\n", count, total, err)
				}else {
					if count== total {
						fmt.Printf("SUCCESS| %v/%v uplaoded\n", count, total)
					}else{
						fmt.Printf("FAIL| %v/%v uplaoded. please view the file: %v\n", count, total, errLogFileName)
					}
				}
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func singleUpload() (err error) {
	// TODO
	return nil
}

func batchUpload(username, password, metaFileName, endpoint string) (count, total int, errorFileName string, err error) {
	// TODO
	var groupID string
	var facelibID int
	var trackingFileName string
	var token string
	var cMetaInfo CreateFaceMetaData
	isNew := true
	already := make([]string,0)
	fmt.Printf("connecting to %v", endpoint)
	token, err = login(username, password, endpoint)
	if err != nil {
		return 0, 0, errorFileName, err
	}
	fmt.Println(" --> login success!")
	// check traking file fist
	fmt.Print("check files .")
	file, err := os.Open(metaFileName)
	if err != nil {
		return 0, 0, errorFileName, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	// read first line
	for scanner.Scan() {
		content := scanner.Bytes()
		err = json.Unmarshal(content, &cMetaInfo)
		if err != nil {
			return 0, 0, errorFileName, err
		}
		fmt.Print(".")
		total = cMetaInfo.Count
		trackingFileName = fmt.Sprintf("%s%v_%v", TrackingPerfix, username, cMetaInfo.Name)
		errorFileName = fmt.Sprintf("%s%v_%v_%v", ErrorLogPerfix, username, cMetaInfo.Name,time.Now().Format("20060102150405"))
		if FileExists(trackingFileName) {
			isNew = false
		} else {
			TouchFile(trackingFileName)
		}
		trackingFile:= NewTrackingLog(trackingFileName)
		defer trackingFile.Close()
		fmt.Print(".")
		if FileExists(errorFileName) {
		} else {
			TouchFile(errorFileName)
		}
		errLogFile := NewErrorLog(errorFileName)
		defer errLogFile.Close()
		fmt.Print(".")

		facelibID, groupID, already, err = newFacelib(token, endpoint, cMetaInfo)
		if err != nil {
			errMsg := fmt.Sprintf(`{"line":%v, "facelib_name":"%v", "err_msg":"%v"}`, 1, cMetaInfo.Name, err)
			ErrorLog.Println(errMsg)
			fmt.Println(errMsg)
			return count, total, errorFileName, err
		}
		totalLine, err := TotalLine(metaFileName)
		if err != nil {
			return count, total, errorFileName, errors.New("need wc to get metadata line number")
		}
		if totalLine != cMetaInfo.Count+1{
			return count, total, errorFileName, errors.New("metadata size value not match images count")
		}
		if isNew {
			trackingMsg:= fmt.Sprintf(`{"facelib_id":%v, "facelib_name":"%v", "group_id":"%v"}`, facelibID, cMetaInfo.Name, groupID)
			TrackingLog.Println(trackingMsg)
		}
		fmt.Print("done\n")
		break
	}

	bar := pb.New(total)
	bar.SetMaxWidth(80)
	bar.ShowTimeLeft = false
	if isNew {
		fmt.Printf("start uploading to %v\n",cMetaInfo.Name)
		lineCount := 1
		bar.Start()
		for scanner.Scan() {
			lineCount++
			bar.Increment()
			content := scanner.Bytes()
			var metaInfo AddFaceMetaData
			err = json.Unmarshal(content, &metaInfo)
			if err != nil {
				log.Error(err)
				errMsg := fmt.Sprintf(`{"line":%v, content":"%v", "err_msg":"%v"}`,lineCount, string(content), err)
				ErrorLog.Println(errMsg)
				continue
			}
			err = addFace(token, groupID, endpoint, facelibID, metaInfo)
			if err != nil {
				errMsg := fmt.Sprintf(`{"line":%v, "human_name":"%v", "file":"%v", "err_msg":"%v"}`, lineCount, metaInfo.Name, metaInfo.Images[0], err)
				ErrorLog.Println(errMsg)
			} else {
				trackingMsg := fmt.Sprintf(`{"line":%v, "human_name":"%v", "file":"%v"}`,lineCount, metaInfo.Name, metaInfo.Images[0])
				TrackingLog.Println(trackingMsg)
				count++
			}
		}
	}else {
		fmt.Printf("continue uploading to %v\n",cMetaInfo.Name)
		count = len(already)
		tmp := map[string]struct{}{}
		for _, item :=range already {
			if _, ok := tmp[item]; !ok {
				tmp[item]= struct{}{}
			}
		}
		lineCount := 1
		bar.Start()
		for scanner.Scan() {
			lineCount++
			bar.Increment()
			content := scanner.Bytes()
			var metaInfo AddFaceMetaData
			err = json.Unmarshal(content, &metaInfo)
			if err != nil {
				errMsg := fmt.Sprintf(`{"line":%v, content":"%v", "err_msg":"%v"}`,lineCount, string(content), err)
				ErrorLog.Println(errMsg)
			}
			if _, ok := tmp[metaInfo.Name];!ok{
				err = addFace(token, groupID, endpoint, facelibID, metaInfo)
				if err != nil {
					errMsg := fmt.Sprintf(`{"line":%v, "human_name":"%v", "file":"%v", "err_msg":"%v"}`, lineCount, metaInfo.Name, metaInfo.Images[0], err)
					ErrorLog.Println(errMsg)
				} else {
					tackingMsg := fmt.Sprintf(`{"line":%v, "human_name":"%v", "file":"%v"}`,lineCount, metaInfo.Name, metaInfo.Images[0])
					TrackingLog.Println(tackingMsg)
					count++
				}
			}
		}
	}
	bar.FinishPrint("Finish!")
	alreadyAdded,err := endAdding(token, endpoint, facelibID, count)
	if err != nil {
		return count, total, errorFileName, err
	}
	count = alreadyAdded
	if alreadyAdded == cMetaInfo.Count {
		err = MoveFile(trackingFileName,fmt.Sprintf(".%v_%v",trackingFileName, time.Now().Format("20060102150405")))
		if err != nil {
			return count, total, errorFileName, err
		}
	}
	return count, total, errorFileName,nil
}

func newFacelib(token, endpoint string, metaInfo CreateFaceMetaData) (facelibID int, groupID string, already []string, err error) {
	if !strings.Contains(endpoint, "http://") {
		endpoint = fmt.Sprintf("http://%s", endpoint)
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		log.Error(err)
		return
	}
	u.Path = path.Join(u.Path, CreateFacelibUrl)
	q := u.Query()
	q.Set("token", token)
	u.RawQuery = q.Encode()
	newFacelibUrl := u.String()
	body := CreateFacelibPostBody{
		Name:        metaInfo.Name,
		Count:       metaInfo.Count,
		Description: metaInfo.Description,
		Source:      metaInfo.Source,
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		log.Error(err)
		return
	}
	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	request, err := http.NewRequest("POST", newFacelibUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Fatal(err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	//resp, err := http.Post(newFacelibUrl, "application/json", bytes.NewBuffer(bodyJson))
	resp, err := client.Do(request)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}
	// err processing
	if resp.StatusCode == http.StatusNotFound {
		err = errors.New(string(bodyBytes[:]))
		return
	}
	if resp.StatusCode != http.StatusOK {
		var errContent ErrBody
		err = json.Unmarshal(bodyBytes, &errContent)
		if err != nil {
			log.Error(err)
			return
		}
		err = errors.New(errContent.Cause)
		return
	}
	// return token
	var content CreateFacelibResponse
	err = json.Unmarshal(bodyBytes, &content)
	if err != nil {
		log.Error(err)
		return
	}
	return content.Data.ID, content.Data.GroupID, content.Data.HumanInfoNames, err
}

func addFace(token, groupID, endpoint string, facelibID int, metaInfo AddFaceMetaData) (err error) {
	// Prepare url
	if !strings.Contains(endpoint, "http://") {
		endpoint = fmt.Sprintf("http://%s", endpoint)
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return
	}
	u.Path = path.Join(u.Path, AddFaceImgUrl)
	q := u.Query()
	q.Set("token", token)
	u.RawQuery = q.Encode()
	addFaceUrl := u.String()
	// Prepare params
	if metaInfo.IdCardNum == "" {
		randIdCardNum, _ := NewUuid()
		metaInfo.IdCardNum = "FAKE:"+randIdCardNum
	}
	param := map[string]string{
		"facelib_id":  strconv.Itoa(facelibID),
		"group_id":    groupID,
		"name":        metaInfo.Name,
		"gender":      strconv.Itoa(metaInfo.Gender),
		"id_card_num": metaInfo.IdCardNum,
	}
	// Prepare request
	req, err := newfileUploadRequest(addFaceUrl, param, "file", metaInfo.Images[0])
	if err != nil {
		return err
	}
	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Process response
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// err processing
	if resp.StatusCode != http.StatusOK {
		var errContent ErrBody
		err = json.Unmarshal(bodyBytes, &errContent)
		if err != nil {
			return err
		}
		err = errors.New(errContent.Cause)
		return err
	}
	var content UploadResponse
	err = json.Unmarshal(bodyBytes, &content)
	if err != nil {
		return err
	}

	return nil
}

func endAdding(token, endpoint string, facelibID, count int) (already int, err error) {
	if !strings.Contains(endpoint, "http://") {
		endpoint = fmt.Sprintf("http://%s", endpoint)
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		log.Error(err)
		return
	}
	u.Path = path.Join(u.Path, EndAddingUrl)
	q := u.Query()
	q.Set("token", token)
	u.RawQuery = q.Encode()
	newFacelibUrl := u.String()

	body := EndAddingPostBody{
		FacelibID: facelibID,
		Count: count,
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		log.Error(err)
		return
	}
	resp, err := http.Post(newFacelibUrl, "application/json", bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}
	// err processing
	if resp.StatusCode == http.StatusNotFound {
		err = errors.New(string(bodyBytes[:]))
		return 0, err
	}
	if resp.StatusCode != http.StatusOK {
		var errContent ErrBody
		err = json.Unmarshal(bodyBytes, &errContent)
		if err != nil {
			return 0, err
		}
		err = errors.New(errContent.Cause)
		return 0, err
	}

	var content EndAddingResponse
	err = json.Unmarshal(bodyBytes, &content)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	return content.Data.Already,nil
}

func login(username, password, endpoint string) (token string, err error) {
	if !strings.Contains(endpoint, "http://") {
		endpoint = fmt.Sprintf("http://%s", endpoint)
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return
	}
	u.Path = path.Join(u.Path, LoginUrl)
	loginUrl := u.String()
	body := LoginPostBody{
		Username: username,
		Password: password,
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return
	}
	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	request, err := http.NewRequest("POST", loginUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Fatal(err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	//resp, err := http.Post(loginUrl, "application/json", bytes.NewBuffer(bodyJson))
	resp, err := client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	// err processing
	if resp.StatusCode != http.StatusOK {
		var errContent ErrBody
		err = json.Unmarshal(bodyBytes, &errContent)
		if err != nil {
			return
		}
		err = errors.New("login fail --> "+ errContent.Cause)
		return
	}
	// return token
	var content LoginResponse
	err = json.Unmarshal(bodyBytes, &content)
	if err != nil {
		return
	}
	return content.Data.Token, nil
}
