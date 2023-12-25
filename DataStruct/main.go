package main

import (
	"DataStruct/man"
	"DataStruct/stack"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
)

const (
	MAXVEX   = 30
	INFINITY = 32768
)

type AdjMatrix struct {
	arcs   [MAXVEX][MAXVEX]int
	vex    [MAXVEX]string
	info   [MAXVEX]string
	vexnum int
	arcnum int
}

var visit [MAXVEX]int

//func WriteFileAdjMatrix(G *AdjMatrix) {
//	// 实现写入邻接矩阵到文件的逻辑
//	fmt.Println("Writing Adjacency Matrix to file...")
//}
//
//func delOldAddress(G *AdjMatrix) {
//	// 实现删除旧地点的逻辑
//	fmt.Println("Deleting old addresses...")
//}
//
//func delOldPath(G *AdjMatrix) {
//	// 实现删除指定路线的逻辑
//	fmt.Println("Deleting specified paths...")
//}
//
//func addNewAddressPath(G *AdjMatrix) {
//	// 实现增加新地点和路径的逻辑
//	fmt.Println("Adding new addresses and paths...")
//}
//
//func displayAddressInfo(G *AdjMatrix) {
//	// 实现显示指定地点信息的逻辑
//	fmt.Println("Displaying address information...")
//}
//
//func displaySimplePath(G *AdjMatrix) {
//	// 实现显示指定两地简单路径的逻辑
//	fmt.Println("Displaying simple path...")
//}
//
//func displayShortestPath(G *AdjMatrix) {
//	// 实现显示指定两地最短路径的逻辑
//	fmt.Println("Displaying shortest path...")
//}

func getVexNo(G *AdjMatrix, str string) int {
	for i := 1; i <= G.vexnum; i++ {
		if G.vex[i] == str {
			return i
		}
	}
	return -1
}

func Create(G *AdjMatrix) {
	var distance, vex1, vex2, arcnum int

	fmt.Println("请输入地点的个数:")
	fmt.Scanln(&G.vexnum)

	for i := 1; i <= G.vexnum; i++ {
		fmt.Printf("第%d个地点:", i)
		fmt.Scanln(&G.vex[i])
		fmt.Println("请输入该地点的介绍：")
		fmt.Scanln(&G.info[i])
	}

	fmt.Println("请输入所有路线,比如宿舍 食堂代表宿舍到食堂的路线，输入0 0结束输入：")
	for {
		var add1, add2 string

		fmt.Println("请输入路线：")
		fmt.Scanln(&add1, &add2)

		if add1 == "0" || add2 == "0" {
			break
		}

		vex1 = getVexNo(G, add1)
		vex2 = getVexNo(G, add2)

		if vex1 == -1 || vex2 == -1 {
			fmt.Println("此地点不存在！请重新输入！！！")
			continue
		}

		fmt.Printf("请输入%s到%s的距离：", add1, add2)
		fmt.Scanln(&distance)

		G.arcs[vex1][vex2] = distance
		G.arcs[vex2][vex1] = distance
		arcnum++
	}

	G.arcnum = arcnum
	WriteFileAdjMatrix(G)
}
func printAdjMatrix(G *AdjMatrix) {
	fmt.Println("\nAdjMatrix:")
	fmt.Print("     ")
	for i := 1; i <= G.vexnum; i++ {
		fmt.Printf("%5s", G.vex[i])
	}
	fmt.Println()
	for i := 1; i <= G.vexnum; i++ {
		fmt.Printf("%-15s", G.vex[i])
		for j := 1; j <= G.vexnum; j++ {
			if G.arcs[i][j] == INFINITY {
				fmt.Print(" ----")
			} else {
				fmt.Printf("%5d", G.arcs[i][j])
			}
		}
		fmt.Println()
	}
	fmt.Scanln()
}

func prime(G *AdjMatrix, start int) {
	type CloseEdge struct {
		adjvex  int
		lowcost int
	}

	type PrimeTree struct {
		vex1     string
		vex2     string
		distance int
	}
	closedge := make([]CloseEdge, MAXVEX)
	primeTree := make([]PrimeTree, MAXVEX)

	for i := 1; i <= G.vexnum; i++ {
		if i != start {
			closedge[i].adjvex = start
			closedge[i].lowcost = G.arcs[start][i]
		}
	}

	for e := 1; e <= G.vexnum; e++ {
		min := INFINITY
		m := 0

		for k := 1; k <= G.vexnum; k++ {
			if closedge[k].lowcost != 0 && closedge[k].lowcost < min {
				m = k
				min = closedge[k].lowcost
			}
		}

		primeTree[e].vex1 = G.vex[closedge[m].adjvex]
		primeTree[e].vex2 = G.vex[m]
		primeTree[e].distance = min

		closedge[m].lowcost = 0

		for i := 1; i <= G.vexnum; i++ {
			if i != m && G.arcs[m][i] < closedge[i].lowcost {
				closedge[i].lowcost = G.arcs[m][i]
				closedge[i].adjvex = m
			}
		}
	}

	fmt.Println("最佳布网路线如下：")
	for i := 1; i <= G.vexnum-1; i++ {
		fmt.Printf("%s-->%s 距离：%d\n", primeTree[i].vex1, primeTree[i].vex2, primeTree[i].distance)
	}
}
func printPath(G *AdjMatrix) {
	clearCommand := exec.Command("clear")
	clearCommand.Stdout = os.Stdout
	clearCommand.Run()

	for i := 1; i <= G.vexnum; i++ {
		for j := 1; j <= G.vexnum; j++ {
			if G.arcs[i][j] != INFINITY {
				fmt.Printf("%s-->%s 距离：%d     ", G.vex[i], G.vex[j], G.arcs[i][j])
			}
		}
		fmt.Println()
	}

	fmt.Scanln()
	fmt.Scanln()
}

func ReadFileCreateAdjMatrix(G *AdjMatrix) {
	file, err := os.Open("AdjMatrix.txt")
	if err != nil {
		fmt.Println("文件访问出错")
		os.Exit(0)
	}
	defer file.Close()

	err = binary.Read(file, binary.BigEndian, G)
	if err != nil {
		fmt.Println("读取文件出错")
		os.Exit(0)
	}
}

func WriteFileAdjMatrix(G *AdjMatrix) {
	file, err := os.OpenFile("AdjMatrix.txt", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("文件访问出错，按任意键返回......")
		fmt.Scanln()
		fmt.Scanln()
		return
	}
	defer file.Close()

	err = binary.Write(file, binary.BigEndian, G)
	if err != nil {
		fmt.Println("写入文件出错")
		fmt.Scanln()
		fmt.Scanln()
		return
	}
}
func bestNetworkPath(G *AdjMatrix) {
	clearCommand := exec.Command("clear")
	clearCommand.Stdout = os.Stdout
	clearCommand.Run()

	var str string
	fmt.Print("请输入布网起点：")
	fmt.Scan(&str)

	num := getVexNo(G, str)
	if num == -1 {
		fmt.Println("该地点不存在！按任意键返回......")
		fmt.Scanln()
		fmt.Scanln()
		return
	} else {
		prime(G, num)
	}

	fmt.Scanln()
	fmt.Scanln()
}
func onlyAddNewPath(G *AdjMatrix) {
	clearCommand := exec.Command("clear")
	clearCommand.Stdout = os.Stdout
	clearCommand.Run()

	var distance, num1, num2 int
	var str1, str2 string

	fmt.Print("请输入路径的起点和终点：")
	fmt.Scan(&str1, &str2)

	num1 = getVexNo(G, str1)
	num2 = getVexNo(G, str2)
	if num1 == -1 || num2 == -1 {
		fmt.Println("该地点不存在！按任意键返回......")
		fmt.Scanln()
		fmt.Scanln()
		return
	}

	fmt.Printf("请输入%s到%s的距离：", str1, str2)
	fmt.Scan(&distance)

	G.arcs[num1][num2] = distance
	WriteFileAdjMatrix(G)

	fmt.Println("添加/修改成功，按任意键返回......")
	fmt.Scanln()
	fmt.Scanln()
}

func menu(G *AdjMatrix) {
	var selectOption int
	for {
		clearCommand := exec.Command("clear")
		clearCommand.Stdout = os.Stdout
		clearCommand.Run()

		fmt.Println("\n\n\t\t\t\t校园导航图")
		fmt.Println("\t\t\t*************************")
		fmt.Println("\t\t\t1.创建新的路线图")
		fmt.Println("\t\t\t2.显示指定地点的信息")
		fmt.Println("\t\t\t3.显示指定两地的简单路径")
		fmt.Println("\t\t\t4.显示指定两地的最短路线")
		fmt.Println("\t\t\t5.增加新的地点和路线")
		fmt.Println("\t\t\t6.删除旧地点")
		fmt.Println("\t\t\t7.删除指定路线")
		fmt.Println("\t\t\t8.显示所有的路线")
		fmt.Println("\t\t\t9.最佳布网方案")
		fmt.Println("\t\t\t10.显示平面图")
		fmt.Println("\t\t\t11.添加新路线或修改路线长度")
		fmt.Println("\t\t\t0.退出")
		fmt.Println("\t\t\t*************************")
		fmt.Print("\t\t\t  请选择：")
		fmt.Scan(&selectOption)

		switch selectOption {
		case 1:
			Create(G)
		case 2:
			displayAddressInfo(G)
		case 3:
			displaySimplePath(G)
		case 4:
			displayShortestPath(G)
		case 5:
			addNewAddressPath(G)
		case 6:
			delOldAddress(G)
		case 7:
			delOldPath(G)
		case 8:
			printPath(G)
		case 9:
			bestNetworkPath(G)
		case 10:
			man.Showman()
		case 11:
			onlyAddNewPath(G)
		case 0:
			os.Exit(0)
		default:
			continue
		}
	}
}
func displayAddressInfo(G *AdjMatrix) {
	clearCommand := exec.Command("clear")
	clearCommand.Stdout = os.Stdout
	clearCommand.Run()

	var str string

	fmt.Print("请输入要显示信息的地点名：")
	fmt.Scan(&str)

	num := getVexNo(G, str)
	if num == -1 {
		fmt.Println("该地点不存在！")
	} else {
		fmt.Printf("%s 简介: %s\n", G.vex[num], G.info[num])
	}

	fmt.Scanln()
	fmt.Scanln()
}
func gettopnextAdj(G *AdjMatrix, top int, st int) int {
	for i := st + 1; i <= G.vexnum; i++ {
		if G.arcs[top][i] != INFINITY {
			return i
		}
	}
	return 0
}
func displaySimplePath(G *AdjMatrix) {
	clearCommand := exec.Command("clear")
	clearCommand.Stdout = os.Stdout
	clearCommand.Run()

	var str1, str2 string
	var start, end, st, top, i int
	count := INFINITY
	mindist := make([]int, MAXVEX)
	S := new(stack.Stack)
	S.Top = -1

	fmt.Print("请输入起点和终点，中间用空格隔开:")
	fmt.Scan(&str1, &str2)

	start = getVexNo(G, str1)
	end = getVexNo(G, str2)
	if start == -1 || end == -1 {
		fmt.Println("该地点不存在，任意键返回……")
		fmt.Scanln()
		fmt.Scanln()
		return
	}

	S.InSert(start)
	visit[start] = 1
	fmt.Printf("\n\n%s到%s的简单路径如下：\n", str1, str2)

	for !S.IsEmpty() {
		S.Gettop(&top)
		st = gettopnextAdj(G, top, st)

		if st == 0 {
			S.Out(&st)
			visit[st] = 0
		} else if st != 0 && visit[st] == 0 {
			S.InSert(st)
			visit[st] = 1

			if st == end {
				if count > S.Top+1 {
					count = S.Top + 1
					for i = 0; i < count; i++ {
						mindist[i] = S.Num[i]
					}
				}
				fmt.Print(G.vex[S.Num[0]])
				for i = 1; i <= S.Top; i++ {
					fmt.Printf("-->%s", G.vex[S.Num[i]])
				}
				fmt.Println()
			}
		}
	}

	fmt.Println("\n\n中转次数最少的路径如下：")
	for i = 0; i < count; i++ {
		if i == 0 {
			fmt.Print(G.vex[mindist[i]])
		} else {
			fmt.Printf("-->%s", G.vex[mindist[i]])
		}
	}
	fmt.Println("\n\n")

	fmt.Scanln()
	fmt.Scanln()
}
func Dijkstra(G *AdjMatrix, start int, dist []int, path [][]int) {
	var mindist, i, j, k, t int
	for i = 1; i <= G.vexnum; i++ {
		dist[i] = G.arcs[start][i]
		if G.arcs[start][i] != INFINITY {
			path[i][1] = start
		}
	}

	for i = 2; i <= G.vexnum; i++ {
		mindist = INFINITY
		for j = 1; j <= G.vexnum; j++ {
			if path[j][0] == 0 && dist[j] < mindist {
				k = j
				mindist = dist[j]
			}
		}
		if mindist == INFINITY {
			return
		}
		path[k][0] = 1
		for j = 1; j <= G.vexnum; j++ {
			if path[j][0] == 0 && G.arcs[k][j] < INFINITY && dist[k]+G.arcs[k][j] < dist[j] {
				dist[j] = dist[k] + G.arcs[k][j]
				t = 1
				for path[k][t] != 0 {
					path[j][t] = path[k][t]
					t++
				}
				path[j][t] = k
				path[j][t+1] = 0
			}
		}
	}
}
func displayShortestPath(G *AdjMatrix) {
	clearCommand := exec.Command("clear")
	clearCommand.Stdout = os.Stdout
	clearCommand.Run()

	var str1, str2 string
	var start, end int
	dist := make([]int, MAXVEX)
	fmt.Print("请输入起点和终点，中间用空格隔开：")
	fmt.Scan(&str1, &str2)

	start = getVexNo(G, str1)
	end = getVexNo(G, str2)
	if start == -1 || end == -1 {
		fmt.Println("该地点不存在!任意键返回……")
		fmt.Scanln()
		fmt.Scanln()
		return
	}

	path := make([][]int, MAXVEX)
	for i := 0; i < MAXVEX; i++ {
		path[i] = make([]int, MAXVEX)
	}

	Dijkstra(G, start, dist, path)

	i := 1
	for path[end][i] != 0 {
		fmt.Printf("%s-->", G.vex[path[end][i]])
		i++
	}
	fmt.Println(G.vex[end])

	fmt.Scanln()
	fmt.Scanln()
}
func addNewAddressPath(G *AdjMatrix) {
	var str string
	var vexnum, num1, i, distance int
	fmt.Print("请输入地点名：")
	fmt.Scan(&str)

	vexnum = G.vexnum
	G.vexnum++
	fmt.Print("请输入地点的介绍：")
	fmt.Scan(&G.info[G.vexnum])

	vexnum++
	G.vex[vexnum] = str
	for i = 1; i <= G.vexnum; i++ {
		G.arcs[vexnum][i] = INFINITY
		G.arcs[i][vexnum] = INFINITY
	}

	fmt.Println("请输入与该地点有路径的地点和它们之间的距离（中间用空格隔开，0 0 代表结束输入）：")
	for {
		fmt.Scan(&str, &distance)
		if str == "0" {
			break
		}
		num1 = getVexNo(G, str)
		if num1 == -1 {
			fmt.Println("该地点不存在！按任意键返回……")
			fmt.Scanln()
			fmt.Scanln()
			return
		} else {
			G.arcs[vexnum][num1] = distance
			G.arcs[num1][vexnum] = distance
		}
	}

	WriteFileAdjMatrix(G)
}
func delOldAddress(G *AdjMatrix) {
	var str string
	fmt.Print("输入要删除的地点：")
	fmt.Scan(&str)

	num := getVexNo(G, str)
	if num == -1 {
		fmt.Println("该地点不存在！！任意键返回……")
		fmt.Scanln()
		fmt.Scanln()
		return
	}

	var i, j int
	for i = num; i <= G.vexnum-1; i++ {
		G.vex[i] = G.vex[i+1]
		G.info[i] = G.info[i+1]
	}

	for i = num; i <= G.vexnum-1; i++ {
		for j = 1; j <= G.vexnum; j++ {
			G.arcs[i][j] = G.arcs[i+1][j]
		}
	}

	for j = num; j <= G.vexnum-1; j++ {
		for i = 1; i <= G.vexnum-1; i++ {
			G.arcs[i][j] = G.arcs[i][j+1]
		}
	}

	G.vexnum--
}
func delOldPath(G *AdjMatrix) {
	clearCommand := exec.Command("clear")
	clearCommand.Stdout = os.Stdout
	clearCommand.Run()

	var str1, str2 string
	var num1, num2 int
	fmt.Print("请输入路径的起点和终点用空格隔开：")
	fmt.Scan(&str1, &str2)

	num1 = getVexNo(G, str1)
	num2 = getVexNo(G, str2)
	if num1 == -1 {
		fmt.Println("起点不存在！任意键返回……")
		fmt.Scanln()
		fmt.Scanln()
		return
	} else if num2 == -1 {
		fmt.Println("终点不存在！任意键返回……")
		fmt.Scanln()
		fmt.Scanln()
		return
	} else if G.arcs[num1][num2] == INFINITY {
		fmt.Println("该路径不存在！任意键返回")
		fmt.Scanln()
		fmt.Scanln()
		return
	} else {
		G.arcs[num1][num2] = INFINITY
	}

	return
}
func addInfo(G *AdjMatrix) {
	var i int
	for i = 1; i <= G.vexnum; i++ {
		fmt.Printf("请输入%s的介绍：", G.vex[i])
		fmt.Scan(&G.info[i])
	}

	WriteFileAdjMatrix(G)
}

func main() {
	G := new(AdjMatrix)
	ReadFileCreateAdjMatrix(G)
	menu(G)
}
