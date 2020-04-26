package types

type SetProgressFunc func(name string, isStarted bool, isCompleted bool)

type ScannerStatus int

const (
	ScannerStatusPending   ScannerStatus = iota
	ScannerStatusRunning   ScannerStatus = iota
	ScannerStatusCompleted ScannerStatus = iota
)

type Progress struct {
	Steps      []string
	StepStatus map[string]ScannerStatus
}

func (p Progress) Get() ([]string, map[string]ScannerStatus) {
	// if p.Steps != nil {
	// 	panic(p)
	// }

	return p.Steps, p.StepStatus
}

func (p *Progress) Set(name string, isStarted bool, isCompleted bool) {
	status := ScannerStatusPending
	if isStarted {
		status = ScannerStatusRunning
	} else if isCompleted {
		status = ScannerStatusCompleted
	} else {
		status = ScannerStatusPending
	}

	if p.StepStatus == nil {
		p.StepStatus = map[string]ScannerStatus{}
	}
	if p.Steps == nil {
		p.Steps = []string{}
	}

	found := false
	for _, step := range p.Steps {
		if step == name {
			found = true
		}
	}
	if !found {
		p.Steps = append(p.Steps, name)
	}

	p.StepStatus[name] = status
}
