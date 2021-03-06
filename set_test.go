package goml_test

import (
	"io/ioutil"
	"os"

	. "github.com/JulzDiverse/goml"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smallfish/simpleyaml"
)

var _ = Describe("Set", func() {
	var yml *simpleyaml.Yaml
	var err error

	Context("Top-level is map", func() {
		BeforeEach(func() {
			yaml := `map:
  name: foo

array:
- bar
- var
- zar

mapArray:
- foo: bar
  zoo: lion
  arr:
  - one
  - two
  - three
- foo: var
  boo: laa`

			yml, err = simpleyaml.NewYaml([]byte(yaml))
			Expect(err).NotTo(HaveOccurred())
		})

		It("should add an value to an array", func() {
			err = Set(yml, "array.2", "bumblebee")
			Expect(err).NotTo(HaveOccurred())
			err = Set(yml, "array.+", "optimusPrime")
			Expect(err).NotTo(HaveOccurred())
			err = Set(yml, "mapArray.0.foo", "wolverine")
			Expect(err).NotTo(HaveOccurred())
			err = Set(yml, "mapArray.foo:var.boo", "baymax")
			Expect(err).NotTo(HaveOccurred())

			Expect(Get(yml, "array.2")).To(Equal("bumblebee"))
			Expect(Get(yml, "array.3")).To(Equal("optimusPrime"))
			Expect(Get(yml, "mapArray.0.foo")).To(Equal("wolverine"))
			Expect(Get(yml, "mapArray.foo:var.boo")).To(Equal("baymax"))

			err = Set(yml, "array.:optimusPrime", "pikachu")
			Expect(Get(yml, "array.:pikachu")).To(Equal("pikachu"))
			err = Set(yml, "mapArray.foo:wolverine.arr.0", "new")
			Expect(Get(yml, "mapArray.foo:wolverine.arr.0")).To(Equal("new"))
		})

		Context("If a path does not exist", func() {
			It("should create the path", func() {
				err = Set(yml, "map.awesome", "bam")
				Expect(err).ToNot(HaveOccurred())

				Expect(Get(yml, "map.awesome")).To(Equal("bam"))

				err = Set(yml, "mapArray.luffy:gomugomuno.beat", "katakuri")
				Expect(err).ToNot(HaveOccurred())

				Expect(Get(yml, "mapArray.luffy:gomugomuno.beat")).To(Equal("katakuri"))
			})
		})

		Context("using the | delimiter", func() {
			It("should add an value to an array", func() {
				err = Set(yml, "mapArray.foo|var.boo", "baymax")
				Expect(err).NotTo(HaveOccurred())
				Expect(Get(yml, "mapArray.foo|var.boo")).To(Equal("baymax"))
				err = Set(yml, "array.+", "pikachu")
				Expect(Get(yml, "array.|pikachu")).To(Equal("pikachu"))
				err = Set(yml, "array.|pikachu", "ruffy")
				Expect(Get(yml, "array.|ruffy")).To(Equal("ruffy"))
				err = Set(yml, "mapArray.foo|bar.arr.0", "new")
				Expect(Get(yml, "mapArray.foo|bar.arr.0")).To(Equal("new"))
			})
		})
	})

	Context("Top-level is array", func() {
		var tempFile *os.File

		BeforeEach(func() {
			tempFile, err = ioutil.TempFile("", "")
			Expect(err).NotTo(HaveOccurred())
			_, err = tempFile.WriteString(`- key1: value1
- key2: value2`)
			Expect(err).NotTo(HaveOccurred())
			Expect(tempFile.Close()).To(Succeed())
		})

		AfterEach(func() {
			os.Remove(tempFile.Name())
		})

		It("changes an element in the array without error", func() {
			err = SetInFile(tempFile.Name(), "0.key1", "bumblebee")
			Expect(err).NotTo(HaveOccurred())
			Expect(GetFromFile(tempFile.Name(), "0.key1")).To(Equal("bumblebee"))
		})
	})
})
