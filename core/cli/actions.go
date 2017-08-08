package cli

import (
	"github.com/urfave/cli"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"

	"github.com/merakiVE/CVDI/src/controllers"
	packageConfig "github.com/merakiVE/CVDI/core/config"
	"github.com/merakiVE/CVDI/core"
	"strings"
	//"github.com/kataras/iris/context"
)

var (
	PORT_SERVER = ":8101"
)

func RunServer(c *cli.Context) error {

	///Iris
	app := iris.New()

	// Init Configuration var
	config := packageConfig.Configuration{}

	//app.Configure(iris.WithConfiguration(iris.YAML("./config_iris.yml")))

	APILogger := logger.New(logger.Config{
		// Status displays status code
		Status: true,
		// IP displays request's remote address
		IP: true,
		// Method displays the http method
		Method: true,
		// Path displays the request path
		Path: true,
	})

	app.Use(APILogger)

	//Load configuration
	config.Load()

	//Context Controller
	contextController := core.ContextController{App: app, Config: config}

	//Init Controllers
	controllers.NewAuthController(contextController)
	controllers.NewNeuronController(contextController)
	controllers.NewUserController(contextController)

	//Verify if pass argument port
	if c.NArg() > 0 {
		PORT_SERVER = c.Args().First()

		if !strings.HasPrefix(PORT_SERVER, ":") {
			PORT_SERVER = ":" + PORT_SERVER
		}
	}

	app.Run(iris.Addr(PORT_SERVER), iris.WithCharset("UTF-8"))

	return nil
}
