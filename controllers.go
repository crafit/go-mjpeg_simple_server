package main

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"mime/multipart"
	"net/textproto"
)

func jpeg(responseWriter http.ResponseWriter, request *http.Request) {
	log.Printf("Start request %s", request.URL)

	log.Printf("Wait source")
	snapshot := <- source
	log.Printf("Write snapshot")

	responseWriter.Header().Add("Content-Type", "image/jpeg")
	responseWriter.Write(snapshot)
	log.Printf("Success request")
}

func mjpeg(responseWriter http.ResponseWriter, request *http.Request) {
	log.Printf("Start request %s", request.URL)

	mimeWriter := multipart.NewWriter(responseWriter)

	log.Printf("Boundary: %s", mimeWriter.Boundary())

	contentType := fmt.Sprintf("multipart/x-mixed-replace;boundary=%s", mimeWriter.Boundary())
	responseWriter.Header().Add("Content-Type", contentType)

	for {
		frameStartTime := time.Now()
		partHeader := make(textproto.MIMEHeader)
		partHeader.Add("Content-Type", "image/jpeg")

		partWriter, partErr := mimeWriter.CreatePart(partHeader)
		if nil != partErr {
			log.Printf(partErr.Error())
			break
		}

		snapshot := <- source
		if _, writeErr := partWriter.Write(snapshot); nil != writeErr {
			log.Printf(writeErr.Error())
		}
		frameEndTime := time.Now()

		frameDuration := frameEndTime.Sub(frameStartTime)
		fps := float64(time.Second) / float64(frameDuration)
		log.Printf("Frame time: %s (%.2f)", frameDuration, fps)
	}

	log.Printf("Success request")
}
