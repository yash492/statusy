package slack

import "fmt"

type IncidentOpenWorker struct{}

func (i IncidentOpenWorker) Do() error {
	fmt.Println("yolo")
	return nil
}

type IncidentInProgressWorker struct{}

func (i IncidentInProgressWorker) Do() error {
	return nil
}

type IncidentClosedWorker struct{}

func (i IncidentClosedWorker) Do() error {
	return nil
}
