package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

var (
	account_addr   string = "0x06f36e8a0fc06518125bbb1c63553e8a7d8597d437f9d56d891b8c7d3c977716"
	privateKey     string = "0x0687bf84896ee63f52d69e6de1b41492abeadc0dc3cb7bd351d0a52116915937"
	public_key     string = "0x58b0824ee8480133cad03533c8930eda6888b3c5170db2f6e4f51b519141963"
	contractMethod string = "0x02f0b3c5710379609eb5495f1ecd348cb28167711b73609fe565a72734550354"
	base           string = "https://starknet-goerli.g.alchemy.com/v2/bbDTAS9Qs68Po3ssBBPVm378Pye6oS2A"
)

func QueryIsendEvents(startblock uint64, endblock uint64) {
	// base := "https://starknet-goerli.g.alchemy.com/v2/bbDTAS9Qs68Po3ssBBPVm378Pye6oS2A"
	fmt.Println("Starting simpeCall example")
	c, err := ethrpc.DialContext(context.Background(), base)
	if err != nil {
		fmt.Println("Failed to connect to the client, did you specify the url in the .env.mainnet?")
		panic(err)
	}
	clientv02 := rpc.NewProvider(c)
	fmt.Println("Established connection with the client")

	contractAddress, _ := utils.HexToFelt("0x0155d74acd3b12a099344c52ae964c74dcb82d7c931d829e645fd4975db04282")

	DepositedEvent, _ := utils.HexToFelt("0x69105484e3b5f553164aa6de1f67321ea2757275a5e614365c90b9ed0a5e9b")

	for currentBlock := startblock; currentBlock <= endblock; currentBlock += 10 {

		eventInput := rpc.EventsInput{
			EventFilter: rpc.EventFilter{
				FromBlock: rpc.BlockID{
					Number: &startblock,
				},
				ToBlock: rpc.BlockID{
					Number: &endblock,
				},
				Address: contractAddress,
				Keys: [][]*felt.Felt{[]*felt.Felt{
					DepositedEvent,
				}},
			},
			ResultPageRequest: rpc.ResultPageRequest{
				ChunkSize: 1,
			},
		}

		events, err := clientv02.Events(context.Background(), eventInput)
		if err != nil {
			fmt.Println("Unsuccessful")
			panic(err)
		}

		for _, emittedEvent := range events.Events {
			dataCount := len(emittedEvent.Event.Data)
			fmt.Println(dataCount)
			var1 := emittedEvent.Event.Data[0]
			fmt.Println(var1, "Value 1")
			var2 := emittedEvent.Event.Data[1]

			fmt.Println("Variable 2 : ", var2)

			// decrypted_value := decryptData(var2.String(), 3)

			var calldata2 []*felt.Felt

			// size_val := "0x02"
			zero_val := "0x00"
			zero_felt, _ := utils.HexToFelt(zero_val)
			// size_felt, _ := utils.HexToFelt(size_val)

			addr_ := "0x03f045eb3AF74bfbAd1cD9767F13257dA110116a0D25B4B82f602759ebc8031b"

			value__1, _ := utils.HexToFelt(addr_)

			// calldata2 = append(calldata2, size_felt)
			calldata2 = append(calldata2, value__1)
			calldata2 = append(calldata2, var1)
			calldata2 = append(calldata2, zero_felt)

			// calldata2 = append(calldata2,var2)

			// hex_address_val := var2

			// hex1_values := []string{var1.String()}

			fmt.Println(calldata2, " : Payload	>>>>>>")

			// felt_data := ConvertHexStringsToFelt(hex1_values)

			// calldata := generic_append(hexArray, felt_data)

			invoke(calldata2)

		}
	}
}

func invoke(value1 []*felt.Felt) {
	c, err := ethrpc.DialContext(context.Background(), base)
	if err != nil {
		fmt.Println("Failed to connect to the client, did you specify the url in the .env.testnet?")
		panic(err)
	}

	clientv02 := rpc.NewProvider(c)

	account_address, err := utils.HexToFelt(account_addr)
	if err != nil {
		panic(err.Error())
	}
	// 	//Initializing the account memkeyStore
	ks := account.NewMemKeystore()
	fakePrivKeyBI, ok := new(big.Int).SetString(privateKey, 0)
	if !ok {
		panic(err.Error())
	}
	ks.Put(public_key, fakePrivKeyBI)

	fmt.Println("Established connection with the client")

	maxfee, _ := utils.HexToFelt("0x127156ff7224")

	// 	//Initializing the account
	accnt, err := account.NewAccount(clientv02, account_address, public_key, ks)
	if err != nil {
		panic(err.Error())
	}

	// 	//Getting the nonce from the account, and then converting it into felt
	nonce_string, _ := accnt.Nonce(context.Background(), rpc.BlockID{Tag: "latest"}, accnt.AccountAddress)
	nonce, _ := utils.HexToFelt(*nonce_string)

	// 	//Building the InvokeTx struct
	InvokeTx := rpc.InvokeTxnV1{
		MaxFee:        maxfee,
		Version:       rpc.TransactionV1,
		Nonce:         nonce,
		Type:          rpc.TransactionType_Invoke,
		SenderAddress: account_address,
	}
	contractAddress_1, _ := utils.HexToFelt("0x0636511044b76a1feb98d6ece83b623113f3cfa5d00e6356163823df58f0e228")

	ContractMethod, _ := utils.HexToFelt(contractMethod)

	fmt.Println(value1, " : Value 1 Array")

	// fmt.Println("Valueeee 2 ::: -> : ", value1)

	FnCall := rpc.FunctionCall{
		ContractAddress:    contractAddress_1,
		EntryPointSelector: ContractMethod,
		Calldata:           value1,
	}

	CairoContractVersion := 2

	InvokeTx.Calldata, err = accnt.FmtCalldata([]rpc.FunctionCall{FnCall}, CairoContractVersion)

	err = accnt.SignInvokeTransaction(context.Background(), &InvokeTx)

	resp, err := accnt.AddInvokeTransaction(context.Background(), InvokeTx)
	fmt.Printf("Response : ", resp)

}

func ConvertHexStringToU128Parts(hexString string) (lowHex, highHex string, err error) {

	bigIntValue := new(big.Int)

	// Convert the hexadecimal string to a big.Int
	_, success := bigIntValue.SetString(hexString, 0)
	if !success {
		return "", "", fmt.Errorf("Error converting the hex string to a big.Int")
	}

	u128 := new(big.Int)
	u128.Exp(big.NewInt(2), big.NewInt(128), nil)

	low := new(big.Int).And(bigIntValue, new(big.Int).Sub(u128, big.NewInt(1)))
	lowHex = fmt.Sprintf("0x%032x", low)

	high := new(big.Int).Rsh(bigIntValue, 128)
	highHex = fmt.Sprintf("0x%04x", high)

	return lowHex, highHex, nil
}

func ConvertHexStringsToFelt(hexStrings []string) []*felt.Felt {
	var feltArr []*felt.Felt

	for _, hexString := range hexStrings {
		lowHex, highHex, err := ConvertHexStringToU128Parts(hexString)
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}

		low_hex, _ := utils.HexToFelt(lowHex)
		high_Hex, _ := utils.HexToFelt(highHex)

		feltArr = append(feltArr, low_hex)
		feltArr = append(feltArr, high_Hex)
	}

	return feltArr
}

func generic_append(val []*felt.Felt, to_append []*felt.Felt) []*felt.Felt {
	// size := utils.Uint64ToFelt(uint64(len(val)))
	// to_append = append(to_append, size)
	to_append = append(to_append, val...)
	return to_append
}

func decryptData(data string, shift int) string {
	result := ""
	for _, char := range data {
		if char >= 'a' && char <= 'z' {
			result += string((char-'a'+26-rune(shift))%26 + 'a')
		} else if char >= 'A' && char <= 'Z' {
			result += string((char-'A'+26-rune(shift))%26 + 'A')
		} else {
			result += string(char)
		}
	}
	return result
}

func main() {
	var startblock, endblock uint64
	startblock = 913060
	endblock =
		913067

	QueryIsendEvents(startblock, endblock)
}

// package main

// import (
// 	"context"
// 	"fmt"
// 	"math/big"
// 	"os"

// 	"github.com/NethermindEth/juno/core/felt"
// 	"github.com/NethermindEth/starknet.go/account"
// 	"github.com/NethermindEth/starknet.go/rpc"
// 	"github.com/NethermindEth/starknet.go/utils"
// 	ethrpc "github.com/ethereum/go-ethereum/rpc"
// 	"github.com/joho/godotenv"
// )

// var (
// 	name         string = "testnet"
// 	account_addr string = "0x06f36e8a0fc06518125bbb1c63553e8a7d8597d437f9d56d891b8c7d3c977716"
// 	privateKey   string = "0x0687bf84896ee63f52d69e6de1b41492abeadc0dc3cb7bd351d0a52116915937"
// 	public_key   string = "0x58b0824ee8480133cad03533c8930eda6888b3c5170db2f6e4f51b519141963"
// 	someContract string = "0x5ff48139a2784ced5ac4a83c6a17913f7c97e407656401503c97e4419a34633" //main contract
// 	// someContract   string = "0x207d4810c08361e42dda527ae2565f1337d35624c5623df650662360e133bf3"
// 	contractMethod string = "0x0154f54d9fe2a16d94da0c42935457deda4e2df970f01c85bdaaf5b11fc341ff" //updateValset
// 	// contractMethod string = "0x004b0981bffc9f82c5b597bc66fac38f60eb3528424c806e4921bfb72720459b"
// )

// func main() {

// 	godotenv.Load(fmt.Sprintf(".env.%s", name))
// 	base := os.Getenv("INTEGRATION_BASE") //please modify the .env.testnet and replace the INTEGRATION_BASE with an starknet goerli RPC.
// 	fmt.Println("Starting simpleInvoke example")

// 	//Initialising the connection
// c, err := ethrpc.DialContext(context.Background(), base)
// if err != nil {
// 	fmt.Println("Failed to connect to the client, did you specify the url in the .env.testnet?")
// 	panic(err)
// }

// 	//Initialising the provider
// clientv02 := rpc.NewProvider(c)

// 	//Here we are converting the account address to felt
// 	account_address, err := utils.HexToFelt(account_addr)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// // 	//Initializing the account memkeyStore
// 	ks := account.NewMemKeystore()
// 	fakePrivKeyBI, ok := new(big.Int).SetString(privateKey, 0)
// 	if !ok {
// 		panic(err.Error())
// 	}
// 	ks.Put(public_key, fakePrivKeyBI)

// 	fmt.Println("Established connection with the client")

// 	maxfee, _ := utils.HexToFelt("0x127156ff7224")

// // 	//Initializing the account
// 	accnt, err := account.NewAccount(clientv02, account_address, public_key, ks)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// // 	//Getting the nonce from the account, and then converting it into felt
// 	nonce_string, _ := accnt.Nonce(context.Background(), rpc.BlockID{Tag: "latest"}, accnt.AccountAddress)
// 	nonce, _ := utils.HexToFelt(*nonce_string)

// // 	//Building the InvokeTx struct
// 	InvokeTx := rpc.InvokeTxnV1{
// 		MaxFee:        maxfee,
// 		Version:       rpc.TransactionV1,
// 		Nonce:         nonce,
// 		Type:          rpc.TransactionType_Invoke,
// 		SenderAddress: account_address,
// 	}

// 	contractAddress, _ := utils.HexToFelt(someContract)

// 	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// 	//Passing it as an array

// 	//Amount - u128
// 	//ContractAddress

// 	validator_array_new := []string{"0x4276bbFE317CED1A1B48084B6a91da31aBd7Ce78"}
// 	validator_felt_new, _ := utils.HexArrToFelt(validator_array_new)

// 	powers_array_new := []uint64{4294967295}
// 	powers_felt_new := Uint64ArrToFelt(powers_array_new)

// 	valsetnonce_single_new := big.NewInt(4)

// 	validator_array := []string{"0x3D3a12a5B194C30858D6295C68F37D309AfDfa5D"}
// 	validator_felt, _ := utils.HexArrToFelt(validator_array)

// 	powers_array := []uint64{4294967295}
// 	powers_felt := Uint64ArrToFelt(powers_array)

// 	valsetnonce_single := big.NewInt(1)
// 	// valsetnonce_felt := B(valsetnonce_array)
// 	// fmt.Println(valsetnonce_felt, "Valset !!")
// 	// fmt.Println(valsetnonce_felt, "valsetnonce_felt")

// 	r1_array := []string{"0xe82b315f409a002560e6e1b1e51c9ff60181523f74528543b9b2aa4fd3621c01"}
// 	r1_felt := ConvertHexStringsToFelt(r1_array)

// 	s1_array := []string{"0xe82b315f409a002560e6e1b1e51c9ff60181523f74528543b9b2aa4fd3621c01"}
// 	s1_felt := ConvertHexStringsToFelt(s1_array)

// 	v1_string := []string{"27"}
// 	v1_felt, _ := utils.HexArrToFelt(v1_string)

// 	var callData1 []*felt.Felt
// 	//validator
// 	callData1 = generic_append(validator_felt_new, callData1)

// 	//powers
// 	callData1 = generic_append(powers_felt_new, callData1)

// 	//valset nonce
// 	callData1 = BigIntToFeltParts(valsetnonce_single_new, callData1)
// 	//validator
// 	callData1 = generic_append(validator_felt, callData1)

// 	//powers
// 	callData1 = generic_append(powers_felt, callData1)

// 	//valset nonce
// 	callData1 = BigIntToFeltParts(valsetnonce_single, callData1)

// 	//r1
// 	callData1 = generic_append_2(r1_felt, callData1)

// 	//s1
// 	callData1 = generic_append_2(s1_felt, callData1)

// 	// //v1
// 	callData1 = generic_append(v1_felt, callData1)

// FnCall := rpc.FunctionCall{
// 	ContractAddress:    contractAddress,
// 	EntryPointSelector: ContractMethod,
// 	Calldata:           callData1,
// }

// CairoContractVersion := 2

// InvokeTx.Calldata, err = accnt.FmtCalldata([]rpc.FunctionCall{FnCall}, CairoContractVersion)

// err = accnt.SignInvokeTransaction(context.Background(), &InvokeTx)

// resp, err := accnt.AddInvokeTransaction(context.Background(), InvokeTx)
// fmt.Printf("Response : ", resp)

// }

// func Uint64ArrToFelt(num []uint64) []*felt.Felt {
// 	var feltArr []*felt.Felt
// 	for _, hex := range num {
// 		felt := utils.Uint64ToFelt(hex)

// 		feltArr = append(feltArr, felt)
// 	}
// 	return feltArr
// }

// func generic_append(val []*felt.Felt, to_append []*felt.Felt) []*felt.Felt {
// 	size := utils.Uint64ToFelt(uint64(len(val)))
// 	to_append = append(to_append, size)
// 	to_append = append(to_append, val...)
// 	return to_append
// }
// func generic_append_2(val []*felt.Felt, to_append []*felt.Felt) []*felt.Felt {
// 	size := utils.Uint64ToFelt(uint64(len(val) / 2))
// 	to_append = append(to_append, size)
// 	to_append = append(to_append, val...)
// 	return to_append
// }

// func BigIntToHexU128Parts(num *big.Int) (lowHex, highHex string) {
// 	// Create a big.Int with a value of 2^128
// 	u128 := new(big.Int)
// 	u128.Exp(big.NewInt(2), big.NewInt(128), nil)

// 	// Calculate the low part
// 	low := new(big.Int).And(num, new(big.Int).Sub(u128, big.NewInt(1)))
// 	lowHex = fmt.Sprintf("0x%032x", low)

// 	// Calculate the high part
// 	high := new(big.Int).Rsh(num, 128)
// 	highHex = fmt.Sprintf("0x%04x", high)

// 	return lowHex, highHex

// }
// func BigIntArrToHexU128Parts(val []*big.Int) []*felt.Felt {
// 	var feltArr []*felt.Felt

// 	for _, hex := range val {
// 		lowHex, highHex := BigIntToHexU128Parts(hex)
// 		low_hex, _ := utils.HexToFelt(lowHex)
// 		high_hex, _ := utils.HexToFelt(highHex)
// 		feltArr = append(feltArr, low_hex)
// 		feltArr = append(feltArr, high_hex)
// 	}

// 	return feltArr
// }

// func BigIntToFeltParts(val *big.Int, arrayToAppend []*felt.Felt) []*felt.Felt {
// 	lowHex, highHex := BigIntToHexU128Parts(val)
// 	low_hex, _ := utils.HexToFelt(lowHex)
// 	high_hex, _ := utils.HexToFelt(highHex)
// 	arrayToAppend = append(arrayToAppend, low_hex)
// 	arrayToAppend = append(arrayToAppend, high_hex)
// 	return arrayToAppend
// }

// func ConvertHexStringsToFelt(hexStrings []string) []*felt.Felt {
// 	var feltArr []*felt.Felt

// 	for _, hexString := range hexStrings {
// 		lowHex, highHex, err := ConvertHexStringToU128Parts(hexString)
// 		if err != nil {
// 			fmt.Println("Error:", err)
// 			return nil
// 		}

// 		low_hex, _ := utils.HexToFelt(lowHex)
// 		high_Hex, _ := utils.HexToFelt(highHex)

// 		feltArr = append(feltArr, low_hex)
// 		feltArr = append(feltArr, high_Hex)
// 	}

// 	return feltArr
// }

// func ConvertHexStringToU128Parts(hexString string) (lowHex, highHex string, err error) {

// 	bigIntValue := new(big.Int)

// 	// Convert the hexadecimal string to a big.Int
// 	_, success := bigIntValue.SetString(hexString, 0)
// 	if !success {
// 		return "", "", fmt.Errorf("Error converting the hex string to a big.Int")
// 	}

// 	u128 := new(big.Int)
// 	u128.Exp(big.NewInt(2), big.NewInt(128), nil)

// 	low := new(big.Int).And(bigIntValue, new(big.Int).Sub(u128, big.NewInt(1)))
// 	lowHex = fmt.Sprintf("0x%032x", low)

// 	high := new(big.Int).Rsh(bigIntValue, 128)
// 	highHex = fmt.Sprintf("0x%04x", high)

// 	return lowHex, highHex, nil
// }

// func BigIntArrToFelt(val []*big.Int) []*felt.Felt {
// 	var feltArr []*felt.Felt
// 	for _, hex := range val {
// 		felt := utils.BigIntToFelt(hex)
// 		feltArr = append(feltArr, felt)
// 	}
// 	return feltArr
// }

// func main() {
// 	ExecuteIncreaseValue()
// }
