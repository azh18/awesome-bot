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

	// options
	driverPath string
	port       int

	needStartDriver = false
)

func InitSelenium(chromeDriverPath string, chromeDriverPort int) {
	driverPath = chromeDriverPath
	port = chromeDriverPort
}

type SeleniumService struct {
	service   *selenium.Service
	WebDriver selenium.WebDriver
	stopCh    chan struct{}
}

func (s *SeleniumService) Stop() {
	s.WebDriver.Quit()
	if s.service != nil {
		s.service.Stop()
	}
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
		klog.Errorf("new selenium service error: %s", err.Error())
		return nil
	}
	klog.Infof("selenium service created for %s", key)
	f.browsers[key] = service
	return service
}

func newSeleniumService() (*SeleniumService, error) {
	if driverPath != "" {
		needStartDriver = true
	}

	var service *selenium.Service
	if needStartDriver {
		ops := []selenium.ServiceOption{}
		var err error
		service, err = selenium.NewChromeDriverService(driverPath, port, ops...)
		if err != nil {
			return nil, err
		}
	}

	wd, err := selenium.NewRemote(selenium.Capabilities{
		"browserName": "chrome",
	}, fmt.Sprintf("http://127.0.0.1:%d/wd/hub", port))

	if err != nil {
		if service != nil {
			service.Stop()
		}
		return nil, err
	}
	stopCh := make(chan struct{})

	return &SeleniumService{
		service:   service,
		WebDriver: wd,
		stopCh:    stopCh,
	}, nil
}
