package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"golang.org/x/crypto/ssh/terminal"
)

type Cookie struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Value  string `json:"value"`
	Domain string `json:"domain"`
	Secure bool   `json:"secure"`
}

var (
	errWrongCredentials = errors.New("Invalid username or password")
	timeout             = 5 * time.Second
	pollInterval        = 500 * time.Millisecond
)

func ExtractCookies() error {
	service, err := selenium.NewChromeDriverService("./chromedriver", 4444)
	if err != nil {
		return err
	}
	defer service.Stop()

	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{
		Args: []string{
			"window-size=1920x1080",
			"--no-sandbox",
			"--disable-dev-shm-usage",
			"disable-gpu",
			// "--headless",
		},
	})

	driver, err := selenium.NewRemote(caps, "")
	if err != nil {
		return err
	}
	defer driver.Quit()

	var username, password string

	done := make(chan struct{})
	errChan := make(chan error)

	go func() {
		username, password, err = getUsernameAndPassword()
		if err != nil {
			errChan <- err
			return
		}
		done <- struct{}{}
	}()

	if err := driver.Get("https://leetcode.com/accounts/login/"); err != nil {
		return err
	}

	select {
	case <-done:
	case err := <-errChan:
		return err
	}

	cookies, err := login(driver, username, password)
	if errors.Is(err, errWrongCredentials) {
		fmt.Println("Wrong username/password, please try again")
		for retry := 0; retry < 2; retry++ {
			username, password, err = getUsernameAndPassword()
			if err != nil {
				return err
			}
			cookies, err = login(driver, username, password)
			if err != nil {
				if errors.Is(err, errWrongCredentials) {
					fmt.Println("Wrong username/password, please try again")
					continue
				}
				return err
			}
			err = nil
			break
		}
		fmt.Println("Retried 3 times. Please check username/password before try again")
		return errWrongCredentials
	}

	if err != nil {
		return err
	}

	httpCookies := make([]*Cookie, 0, len(cookies))
	for _, cookie := range cookies {
		httpCookies = append(httpCookies, &Cookie{
			Name:   cookie.Name,
			Path:   cookie.Path,
			Value:  cookie.Value,
			Domain: cookie.Domain,
			Secure: cookie.Secure,
		})
	}

	cookieJson, err := json.Marshal(httpCookies)
	if err != nil {
		return err
	}

	if err := os.WriteFile(_cookiePath, cookieJson, 0777); err != nil {
		return err
	}

	return nil
}

func getUsernameAndPassword() (username, password string, err error) {
	fmt.Println("Please input username")
	fmt.Scan(&username)
	fmt.Println("Please input password")
	passwordByte, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}
	password = string(passwordByte)
	return
}

func login(driver selenium.WebDriver, usernameText, passwordText string) ([]selenium.Cookie, error) {
	if err := driver.WaitWithTimeoutAndInterval(func(wd selenium.WebDriver) (bool, error) {
		username, err := driver.FindElement(selenium.ByID, "id_login")
		if err != nil || username == nil {
			return false, nil
		}
		return username.IsDisplayed()
	}, timeout, pollInterval); err != nil {
		return nil, err
	}

	username, err := driver.FindElement(selenium.ByID, "id_login")
	if err != nil {
		return nil, err
	}
  // Clearing username & password fields if they have values
	if err := username.Click(); err != nil {
		return nil, err
	}
	if err := username.SendKeys(selenium.ControlKey + "a" + selenium.BackspaceKey); err != nil {
		return nil, err
	}
	if err := username.SendKeys(usernameText); err != nil {
		return nil, err
	}

	password, err := driver.FindElement(selenium.ByID, "id_password")
	if err != nil {
		return nil, err
	}
	if err := password.Click(); err != nil {
		return nil, err
	}
	if err := password.SendKeys(selenium.ControlKey + "a" + selenium.BackspaceKey); err != nil {
		return nil, err
	}
  time.Sleep(50 * time.Millisecond)
	if err := password.SendKeys(passwordText); err != nil {
		return nil, err
	}

	signInButton, err := driver.FindElement(selenium.ByID, "signin_btn")
	if err != nil {
		return nil, err
	}
	if err := signInButton.Click(); err != nil {
		return nil, err
	}

	if err = driver.WaitWithTimeoutAndInterval(
		func(wd selenium.WebDriver) (bool, error) {
			navbar, err := driver.FindElement(selenium.ByID, "home-app")
			fmt.Println("Login...")
			if err != nil || navbar == nil {
				return false, nil
			}
			return navbar.IsDisplayed()
		}, timeout, pollInterval); err != nil {
		return nil, errWrongCredentials
	}

	cookies, err := driver.GetCookies()
	if err != nil {
		return nil, err
	}

	return cookies, nil
}
