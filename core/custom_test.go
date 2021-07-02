package core

import (
	"errors"
	"os"

	"github.com/manifoldco/promptui"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/miniscruff/changie/testutils"
)

var _ = Describe("Custom", func() {
	It("returns error on invalid prompt type", func() {
		_, err := Custom{
			Type: "invalid type",
		}.CreatePrompt(os.Stdin)
		Expect(errors.Is(err, ErrInvalidPromptType)).To(BeTrue())
	})

	It("can create custom string prompt", func() {
		prompt, err := Custom{
			Type:  CustomString,
			Label: "a label",
		}.CreatePrompt(os.Stdin)
		Expect(err).To(BeNil())

		underPrompt, ok := prompt.(*promptui.Prompt)
		Expect(ok).To(BeTrue())
		Expect(underPrompt.Label).To(Equal("a label"))
	})

	It("can create custom int prompt", func() {
		prompt, err := Custom{
			Key:  "name",
			Type: CustomInt,
		}.CreatePrompt(os.Stdin)
		Expect(err).To(BeNil())

		underPrompt, ok := prompt.(*promptui.Prompt)
		Expect(ok).To(BeTrue())
		Expect(underPrompt.Validate("12")).To(BeNil())
		Expect(underPrompt.Validate("not an int")).NotTo(BeNil())
		Expect(underPrompt.Validate("12.5")).NotTo(BeNil())
		Expect(underPrompt.Label).To(Equal("name"))
	})

	It("can create custom int prompt with a min and max", func() {
		var min int64 = 5
		var max int64 = 10
		prompt, err := Custom{
			Type:   CustomInt,
			MinInt: &min,
			MaxInt: &max,
		}.CreatePrompt(os.Stdin)
		Expect(err).To(BeNil())

		underPrompt, ok := prompt.(*promptui.Prompt)
		Expect(ok).To(BeTrue())
		Expect(underPrompt.Validate("7")).To(BeNil())
		Expect(errors.Is(underPrompt.Validate("3"), errIntTooLow)).To(BeTrue())
		Expect(errors.Is(underPrompt.Validate("25"), errIntTooHigh)).To(BeTrue())
	})

	It("can create custom enum prompt", func() {
		prompt, err := Custom{
			Type:        CustomEnum,
			EnumOptions: []string{"a", "b", "c"},
		}.CreatePrompt(os.Stdin)
		Expect(err).To(BeNil())

		underPrompt, ok := prompt.(*enumWrapper)
		Expect(ok).To(BeTrue())
		Expect(underPrompt.Select.Items).To(Equal([]string{"a", "b", "c"}))
	})

	It("can run enum prompt", func() {
		stdinReader, stdinWriter, err := os.Pipe()
		Expect(err).To(BeNil())

		defer stdinReader.Close()
		defer stdinWriter.Close()

		prompt, err := Custom{
			Type:        CustomEnum,
			EnumOptions: []string{"a", "b", "c"},
		}.CreatePrompt(stdinReader)
		Expect(err).To(BeNil())

		go func() {
			DelayWrite(stdinWriter, []byte{106, 13})
		}()

		out, err := prompt.Run()
		Expect(err).To(BeNil())
		Expect(out).To(Equal("b"))
	})
})
