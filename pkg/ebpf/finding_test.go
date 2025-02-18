package ebpf

import (
	"sort"
	"testing"

	"github.com/aquasecurity/tracee/pkg/events"
	"github.com/aquasecurity/tracee/types/detect"
	"github.com/aquasecurity/tracee/types/protocol"
	"github.com/aquasecurity/tracee/types/trace"
	"github.com/stretchr/testify/assert"
)

func TestFindingToEvent(t *testing.T) {
	expected := &trace.Event{
		EventID:             int(events.StartSignatureID),
		EventName:           "fake_signature_event",
		ProcessorID:         1,
		ProcessID:           2,
		CgroupID:            3,
		ThreadID:            4,
		ParentProcessID:     5,
		HostProcessID:       6,
		HostThreadID:        7,
		HostParentProcessID: 8,
		UserID:              9,
		MountNS:             10,
		PIDNS:               11,
		ProcessName:         "process",
		HostName:            "host",
		ContainerID:         "containerID",
		ContainerImage:      "image",
		ContainerName:       "container",
		PodName:             "pod",
		PodNamespace:        "namespace",
		PodUID:              "uid",
		ReturnValue:         0,
		MatchedScopes:       1,
		ArgsNum:             2,
		Args: []trace.Argument{
			{
				ArgMeta: trace.ArgMeta{
					Name: "arg1",
					Type: "const char *",
				},
				Value: "value1",
			},
			{
				ArgMeta: trace.ArgMeta{
					Name: "arg2",
					Type: "int",
				},
				Value: 1,
			},
		},
		Metadata: &trace.Metadata{
			Version:     "1",
			Description: "description",
			Tags:        []string{"tag1", "tag2"},
			Properties: map[string]interface{}{
				"prop1":         "value1",
				"prop2":         1,
				"signatureID":   "fake_signature_id",
				"signatureName": "fake_signature_event",
			},
		},
	}

	finding := createFakeEventAndFinding()
	got, err := FindingToEvent(finding)

	assert.NoError(t, err)

	// sort arguments to avoid flaky tests
	sort.Slice(got.Args, func(i, j int) bool { return got.Args[i].Name < got.Args[j].Name })
	sort.Slice(expected.Args, func(i, j int) bool { return expected.Args[i].Name < expected.Args[j].Name })

	assert.Equal(t, got, expected)
}

func createFakeEventAndFinding() detect.Finding {
	eventName := "fake_signature_event"
	event := events.NewEventDefinition(eventName, []string{"signatures"}, []events.ID{events.Ptrace})

	events.Definitions.Add(events.StartSignatureID, event)

	return detect.Finding{
		SigMetadata: detect.SignatureMetadata{
			ID:          "fake_signature_id",
			Name:        eventName,
			EventName:   eventName,
			Version:     "1",
			Description: "description",
			Tags:        []string{"tag1", "tag2"},
			Properties: map[string]interface{}{
				"prop1": "value1",
				"prop2": 1,
			},
		},
		Data: map[string]interface{}{
			"arg1": "value1",
			"arg2": 1,
		},
		Event: protocol.Event{
			Headers: protocol.EventHeaders{},
			Payload: trace.Event{
				EventID:             int(events.Ptrace),
				EventName:           "ptrace",
				ProcessorID:         1,
				ProcessID:           2,
				CgroupID:            3,
				ThreadID:            4,
				ParentProcessID:     5,
				HostProcessID:       6,
				HostThreadID:        7,
				HostParentProcessID: 8,
				UserID:              9,
				MountNS:             10,
				PIDNS:               11,
				ProcessName:         "process",
				HostName:            "host",
				ContainerID:         "containerID",
				ContainerImage:      "image",
				ContainerName:       "container",
				PodName:             "pod",
				PodNamespace:        "namespace",
				PodUID:              "uid",
				ReturnValue:         0,
				MatchedScopes:       1,
				ArgsNum:             0,
			},
		},
	}
}
