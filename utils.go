package artifacts

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/getsentry/sentry-go"
)

var analysisEventsLogFmt = "%s,%s,%s,%s,%s,%s,%s,%s,%s,%d\n"

// Returns bearer token that is used to authenticate while
// interacting with the k8s REST API
// Utilized by janus and atlas.
func GetNewBearerToken(tokenFilePath string) (string, error) {
	authToken, err := ioutil.ReadFile(tokenFilePath)
	if err != nil {
		return "", err
	}
	bearer := "Bearer " + string(authToken)
	return bearer, nil
}

// Utility function to retry any routine
func Retry(attempts int, sleep time.Duration, f func() error) (err error) {
	for i := 0; ; i++ {
		err = f()
		if err == nil {
			return
		}
		if i >= (attempts - 1) {
			break
		}
		time.Sleep(sleep)
		log.Println("Retrying job after error:", err)
	}
	sentry.CaptureException(err)
	return fmt.Errorf("Failed to trigger the job aftert %d attempts, last error: %s", attempts, err)
}

// Get a HTTP client for interacting with k8s REST API
func GetNewHTTPClient(certFilePath string) (*http.Client, error) {
	var httpClient *http.Client

	// if auth is true, send request with the certificate
	caCert, err := ioutil.ReadFile(certFilePath)
	if err != nil {
		return httpClient, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	httpClient = &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
				RootCAs:            caCertPool,
				MinVersion:         tls.VersionTLS12,
				MaxVersion:         0,
			},
		},
	}

	return httpClient, nil
}

// Fibonacci returns successive Fibonacci numbers starting from 1
func fibonacci() func() int {
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}
}

// FibonacciNext returns next number in Fibonacci sequence greater than start
func fibonacciNextNum(start int) int {
	fib := fibonacci()
	num := fib()
	for num <= start {
		num = fib()
	}
	return num
}

// Gets the duration interval of retries based on fibonacci series
func GetRetryTimeout(currentTimeout int) (int, time.Duration) {
	retryTimeout := fibonacciNextNum(currentTimeout)
	durationString := fmt.Sprintf("%vs", retryTimeout)
	duration, _ := time.ParseDuration(durationString)
	return retryTimeout, duration
}

// Watches the broker config for any "WRITE" events
// Exits with status code 1 after any changes to config
func WatchBrokerConfigForChanges(filePath string, reloadFunc func() error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println(err)
		sentry.CaptureException(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println(event)
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("Received event:", event)
					log.Println("modified file:", event.Name)

					if strings.HasSuffix(event.Name, filePath) {
						watcher.Close()
						if reloadErr := reloadFunc(); reloadErr != nil {
							log.Fatalln("Failed to reload janus. Exiting.") // skipcq
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Watcher error:", err)
			}
		}
	}()

	err = watcher.Add(filePath)
	if err != nil {
		log.Println(err)
		sentry.CaptureException(err)
	}
	<-done
}

// Utility to do a rolling restart deployment of a kubernetes pod
func TriggerDeploymentRestart(auth bool, podName, patchData, baseURL, tokenPath string, httpClient *http.Client) error {
	// Configuring the request
	restartDeploymentRequest, _ := http.NewRequest("PATCH", baseURL+"/apis/apps/v1/namespaces/default/deployments/"+podName, bytes.NewBuffer([]byte(patchData)))
	restartDeploymentRequest.Header.Set("Content-Type", "application/strategic-merge-patch+json")

	// If auth is true, set Authorization header
	if auth {
		bearer, err := GetNewBearerToken(tokenPath)
		if err != nil {
			log.Println(err)
			return err
		}
		restartDeploymentRequest.Header.Add("Authorization", bearer)
	}

	// If not, try anyways without the Bearer Token in the Header
	err := Retry(5, 2*time.Second, func() (err error) {
		restartDeploymentResponse, err := httpClient.Do(restartDeploymentRequest)
		if err != nil {
			log.Println(err)
			return err
		}
		if restartDeploymentResponse.StatusCode != 200 && restartDeploymentResponse.StatusCode != 201 {
			responseBody, _ := ioutil.ReadAll(restartDeploymentResponse.Body)
			log.Println(string(responseBody))
			return errors.New("Unable to create a kubernetes job")
		}
		defer restartDeploymentResponse.Body.Close()
		return nil
	})
	if err != nil {
		log.Print(err.Error())
		return err
	}
	return nil
}

// AnalysisEventsLog represents the struct that contains the fields
// that need to be logged for tracking analysis pipeline's performance.
type AnalysisEventsLog struct {
	RunID         string
	RunSerial     string
	CheckSequence string
	Repository    string
	Shortcode     string
	CommitSHA     string
	IsFullRun     string
}

// logAnalysisEventTimestamp logs the timestamp at various analysis stages.
func (a *AnalysisEventsLog) LogAnalysisEventTimestamp(runType, stage string) {
	fmt.Printf(analysisEventsLogFmt,
		runType,
		a.RunID,
		a.RunSerial,
		a.CheckSequence,
		a.Shortcode,
		a.Repository,
		a.CommitSHA,
		a.IsFullRun,
		stage,
		time.Now().Unix(),
	)
}
