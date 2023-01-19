package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/teris-io/shortid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var ctx = context.TODO()
var baseUrl = "http://localhost:8080/"

type shortenBody struct {
	LongUrl string `json:"long_url"`
}

type SubmissionBody struct {
	Languageid string `json:"language_id"`
	Sourcecode string `json:"source_code"`
	StdInput   string `json:"std_in"`
}
type Request struct {
	LanguageId string `json:"language_id"`
	SourceCode string `json:"source_code"`
	Std_Input  string `json:"std_in"`
}

type UrlDoc struct {
	ID        primitive.ObjectID `bson:"_id"`
	UrlId     string             `bson:"url_id"`
	UrlCode   string             `bson:"urlCode"`
	LongUrl   string             `bson:"longUrl"`
	ShortUrl  string             `bson:"shortUrl"`
	CreatedAt time.Time          `bson:"createdAt"`
}

type ResponseBody struct {
	Message string `json:"message"`
	Error   bool   `json:"error"`
}
type TokenBody struct {
	TokenID primitive.ObjectID `bson:"token_id"`
	Token   string             `bson:"token"`
}
type ReqTokenBody struct {
	Token string `bson:"token"`
}

type SubmissionResponse struct {
	LanguageId int    `json:"language_id"`
	SourceCode string `json:"source_code"`
	StdIn      string `json:"std_in"`
}

func init() {
	clientOptons := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptons)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("project").Collection("urls")
	log.Print("DB connected")
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{"message": "shorten your url"})
	})
	r.GET("/:code", redirect)
	r.POST("/shorten", shorten)
	r.PUT("/:id", updateOneUrl)
	r.DELETE("/:id", deleteOneUrl)
	r.POST("/submission", Submission)
	r.GET("/submissions/", getSubmission)

	r.Run(":8080")
}

func Submission(c *gin.Context) {
	var body SubmissionBody
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	url := "https://judge0-ce.p.rapidapi.com/submissions?base64_encoded=false&fields=*"

	payload := map[string]string{"language_id": body.Languageid, "source_code": body.Sourcecode, "stdin": body.StdInput}

	json_data, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error while marshaling")
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(json_data))

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-RapidAPI-Key", "fe4fbadcaemsh1ce5a676e592147p1970c7jsneaa9bc45bb93")
	req.Header.Add("X-RapidAPI-Host", "judge0-ce.p.rapidapi.com")
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	fmt.Println(res.Body)
	decoder := json.NewDecoder(res.Body)
	var tokenBody ReqTokenBody

	err = decoder.Decode(&tokenBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(tokenBody.Token)

	c.JSON(http.StatusCreated, gin.H{
		"Message":        "Token has been created",
		"Mongo Database": "Token has been created in MongoDB",
		"token":          tokenBody.Token,
	})
	clientOptons := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptons)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("project").Collection("token")
	var docId = primitive.NewObjectID()
	newDoc := &TokenBody{

		TokenID: docId,
		Token:   tokenBody.Token,
	}
	s, err := collection.InsertOne(ctx, newDoc)
	fmt.Println(s)

	log.Print("DB set for token is connected")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}

func getSubmission(c *gin.Context) {
	token := c.Query("token")
	fmt.Println("token =>", token)

	url := fmt.Sprintf("https://judge0-ce.p.rapidapi.com/submissions/%s?base64_encoded=true&fields=*", token)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", "fe4fbadcaemsh1ce5a676e592147p1970c7jsneaa9bc45bb93")
	req.Header.Add("X-RapidAPI-Host", "judge0-ce.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println("body", &body)

	// fmt.Println("Response body", res.Body)
	// decoder := json.NewDecoder(res.Body)

	// var getReq Request
	// fmt.Println(getReq)

	// err := decoder.Decode(&getReq)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	// fmt.Println(getReq)

	var submissionResponse SubmissionResponse
	var response ResponseBody
	fmt.Printf("Type", "%T", string(body))

	err := json.Unmarshal([]byte(body), &submissionResponse)

	if err != nil {
		response.Error = true
		response.Message = "Error in decod"
		fmt.Println("Error while unmarshalling")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	fmt.Println("submission", submissionResponse)
	// submissionResponse.LanguageId = 1
	// submissionResponse.SourceCode = "print(`hello`)"
	// submissionResponse.StdIn = "sumit"
	c.JSON(200, &submissionResponse)

	clientOptons := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptons)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("project").Collection("Response")
	// var docId = primitive.NewObjectID()

	newDoc := &SubmissionResponse{
		LanguageId: submissionResponse.LanguageId,
		SourceCode: submissionResponse.SourceCode,
		StdIn:      submissionResponse.StdIn,
	}
	s, err := collection.InsertOne(ctx, newDoc)
	fmt.Println(s)

	log.Print("DB set for token is connected")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}
func shorten(c *gin.Context) {
	var body shortenBody

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, urlErr := url.ParseRequestURI(body.LongUrl)
	if urlErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": urlErr.Error()})
		return
	}

	shorturlid, idErr := shortid.Generate()
	urlCode := shorturlid[0:4]

	if idErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": idErr.Error()})
		return
	}

	var result bson.M

	queryErr := collection.FindOne(ctx, bson.D{{Key: "longUrl", Value: body.LongUrl}}).Decode(&result)

	if queryErr != nil {
		if queryErr != mongo.ErrNoDocuments {
			c.JSON(http.StatusInternalServerError, gin.H{"error": queryErr.Error()})
			return
		}
	}

	var newUrl = baseUrl + urlCode
	var docId = primitive.NewObjectID()

	newDoc := &UrlDoc{
		ID:        docId,
		UrlCode:   urlCode,
		LongUrl:   body.LongUrl,
		ShortUrl:  newUrl,
		CreatedAt: time.Now(),
		UrlId:     uuid.New().String(),
	}

	_, err := collection.InsertOne(ctx, newDoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"newUrl": newUrl,

		"db_id": docId,
	})
}

func redirect(c *gin.Context) {
	code := c.Param("code")
	var result bson.M
	queryErr := collection.FindOne(ctx, bson.D{{Key: "urlCode", Value: code}}).Decode(&result)

	if queryErr != nil {
		if queryErr == mongo.ErrNoDocuments {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("No URL with code: %s", code)})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": queryErr.Error()})
			return
		}
	}
	log.Print(result["longUrl"])
	var longUrl = fmt.Sprint(result["longUrl"])
	c.Redirect(http.StatusPermanentRedirect, longUrl)
}
func updateOneUrl(ctx *gin.Context) {
	var body shortenBody

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Long Url", body.LongUrl)

	_, urlErr := url.ParseRequestURI(body.LongUrl)
	if urlErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": urlErr.Error()})
		return
	}

	id := ctx.Param("id")
	fmt.Println("ID=>", id)
	filter := bson.D{{Key: "url_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "longUrl", Value: body.LongUrl}}}}
	_, err := collection.UpdateOne(ctx, filter, update)

	if err != nil {
		fmt.Println("Error while updating..", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Long Url has been Updated", "error": false})

}
func deleteOneUrl(ctx *gin.Context) {
	id := ctx.Param("id")
	fmt.Println("ID=>", id)
	filter := bson.D{{Key: "url_id", Value: id}}
	_, err := collection.DeleteOne(ctx, filter)

	if err != nil {
		fmt.Println("Error while deleting..", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Long Url is Deleted", "error": false})

}
