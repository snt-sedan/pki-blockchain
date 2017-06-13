// This code parses the x.509 certificates that include extension for blockchain
// Author Wazen SHBAIR
// Snt, university of Luxembourg

package main
import (
	"crypto/x509"
	"io/ioutil"
	"fmt"
	"os"
)

func main() {
	// Certificate types
	//BlockchainTestRootCA3 % BlockchainEndUser % BlockchainTestIssuingCA3
	fmt.Println("\t--------------------------------")
	cert := os.Args[1]
	//cert:="BlockchainEndUser.cer"
	fmt.Println("\tParsing : " + cert)
	fmt.Println("\t--------------------------------")
	certPEMBlock, _ := ioutil.ReadFile(cert)
	ca, _ := x509.ParseCertificate(certPEMBlock)
	// iterate in the extension to get the information
	for _, element := range ca.Extensions {
		if element.Id.String() == "1.2.752.115.33.3"{ //<-- this string comes from the certifcate it self
			fmt.Printf("\tBlockChain Name: %+#+x\n", element.Value)
		}
		if element.Id.String() == "1.2.752.115.33.2"{
			fmt.Printf("\tCaContractIdentifier: %+#+x\n", element.Value)
		}
		if element.Id.String() == "1.2.752.115.33.1"{
			fmt.Printf("\tIssuerCaContractIdentifier: %+#+x\n", element.Value)
		}
		if element.Id.String() == "1.2.752.115.33.0"{
			fmt.Printf("\tHash Algo: %+#+x\n", element.Value)
		}
	}
}
