package main

import (
	"bufio"
	hp "container/Heap"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Path struct {
	Value int
	Rooms []string
}

type minPath []Path

func (h minPath) Len() int           { return len(h) }
func (h minPath) Less(i, j int) bool { return h[i].Value < h[j].Value }
func (h minPath) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *minPath) Push(x interface{}) {
	*h = append(*h, x.(Path))
}

func (h *minPath) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type Heap struct {
	values *minPath
}

func newHeap() *Heap {
	return &Heap{values: &minPath{}}
}

func (h *Heap) push(p Path) {
	hp.Push(h.values, p)
}

func (h *Heap) pop() Path {
	i := hp.Pop(h.values)
	return i.(Path)
}

type Room struct {
	RoomName  string
	Weight    int
	StartRoom string
	EndRoom   string
	Links     [][]string
	Dijkstra  [][]string
}

type Ant struct {
	NumberOfAnts int
	Next         *Room
}

type Graph struct {
	Rooms map[string][]Room
	Paths []string
}

func newGraph() *Graph {
	return &Graph{Rooms: make(map[string][]Room)}
}

func (g *Graph) addRoom(start, end string, Weight int) {
	g.Rooms[start] = append(g.Rooms[start], Room{RoomName: end, Weight: Weight})

	g.Rooms[end] = append(g.Rooms[end], Room{RoomName: start, Weight: Weight})
}

func (g *Graph) getRooms(RoomName string) []Room {
	return g.Rooms[RoomName]
}

func (g *Graph) getPath(start, end string) []string { // int, []string
	h := newHeap()
	h.push(Path{Value: 0, Rooms: []string{start}})
	visited := make(map[string]bool)

	for len(*h.values) > 0 {
		// Find the nearest yet to visit RoomName
		p := h.pop()
		RoomName := p.Rooms[len(p.Rooms)-1]

		if visited[RoomName] {
			continue
		}

		if RoomName == end {
			g.Paths = p.Rooms
			return p.Rooms
			// return p.Value, p.Rooms
		}

		for _, e := range g.getRooms(RoomName) {
			if !visited[e.RoomName] {
				// We calculate the total spent so far plus the cost and the Path of getting here
				h.push(Path{Value: p.Value + e.Weight, Rooms: append([]string{}, append(p.Rooms, e.RoomName)...)})
			}
		}

		visited[RoomName] = true
	}
	// return 0, nil
	return nil
}

// func (a *Ant) checkRoom()

func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func Chunk(slice []string, size int) [][]string {
	var chunk [][]string
	for i := 0; i < len(slice); i += size {

		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunk = append(chunk, slice[i:end])
	}

	return chunk
}

// Queue represents a queue that holds a slice
type Queue struct {
	ant []string
}

// Enqueue adds a value at the end
func (q *Queue) Enqueue(i string) {
	q.ant = append(q.ant, i)
}

// Dequeue
func (q *Queue) Dequeue() string {
	toRemove := q.ant[len(q.ant)-1]
	q.ant = q.ant[1:]
	return toRemove
}

func main() {
	// start := time.Now()

	var strArr []string
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)

	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		strArr = append(strArr, fileScanner.Text())
	}

	room := &Room{}
	ant := &Ant{}

	if strArr[0] <= "0" {
		fmt.Println("Error, ant colony has died! (Number of ants must be at least 1.)")
		os.Exit(0)
	}
	NumberOfAnts, err := strconv.Atoi(strArr[0])
	if err == nil {
		ant.NumberOfAnts = NumberOfAnts
	}

	RemoveIndex(strArr, 0)
	strArr = strArr[:len(strArr)-1]
	for i := 0; i < len(strArr); i++ {
		if strArr[i] == "##start" {
			strArr = append(strArr, strArr[i+1])
			RemoveIndex(strArr, i+1)
			strArr = strArr[:len(strArr)-1]
		}
		if strArr[i] == "##end" {
			strArr = append(strArr, strArr[i+1])
			RemoveIndex(strArr, i+1)
			strArr = strArr[:len(strArr)-1]
		}
	}
	replaceWordHyphenWord := regexp.MustCompile(`\w+\-+\w+`)
	replaceHashtagWord := regexp.MustCompile(`\#+\w+`)
	deleteComment := regexp.MustCompile(`comment`)
	joinStrArr := strings.Join(strArr, " ")
	result := replaceWordHyphenWord.ReplaceAllString(joinStrArr, "")
	result = replaceHashtagWord.ReplaceAllString(result, "")
	result = deleteComment.ReplaceAllString(result, "")
	roomsWithCoordinates := strings.Fields(result)

	rooms := Chunk(roomsWithCoordinates, 3)

	for i := 0; i < len(rooms); i++ {
		room.StartRoom = rooms[len(rooms)-2][0]
		room.EndRoom = rooms[len(rooms)-1][0]
	}

	var links []string
	for index := range strArr {
		if strings.Contains(strArr[index], "-") {
			link := strings.Split(strArr[index], "-")
			for _, eachLine := range link {
				links = append(links, eachLine)
			}
		}
	}

	room.Links = Chunk(links, 2)
	Graph := newGraph()
	for i := 0; i < len(room.Links); i++ {
		if room.Links[i][0] != room.Links[i][1] {
			Graph.addRoom(room.Links[i][0], room.Links[i][1], 1)
		} else {
			fmt.Println("Ants cannot process pheromones! (Room cannot link to itself e.g. 3-3)")
			os.Exit(0)

		}
	}

	fmt.Println("Dijkstra")

	for i := 0; i < len(room.Links); i++ {
		if room.StartRoom == room.Links[i][0] {
			path := (Graph.getPath(room.Links[i][1], room.EndRoom))
			room.Dijkstra = append(room.Dijkstra, path)
		}
		if room.StartRoom == room.Links[i][1] {
			path := (Graph.getPath(room.Links[i][1], room.EndRoom))
			room.Dijkstra = append(room.Dijkstra, path)
		}

	}

	q := Queue{}

	// shortest path len
	min := len(room.Dijkstra[0])
	for _, v := range room.Dijkstra {
		if len(v) < min {
			min = len(v)
		}
	}
	// to prevent index out of range
	var empty []string
	room.Dijkstra = append(room.Dijkstra, empty)

	for i := 1; i <= ant.NumberOfAnts; i++ {
		a := strconv.Itoa(i)

		for j := 0; j < len(room.Dijkstra)-1; j++ {
			for k := 0; k < min; k++ {
				if len(room.Dijkstra) == 2 || len(room.Dijkstra[j]) == min {
					q.Enqueue("L" + a + "-" + (room.Dijkstra[j][k]))
				}
			}
		}

	}
	for _, i := range q.ant {
		fmt.Println(i + " ")
	}

	// // Code to measure
	// duration := time.Since(start)

	// // Formatted string, such as "2h3m0.5s" or "4.503μs"
	// fmt.Println(duration.Seconds())

	// // Nanoseconds as int64
	// fmt.Println(duration.Nanoseconds())
}
