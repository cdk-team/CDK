package kubectl

import (
	"bytes"
	_ "embed"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed assets/kubectl-amd64
var kubectlBinary []byte

func ExtractKubectl() (string, error) {
	tmpDir, err := os.MkdirTemp("", ".bin")
	if err != nil {
		return "", err
	}

	kubectlPath := filepath.Join(tmpDir, "kubectl")

	err = os.WriteFile(kubectlPath, kubectlBinary, 0755)
	if err != nil {
		return "", err
	}

	return kubectlPath, nil
}

func ExecKubectl(kubectlPath string, args []string) (out string, errStr string) {

	// Example: Run "kubectl version --client"
	cmd := exec.Command(kubectlPath, args...)

	//	执行命令并将结果放到 out 和 err 中
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	out = stdout.String()
	errStr = stderr.String()

	if err != nil {
		errStr = err.Error() + "\n" + errStr
		return out, errStr
	}

	return out, errStr

}
