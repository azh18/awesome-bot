package browser

import (
	"github.com/tebeka/selenium"
)

func ExtractTextValue(wd selenium.WebDriver, xpath string) (string, error) {
	elem, err := wd.FindElement(selenium.ByXPATH, xpath)
	if err != nil {
		return "", err
	}
	value, err := elem.Text()
	if err != nil {
		return "", err
	}
	return value, nil
}
