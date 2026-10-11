package domain

type Effect string

const (
	EffectAllow Effect = "ALLOW"
	EffectDeny  Effect = "DENY"
)

func ToEffect(dbEffect string) Effect {
	switch Effect(dbEffect) {
	case EffectAllow:
		return EffectAllow
	case EffectDeny:
		return EffectDeny
	default:
		return ""
	}
}
