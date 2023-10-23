package glue

import (
	awsv1alpha1 "github.com/90poe/glue-jobs-operator/api/v1alpha1"
)

type GlueJob struct {
	name string
}

func NewGlueJob(job awsv1alpha1.GlueJobSpec) (*GlueJob, error) {
	ret := &GlueJob{
		name: job.Name,
	}
	return ret, nil
}
