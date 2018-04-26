# PKI-blockchain

Public-Key Infrastructure (PKI) is the cornerstone technology that facilitates secure information exchange over the Internet. However, PKI is exposed to risks due to potential failures of Certificate Authorities (CAs) that may be used to issue unauthorized certificates for end-users. Many recent breaches show that if a CA is compromised, the security of the corresponding end-users will be in risk. 

As an emerging solution, Blockchain technology potentially resolves the problems of traditional PKI systems - in particular, elimination of single point-of-failure and rapid reaction to CAs shortcomings. Blockchain has the ability to store and manage digital certificates within a public and immutable ledger, resulting in a fully traceable history log. 

We designed and developed a blockchain-based PKI management framework for issuing, validating and revoking X.509 certificates. Evaluation and experimental results confirm that the proposed framework provides more reliable and robust PKI systems with modest maintenance costs.

------------------------------------------
PKI: Proof-of-Concept of blockchain PKI implementation (Ethereum).


Web server (pki-web.go) - testing interface to interract with REST service.
It allows to:
1. Navigate within CA tree
2. Generate testing X.509 cetrificates
3. With REST service (pki-rest.go) it adds certificates to the CA's white lists and black lists stored in the blockchain, validates certificates, etc


-----------------------------------------
REST API for Blockchain PKI (pki-rest.go):


A. ENROLL
/enroll_user, all parameters in POST

Parameters:
1. Hash or UplFiles (hash is a hex string without a leading "0x")
2. UplFiles : uploaded certificate
3. ParentAddr: the address of the CA smart contract where the certificate's hash is stored.This address of this contract should be called at user account CurrentUserAddr
4. CurrentUserAddr: the ID (address) of the user who has the privilage to modify the parent smart contract. The key of this user should be available in key storage

Returns: 200 and "OK" in the html body in case of success

Errors (details are in html body):
1. 480 : hash has the wrong length or hash is incorrect
2. 481 : hash is already enrolled
3. 482 : Certificate errors in case it was provided instead of hash
4. 484 : ParentAddr is incorrect
5. 485 : CurrentUserAddr is incorrect
6. 580 : Ethereum execution error (out of gas and others)
7. 581 : Ethereum connection error
8. 500 : Other error


B. BLACKLIST 
/blacklist-user, all parameters are in POST

Puts certificate (either ordinary or CA) from the white list to the black list

Parameters:
1. ParentAddr: the address of the CA smart contract where the certificate's hash is stored
2. UserAddr: the ID (address) of the user who has the privilage to modify the smart contract. The key of this user should be available in key storage
3. Deletion: array of strings with IDs of the items to be deleted in the user list. It is produced with checkbox HTML forms

Returns: 200 and "OK" in the html body in case of success
		
Errors (details are in html body):
1. 484 : ParentAddr is incorrect
2. 485 : Deletion is incorrect
3. 580 : Ethereum executionn error (out of gas and others)
4. 581 : Ethereum connection error
5. 500 : Other error


C. CREATE CONTRACT
/create_contract, all params as POST

Creation of the "empty" CA smart contract:
1. CA certificate should be added to smart contract through population procedure
2. the right to execute the smart contract should be changed to the CA account with population procedure as well
	
Params:
1. ParentAddr: the address of the CA smart contract which is used for creation (it has the bin code). This address of this contract should be called at user account CurrentUserAddr
2. NewUserAddr - owner is set to this address at the end of the proc. If empty, then new owner is not set. At the end of the population procedure only the NewUserAddr can modify the smart contract in the future
3. CurrentUserAddr: - the user addr to connect to Ethereum. If empty, then set to root user addr

Returns: 200 and the smart contract address WITH heading "0x" in the html body in case of success

Errors (details are in html body):
1. 480 : Current user does not have rights to execute the creation of the CA certificate
2. 483 : NewUserAddr is incorrect
3. 484 : ParentAddr is incorrect
4. 485 : CurrentUserAddr is incorrect
5. 580 : Ethereum execution error (out of gas and others)
6. 581 : Ethereum connection error
7. 500 : Other error


D. POPULATE
/populate_contract, all parameters in POST

Population of the CA smart contract:
1. putting a certificate into the contract referencing its parent, and
2. setting ownership of the smartcontract to the user

Params:
1. UplFiles : uploaded certificate
2. NewUserAddr - owner is set to this address at the end of the proc. If empty, then new owner is not set. At the end of the population procedure only the NewUserAddr can modify the smart contract in the future
3. CurrentUserAddr: - the user addr to connect to Ethereum. If empty, then set to root user addr
4. ContrAddr: the address of the CA smart contract which should be populated. This address of this contract should be called at user account CurrentUserAddr

Returns: 200 and hash string WITHOUT heading "0x" in the html body in case of success

Errors (details are in html body):
1. 482 : Certificate errors
2. 483 : NewUserAddr is incorrect
3. 484 : ContrAddr is incorrect
4. 485 : CurrentUserAddr is incorrect
5. 580 : Ethereum execution error (out of gas and others)
6. 581 : Ethereum connection error
7. 500 : Other error


E. DOWNLOADING OF CA CERTIFICATE FROM BLOCKCHAIN
/download_cacert
Extracting (download) of certificate from CA smart contract

Params:
1. ContrAddr: the address of the CA smart contract

Returns:
200 and the smart contract address WITH heading "0x" in the html body in case of success

Errors (details are in html body):
1. 484 : ContrAddr is incorrect
2. 580 : Ethereum execution error (out of gas and others)
3. 581 : Ethereum connection error
4. 500 : Other error


F. CERTIFICATE VALIDATION
/validate_cert, all params as POST
	
Parameters:
1. Hash or UplFiles (hash is a hex string without a leading "0x")
2. UplFiles : uploaded certificate
3. ParentAddr: the address of the CA smart contract where the certificate's hash is stored. If certificate is uploaded through UplFiles, ParentAddr may not be specified

Returns: 200 and JSON with the validation results in the html body in case of success
		
Errors (details are in html body):
1. 480 : hash has wrong length or hash is incorrect
2. 482 : Certificate errors in case it was provided instead of hash
3. 484 : ParentAddr is incorrect
4. 580 : Ethereum execution error (out of gas and others)
5. 581 : Ethereum connection error
6. 500 : Other error

Validation result codes (used in smart contract validation, for instance):
1. 0  - OK, the certificate is valid
2. 1  - certificate not found
3. 2  - certificate revoked
4. 11 - error in parsing
5. 12 - CA addr in the certificate does not correspond to _addrCA
6. 13 - empty cert received for this CA
7. 14 - empty addrCA parsed in this CA cert
8. 15 - parent addr is null, but CA addr does not correspond to Root addr  
9. 16 - too many iterations: a certain limit (100?) is exceeded

---------------------------------------------------
SMART CONTRACT VERIFICATION:
1. CheckCert(Hash, CurrentContract, Root Contract)  – the verification itself
2. DecodeReturnErr(returnCode_from_CheckCert) – the validation result code
3. DecodeReturnIter(returnCode_from_CheckCert) – the level from leaf to the root (starting from 0) where the code should be applied. If certificate is valid the number of levels to the root from the leaf

Example: 

verif = eth.contract({ABI}).at({VERIF_CONTR_ADDRESS})

// To get verification result (0 if valid):
verif.DecodeReturnErr(verif.CheckCert("0x4c39a4efe6a1266bb4d479716fc0a674128c5437ba6ddafe63ba326307c430f9","0x27290fea2bf264b221ba1e97518650fcce1cf0d5","0x778d81a6563d3bd442b844849abde2959e8a0dc7"))

// To get level to the root at which the verification was not valid or number of iterations to the root if certificate is valid
verif.DecodeReturnIter(verif.CheckCert("0x4c39a4efe6a1266bb4d479716fc0a674128c5437ba6ddafe63ba326307c430f9","0x27290fea2bf264b221ba1e97518650fcce1cf0d5","0x778d81a6563d3bd442b844849abde2959e8a0dc7"))


---------------------------------------------------------------
Installation:

A. Compile smart contracts 
1. cd ./scontract
2. solc --bin pki_scont.sol > bin/pki_scont.bin
3. solc --bin pki_scont_web.sol > bin/pki_scont_web.bin
4. solc --bin pki_scont_valid.sol > bin/pki_scont_valid.bin
5. solc --abi pki_scont.sol > abi/pki_abi.json
6. solc --abi pki_scont_web.sol > abi/pki_abi_web.json
7. solc --abi pki_scont_valid.sol > abi/pki_abi_valid.json


B. Deploy smart contracts to Ethereum (public or private) 

C. Generate bindings for smart contracts in Golang
1. cd ../
2. abigen --abi ./scontract/abi/pki_abi.json --pkg main --type LuxUni_PKI --out bind_pki.go --bin ./scontract/bin/pki_scont.bin
3. abigen --abi ./scontract/abi/pki_abi_web.json --pkg main --type LuxUni_PKI_web --out bind_pki_web.go
4. abigen --abi ./scontract/abi/pki_abi_valid.json --pkg main --type LuxUni_PKI_valid --out bind_pki_valid.go

D. Compile REST service (pki-rest) and testing web server (pki-web)  
1. cd {go-ethereum dir}
2. godep go build {current PKI dir}/pki-rest.go {current PKI dir}/bind_pki.go {current PKI dir}/pki_conf.go
3. godep go build {current PKI dir}/pki-web.go {current PKI dir}/bind_pki.go {current PKI dir}/bind_pki_web.go {current PKI dir}/pki_conf.go

E. Configure pki-rest and pki-web, if needed (update of smart contract address, etc)
By default, configuration file is ./config/pki-conf.json



