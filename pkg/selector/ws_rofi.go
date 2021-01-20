package selector

import (
	"fmt"
	"github.com/geshido/i3-tools/pkg/i3tools"
	"go.i3wm.org/i3/v4"
	"io/ioutil"
	"os/exec"
	"sort"
	"strings"
)

type rofi struct {
	prompt        string
	emptyChoices  bool
	choicesGetter func() ([]string, error)
}

type RofiOption func(r *rofi)

func NewRofi(prompt string, opts ...RofiOption) *rofi {
	rofi := rofi{
		prompt:        prompt,
		choicesGetter: SortedI3Workspaces,
	}
	for _, opt := range opts {
		opt(&rofi)
	}
	return &rofi
}
func RofiOptionChoicesGetter(fn func() ([]string, error)) RofiOption {
	return func(r *rofi) {
		r.choicesGetter = fn
	}
}
func RofiOptionNoChoices() RofiOption {
	return func(r *rofi) {
		r.emptyChoices = true
	}
}

func (r rofi) Select() (string, error) {
	arguments := []string{
		"-dmenu",
		"-i",
		"-p",
		r.prompt,
	}
	if r.emptyChoices {
		arguments = append(arguments, "-lines")
		arguments = append(arguments, "0")
	}
	cmd := exec.Command("rofi", arguments...)
	if !r.emptyChoices {
		choices, err := r.choicesGetter()
		if err != nil {
			return "", fmt.Errorf("cannot get choices: %w", err)
		}
		in, err := cmd.StdinPipe()
		if err != nil {
			return "", fmt.Errorf("cannot open stdin pipe: %w", err)
		}
		_, err = in.Write([]byte(strings.Join(choices, "\n")))
		if err != nil {
			return "", fmt.Errorf("cannot write to rofi stdin: %w", err)
		}
		if err := in.Close(); err != nil {
			return "", fmt.Errorf("cannot close rofi's stdin: %w", err)
		}
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("cannot open stdout pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("cannot start rofi: %w", err)
	}

	b, err := ioutil.ReadAll(out)
	if err != nil {
		return "", fmt.Errorf("cannot get rofi's output: %w", err)
	}
	if err := out.Close(); err != nil {
		return "", fmt.Errorf("cannot close rofi's stdout: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return "", nil
		}
		return "", fmt.Errorf("rofi fails to run: %w", err)
	}

	return string(b), nil
}

func SortedI3Workspaces() ([]string, error) {
	wsList, err := i3.GetWorkspaces()
	if err != nil {
		return nil, fmt.Errorf("cannot get workspace list: %w", err)
	}
	sortable := i3tools.MakeSortableBy(wsList, i3tools.DefaultWorkspaceLessFn)
	sort.Sort(sortable)
	var choices []string
	for _, ws := range sortable.Data {
		choices = append(choices, ws.Name)
	}
	return choices, nil
}
