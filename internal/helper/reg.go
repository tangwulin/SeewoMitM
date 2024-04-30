package helper

import (
	"golang.org/x/sys/windows/registry"
)

func WriteMitMPortToRegistry(port int) error {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Zeus\Rpc\SeewoProxyHttp`, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer key.Close()
	err = key.SetDWordValue("port", uint32(port))
	if err != nil {
		return err
	}
	return nil
}

func GetUpstreamPortFromRegistry() (int, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Zeus\Rpc\SeewoProxyHttp`, registry.READ)
	if err != nil {
		return 0, err
	}

	defer key.Close()

	port, _, err := key.GetIntegerValue("port")
	if err != nil {
		return 0, err
	}

	return int(port), nil

}

func GetSeewoServiceAssisantPath() (string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Zeus\Rpc`, registry.READ)
	if err != nil {
		return "", err
	}

	defer key.Close()

	path, _, err := key.GetStringValue("SeewoServiceAssistant.exe")
	if err != nil {
		return "", err
	}

	return path, nil
}

func GetSeewoCorePath() (string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Zeus\Rpc`, registry.READ)
	if err != nil {
		return "", err
	}

	defer key.Close()

	path, _, err := key.GetStringValue("SeewoCore.exe")
	if err != nil {
		return "", err
	}

	return path, nil
}
