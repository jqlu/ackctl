package client

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/denverdino/aliyungo/cs"
	"github.com/jqlu/ackctl/config"
	"sync"
)

var once1, once2 sync.Once
var csClient *cs.Client
var sdkClient *SdkClient

func GetCsClient() *cs.Client {
	once1.Do(func() {
		akConfig := config.MustLoadConfig()
		if akConfig != nil {
			csClient = cs.NewClient(akConfig.AccessKeyId, akConfig.AccessKeySecret)
		}
	})
	return csClient
}

func GetSdkClient() *SdkClient {
	once2.Do(func() {
		akConfig := config.MustLoadConfig()
		client, err := sdk.NewClientWithAccessKey("cn-hangzhou", akConfig.AccessKeyId, akConfig.AccessKeySecret)
		if err != nil {
			panic(err)
		}
		sdkClient = &SdkClient{client: client}
	})

	return sdkClient
}

type SdkClient struct {
	client *sdk.Client
}

func (c *SdkClient) Request(method string, path string, query map[string]string, body *[]byte, result interface{}) error {
	request := requests.NewCommonRequest()
	request.Method = method
	request.Product = "CS"
	request.Domain = "cs.aliyuncs.com"
	request.Version = "2015-12-15"
	request.PathPattern = path
	request.Scheme = "http"
	request.QueryParams = query

	if body != nil {
		request.SetContentType("application/json")
		request.SetContent(*body)
	}

	request.TransToAcsRequest()
	response, err := c.client.ProcessCommonRequest(request)
	if err != nil {
		return err
	}

	if result != nil {
		err := json.Unmarshal(response.GetHttpContentBytes(), result)
		if err != nil {
			return err
		}
	}

	return nil
}