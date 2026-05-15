package event

import "testing"

func TestMapEvent(t *testing.T) {
	tests := []struct {
		name string
		ev   HookEvent
		want Category
	}{
		{
			name: "SessionStart maps to session.start",
			ev:   HookEvent{EventName: "SessionStart"},
			want: SessionStart,
		},
		{
			name: "SessionStart with source=compact is suppressed",
			ev:   HookEvent{EventName: "SessionStart", Source: "compact"},
			want: None,
		},
		{
			name: "Stop maps to task.complete",
			ev:   HookEvent{EventName: "Stop"},
			want: TaskComplete,
		},
		{
			name: "PermissionRequest maps to input.required",
			ev:   HookEvent{EventName: "PermissionRequest"},
			want: InputRequired,
		},
		{
			name: "Notification with idle_prompt maps to task.complete",
			ev:   HookEvent{EventName: "Notification", NotificationType: "idle_prompt"},
			want: TaskComplete,
		},
		{
			name: "Notification with elicitation_dialog maps to input.required",
			ev:   HookEvent{EventName: "Notification", NotificationType: "elicitation_dialog"},
			want: InputRequired,
		},
		{
			name: "Notification with permission_prompt is suppressed",
			ev:   HookEvent{EventName: "Notification", NotificationType: "permission_prompt"},
			want: None,
		},
		{
			name: "Notification with unknown type is suppressed",
			ev:   HookEvent{EventName: "Notification", NotificationType: "something_else"},
			want: None,
		},
		{
			name: "UserPromptSubmit maps to task.acknowledge",
			ev:   HookEvent{EventName: "UserPromptSubmit"},
			want: TaskAcknowledge,
		},
		{
			name: "PostToolUseFailure with Bash and error maps to task.error",
			ev:   HookEvent{EventName: "PostToolUseFailure", ToolName: "Bash", Error: "exit code 1"},
			want: TaskError,
		},
		{
			name: "PostToolUseFailure without Bash tool is suppressed",
			ev:   HookEvent{EventName: "PostToolUseFailure", ToolName: "Read", Error: "not found"},
			want: None,
		},
		{
			name: "PostToolUseFailure with Bash but no error is suppressed",
			ev:   HookEvent{EventName: "PostToolUseFailure", ToolName: "Bash"},
			want: None,
		},
		{
			name: "PreCompact maps to resource.limit",
			ev:   HookEvent{EventName: "PreCompact"},
			want: ResourceLimit,
		},
		{
			name: "SessionEnd is suppressed",
			ev:   HookEvent{EventName: "SessionEnd"},
			want: None,
		},
		{
			name: "PreToolUse is suppressed",
			ev:   HookEvent{EventName: "PreToolUse"},
			want: None,
		},
		{
			name: "PostToolUse is suppressed",
			ev:   HookEvent{EventName: "PostToolUse"},
			want: None,
		},
		{
			name: "SubagentStart is suppressed",
			ev:   HookEvent{EventName: "SubagentStart"},
			want: None,
		},
		{
			name: "SubagentStop is suppressed",
			ev:   HookEvent{EventName: "SubagentStop"},
			want: None,
		},
		{
			name: "unknown event is suppressed",
			ev:   HookEvent{EventName: "SomethingNew"},
			want: None,
		},
		{
			name: "empty event name is suppressed",
			ev:   HookEvent{},
			want: None,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapEvent(&tt.ev)
			if got != tt.want {
				t.Errorf("MapEvent() = %q, want %q", got, tt.want)
			}
		})
	}
}
