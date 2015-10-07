package main

import (
	"log"
	"os"
	"text/template"

	"github.com/codegangsta/cli"
)

func requiredStringFlag(c *cli.Context, name string) string {
	var val string
	if val = c.String(name); val == "" {
		log.Fatalf("missing required flag: --%s", name)
	}
	return val
}

func main() {
	app := cli.NewApp()
	app.Name = "convox-install"
	app.Usage = "runs `convox install` with our patched formation.json"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "ami",
			Value: "",
			Usage: "Our custom AMI for us-east-1",
		},
		cli.StringFlag{
			Name:  "sumo-token",
			Value: "",
			Usage: "Sumo Logic Collector Token",
		},
	}

	app.Action = func(c *cli.Context) {
		var params = make(map[string]string)
		params["sumo_collector_token"] = requiredStringFlag(c, "sumo-token")
		params["ami_us_east_1"] = requiredStringFlag(c, "ami")

		t, err := template.ParseFiles("./add-logstash-to-formation.patch.tmpl")
		if err != nil {
			log.Fatalf("parsing: %s", err)
		}

		err = t.Execute(os.Stdout, params)
		if err != nil {
			log.Fatalf("execution: %s", err)
		}
	}

	app.Run(os.Args)
}
