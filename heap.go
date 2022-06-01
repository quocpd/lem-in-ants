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
	NumberOfAnts int
	RoomName     string
	Weight       int
	StartRoom    string
	EndRoom      string
	Paths        []string
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

func (g *Graph) getPath(start, end string) (int, []string) {
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

			return p.Value, p.Rooms
		}

		for _, e := range g.getRooms(RoomName) {
			if !visited[e.RoomName] {
				// We calculate the total spent so far plus the cost and the Path of getting here
				h.push(Path{Value: p.Value + e.Weight, Rooms: append([]string{}, append(p.Rooms, e.RoomName)...)})
			}
		}

		visited[RoomName] = true
	}

	return 0, nil
}

func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

// Queue represents a queue that holds a slice
type Queue struct {
	items []string
}

// Enqueue adds a value at the end
func (q *Queue) Enqueue(i string) {
	q.items = append(q.items, i)
}

// Dequeue
func (q *Queue) Dequeue() string {
	toRemove := q.items[len(q.items)-1]
	q.items = q.items[1:]
	return toRemove
}

func main() {
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
	if strArr[0] <= "0" {
		fmt.Println("Error, ant colony has died! (Number of ants must be at least 1.)")
		os.Exit(0)
	} else if NumberOfAnts, err := strconv.Atoi(strArr[0]); err != nil {
		room.NumberOfAnts = NumberOfAnts
		fmt.Println("Number of ants must be a positive integer.")
		os.Exit(0)

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

	var chunk [][]string
	for i := 0; i < len(roomsWithCoordinates); i += 3 {

		end := i + 3
		if end > len(roomsWithCoordinates) {
			end = len(roomsWithCoordinates)
		}
		chunk = append(chunk, roomsWithCoordinates[i:end])
	}

	for i := 0; i < len(chunk); i++ {
		room.StartRoom = chunk[len(chunk)-2][0]
		room.EndRoom = chunk[len(chunk)-1][0]
	}

	Graph := newGraph()
	for i, str := range strArr {
		if strings.Contains(str, "-") {

			links := strings.Split(strArr[i], "-")

			if links[0] != links[1] {
				Graph.addRoom(links[0], links[1], 1)
				fmt.Println(links[0], links[1])
			} else {
				fmt.Println("Ants cannot process pheromones! (Room cannot link to itself e.g. 3-3)")
				os.Exit(0)
			}

		}
	}

	fmt.Println("Dijkstra")
	fmt.Println()
	// fmt.Println(Graph.getPath(room.StartRoom, room.EndRoom))

	for i := 0; i < len(chunk)-2; i++ {
		fmt.Println(Graph.getPath(chunk[i][0], room.EndRoom))

		// Graph.getPath(chunk[i][0], room.EndRoom)
	}
}
