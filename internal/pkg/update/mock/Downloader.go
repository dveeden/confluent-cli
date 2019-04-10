// Code generated by mocker. DO NOT EDIT.
// github.com/travisjeffery/mocker
// Source: private.go

package mock

import (
	io "io"
	sync "sync"

	github_com_aws_aws_sdk_go_service_s3 "github.com/aws/aws-sdk-go/service/s3"
	github_com_aws_aws_sdk_go_service_s3_s3manager "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Downloader is a mock of Downloader interface
type Downloader struct {
	lockDownload sync.Mutex
	DownloadFunc func(w io.WriterAt, input *github_com_aws_aws_sdk_go_service_s3.GetObjectInput, options ...func(*github_com_aws_aws_sdk_go_service_s3_s3manager.Downloader)) (int64, error)

	calls struct {
		Download []struct {
			W       io.WriterAt
			Input   *github_com_aws_aws_sdk_go_service_s3.GetObjectInput
			Options []func(*github_com_aws_aws_sdk_go_service_s3_s3manager.Downloader)
		}
	}
}

// Download mocks base method by wrapping the associated func.
func (m *Downloader) Download(w io.WriterAt, input *github_com_aws_aws_sdk_go_service_s3.GetObjectInput, options ...func(*github_com_aws_aws_sdk_go_service_s3_s3manager.Downloader)) (int64, error) {
	m.lockDownload.Lock()
	defer m.lockDownload.Unlock()

	if m.DownloadFunc == nil {
		panic("mocker: Downloader.DownloadFunc is nil but Downloader.Download was called.")
	}

	call := struct {
		W       io.WriterAt
		Input   *github_com_aws_aws_sdk_go_service_s3.GetObjectInput
		Options []func(*github_com_aws_aws_sdk_go_service_s3_s3manager.Downloader)
	}{
		W:       w,
		Input:   input,
		Options: options,
	}

	m.calls.Download = append(m.calls.Download, call)

	return m.DownloadFunc(w, input, options...)
}

// DownloadCalled returns true if Download was called at least once.
func (m *Downloader) DownloadCalled() bool {
	m.lockDownload.Lock()
	defer m.lockDownload.Unlock()

	return len(m.calls.Download) > 0
}

// DownloadCalls returns the calls made to Download.
func (m *Downloader) DownloadCalls() []struct {
	W       io.WriterAt
	Input   *github_com_aws_aws_sdk_go_service_s3.GetObjectInput
	Options []func(*github_com_aws_aws_sdk_go_service_s3_s3manager.Downloader)
} {
	m.lockDownload.Lock()
	defer m.lockDownload.Unlock()

	return m.calls.Download
}

// Reset resets the calls made to the mocked methods.
func (m *Downloader) Reset() {
	m.lockDownload.Lock()
	m.calls.Download = nil
	m.lockDownload.Unlock()
}
