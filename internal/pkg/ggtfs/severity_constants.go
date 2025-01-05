package ggtfs

type ValidationNoticeSeverity int

const (
	SeverityInfo           ValidationNoticeSeverity = 1
	SeverityRecommendation ValidationNoticeSeverity = 2
	SeverityViolation      ValidationNoticeSeverity = 3
)
