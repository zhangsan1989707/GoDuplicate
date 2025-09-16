package core

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ScanPreset struct {
	Name   string
	Config ScanConfig
}

type PolicyPreset struct {
	Name   string
	Policy Policy
}

func presetsDir() string {
	return filepath.Join(os.TempDir(), "haste_presets")
}

func SaveScanPreset(p ScanPreset) (string, error) {
	if err := os.MkdirAll(presetsDir(), 0o755); err != nil {
		return "", err
	}
	path := filepath.Join(presetsDir(), p.Name+"_scan.json")
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(p); err != nil {
		return "", err
	}
	return path, nil
}

func LoadScanPreset(name string) (ScanPreset, error) {
	path := filepath.Join(presetsDir(), name+"_scan.json")
	f, err := os.Open(path)
	if err != nil {
		return ScanPreset{}, err
	}
	defer f.Close()
	var p ScanPreset
	dec := json.NewDecoder(f)
	err = dec.Decode(&p)
	return p, err
}

func SavePolicyPreset(p PolicyPreset) (string, error) {
	if err := os.MkdirAll(presetsDir(), 0o755); err != nil {
		return "", err
	}
	path := filepath.Join(presetsDir(), p.Name+"_policy.json")
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(p); err != nil {
		return "", err
	}
	return path, nil
}

func LoadPolicyPreset(name string) (PolicyPreset, error) {
	path := filepath.Join(presetsDir(), name+"_policy.json")
	f, err := os.Open(path)
	if err != nil {
		return PolicyPreset{}, err
	}
	defer f.Close()
	var p PolicyPreset
	dec := json.NewDecoder(f)
	err = dec.Decode(&p)
	return p, err
}

// ListScanPresets returns available scan preset names.
func ListScanPresets() ([]string, error) {
	dir := presetsDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if len(name) > len("_scan.json") && name[len(name)-len("_scan.json"):] == "_scan.json" {
			names = append(names, name[:len(name)-len("_scan.json")])
		}
	}
	return names, nil
}

// DeleteScanPreset removes a scan preset file by name.
func DeleteScanPreset(name string) error {
	path := filepath.Join(presetsDir(), name+"_scan.json")
	return os.Remove(path)
}
