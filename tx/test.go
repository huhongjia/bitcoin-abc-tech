package main

import "fmt"

func main() {

	sup := int64(50 * 1e8)
	fmt.Println(sup >> 64)

	//temp := "76a91419650240e343a3fba20abf37cc6dbfba9cdc0f1288ac"
	//
	//size := len(temp)
	//
	//splitRes := ""
	//list := make([]string, 0)
	//for i := 0; i < size; i += 2 {
	//
	//	str := temp[i : i+2]
	//	list = append(list, str)
	//	splitRes = splitRes + str + " "
	//}
	//
	//length := len(list)
	//for i := 0; i < length/2; i++ {
	//	list[i], list[length-i-1] = list[length-i-1 ], list[i]
	//}
	//
	//fmt.Println(splitRes)
	//revertRes := ""
	//for _, t := range list {
	//	revertRes = revertRes + t + " "
	//}
	//fmt.Println(revertRes)
	//fmt.Println(length)
}
