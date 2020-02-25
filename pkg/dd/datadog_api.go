package dd

import (
	"os"

	datadogv1alpha1 "github.com/abevier/datadog-monitor-operator/pkg/apis/datadog/v1alpha1"
	"github.com/zorkian/go-datadog-api"
)

type DataDogClient struct {
	client *datadog.Client
}

func NewClient() *DataDogClient {
	return &DataDogClient{
		client: datadog.NewClient(os.Getenv("DD_API_KEY"), os.Getenv("DD_APP_KEY")),
	}
}

func (c *DataDogClient) CreateMonitor(monitorSpec *datadogv1alpha1.MonitorSpec) (int, error) {
	monitor := monitorSpec.ToDDMonitor()

	result, err := c.client.CreateMonitor(monitor)
	if err != nil {
		return 0, err
	}

	return result.GetId(), nil
}

func (c *DataDogClient) GetMonitor(id int) (*datadogv1alpha1.MonitorSpec, error) {
	monitor, err := c.client.GetMonitor(id)
	if err != nil {
		//TODO: wrap error?
		return nil, err
	}

	//TODO: map function?
	return &datadogv1alpha1.MonitorSpec{
		Type:    monitor.Type,
		Query:   monitor.Query,
		Name:    monitor.Name,
		Message: monitor.Message,
	}, nil
}

func (c *DataDogClient) UpdateMonitor(id int, monitorSpec *datadogv1alpha1.MonitorSpec) error {
	monitor := monitorSpec.ToDDMonitor()
	monitor.SetId(id)
	return c.client.UpdateMonitor(monitor)
}

func (c *DataDogClient) DeleteMonitor(id int) error {
	return c.client.DeleteMonitor(id)
}
