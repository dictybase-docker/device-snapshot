package command

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/dictyBase-docker/device-snapshot/internal/logger"
	"github.com/urfave/cli"
)

var snapDevices []chromedp.Device = []chromedp.Device{
	device.IPhone8,
	device.IPhone8landscape,
	device.IPad,
	device.IPadlandscape,
	device.GalaxyS5,
	device.GalaxyS5landscape,
}

type ChromeWebsocketInfo struct {
	Url             string `json:"webSocketDebuggerUrl"`
	WebkitVersion   string `json:"WebKit-Version"`
	V8Version       string `json:"V8-Version"`
	UserAgent       string `json:"User-Agent"`
	ProtocolVersion string `json:"Protocol-Version"`
	Browser         string `json:"Browser"`
}

func GenerateSnapshot(c *cli.Context) error {
	r, err := http.Get(
		fmt.Sprintf(
			"http://%s:%s/json/version",
			c.String("remote-chrome-host"), c.String("remote-chrome-port")),
	)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in http response from remote chrome %s", err),
			2,
		)
	}
	defer r.Body.Close()
	chrInfo := new(ChromeWebsocketInfo)
	if err := json.NewDecoder(r.Body).Decode(chrInfo); err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in decoding json from remote chrome %s", err),
			2,
		)
	}
	// create context
	actxt, cancelActxt := chromedp.NewRemoteAllocator(context.Background(), chrInfo.Url)
	defer cancelActxt()
	ctx, cancel := chromedp.NewContext(actxt)
	defer cancel()

	// capture screenshot of an element
	//var buf []byte
	l := logger.GetLogger(c)
	for _, p := range c.StringSlice("path") {
		url := fmt.Sprintf("%s/%s", c.String("host"), p)
		for _, d := range snapDevices {
			var b []byte
			err := chromedp.Run(ctx,
				fullScreenshot(url, 90, &b, d),
			)
			fname := fmt.Sprintf("snapshot-%s-%s.png", strings.Replace(p, "/", "-", -1), d.Device().Name)
			if err != nil {
				return cli.NewExitError(
					fmt.Sprintf("error in running remote chrome for url %s %s", url, err),
					2,
				)
			}
			l.Debugf("took snapshot of url %s", url)

			fname = filepath.Join(c.String("output"), fname)
			if err := ioutil.WriteFile(fname, b, 0644); err != nil {
				return cli.NewExitError(
					fmt.Sprintf("error in saving file %s %s", fname, err),
					2,
				)
			}
			l.Debugf("saved file of snapshot %s", fname)
		}
	}
	return nil
}

// fullScreenshot takes a screenshot of the entire browser viewport.
func fullScreenshot(urlstr string, quality int64, res *[]byte, d chromedp.Device) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Emulate(d),
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
		chromedp.Emulate(device.Reset),
	}
}
