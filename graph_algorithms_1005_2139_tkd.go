// 代码生成时间: 2025-10-05 21:39:55
package main

import (
    "fmt"
    "math"
    "net/http"
    "strings"
    "golang.org/x/exp/constraints"
)

// Graph represents an adjacency list representation of a graph.
type Graph[T constraints.Ordered] struct {
    AdjList map[T]*list.List
    V       int
    E       int
}

// NewGraph creates a new graph.
func NewGraph[T constraints.Ordered]() *Graph[T] {
    return &Graph[T]{
        AdjList: make(map[T]*list.List),
        V:       0,
        E:       0,
    }
}

// AddEdge adds an edge to the graph.
func (g *Graph[T]) AddEdge(v, w T) {
    g.V++
    g.V++
    g.E++
    g.get(v).PushBack(w)
    g.get(w).PushBack(v)
}

// get returns a list for a given vertex.
func (g *Graph[T]) get(v T) *list.List {
    if _, present := g.AdjList[v]; !present {
        g.AdjList[v] = list.New()
    }
    return g.AdjList[v]
}

// BFS performs a breadth-first search on the graph.
func (g *Graph[T]) BFS(s T) []T {
    visited := make(map[T]bool)
    queue := list.New()
    parent := make(map[T]T)
    visited[s] = true
    queue.PushBack(s)
    parent[s] = s

    var path []T
    for queue.Len() > 0 {
        u := queue.Front().Value.(T)
        fmt.Printf("Visited: %v
", u)
        queue.Remove(queue.Front())
        for _, w := range g.AdjList[u] {
            if !visited[w] {
                visited[w] = true
                parent[w] = u
                queue.PushBack(w)
            }
        }
    }

    for v := range g.AdjList {
        if v != s && parent[v] == nil {
            fmt.Printf("Vertex %v has no path to %v
", v, s)
        } else if v != s {
            fmt.Printf("Path to %v: %v -> %v
", v, parent[v], v)
        }
    }
    return path
}

// DFS performs a depth-first search on the graph.
func (g *Graph[T]) DFS(s T) {
    visited := make(map[T]bool)
    stack := list.New()
    visited[s] = true
    stack.PushBack(s)
    fmt.Printf("Visited: %v
", s)

    for stack.Len() > 0 {
        u := stack.Back().Value.(T)
        stack.Remove(stack.Back())
        for _, w := range g.AdjList[u] {
            if !visited[w] {
                visited[w] = true
                stack.PushBack(w)
                fmt.Printf("Visited: %v
", w)
            }
        }
    }
}

// shortestPath returns the shortest path from a source to a destination in the graph.
func (g *Graph[T]) ShortestPath(s, d T) []T {
    visited := make(map[T]bool)
    pred := make(map[T]T)
    var stack []T
    visited[s] = true
    stack = append(stack, s)

    for len(stack) > 0 {
        u := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        for _, v := range g.AdjList[u] {
            if !visited[v] {
                visited[v] = true
                pred[v] = u
                stack = append(stack, v)
                if v == d {
                    break
                }
            }
        }
        if len(stack) == 0 && u == d {
            break
        }
    }

    path := []T{d}
    for u := d; u != s; {
        u = pred[u] // will error if u not in pred
        path = append([]T{u}, path...)
    }
    return path
}

func main() {
    g := NewGraph[int]()
    g.AddEdge(0, 1)
    g.AddEdge(0, 2)
    g.AddEdge(1, 2)
    g.AddEdge(2, 0)
    g.AddEdge(2, 3)
    g.AddEdge(3, 3)

    fmt.Println("BFS: ")
    g.BFS(0)

    fmt.Println("DFS: ")
    g.DFS(0)

    path := g.ShortestPath(0, 3)
    fmt.Println("Shortest Path from 0 to 3: ", path)
}
