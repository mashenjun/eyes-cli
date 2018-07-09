package main

import (
	"testing"
	"qiniupkg.com/x/log.v7"
)

func TestBatchUpload(t *testing.T) {
	endpoint := "localhost:8080"
	token, err := login("qiniu","123123",endpoint)
	log.Println(err)
	log.Println(token)
}

func TestCreatFacelib(t *testing.T) {
	endpoint := "localhost:8080"
	token, err := login("qiniu","123123",endpoint)
	log.Println(err)
	metaData := CreateFaceMetaData{
		Name: "scriptTest",
		Count: 10,
		Description:"scriptTest",
		Source:"qiniuStaff",
	}
	facelibID, groupID, _, err := newFacelib(token,endpoint,metaData)
	log.Println(err)
	log.Println(facelibID)
	log.Println(groupID)
}

func TestAddFace(t *testing.T) {
	endpoint := "localhost:8080"
	token, err := login("qiniu","123123",endpoint)
	log.Println(err)
	metaData := AddFaceMetaData{
		Name: "小马",
		IdCardNum:"000000000000000000",
		Gender: 1,
		Images: []string{"face.jpg"},

	}
	err = addFace(token, "9cf54c04bed84bf3806f2b4eef4554aa",endpoint,139, metaData)
	log.Println(err)
}
