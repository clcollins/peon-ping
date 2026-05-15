package event

func MapEvent(ev *HookEvent) Category {
	switch ev.EventName {
	case "SessionStart":
		if ev.Source == "compact" {
			return None
		}
		return SessionStart

	case "UserPromptSubmit":
		return TaskAcknowledge

	case "Stop":
		return TaskComplete

	case "Notification":
		switch ev.NotificationType {
		case "idle_prompt":
			return TaskComplete
		case "elicitation_dialog":
			return InputRequired
		default:
			return None
		}

	case "PermissionRequest":
		return InputRequired

	case "PostToolUseFailure":
		if ev.ToolName == "Bash" && ev.Error != "" {
			return TaskError
		}
		return None

	case "PreCompact":
		return ResourceLimit

	default:
		return None
	}
}
