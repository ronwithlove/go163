package main

func main(){

	//bc:=BLC.CreateBlockCHainWithGenesisBlock()
	//defer bc.DB.Close()
	////上链
	//bc.AddBlock([]byte("ron sent 10 tc to aaron"))
	//bc.AddBlock([]byte("jacky sent 10 tc to aaron"))
	//
	//bc.PrintChain()

	cli:=BLC.CLI{}
	cli.Run()

	//result:=BLC.Base58Encode([]byte("this is the example"))
	//fmt.Printf("result:%s\n",result)
	//
	//decodeResult:=BLC.Base58Decode([]byte("NK2smnfSzALMcNJ8YHsxUJMrfN"))
	//fmt.Printf("decode result:%s\n",decodeResult)

}