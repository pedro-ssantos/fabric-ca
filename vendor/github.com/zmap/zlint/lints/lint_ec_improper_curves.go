package lints

/*
 * ZLint Copyright 2018 Regents of the University of Michigan
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License. You may obtain a copy
 * of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
 * implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

/************************************************
BRs: 6.1.5
Certificates MUST meet the following requirements for algorithm type and key size.
ECC Curve: NIST P-256, P-384, or P-521
************************************************/

import (
	"crypto/ecdsa"

	"github.com/zmap/zcrypto/x509"
	"github.com/zmap/zlint/util"
)

type ecImproperCurves struct{}

func (l *ecImproperCurves) Initialize() error {
	return nil
}

func (l *ecImproperCurves) CheckApplies(c *x509.Certificate) bool {
	return c.PublicKeyAlgorithm == x509.ECDSA
}

func (l *ecImproperCurves) Execute(c *x509.Certificate) *LintResult {
	/* Declare theKey to be a ECDSA Public Key */
	var theKey *ecdsa.PublicKey
	/* Need to do different things based on what c.PublicKey is */
	switch c.PublicKey.(type) {
	case *x509.AugmentedECDSA:
		temp := c.PublicKey.(*x509.AugmentedECDSA)
		theKey = temp.Pub
	case *ecdsa.PublicKey:
		theKey = c.PublicKey.(*ecdsa.PublicKey)
	}
	/* Now can actually check the params */
	theParams := theKey.Curve.Params()
	switch theParams.Name {
	case "P-256", "P-384", "P-521":
		return &LintResult{Status: Pass}
	default:
		return &LintResult{Status: Error}
	}
}

func init() {
	RegisterLint(&Lint{
		Name:        "e_ec_improper_curves",
		Description: "Only one of NIST P‐256, P‐384, or P‐521 can be used",
		Citation:    "BRs: 6.1.5",
		Source:      CABFBaselineRequirements,
		// Refer to BRs: 6.1.5, taking the statement "Before 31 Dec 2010" literally
		EffectiveDate: util.ZeroDate,
		Lint:          &ecImproperCurves{},
	})
}
