package main

import (
	"DataStruct/man"
	"DataStruct/stack"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

const (
	MAXVex   = 30
	INFINITY = 32768
)

type CloseEdge struct {
	adjVex  int
	lowcost int
}

type PrimeTreeNode struct {
	Vex1     string
	Vex2     string
	distance int
}

type PrimeTree struct {
	Node   PrimeTreeNode
	Lchild *PrimeTree
	Rchild *PrimeTree
}

type AdjMatrix struct {
	Arcs   [MAXVex][MAXVex]int `json:"Arcs"`   //邻接矩阵
	Vex    [MAXVex]string      `json:"Vex"`    //顶点
	Info   [MAXVex]string      `json:"Info"`   //描述信息
	Vexnum int                 `json:"Vexnum"` //顶点数量
	Arcnum int                 `json:"Arcnum"` //边的数量
}

var visit [MAXVex]int

func getVexNo(G *AdjMatrix, str string) int {
	for i := 1; i <= G.Vexnum; i++ {
		if G.Vex[i] == str {
			return i
		}
	}
	return -1
} //返回顶点的下标

func Create(G *AdjMatrix) {
	var distance, Vex1, Vex2, Arcnum int

	fmt.Println("请输入地点的个数:")
	fmt.Scanln(&G.Vexnum)

	for i := 1; i <= G.Vexnum; i++ {
		fmt.Printf("第%d个地点:", i)
		fmt.Scanln(&G.Vex[i])
		fmt.Println("请输入该地点的介绍：")
		fmt.Scanln(&G.Info[i])
	}
	for i := 1; i <= G.Vexnum; i++ {
		for j := 1; j <= G.Vexnum; j++ {
			if i == j {
				G.Arcs[i][j] = 0
			} else {
				G.Arcs[i][j] = INFINITY // 设置为一个足够大的数值，表示无穷大
			}
		}
	}

	fmt.Println("请输入所有路线,比如宿舍 食堂代表宿舍到食堂的路线，输入0 0结束输入：")
	for {
		var add1, add2 string

		fmt.Println("请输入路线：")
		fmt.Scanln(&add1, &add2)

		if add1 == "0" || add2 == "0" {
			break
		}

		Vex1 = getVexNo(G, add1)
		Vex2 = getVexNo(G, add2)

		if Vex1 == -1 || Vex2 == -1 {
			fmt.Println("此地点不存在！请重新输入！！！")
			continue
		}

		fmt.Printf("请输入%s到%s的距离：", add1, add2)
		fmt.Scanln(&distance)

		G.Arcs[Vex1][Vex2] = distance
		G.Arcs[Vex2][Vex1] = distance
		Arcnum++
	}
	G.Arcnum = Arcnum
	//fmt.Println(G)
	WriteFileAdjMatrix("AdjMatrix.txt", G)
}
func printTree(root *PrimeTree) {
	if root == nil {
		return
	}
	fmt.Println(root.Node.Vex1, "----------", root.Node.distance, "------------", root.Node.Vex2)
	printTree(root.Rchild)
	printTree(root.Lchild)
}
func prime(G *AdjMatrix, start int) *PrimeTree {
	closedge := make([]CloseEdge, MAXVex)
	// 初始化 closedge 数组
	for i := 1; i <= G.Vexnum; i++ {
		if i != start {
			closedge[i].adjVex = start
			closedge[i].lowcost = G.Arcs[start][i]
		}
	}
	// 根节点
	root := &PrimeTree{}
	var buildTree func(parent *PrimeTree, e int)
	buildTree = func(parent *PrimeTree, e int) {
		min := INFINITY
		m := 0 // 最近的顶点
		// 在 closedge 数组中找到最小权值的边
		for k := 1; k <= G.Vexnum; k++ {
			if closedge[k].lowcost != 0 && closedge[k].lowcost < min {
				m = k
				min = closedge[k].lowcost
			}
		}
		// 将找到的边添加到最小生成树中
		node := PrimeTreeNode{
			Vex1:     G.Vex[closedge[m].adjVex],
			Vex2:     G.Vex[m],
			distance: min,
		}
		current := &PrimeTree{Node: node}
		if parent != nil {
			if parent.Lchild == nil {
				parent.Lchild = current
			} else {
				parent.Rchild = current
			}
		} else {
			parent = current
		}
		closedge[m].lowcost = 0
		for i := 1; i <= G.Vexnum; i++ {
			if i != m && G.Arcs[m][i] < closedge[i].lowcost {
				closedge[i].lowcost = G.Arcs[m][i]
				closedge[i].adjVex = m
			}
		}
		if e < G.Vexnum {
			buildTree(current, e+1)
		}
	}

	// 从根节点开始构建最小生成树
	buildTree(root, 1)

	// 打印最小生成树的结果
	fmt.Println("最佳布网路线如下：")
	printTree(root)

	return root
}

func printPath(G *AdjMatrix) {
	clearCommand := exec.Command("clear")
	clearCommand.Stdout = os.Stdout
	clearCommand.Run()

	for i := 1; i <= G.Vexnum; i++ {
		for j := 1; j <= G.Vexnum; j++ {
			if G.Arcs[i][j] != INFINITY {
				fmt.Printf("%s-->%s 距离：%d     ", G.Vex[i], G.Vex[j], G.Arcs[i][j])
			}
		}
		fmt.Println()
	}

	fmt.Scanln()
	fmt.Scanln()
}

func WriteFileAdjMatrix(filePath string, G *AdjMatrix) error {
	data, err := json.MarshalIndent(G, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Adjacency Matrix written to %s\n", filePath)
	return nil
}

func ReadFileCreateAdjMatrix(filePath string, G *AdjMatrix) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, G)
	if err != nil {
		return err
	}

	fmt.Printf("Adjacency Matrix read from %s\n", filePath)
	return nil
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

	G.Arcs[num1][num2] = distance
	WriteFileAdjMatrix("AdjMatrix.txt", G)

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
		fmt.Println("\t\t                ┌───────────────────┐                ")
		fmt.Println("\t\t┌──────├── 欢迎使用西邮导航系统 ──┤──────┐")
		fmt.Println("\t\t│              └───────────────────┘              │")
		fmt.Println("\t\t│                                                                      │")
		fmt.Println("\t\t│                                                                      │")
		fmt.Println("\t\t│   ┌────────────┐        ┌────────────┐   │")
		fmt.Println("\t\t│   │ 1.创建新的路线图 │        │ 2.显示指定地点的信息 │   │")
		fmt.Println("\t\t│   └────────────┘        └────────────┘   │")
		fmt.Println("\t\t│                                                                      │")
		fmt.Println("\t\t│   ┌────────────┐        ┌────────────┐   │")
		fmt.Println("\t\t│   │ 3.显示指定两地的简单路径 │  │ 4.显示指定两地的最短路线 │   │")
		fmt.Println("\t\t│   └────────────┘        └────────────┘   │")
		fmt.Println("\t\t│                                                                      │")
		fmt.Println("\t\t│   ┌────────────┐        ┌────────────┐   │")
		fmt.Println("\t\t│   │ 5.增加新的地点和路线 │     │ 6.删除旧地点 │   │")
		fmt.Println("\t\t│   └────────────┘        └────────────┘   │")
		fmt.Println("\t\t│   ┌────────────┐        ┌────────────┐   │")
		fmt.Println("\t\t│   │ 7.删除指定路线 │  │ 8.显示所有的路线 │   │")
		fmt.Println("\t\t│   └────────────┘        └────────────┘   │")
		fmt.Println("\t\t│   ┌────────────┐        ┌────────────┐   │")
		fmt.Println("\t\t│   │ 9.最佳布网方案 │  │ 10. │ 显示平面图  │")
		fmt.Println("\t\t│   └────────────┘        └────────────┘   │")
		fmt.Println("\t\t│   ┌────────────┐        ┌────────────┐   │")
		fmt.Println("\t\t│   │ 11. │ 添加新路线或修改路线长度 │ 0. │ 退出 │")
		fmt.Println("\t\t│   └────────────┘        └────────────┘   │")
		fmt.Println("\t\t│                                                                      │")
		fmt.Println("\t\t│                                                                      │")
		fmt.Println("\t\t│                                                                      │")
		fmt.Println("\t\t└───────────────────────────────────┘\n\n")
		fmt.Print("\t\t请根据你的需求选择操作：")
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
	//实现显示指定地点信息的逻辑
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
		fmt.Printf("%s 简介: %s\n", G.Vex[num], G.Info[num])
	}

	fmt.Scanln()
	fmt.Scanln()
}
func gettopnextAdj(G *AdjMatrix, top int, st int) int {
	for i := st + 1; i <= G.Vexnum; i++ {
		if G.Arcs[top][i] != INFINITY {
			return i
		}
	}
	return 0
}
func displaySimplePath(G *AdjMatrix) {
	//实现显示指定两地简单路径的逻辑
	clearCommand := exec.Command("clear")
	clearCommand.Stdout = os.Stdout
	clearCommand.Run()

	var str1, str2 string
	var start, end, st, top, i int
	count := INFINITY
	mindist := make([]int, MAXVex)
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
				fmt.Print(G.Vex[S.Num[0]])
				for i = 1; i <= S.Top; i++ {
					fmt.Printf("-->%s", G.Vex[S.Num[i]])
				}
				fmt.Println()
			}
		}
	}

	fmt.Println("\n\n中转次数最少的路径如下：")
	for i = 0; i < count; i++ {
		if i == 0 {
			fmt.Print(G.Vex[mindist[i]])
		} else {
			fmt.Printf("-->%s", G.Vex[mindist[i]])
		}
	}
	fmt.Println("\n\n")

	fmt.Scanln()
	fmt.Scanln()
}
func Dijkstra(G *AdjMatrix, start int, dist []int, path [][]int) {
	var mindist, i, j, k, t int
	for i = 1; i <= G.Vexnum; i++ {
		dist[i] = G.Arcs[start][i]
		if G.Arcs[start][i] != INFINITY {
			path[i][1] = start
		}
	}

	for i = 2; i <= G.Vexnum; i++ {
		mindist = INFINITY
		for j = 1; j <= G.Vexnum; j++ {
			if path[j][0] == 0 && dist[j] < mindist {
				k = j
				mindist = dist[j]
			}
		}
		if mindist == INFINITY {
			return
		}
		path[k][0] = 1
		//for j = 1; j <= G.Vexnum; j++ {
		//	if path[j][0] == 0 && G.Arcs[k][j] < INFINITY && dist[k]+G.Arcs[k][j] < dist[j] {
		//		dist[j] = dist[k] + G.Arcs[k][j]
		//		t = 1
		//		for path[k][t] != 0 {
		//			path[j][t] = path[k][t]
		//			t++
		//		}
		//		path[j][t] = k
		//		path[j][t+1] = 0
		//	}
		//}
		for j = 1; j <= G.Vexnum; j++ {
			if path[j][0] == 0 && G.Arcs[k][j] < INFINITY && dist[k]+G.Arcs[k][j] < dist[j] {
				dist[j] = dist[k] + G.Arcs[k][j]
				t = 1
				for path[k][t] != 0 {
					path[j][t] = path[k][t]
					t++
				}
				path[j][t] = k
				path[j][t+1] = 0
			}
		}

		// 更新路径信息
		t = 1
		for path[k][t] != 0 {
			path[j][t] = path[k][t]
			t++
		}
		path[j][t] = k
		path[j][t+1] = 0
	}
}
func displayShortestPath(G *AdjMatrix) {
	//实现显示指定两地最短路径的逻辑
	clearCommand := exec.Command("clear")
	clearCommand.Stdout = os.Stdout
	clearCommand.Run()

	var str1, str2 string
	var start, end int
	dist := make([]int, MAXVex)
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

	path := make([][]int, MAXVex)
	for i := 0; i < MAXVex; i++ {
		path[i] = make([]int, MAXVex)
	}

	Dijkstra(G, start, dist, path)

	i := 1
	for path[end][i] != 0 {
		fmt.Printf("%s-->", G.Vex[path[end][i]])
		i++
	}
	fmt.Println(G.Vex[end])

	fmt.Scanln()
	fmt.Scanln()
}
func addNewAddressPath(G *AdjMatrix) {
	//实现增加新地点和路径的逻辑
	var str string
	var Vexnum, num1, i, distance int
	fmt.Print("请输入地点名：")
	fmt.Scan(&str)

	Vexnum = G.Vexnum
	G.Vexnum++
	fmt.Print("请输入地点的介绍：")
	fmt.Scan(&G.Info[G.Vexnum])

	Vexnum++
	G.Vex[Vexnum] = str
	for i = 1; i <= G.Vexnum; i++ {
		G.Arcs[Vexnum][i] = INFINITY
		G.Arcs[i][Vexnum] = INFINITY
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
			G.Arcs[Vexnum][num1] = distance
			G.Arcs[num1][Vexnum] = distance
		}
	}

	WriteFileAdjMatrix("AdjMatrix.txt", G)
}
func delOldAddress(G *AdjMatrix) {
	//实现删除旧地点的逻辑
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
	for i = num; i <= G.Vexnum-1; i++ {
		G.Vex[i] = G.Vex[i+1]
		G.Info[i] = G.Info[i+1]
	}

	for i = num; i <= G.Vexnum-1; i++ {
		for j = 1; j <= G.Vexnum; j++ {
			G.Arcs[i][j] = G.Arcs[i+1][j]
		}
	}

	for j = num; j <= G.Vexnum-1; j++ {
		for i = 1; i <= G.Vexnum-1; i++ {
			G.Arcs[i][j] = G.Arcs[i][j+1]
		}
	}

	G.Vexnum--
}
func delOldPath(G *AdjMatrix) {
	//实现删除指定路线的逻辑
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
	} else if G.Arcs[num1][num2] == INFINITY {
		fmt.Println("该路径不存在！任意键返回")
		fmt.Scanln()
		fmt.Scanln()
		return
	} else {
		G.Arcs[num1][num2] = INFINITY
	}

	return
}
func addInfo(G *AdjMatrix) {
	var i int
	for i = 1; i <= G.Vexnum; i++ {
		fmt.Printf("请输入%s的介绍：", G.Vex[i])
		fmt.Scan(&G.Info[i])
	}

	WriteFileAdjMatrix("AdjMatrix.txt", G)
}

func main() {
	G := new(AdjMatrix)
	ReadFileCreateAdjMatrix("AdjMatrix.txt", G)
	menu(G)
}
