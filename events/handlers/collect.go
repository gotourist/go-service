package handlers

import (
	"encoding/json"
	"github.com/iman_task/go-service/domain/entities"
	brokerpb "github.com/iman_task/go-service/genproto/broker/post"
	loggerPkg "github.com/iman_task/go-service/pkg/logger"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

func (e *EventHandler) Collect(wg *sync.WaitGroup, page int) {
	defer wg.Done()

	request, err := http.NewRequest(http.MethodGet, "https://gorest.co.in/public/v1/posts", nil)
	if err != nil {
		e.logger.Error("failed to create new request", loggerPkg.Error(err))
		return
	}

	values := request.URL.Query()
	values.Add("page", strconv.Itoa(page))

	request.URL.RawQuery = values.Encode()

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		e.logger.Error("failed to request", loggerPkg.Error(err))
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			e.logger.Error("failed to close response body", loggerPkg.Error(err))
		}
	}(response.Body)

	body, _ := ioutil.ReadAll(response.Body)

	var postList entities.PostList

	err = json.Unmarshal(body, &postList)
	if err != nil {
		e.logger.Error("failed to unmarshal data", loggerPkg.Error(err))
		return
	}

	for _, post := range postList.Data {
		e.CreatePost(&post)
	}

	return
}

func (e *EventHandler) CreatePost(data *entities.Post) {

	post, err := e.storage.Post().CreatePost(data)
	if err != nil {
		e.logger.Error("failed to create post in db", loggerPkg.Error(err))
		return
	}

	err = e.publishCreatePostMessage(post)

	return
}

func (e *EventHandler) publishCreatePostMessage(post *entities.Post) error {

	var postPb brokerpb.Post

	postPb.Id = post.Id
	postPb.Title = post.Title
	postPb.Body = post.Body
	postPb.IsDeleted = post.IsDeleted
	postPb.CreatedAt = post.CreatedAt

	data, err := postPb.Marshal()
	if err != nil {
		return err
	}

	logBody := postPb.String()

	err = e.publisher[PostAddTopic].Publish([]byte("post_add"), data, logBody)
	if err != nil {
		return err
	}

	return nil
}
