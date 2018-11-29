package appManagers

import "os/exec"

func runExe(exepath string, args1 string, args2 string, args3 string) error {
	cmd := exec.Command(exepath, args1, args2, args3)
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}
