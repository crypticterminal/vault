package main

import (
	"github.com/hashicorp/vault/vault"
)

func sysInit() Path {
	p := NewPath("/sys/init")

	m := NewMethod(vault.SysHelp["init"][0], "sys")
	m.addResponse(200, `{"initialized": true}`)
	p.Methods["get"] = m

	m = NewMethod(vault.SysHelp["init"][0], "sys")
	m.BodyProps = []Property{
		NewProperty("pgp_keys", "array/string",
			"Specifies an array of PGP public keys used to encrypt the output unseal keys. Ordering is preserved. The keys must be base64-encoded from their original binary representation. The size of this array must be the same as secret_shares."),
		NewProperty("root_token_pgp_key", "string",
			"Specifies a PGP public key used to encrypt the initial root token. The key must be base64-encoded from its original binary representation."),
		NewProperty("secret_shares", "number",
			"Specifies the number of shares to split the master key into."),
		NewProperty("secret_threshold", "number",
			"Specifies the number of shares required to reconstruct the master key. This must be less than or equal secret_shares. If using Vault HSM with auto-unsealing, this value must be the same as secret_shares."),
		NewProperty("stored_shares", "number",
			"Specifies the number of shares that should be encrypted by the HSM and stored for auto-unsealing. Currently must be the same as secret_shares."),
		NewProperty("recovery_pgp_keys", "array/string",
			"Specifies an array of PGP public keys used to encrypt the output recovery keys. Ordering is preserved. The keys must be base64-encoded from their original binary representation. The size of this array must be the same as recovery_shares."),
	}
	m.addResponse(200, `
		{
		  "keys": ["one", "two", "three"],
		  "keys_base64": ["cR9No5cBC", "F3VLrkOo", "zIDSZNGv"],
		  "root_token": "foo"
		}`)

	p.Methods["put"] = m

	return p
}

func sysLeader() Path {
	p := NewPath("/sys/leader")
	m := NewMethod("Check the high availability status and current leader of Vault", "sys")
	m.addResponse(200, `
		{
            "ha_enabled": true,
            "is_self": false,
            "leader_address": "https://127.0.0.1:8200/",
            "leader_cluster_address": "https://127.0.0.1:8201/"
        }`)
	p.Methods["get"] = m

	return p
}

func sealStatus() Path {
	p := NewPath("/sys/seal-status")
	m := NewMethod(vault.SysHelp["seal-status"][0], "sys")
	m.addResponse(200, `
		{
			  "type": "shamir",
			  "sealed": false,
			  "t": 3,
			  "n": 5,
			  "progress": 0,
			  "version": "0.9.0",
			  "cluster_name": "vault-cluster-d6ec3c7f",
			  "cluster_id": "3e8b3fec-3749-e056-ba41-b62a63b997e8",
			  "nonce": "ef05d55d-4d2c-c594-a5e8-55bc88604c24"
		}`)
	p.Methods["get"] = m

	return p
}

func seal() Path {
	p := NewPath("/sys/seal")
	m := NewMethod(vault.SysHelp["seal"][0], "sys")
	m.Responses = []Response{StdRespNoContent}
	p.Methods["get"] = m

	return p
}

func stepDown() Path {
	p := NewPath("/sys/step-down")
	m := NewMethod("Causes the node to give up active status.", "sys")
	m.Responses = []Response{StdRespNoContent}
	p.Methods["put"] = m

	return p
}

func sysHealth() Path {
	p := NewPath("/sys/health")
	m := NewMethod("Returns the health status of Vault.", "sys")
	m.Responses = []Response{
		NewResponse(200, "initialized, unsealed, and active", ""),
		NewResponse(429, "unsealed and standby", ""),
		NewResponse(472, "data recovery mode replication secondary and active", ""),
		NewResponse(501, "not initialized", ""),
		NewResponse(503, "sealed", ""),
	}

	p.Methods["get"] = m
	p.Methods["head"] = m

	return p
}

func unseal() Path {
	p := NewPath("/sys/unseal")
	m := NewMethod(vault.SysHelp["unseal"][0], "sys")
	m.BodyProps = []Property{
		NewProperty("key", "string", "Specifies a single master key share. This is required unless reset is true."),
		NewProperty("reset", "boolean", "Specifies if previously-provided unseal keys are discarded and the unseal process is reset."),
	}
	m.addResponse(200, `
		{
		  "sealed": false,
		  "t": 3,
		  "n": 5,
		  "progress": 0,
		  "version": "0.6.2",
		  "cluster_name": "vault-cluster-d6ec3c7f",
		  "cluster_id": "3e8b3fec-3749-e056-ba41-b62a63b997e8"
		}`)

	p.Methods["put"] = m

	return p
}

// This is so horrible
var sysPaths = []Path{
	sysLeader(),
	sysInit(),
	sysHealth(),
	sealStatus(),
	seal(),
	stepDown(),
	unseal(),
}
