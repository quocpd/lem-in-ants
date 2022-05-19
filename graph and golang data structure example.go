package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Test struct {
	NumberOfAnts int
	TheRooms     [][]string
	TheLinks     [][]string
	StartRoom    string
	EndRoom      string
}

func ByteToString(fromByte []byte) string {
	toString := string(fromByte)
	return toString
}

func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func main() {
	var strArr []string
	file, err := os.Open("example02.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanWords)
	for fileScanner.Scan() {
		strArr = append(strArr, fileScanner.Text())
	}

	NumberOfAnts, _ := strconv.Atoi(strArr[0])

	test := &Test{
		NumberOfAnts: NumberOfAnts,
	}
	// fmt.Println(test.NumberOfAnts)

	// fmt.Println(len(strArr))

	RemoveIndex(strArr, 0)
	strArr = strArr[:len(strArr)-1]

	for i := 0; i < len(strArr); i++ {

		if strArr[i] == "##start" {
			test.StartRoom = strArr[i+1]

			RemoveIndex(strArr, i)
			strArr = strArr[:len(strArr)-1]

		}
		if strArr[i] == "##end" {
			test.EndRoom = strArr[i+1]
			RemoveIndex(strArr, i)
			strArr = strArr[:len(strArr)-1]
		}

	}
	var chunk [][]string
	for i := 0; i < len(strArr); i += 3 {
		end := i + 3
		if end > len(strArr) {
			end = len(strArr)
		}
		chunk = append(chunk, strArr[i:end])
	}

	// for _, i:= range strArr{
	// 	fmt.Println(strings.Index("IndeX:" ,i ))
	// }
	testGraph := &Graph{}
	for i := 0; i < len(chunk); i++ {
		if len(chunk[i][0]) == 1 {
			testGraph.AddRoom(chunk[i][0])
		}
	} // room

	for _, i := range strArr {
		if strings.Contains(i, "-") {
			link := strings.ReplaceAll(string(i), "-", "")

			links := strings.Split(link, "")
			testGraph.AddTunnel(links[0], links[1])
		}
	}

	testGraph.Print()

	// for _, i := range strArr {
	// 	fmt.Println(strings.Index("IndeX:", i))
	// }
}

// Graph represents an adjacency list graph
type Graph struct {
	Rooms []*Room
}

// Room represents a graphh Room
type Room struct {
	RoomName string // RoomName
	Links    []*Room
}

// Add Room adds a Room to the Graph
func (g *Graph) AddRoom(k string) {
	if contains(g.Rooms, k) {
		err := fmt.Errorf("Room %v not added because it is an existing Room Name", k)
		fmt.Println(err.Error())
	} else {
		g.Rooms = append(g.Rooms, &Room{RoomName: k})
	}
}

// Add Tunnel adds an Tunnel to the graph
func (g *Graph) AddTunnel(from, to string) {
	// get Room
	fromRoom := g.getRoom(from)
	toRoom := g.getRoom(to)
	// check error
	if fromRoom == nil || toRoom == nil {
		err := fmt.Errorf("Invalid Tunnel (%v-->%v)", from, to)
		fmt.Println(err.Error())
	} else if contains(fromRoom.Links, to) {
		err := fmt.Errorf("Existing Tunnel (%v-->%v)", from, to)
		fmt.Println(err.Error())
	} else {
		// add Tunnel
		fromRoom.Links = append(fromRoom.Links, toRoom)
	}
}

// getRoom
func (g *Graph) getRoom(k string) *Room {
	for i, v := range g.Rooms {
		if v.RoomName == k {
			return g.Rooms[i]
		}
	}
	return nil
}

// contains
func contains(s []*Room, k string) bool {
	for _, v := range s {
		if k == v.RoomName {
			return true
		}
	}
	return false
}

// Print will print the Links list for each cRoom of the graph
func (g *Graph) Print() {
	for _, v := range g.Rooms {
		fmt.Printf("\nRoom %v : ", v.RoomName)
		for _, v := range v.Links {
			fmt.Printf("%v", v.RoomName)
		}
	}
	fmt.Println()
}

// txt := strings.Fields(string(content))
// fmt.Println(txt)
// NumberOfAnts := txt[0]
// StartRoom := txt[2]
// fmt.Println(NumberOfAnts)
// fmt.Println(StartRoom)

// for _,i:= range txt{

// 	if i == "##end"{
// 		EndRoom := len(i)+1
// 		fmt.Println(EndRoom)
// 	}

// }

// txt := strings.Fields(string(text))

// fmt.Print(txt[0])
// for i := 0; i < len(strArr); i++ {
// 	if strArr[0] == "##start" {
// 		fmt.Println("test")
// 		RemoveIndex(strArr, i)

// 	}

// 	if strArr[i] == "##end" {
// 		RemoveIndex(strArr, i)
// 	}

// }

// testGraph := &Graph{}
// for _, i := range chunk {
// 	fmt.Println(chunk)
// 	if len(chunk[0][0]) == 1 {

// 		numOfRooms, _ := strconv.Atoi(i[0])

// 		for j := 0; j < numOfRooms; j++ {
// 			testGraph.AddRoom(i[0])
// 		}
// 	}
// }

// testGraph :=&Graph{}
// for i:= 0; i<0; i++{
// 	if len(chunk[i][i]) == 1{
// 		testGraph.AddRoom(chunk[i][i])
// 	}
// }

//
//
// fmt.Println(testGraph.Rooms)

// test := append(array, txt)
// fmt.Println(test)
// for i,_ := range txt{
// 	array = append(array, txt)
// 	fmt.Println(i)
//
// fmt.Println(array)

// fmt.Println(a )
// for i := 0; i < 5; i++ {
// 	test.AddRoom(strconv.Itoa(i))

//
// temp1 := "1 23 3"
// Temp2 := "2 16 7"
// var a [][]string
// b := strings.Split(temp1, " ")

// c := strings.Split(Temp2, " ")
// fmt.Println(b)
// fmt.Println(c)
// a = append(a, b)
// a = append(a, c)
// fmt.Println(a[0][1])
