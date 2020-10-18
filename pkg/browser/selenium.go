package browser

import (
	"fmt"
	"sync"

	"github.com/tebeka/selenium"
	"k8s.io/klog/v2"
)

var (
	once    sync.Once
	factory *SeleniumFactory
)

const (
	// chromeDriverPath = "/opt/selenium/chromedriver-86.0.4240.22"
	port = 4444
)

type SeleniumService struct {
	// service   *selenium.Service
	WebDriver selenium.WebDriver
	stopCh    chan struct{}
}

func (s *SeleniumService) Stop() {
	s.WebDriver.Quit()
	// s.service.Stop()
}

type SeleniumFactory struct {
	browsers map[string]*SeleniumService
	sync.RWMutex
}

func GetSeleniumFactory() *SeleniumFactory {
	once.Do(func() {
		factory = &SeleniumFactory{
			browsers: map[string]*SeleniumService{},
		}
	})
	return factory
}

func (f *SeleniumFactory) Stop() {
	f.Lock()
	defer f.Unlock()
	for _, b := range f.browsers {
		b.Stop()
	}
}

func (f *SeleniumFactory) GetSelenium(key string) *SeleniumService {
	f.RLock()
	if f.browsers[key] != nil {
		f.RUnlock()
		return f.browsers[key]
	}
	f.RUnlock()
	f.Lock()
	defer f.Unlock()

	service, err := newSeleniumService()
	if err != nil {
		klog.Infof("new selenium service error: %s", err.Error())
		return nil
	}
	f.browsers[key] = service
	return service
}

func newSeleniumService() (*SeleniumService, error) {
	// ops := []selenium.ServiceOption{}

	// service, err := selenium.NewChromeDriverService(chromeDriverPath, port, ops...)
	// if err != nil {
	// 	return nil, err
	// }

	wd, err := selenium.NewRemote(selenium.Capabilities{
		"browserName": "chrome",
	}, fmt.Sprintf("http://127.0.0.1:%d/wd/hub", port))
	if err != nil {
		// service.Stop()
		return nil, err
	}
	stopCh := make(chan struct{})

	return &SeleniumService{
		// service:   service,
		WebDriver: wd,
		stopCh:    stopCh,
	}, nil
}
