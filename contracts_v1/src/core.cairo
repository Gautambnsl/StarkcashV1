use starknet::ContractAddress;

#[starknet::interface]
trait ICoreContract<TContractState> {
    fn deposit(ref self: TContractState, amount: u128, value: felt252);
    fn initialize(ref self: TContractState, token_address: ContractAddress);
}

#[starknet::interface]
trait IERC20<TContractState> {
    fn transfer_from(
        ref self: TContractState, sender: ContractAddress, recipient: ContractAddress, amount: u256
    );
}


#[starknet::contract]
mod core_contract {
    use starknet::ContractAddress;
    use super::{
        ICoreContractDispatcher, ICoreContractDispatcherTrait, ICoreContract, IERC20,
        IERC20Dispatcher, IERC20DispatcherTrait
    };
    use starknet::get_caller_address;
    use integer::u128_to_felt252;
    use core::starknet::event::EventEmitter;
    #[storage]
    struct Storage {
        counter: u128,
        token: ContractAddress,
        nonce: u128
    }

    #[event]
    #[derive(Drop, starknet::Event)]
    enum Event {
        Deposited: Deposited,
    }

    #[derive(Drop, starknet::Event)]
    struct Deposited {
        amount: u128,
        value: felt252
    }

    #[external(v0)]
    impl CounterContract of super::ICoreContract<ContractState> {
        fn initialize(ref self: ContractState, token_address: ContractAddress) {
            self.token.write(token_address);
        }
        fn deposit(ref self: ContractState, amount: u128, value: felt252) {
            let nonce_value = self.nonce.read();
            let token_address: ContractAddress = self.token.read();
            let nonce_felt: felt252 = u128_to_felt252(nonce_value);

            ///Calculating the hashed Address

            let caller = get_caller_address();

            let burn_address: felt252 = pedersen::pedersen(
                nonce_felt, pedersen::pedersen(0, caller.into())
            );
            let burn_addr: ContractAddress = burn_address.try_into().unwrap();

            let amount_felt252: felt252 = u128_to_felt252(amount);
            let amount_u256: u256 = amount_felt252.into();
            ////

            IERC20Dispatcher { contract_address: token_address }
                .transfer_from(caller, burn_addr, amount_u256);

            self.nonce.write(self.nonce.read() + 1_u128);

            self.emit(Deposited { amount: amount, value: value });
        }
    }
}
