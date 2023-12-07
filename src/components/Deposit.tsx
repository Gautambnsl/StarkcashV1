"use client";
//@ts-ignore
import React, { useState } from "react";
//@ts-ignore
import {encrypt} from "../components/middleware/encrypt.js"
//@ts-ignore
import {
	//@ts-ignore
	Provider,
	//@ts-ignore
	Account,
	//@ts-ignore
	Contract,
	//@ts-ignore
	json,
	//@ts-ignore
	hash,
	//@ts-ignore
	Calldata,
	//@ts-ignore
	num,
	//@ts-ignore
	RawCalldata,
	//@ts-ignore
	RawArgsObject,
	cairo,
	//@ts-ignore
	uint256,
	//@ts-ignore
	constants,
	//@ts-ignore
	RawArgsArray,
  } from "starknet";
  //@ts-ignore
import { hash_message,} from "./middleware/Interaction_script";
import { CallData } from "starknet";
//@ts-ignore

const erc20_address =
	"0x0636511044b76a1feb98d6ece83b623113f3cfa5d00e6356163823df58f0e228";
const core_address =
	"0x0155d74acd3b12a099344c52ae964c74dcb82d7c931d829e645fd4975db04282";


interface DepositProps {
	isConnect: boolean;
	connection : any;
	walletHandle: () => void;
}

interface UserInput {
	address: string;
	amount: number;
}

interface Error {
	status: boolean;
	address: string;
	amount: string;
}

export default function Deposit({ isConnect, walletHandle , connection}: DepositProps) {
	const [userInput, setUserInput] = useState<UserInput>({
		address: "",
		amount: 0,
	});

	const [error, setError] = useState<Error>({
		status: false,
		address: "",
		amount: "",
	});

	const handleChange = (key: string, value: string | number) => {
		setUserInput((prev) => ({ ...prev, [key]: value }));
	};

	const handleDeposit = async () => {

		console.log(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		// handle deposit;
		const { address, amount } = userInput;

		// handling errors
		if (address.length === 0) {
			let msg = "Enter valid address!!";
			setError({
				status: true,
				address: msg,
				amount: "",
			});
			return;
		}

		if (amount <= 0) {
			let msg = "Enter valid amount!!";
			setError({
				status: true,
				address: "",
				amount: msg,
			});
			return;
		}

		if (error.status) {
			// clearing error messages
			setError({
				status: false,
				address: "",
				amount: "",
			});
		}

		
		//userInput.address
		console.log(userInput.address)
		console.log(encrypt(userInput.address))
		let hash = encrypt(userInput.address)
		// hash = "0x" + hash
		// console.log(hash)


		const modifiedString = "0x" + hash.slice(4);
		console.log(modifiedString,"this is final arg")

		let success = await connection.execute([
			{
				contractAddress: erc20_address,
				entrypoint: 'approve',
				calldata: CallData.compile({
				recipient : core_address,
				amount: cairo.uint256(amount)	
				}),
			}, 
			{
				contractAddress: core_address,
				entrypoint: 'deposit',
				calldata: CallData.compile({
				amount: amount	,
				hash : modifiedString,
				
				}),
			}
		])

		// await approve(userInput.address, userInput.amount,connection );
		// await deposit(userInput.amount,userInput.address,connection );
		console.log(success,"success")
		alert(`Transaction Success, Tx Hash : ${success.transaction_hash} `)
	};

	return (
		<>
			<div className="input">
				<label htmlFor="address">Withdrawer's Address:</label>

				<input
					type="text"
					value={userInput.address}
					onChange={(e) => handleChange("address", e.target.value)}
					disabled={!isConnect}
					id="address"
					placeholder="address..."
				/>

				{error.status && error.address && <p>{error.address}</p>}
			</div>

			<div className="input">
    <label htmlFor="token">Token: </label>
    <select
        value="Enter the Token"
        id="token"
    >
        <option value="eth">wETH</option>
        {/* Add more options as needed */}
    </select>
</div>

			

			<div className="input">
				<label htmlFor="amount">Amount:</label>

				<input
					type="number"
					value={userInput.amount}
					onChange={(e) => handleChange("amount", +e.target.value)}
					disabled={!isConnect}
					min={0}
					id="amount"
					placeholder="amount..."
				/>
				{error.status && error.amount && <p>{error.amount}</p>}
			</div>

			{!isConnect ? (
				<button onClick={walletHandle}>Connect Wallet</button>
			) : (
				<button onClick={handleDeposit}>Deposit</button>
			)}
		</>
	);
}
