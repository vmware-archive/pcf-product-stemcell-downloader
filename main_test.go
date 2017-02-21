// Copyright 2017-Present Pivotal Software, Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main_test

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const metadataContents = `---
stemcell_criteria:
  version: "3262.4"
  os: ubuntu-trusty
  requires_cpi: false
`

var _ = Describe("download-bosh-io-stemcell", func() {
	var (
		ertFilepath string
		downloadDir string
	)

	BeforeEach(func() {
		ertFile, err := ioutil.TempFile("", "ert")
		Expect(err).NotTo(HaveOccurred())
		ertFilepath = ertFile.Name()

		zipWriter := zip.NewWriter(ertFile)

		fh, err := zipWriter.Create("metadata/cf.yml")
		Expect(err).NotTo(HaveOccurred())

		_, err = fh.Write([]byte(metadataContents))
		Expect(err).NotTo(HaveOccurred())

		err = zipWriter.Close()
		Expect(err).NotTo(HaveOccurred())

		downloadDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := os.Remove(ertFilepath)
		Expect(err).NotTo(HaveOccurred())

		err = os.RemoveAll(downloadDir)
		Expect(err).NotTo(HaveOccurred())
	})

	It("downloads a stemcell from bosh.io, based on the version specified in an ERT", func() {
		command := exec.Command(pathToMain,
			"--product-file", ertFilepath,
			"--product-name", "cf",
			"--iaas-type", "aws",
			"--download-dir", downloadDir)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session, "30s").Should(gexec.Exit(0))

		Eventually(session.Out).Should(gbytes.Say(fmt.Sprintf("Extracting 'metadata/cf.yml' from '%s'...", ertFilepath)))
		Eventually(session.Out).Should(gbytes.Say(`Downloading stemcell`))
		Eventually(session.Out).Should(gbytes.Say("Done."))

		Expect(filepath.Join(downloadDir, "bosh-stemcell-3262.4-aws-xen-hvm-ubuntu-trusty-go_agent.tgz")).To(BeARegularFile())

		fh, err := os.Open(filepath.Join(downloadDir, "bosh-stemcell-3262.4-aws-xen-hvm-ubuntu-trusty-go_agent.tgz"))
		Expect(err).NotTo(HaveOccurred())

		fi, err := fh.Stat()
		Expect(err).NotTo(HaveOccurred())

		Expect(fi.Size()).NotTo(BeZero())

		err = fh.Close()
		Expect(err).NotTo(HaveOccurred())
	})

	It("downloads the raw openstack stemcell", func() {
		command := exec.Command(pathToMain,
			"--product-file", ertFilepath,
			"--product-name", "cf",
			"--iaas-type", "openstack",
			"--download-dir", downloadDir)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session, "300s").Should(gexec.Exit(0))

		Eventually(session.Out).Should(gbytes.Say(fmt.Sprintf("Extracting 'metadata/cf.yml' from '%s'...", ertFilepath)))
		Eventually(session.Out).Should(gbytes.Say(`Downloading stemcell`))
		Eventually(session.Out).Should(gbytes.Say("Done."))

		fh, err := os.Open(filepath.Join(downloadDir, "bosh-stemcell-3262.4-openstack-kvm-ubuntu-trusty-go_agent-raw.tgz"))
		Expect(err).NotTo(HaveOccurred())

		fi, err := fh.Stat()
		Expect(err).NotTo(HaveOccurred())

		Expect(fi.Size()).NotTo(BeZero())

		err = fh.Close()
		Expect(err).NotTo(HaveOccurred())
	})

	Context("Failure cases", func() {
		Context("when an invalid iaas is specified", func() {
			It("outputs an error and exits with a non-zero status", func() {
				command := exec.Command(pathToMain,
					"--product-file", ertFilepath,
					"--product-name", "cf",
					"--iaas-type", "foo",
					"--download-dir", downloadDir)
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, "30s").Should(gexec.Exit(1))

				Eventually(session.Err).Should(gbytes.Say("Invalid IaaS: foo"))
			})
		})
	})
})
