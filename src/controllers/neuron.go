package controllers

import (
	"github.com/kataras/iris/context"
	"github.com/kataras/iris"

	"github.com/merakiVE/CVDI/src/models"
	"github.com/merakiVE/CVDI/core/types"
	"github.com/merakiVE/CVDI/core"

	"github.com/merakiVE/koinos/db"

	arangoDB "github.com/hostelix/aranGO"
)

type NeuronController struct {
	context core.ContextController
}

func NewNeuronController(cc core.ContextController) (NeuronController) {
	controller := NeuronController{}
	controller.SetContext(cc)
	controller.RegisterRouters()
	return controller
}

func (this *NeuronController) RegisterRouters() {
	app := this.context.App

	routerNeuron := app.Party("/neurons")
	{
		routerNeuron.Get("/", this.ListNeurons)
		routerNeuron.Get("/{neuronKey:string}", this.GetNeuron)
		routerNeuron.Post("/", this.Subscribe)
		//Action Neuron
		routerNeuron.Get("/{neuronKey:string}/actions", this.ListActions)
		routerNeuron.Post("/{neuronKey:string}/actions", this.CreateAction)
	}
}

func (this *NeuronController) SetContext(cc core.ContextController) {
	this.context = cc
}

func (this NeuronController) GetNeuron(_context context.Context) {

	var neuron models.NeuronModel

	key_neuron := _context.Params().Get("neuronKey")

	neuron.SetKey(key_neuron)
	success := db.GetModel(db.GetCurrentDatabase(), &neuron)

	if !success {

		_context.StatusCode(iris.StatusNotFound)
		_context.JSON(types.ResponseAPI{
			Message: "Error get neuron, key not found",
			Data:    nil,
			Errors:  nil,
		})
		return
	}

	_context.StatusCode(iris.StatusOK)
	_context.JSON(types.ResponseAPI{
		Message: "Neuron " + key_neuron,
		Data:    neuron,
		Errors:  nil,
	})
}

func (this NeuronController) ListNeurons(_context context.Context) {

	result := make([]models.NeuronModel, 0)
	var err error

	q := arangoDB.NewQuery(`
		FOR neuron in neurons
		RETURN neuron
	`)
	cur, err := db.GetCurrentDatabase().Execute(q)

	if err != nil {

		_context.StatusCode(iris.StatusInternalServerError)
		_context.JSON(types.ResponseAPI{
			Message: "Fail get data from database",
			Data:    nil,
			Errors:  nil,
		})
		return
	}

	err = cur.FetchBatch(&result)

	if err != nil {
		//Log Error
		_context.Application().Logger().Error(err.Error())

		_context.StatusCode(iris.StatusInternalServerError)
		_context.JSON(types.ResponseAPI{
			Message: "Fail unmarshal",
			Data:    nil,
			Errors:  nil,
		})
		return
	}

	_context.StatusCode(iris.StatusOK)
	_context.JSON(types.ResponseAPI{
		Message: "Success",
		Data:    result,
		Errors:  nil,
	})
}

func (this NeuronController) ListActions(_context context.Context) {

	var neuron models.NeuronModel

	key_neuron := _context.Params().Get("neuronKey")

	neuron.SetKey(key_neuron)
	success := db.GetModel(db.GetCurrentDatabase(), &neuron)

	if !success {

		_context.StatusCode(iris.StatusNotFound)
		_context.JSON(types.ResponseAPI{
			Message: "Error get actions, key not found",
			Data:    nil,
			Errors:  nil,
		})
		return
	}

	_context.StatusCode(iris.StatusOK)
	_context.JSON(types.ResponseAPI{
		Message: "List Actions Neuron",
		Data:    neuron.Actions,
		Errors:  nil,
	})
}

func (this NeuronController) CreateAction(_context context.Context) {
	var new_action models.ActionNeuron
	var neuron models.NeuronModel
	var err error

	key_neuron := _context.Params().Get("neuronKey")

	neuron.SetKey(key_neuron)
	success := db.GetModel(db.GetCurrentDatabase(), &neuron)

	if !success {

		_context.StatusCode(iris.StatusNotFound)
		_context.JSON(types.ResponseAPI{
			Message: "Error get neuron, key not found",
			Data:    nil,
			Errors:  nil,
		})
		return
	}

	err = _context.ReadJSON(&new_action)

	if err != nil {

		_context.StatusCode(iris.StatusInternalServerError)
		_context.JSON(types.ResponseAPI{
			Message: "Invalid data Neuron",
			Data:    nil,
			Errors:  nil,
		})
		return
	}

	//Validate action
	mesg_validation := new_action.Validate()

	if len(mesg_validation) > 0 {
		_context.StatusCode(iris.StatusInternalServerError)
		_context.JSON(types.ResponseAPI{
			Message: "Error to the add new action, invalid data",
			Data:    nil,
			Errors:  mesg_validation,
		})
		return
	}

	//Add new action to Actions model
	neuron.Actions = append(neuron.Actions, new_action)

	//Update actions model in database
	err = db.ReplaceModel(db.GetCurrentDatabase(), neuron)

	if err != nil {
		_context.StatusCode(iris.StatusInternalServerError)
		_context.JSON(types.ResponseAPI{
			Message: "Error to the add new action to neuron",
			Data:    nil,
			Errors:  nil,
		})
		return
	}

	_context.StatusCode(iris.StatusOK)
	_context.JSON(types.ResponseAPI{
		Message: "Added new action success",
		Data:    new_action,
		Errors:  nil,
	})
}

func (this NeuronController) Subscribe(_context context.Context) {
	var _neuron models.NeuronModel

	var err error

	err = _context.ReadJSON(&_neuron)

	if err != nil {

		//Log error
		_context.Application().Logger().Error(err.Error())

		_context.StatusCode(iris.StatusInternalServerError)
		_context.JSON(types.ResponseAPI{
			Message: "Invalid data Neuron",
			Data:    nil,
			Errors:  nil,
		})
		return
	}

	success := db.SaveModel(db.GetCurrentDatabase(), &_neuron)

	if success {
		_context.StatusCode(iris.StatusCreated)
		_context.JSON(types.ResponseAPI{
			Message: "Neuron subscribe successfully",
			Data:    nil,
			Errors:  nil,
		})

	} else {
		_context.StatusCode(iris.StatusOK)
		_context.JSON(types.ResponseAPI{
			Message: "Error subscribing neuron, invalid data",
			Data:    nil,
			Errors:  _neuron.GetValidationErrors(),
		})
	}
}
