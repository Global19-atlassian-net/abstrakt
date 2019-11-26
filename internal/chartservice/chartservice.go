package chartservice

import (
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/downloader"
	"io/ioutil"
	"os"
	"path"
)

const startMeta = `apiVersion: v1
name: wormhole_constellation
description: A Helm chart for Kubernetes
version: 4.3.2
home: ""`

//CreateChart makes a new chart at the specified location
func CreateChart(name string, dir string) (*chart.Chart, error) {
	tdir, err := ioutil.TempDir("./", "output-")

	if err != nil {
		return nil, err
	}

	defer os.RemoveAll(tdir)

	err = os.Mkdir(path.Join(tdir, "templates"), 0777)

	if err != nil {
		return nil, err
	}

	files := []string{"values.yaml", path.Join("templates", "NOTES.txt")}

	for _, k := range files {
		err = ioutil.WriteFile(path.Join(tdir, k), []byte(""), 0777)
		if err != nil {
			return nil, err
		}
	}

	err = ioutil.WriteFile(path.Join(tdir, "Chart.yaml"), []byte(startMeta), 0777)

	if err != nil {
		return nil, err
	}

	cfile := &chart.Metadata{
		Name:        name,
		Description: "A Helm chart for Kubernetes",
		Type:        "application",
		Version:     "0.1.0",
		AppVersion:  "0.1.0",
		APIVersion:  chart.APIVersionV2,
	}

	err = chartutil.CreateFrom(cfile, dir, tdir)

	if err != nil {
		return nil, err
	}

	chart, err := loader.LoadDir(path.Join(dir, name))

	if err != nil {
		return nil, err
	}

	return chart, nil
}

// LoadChartFromDir loads a Helm chart from the specified director
func LoadChartFromDir(dir string) (*chart.Chart, error) {
	h, err := loader.LoadDir(dir)

	if err != nil {
		return nil, err
	}

	return h, nil
}

// SaveChartToDir takes the chart object and saves it as a set of files in the specified director
func SaveChartToDir(chart *chart.Chart, dir string) error {
	return chartutil.SaveDir(chart, dir)
}

// ZipChartToDir compresses the chart and saves it in compiled format
func ZipChartToDir(chart *chart.Chart, dir string) (string, error) {
	return chartutil.Save(chart, dir)
}

// BuildChart download charts
func BuildChart(dir string) (err error) {

	manager := downloader.Manager{
		Out:       os.Stdout,
		ChartPath: dir,
	}

	err = manager.Build()
	return
}
