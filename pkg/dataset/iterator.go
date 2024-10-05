package dataset

type IterationMethod string

const (
	LineMethod  IterationMethod = "line"
	SplitMethod IterationMethod = "split"
	WholeMethod IterationMethod = "whole"
)
