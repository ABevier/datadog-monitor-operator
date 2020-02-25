package main

import (
	"fmt"

	datadogv1alpha1 "github.com/abevier/datadog-monitor-operator/pkg/apis/datadog/v1alpha1"
	"github.com/abevier/datadog-monitor-operator/pkg/dd"
)

func main() {
	client := dd.NewClient()

	monitorType := "metric alert"
	name := "alan test"
	query := "sum(last_1m):sum:argo.router_store_size{environment:aws_abevier} > 50000"

	spec := &datadogv1alpha1.MonitorSpec{
		Type:  &monitorType,
		Name:  &name,
		Query: &query,
	}

	id, err := client.CreateMonitor(spec)
	if err != nil {
		fmt.Printf("Failed to make monitor: %v\n", err)
		return
	}
	fmt.Printf("Created monitor. Id=%v\n", id)

	retMonitor, err := client.GetMonitor(id)
	if err != nil {
		fmt.Printf("Failed to get monitor:%v\n", err)
		return
	}
	fmt.Printf("Fetched Monitor: %v\n", retMonitor)

	newName := "new name - alan test updated"
	retMonitor.Name = &newName
	err = client.UpdateMonitor(id, retMonitor)
	if err != nil {
		fmt.Printf("Failed to update monitor:%v\n", err)
		return
	}

	retMonitor, err = client.GetMonitor(id)
	if err != nil {
		fmt.Printf("Failed to get monitor:%v\n", err)
		return
	}
	fmt.Printf("ReFetched Monitor: %v\n", retMonitor)

	err = client.DeleteMonitor(id)
	if err != nil {
		fmt.Printf("failed to delete monitor: %v\n", err)
		return
	}
	fmt.Printf("Successfully deleted monitor %v\n", id)
}
