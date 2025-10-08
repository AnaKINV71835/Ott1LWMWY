// 代码生成时间: 2025-10-08 23:43:48
package main

import (
    "math/rand"
    "time"
    "fmt"
    "net/http"
    "github.com/gofiber/fiber/v2"
)

// Neuron represents a single node in a neural network
type Neuron struct {
    weights   []float64
    bias      float64
    activation float64
    output    float64
}

// NeuralNetwork represents a simple neural network with one input layer,
// one hidden layer, and one output layer.
type NeuralNetwork struct {
    inputLayer   []Neuron
    hiddenLayer []Neuron
    outputLayer []Neuron
}

// NewNeuron creates a new neuron with random weights and a bias
func NewNeuron(numberOfInputs int) *Neuron {
    neuron := &Neuron{
        activation: 0.0,
        output: 0.0,
        weights:   make([]float64, numberOfInputs),
        bias:      rand.Float64() * 2 - 1,
    }
    for i := range neuron.weights {
        neuron.weights[i] = rand.Float64() * 2 - 1
    }
    return neuron
}

// NewNeuralNetwork initializes a new neural network with a specified number of inputs,
// hidden neurons, and outputs.
func NewNeuralNetwork(inputs, hidden, outputs int) *NeuralNetwork {
    nn := &NeuralNetwork{
        inputLayer:   make([]Neuron, inputs),
        hiddenLayer: make([]Neuron, hidden),
        outputLayer: make([]Neuron, outputs),
    }
    for i := range nn.inputLayer {
        nn.inputLayer[i] = *NewNeuron(inputs)
    }
    for i := range nn.hiddenLayer {
        nn.hiddenLayer[i] = *NewNeuron(hidden)
    }
    for i := range nn.outputLayer {
        nn.outputLayer[i] = *NewNeuron(1) // Output layer has one input from the hidden layer
    }
    return nn
}

// Sigmoid is the activation function used in the neural network
func Sigmoid(x float64) float64 {
    return 1 / (1 + math.Exp(-x))
}

// FeedForward propagates inputs through the network and calculates the output
func (nn *NeuralNetwork) FeedForward(input []float64) []float64 {
    // Calculate input layer outputs
    for i := range nn.inputLayer {
        nn.inputLayer[i].output = input[i]
    }
    // Calculate hidden layer outputs
    for i := range nn.hiddenLayer {
        sum := 0.0
        for j := range nn.inputLayer {
            sum += nn.inputLayer[j].output * nn.hiddenLayer[i].weights[j]
        }
        sum += nn.hiddenLayer[i].bias
        nn.hiddenLayer[i].activation = Sigmoid(sum)
        nn.hiddenLayer[i].output = nn.hiddenLayer[i].activation
    }
    // Calculate output layer outputs
    for i := range nn.outputLayer {
        sum := 0.0
        for j := range nn.hiddenLayer {
            sum += nn.hiddenLayer[j].output * nn.outputLayer[i].weights[j]
        }
        sum += nn.outputLayer[i].bias
        nn.outputLayer[i].activation = Sigmoid(sum)
        nn.outputLayer[i].output = nn.outputLayer[i].activation
    }
    return []float64{nn.outputLayer[0].output}
}

// StartServer starts the Fiber web server
func StartServer() {
    app := fiber.New()

    app.Get("/predict", func(c *fiber.Ctx) error {
        // Example usage of the neural network
        nn := NewNeuralNetwork(2, 2, 1)
        input := []float64{0.5, 0.5}
        output := nn.FeedForward(input)
        return c.JSON(fiber.Map{
            "input":  input,
            "output": output,
        })
    })

    app.Listen(":3000")
}

func main() {
    rand.Seed(time.Now().UnixNano())
    StartServer()
}
