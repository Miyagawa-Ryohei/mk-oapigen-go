// Package entity
// this file is generated by swag2go. Implement all methods
//
//	authorized : miyagawa.ryohei@bell-face.com
package entity

import (
	"fmt"
	"mk-oapigen-go/gen"
)

// JobService
// /job Handler
type JobService struct {
	// TODO edit here
}

func (s *JobService) PostTranscripitonJobStatus(req gen.JobStateSchema) (*gen.JobStateSchema, error) {
	return nil, fmt.Errorf("not implemented")
}

// JobmeetingIdService
// /job/{meeting_id} Handler
type JobmeetingIdService struct {
	// TODO edit here
}

func (s *JobmeetingIdService) GetTranscripitonJobStatus(req gen.GetTranscripitonJobStatusParam) (*gen.JobStateSchema, error) {
	return nil, fmt.Errorf("not implemented")
}

// JobsService
// /jobs Handler
type JobsService struct {
	// TODO edit here
}

func (s *JobsService) GetTranscripitonJobStatusList(req gen.GetTranscripitonJobStatusListParam) (*gen.JobStateSchemaList, error) {
	return nil, fmt.Errorf("not implemented")
}

// TranscriptionService
// /transcription Handler
type TranscriptionService struct {
	// TODO edit here
}

func (s *TranscriptionService) PostTranscriptionRequest(req gen.TranscriptionRequestInput) (*gen.TranscriptionRequestResult, error) {
	return nil, fmt.Errorf("not implemented")
}

// RootService
// / Handler
type RootService struct {
	// TODO edit here
}

func (s *RootService) Root() (*gen.RootResult, error) {
	return nil, fmt.Errorf("not implemented")
}

// ChkService
// /_chk Handler
type ChkService struct {
	// TODO edit here
}

func (s *ChkService) HealthCheck() (*gen.HealthCheckResult, error) {
	return nil, fmt.Errorf("not implemented")
}

func GetJobServiceProvider(
// TODO edit here
) gen.JobServiceProvider {
	return func() (gen.JobService, error) {
		//TODO edit here
		return &JobService{
			//TODO edit here
		}, nil
	}
}

func GetJobmeetingIdServiceProvider(
// TODO edit here
) gen.JobmeetingIdServiceProvider {
	return func() (gen.JobmeetingIdService, error) {
		//TODO edit here
		return &JobmeetingIdService{
			//TODO edit here
		}, nil
	}
}

func GetJobsServiceProvider(
// TODO edit here
) gen.JobsServiceProvider {
	return func() (gen.JobsService, error) {
		//TODO edit here
		return &JobsService{
			//TODO edit here
		}, nil
	}
}

func GetTranscriptionServiceProvider(
// TODO edit here
) gen.TranscriptionServiceProvider {
	return func() (gen.TranscriptionService, error) {
		//TODO edit here
		return &TranscriptionService{
			//TODO edit here
		}, nil
	}
}

func GetRootServiceProvider(
// TODO edit here
) gen.RootServiceProvider {
	return func() (gen.RootService, error) {
		//TODO edit here
		return &RootService{
			//TODO edit here
		}, nil
	}
}

func GetChkServiceProvider(
// TODO edit here
) gen.ChkServiceProvider {
	return func() (gen.ChkService, error) {
		//TODO edit here
		return &ChkService{
			//TODO edit here
		}, nil
	}
}
