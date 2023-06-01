package main

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"time"

	"github.com/JeremyLoy/config"
	chrome "github.com/bjornpagen/goplay"
	openai "github.com/sashabaranov/go-openai"
)

type Env struct {
	OpenAIKey   string `config:"OPENAI_KEY"`
	BrowserPort string `config:"BROWSER_PORT"`
}

func validate(e Env) error {
	t := reflect.TypeOf(e)
	var missing []string
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("config")
		if tag == "" {
			continue
		}
		v := reflect.ValueOf(e).Field(i).String()
		if v == "" {
			missing = append(missing, tag)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("required env variables: %v", missing)
	}

	return nil
}

func main() {
	e := Env{}
	err := config.FromEnv().To(&e)
	if err != nil {
		panic(err)
	}

	err = validate(e)
	if err != nil {
		panic(err)
	}

	// check port is valid uint16
	var port uint16
	_, err = fmt.Sscanf(e.BrowserPort, "%d", &port)
	if err != nil {
		panic(err)
	}

	c := New(e)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := c.run(ctx); err != nil {
		panic(err)
	}
}

type client struct {
	Env

	oc *openai.Client
}

func New(e Env) *client {
	return &client{
		Env: e,
		oc:  openai.NewClient(e.OpenAIKey),
	}
}

func (c *client) run(ctx context.Context) error {
	port, _ := strconv.ParseUint(c.Env.BrowserPort, 10, 16)
	b, err := chrome.New(chrome.AttachToExisting(uint16(port)))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err = b.Start(ctx)
	if err != nil {
		return err
	}

	firstone := "https://boards.greenhouse.io/excelsportsmanagement/jobs/4234907005"
	otherone := "https://boards.greenhouse.io/clubmonaco/jobs/4929420"
	_ = firstone
	err = b.Navigate(ctx, otherone)
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	err = fillMainFields(ctx, b)
	if err != nil {
		return err
	}

	time.Sleep(time.Duration(20+rand.Intn(40)) * time.Second)

	return nil
}

func fillMainFields(ctx context.Context, b *chrome.Browser) error {
	err := fillField(ctx, b, `#first_name`, "Lu")
	if err != nil {
		return err
	}

	return nil

	err = fillField(ctx, b, `#last_name`, "Cifer")
	if err != nil {
		return err
	}

	err = fillField(ctx, b, `#email`, "lucifer@gmail.com")
	if err != nil {
		return err
	}

	err = fillField(ctx, b, `#phone`, "+16666666666")
	if err != nil {
		return err
	}

	err = fillFile(ctx, b, `#s3_upload_for_resume input[type="file"]`, "/home/satan/resume.pdf")
	if err != nil {
		return err
	}

	return nil
}
