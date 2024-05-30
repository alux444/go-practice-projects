package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("domain // hasMx // hasSpf // spfRecord // hasDmarc // dmarcRecord\n")
	fmt.Printf("hasMx - Whether domain has Mail Exchange records set up\n")
	fmt.Printf("hasSpf - Whether domain has Sender Policy Framework set up\n")
	fmt.Printf("hasDmarc - Whether domain has Domain-based Message Authentication, Response and Conformance records\n")

	for scanner.Scan() {
		checkDomain(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Fatal error: %v\n", err)
	}
}

func checkDomain(domain string) {
	var hasMx, hasSpf, hasDmarc bool
	var spfRecord, dmarcRecord string

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error getting MX records: %v\n", err)
	}
	if len(mxRecords) > 0 {
		hasMx = true
	}

	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error getting txt records: %v\n", err)
	}
	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSpf = true
			spfRecord = record
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error getting dmarc records: %v\n", err)
	}
	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDmarc = true
			dmarcRecord = record
			break
		}
	}

	fmt.Printf("%s hasMx: %v\n", domain, hasMx)
	fmt.Printf("%s hasSpf: %v\n", domain, hasSpf)
	fmt.Printf("%s spfRecord: %v\n", domain, spfRecord)
	fmt.Printf("%s hasDmarc: %v\n", domain, hasDmarc)
	fmt.Printf("%s dmarcRecord: %v\n", domain, dmarcRecord)
}
