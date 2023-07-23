package e2pg

import "github.com/indexsupply/x/pgmig"

var Migrations = map[int]pgmig.Migration{
	0: pgmig.Migration{
		SQL: `
			create table task (
				id smallint not null,
				number bigint,
				hash bytea,
				insert_at timestamptz default now()
			);
			create table nft_transfers (
				contract bytea,
				token_id numeric,
				quantity numeric,
				f bytea,
				t bytea,
				tx_sender bytea,
				eth numeric,
				task_id numeric,
				chain_id numeric,
				block_hash bytea,
				block_number numeric,
				transaction_hash bytea,
				transaction_index numeric,
				log_index numeric
			);
			create table erc20_transfers (
				contract bytea,
				f bytea,
				t bytea,
				value numeric,
				tx_sender bytea,
				eth numeric,
				task_id numeric,
				chain_id numeric,
				block_hash bytea,
				block_number numeric,
				transaction_hash bytea,
				transaction_index numeric,
				log_index numeric
			);
		`,
	},
	1: pgmig.Migration{
		SQL: `
			create table erc4337_userops (
				contract bytea,
				op_hash bytea,
				op_sender bytea,
				op_paymaster bytea,
				op_nonce numeric,
				op_success boolean,
				op_actual_gas_cost numeric,
				op_actual_gas_used numeric,
		
				tx_sender bytea,
				eth numeric,
				task_id numeric,
				chain_id numeric,
				block_hash bytea,
				block_number numeric,
				transaction_hash bytea,
				transaction_index numeric,
				log_index numeric
		);
		`,
	},
	2: pgmig.Migration{
		SQL: `
			alter table erc4337_userops
				add primary key (transaction_hash, log_index);
			alter table erc20_transfers
				add primary key (transaction_hash, log_index);
			alter table nft_transfers
				add primary key (transaction_hash, log_index);
		`,
	},
}
