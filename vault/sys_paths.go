package vault

import (
	"github.com/hashicorp/vault/apidoc/apidoc"
	"github.com/hashicorp/vault/logical/framework"
)

func sysGenerateRootAttempt() apidoc.Path {
	p := apidoc.NewPath("/sys/generate-root/attempt")

	// GET
	m := apidoc.NewMethod("Reads the configuration and process of the current root generation attempt.", "sys")
	m.AddResponse(200, `
	{
	  "started": true,
	  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
	  "progress": 1,
	  "required": 3,
	  "encoded_token": "",
	  "pgp_fingerprint": "816938b8a29146fbe245dd29e7cbaf8e011db793",
	  "complete": false
}`)
	p.Methods["get"] = m

	// PUT
	m = apidoc.NewMethod("Initializes a new root generation attempt", "sys")
	m.BodyProps = []apidoc.Property{
		apidoc.NewProperty("otp", "string", "Specifies a base64-encoded 16-byte value."),
		apidoc.NewProperty("pgp_key", "string", "Specifies a base64-encoded PGP public key."),
	}
	m.AddResponse(200, `
	{
	    "started": true,
	    "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
	    "progress": 1,
	    "required": 3,
	    "encoded_token": "",
	    "pgp_fingerprint": "",
	    "complete": false
	}`)
	p.Methods["put"] = m

	// DELETE
	m = apidoc.NewMethod("Cancels any in-progress root generation attempt.", "sys")
	m.Responses = []apidoc.Response{apidoc.StdRespNoContent}
	p.Methods["delete"] = m

	return p
}

func sysGenerateRootUpdate() apidoc.Path {
	p := apidoc.NewPath("/sys/generate-root/update")

	// PUT
	m := apidoc.NewMethod("Enter a single master key share to progress the root generation attempt.", "sys")
	m.BodyProps = []apidoc.Property{
		apidoc.NewProperty("key", "string", "Specifies a single master key share."),
		apidoc.NewProperty("nonce", "string", "Specifies the nonce of the attempt."),
	}
	m.AddResponse(200, `
	{
	  "started": true,
	  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
	  "progress": 3,
	  "required": 3,
	  "pgp_fingerprint": "",
	  "complete": true,
	  "encoded_token": "FPzkNBvwNDeFh4SmGA8c+w=="
	}`)
	p.Methods["put"] = m

	return p
}

func sysInit() apidoc.Path {
	p := apidoc.NewPath("/sys/init")

	m := apidoc.NewMethod(SysHelp["init"][0], "sys")
	m.AddResponse(200, `{"initialized": true}`)
	p.Methods["get"] = m

	m = apidoc.NewMethod(SysHelp["init"][0], "sys")
	m.BodyProps = []apidoc.Property{
		apidoc.NewProperty("pgp_keys", "array/string",
			"Specifies an array of PGP public keys used to encrypt the output unseal keys. Ordering is preserved. The keys must be base64-encoded from their original binary representation. The size of this array must be the same as secret_shares."),
		apidoc.NewProperty("root_token_pgp_key", "string",
			"Specifies a PGP public key used to encrypt the initial root token. The key must be base64-encoded from its original binary representation."),
		apidoc.NewProperty("secret_shares", "number",
			"Specifies the number of shares to split the master key into."),
		apidoc.NewProperty("secret_threshold", "number",
			"Specifies the number of shares required to reconstruct the master key. This must be less than or equal secret_shares. If using Vault HSM with auto-unsealing, this value must be the same as secret_shares."),
		apidoc.NewProperty("stored_shares", "number",
			"Specifies the number of shares that should be encrypted by the HSM and stored for auto-unsealing. Currently must be the same as secret_shares."),
		apidoc.NewProperty("recovery_pgp_keys", "array/string",
			"Specifies an array of PGP public keys used to encrypt the output recovery keys. Ordering is preserved. The keys must be base64-encoded from their original binary representation. The size of this array must be the same as recovery_shares."),
	}
	m.AddResponse(200, `
		{
		  "keys": ["one", "two", "three"],
		  "keys_base64": ["cR9No5cBC", "F3VLrkOo", "zIDSZNGv"],
		  "root_token": "foo"
		}`)

	p.Methods["put"] = m

	return p
}

func sysLeader() apidoc.Path {
	p := apidoc.NewPath("/sys/leader")
	m := apidoc.NewMethod("Check the high availability status and current leader of Vault", "sys")
	m.AddResponse(200, `
		{
            "ha_enabled": true,
            "is_self": false,
            "leader_address": "https://127.0.0.1:8200/",
            "leader_cluster_address": "https://127.0.0.1:8201/"
        }`)
	p.Methods["get"] = m

	return p
}

func sealStatus() apidoc.Path {
	p := apidoc.NewPath("/sys/seal-status")
	m := apidoc.NewMethod(SysHelp["seal-status"][0], "sys")
	m.AddResponse(200, `
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

func seal() apidoc.Path {
	p := apidoc.NewPath("/sys/seal")
	m := apidoc.NewMethod(SysHelp["seal"][0], "sys")
	m.Responses = []apidoc.Response{apidoc.StdRespNoContent}
	p.Methods["get"] = m

	return p
}

func stepDown() apidoc.Path {
	p := apidoc.NewPath("/sys/step-down")
	m := apidoc.NewMethod("Causes the node to give up active status.", "sys")
	m.Responses = []apidoc.Response{apidoc.StdRespNoContent}
	p.Methods["put"] = m

	return p
}

func sysHealth() apidoc.Path {
	p := apidoc.NewPath("/sys/health")
	m := apidoc.NewMethod("Returns the health status of Vault.", "sys")
	m.Responses = []apidoc.Response{
		apidoc.NewResponse(200, "initialized, unsealed, and active", ""),
		apidoc.NewResponse(429, "unsealed and standby", ""),
		apidoc.NewResponse(472, "data recovery mode replication secondary and active", ""),
		apidoc.NewResponse(501, "not initialized", ""),
		apidoc.NewResponse(503, "sealed", ""),
	}

	p.Methods["get"] = m
	p.Methods["head"] = m

	return p
}

func sysRekeyInit() apidoc.Path {
	p := apidoc.NewPath("/sys/rekey/init")

	// GET
	m := apidoc.NewMethod("Read the configuration and progress of the current rekey attempt.", "sys")
	m.AddResponse(200, `
	{
	  "started": true,
	  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
	  "t": 3,
	  "n": 5,
	  "progress": 1,
	  "required": 3,
	  "pgp_fingerprints": ["abcd1234"],
	  "backup": true
	}`)

	p.Methods["get"] = m

	// PUT
	m = apidoc.NewMethod("Initializes a new rekey attempt", "sys")
	m.BodyProps = []apidoc.Property{
		apidoc.NewProperty("secret_shares", "number",
			"Specifies the number of shares to split the master key into."),
		apidoc.NewProperty("secret_threshold", "number",
			"Specifies the number of shares required to reconstruct the master key. This must be less than or equal secret_shares. If using Vault HSM with auto-unsealing, this value must be the same as secret_shares."),
		apidoc.NewProperty("pgp_keys", "array/string",
			"Specifies an array of PGP public keys used to encrypt the output unseal keys. Ordering is preserved. The keys must be base64-encoded from their original binary representation. The size of this array must be the same as secret_shares."),
		apidoc.NewProperty("backup", "boolean", "Specifies if using PGP-encrypted keys, whether Vault should also store a plaintext backup of the PGP-encrypted keys."),
	}
	m.Responses = []apidoc.Response{apidoc.StdRespNoContent}
	p.Methods["put"] = m

	// DELETE
	m = apidoc.NewMethod("Cancels any in-progress rekey.", "sys")
	m.Responses = []apidoc.Response{apidoc.StdRespNoContent}
	p.Methods["delete"] = m
	return p
}

func sysRekeyUpdate() apidoc.Path {
	p := apidoc.NewPath("/sys/rekey/update")

	// PUT
	m := apidoc.NewMethod("Enter a single master key share to progress the rekey of the Vault.", "sys")
	m.BodyProps = []apidoc.Property{
		apidoc.NewProperty("key", "string", "Specifies a single master key share."),
		apidoc.NewProperty("nonce", "string", "Specifies the nonce of the rekey attempt."),
	}
	m.AddResponse(200, `
	{
	  "complete": true,
	  "keys": ["one", "two", "three"],
	  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
	  "pgp_fingerprints": ["abcd1234"],
	  "keys_base64": ["base64keyvalue"],
	  "backup": true
	}`)
	p.Methods["put"] = m

	return p
}

func sysRekeyBackup() apidoc.Path {
	p := apidoc.NewPath("/sys/rekey/backup")

	// GET
	m := apidoc.NewMethod("Return the backup copy of PGP-encrypted unseal keys.", "sys")
	m.AddResponse(200, `
	{
	  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
	  "keys": {
		"abcd1234": "..."
	  }
	}`)

	p.Methods["get"] = m

	// DELETE
	m = apidoc.NewMethod("Deletes the backup copy of PGP-encrypted unseal keys.", "sys")
	m.Responses = []apidoc.Response{apidoc.StdRespNoContent}
	p.Methods["delete"] = m

	return p
}

func sysRekeyRecoveryBackup() apidoc.Path {
	p := apidoc.NewPath("/sys/rekey-recovery-key/backup")

	// GET
	m := apidoc.NewMethod("Return the backup copy of PGP-encrypted recovery key shares.", "sys")
	m.AddResponse(200, `
	{
	  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
	  "keys": {
		"abcd1234": "..."
	  }
	}`)

	p.Methods["get"] = m

	// DELETE
	m = apidoc.NewMethod("Deletes the backup copy of PGP-encrypted recovery key shares.", "sys")
	m.Responses = []apidoc.Response{apidoc.StdRespNoContent}
	p.Methods["delete"] = m

	return p
}

func sysRekeyRecoveryInit() apidoc.Path {
	p := apidoc.NewPath("/sys/rekey-recovery-key/init")

	// GET
	m := apidoc.NewMethod("Read the configuration and progress of the current rekey attempt.", "sys")
	m.AddResponse(200, `
	{
	  "started": true,
	  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
	  "t": 3,
	  "n": 5,
	  "progress": 1,
	  "required": 3,
	  "pgp_fingerprints": ["abcd1234"],
	  "backup": true
	}`)

	p.Methods["get"] = m

	// PUT
	m = apidoc.NewMethod("Initializes a new rekey attempt", "sys")
	m.BodyProps = []apidoc.Property{
		apidoc.NewProperty("secret_shares", "number",
			"Specifies the number of shares to split the recovery key into."),
		apidoc.NewProperty("secret_threshold", "number",
			"Specifies the number of shares required to reconstruct the recovery key. This must be less than or equal secret_shares. If using Vault HSM with auto-unsealing, this value must be the same as secret_shares."),
		apidoc.NewProperty("pgp_keys", "array/string",
			"Specifies an array of PGP public keys used to encrypt the output unseal keys. Ordering is preserved. The keys must be base64-encoded from their original binary representation. The size of this array must be the same as secret_shares."),
		apidoc.NewProperty("backup", "boolean", "Specifies if using PGP-encrypted keys, whether Vault should also store a plaintext backup of the PGP-encrypted keys."),
	}
	m.Responses = []apidoc.Response{apidoc.StdRespNoContent}
	p.Methods["put"] = m

	// DELETE
	m = apidoc.NewMethod("Cancels any in-progress rekey.", "sys")
	m.Responses = []apidoc.Response{apidoc.StdRespNoContent}
	p.Methods["delete"] = m
	return p
}

func sysRekeyRecoveryUpdate() apidoc.Path {
	p := apidoc.NewPath("/sys/rekey-recovery-key/update")

	// PUT
	m := apidoc.NewMethod("Enter a single master key share to progress the rekey of the Vault.", "sys")
	m.BodyProps = []apidoc.Property{
		apidoc.NewProperty("key", "string", "Specifies a single master key share."),
		apidoc.NewProperty("nonce", "string", "Specifies the nonce of the rekey attempt."),
	}
	m.AddResponse(200, `
	{
	  "complete": true,
	  "keys": ["one", "two", "three"],
	  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
	  "pgp_fingerprints": ["abcd1234"],
	  "keys_base64": ["base64keyvalue"],
	  "backup": true
	}`)
	p.Methods["put"] = m

	return p
}

func sysWrappingLookup() apidoc.Path {
	p := apidoc.NewPath("/sys/wrapping/lookup")

	// POST
	m := apidoc.NewMethod("Look up wrapping properties for the given token.", "sys")
	m.BodyProps = []apidoc.Property{
		apidoc.NewProperty("token", "string", "Specifies the wrapping token ID."),
	}
	m.AddResponse(200, `
	{
	  "request_id": "481320f5-fdf8-885d-8050-65fa767fd19b",
	  "lease_id": "",
	  "lease_duration": 0,
	  "renewable": false,
	  "data": {
		"creation_path": "sys/wrapping/wrap",
		"creation_time": "2016-09-28T14:16:13.07103516-04:00",
		"creation_ttl": 300
	  },
	  "wrap_info": null,
	  "warnings": null,
	  "auth": null
	}`)

	p.Methods["post"] = m

	return p
}

func unseal() apidoc.Path {
	p := apidoc.NewPath("/sys/unseal")
	m := apidoc.NewMethod(SysHelp["unseal"][0], "sys")
	m.BodyProps = []apidoc.Property{
		apidoc.NewProperty("key", "string", "Specifies a single master key share. This is required unless reset is true."),
		apidoc.NewProperty("reset", "boolean", "Specifies if previously-provided unseal keys are discarded and the unseal process is reset."),
	}
	m.AddResponse(200, `
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

type DocExport struct {
}

func (d DocExport) BackendPaths() []*framework.Path {
	return NewSystemBackend(&Core{}, nil).Paths
}

func (d DocExport) ManualPaths() []apidoc.Path {
	return []apidoc.Path{
		sysGenerateRootAttempt(),
		sysGenerateRootUpdate(),
		sysLeader(),
		sysInit(),
		sysHealth(),
		sysRekeyInit(),
		sysRekeyUpdate(),
		//sysRekeyBackup(),
		sysRekeyRecoveryInit(),
		sysRekeyRecoveryUpdate(),
		sysRekeyRecoveryBackup(),
		//sysWrappingLookup(),
		//sysWrappingRewrap(),
		//sysWrappingUnwrap(),
		sealStatus(),
		seal(),
		stepDown(),
		unseal(),
	}
}
