/*
 * Minio Cloud Storage, (C) 2015, 2016 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// createCertsPath create certs path.
func createCertsPath() error {
	certsPath, err := getCertsPath()
	if err != nil {
		return err
	}
	return os.MkdirAll(certsPath, 0700)
}

// getCertsPath get certs path.
func getCertsPath() (string, error) {
	var certsPath string
	configDir, err := getConfigPath()
	if err != nil {
		return "", err
	}
	certsPath = filepath.Join(configDir, globalMinioCertsDir)
	return certsPath, nil
}

// mustGetCertsPath must get certs path.
func mustGetCertsPath() string {
	certsPath, err := getCertsPath()
	fatalIf(err, "Failed to get certificate path.")
	return certsPath
}

// mustGetCertFile must get cert file.
func mustGetCertFile() string {
	return filepath.Join(mustGetCertsPath(), globalMinioCertFile)
}

// mustGetKeyFile must get key file.
func mustGetKeyFile() string {
	return filepath.Join(mustGetCertsPath(), globalMinioKeyFile)
}

// isCertFileExists verifies if cert file exists, returns true if
// found, false otherwise.
func isCertFileExists() bool {
	st, e := os.Stat(filepath.Join(mustGetCertsPath(), globalMinioCertFile))
	// If file exists and is regular return true.
	if e == nil && st.Mode().IsRegular() {
		return true
	}
	return false
}

// isKeyFileExists verifies if key file exists, returns true if found,
// false otherwise.
func isKeyFileExists() bool {
	st, e := os.Stat(filepath.Join(mustGetCertsPath(), globalMinioKeyFile))
	// If file exists and is regular return true.
	if e == nil && st.Mode().IsRegular() {
		return true
	}
	return false
}

// isSSL - returns true with both cert and key exists.
func isSSL() bool {
	return isCertFileExists() && isKeyFileExists()
}

// Reads certificated file and returns a list of parsed certificates.
func readCertificateChain() ([]*x509.Certificate, error) {
	file, err := os.Open(mustGetCertFile())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the cert successfully.
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Proceed to parse the certificates.
	return parseCertificateChain(bytes)
}

// Parses certificate chain, returns a list of parsed certificates.
func parseCertificateChain(bytes []byte) ([]*x509.Certificate, error) {
	var certs []*x509.Certificate
	var block *pem.Block
	current := bytes

	// Parse all certs in the chain.
	for len(current) > 0 {
		block, current = pem.Decode(current)
		if block == nil {
			return nil, errors.New("Could not PEM block")
		}
		// Parse the decoded certificate.
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		certs = append(certs, cert)

	}
	return certs, nil
}
