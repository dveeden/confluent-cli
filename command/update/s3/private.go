package s3

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-version"

	"github.com/confluentinc/cli/log"
)

type PrivateRepoParams struct {
	S3BinBucket string
	S3BinRegion string
	S3BinPrefix string
	AWSProfiles []string
	Logger      *log.Logger
}

type PrivateRepo struct {
	*PrivateRepoParams
	session *session.Session
	s3svc   *s3.S3
}

func NewPrivateRepo(params *PrivateRepoParams) (*PrivateRepo, error) {
	if err := validate(params); err != nil {
		return nil, err
	}

	creds, err := GetCredentials(params.AWSProfiles)
	if err != nil {
		return nil, err
	}

	s, err := session.NewSession(&aws.Config{
		Region:      aws.String(params.S3BinRegion),
		Credentials: creds,
	})
	if err != nil {
		return nil, err
	}
	return &PrivateRepo{
		PrivateRepoParams: params,
		session: s,
		s3svc:   s3.New(s),
	}, nil
}

func validate(params *PrivateRepoParams) error {
	var err *multierror.Error
	if params.S3BinRegion == "" {
		err = multierror.Append(err, fmt.Errorf("missing required parameter: S3BinRegion"))
	}
	if params.S3BinBucket == "" {
		err = multierror.Append(err, fmt.Errorf("missing required parameter: S3BinBucket"))
	}
	if params.S3BinPrefix == "" {
		err = multierror.Append(err, fmt.Errorf("missing required parameter: S3BinPrefix"))
	}
	return err.ErrorOrNil()
}

func (r *PrivateRepo) GetAvailableVersions(name string) (version.Collection, error) {
	result, err := r.s3svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(r.S3BinBucket),
		Prefix: aws.String(r.S3BinPrefix + "/"),
	})
	if err != nil {
		return nil, fmt.Errorf("error listing s3 bucket %s", err)
	}

	var availableVersions version.Collection
	for _, c := range result.Contents {
		// Format: S3BinPrefix/NAME-v0.0.0-OS-ARCH
		split := strings.Split(*c.Key, "-")

		// Skip files that don't match our naming standards for binaries
		if len(split) != 4 {
			continue
		}

		// Skip non-matching binaries
		if split[0] != fmt.Sprintf("%s/%s", r.S3BinPrefix, name) {
			continue
		}

		// Skip binaries not for this OS
		if split[2] != runtime.GOOS {
			continue
		}

		// Skip binaries not for this Arch
		if split[3] != runtime.GOARCH {
			continue
		}

		v, err := version.NewVersion(split[1])
		if err != nil {
			r.Logger.Warnf("WARNING: Unable to parse version %s - %s", split[1], err)
			continue
		}
		availableVersions = append(availableVersions, v)
	}

	if len(availableVersions) <= 0 {
		return nil, fmt.Errorf("no versions found, that's pretty weird")
	}

	sort.Sort(availableVersions)

	return availableVersions, nil
}

func (r *PrivateRepo) DownloadVersion(name, version, downloadDir string) (string, int64, error) {
	binName := fmt.Sprintf("%s-v%s-%s-%s", name, version, runtime.GOOS, runtime.GOARCH)
	downloader := s3manager.NewDownloader(r.session)

	downloadBinPath := filepath.Join(downloadDir, binName)
	downloadBin, err := os.Create(downloadBinPath)
	if err != nil {
		return "", 0, err
	}
	defer downloadBin.Close()

	bytes, err := downloader.Download(downloadBin, &s3.GetObjectInput{
		Bucket: aws.String(r.S3BinBucket),
		Key:    aws.String(fmt.Sprintf("%s/%s", name, binName)),
	})
	if err != nil {
		return "", 0, err
	}

	return downloadBinPath, bytes, nil
}

func GetCredentials(allProfiles []string) (*credentials.Credentials, error) {
	envProfile := os.Getenv("AWS_PROFILE")
	if envProfile != "" {
		allProfiles = append(allProfiles, envProfile)
	}

	var creds *credentials.Credentials
	var allErrors *multierror.Error
	for _, profile := range allProfiles {
		profileCreds := credentials.NewSharedCredentials("", profile)
		val, err := profileCreds.Get()
		if err != nil {
			allErrors = multierror.Append(allErrors, fmt.Errorf("error while finding creds: %s", err))
			continue
		}

		if val.AccessKeyID == "" {
			allErrors = multierror.Append(allErrors, fmt.Errorf("error: access key id is empty for %s", profile))
			continue
		}

		if profileCreds.IsExpired() {
			allErrors = multierror.Append(allErrors, fmt.Errorf("error: aws creds in profile %s are expired", profile))
			continue
		}

		creds = profileCreds
		break
	}

	if creds == nil {
		return nil, formatError(allProfiles, allErrors)
	}
	return creds, nil
}

func formatError(profiles []string, origErrors error) error {
	var newErrors *multierror.Error
	if e, ok := (origErrors).(*multierror.Error); ok {
		newErrors = multierror.Append(newErrors, fmt.Errorf("failed to find aws credentials in profiles: %s",
			strings.Join(profiles, ", ")),
		)
		for _, errMsg := range e.Errors {
			/*
				aws error puts a newline into the message; idk why but it looks
				ugly so remove it

				2019/01/17 09:25:40 failed to find aws credentials in profiles: confluent-dev, confluent, default
				2019/01/17 09:25:40   error while finding creds: SharedCredsLoad: failed to get profile
				caused by: section 'confluent-dev' does not exist
				2019/01/17 09:25:40   error while finding creds: SharedCredsLoad: failed to get profile
				caused by: section 'confluent' does not exist
				2019/01/17 09:25:40   error while finding creds: SharedCredsLoad: failed to get profile
				caused by: section 'default' does not exist
				2019/01/17 09:25:40 Checking for updates...

				vs

				2019/01/17 09:27:12 failed to find aws credentials in profiles: confluent-dev, confluent, default
				2019/01/17 09:27:12   error while finding creds: SharedCredsLoad: failed to get profile caused by: section 'confluent-dev' does not exist
				2019/01/17 09:27:12   error while finding creds: SharedCredsLoad: failed to get profile caused by: section 'confluent' does not exist
				2019/01/17 09:27:12   error while finding creds: SharedCredsLoad: failed to get profile caused by: section 'default' does not exist
				2019/01/17 09:27:12 Checking for updates...
			*/
			newErrors = multierror.Append(newErrors, fmt.Errorf("  %s", strings.Replace(errMsg.Error(), "\n", " ", -1)))
		}
	}
	return newErrors.ErrorOrNil()
}
