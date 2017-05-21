package main

import (
    "runtime"
    "fmt"
    "math"
    "bufio"
    "os"
)

type polar struct {
    radius float64
    θ      float64
}

type cartesian struct {
    x float64
    y float64
}

var prompt = "Enter a radius and an angle (in degree), e.g., 12.5 90, " + "or %s to quit."

func init() {
    if runtime.GOOS == "windows" {
        prompt = fmt.Sprintf(prompt, "Ctrl + Z, Enter")
    } else {
        prompt = fmt.Sprintf(prompt, "Ctrl + D")
    }
}

func main() {
    questions := make(chan polar)
    defer close(questions)
    answers := createSolver(questions)
    defer close(answers)
    interact(questions, answers)
}

func createSolver(questions chan polar) chan cartesian {
    answer := make(chan cartesian)
    go func() {
        for {
            polarCoord := <-questions
            θ := polarCoord.θ * math.Pi / 180.0
            x := polarCoord.radius * math.Cos(θ)
            y := polarCoord.radius * math.Sin(θ)
            answer <- cartesian{x, y}
        }
    }()
    return answer
}

const result = "Polar radius=%.02f θ=%.02f° ===> Cartesian x=%.02f y=%.02f\n"

func interact(questions chan polar, answer chan cartesian) {
    reader := bufio.NewReader(os.Stdin)
    fmt.Println(prompt)
    for {
        fmt.Print("Redius and angle: ")
        line, err := reader.ReadString('\n')
        if err != nil {
            break
        }
        var radius, θ float64
        if _, err := fmt.Sscanf(line, "%f %f", &radius, &θ); err != nil {
            fmt.Println(os.Stderr, "Invalid input")
            continue
        }
        questions <- polar{radius, θ}
        corrd := <-answer
        fmt.Printf(result, radius, θ, corrd.x, corrd.y)
    }
    fmt.Println()
}
