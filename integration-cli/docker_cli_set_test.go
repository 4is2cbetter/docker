package main

import (
	"strconv"
	"strings"

	"os/exec"
	"testing"
)

func TestSetContainer(t *testing.T) {
	cmd := exec.Command(dockerBinary, "run", "-d", "--name", "test-set-container", "-m", "300M", "busybox", "true")
	_, err := runCommand(cmd)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteAllContainers()

	cmd = exec.Command(dockerBinary, "set", "-m", "500M", "test-set-container")
	_, err = runCommand(cmd)
	if err != nil {
		t.Fatal(err)
	}

	cmd = exec.Command(dockerBinary, "inspect", "-f", "{{.HostConfig.Memory}}", "test-set-container")
	memory, _, err := runCommandWithOutput(cmd)
	if err != nil {
		t.Fatal(err)
	}

	float_mem, err := strconv.ParseFloat(strings.Trim(memory, "\n"), 64)
	if err != nil {
		t.Fatal(err)
	}
	int_mem := int(float_mem)

	if int_mem != 524288000 {
		t.Fatalf("Got the wrong memory value, we got %d, expected 524288000(500M).", int_mem)
	}
}

func TestSetContainerInvalidValue(t *testing.T) {
	cmd := exec.Command(dockerBinary, "run", "-d", "--name", "test-set-container", "-m", "300M", "busybox", "true")
	_, err := runCommand(cmd)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteAllContainers()

	cmd = exec.Command(dockerBinary, "set", "-m", "2M", "test-set-container")
	_, err = runCommand(cmd)
	if err == nil {
		t.Fatal("[set] should failed if we tried to set invalid value.")
	}
}
